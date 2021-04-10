package kubectl

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/cjlapao/common-go/commands"
	"github.com/cjlapao/common-go/helper"
	"github.com/cjlapao/common-go/log"
	"github.com/cjlapao/deployment-tools-go/common"
	"gopkg.in/yaml.v2"
)

var logger = log.Get()

// Module Entity
type Module struct {
	Exists         bool
	ClientVersion  string
	ServerVersion  string
	Version        Version
	Exec           string
	Path           string
	CurrentContext string
	KubeConfig     KubeConfig
}

// NewModule creates a new Istio Module
func NewModule() Module {
	module := Module{
		Exists: false,
	}
	module.SetExists()

	if !module.Exists {
		module.GetModule()
	}

	// module.GetKubeConfig()

	return module
}

// SetExists checks if Istio module exists
func (m *Module) SetExists() bool {
	logger.Notice("Checking for installed kubectl version...")
	var kubectlModuleFolder string
	toolsFolder := common.ToolsFolder()

	files, err := ioutil.ReadDir(toolsFolder)

	if err != nil {
		return false
	}

	for _, f := range files {
		matched, _ := regexp.Match("^kubectl", []byte(f.Name()))
		logger.Debug("Module Item:", f.Name())
		if matched && f.IsDir() {
			kubectlModuleFolder = f.Name()
			logger.Debug("Matched On:", f.Name())
		}
	}

	if len(kubectlModuleFolder) > 0 {
		m.Path = helper.JoinPath(toolsFolder, kubectlModuleFolder)

		switch helper.GetOperatingSystem() {
		case helper.WindowsOs:
			m.Exec = helper.JoinPath(m.Path, "kubectl.exe")
		case helper.LinuxOs:
			m.Exec = helper.JoinPath(m.Path, "kubectl")
		}

		if helper.FileExists(m.Exec) {
			err := m.GetVersion()
			if err == nil {
				logger.Success("Found kubectl Module version", m.ClientVersion)
				m.Exists = true
				return true
			}
		}
	}

	logger.Notice("Kubectl Module was not found under tools")
	return false
}

// GetVersion gets Istio latest version from github
func (m *Module) GetVersion() error {
	if len(m.Exec) > 0 {
		logger.Debug("Found Kubectl executable, getting the version")
		out, err := commands.Execute(m.Exec, "version", "-o", "json", "--client=true")

		if err != nil && !strings.ContainsAny(err.Error(), "Unable to connect to the server: dial tcp [::1]:8080: connectex: No connection could be made because the target machine actively refused it.") {
			logger.Error(err.Error())
			return err
		}

		logger.Debug("Kubectl version output")
		logger.Debug(out)
		version := Version{}
		jsonError := json.Unmarshal([]byte(out), &version)
		if jsonError != nil {
			logger.Error(jsonError.Error())
			return jsonError
		}

		m.Version = version
		m.ClientVersion = version.ClientVersion.GitVersion
		m.ServerVersion = version.ServerVersion.GitVersion

		return nil
	}

	return errors.New("Executable was not found, cannot get a version")
}

// GetModule GetModule Istio into the tools folder
func (m *Module) GetModule() {
	logger.Notice("Checking online for latest version")
	latestVersion := getLatestVersion()
	toolsFolder := common.ToolsFolder()

	if m.Exists && getVersionFromString(m.ClientVersion) < getVersionFromString(latestVersion) {
		logger.Warn("There is a new version (" + latestVersion + ") online...")
	}

	if !m.Exists {
		if len(latestVersion) > 0 {
			filename, err := downloadModule(latestVersion)
			if err == nil {
				logger.Notice("Installing Kubectl", latestVersion)

				switch helper.GetOperatingSystem() {
				case helper.WindowsOs:
					m.Path = helper.JoinPath(toolsFolder, "kubectl-"+latestVersion, "bin")
					m.Exec = helper.JoinPath(m.Path, "kubectl.exe")
				case helper.LinuxOs:
					commands.Execute("chmod", "+x", filename)
					m.Path = helper.JoinPath(toolsFolder, "kubectl-"+latestVersion, "bin")
					m.Exec = helper.JoinPath(m.Path, "kubectl")
				}
			}

			m.SetExists()

			commands.Execute("export", "PATH="+m.Path+"/bin:$PATH")
		}
	} else {
		logger.Notice("Found version", m.ClientVersion, "on", toolsFolder)
	}
}

// GetKubeConfig gets Kubectl config from system
func (m *Module) GetKubeConfig() *KubeConfig {
	logger.Notice("Getting Kubectl config from system")
	out, err := commands.Execute(m.Exec, "config", "view", "--raw")

	if err != nil {
		logger.Error(err.Error())
		return nil
	}

	if err == nil && len(out) > 0 {

		kubeConfig, err := parseKubeConfig(out)
		if err != nil {
			logger.LogError(err)
		}

		m.KubeConfig = kubeConfig
		m.CurrentContext = kubeConfig.CurrentContext

		kubeconfigBytes, err := yaml.Marshal(kubeConfig)
		if err != nil {
			logger.LogError(err)
		}

		helper.WriteToFile(string(kubeconfigBytes), common.KubeConfig())

		return &m.KubeConfig
	}
	return nil
}

// AddKubeConfig adds a kubeconfig string into the Kubectl config
func (m *Module) AddKubeConfig(kubeConfig string) error {
	if len(kubeConfig) == 0 {
		return errors.New("The kubeconfig was null or empty")
	}

	logger.Notice("Adding kube config to Kubectl")
	config, err := parseKubeConfig(kubeConfig)
	if err != nil {
		logger.LogError(err)
		return err
	}

	// Adding the cluster
	for _, cluster := range config.Clusters {
		err := m.AddCluster(cluster)

		if err != nil {
			return err
		}
	}

	// Adding the users
	for _, user := range config.Users {
		err := m.AddCredentials(user)

		if err != nil {
			return err
		}
	}

	// Adding the contexts
	for _, context := range config.Contexts {
		err := m.AddContext(context)
		if err != nil {
			return err
		}
	}

	// Adding the current context
	if len(config.CurrentContext) > 0 {
		err := m.SetContext(config.CurrentContext)
		if err != nil {
			return err
		}
	}

	m.GetKubeConfig()

	return nil
}

// AddCluster Adds a cluster to the kubectl config
func (m *Module) AddCluster(cluster ClusterElement) error {
	var cmdErr error
	var out string
	logger.Debug("AddCluster")

	certTempFile := helper.JoinPath(common.TempFolder(), "ca.cert")
	if len(cluster.Cluster.CertificateAuthorityData) > 0 {
		bytes, err := base64.StdEncoding.DecodeString(cluster.Cluster.CertificateAuthorityData)

		if err != nil {
			return err
		}

		helper.WriteToFile(string(bytes), certTempFile)
		cluster.Cluster.CertificateAuthority = certTempFile
		out, cmdErr = commands.Execute(m.Exec, "config", "set-cluster", cluster.Name, "--server="+cluster.Cluster.Server, "--certificate-authority="+cluster.Cluster.CertificateAuthority, "--embed-certs=true")
	} else {
		out, cmdErr = commands.Execute(m.Exec, "config", "set-cluster", cluster.Name, "--server="+cluster.Cluster.Server, "--certificate-authority="+cluster.Cluster.CertificateAuthority)
	}
	logger.Debug(out)

	if cmdErr != nil {
		logger.LogError(cmdErr)
		return cmdErr
	}

	helper.DeleteFile(certTempFile)
	return nil
}

// AddCredentials Adds a credential to the kubectl config
func (m *Module) AddCredentials(user UserElement) error {
	var cmdErr error
	var out string
	logger.Debug("AddCredentials")
	if len(user.User.Token) > 0 {
		out, cmdErr = commands.Execute(m.Exec, "config", "set-credentials", user.Name, "--token="+user.User.Token)

	} else if len(user.User.Username) > 0 {
		out, cmdErr = commands.Execute(m.Exec, "config", "set-credentials", user.Name, "--username="+user.User.Username, "--password="+user.User.Password)
	} else {
		out, cmdErr = commands.Execute(m.Exec, "config", "set-credentials", user.Name, "--client-certificate="+user.User.ClientCertificate, "--client-key="+user.User.ClientKey)
	}
	logger.Debug(out)

	if cmdErr != nil {
		logger.LogError(cmdErr)
		return cmdErr
	}

	return nil
}

// AddContext Adds a context to the kubectl config
func (m *Module) AddContext(context ContextElement) error {
	logger.Debug("AddContext")
	out, cmdErr := commands.Execute(m.Exec, "config", "set-context", context.Name, "--cluster="+context.Context.Cluster, "--user="+context.Context.User, "--namespace="+context.Context.Namespace)
	logger.Debug(out)

	if cmdErr != nil {
		logger.LogError(cmdErr)
		return cmdErr
	}

	return nil
}

// SetContext Sets the default context for kubectl
func (m *Module) SetContext(contextName string) error {
	logger.Debug("SetContext")
	if len(contextName) == 0 {
		return errors.New("Context name is empty")
	}

	// Adding the current context
	out, cmdErr := commands.Execute(m.Exec, "config", "use-context", contextName)
	logger.Debug(out)

	if cmdErr != nil {
		logger.LogError(cmdErr)
		return cmdErr
	}

	m.CurrentContext = contextName

	return nil
}

// Delete Deletes a manifest from the default context
func (m *Module) Delete(manifestPath string) error {
	logger.Notice("Deleting " + manifestPath + " from " + m.CurrentContext)

	if len(manifestPath) == 0 {
		return errors.New("Path is empty")
	}

	if !helper.FileExists(manifestPath) {
		return errors.New("Manifes file " + manifestPath + "Does not exists in system")
	}

	out, cmdErr := commands.Execute(m.Exec, "delete", "--ignore-not-found=true", "-f", manifestPath, "--kubeconfig", common.KubeConfig(), "--context", m.CurrentContext)

	if cmdErr != nil {
		return cmdErr
	}

	logger.Debug(out)

	logger.Success("Manifest " + manifestPath + " deleted sucessfully")

	return nil
}

// DeleteNamespace Deletes a manifest from the default context
func (m *Module) DeleteNamespace(namespace string) error {
	logger.Notice("Deleting Namespace " + namespace + " from " + m.CurrentContext)

	if len(namespace) == 0 {
		return errors.New("Namespace is empty")
	}

	out, cmdErr := commands.Execute(m.Exec, "delete", "namespace", namespace, "--kubeconfig", common.KubeConfig(), "--context", m.CurrentContext)

	if cmdErr != nil {
		return cmdErr
	}

	logger.Debug(out)

	logger.Success("Namespace " + namespace + " deleted successfully")

	return nil
}

// GetAllNamespace Deletes a manifest from the default context
func (m *Module) GetAllNamespace() error {
	logger.Notice("Getting all namespaces from " + m.CurrentContext)

	out, cmdErr := commands.Execute(m.Exec, "get", "namespace", "--kubeconfig", common.KubeConfig(), "--context", m.CurrentContext)

	if cmdErr != nil {
		return cmdErr
	}

	logger.Debug(out)

	logger.Success("Got all namespaces successfully")

	return nil
}

func getLatestVersion() string {
	baseURL := "https://storage.googleapis.com/kubernetes-release/release/stable.txt"

	logger.Debug("Listing all available versions of kubectl to download")
	response, err := http.Get(baseURL)

	helper.CheckError(err)

	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(response.Body)

		helper.CheckError(err)

		latestVersion := string(bodyBytes)
		logger.Debug("Version:", latestVersion)

		return latestVersion
	}

	return ""
}

func downloadModule(version string) (string, error) {
	var filename string
	toolsFolder := common.ToolsFolder()
	foldername := helper.JoinPath(toolsFolder, "kubectl-"+strings.ReplaceAll(version, "v", ""))
	switch helper.GetOperatingSystem() {
	case helper.WindowsOs:
		filename = "kubectl.exe"
	case helper.LinuxOs:
		filename = "kubectl"
	}
	url := "https://storage.googleapis.com/kubernetes-release/release/" + version + "/bin/" + strings.ToLower(helper.GetOperatingSystem().String()) + "/amd64/" + filename

	if _, err := os.Stat(toolsFolder); os.IsNotExist(err) {
		logger.Debug("Tools folder doesn't exist, creating")
		os.Mkdir(toolsFolder, os.ModeDir)
	}

	if _, err := os.Stat(foldername); os.IsNotExist(err) {
		logger.Debug("Kubectl folder doesn't exist, creating")
		os.Mkdir(foldername, os.ModeDir)
	}

	downloadPath := helper.JoinPath(foldername, filename)
	logger.Debug("Downloading file into", downloadPath)
	err := helper.DownloadFile(url, downloadPath)

	logger.LogError(err)
	return filename, err
}

func getVersionFromString(value string) int {
	v, err := strconv.Atoi(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(value, "kubectl", ""), "-", ""), ".", ""))

	if err != nil {
		return -1
	}

	logger.Debug("Istio int version: ", fmt.Sprint(v))
	return v
}

func parseKubeConfig(value string) (KubeConfig, error) {
	var kubeConfig KubeConfig

	yamlErr := yaml.Unmarshal([]byte(value), &kubeConfig)
	if yamlErr != nil {
		logger.LogError(yamlErr)
		return KubeConfig{}, yamlErr
	}

	return kubeConfig, nil
}

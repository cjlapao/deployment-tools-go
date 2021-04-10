package istio

import (
	"errors"
	"strings"

	"github.com/cjlapao/common-go/commands"
	"github.com/cjlapao/common-go/helper"
	"github.com/cjlapao/common-go/log"
	"github.com/cjlapao/deployment-tools-go/executioncontext"
	"github.com/cjlapao/deployment-tools-go/kubectl"
	"github.com/cjlapao/deployment-tools-go/module"
	"github.com/cjlapao/deployment-tools-go/startup"
)

var ctx = executioncontext.Get()

// Module entity
type Module struct {
	Exists        bool
	ClientVersion string
	Version       Version
	Exec          string
	BinPath       string
	FolderPath    string
	Profile       Profile
	Dependencies  Dependencies
	Module        module.DeploymentToolsModule
}

// Dependencies entity
type Dependencies struct {
	Kubectl kubectl.Module
}

// Profile Enum
type Profile int

var logger = log.Get()
var globalIstioModule *Module

// IstioProfile Enum definition
const (
	DefaultProfile Profile = iota
	DemoProfile
	EmptyProfile
	MinimalProfile
	OpenshiftProfile
	PreviewProfile
	RemoteProfile
)

// IstioProfile ToString convertion
func (o Profile) String() string {
	return [...]string{"default", "demo", "empty", "minimal", "openshift", "preview", "remote"}[o]
}

// ProfileFromString convertion
func ProfileFromString(value string) Profile {
	switch strings.ToLower(value) {
	case "demo":
		return DemoProfile
	case "empty":
		return EmptyProfile
	case "minimal":
		return MinimalProfile
	case "openshift":
		return OpenshiftProfile
	case "preview":
		return PreviewProfile
	case "remote":
		return RemoteProfile
	default:
		return DefaultProfile
	}
}

// Create creates a new Istio Module
func Create() *Module {
	if globalIstioModule != nil {
		return globalIstioModule
	}

	module := Module{
		Exists:  false,
		Profile: DefaultProfile,
		Module: module.DeploymentToolsModule{
			Name:                  "Istio",
			WindowsExecutableName: "istioctl.exe",
			LinuxExecutableName:   "istioctl",
			ExecPath:              "bin",
			Download:              Download,
			GetLocalVersion:       GetLocalVersion,
			GetOnlineVersions:     GetOnlineVersions,
			Process:               Processor,
			PostInstall:           PostInstall,
		},
	}

	// module.SetDependencies()
	module.Module.GetLatestCachedVersion()

	globalIstioModule = &module
	return globalIstioModule
}

// SetDependencies Sets Module dependencies
func (m *Module) SetDependencies() {
	logger.Command("Checking Istio module dependencies")
	m.Dependencies.Kubectl = kubectl.NewModule()
}

// Remove removes Istio from cluster
// This is using the uninstall process stated in the istio homepage
func (m *Module) Remove(context string) error {
	if len(context) == 0 {
		return errors.New("Context is null or empty")
	}

	if m.Dependencies.Kubectl.CurrentContext != context {
		err := m.Dependencies.Kubectl.SetContext(context)

		if err != nil {
			return err
		}
	}

	// Deleting the addons used during install
	addonsManifest := helper.JoinPath(m.FolderPath, "samples/addons")

	err := m.Dependencies.Kubectl.Delete(addonsManifest)

	if err != nil {
		return err
	}

	// Generating the manifest used for installation
	profileTempManifest := helper.JoinPath(m.FolderPath, m.Profile.String()+"-temp.yaml")

	out, err := commands.Execute(m.Exec, "manifest", "generate", "--set", "profile="+m.Profile.String(), "--kubeconfig "+startup.KubeConfig())
	if err != nil {
		return err
	}

	err = helper.WriteToFile(out, profileTempManifest)
	if err != nil {
		return err
	}

	// Deleting the istio using kubectl
	err = m.Dependencies.Kubectl.Delete(profileTempManifest)

	if err != nil {
		return err
	}

	// Deleting the namespace
	err = m.Dependencies.Kubectl.DeleteNamespace("istio-system")

	if err != nil {
		return err
	}

	return nil
}

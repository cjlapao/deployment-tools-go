package module

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/cjlapao/common-go/commands"
	"github.com/cjlapao/common-go/helper"
	"github.com/cjlapao/common-go/log"
	"github.com/cjlapao/common-go/version"
	"github.com/cjlapao/deployment-tools-go/common"
	"github.com/cjlapao/deployment-tools-go/executioncontext"
)

var logger = log.Get()
var versionSvc = version.Get()
var ctx = executioncontext.Get()

// DeploymentToolsModule Entity
type DeploymentToolsModule struct {
	Name                  string
	WindowsExecutableName string
	LinuxExecutableName   string
	FolderPath            string
	ModulePath            string
	ExecPath              string
	Exec                  string
	Exists                bool
	Latest                bool
	Version               string
	GetLocalVersion       GetLocalVersion
	GetOnlineVersions     GetOnlineVersions
	Download              Download
	Process               Process
	PreInstall            PreInstall
	PostInstall           PostInstall
	PreUninstall          PreUninstall
	PostUninstall         PostUninstall
}

// Version Entity
type Version struct {
	Version  string `json:"version" yaml:"version"`
	Revision string `json:"revision" yaml:"revision"`
	Path     string `json:"path" yaml:"path"`
}

// GetLocalVersion Local Module Version function
type GetLocalVersion func(filePath string) (string, error)

// GetOnlineVersions Local Module Get Online Versions function
type GetOnlineVersions func() ([]string, error)

// Download Local Module Download function
type Download func(version string) (string, error)

// Process Process the module variables and functions
type Process func()

// PreInstall Instructions to be executed before installing any version
type PreInstall func(version string) error

// PostInstall Instructions to be executed after installing any version
type PostInstall func(version string) error

// PreUninstall Instructions to be executed before uninstalling any version
type PreUninstall func(version string) error

// PostUninstall Instructions to be executed after uninstalling any version
type PostUninstall func(version string) error

// SetExists checks if Terraform module exists
func (m *DeploymentToolsModule) SetExists(version string) bool {
	logger.Notice("Checking for installed %v version...", m.Name)

	chachedVersions := m.GetCachedVersions()
	if len(chachedVersions) > 0 {
		var setVersion Version
		if version == "" {
			setVersion = chachedVersions[len(chachedVersions)-1]
		} else {
			for _, v := range chachedVersions {
				if m.VersionToInt(v.Version) == m.VersionToInt(version) {
					setVersion = v
					break
				}
			}
		}

		m.SetPathsFromVersion(setVersion)

		if setVersion.Version == "" {
			logger.Debug("Version %v was not found in the cache system")
			return false
		}

		if helper.FileExists(m.Exec) {
			version, err := m.GetLocalVersion(m.Exec)
			m.Version = version
			logger.LogError(err)
			if err == nil {
				logger.Success("Found %v Module version %v", m.Name, m.Version)
				m.Exists = true
				m.Latest = m.IsLatest()
				return true
			}
		}
	}

	logger.Notice("%v Module was not found under tools", m.Name)
	return false
}

// SetWorkingVersion Sets a working version for the module
func (m *DeploymentToolsModule) SetWorkingVersion(version string) {
	logger.Notice("Setting %v version to %v", m.Name, version)
	cachedVersions := m.GetCachedVersions()
	foundCachedVersion := false
	for _, v := range cachedVersions {
		if m.VersionToInt(v.Version) == m.VersionToInt(version) {
			logger.Debug("Found version %v cached in the system, using it", version)
			foundCachedVersion = true
			m.SetPathsFromVersion(v)
			break
		}
	}
	if !foundCachedVersion {
		m.Exists = false
		m.InstallVersion(version)
	}
}

// GetCachedVersions Get all the cached version of a module
func (m *DeploymentToolsModule) GetCachedVersions() []Version {
	installedVersions := make([]Version, 0)
	moduleName := strings.ToLower(m.Name)
	toolsFolder := common.ToolsFolder()

	files, err := ioutil.ReadDir(toolsFolder)

	if err != nil {
		return installedVersions
	}

	for _, f := range files {
		matched, _ := regexp.Match("^"+moduleName, []byte(f.Name()))
		logger.Debug("Module Item: %v", f.Name())
		if matched && f.IsDir() {
			installedVersion := Version{
				Version: strings.ReplaceAll(f.Name(), moduleName+"-", ""),
				Path:    helper.JoinPath(toolsFolder, f.Name()),
			}
			installedVersions = append(installedVersions, installedVersion)
			logger.Debug("Matched On: %v", f.Name())
		}
	}

	return installedVersions
}

// GetLatestCachedVersion Gets the latest cached version of the module
func (m *DeploymentToolsModule) GetLatestCachedVersion() Version {
	logger.Debug("Getting latest cached version for %v", m.Name)
	cachedVersions := m.GetCachedVersions()
	if len(cachedVersions) > 0 {
		version := cachedVersions[len(cachedVersions)-1]
		m.SetPathsFromVersion(version)
		m.Exists = true
		logger.Debug("Found local cached version %v for %v", version.Version, m.Name)
		return version
	}
	return Version{}
}

// UninstallAll Removes all cached versions of the module
func (m *DeploymentToolsModule) UninstallAll() error {
	logger.Notice("Starting %v Uninstall", m.Name)
	toolsFolder := common.ToolsFolder()
	cachedVersions := m.GetCachedVersions()
	if len(cachedVersions) > 0 {
		for _, cVer := range cachedVersions {
			err := os.RemoveAll(helper.JoinPath(toolsFolder, cVer.Path))
			if err != nil {
				logger.LogError(err)
				return err
			}
		}
	}

	logger.Success("Successfully remove all %v Modules", m.Name)
	return nil
}

// UninstallVersion Removes a specific version of a module
func (m *DeploymentToolsModule) UninstallVersion(version string) error {
	logger.Notice("Starting %v version %v Uninstall", m.Name, version)
	cachedVersions := m.GetCachedVersions()
	var versionFound Version
	cached := false
	if len(cachedVersions) > 0 {
		for _, cVer := range cachedVersions {
			if m.VersionToInt(cVer.Version) == m.VersionToInt(version) {
				cached = true
				versionFound = cVer
				break
			}
		}
	}

	if cached {
		if m.PreUninstall != nil {
			logger.Notice("Starting Pre-Uninstall for %v version %v", m.Name, version)
			err := m.PreUninstall(versionFound.Path)
			if err != nil {
				return err
			}
			logger.Success("Finished Pre-Uninstall for %v version %v", m.Name, version)
		}
		err := os.RemoveAll(versionFound.Path)
		if err != nil {
			logger.LogError(err)
			return err
		}
		if m.PostUninstall != nil {
			logger.Notice("Starting Post-Uninstall for %v version %v", m.Name, version)
			err := m.PostUninstall(versionFound.Path)
			if err != nil {
				return err
			}
			logger.Success("Finished Post-Uninstall for %v version %v", m.Name, version)
		}

		logger.Success("Successfully remove Module %v version %v", m.Name, version)
	} else {
		logger.Notice("Could not find Module %v version %v in the cached folder", m.Name, version)
	}
	return nil
}

// InstallVersion Installs a specific version of a module
func (m *DeploymentToolsModule) InstallVersion(version string) error {
	toolsFolder := common.ToolsFolder()
	moduleName := strings.ToLower(m.Name)
	logger.Notice("Checking for installed %v version %v...", m.Name, version)
	cachedVersions := m.GetCachedVersions()
	cached := false
	if len(cachedVersions) > 0 {
		for _, cVer := range cachedVersions {
			if m.VersionToInt(cVer.Version) == m.VersionToInt(version) {
				cached = true
				break
			}
		}
	}

	if cached {
		logger.Success("Found cached version %v of %v module", version, m.Name)
	} else {
		versions, err := m.GetOnlineVersions()
		if err != nil {
			logger.Error("There was an error getting the latest versions online")
			return err
		}

		for _, v := range versions {
			if m.VersionToInt(v) == m.VersionToInt(version) {
				if m.PreInstall != nil {
					logger.Notice("Starting Pre-install for %v version %v", m.Name, version)
					err = m.PreInstall(version)
					if err != nil {
						return err
					}
					logger.Success("Finished Pre-install for %v version %v", m.Name, version)
				}
				if len(version) > 0 {
					filename, err := m.Download(version)
					if err == nil {
						logger.Notice("Installing %v %v", m.Name, version)
						if _, err := os.Stat(toolsFolder); os.IsNotExist(err) {
							logger.Debug("Tools folder doesn't exist, creating")
							os.Mkdir(toolsFolder, os.ModeDir)
						}

						switch helper.GetOperatingSystem() {
						case helper.WindowsOs:
							logger.Debug("Extracting %v to the tools folder", m.Name)
							_, err := helper.Unzip(helper.JoinPath(helper.GetExecutionPath(), filename), toolsFolder+"/"+moduleName+"-"+version)
							if err != nil {
								logger.LogError(err)
							}
							os.Remove(helper.JoinPath(helper.GetExecutionPath(), filename))
							m.ModulePath = helper.JoinPath(toolsFolder, moduleName+"-"+version)
							m.Exec = helper.JoinPath(m.ModulePath, m.ExecPath, m.WindowsExecutableName)
							commands.Execute("set", "PATH=%PATH%;"+m.ModulePath)
						case helper.LinuxOs:
							logger.Debug("Extracting %v to the tools folder", m.Name)
							_, err := helper.Unzip(helper.JoinPath(helper.GetExecutionPath(), filename), toolsFolder+"/"+moduleName+"-"+version)

							if err != nil {
								logger.LogError(err)
							}

							os.Remove(helper.JoinPath(helper.GetExecutionPath(), filename))
							m.ModulePath = helper.JoinPath(toolsFolder, moduleName+"-"+version)
							m.Exec = helper.JoinPath(m.ModulePath, m.ExecPath, m.LinuxExecutableName)
							commands.Execute("export", "PATH="+m.ModulePath+":$PATH")
						}
					}

					if m.PostInstall != nil {
						logger.Notice("Starting Post-install for %v version %v", m.Name, version)
						err = m.PostInstall(version)
						if err != nil {
							return err
						}
						logger.Success("Finished Post-install for %v version %v", m.Name, version)
					}
					m.SetExists(version)
				}
				logger.Debug("Module %v version is now %v", m.Name, v)
			}
		}
	}

	return nil
}

// GetLatestModule Get the latest version of a specific module
func (m *DeploymentToolsModule) GetLatestModule() {
	logger.Notice("Checking online for latest version of %v", m.Name)
	latestVersion := m.GetLatestVersion()
	cachedVersions := m.GetCachedVersions()
	foundLatestCached := false
	for _, v := range cachedVersions {
		if m.VersionToInt(v.Version) == m.VersionToInt(latestVersion) {
			logger.Debug("Found version %v cached in the system, using it", latestVersion)
			foundLatestCached = true
			m.SetPathsFromVersion(v)
			break
		}
	}
	if !foundLatestCached {
		m.Exists = false
		m.GetModule(latestVersion)
	}
}

// GetModule GetModule Terraform into the tools folder
func (m *DeploymentToolsModule) GetModule(version string) {
	toolsFolder := common.ToolsFolder()

	if m.Exists && m.VersionToInt(m.Version) < m.VersionToInt(version) {
		logger.Warn("There is a new version (%v) online...", version)
	}

	if !m.Exists {
		m.InstallVersion(version)
	} else {
		logger.Notice("Found version %v on %v", m.Version, toolsFolder)
	}
}

// GetLatestVersion Gets a module latest version
func (m *DeploymentToolsModule) GetLatestVersion() string {
	versions, err := m.GetOnlineVersions()
	if err != nil {
		return ""
	}

	sort.Strings(versions)
	for _, v := range versions {
		logger.Debug("Version: %v", v)
	}

	latest := versions[len(versions)-1]
	logger.Debug("Latest Version: %v", latest)

	return latest
}

// VersionToInt Gets a integer of a version for version comparing
func (m *DeploymentToolsModule) VersionToInt(value string) int {
	v, err := strconv.Atoi(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(value, m.Name, ""), "-", ""), ".", ""))

	if err != nil {
		return -1
	}

	logger.Trace("%v ToInt: %v", m.Name, fmt.Sprint(v))
	return v
}

// IsLatest Check if the installed version is the latest
func (m *DeploymentToolsModule) IsLatest() bool {
	logger.Debug("Checking for latest version")
	latestVersion := m.GetLatestVersion()
	logger.Debug("Found %v version online as the latest", latestVersion)
	logger.Debug("Found %v version being used locally", m.Version)
	if m.Exists && m.VersionToInt(m.Version) < m.VersionToInt(latestVersion) {
		logger.Warn("There is a new version (%v) online...", latestVersion)
		return false
	}
	logger.Info("Version %v is up to date", latestVersion)
	return true
}

// SetPathsFromVersion Sets the path od the module from a version
func (m *DeploymentToolsModule) SetPathsFromVersion(version Version) {
	m.ModulePath = version.Path
	m.Version = version.Version
	switch helper.GetOperatingSystem() {
	case helper.WindowsOs:
		m.Exec = helper.JoinPath(m.ModulePath, m.ExecPath, m.WindowsExecutableName)
	case helper.LinuxOs:
		m.Exec = helper.JoinPath(m.ModulePath, m.ExecPath, m.LinuxExecutableName)
	}
	logger.Debug("Module %v was set to version %v in path %v using os exec %v", m.Name, m.Version, m.ModulePath, m.Exec)
}

func (m *DeploymentToolsModule) GetExec() string {
	if m.Exec == "" {
		return ""
	}

	switch helper.GetOperatingSystem() {
	case helper.WindowsOs:
		return helper.JoinPath(m.ModulePath, m.ExecPath, m.WindowsExecutableName)
	case helper.LinuxOs:
		return helper.JoinPath(m.ModulePath, m.ExecPath, m.LinuxExecutableName)
	default:
		return ""
	}
}

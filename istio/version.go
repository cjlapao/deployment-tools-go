package istio

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/cjlapao/common-go/commands"
	"github.com/cjlapao/common-go/helper"
	"github.com/cjlapao/deployment-tools-go/common"
)

// Version entity
type Version struct {
	ClientVersion    ClientVersion      `json:"clientVersion"`
	MeshVersion      []MeshVersion      `json:"meshVersion"`
	DataPlaneVersion []DataPlaneVersion `json:"dataPlaneVersion"`
}

// ClientVersion entity
type ClientVersion struct {
	Version       string `json:"version"`
	Revision      string `json:"revision"`
	GolangVersion string `json:"golang_version"`
	Status        string `json:"status"`
	Tag           string `json:"tag"`
}

// DataPlaneVersion entity
type DataPlaneVersion struct {
	ID           string `json:"ID"`
	IstioVersion string `json:"IstioVersion"`
}

// MeshVersion entity
type MeshVersion struct {
	Component string        `json:"Component"`
	Info      ClientVersion `json:"Info"`
}

// GetLocalVersion gets Module latest version from github
func GetLocalVersion(filePath string) (string, error) {
	if len(filePath) > 0 {
		logger.Debug("Found Istio executable on path %v, getting the version", filePath)
		out, err := commands.Execute(filePath, "version", "--remote=false", "-o", "json")

		if err != nil {
			logger.Error(err.Error())
			return "", err
		}

		version := Version{}
		jsonError := json.Unmarshal([]byte(out), &version)
		if jsonError != nil {
			logger.Error(jsonError.Error())
			return "", jsonError
		}

		logger.Debug("Found Istio version %v", version.ClientVersion.Version)
		return version.ClientVersion.Version, nil
	}

	return "", errors.New("Executable was not found, cannot get a version")
}

// GetOnlineVersions Get the available online versions
func GetOnlineVersions() ([]string, error) {
	versions := make([]string, 0)
	baseURL := "https://github.com/istio/istio/releases"

	logger.Debug("Listing all available versions of isto to download")
	response, err := http.Get(baseURL)

	if err != nil {
		return versions, err
	}

	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(response.Body)

		if err != nil {
			return versions, err
		}

		bodyString := string(bodyBytes)

		re := regexp.MustCompile("releases/[0-9]*.[0-9]*.[0-9]*/")
		match := re.FindAllStringSubmatch(bodyString, 10)
		var versions []string
		for _, v := range match {
			versions = append(versions, strings.ReplaceAll(strings.ReplaceAll(v[0], "releases/", ""), "/", ""))
		}

		return versions, nil
	}

	return versions, nil
}

// Download Downloads a version
func Download(version string) (string, error) {
	var filename string
	switch helper.GetOperatingSystem() {
	case helper.WindowsOs:
		filename = "istio-" + version + "-win.zip"
	case helper.LinuxOs:
		filename = "istio-" + version + "-linux-arm64.tar.gz"
	}
	url := "https://github.com/istio/istio/releases/download/" + version + "/" + filename

	downloadPath := helper.JoinPath(helper.GetExecutionPath(), filename)
	logger.Debug("Downloading file into %v", downloadPath)
	err := helper.DownloadFile(url, downloadPath)

	if err != nil {
		logger.LogError(err)
		return "", err
	}

	return filename, nil
}

// PostInstall Istio post install instructions
func PostInstall(version string) error {
	logger.Notice("Executing post install for Istio %v", version)
	toolsFolder := common.ToolsFolder()
	istioDir := "istio-" + version
	baseDir := helper.JoinPath(toolsFolder, istioDir)
	istioDownloadDir := helper.JoinPath(baseDir, istioDir)
	logger.Debug("Starting to move Istio from %v to %v", istioDownloadDir, baseDir)
	err := helper.CopyDir(istioDownloadDir, baseDir)
	if err != nil {
		logger.LogError(err)
		return err
	}
	logger.Debug("Finished moving Istio to %v", baseDir)
	err = os.RemoveAll(istioDownloadDir)
	if err != nil {
		logger.LogError(err)
		return err
	}
	logger.Debug("Finished removing Istio %v directory", istioDownloadDir)

	return nil
}

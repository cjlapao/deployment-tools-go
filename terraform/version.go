package terraform

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"github.com/cjlapao/common-go/commands"
	"github.com/cjlapao/common-go/helper"

	"github.com/pkg/errors"
)

// Version entity
type Version struct {
	Version            string      `json:"terraform_version"`
	Revision           string      `json:"terraform_revision"`
	ProviderSelections interface{} `json:"provider_selections"`
	TerraformOutdated  bool        `json:"terraform_outdated"`
}

// GetLocalVersion gets Module latest version from github
func GetLocalVersion(filePath string) (string, error) {
	if len(filePath) > 0 {
		logger.Debug("Found Terraform executable on path %v, getting the version", filePath)
		out, err := commands.Execute(filePath, "-v", "-json")

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

		logger.Debug("Found Terraform version %v", version.Version)
		return version.Version, nil
	}

	return "", errors.New("Executable was not found, cannot get a version")
}

// GetOnlineVersions Get the available online versions
func GetOnlineVersions() ([]string, error) {
	versions := make([]string, 0)
	baseURL := "https://releases.hashicorp.com/terraform/"

	logger.Debug("Listing all available versions of kubectl to download")
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
		re := regexp.MustCompile("/terraform/[0-9]*.[0-9]*.[0-9]*/")
		match := re.FindAllStringSubmatch(bodyString, 20)
		for _, v := range match {
			versions = append(versions, strings.ReplaceAll(strings.ReplaceAll(v[0], "terraform/", ""), "/", ""))
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
		filename = "terraform_" + version + "_windows_amd64.zip"
	case helper.LinuxOs:
		filename = "terraform_" + version + "_linux-_arm64.zip"
	}

	url := "https://releases.hashicorp.com/terraform/" + version + "/" + filename

	downloadPath := helper.JoinPath(helper.GetExecutionPath(), filename)
	logger.Debug("Downloading file into %v", downloadPath)
	err := helper.DownloadFile(url, downloadPath)
	if err != nil {
		logger.Error(err)
		return "", err
	}

	return filename, nil
}

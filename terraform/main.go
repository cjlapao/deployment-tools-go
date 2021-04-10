package terraform

import (
	"github.com/cjlapao/common-go/fileproc"
	"github.com/cjlapao/common-go/helper"
	"github.com/cjlapao/common-go/log"
	"github.com/cjlapao/deployment-tools-go/module"
)

// Module Terraform Module
type Module struct {
	VariableFileName string
	Verbose          bool
	Exists           bool
	Version          Version
	Module           *module.DeploymentToolsModule
}

var globalTerraformModule *Module
var logger = log.Get()

// Create Creates a terraform module to be used in the deployment tools
func Create() *Module {
	if globalTerraformModule != nil {
		return globalTerraformModule
	}

	module := Module{
		Verbose: false,
		Exists:  false,
		Module: &module.DeploymentToolsModule{
			Name:                  "Terraform",
			WindowsExecutableName: "terraform.exe",
			LinuxExecutableName:   "terraform",
			GetOnlineVersions:     GetOnlineVersions,
			Download:              Download,
			GetLocalVersion:       GetLocalVersion,
		},
	}

	module.Module.GetLatestCachedVersion()

	globalTerraformModule = &module

	return globalTerraformModule
}

func (m *Module) ImportTest(filePath string, format ImportFormat, variables ...fileproc.Variable) {
	if helper.FileExists(filePath) {
		logger.Info("Found file %v, reading content", filePath)
		content, err := helper.ReadFromFile(filePath)
		if err != nil {
			logger.LogError(err)
			return
		}
		if len(variables) > 0 {
			content = fileproc.Process(content, variables...)
		}

		imp, err := m.ReadImportContent(content, format)
		if err != nil {
			logger.LogError(err)
		}

		logger.Debug("ID: %v", imp.Azure.SubscriptionID)
		for _, topic := range imp.Azure.ServiceBus.Topics {
			logger.Debug("Topic Name: %v", topic.Name)
		}
	}
}

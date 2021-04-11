package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/cjlapao/deployment-tools-go/azurecli"
	"github.com/cjlapao/deployment-tools-go/istio"
	"github.com/cjlapao/deployment-tools-go/module"
	"github.com/cjlapao/deployment-tools-go/service"
	"github.com/cjlapao/deployment-tools-go/startup"
	"github.com/cjlapao/deployment-tools-go/terraform"

	"github.com/cjlapao/common-go/helper"
	"github.com/cjlapao/common-go/log"
	"github.com/cjlapao/common-go/version"
)

var startupSvc = service.CreateProvider()
var logger = log.Get()
var versionSvc = version.Get()

func main() {
	getVersion := helper.GetFlagSwitch("version", false)
	if getVersion {
		format := helper.GetFlagValue("o", "json")
		switch strings.ToLower(format) {
		case "json":
			fmt.Println(versionSvc.PrintVersion(int(version.JSON)))
		case "yaml":
			fmt.Println(versionSvc.PrintVersion(int(version.JSON)))
		default:
			fmt.Println("Please choose a valid format, this can be either json or yaml")
		}
		os.Exit(0)
	}

	versionSvc.PrintAnsiHeader()

	startup.Start()

	defer func() {
		startup.Exit(0)
	}()

	cs := os.Getenv("mongoConnectionString")
	logger.Info("Connection string %v", cs)

	moduleName := module.GetModuleArgument()

	switch moduleName {
	case "sandbox":
		sandboxModule()
	case "api":
		module.RestApiModuleProcessor()
	case "azure":
		module.AzurecliModuleProcessor()
	case "istio":
		istioModule := istio.Create()
		istioModule.Module.Process()
	case "servicebus":
		module.ServiceBusCliModuleProcessor()
	case "terraform":
		terraformModule := terraform.Create()
		terraformModule.Module.Process()
	default:
		module.PrintCommandHelper()
	}

	startup.Exit(0)
}

func sandboxModule() {
	command := module.GetCommandArgument()
	subCommand := module.GetSubCommandArgument()
	g := azurecli.CreateStorageAccountClient()
	g.Test()
	logger.Info("command %v, subcommand %v", command, subCommand)
	if command == "storage" {
		switch subCommand {
		case "upload":
			azurecli.UploadBlob()
		case "download":
			azurecli.DownloadBlob()
		case "list":
			items, _ := azurecli.ListFilesInContainer("")
			for _, item := range items {
				fmt.Println(item.Name)
			}
		case "create-container":
			err := azurecli.CreateContainer()
			if err != nil {
				logger.Error(err)
			}
		case "delete-container":
			err := azurecli.DeleteContainer()
			if err != nil {
				logger.Error(err)
			}
		}
	}
	// cli := terraform.Create()
	// file := helper.GetFlagValue("file", "")
	// test := helper.GetFlagValue("test", "")
	// json := helper.GetFlagSwitch("json", false)
	// if test != "" {
	// 	logger.Debug("Found test flag value %v", test)
	// } else {
	// 	logger.Debug("Did not find any test flag")
	// }
	// // cli.Module.SetWorkingVersion("0.13.0")
	// logger.Debug("Terraform Module Version: %v", cli.Module.Version)
	// logger.Debug("Terraform Module Path: %v", cli.Module.ModulePath)
	// // cli.Module.UninstallVersion("0.13.0")
	// if file == "" {
	// 	logger.Error("File was not found, please use --file={{file}} to import")
	// }

	// variable1 := fileproc.Variable{
	// 	Name:  "terraformModule",
	// 	Value: "module.topics[test]",
	// }
	// variable2 := fileproc.Variable{
	// 	Name:  "azureSubscriptionId",
	// 	Value: "e7ebda2e-82c4-458f-8ce7-47677990870e",
	// }

	// if json {
	// 	cli.ImportTest(file, terraform.JSON, variable1, variable2)
	// } else {
	// 	cli.ImportTest(file, terraform.Yaml, variable1, variable2)
	// }
}

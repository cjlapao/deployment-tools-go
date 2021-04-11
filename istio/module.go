package istio

import (
	"fmt"

	"github.com/cjlapao/common-go/helper"
	"github.com/cjlapao/deployment-tools-go/kubectl"
	"github.com/cjlapao/deployment-tools-go/module"
	"github.com/cjlapao/deployment-tools-go/startup"
)

// Processor Process the module commands for istio
func Processor() {
	defer func() {
		startup.Exit(0)
	}()

	command := module.GetCommandArgument()
	if len(command) == 0 {
		istioModuleCommandHelper()
		startup.Exit(0)
	}

	switch command {
	case "version":
		istio := Create()
		logger.Debug("Version %v", istio.Module.Version)
	case "get-latest":
		istio := Create()
		istio.Module.GetLatestModule()
	case "set-version":
		istio := Create()
		version := helper.GetFlagValue("version", "")
		if version == "" {
			istioModuleCommandHelper()
			return
		}

		istio.Module.SetWorkingVersion(version)
		logger.Debug("Current module version: %v using %v", istio.Module.Version, istio.Module.Exec)
	case "remove":
		istioModule := Create()
		if ctx.ShowHelp {
			istioModuleRemoveSubCommandHelper()
			startup.Exit(0)
		}
		if len(ctx.Istio.Profile) > 0 {
			istioModule.Profile = ProfileFromString(ctx.Istio.Profile)
		}

		var context string
		if len(istioModule.Dependencies.Kubectl.CurrentContext) > 0 {
			context = istioModule.Dependencies.Kubectl.CurrentContext
		}

		if len(ctx.Kubernetes.Context) > 0 {
			context = ctx.Kubernetes.Context
			istioModule.Dependencies.Kubectl.CurrentContext = context
		}

		if len(context) == 0 && !ctx.ShowHelp {
			fmt.Println("Remove subcommand requires a context to work with")
			fmt.Println("Use DeploymentToos istio remove --help to show possible arguments")
			startup.Exit(1)
		}

		istioModule.Dependencies.Kubectl.GetAllNamespace()
		// istioModule.Remove(context)
	case "install":
		// istio.NewModule()
		kubeModule := kubectl.NewModule()

		kubeerr := kubeModule.AddKubeConfig(ctx.Kubernetes.KubeConfig)
		if kubeerr != nil {
			logger.Error(kubeerr)
		}
	default:
		istioModuleCommandHelper()
	}
}

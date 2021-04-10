package module

import (
	"fmt"
	"os"
	"strings"

	"github.com/cjlapao/deployment-tools-go/startup"
)

func GetModuleArgument() string {
	args := os.Args[1:]

	if len(args) == 0 || strings.HasPrefix(args[0], "-") {
		PrintCommandHelper()
		startup.Exit(0)
	}

	return args[0]
}

func GetCommandArgument() string {
	args := os.Args[2:]

	if len(args) == 0 || strings.HasPrefix(args[0], "-") {
		return ""
	}

	return args[0]
}

func GetSubCommandArgument() string {
	args := os.Args[3:]

	if len(args) == 0 || strings.HasPrefix(args[0], "-") {
		return ""
	}

	return args[0]
}

func PrintCommandHelper() {
	fmt.Println("Please choose a command:")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  DeploymentTools [command]")
	fmt.Println()
	fmt.Println("Available Commands:")
	fmt.Println("  api       Api module starting the RestApi")
	fmt.Println("  azurecli  Azure client module")
	fmt.Println("  isto      Istio module to control istio installation and removal")
	fmt.Println("  kubectl   Kubectl module to deploy kubernetes manifest into a cluster")
}

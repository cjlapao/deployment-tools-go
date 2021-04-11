package istio

import "fmt"

func istioModuleCommandHelper() {
	fmt.Println("Please choose a sub command:")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  DeploymentTools istio [subcommand]")
	fmt.Println()
	fmt.Println("Available Commands:")
	fmt.Println("  Install         Installs Istio in a Kubernetes cluster")
	fmt.Println("  Remove          Removes Istio from a Kubernetes cluster")
	fmt.Println("  Install-module  Installs Istio Module in the machine")
}

func istioModuleRemoveSubCommandHelper() {
	fmt.Println("Istio Remove:")
	fmt.Println()
	fmt.Println("Example:")
	fmt.Println("  DeploymentTools istio remove --kube-context=server-ctx")
	fmt.Println()
	fmt.Println("Available Flags:")
	fmt.Println("  --kube-config   Adds the kube config to the kubectl for usage with istio")
	fmt.Println("  --kube-context  Sets the default kubectl context for istio to use")
}

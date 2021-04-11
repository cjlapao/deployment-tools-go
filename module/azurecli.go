package module

import (
	"github.com/cjlapao/deployment-tools-go/azurecli"
	"github.com/cjlapao/deployment-tools-go/startup"
)

func AzurecliModuleProcessor() {
	command := GetCommandArgument()

	if len(command) == 0 {
		AzureCliModuleCommandHelper()
		startup.Exit(0)
	}

	out := azurecli.Login()
	if out == nil {
		logger.Fatal("There was an error login in to Azure Portal")
	}
	// download := azurecli.DownloadBlob()
	// if download == nil {
	// 	logger.Fatal("There was an error downloading the file")
	// }
}

func AzureCliModuleCommandHelper() {

}

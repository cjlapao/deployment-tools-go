package azurecli

import (
	"encoding/json"
	"os"

	"github.com/cjlapao/common-go/commands"
	"github.com/cjlapao/common-go/log"
	"github.com/cjlapao/deployment-tools-go/entities"
)

var logger = log.Get()

func Login() *entities.AzureLogin {
	loggedIn := os.Getenv("DEVTOOLS_LOGGED")
	if loggedIn != "true" {
		logger.Notice("Starting Login to Azure Portal")
		out, err := commands.Execute("az", "login", "--service-principal", "-u", ctx.ClientID, "-p", ctx.ClientSecret, "-t", ctx.TenantID)

		if err != nil {
			logger.Error(err)
			return nil
		}

		var response []entities.AzureLogin

		json.Unmarshal([]byte(out), &response)

		if len(response) == 0 {
			return nil
		}

		os.Setenv("DEVTOOLS_LOGGED", "true")

		logger.Success("Logged in successfully to Azure Portal")

		jsonResponse, jsonEncodeError := json.Marshal(&response[0])

		if jsonEncodeError != nil {
			logger.FatalError(jsonEncodeError, "There was an error serializing the azure response")
		}
		os.Setenv("DEVTOOLS_LOGGED_RESPONSE", string(jsonResponse))
		return &response[0]
	}
	logger.Notice("Already Logged In to Azure Portal")

	var envResponse entities.AzureLogin
	json.Unmarshal([]byte(os.Getenv("DEVTOOLS_LOGGED")), &envResponse)

	return &envResponse
}

func Logoff() {
	loggedIn := os.Getenv("DEVTOOLS_LOGGED")
	if loggedIn == "true" {
		logger.Notice("Starting logout from Azure Portal")
		out, err := commands.Execute("az", "logout")

		if err != nil {
			logger.Error(err)
		}

		os.Unsetenv("DEVTOOLS_LOGGED")
		os.Unsetenv("ARM_CLIENT_ID")
		os.Unsetenv("ARM_CLIENT_SECRET")
		os.Unsetenv("ARM_TENANT_ID")

		if len(out) == 0 {
			logger.Success("Successfully logged out of azure")
		}
	} else {
		logger.Notice("Not logged in")
	}
}

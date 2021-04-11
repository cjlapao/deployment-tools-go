package azurecli

import (
	"errors"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure/auth"
)

// GetClientAuthorizer Gets a client Authorizer for Azure
func GetClientAuthorizer() (autorest.Authorizer, error) {
	if !ctx.IsValid() {
		return nil, errors.New("The Azure content is Invalid, please check your flags or environment variables")
	}

	clientAuthorizer := auth.NewClientCredentialsConfig(ctx.ClientID, ctx.ClientSecret, ctx.TenantID)
	return clientAuthorizer.Authorizer()
}

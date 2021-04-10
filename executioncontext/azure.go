package executioncontext

import (
	"os"

	"github.com/cjlapao/common-go/helper"
)

// AzureClientContext entity
type AzureClientContext struct {
	TenantID       string               `json:"tenantId"`
	SubscriptionID string               `json:"subscriptionId"`
	ClientID       string               `json:"clientId"`
	ClientSecret   string               `json:"clientSecret"`
	ResourceGroup  string               `json:"resourceGroup"`
	Storage        *AzureStorageContext `json:"storage"`
}

// AzureStorageContext entity
type AzureStorageContext struct {
	PrimaryAccountKey   string `json:"primaryAccountKey"`
	SecondaryAccountKey string `json:"secondaryAccountKey"`
	AccountName         string `json:"storageAccount"`
	ContainerName       string `json:"storageContainer"`
	FileName            string `json:"fileName"`
	FromPath            string `json:"fromPath"`
	ToFileName          string `json:"toFileName"`
	ToPath              string `json:"toPath"`
}

func (c *AzureClientContext) IsValid() bool {
	if c.ClientID != "" && c.ClientSecret != "" && c.TenantID != "" && c.SubscriptionID != "" {
		return true
	}
	return false
}

func (e *Context) GetAzureContext() {
	azureContext := AzureClientContext{
		TenantID:       helper.GetFlagValue("tenantId", ""),
		SubscriptionID: helper.GetFlagValue("subscriptionId", ""),
		ClientID:       helper.GetFlagValue("clientId", ""),
		ClientSecret:   helper.GetFlagValue("clientSecret", ""),
		ResourceGroup:  helper.GetFlagValue("resourceGroup", ""),
	}
	e.AzureClient = azureContext
	e.SetAzureEnvironmentKeys()

	e.SetAzureStorageContext(nil)
}

func (e *Context) SetAzureEnvironmentKeys() {
	subscriptionID := os.Getenv("DT_AZURE_SUBSCRIPTION_ID")
	tenantID := os.Getenv("DT_AZURE_TENANT_ID")
	clientID := os.Getenv("DT_AZURE_CLIENT_ID")
	clientSecret := os.Getenv("DT_AZURE_CLIENT_SECRET")
	resourceGroup := os.Getenv("DT_AZURE_RESOURCE_GROUP")

	if len(subscriptionID) > 0 {
		e.AzureClient.SubscriptionID = subscriptionID
	}
	if len(tenantID) > 0 {
		e.AzureClient.TenantID = tenantID
	}
	if len(clientID) > 0 {
		e.AzureClient.ClientID = clientID
	}
	if len(clientSecret) > 0 {
		e.AzureClient.ClientSecret = clientSecret
	}
	if len(resourceGroup) > 0 {
		e.AzureClient.ResourceGroup = resourceGroup
	}
}

func (e *Context) SetAzureStorageContext(context *AzureStorageContext) {
	if context != nil {
		e.AzureClient.Storage = context
	} else {
		storageCtx := AzureStorageContext{
			AccountName:   helper.GetFlagValue("storageAccount", ""),
			ContainerName: helper.GetFlagValue("storageContainer", ""),
			FileName:      helper.GetFlagValue("blobName", ""),
			ToFileName:    helper.GetFlagValue("downloadBlobHas", ""),
			FromPath:      helper.GetFlagValue("uploadFrom", ""),
			ToPath:        helper.GetFlagValue("downloadBlobTo", ""),
		}
		e.AzureClient.Storage = &storageCtx
	}
}

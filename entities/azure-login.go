package entities

type AzureLogin struct {
	CloudName        string                    `json:"cloudName"`
	HomeTenantID     string                    `json:"homeTenantId"`
	ID               string                    `json:"id"`
	IsDefault        bool                      `json:"isDefault"`
	ManagedByTenants []interface{}             `json:"managedByTenants"`
	Name             string                    `json:"name"`
	State            string                    `json:"state"`
	TenantID         string                    `json:"tenantId"`
	User             AzureServicePrincipalUser `json:"user"`
}

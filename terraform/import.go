package terraform

import (
	"encoding/json"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

// ImportEntity structure
type ImportEntity struct {
	Azure AzureRmImportEntity `json:"azure" yaml:"azure"`
}

// AzureRmImportEntity structure
type AzureRmImportEntity struct {
	SubscriptionID string           `json:"subscriptionId" yaml:"subscriptionId"`
	ResourceGroup  string           `json:"resourceGroup" yaml:"resourceGroup"`
	TestNumber     int64            `json:"testNumber" yaml:"testNumber"`
	TestBoolean    bool             `json:"testBoolean" yaml:"testBoolean"`
	ServiceBus     ServiceBusEntity `json:"serviceBus" yaml:"serviceBus"`
}

// ServiceBusEntity structure
type ServiceBusEntity struct {
	Namespace string        `json:"namespace" yaml:"namespace"`
	Topics    []TopicEntity `json:"topics" yaml:"topics"`
	Queues    []QueueEntity `json:"queues" yaml:"queues"`
}

// TopicEntity structure
type TopicEntity struct {
	TerraformResourceName string               `json:"terraformResourceName" yaml:"terraformResourceName"`
	TerraformModule       string               `json:"terraformModule" yaml:"terraformModule"`
	Name                  string               `json:"name" yaml:"name"`
	Subscriptions         []SubscriptionEntity `json:"subscriptions" yaml:"subscriptions"`
}

// SubscriptionEntity structure
type SubscriptionEntity struct {
	TerraformResourceName string                   `json:"terraformResourceName" yaml:"terraformResourceName"`
	TerraformModule       string                   `json:"terraformModule" yaml:"terraformModule"`
	Name                  string                   `json:"name" yaml:"name"`
	SubscriptionRules     []SubscriptionRuleEntity `json:"subscriptionRules" yaml:"subscriptionRules"`
}

// SubscriptionRuleEntity structure
type SubscriptionRuleEntity struct {
	TerraformResourceName string `json:"terraformResourceName" yaml:"terraformResourceName"`
	TerraformModule       string `json:"terraformModule" yaml:"terraformModule"`
	Name                  string `json:"name" yaml:"name"`
}

// QueueEntity structure
type QueueEntity struct {
	TerraformResourceName string `json:"terraformResourceName" yaml:"terraformResourceName"`
	TerraformModule       string `json:"terraformModule" yaml:"terraformModule"`
	Name                  string `json:"name" yaml:"name"`
}

// ImportFormat Enum
type ImportFormat int

// ImportOperationEntity Terraform Import Operation entity
type ImportOperationEntity struct {
	Name       string
	Module     string
	ModuleName string
	ResourceID string
}

// Terraform Import Enum Definition
const (
	JSON ImportFormat = iota
	Yaml
)

// ReadImportContent imports a terraform import instruction
func (m *Module) ReadImportContent(content []byte, importFormat ImportFormat) (*ImportEntity, error) {
	var imp ImportEntity

	if len(content) > 0 {
		switch importFormat {
		case JSON:
			logger.Debug("Converting from JSON...")
			err := json.Unmarshal(content, &imp)
			if err != nil {
				return nil, err
			}
			logger.Success("Imported successfully")
			return &imp, nil
		case Yaml:
			logger.Debug("Converting from Yaml...")
			err := yaml.Unmarshal(content, &imp)
			if err != nil {
				return nil, err
			}
			logger.Success("Imported successfully")
			return &imp, nil
		}
	}
	return nil, errors.New("There was an error importing file")
}

// RunImport Starts a terraform import
func (m *Module) RunImport(operation ImportOperationEntity) error {

	return nil
}

package azurecli

import (
	"github.com/cjlapao/deployment-tools-go/executioncontext"
)

var ctx = executioncontext.Get().AzureClient

type AzureClient struct {
}

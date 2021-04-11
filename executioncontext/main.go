package executioncontext

import (
	"os"

	"github.com/cjlapao/common-go/helper"
	"github.com/cjlapao/common-go/log"
	"github.com/cjlapao/deployment-tools-go/common"
	"github.com/cjlapao/deployment-tools-go/kubectl"

	"gopkg.in/yaml.v2"
)

var globalContext *Context
var logger = log.Get()

// Context entity
type Context struct {
	Module      string              `json:"module"`
	Operation   string              `json:"operation"`
	TenantID    string              `json:"tenantId"`
	Database    string              `json:"database"`
	ShowHelp    bool                `json:"help"`
	AzureClient AzureClientContext  `json:"azureClient"`
	Blob        AzureStorageContext `json:"blob"`
	Kubernetes  KubernetesContext   `json:"kubernetes"`
	Istio       IstioContext        `json:"istio"`
	ServiceBus  ServiceBusContext   `json:"service-bus"`
}

// KubernetesContext entity
type KubernetesContext struct {
	KubeConfig string `json:"kube-config"`
	Context    string `json:"context"`
}

type IstioContext struct {
	Profile string `json:"profile"`
	Context string `json:"context"`
}

type ServiceBusContext struct {
	Topic        string `json:"topic"`
	Subscription string `json:"subscription"`
	Queue        string `json:"queue"`
	Message      string `json:"message"`
}

func Get() *Context {
	if globalContext != nil {
		return globalContext
	}

	logger.Debug("Creating Execution Context")
	globalContext = &Context{
		Operation: helper.GetFlagValue("operation", "api"),
		ShowHelp:  helper.GetFlagSwitch("help", false),
		Kubernetes: KubernetesContext{
			KubeConfig: helper.GetFlagValue("kube-config", ""),
			Context:    helper.GetFlagValue("kube-Context", ""),
		},
		Istio: IstioContext{
			Profile: helper.GetFlagValue("istio-profile", "default"),
			Context: helper.GetFlagValue("kube-Context", ""),
		},
		ServiceBus: ServiceBusContext{
			Topic:        helper.GetFlagValue("sb-topic", "default"),
			Subscription: helper.GetFlagValue("sb-subscription", ""),
			Queue:        helper.GetFlagValue("sb-queue-name", ""),
			Message:      helper.GetFlagValue("sb-message", ""),
		},
	}

	globalContext.Getenv()

	if len(globalContext.Kubernetes.KubeConfig) > 0 {
		var kubeConfig kubectl.KubeConfig
		yaml.Unmarshal([]byte(globalContext.Kubernetes.KubeConfig), &kubeConfig)
		configBytes, err := yaml.Marshal(kubeConfig)

		if err != nil {
			os.Exit(1)
		}

		helper.WriteToFile(string(configBytes), common.KubeConfig())
	}

	return globalContext
}

// Getenv gets the environment variables for the entities
func (e *Context) Getenv() {
	e.GetAzureContext()

	op := os.Getenv("DO_OPERATION")
	module := os.Getenv("DT_MODULE")
	kubeConfig := os.Getenv("DT_KUBE_CONFIG")
	kubeContext := os.Getenv("DT_KUBE_CONTEXT")

	istioProfile := os.Getenv("DT_ISTIO_PROFILE")

	if len(op) > 0 {
		e.Operation = op
	}

	if len(module) > 0 {
		e.Module = module
	}

	if len(kubeConfig) > 0 {
		e.Kubernetes.KubeConfig = kubeConfig
	}
	if len(kubeContext) > 0 {
		e.Kubernetes.Context = kubeContext
		e.Istio.Context = kubeContext
	}

	if len(istioProfile) > 0 {
		e.Kubernetes.Context = kubeContext
	}
}

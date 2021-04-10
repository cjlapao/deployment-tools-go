package startup

import "github.com/cjlapao/common-go/helper"

// ToolsFolder Gets the modules tools folder
func KubeConfig() string {
	return helper.JoinPath(helper.GetExecutionPath(), ".temp", "kubeconfig")
}

package common

import "github.com/cjlapao/common-go/helper"

// ToolsFolder Gets the modules tools folder
func KubeConfig() string {
	return helper.JoinPath(helper.GetExecutionPath(), ".temp", "kubeconfig")
}

// ToolsFolder Gets the modules tools folder
func ToolsFolder() string {
	return helper.JoinPath(helper.GetExecutionPath(), ".tools")
}

// ToolsFolder Gets the modules tools folder
func TempFolder() string {
	return helper.JoinPath(helper.GetExecutionPath(), ".temp")
}

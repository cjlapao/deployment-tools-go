package entities

type AzurePipeline struct {
	Name    string  `json:"name"`
	Trigger Trigger `json:"trigger"`
	Pool    string  `json:"pool"`
	Stages  []Stage `json:"stages"`
}

type Stage struct {
	Stage       string `json:"stage"`
	DisplayName string `json:"displayName"`
	Jobs        []Job  `json:"jobs"`
}

type Job struct {
	Job         interface{} `json:"job"`
	DisplayName string      `json:"displayName"`
	Steps       []Step      `json:"steps"`
}

type Step struct {
	Bash             string `json:"bash"`
	DisplayName      string `json:"displayName"`
	WorkingDirectory string `json:"workingDirectory"`
	Env              Env    `json:"env"`
}

type Env struct {
	Debug           bool   `json:"debug"`
	ArmClientID     string `json:"ARM_CLIENT_ID"`
	ArmClientSecret string `json:"ARM_CLIENT_SECRET"`
	ArmTenantID     string `json:"ARM_TENANT_ID"`
}

type Trigger struct {
	Branches Branches `json:"branches"`
}

type Branches struct {
	Include []string `json:"include"`
}

package kubectl

// Version entity
type Version struct {
	ClientVersion VersionInfo `json:"clientVersion"`
	ServerVersion VersionInfo `json:"serverVersion"`
}

// VersionInfo entity
type VersionInfo struct {
	Major        string `json:"major"`
	Minor        string `json:"minor"`
	GitVersion   string `json:"gitVersion"`
	GitCommit    string `json:"gitCommit"`
	GitTreeState string `json:"gitTreeState"`
	BuildDate    string `json:"buildDate"`
	GoVersion    string `json:"goVersion"`
	Compiler     string `json:"compiler"`
	Platform     string `json:"platform"`
}

// KubeConfig entity
type KubeConfig struct {
	APIVersion     string           `json:"apiVersion" yaml:"apiVersion"`
	Kind           string           `json:"kind" yaml:"kind"`
	Preferences    Preferences      `json:"preferences" yaml:"preferences"`
	Clusters       []ClusterElement `json:"clusters" yaml:"clusters"`
	Users          []UserElement    `json:"users" yaml:"users"`
	Contexts       []ContextElement `json:"contexts" yaml:"contexts"`
	CurrentContext string           `json:"current-context" yaml:"current-context"`
}

// ClusterElement entity
type ClusterElement struct {
	Cluster ClusterCluster `json:"cluster" yaml:"cluster"`
	Name    string         `json:"name" yaml:"name"`
}

// ClusterCluster entity
type ClusterCluster struct {
	CertificateAuthorityData string `json:"certificate-authority-data" yaml:"certificate-authority-data"`
	CertificateAuthority     string `json:"certificate-authority" yaml:"certificate-authority"`
	Server                   string `json:"server" yaml:"server"`
}

// ContextElement entity
type ContextElement struct {
	Context ContextContext `json:"context" yaml:"context"`
	Name    string         `json:"name" yaml:"name"`
}

// ContextContext entity
type ContextContext struct {
	Cluster   string `json:"cluster" yaml:"cluster"`
	Namespace string `json:"namespace" yaml:"namespace"`
	User      string `json:"user" yaml:"user"`
}

// Preferences entity
type Preferences struct {
}

// UserElement entity
type UserElement struct {
	Name string   `json:"name" yaml:"name"`
	User UserUser `json:"user" yaml:"user"`
}

// UserUser entity
type UserUser struct {
	AsUserExtra       Preferences `json:"as-user-extra" yaml:"as-user-extra"`
	Token             string      `json:"token" yaml:"token"`
	ClientCertificate string      `json:"client-certificate" yaml:"client-certificate"`
	ClientKey         string      `json:"client-key" yaml:"client-key"`
	Username          string      `json:"username" yaml:"username"`
	Password          string      `json:"password" yaml:"password"`
}

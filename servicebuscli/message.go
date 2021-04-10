package servicebuscli

// MessageEntity Entity
type MessageEntity struct {
	Label      string                 `json:"label"`
	Message    map[string]interface{} `json:"message"`
	Properties map[string]interface{} `json:"properties"`
}

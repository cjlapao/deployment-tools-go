package entities

// LoginErrorResponse entity
type LoginErrorResponse struct {
	Code    string `json:"code"`
	Error   string `json:"error"`
	Message string `json:"message"`
}

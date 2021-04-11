package controllers

import (
	"encoding/json"
	"net/http"
)

type TestResponse struct {
	Message string
}

// Login Generate a token for a valid user
func (c *Controllers) Test(w http.ResponseWriter, r *http.Request) {
	logger.Info("Test Endpoint Hit")
	response := TestResponse{
		Message: "Hello World",
	}

	json.NewEncoder(w).Encode(response)
}

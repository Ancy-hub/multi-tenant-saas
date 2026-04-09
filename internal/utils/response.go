package utils

import (
	"encoding/json"
	"net/http"
)

// APIResponse represents the standard API response structure.
type APIResponse struct {
	// Success indicates if the request was successful.
	Success bool `json:"success"`
	// Data contains the response data when successful.
	Data any `json:"data,omitempty"`
	// Error contains the error message when the request fails.
	Error string `json:"error,omitempty"`
}

// WriteSuccess writes a successful JSON response to the HTTP response writer.
func WriteSuccess(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(APIResponse{
		Success: true,
		Data:    data,
	})
}

// WriteError writes an error JSON response to the HTTP response writer.
func WriteError(w http.ResponseWriter, status int, errMsg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	json.NewEncoder(w).Encode(APIResponse{
		Success: false,
		Error:   errMsg,
	})
}

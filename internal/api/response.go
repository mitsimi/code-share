package api

import (
	"encoding/json"
	"net/http"
	"time"

	"mitsimi.dev/codeShare/internal/api/dto"
)

// WriteSuccess writes a standardized success response.
func WriteSuccess(w http.ResponseWriter, statusCode int, message string, data any) {
	resp := dto.APIResponse{
		StatusCode: statusCode,
		Message:    message,
		Data:       data,
		Timestamp:  time.Now().UTC().Format(time.RFC3339),
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(resp)
}

// WriteError writes a standardized error response.
func WriteError(w http.ResponseWriter, statusCode int, message string) {
	resp := dto.APIResponse{
		StatusCode: statusCode,
		Error:      http.StatusText(statusCode),
		Message:    message,
		Timestamp:  time.Now().UTC().Format(time.RFC3339),
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(resp)
}

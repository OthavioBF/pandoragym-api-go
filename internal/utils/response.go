package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

type SuccessResponse[T any] struct {
	Data T `json:"data,omitempty"`
}

func EncodeJSON[T any](w http.ResponseWriter, statusCode int, data T) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		return fmt.Errorf("failed to encode json: %w", err)
	}
	return nil
}

func WriteJSONResponse[T any](w http.ResponseWriter, statusCode int, data T) {
	if err := EncodeJSON(w, statusCode, data); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func WriteErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	response := ErrorResponse{
		Error:   http.StatusText(statusCode),
		Message: message,
	}
	WriteJSONResponse(w, statusCode, response)
}

func WriteSuccessResponse[T any](w http.ResponseWriter, statusCode int, data T) {
	WriteJSONResponse(w, statusCode, data)
}

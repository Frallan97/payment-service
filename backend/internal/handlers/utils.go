package handlers

import (
	"encoding/json"
	"net/http"
	"payment-service/internal/models"
)

// WriteJSON writes a JSON response
func WriteJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// WriteError writes an error response
func WriteError(w http.ResponseWriter, err error) {
	if apiErr, ok := err.(*models.APIError); ok {
		WriteJSON(w, apiErr.StatusCode, map[string]any{
			"error": map[string]any{
				"code":    apiErr.Code,
				"message": apiErr.Message,
				"details": apiErr.Details,
			},
		})
		return
	}

	// Generic error
	WriteJSON(w, http.StatusInternalServerError, map[string]any{
		"error": map[string]any{
			"code":    models.ErrCodeProviderError,
			"message": "An unexpected error occurred",
		},
	})
}

// DecodeJSON decodes JSON from request body
func DecodeJSON(r *http.Request, v any) error {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		return models.NewAPIError(
			models.ErrCodeInvalidRequest,
			"Invalid JSON in request body",
			http.StatusBadRequest,
		)
	}
	return nil
}

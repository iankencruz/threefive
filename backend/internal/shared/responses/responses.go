package responses

import (
	"encoding/json"
	"net/http"

	"github.com/iankencruz/threefive/internal/shared/errors"
)

// WriteErr writes an error response
func WriteErr(w http.ResponseWriter, err error) {
	var appErr *errors.AppError

	// Check if it's our custom error
	if e, ok := err.(*errors.AppError); ok {
		appErr = e
	} else {
		// Default to internal server error
		appErr = errors.Internal("An unexpected error occurred", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(appErr.StatusCode)

	// If the wrapped error implements json.Marshaler, use it
	if marshaler, ok := appErr.Err.(json.Marshaler); ok {
		json.NewEncoder(w).Encode(marshaler)
	} else {
		// Otherwise, return the standard AppError
		json.NewEncoder(w).Encode(appErr)
	}
}

// WriteJSON writes a JSON response with custom status code
func WriteJSON(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

// Convenience functions
func WriteCreated(w http.ResponseWriter, data any) {
	WriteJSON(w, http.StatusCreated, data)
}

func WriteOK(w http.ResponseWriter, data any) {
	WriteJSON(w, http.StatusOK, data)
}

func WriteNoContent(w http.ResponseWriter) {
	WriteJSON(w, http.StatusNoContent, nil)
}

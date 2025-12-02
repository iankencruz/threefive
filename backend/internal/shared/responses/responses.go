package responses

import (
	"encoding/json"
	"net/http"

	"github.com/iankencruz/threefive/internal/shared/errors"
)

// WriteErr writes an error response.
// It accepts an optional message argument to override the default error message.
// Usage: WriteErr(w, err) or WriteErr(w, err, "Custom user-friendly message")
func WriteErr(w http.ResponseWriter, err error, messages ...string) {
	var appErr *errors.AppError

	// Check if it's our custom error
	if e, ok := err.(*errors.AppError); ok {
		// Create a shallow copy to avoid mutating global sentinels
		shallow := *e
		appErr = &shallow
	} else {
		// Default to internal server error
		appErr = errors.Internal("An unexpected error occurred", err)
	}

	// If a custom message was passed, override the AppError message
	if len(messages) > 0 {
		appErr.Message = messages[0]
	}

	// --- LOGGING LOGIC ---
	// We log both Server Errors (5xx) and Client Errors (4xx) to the terminal.
	// if appErr.StatusCode >= 500 {
	// 	// 500s: Log the underlying error (appErr.Err) which usually contains the stack trace or DB error
	// 	log.Printf("🛑 SERVER ERROR [%s]: %v\n", appErr.Code, appErr.Err)
	// } else {
	// 	// 400s: Log the message so you know why the request failed (e.g., "User not found")
	// 	log.Printf("⚠️  CLIENT ERROR [%d] [%s]: %s\n", appErr.StatusCode, appErr.Code, appErr.Message)
	// }

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

// --- Convenience Success Functions ---

func WriteCreated(w http.ResponseWriter, data any) {
	WriteJSON(w, http.StatusCreated, data)
}

func WriteOK(w http.ResponseWriter, data any) {
	WriteJSON(w, http.StatusOK, data)
}

func WriteNoContent(w http.ResponseWriter) {
	WriteJSON(w, http.StatusNoContent, nil)
}

// --- Convenience Error Functions ---

func WriteBadRequest(w http.ResponseWriter, message, code string) {
	WriteErr(w, errors.BadRequest(message, code))
}

func WriteNotFound(w http.ResponseWriter, message, code string) {
	WriteErr(w, errors.NotFound(message, code))
}

func WriteUnauthorized(w http.ResponseWriter, message, code string) {
	WriteErr(w, errors.Unauthorized(message, code))
}

func WriteForbidden(w http.ResponseWriter, message, code string) {
	WriteErr(w, errors.Forbidden(message, code))
}

func WriteConflict(w http.ResponseWriter, message, code string) {
	WriteErr(w, errors.Conflict(message, code))
}

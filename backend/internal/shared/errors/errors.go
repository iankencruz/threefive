// backend/internal/shared/errors/errors.go
package errors

import (
	"fmt"
	"net/http"
)

// AppError represents application errors with HTTP status codes
type AppError struct {
	Code       string `json:"code"`
	Message    string `json:"message"`
	StatusCode int    `json:"-"`
	Err        error  `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

// Common error constructors
func BadRequest(message, code string) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		StatusCode: http.StatusBadRequest,
	}
}

func Unauthorized(message, code string) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		StatusCode: http.StatusUnauthorized,
	}
}

func Forbidden(message, code string) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		StatusCode: http.StatusForbidden,
	}
}

func NotFound(message, code string) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		StatusCode: http.StatusNotFound,
	}
}

func Conflict(message, code string) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		StatusCode: http.StatusConflict,
	}
}

func Internal(message string, err error) *AppError {
	return &AppError{
		Code:       "internal_error",
		Message:    message,
		StatusCode: http.StatusInternalServerError,
		Err:        err,
	}
}

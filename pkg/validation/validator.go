// pkg/validation/validator.go
package validation

import (
	"fmt"
	"net/url"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"
)

// FieldErrors maps field names to error messages
type FieldErrors map[string]string

// Validator is an interface that any request struct can implement
type Validator interface {
	Validate() (FieldErrors, error)
}

// ValidationRule is a function that validates a value and returns an error message if invalid
type ValidationRule func(value string) string

// Common validation rules
var (
	slugRegex  = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
)

// Rule builders - return validation functions
func Required(message string) ValidationRule {
	if message == "" {
		message = "This field is required"
	}
	return func(value string) string {
		if strings.TrimSpace(value) == "" {
			return message
		}
		return ""
	}
}

func MinLength(min int, message string) ValidationRule {
	if message == "" {
		message = fmt.Sprintf("Must be at least %d characters", min)
	}
	return func(value string) string {
		if len(strings.TrimSpace(value)) < min {
			return message
		}
		return ""
	}
}

func MaxLength(max int, message string) ValidationRule {
	if message == "" {
		message = fmt.Sprintf("Must be at most %d characters", max)
	}
	return func(value string) string {
		if len(value) > max {
			return message
		}
		return ""
	}
}

func IsSlug(message string) ValidationRule {
	if message == "" {
		message = "Must be lowercase, alphanumeric with hyphens only"
	}
	return func(value string) string {
		if value != "" && !slugRegex.MatchString(value) {
			return message
		}
		return ""
	}
}

func IsEmail(message string) ValidationRule {
	if message == "" {
		message = "Invalid email address"
	}
	return func(value string) string {
		if value != "" && !emailRegex.MatchString(value) {
			return message
		}
		return ""
	}
}

func IsURL(message string) ValidationRule {
	if message == "" {
		message = "Invalid URL"
	}
	return func(value string) string {
		if value == "" {
			return ""
		}
		if !strings.HasPrefix(value, "http://") && !strings.HasPrefix(value, "https://") {
			return "URL must start with http:// or https://"
		}
		if _, err := url.ParseRequestURI(value); err != nil {
			return message
		}
		return ""
	}
}

func IsYear(message string) ValidationRule {
	if message == "" {
		message = "Invalid year"
	}
	return func(value string) string {
		if value == "" {
			return ""
		}
		year, err := strconv.Atoi(value)
		if err != nil {
			return message
		}
		if year < 1900 || year > 2100 {
			return "Year must be between 1900 and 2100"
		}
		return ""
	}
}

func IsDate(message string) ValidationRule {
	if message == "" {
		message = "Invalid date format (YYYY-MM-DD)"
	}
	return func(value string) string {
		if value == "" {
			return ""
		}
		if _, err := time.Parse("2006-01-02", value); err != nil {
			return message
		}
		return ""
	}
}

func OneOf(allowed []string, message string) ValidationRule {
	if message == "" {
		message = fmt.Sprintf("Must be one of: %s", strings.Join(allowed, ", "))
	}
	return func(value string) string {
		if value == "" {
			return ""
		}
		if slices.Contains(allowed, value) {
			return ""
		}
		return message
	}
}

// Field represents a field to validate with its rules
type Field struct {
	Name  string
	Value string
	Rules []ValidationRule
}

// ValidateFields validates multiple fields and returns all errors
func ValidateFields(fields []Field) FieldErrors {
	errors := make(FieldErrors)

	for _, field := range fields {
		for _, rule := range field.Rules {
			if errMsg := rule(field.Value); errMsg != "" {
				errors[field.Name] = errMsg
				break // Stop at first error for this field
			}
		}
	}

	return errors
}

// HasErrors returns true if there are any validation errors
func (fe FieldErrors) HasErrors() bool {
	return len(fe) > 0
}

// Error returns a general error message
func (fe FieldErrors) Error() string {
	if !fe.HasErrors() {
		return ""
	}
	return "validation failed"
}

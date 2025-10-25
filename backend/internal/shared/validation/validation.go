// backend/internal/shared/validation/validation.go
package validation

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/iankencruz/threefive/internal/shared/errors"
)

// Common regex patterns
var (
	EmailRegex    = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`)
	SlugRegex     = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*`)
	UsernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]{3,30}`)
	PhoneRegex    = regexp.MustCompile(`^\+?[1-9]\d{1,14}`) // Basic international phone format
	URLRegex      = regexp.MustCompile(`^https?://[^\s/$.?#].[^\s]*`)

	// Password strength patterns
	UppercaseRegex = regexp.MustCompile(`[A-Z]`)
	LowercaseRegex = regexp.MustCompile(`[a-z]`)
	NumberRegex    = regexp.MustCompile(`[0-9]`)
	SpecialRegex   = regexp.MustCompile(`[^a-zA-Z0-9]`)

	// Common validation patterns
	AlphaRegex      = regexp.MustCompile(`^[a-zA-Z]+`)
	AlphaNumRegex   = regexp.MustCompile(`^[a-zA-Z0-9]+`)
	AlphaSpaceRegex = regexp.MustCompile(`^[a-zA-Z\s]+`)
)

// ValidationError represents a single validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidationErrors represents multiple validation errors
type ValidationErrors struct {
	Errors []ValidationError `json:"errors"`
}

func (v ValidationErrors) Error() string {
	var messages []string
	for _, err := range v.Errors {
		messages = append(messages, fmt.Sprintf("%s: %s", err.Field, err.Message))
	}
	return "validation failed: " + strings.Join(messages, ", ")
}

// MarshalJSON implements json.Marshaler interface
func (v ValidationErrors) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Code    string            `json:"code"`
		Message string            `json:"message"`
		Errors  []ValidationError `json:"errors"`
	}{
		Code:    "validation_error",
		Message: "Validation failed",
		Errors:  v.Errors,
	})
}

// Validator holds validation rules and errors
type Validator struct {
	errors []ValidationError
}

// New creates a new validator
func New() *Validator {
	return &Validator{}
}

// AddError adds a validation error
func (v *Validator) AddError(field, message string) {
	v.errors = append(v.errors, ValidationError{
		Field:   field,
		Message: message,
	})
}

// HasErrors returns true if there are validation errors
func (v *Validator) HasErrors() bool {
	return len(v.errors) > 0
}

// Errors returns all validation errors
func (v *Validator) Errors() []ValidationError {
	return v.errors
}

// ToAppError converts validation errors to AppError
func (v *Validator) ToAppError() *errors.AppError {
	return &errors.AppError{
		Code:       "validation_error",
		Message:    "Validation failed",
		StatusCode: http.StatusBadRequest,
		Err:        ValidationErrors{Errors: v.errors},
	}
}

// Common validation rules

// Required checks if a string field is not empty
func (v *Validator) Required(field, value string) {
	if strings.TrimSpace(value) == "" {
		v.AddError(field, "This field is required")
	}
}

// Email validates email format
func (v *Validator) Email(field, value string) {
	if value == "" {
		return // Don't validate empty values, use Required() for that
	}

	if !EmailRegex.MatchString(value) {
		v.AddError(field, "Must be a valid email address")
	}
}

// Slug validates URL-friendly slug format
func (v *Validator) Slug(field, value string) {
	if value == "" {
		return
	}

	if !SlugRegex.MatchString(value) {
		v.AddError(field, "Must be a valid slug (lowercase letters, numbers, and hyphens only)")
	}
}

// Username validates username format
func (v *Validator) Username(field, value string) {
	if value == "" {
		return
	}

	if !UsernameRegex.MatchString(value) {
		v.AddError(field, "Username must be 3-30 characters and contain only letters, numbers, underscores, and hyphens")
	}
}

// Phone validates phone number format
func (v *Validator) Phone(field, value string) {
	if value == "" {
		return
	}

	if !PhoneRegex.MatchString(value) {
		v.AddError(field, "Must be a valid phone number")
	}
}

// URL validates URL format
func (v *Validator) URL(field, value string) {
	if value == "" {
		return
	}

	if !URLRegex.MatchString(value) {
		v.AddError(field, "Must be a valid URL")
	}
}

// Alpha validates that value contains only letters
func (v *Validator) Alpha(field, value string) {
	if value == "" {
		return
	}

	if !AlphaRegex.MatchString(value) {
		v.AddError(field, "Must contain only letters")
	}
}

// AlphaNum validates that value contains only letters and numbers
func (v *Validator) AlphaNum(field, value string) {
	if value == "" {
		return
	}

	if !AlphaNumRegex.MatchString(value) {
		v.AddError(field, "Must contain only letters and numbers")
	}
}

// AlphaSpace validates that value contains only letters and spaces
func (v *Validator) AlphaSpace(field, value string) {
	if value == "" {
		return
	}

	if !AlphaSpaceRegex.MatchString(value) {
		v.AddError(field, "Must contain only letters and spaces")
	}
}

// MinLength validates minimum string length
func (v *Validator) MinLength(field, value string, min int) {
	if value == "" {
		return // Don't validate empty values
	}

	if len(value) < min {
		v.AddError(field, fmt.Sprintf("Must be at least %d characters long", min))
	}
}

// MaxLength validates maximum string length
func (v *Validator) MaxLength(field, value string, max int) {
	if len(value) > max {
		v.AddError(field, fmt.Sprintf("Must be at most %d characters long", max))
	}
}

// Length validates exact string length
func (v *Validator) Length(field, value string, length int) {
	if len(value) != length {
		v.AddError(field, fmt.Sprintf("Must be exactly %d characters long", length))
	}
}

// In validates that value is in allowed list
func (v *Validator) In(field, value string, allowed []string) {
	if value == "" {
		return
	}

	for _, item := range allowed {
		if value == item {
			return
		}
	}

	v.AddError(field, fmt.Sprintf("Must be one of: %s", strings.Join(allowed, ", ")))
}

// Match validates that value matches a regex pattern
func (v *Validator) Match(field, value string, pattern *regexp.Regexp, message string) {
	if value == "" {
		return
	}

	if !pattern.MatchString(value) {
		v.AddError(field, message)
	}
}

// StrongPassword validates password strength
func (v *Validator) StrongPassword(field, value string) {
	if value == "" {
		return
	}

	if len(value) < 8 {
		v.AddError(field, "Password must be at least 8 characters long")
		return
	}

	if !UppercaseRegex.MatchString(value) {
		v.AddError(field, "Password must contain at least one uppercase letter")
	}
	if !LowercaseRegex.MatchString(value) {
		v.AddError(field, "Password must contain at least one lowercase letter")
	}
	if !NumberRegex.MatchString(value) {
		v.AddError(field, "Password must contain at least one number")
	}
	if !SpecialRegex.MatchString(value) {
		v.AddError(field, "Password must contain at least one special character")
	}
}

// ParseAndValidateJSON parses JSON request body and validates it
func ParseAndValidateJSON(r *http.Request, dest any, validateFn func(*Validator)) error {
	// Parse JSON
	if err := json.NewDecoder(r.Body).Decode(dest); err != nil {
		return errors.BadRequest("Invalid JSON format", "invalid_json")
	}

	// Validate
	validator := New()
	validateFn(validator)

	if validator.HasErrors() {
		return validator.ToAppError()
	}

	return nil
}

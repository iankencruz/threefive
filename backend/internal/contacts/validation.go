// internal/contacts/validation.go
package contacts

import (
	"github.com/iankencruz/threefive/internal/shared/validation"
)

// Valid contact statuses
var ValidContactStatuses = []string{"new", "read", "archived"}

// Validate validates a ContactRequest
func (r *ContactRequest) Validate(v *validation.Validator) {
	// Required fields
	v.Required("name", r.Name)
	v.Required("email", r.Email)
	v.Required("message", r.Message)

	// Name validation
	if r.Name != "" {
		v.MinLength("name", r.Name, 2)
		v.MaxLength("name", r.Name, 100)
	}

	// Email validation
	if r.Email != "" {
		v.Email("email", r.Email)
		v.MaxLength("email", r.Email, 255)
	}

	// Subject validation (optional)
	if r.Subject != "" {
		v.MaxLength("subject", r.Subject, 200)
	}

	// Message validation
	if r.Message != "" {
		v.MinLength("message", r.Message, 10)
		v.MaxLength("message", r.Message, 5000)
	}
}

// Validate validates an UpdateContactStatusRequest
func (r *UpdateContactStatusRequest) Validate(v *validation.Validator) {
	v.Required("status", r.Status)
	v.In("status", r.Status, ValidContactStatuses)
}

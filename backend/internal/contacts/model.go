// internal/contacts/model.go
package contacts

import (
	"time"

	"github.com/google/uuid"
)

// ============================================
// Request Models
// ============================================

// ContactRequest represents the incoming contact form submission
type ContactRequest struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

// UpdateContactStatusRequest represents the request to update contact status
type UpdateContactStatusRequest struct {
	Status string `json:"status"`
}

// ListContactsParams represents pagination and filtering parameters
type ListContactsParams struct {
	Status  string `json:"status"`
	Limit   int32  `json:"limit"`
	Offset  int32  `json:"offset"`
	OrderBy string `json:"order_by"`
}

// ============================================
// Response Models
// ============================================

// ContactResponse represents the contact data returned to clients
type ContactResponse struct {
	ID        uuid.UUID  `json:"id"`
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	Subject   *string    `json:"subject"`
	Message   string     `json:"message"`
	Status    string     `json:"status"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

// ContactListResponse represents paginated contact list
type ContactListResponse struct {
	Contacts   []ContactResponse `json:"contacts"`
	Total      int64             `json:"total"`
	Limit      int32             `json:"limit"`
	Offset     int32             `json:"offset"`
	TotalPages int64             `json:"total_pages"`
}

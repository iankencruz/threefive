package auth

import (
	"time"

	"github.com/google/uuid"
	"github.com/iankencruz/threefive/internal/shared/validation"
)

// SessionWithUser represents a user session with extended user information
// This is needed because our GetSessionByToken query joins session + user data
type SessionWithUser struct {
	// Session fields
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Token     string    `json:"-"` // Never serialize the token
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	IPAddress *string   `json:"ip_address,omitempty"`
	UserAgent *string   `json:"user_agent,omitempty"`
	IsActive  bool      `json:"is_active"`

	// User fields (from JOIN)
	Email         string    `json:"email"`
	FirstName     string    `json:"first_name"`
	LastName      string    `json:"last_name"`
	UserCreatedAt time.Time `json:"user_created_at"`
}

// FullName returns the user's full name
func (s *SessionWithUser) FullName() string {
	return s.FirstName + " " + s.LastName
}

// LoginRequest represents the login request payload
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Validate validates the login request
func (r *LoginRequest) Validate(v *validation.Validator) {
	v.Required("email", r.Email)
	v.Email("email", r.Email)
	v.Required("password", r.Password)
	v.MinLength("password", r.Password, 6)
}

// RegisterRequest represents the registration request payload
type RegisterRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

// Validate validates the registration request
func (r *RegisterRequest) Validate(v *validation.Validator) {
	v.Required("first_name", r.FirstName)
	v.AlphaSpace("first_name", r.FirstName)
	v.MinLength("first_name", r.FirstName, 2)
	v.MaxLength("first_name", r.FirstName, 50)

	v.Required("last_name", r.LastName)
	v.AlphaSpace("last_name", r.LastName)
	v.MinLength("last_name", r.LastName, 2)
	v.MaxLength("last_name", r.LastName, 50)

	v.Required("email", r.Email)
	v.Email("email", r.Email)

	v.Required("password", r.Password)
	v.StrongPassword("password", r.Password)
}

// PasswordResetRequest represents password reset request
type PasswordResetRequest struct {
	Email string `json:"email"`
}

// Validate validates the password reset request
func (r *PasswordResetRequest) Validate(v *validation.Validator) {
	v.Required("email", r.Email)
	v.Email("email", r.Email)
}

// PasswordResetConfirmRequest represents password reset confirmation
type PasswordResetConfirmRequest struct {
	Token       string `json:"token"`
	NewPassword string `json:"new_password"`
}

// Validate validates the password reset confirmation request
func (r *PasswordResetConfirmRequest) Validate(v *validation.Validator) {
	v.Required("token", r.Token)
	v.Required("new_password", r.NewPassword)
	v.StrongPassword("new_password", r.NewPassword)
}

// ChangePasswordRequest represents password change request
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

// Validate validates the change password request
func (r *ChangePasswordRequest) Validate(v *validation.Validator) {
	v.Required("current_password", r.CurrentPassword)
	v.Required("new_password", r.NewPassword)
	v.StrongPassword("new_password", r.NewPassword)
}

// SessionConfig holds session configuration
type SessionConfig struct {
	Duration    time.Duration
	IdleTimeout time.Duration
	CookieName  string
	Domain      string
	Secure      bool
	HTTPOnly    bool
	SameSite    string
}

// DefaultSessionConfig returns default session configuration
func DefaultSessionConfig() SessionConfig {
	return SessionConfig{
		Duration:    24 * time.Hour, // 24 hours
		IdleTimeout: 2 * time.Hour,  // 2 hours of inactivity
		CookieName:  "session_token",
		Domain:      "",    // Will be set based on environment
		Secure:      true,  // HTTPS only in production
		HTTPOnly:    true,  // Prevent XSS
		SameSite:    "lax", // Good for SSR
	}
}

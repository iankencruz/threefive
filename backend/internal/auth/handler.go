// backend/internal/auth/handler.go
package auth

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/iankencruz/threefive/internal/shared/errors"
	"github.com/iankencruz/threefive/internal/shared/responses"
	"github.com/iankencruz/threefive/internal/shared/session"
	"github.com/iankencruz/threefive/internal/shared/sqlc"
	"github.com/iankencruz/threefive/internal/shared/validation"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Handler struct {
	service        *Service
	sessionManager *session.Manager
}

// NewHandler creates a new auth handler with its own service
func NewHandler(db *pgxpool.Pool, queries *sqlc.Queries, sessionManager *session.Manager) *Handler {
	service := NewService(db, queries, sessionManager)
	return &Handler{
		service:        service,
		sessionManager: sessionManager,
	}
}

// Register handles user registration
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest

	// Parse and validate request
	err := validation.ParseAndValidateJSON(r, &req, func(v *validation.Validator) {
		req.Validate(v)
	})
	if err != nil {
		responses.WriteErr(w, err)
		return
	}

	// Register user
	user, session, err := h.service.Register(r.Context(), req, r)
	if err != nil {
		responses.WriteErr(w, err)
		return
	}

	// Set session cookie
	h.sessionManager.SetSessionCookie(w, session.Token)

	// Return user data (without sensitive info)
	userResponse := map[string]any{
		"id":         user.ID,
		"email":      user.Email,
		"first_name": user.FirstName,
		"last_name":  user.LastName,
		"created_at": user.CreatedAt,
	}

	response := map[string]any{
		"user":    userResponse,
		"message": "Registration successful",
	}

	responses.WriteCreated(w, response)
}

// Login handles user login
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest

	// Parse and validate request
	err := validation.ParseAndValidateJSON(r, &req, func(v *validation.Validator) {
		req.Validate(v)
	})
	if err != nil {
		responses.WriteErr(w, err)
		return
	}

	// Login user
	user, session, err := h.service.Login(r.Context(), req, r)
	if err != nil {
		fmt.Print("Error Here")
		responses.WriteErr(w, err)
		return
	}

	// Set session cookie
	h.sessionManager.SetSessionCookie(w, session.Token)

	// Return user data (without sensitive info)
	userResponse := map[string]any{
		"id":         user.ID,
		"email":      user.Email,
		"first_name": user.FirstName,
		"last_name":  user.LastName,
		"created_at": user.CreatedAt,
	}

	response := map[string]any{
		"user":    userResponse,
		"message": "Login successful",
	}

	responses.WriteOK(w, response)
}

// Logout handles user logout
func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	// Get session token from request
	token := h.sessionManager.GetSessionTokenFromRequest(r)

	// Logout (deactivate session)
	err := h.service.Logout(r.Context(), token)
	if err != nil {
		responses.WriteErr(w, err)
		return
	}

	// Clear session cookie
	h.sessionManager.ClearSessionCookie(w)

	response := map[string]any{
		"message": "Logout successful",
	}

	responses.WriteOK(w, response)
}

// LogoutAll handles logout from all devices
func (h *Handler) LogoutAll(w http.ResponseWriter, r *http.Request) {
	// Get current user from context (middleware ensures this exists)
	user := MustGetUserFromContext(r.Context())

	// Logout from all devices
	err := h.service.LogoutAll(r.Context(), user.ID)
	if err != nil {
		responses.WriteErr(w, err)
		return
	}

	// Clear current session cookie
	h.sessionManager.ClearSessionCookie(w)

	response := map[string]any{
		"message": "Logged out from all devices",
	}

	responses.WriteOK(w, response)
}

// Me returns current user information
func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	// Get current user from context (middleware ensures this exists)
	user := MustGetUserFromContext(r.Context())

	userResponse := map[string]any{
		"id":         user.ID,
		"email":      user.Email,
		"first_name": user.FirstName,
		"last_name":  user.LastName,
	}

	responses.WriteOK(w, userResponse)
}

// ChangePassword handles password change
func (h *Handler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	var req ChangePasswordRequest

	// Parse and validate request
	err := validation.ParseAndValidateJSON(r, &req, func(v *validation.Validator) {
		req.Validate(v)
	})
	if err != nil {
		responses.WriteErr(w, err)
		return
	}

	// Get current user from context
	user := MustGetUserFromContext(r.Context())

	// Change password
	err = h.service.ChangePassword(r.Context(), user.ID, req)
	if err != nil {
		responses.WriteErr(w, err)
		return
	}

	// Clear session cookie (user needs to login again)
	h.sessionManager.ClearSessionCookie(w)

	response := map[string]any{
		"message": "Password changed successfully. Please login again.",
	}

	responses.WriteOK(w, response)
}

// RequestPasswordReset handles password reset request
func (h *Handler) RequestPasswordReset(w http.ResponseWriter, r *http.Request) {
	var req PasswordResetRequest

	// Parse and validate request
	err := validation.ParseAndValidateJSON(r, &req, func(v *validation.Validator) {
		req.Validate(v)
	})
	if err != nil {
		responses.WriteErr(w, err)
		return
	}

	// Request password reset
	err = h.service.RequestPasswordReset(r.Context(), req)
	if err != nil {
		responses.WriteErr(w, err)
		return
	}

	response := map[string]any{
		"message": "If an account with that email exists, a password reset link has been sent.",
	}

	responses.WriteOK(w, response)
}

// ResetPassword handles password reset confirmation
func (h *Handler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var req PasswordResetConfirmRequest

	// Parse and validate request
	err := validation.ParseAndValidateJSON(r, &req, func(v *validation.Validator) {
		req.Validate(v)
	})
	if err != nil {
		responses.WriteErr(w, err)
		return
	}

	// Reset password
	err = h.service.ResetPassword(r.Context(), req)
	if err != nil {
		responses.WriteErr(w, err)
		return
	}

	response := map[string]any{
		"message": "Password reset successful. Please login with your new password.",
	}

	responses.WriteOK(w, response)
}

// GetSessions returns all active sessions for the current user
func (h *Handler) GetSessions(w http.ResponseWriter, r *http.Request) {
	// Get current user from context
	user := MustGetUserFromContext(r.Context())

	// Get user sessions
	sessions, err := h.service.GetUserSessions(r.Context(), user.ID)
	if err != nil {
		responses.WriteErr(w, err)
		return
	}

	// Format sessions for response (hide sensitive info)
	var sessionResponses []map[string]any
	for _, session := range sessions {
		sessionResponse := map[string]any{
			"id":         session.ID,
			"created_at": session.CreatedAt,
			"updated_at": session.UpdatedAt,
			"expires_at": session.ExpiresAt,
		}

		// Add IP address if available
		if session.IpAddress != nil {
			sessionResponse["ip_address"] = session.IpAddress.String()
		}

		// Add user agent if available
		if session.UserAgent.Valid {
			sessionResponse["user_agent"] = session.UserAgent.String
		}

		sessionResponses = append(sessionResponses, sessionResponse)
	}

	response := map[string]any{
		"sessions": sessionResponses,
	}

	responses.WriteOK(w, response)
}

// Updated RevokeSession handler to extract sessionID from URL
func (h *Handler) RevokeSession(w http.ResponseWriter, r *http.Request) {
	// Extract session ID from URL parameter
	sessionIDStr := chi.URLParam(r, "sessionID")
	if sessionIDStr == "" {
		responses.WriteErr(w, errors.BadRequest("Session ID is required", "missing_session_id"))
		return
	}

	// Parse session ID
	sessionID, err := uuid.Parse(sessionIDStr)
	if err != nil {
		responses.WriteErr(w, errors.BadRequest("Invalid session ID format", "invalid_session_id"))
		return
	}

	// Get current user from context
	user := MustGetUserFromContext(r.Context())

	// Revoke the session
	err = h.service.RevokeSession(r.Context(), sessionID, user.ID)
	if err != nil {
		responses.WriteErr(w, err)
		return
	}

	response := map[string]any{
		"message": "Session revoked successfully",
	}

	responses.WriteOK(w, response)
}

// RefreshSession extends the current session
func (h *Handler) RefreshSession(w http.ResponseWriter, r *http.Request) {
	// Get session token from request
	token := h.sessionManager.GetSessionTokenFromRequest(r)

	// Refresh session
	err := h.service.RefreshSession(r.Context(), token)
	if err != nil {
		responses.WriteErr(w, err)
		return
	}

	response := map[string]any{
		"message": "Session refreshed successfully",
	}

	responses.WriteOK(w, response)
}

package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/iankencruz/threefive/internal/shared/errors"
	"github.com/iankencruz/threefive/internal/shared/session"
	"github.com/iankencruz/threefive/internal/shared/sqlc"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

// Service handles authentication business logic
type Service struct {
	db             *pgxpool.Pool
	queries        *sqlc.Queries
	sessionManager *session.Manager
}

// NewService creates a new auth service
func NewService(db *pgxpool.Pool, queries *sqlc.Queries, sessionManager *session.Manager) *Service {
	return &Service{
		db:             db,
		queries:        queries,
		sessionManager: sessionManager,
	}
}

// Register creates a new user account
func (s *Service) Register(ctx context.Context, req *RegisterRequest, r *http.Request) (*sqlc.Users, sqlc.Sessions, error) {
	// Check if user already exists
	_, err := s.queries.GetUserByEmail(ctx, req.Email)
	if err == nil {
		return nil, sqlc.Sessions{}, errors.Conflict("User with this email already exists", "email_exists")
	}
	if err != pgx.ErrNoRows {
		return nil, sqlc.Sessions{}, errors.Internal("Failed to check existing user", err)
	}

	// Hash password
	hashedPassword, err := s.hashPassword(req.Password)
	if err != nil {
		return nil, sqlc.Sessions{}, errors.Internal("Failed to hash password", err)
	}

	// Create user
	user, err := s.queries.CreateUser(ctx, sqlc.CreateUserParams{
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Email:        req.Email,
		PasswordHash: hashedPassword,
	})
	if err != nil {
		return nil, sqlc.Sessions{}, errors.Internal("Failed to create user", err)
	}

	// Create session
	session, err := s.sessionManager.CreateSession(ctx, user.ID, r)
	if err != nil {
		return nil, sqlc.Sessions{}, errors.Internal("Failed to create session", err)
	}

	return &user, session, nil
}

// Login authenticates a user and creates a session
func (s *Service) Login(ctx context.Context, req *LoginRequest, r *http.Request) (*sqlc.Users, sqlc.Sessions, error) {
	// Get user by email
	user, err := s.queries.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, sqlc.Sessions{}, errors.Unauthorized("Invalid email or password", "invalid_credentials")
		}
		return nil, sqlc.Sessions{}, errors.Internal("Failed to get user", err)
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, sqlc.Sessions{}, errors.Unauthorized("Invalid email or password", "invalid_credentials")
	}

	// Create session
	session, err := s.sessionManager.CreateSession(ctx, user.ID, r)
	if err != nil {
		return nil, sqlc.Sessions{}, errors.Internal("Failed to create session", err)
	}

	return &user, session, nil
}

// Logout deactivates the current session
func (s *Service) Logout(ctx context.Context, token string) error {
	if token == "" {
		return errors.BadRequest("No session token provided", "missing_token")
	}

	err := s.sessionManager.DeactivateSession(ctx, token)
	if err != nil {
		return errors.Internal("Failed to logout", err)
	}

	return nil
}

// LogoutAll deactivates all sessions for a user
func (s *Service) LogoutAll(ctx context.Context, userID uuid.UUID) error {
	err := s.sessionManager.DeactivateAllUserSessions(ctx, userID)
	if err != nil {
		return errors.Internal("Failed to logout from all devices", err)
	}

	return nil
}

// GetCurrentUser returns the current authenticated user from session
func (s *Service) GetCurrentUser(ctx context.Context, token string) (sqlc.GetSessionByTokenRow, error) {
	if token == "" {
		return sqlc.GetSessionByTokenRow{}, errors.Unauthorized("No session token provided", "missing_token")
	}

	sessionWithUser, err := s.sessionManager.GetSession(ctx, token)
	if err != nil {
		return sqlc.GetSessionByTokenRow{}, errors.Unauthorized("Invalid or expired session", "invalid_session")
	}

	return sessionWithUser, nil
}

// ChangePassword changes a user's password
func (s *Service) ChangePassword(ctx context.Context, userID uuid.UUID, req *ChangePasswordRequest) error {
	// Get current user
	user, err := s.queries.GetUserByID(ctx, userID)
	if err != nil {
		return errors.NotFound("User not found", "user_not_found")
	}

	// Verify current password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.CurrentPassword)); err != nil {
		return errors.BadRequest("Current password is incorrect", "invalid_current_password")
	}

	// Hash new password
	newHashedPassword, err := s.hashPassword(req.NewPassword)
	if err != nil {
		return errors.Internal("Failed to hash new password", err)
	}

	// Update password
	err = s.queries.UpdateUserPassword(ctx, sqlc.UpdateUserPasswordParams{
		PasswordHash: newHashedPassword,
		UserID:       userID,
	})
	if err != nil {
		return errors.Internal("Failed to update password", err)
	}

	// Deactivate all user sessions (they'll need to login again)
	err = s.sessionManager.DeactivateAllUserSessions(ctx, userID)
	if err != nil {
		// Log but don't fail - password was changed successfully
		// You might want to add proper logging here
	}

	return nil
}

// RequestPasswordReset creates a password reset token
func (s *Service) RequestPasswordReset(ctx context.Context, req *PasswordResetRequest) error {
	// Check if user exists
	user, err := s.queries.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if err == pgx.ErrNoRows {
			// Don't reveal whether email exists or not - always return success
			return nil
		}
		return errors.Internal("Failed to check user", err)
	}

	// Generate reset token
	token, err := s.generateSecureToken()
	if err != nil {
		return errors.Internal("Failed to generate reset token", err)
	}

	// Create password reset token (valid for 1 hour)
	expiresAt := time.Now().Add(1 * time.Hour)
	_, err = s.queries.CreatePasswordResetToken(ctx, sqlc.CreatePasswordResetTokenParams{
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: expiresAt,
	})
	if err != nil {
		return errors.Internal("Failed to create reset token", err)
	}

	// TODO: Send email with reset link
	// In a real application, you would send an email here
	// For now, you might want to log the token for testing
	// log.Printf("Password reset token for %s: %s", user.Email, token)

	return nil
}

// ResetPassword resets a user's password using a reset token
func (s *Service) ResetPassword(ctx context.Context, req *PasswordResetConfirmRequest) error {
	// Get and validate reset token
	resetToken, err := s.queries.GetPasswordResetToken(ctx, req.Token)
	if err != nil {
		if err == pgx.ErrNoRows {
			return errors.BadRequest("Invalid or expired reset token", "invalid_reset_token")
		}
		return errors.Internal("Failed to get reset token", err)
	}

	// Hash new password
	newHashedPassword, err := s.hashPassword(req.NewPassword)
	if err != nil {
		return errors.Internal("Failed to hash new password", err)
	}

	// Start transaction
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return errors.Internal("Failed to start transaction", err)
	}
	defer tx.Rollback(ctx)

	qtx := s.queries.WithTx(tx)

	// Update password
	err = qtx.UpdateUserPassword(ctx, sqlc.UpdateUserPasswordParams{
		PasswordHash: newHashedPassword,
		UserID:       resetToken.UserID,
	})
	if err != nil {
		return errors.Internal("Failed to update password", err)
	}

	// Mark reset token as used
	err = qtx.UsePasswordResetToken(ctx, req.Token)
	if err != nil {
		return errors.Internal("Failed to mark token as used", err)
	}

	// Deactivate all user sessions
	err = qtx.DeactivateAllUserSessions(ctx, resetToken.UserID)
	if err != nil {
		return errors.Internal("Failed to deactivate sessions", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return errors.Internal("Failed to commit transaction", err)
	}

	return nil
}

// GetUserSessions returns all active sessions for a user
func (s *Service) GetUserSessions(ctx context.Context, userID uuid.UUID) ([]sqlc.Sessions, error) {
	sessions, err := s.sessionManager.GetUserActiveSessions(ctx, userID)
	if err != nil {
		return nil, errors.Internal("Failed to get user sessions", err)
	}
	return sessions, nil
}

// RevokeSession revokes a specific session
func (s *Service) RevokeSession(ctx context.Context, sessionID uuid.UUID, userID uuid.UUID) error {
	// First verify the session belongs to the user
	sessions, err := s.sessionManager.GetUserActiveSessions(ctx, userID)
	if err != nil {
		return errors.Internal("Failed to get user sessions", err)
	}

	// Find the session and get its token
	var sessionToken string
	for _, session := range sessions {
		if session.ID == sessionID {
			sessionToken = session.Token
			break
		}
	}

	if sessionToken == "" {
		return errors.NotFound("Session not found or already inactive", "session_not_found")
	}

	// Deactivate the session
	err = s.sessionManager.DeactivateSession(ctx, sessionToken)
	if err != nil {
		return errors.Internal("Failed to revoke session", err)
	}

	return nil
}

// RefreshSession extends the current session
func (s *Service) RefreshSession(ctx context.Context, token string) error {
	if token == "" {
		return errors.BadRequest("No session token provided", "missing_token")
	}

	err := s.sessionManager.RefreshSession(ctx, token)
	if err != nil {
		return errors.Internal("Failed to refresh session", err)
	}

	return nil
}

// Helper methods

func (s *Service) hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (s *Service) generateSecureToken() (string, error) {
	bytes := make([]byte, 32) // 256-bit token
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

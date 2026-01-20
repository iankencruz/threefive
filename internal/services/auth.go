// internal/services/auth.go
package services

import (
	"context"
	"log/slog"

	"github.com/iankencruz/threefive/database/generated"
	"github.com/iankencruz/threefive/pkg/errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	db      *pgxpool.Pool
	queries *generated.Queries
	logger  *slog.Logger
}

func NewAuthService(db *pgxpool.Pool, queries *generated.Queries, logger *slog.Logger) *AuthService {
	return &AuthService{
		db:      db,
		queries: queries,
		logger:  logger.With("component", "auth_service"),
	}
}

// Authenticate verifies user credentials and returns the user ID
func (s *AuthService) Authenticate(ctx context.Context, email, password string) (string, error) {
	s.logger.Debug("attempting to authenticate user",
		"email", email,
	)
	// Get user by email
	user, err := s.queries.GetUserByEmail(ctx, email)
	if err != nil {
		if err == pgx.ErrNoRows {
			s.logger.Warn("authentication failed - user not found",
				"email", email,
			)
			return "", errors.Unauthorized("Invalid email or password")
		}
		s.logger.Error("failed to query user by email",
			"error", err,
			"email", email,
		)
		return "", errors.Internal("Failed to query user", err)
	}

	// Verify password
	s.logger.Info("password comparison details",
		"email", email,
		"stored_hash_length", len(user.PasswordHash),
		"stored_hash_preview", user.PasswordHash[:20]+"...",
		"input_password_length", len(password),
		"input_password", password, // Remove this after debugging!
	)

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		s.logger.Warn("authentication failed - invalid password",
			"email", email,
			"user_id", user.ID,
			"bcrypt_error", err.Error(),
		)
		return "", errors.Unauthorized("Invalid email or password")
	}

	s.logger.Info("user authenticated successfully",
		"email", email,
		"user_id", user.ID,
	)

	return user.ID.String(), nil
}

// HashPassword hashes a password using bcrypt
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", errors.Internal("Failed to hash password", err)
	}
	return string(bytes), nil
}

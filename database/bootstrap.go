// database/bootstrap.go
package database

import (
	"context"
	"log/slog"
	"os"

	"github.com/iankencruz/threefive/database/generated"
	"github.com/iankencruz/threefive/internal/services"
	"github.com/jackc/pgx/v5"
)

// Bootstrap ensures required data exists in the database
func Bootstrap(ctx context.Context, queries *generated.Queries, logger *slog.Logger) error {
	logger.Info("running database bootstrap")

	// Ensure admin user exists
	if err := ensureAdminUser(ctx, queries, logger); err != nil {
		return err
	}

	logger.Info("database bootstrap complete")
	return nil
}

// ensureAdminUser creates the default admin user if it doesn't exist
func ensureAdminUser(ctx context.Context, queries *generated.Queries, logger *slog.Logger) error {
	// Get admin credentials from environment with defaults
	adminEmail := getEnv("ADMIN_EMAIL", "admin@example.com")
	adminPassword := getEnv("ADMIN_PASSWORD", "Password!123")
	adminFirstName := getEnv("ADMIN_FIRST_NAME", "Admin")
	adminLastName := getEnv("ADMIN_LAST_NAME", "User")

	// Check if admin user already exists
	_, err := queries.GetUserByEmail(ctx, adminEmail)
	if err == nil {
		logger.Info("admin user already exists, skipping creation",
			"email", adminEmail,
		)
		return nil
	}

	// If error is not "no rows", something went wrong
	if err != pgx.ErrNoRows {
		logger.Error("failed to check for admin user",
			"error", err,
		)
		return err
	}

	// Admin user doesn't exist, create it
	logger.Info("creating default admin user",
		"email", adminEmail,
	)

	// Hash the password
	passwordHash, err := services.HashPassword(adminPassword)
	if err != nil {
		logger.Error("failed to hash admin password",
			"error", err,
		)
		return err
	}

	// Create the user
	user, err := queries.CreateUser(ctx, generated.CreateUserParams{
		FirstName:    adminFirstName,
		LastName:     adminLastName,
		Email:        adminEmail,
		PasswordHash: passwordHash,
	})
	if err != nil {
		logger.Error("failed to create admin user",
			"error", err,
		)
		return err
	}

	logger.Info("default admin user created successfully",
		"email", adminEmail,
		"user_id", user.ID,
	)

	return nil
}

// getEnv gets an environment variable with a fallback default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

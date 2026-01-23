// internal/services/auth_test.go
package services

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/iankencruz/threefive/database/generated"
	"github.com/iankencruz/threefive/pkg/errors"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"golang.org/x/crypto/bcrypt"
)

// runTestMigrations runs migrations (copied here to avoid import cycle)
func runTestMigrations(pool *pgxpool.Pool) error {
	connConfig := pool.Config().ConnConfig
	connStr := stdlib.RegisterConnConfig(connConfig)

	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return fmt.Errorf("failed to open sql connection: %w", err)
	}
	defer db.Close()

	migrationsPath, err := findMigrationsDir()
	if err != nil {
		return fmt.Errorf("failed to find migrations directory: %w", err)
	}

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set dialect: %w", err)
	}

	if err := goose.Up(db, migrationsPath); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}

// findMigrationsDir locates the migrations directory
func findMigrationsDir() (string, error) {
	candidates := []string{
		"../../migrations",
		"../../../migrations",
		"migrations",
		"./migrations",
	}

	for _, candidate := range candidates {
		absPath, err := filepath.Abs(candidate)
		if err != nil {
			continue
		}

		if _, err := os.Stat(absPath); err == nil {
			return absPath, nil
		}
	}

	return "", fmt.Errorf("migrations directory not found")
}

// setupTestContainer creates a PostgreSQL container and runs migrations
func setupTestContainer(t *testing.T) (*pgxpool.Pool, *generated.Queries, *slog.Logger, func()) {
	t.Helper()

	ctx := context.Background()

	// Start PostgreSQL container
	pgContainer, err := postgres.Run(
		ctx,
		"postgres:18-alpine",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second),
		),
	)
	require.NoError(t, err, "failed to start postgres container")

	// Get connection string
	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err, "failed to get connection string")

	// Create connection pool
	pool, err := pgxpool.New(ctx, connStr)
	require.NoError(t, err, "failed to create connection pool")

	// Run migrations from filesystem
	err = runTestMigrations(pool)
	require.NoError(t, err, "failed to run migrations")

	// Create queries and logger
	queries := generated.New(pool)
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	// Cleanup function
	cleanup := func() {
		pool.Close()
		if err := pgContainer.Terminate(ctx); err != nil {
			t.Logf("failed to terminate container: %s", err)
		}
	}

	return pool, queries, logger, cleanup
}

// createTestUser creates a user with known credentials for testing
func createTestUser(t *testing.T, ctx context.Context, queries *generated.Queries, email, password, firstName, lastName string) generated.User {
	t.Helper()

	hashedPassword, err := HashPassword(password)
	require.NoError(t, err, "failed to hash password")

	user, err := queries.CreateUser(ctx, generated.CreateUserParams{
		FirstName:    firstName,
		LastName:     lastName,
		Email:        email,
		PasswordHash: hashedPassword,
	})
	require.NoError(t, err, "failed to create user")

	// Clean up after test
	t.Cleanup(func() {
		_ = queries.HardDeleteUser(ctx, user.ID)
	})

	return user
}

// TestAuthService_Authenticate tests the authentication logic with various scenarios
func TestAuthService_Authenticate(t *testing.T) {
	// Setup: Create test container
	pool, queries, logger, cleanup := setupTestContainer(t)
	defer cleanup()

	ctx := context.Background()
	authService := NewAuthService(pool, queries, logger)

	// Define test scenarios
	tests := []struct {
		name      string                        // Description of what we're testing
		setupData func() string                 // Returns expected user ID if user created
		email     string                        // Input: email to authenticate
		password  string                        // Input: password to authenticate
		wantErr   bool                          // Should authentication fail?
		checkErr  func(t *testing.T, err error) // Custom error validation
	}{
		{
			name: "success: valid credentials authenticate successfully",
			setupData: func() string {
				user := createTestUser(t, ctx, queries,
					"valid@example.com", "ValidPass123!", "Test", "User")
				return user.ID.String()
			},
			email:    "valid@example.com",
			password: "ValidPass123!",
			wantErr:  false,
		},
		{
			name: "failure: wrong password is rejected",
			setupData: func() string {
				user := createTestUser(t, ctx, queries,
					"user@example.com", "CorrectPass123!", "Test", "User")
				return user.ID.String()
			},
			email:    "user@example.com",
			password: "WrongPass123!",
			wantErr:  true,
			checkErr: func(t *testing.T, err error) {
				assert.Contains(t, err.Error(), "Invalid email or password")
				// Verify it's an Unauthorized error
				appErr, ok := err.(*errors.AppError)
				require.True(t, ok, "error should be AppError")
				assert.Equal(t, 401, appErr.Code)
			},
		},
		{
			name:      "failure: non-existent user is rejected",
			setupData: nil, // No user created
			email:     "ghost@example.com",
			password:  "AnyPass123!",
			wantErr:   true,
			checkErr: func(t *testing.T, err error) {
				assert.Contains(t, err.Error(), "Invalid email or password")
			},
		},
		{
			name:      "failure: empty email is rejected",
			setupData: nil,
			email:     "",
			password:  "Pass123!",
			wantErr:   true,
			checkErr: func(t *testing.T, err error) {
				assert.Contains(t, err.Error(), "Invalid email or password")
			},
		},
		{
			name: "failure: empty password is rejected",
			setupData: func() string {
				user := createTestUser(t, ctx, queries,
					"test@example.com", "CorrectPass123!", "Test", "User")
				return user.ID.String()
			},
			email:    "test@example.com",
			password: "",
			wantErr:  true,
			checkErr: func(t *testing.T, err error) {
				assert.Contains(t, err.Error(), "Invalid email or password")
			},
		},
		{
			name:      "security: SQL injection attempt in email is handled safely",
			setupData: nil,
			email:     "admin' OR '1'='1",
			password:  "Pass123!",
			wantErr:   true,
			checkErr: func(t *testing.T, err error) {
				assert.Contains(t, err.Error(), "Invalid email or password")
			},
		},
		{
			name: "security: case sensitivity is preserved in email",
			setupData: func() string {
				user := createTestUser(t, ctx, queries,
					"Test@Example.com", "Pass123!", "Test", "User")
				return user.ID.String()
			},
			email:    "test@example.com", // Different case
			password: "Pass123!",
			wantErr:  true, // Should fail because email case doesn't match
		},
	}

	// Execute all test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange: Set up test data
			var expectedUserID string
			if tt.setupData != nil {
				expectedUserID = tt.setupData()
			}

			// Act: Execute the authentication
			userID, err := authService.Authenticate(ctx, tt.email, tt.password)

			// Assert: Verify the results
			if tt.wantErr {
				assert.Error(t, err, "expected authentication to fail")
				assert.Empty(t, userID, "userID should be empty on failure")

				// Run custom error checks if provided
				if tt.checkErr != nil {
					tt.checkErr(t, err)
				}
			} else {
				assert.NoError(t, err, "expected authentication to succeed")
				assert.NotEmpty(t, userID, "userID should not be empty on success")
				assert.Equal(t, expectedUserID, userID, "userID should match created user")
			}
		})
	}
}

// ============================================================================
// TESTS - Password Hashing
// ============================================================================

func TestHashPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
		checkErr func(t *testing.T, err error)
	}{
		{
			name:     "success: normal password is hashed",
			password: "Password123!",
			wantErr:  false,
		},
		{
			name:     "success: long password is hashed",
			password: "ThisIsAVeryLongPasswordWithLotsOfCharacters123!@#$%^&*()",
			wantErr:  false,
		},
		{
			name:     "success: password with special characters",
			password: "P@ssw0rd!#$%^&*()",
			wantErr:  false,
		},
		{
			name:     "success: unicode characters in password",
			password: "Пароль123!",
			wantErr:  false,
		},
		{
			name:     "success: emoji in password",
			password: "Pass🔒word123!",
			wantErr:  false,
		},
		{
			name:     "success: very long password (72 chars - bcrypt limit)",
			password: "ThisIsExactlySeventyTwoCharactersLongPasswordToTestBcryptMaxLength!!",
			wantErr:  false,
		},
		{
			name:     "success: minimum length password",
			password: "a",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			hash, err := HashPassword(tt.password)

			// Assert
			if tt.wantErr {
				assert.Error(t, err, "expected password hashing to fail")
				assert.Empty(t, hash, "hash should be empty on error")

				if tt.checkErr != nil {
					tt.checkErr(t, err)
				}
			} else {
				assert.NoError(t, err, "expected password hashing to succeed")
				assert.NotEmpty(t, hash, "hash should not be empty")
				assert.NotEqual(t, tt.password, hash, "hash should differ from plaintext")

				// Verify hash is valid
				err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(tt.password))
				assert.NoError(t, err, "hash should validate against original password")

				// Verify bcrypt format
				assert.Contains(t, hash, "$2a$", "should be bcrypt format")

				// Verify salt uniqueness (same password produces different hashes)
				hash2, err := HashPassword(tt.password)
				assert.NoError(t, err)
				assert.NotEqual(t, hash, hash2, "same password should produce different hashes (salt is working)")
			}
		})
	}
}

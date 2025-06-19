package auth

import (
	"context"
	"log"
	"path/filepath"
	"testing"

	"github.com/iankencruz/threefive/internal/generated"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

func runMigrations(t *testing.T, db *pgxpool.Pool) {
	t.Helper()

	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			first_name TEXT NOT NULL,
			last_name TEXT NOT NULL,
			email text NOT NULL CHECK (TRIM(BOTH FROM email) <> ''::text) UNIQUE,
			password_hash TEXT NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
		);`,
		// Add other schema setup if needed
	}

	for _, q := range queries {
		_, err := db.Exec(context.Background(), q)
		require.NoError(t, err)
	}
}

func SetupPostgres(t *testing.T) error {

	ctx := context.Background()

	dbName := "users"
	dbUser := "user"
	dbPassword := "password"

	postgresContainer, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithInitScripts(filepath.Join("testdata", "init-user-db.sh")),
		postgres.WithConfigFile(filepath.Join("testdata", "my-postgres.conf")),
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		postgres.BasicWaitStrategies(),
	)
	defer func() {
		if err := testcontainers.TerminateContainer(postgresContainer); err != nil {
			log.Printf("failed to terminate container: %s", err)
		}
	}()
	if err != nil {
		log.Printf("failed to start container: %s", err)
		return err
	}

	return nil
}

func TestCreateUser(t *testing.T) {
	tests := []struct {
		name       string
		input      generated.CreateUserParams
		wantErr    bool
		assertFunc func(t *testing.T, user *generated.User)
	}{
		{
			name: "valid user",
			input: generated.CreateUserParams{
				FirstName:    "Alice",
				LastName:     "Test",
				Email:        "alice@example.com",
				PasswordHash: "password123",
			},
			wantErr: false,
			assertFunc: func(t *testing.T, user *generated.User) {
				assert.NotZero(t, user.ID)
				assert.Equal(t, "Alice", user.FirstName)
				assert.Equal(t, "alice@example.com", user.Email)
			},
		},
		{
			name: "missing email",
			input: generated.CreateUserParams{
				FirstName:    "Bob",
				LastName:     "Test",
				Email:        "",
				PasswordHash: "password123",
			},
			wantErr: true,
		},
		{
			name: "duplicate email",
			input: generated.CreateUserParams{
				FirstName:    "Charlie",
				LastName:     "Test",
				Email:        "duplicate@example.com",
				PasswordHash: "password123",
			},
			wantErr: false, // first insert
		},
		{
			name: "duplicate email again",
			input: generated.CreateUserParams{
				FirstName:    "Charlie2",
				LastName:     "Test",
				Email:        "duplicate@example.com",
				PasswordHash: "anotherpass",
			},
			wantErr: true, // second insert should fail
		},
	}

	for _, tt := range tests {
		tt := tt // capture range var
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()

			ctr, err := postgres.Run(
				ctx,
				"postgres:16-alpine",
				postgres.WithDatabase("users"),
				postgres.WithUsername("users"),
				postgres.WithPassword("password"),
				postgres.BasicWaitStrategies(),
				postgres.WithSQLDriver("pgx"),
			)
			testcontainers.CleanupContainer(t, ctr)
			require.NoError(t, err)

			connURI, err := ctr.ConnectionString(ctx, "sslmode=disable")
			assert.NoError(t, err)

			/// create a pgx connection pool
			db, err := pgxpool.New(ctx, connURI)
			assert.NoError(t, err)

			queries := generated.New(db) // ✅ Initialize sqlc Queries
			repo := NewAuthRepository(queries)

			runMigrations(t, db) // << Add this line

			// Handle duplicate email setup manually
			if tt.name == "duplicate email again" {
				_, _ = repo.CreateUser(ctx, generated.CreateUserParams{
					FirstName:    "Seed",
					LastName:     "User",
					Email:        "duplicate@example.com",
					PasswordHash: "seedpass",
				})
			}

			user, err := repo.CreateUser(ctx, tt.input)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				if tt.assertFunc != nil {
					tt.assertFunc(t, user)
				}
			}
		})
	}
}

func TestGetUserByEmail(t *testing.T) {
	ctx := context.Background()

	ctr, err := postgres.Run(
		ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("users"),
		postgres.WithUsername("users"),
		postgres.WithPassword("password"),
		postgres.BasicWaitStrategies(),
		postgres.WithSQLDriver("pgx"),
	)
	testcontainers.CleanupContainer(t, ctr)
	require.NoError(t, err)

	connURI, err := ctr.ConnectionString(ctx, "sslmode=disable")
	assert.NoError(t, err)

	/// create a pgx connection pool
	db, err := pgxpool.New(ctx, connURI)
	assert.NoError(t, err)

	testUser := generated.CreateUserParams{
		FirstName:    "Ian",
		LastName:     "Cruz",
		Email:        "iankencruz@gmail.com",
		PasswordHash: "secret-hash",
	}

	runMigrations(t, db) // << Add this line

	queries := generated.New(db) // ✅ Initialize sqlc Queries
	repo := NewAuthRepository(queries)

	created, err := repo.CreateUser(ctx, testUser)
	assert.NoError(t, err)

	tests := []struct {
		name      string
		email     string
		wantError bool
	}{
		{"existing_user", "iankencruz@gmail.com", false},
		{"nonexistent_user", "nobody@example.com", true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			user, err := repo.GetUserByEmail(ctx, tc.email)
			if tc.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, created.ID, user.ID)
			}
		})
	}
}

func TestGetUserByID(t *testing.T) {

	ctx := context.Background()

	ctr, err := postgres.Run(
		ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("users"),
		postgres.WithUsername("users"),
		postgres.WithPassword("password"),
		postgres.BasicWaitStrategies(),
		postgres.WithSQLDriver("pgx"),
	)
	testcontainers.CleanupContainer(t, ctr)
	require.NoError(t, err)

	connURI, err := ctr.ConnectionString(ctx, "sslmode=disable")
	assert.NoError(t, err)

	/// create a pgx connection pool
	db, err := pgxpool.New(ctx, connURI)
	assert.NoError(t, err)

	testUser := generated.CreateUserParams{
		FirstName:    "Ian",
		LastName:     "Cruz",
		Email:        "idtest@example.com",
		PasswordHash: "somehash",
	}

	runMigrations(t, db) // << Add this line

	queries := generated.New(db) // ✅ Initialize sqlc Queries
	repo := NewAuthRepository(queries)

	created, err := repo.CreateUser(ctx, testUser)
	assert.NoError(t, err)

	tests := []struct {
		name      string
		id        int32
		wantError bool
	}{
		{"existing_id", created.ID, false},
		{"missing_id", 99999, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			user, err := repo.GetUserByID(ctx, tc.id)
			if tc.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testUser.Email, user.Email)
			}
		})
	}
}

func TestDeleteUser(t *testing.T) {
	ctx := context.Background()

	ctr, err := postgres.Run(
		ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("users"),
		postgres.WithUsername("users"),
		postgres.WithPassword("password"),
		postgres.BasicWaitStrategies(),
		postgres.WithSQLDriver("pgx"),
	)
	testcontainers.CleanupContainer(t, ctr)
	require.NoError(t, err)

	connURI, err := ctr.ConnectionString(ctx, "sslmode=disable")
	assert.NoError(t, err)

	/// create a pgx connection pool
	db, err := pgxpool.New(ctx, connURI)
	assert.NoError(t, err)

	testUser := generated.CreateUserParams{
		FirstName:    "Delete",
		LastName:     "Me",
		Email:        "delete@example.com",
		PasswordHash: "hash",
	}

	queries := generated.New(db) // ✅ Initialize sqlc Queries
	repo := NewAuthRepository(queries)

	runMigrations(t, db) // << Add this line

	created, err := repo.CreateUser(ctx, testUser)
	assert.NoError(t, err)

	err = repo.DeleteUserByID(ctx, created.ID)
	// err = repo.DeleteUserByID(ctx, created.ID)
	assert.NoError(t, err)

	_, err = repo.GetUserByID(ctx, created.ID)
	assert.Error(t, err)
}

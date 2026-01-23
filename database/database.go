// database/database.go
package database

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

type Service interface {
	Health() map[string]string
	Close()
	Pool() *pgxpool.Pool
	RunMigrations(migrationsDir string) error
}

type service struct {
	db     *pgxpool.Pool
	logger *slog.Logger
}

func New(logger *slog.Logger) Service {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		logger.Error("DATABASE_URL environment variable is not set")
		// panic("DATABASE_URL is required")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Configure the database connection pool
	config, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		logger.Error("failed to parse database config", "error", err)
		panic(err)
	}

	config.MaxConns = 10
	config.MinConns = 2
	config.MaxConnIdleTime = 5 * time.Minute

	// Create the connection pool
	dbPool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		logger.Error("failed to create database connection pool", "error", err)
		panic(err)
	}

	// Verify the connection
	if err := dbPool.Ping(ctx); err != nil {
		logger.Error("failed to ping database", "error", err)
		panic(err)
	}

	logger.Info("connected to PostgreSQL successfully")
	return &service{db: dbPool, logger: logger}
}

func (s *service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err := s.db.Ping(ctx)
	if err != nil {
		return map[string]string{"status": "down", "error": err.Error()}
	}

	return map[string]string{"status": "up"}
}

func (s *service) Pool() *pgxpool.Pool {
	return s.db
}

func (s *service) Close() {
	s.logger.Info("closing PostgreSQL connection pool")
	s.db.Close()
}

// RunMigrations runs all pending migrations
func (s *service) RunMigrations(migrationsDir string) error {
	// Convert pgxpool to database/sql (goose requires *sql.DB)
	connConfig := s.db.Config().ConnConfig
	connStr := stdlib.RegisterConnConfig(connConfig)

	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return fmt.Errorf("failed to open sql connection: %w", err)
	}
	defer db.Close()

	// Set dialect
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set goose dialect: %w", err)
	}

	// Get current version
	current, err := goose.GetDBVersion(db)
	if err != nil {
		return fmt.Errorf("failed to get current DB version: %w", err)
	}
	s.logger.Info("current database version", "version", current)

	// Run migrations
	if err := goose.Up(db, migrationsDir); err != nil {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	// Get new version
	newVersion, err := goose.GetDBVersion(db)
	if err != nil {
		return fmt.Errorf("failed to get new DB version: %w", err)
	}

	if newVersion > current {
		s.logger.Info("database migrated",
			"from_version", current,
			"to_version", newVersion,
		)
	} else {
		s.logger.Info("database is up-to-date", "version", current)
	}

	return nil
}

package database

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

type Service interface {
	Health() map[string]string
	Close()
	Pool() *pgxpool.Pool
}

type service struct {
	db     *pgxpool.Pool
	logger *slog.Logger
}

func New(logger *slog.Logger) Service {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		logger.Info("DATABASE_URL environment variable is not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Configure the database connection pool
	config, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		logger.Error("failed to parse database config")
	}

	config.MaxConns = 10
	config.MinConns = 2
	config.MaxConnIdleTime = 5 * time.Minute

	// Create the connection pool
	dbPool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		logger.Error("failed to create database connection pool")
	}

	// Verify the connection
	if err := dbPool.Ping(ctx); err != nil {
		logger.Error("failed to ping database")
	}

	logger.Info("Connected to PostgreSQL successfully")

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
	log.Info().Msg("Closing PostgreSQL connection pool...")
	s.db.Close()
}

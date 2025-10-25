package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Connect establishes a connection pool to PostgreSQL using pgxpool
func Connect(databaseURL string) (*pgxpool.Pool, error) {
	// Parse the database URL and create pool config
	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database URL: %w", err)
	}

	// Configure connection pool settings
	config.MaxConns = 25
	config.MinConns = 5
	config.MaxConnLifetime = 5 * time.Minute
	config.MaxConnIdleTime = 1 * time.Minute
	config.HealthCheckPeriod = 1 * time.Minute

	// Configure individual connection settings
	config.ConnConfig.ConnectTimeout = 30 * time.Second

	// Optional: Add connection logging in development
	if isDevelopment() {
		config.ConnConfig.Tracer = &queryTracer{}
	}

	// Create connection pool
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Test the connection
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Printf("✅ Database connected (pool: %d-%d connections)", config.MinConns, config.MaxConns)
	return pool, nil
}

// isDevelopment checks if we're in development mode
func isDevelopment() bool {
	// You could also check ENV environment variable here
	return true // For now, always enable in this basic setup
}

// Optional: Query tracer for development logging
type queryTracer struct{}

func (t *queryTracer) TraceQueryStart(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	// log.Printf("🔍 SQL: %s", data.SQL)
	// if len(data.Args) > 0 {
	// 	log.Printf("🔍 Args: %v", data.Args)
	// }
	return ctx
}

func (t *queryTracer) TraceQueryEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryEndData) {
	if data.Err != nil {
		log.Printf("❌ SQL Error: %v", data.Err)
	}
}

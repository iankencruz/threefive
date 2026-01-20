// internal/session/store.go
package session

import (
	"context"
	"log/slog"
	"time"

	"github.com/iankencruz/threefive/database/generated"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresStore struct {
	db      *pgxpool.Pool
	queries *generated.Queries
	logger  *slog.Logger
}

func NewPostgresStore(db *pgxpool.Pool, queries *generated.Queries, logger *slog.Logger) *PostgresStore {
	return &PostgresStore{
		db:      db,
		queries: queries,
		logger:  logger.With("component", "session_store"),
	}
}

// Find returns the data for a given session token from the PostgresStore instance.
// If the session token is not found or is expired, the returned exists flag will be false.
func (p *PostgresStore) Find(ctx context.Context, token string) (b []byte, exists bool, err error) {
	data, err := p.queries.GetSession(ctx, token)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, false, nil
		}
		p.logger.Error("failed to find session",
			"error", err,
			"token", token[:10]+"...",
		)
		return nil, false, err
	}

	return data, true, nil
}

// Commit adds a session token and data to the PostgresStore instance with the given expiry time.
// If the session token already exists, then the data and expiry time are updated.
func (p *PostgresStore) Commit(ctx context.Context, token string, b []byte, expiry time.Time) error {
	err := p.queries.CommitSession(ctx, generated.CommitSessionParams{
		Token:  token,
		Data:   b,
		Expiry: expiry,
	})
	if err != nil {
		p.logger.Error("failed to commit session",
			"error", err,
			"token", token[:10]+"...",
		)
		return err
	}

	p.logger.Debug("session committed",
		"token", token[:10]+"...",
		"expiry", expiry,
	)

	return nil
}

// Delete removes a session token and corresponding data from the PostgresStore instance.
func (p *PostgresStore) Delete(ctx context.Context, token string) error {
	err := p.queries.DeleteSession(ctx, token)
	if err != nil {
		p.logger.Error("failed to delete session",
			"error", err,
			"token", token[:10]+"...",
		)
		return err
	}

	p.logger.Debug("session deleted",
		"token", token[:10]+"...",
	)

	return nil
}

// All returns a map containing the token and data for all active (i.e. not expired) sessions.
func (p *PostgresStore) All(ctx context.Context) (map[string][]byte, error) {
	// This is typically not needed for most applications
	// We can implement it if needed later
	return nil, nil
}

// Cleanup deletes expired session data from the PostgresStore instance.
func (p *PostgresStore) Cleanup(ctx context.Context) error {
	err := p.queries.DeleteExpiredSessions(ctx)
	if err != nil {
		p.logger.Error("failed to cleanup expired sessions",
			"error", err,
		)
		return err
	}

	p.logger.Debug("cleaned up expired sessions")
	return nil
}

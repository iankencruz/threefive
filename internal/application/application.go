package application

import (
	"context"
	"log/slog"

	"github.com/alexedwards/scs/v2"
	models "github.com/iankencruz/threefive/internal/models/users"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Application struct {
	Logger         *slog.Logger
	DB             *pgxpool.Pool
	SessionManager *scs.SessionManager
	Users          models.Repository
}

func New(ctx context.Context, db *pgxpool.Pool, sm *scs.SessionManager, logger *slog.Logger) *Application {
	return &Application{
		Logger:         logger,
		DB:             db,
		SessionManager: sm,
		Users:          models.NewUser(db),
	}
}

package application

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/iankencruz/threefive/backend/internal/auth"
	"github.com/iankencruz/threefive/backend/internal/core/sessions"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Application struct {
	Config         *Config
	Logger         *slog.Logger
	DB             *pgxpool.Pool
	SessionManager *sessions.Manager
	AuthHandler    *auth.Handler
}

func New(ctx context.Context, cfg *Config, db *pgxpool.Pool, sm *sessions.Manager, logger *slog.Logger) *Application {

	authRepo := auth.NewPgxRepository(db)
	authService := auth.NewService(authRepo)
	authHandler := &auth.Handler{
		Service:        authService,
		SessionManager: sm,
		Logger:         logger,
	}

	return &Application{
		Config:         cfg,
		Logger:         logger,
		DB:             db,
		SessionManager: sm,
		AuthHandler:    authHandler,
	}
}

func (app *Application) GetUserID(r *http.Request) (int32, error) {
	id, err := app.SessionManager.GetUserID(r)
	if err != nil || id == 0 {
		return 0, errors.New("no session user")
	}
	return id, nil
}

func (app *Application) LoadUser(ctx context.Context, userID int) (any, error) {
	return app.AuthHandler.Service.GetUserByID(ctx, int32(userID))
}

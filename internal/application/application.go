package application

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/alexedwards/scs/v2"
	models "github.com/iankencruz/threefive/internal/models/users"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Application struct {
	Logger         *slog.Logger
	DB             *pgxpool.Pool
	SessionManager *scs.SessionManager
	UserRepo       models.Repository
}

func New(ctx context.Context, db *pgxpool.Pool, sm *scs.SessionManager, logger *slog.Logger) *Application {
	return &Application{
		Logger:         logger,
		DB:             db,
		SessionManager: sm,
		UserRepo:       models.NewUser(db),
	}
}

func (app *Application) GetUserID(r *http.Request) (int, error) {
	id := app.SessionManager.GetInt(r.Context(), "userID")
	if id == 0 {
		return 0, errors.New("no session user")
	}
	return id, nil
}

func (app *Application) LoadUser(ctx context.Context, userID int) (any, error) {
	return app.UserRepo.GetUserByID(ctx, int32(userID))
}

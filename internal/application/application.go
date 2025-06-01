package application

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/a-h/templ"
	"github.com/iankencruz/threefive/internal/models/users"
	"github.com/iankencruz/threefive/internal/templates"
	"github.com/iankencruz/threefive/ui/templates/layouts"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Application struct {
	Logger *slog.Logger
	DB     *pgxpool.Pool
	Users  users.Repository
}

func New(ctx context.Context, db *pgxpool.Pool, logger *slog.Logger) *Application {
	return &Application{
		Logger: logger,
		DB:     db,
		Users:  users.NewRepository(db),
	}
}

func (app *Application) Render(
	w http.ResponseWriter,
	r *http.Request,
	title string,
	content templ.Component,
) {
	templates.RenderTempl(w, r, layouts.BaseLayout, title, content)
}

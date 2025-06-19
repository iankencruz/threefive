package application

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/iankencruz/threefive/internal/auth"
	"github.com/iankencruz/threefive/internal/core/s3"
	"github.com/iankencruz/threefive/internal/core/sessions"
	"github.com/iankencruz/threefive/internal/generated"
	"github.com/iankencruz/threefive/internal/media"
	project "github.com/iankencruz/threefive/internal/projects"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Application struct {
	Config         *Config
	Logger         *slog.Logger
	DB             *pgxpool.Pool
	SessionManager *sessions.Manager
	AuthHandler    *auth.Handler
	MediaHandler   *media.Handler
	ProjectHandler *project.Handler
}

func New(
	ctx context.Context,
	cfg *Config,
	db *pgxpool.Pool,
	sm *sessions.Manager,
	logger *slog.Logger) *Application {

	queries := generated.New(db) // ✅ Initialize sqlc Queries

	// ✅ Initialise S3 uploader here
	uploader, err := s3.NewUploader(
		cfg.S3.Endpoint,
		cfg.S3.AccessKey,
		cfg.S3.SecretKey,
		cfg.S3.Bucket,
		cfg.S3.UseSSL,
		cfg.S3.BaseURL,
	)

	if err != nil {
		logger.Error("failed to initialise S3", "err", err)
		panic(err) // or return error if you propagate
	}

	authHandler := auth.NewHandler(queries, sm, logger)
	mediaHandler := media.NewHandler(queries, logger, uploader)
	projectHandler := project.NewHandler(queries, logger)

	return &Application{
		Config:         cfg,
		Logger:         logger,
		DB:             db,
		SessionManager: sm,
		AuthHandler:    authHandler,
		MediaHandler:   mediaHandler,
		ProjectHandler: projectHandler,
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

package application

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/iankencruz/threefive/internal/auth"
	"github.com/iankencruz/threefive/internal/blocks"
	"github.com/iankencruz/threefive/internal/core/s3"
	"github.com/iankencruz/threefive/internal/core/sessions"
	"github.com/iankencruz/threefive/internal/generated"
	"github.com/iankencruz/threefive/internal/media"
	"github.com/iankencruz/threefive/internal/pages"
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
	PageHandler    *pages.Handler
	BlockHandler   *blocks.Handler
}

func New(
	ctx context.Context,
	cfg *Config,
	db *pgxpool.Pool,
	sm *sessions.Manager,
	logger *slog.Logger,
) *Application {
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

	blockRepo := blocks.NewRepository(queries)
	blockService := blocks.NewService(blockRepo)

	authHandler := auth.NewHandler(queries, sm, logger)
	mediaHandler := media.NewHandler(queries, logger, uploader)
	projectHandler := project.NewHandler(queries, logger)
	pageHandler := *pages.NewHandler(queries, blockRepo, blockService, logger)
	blockHandler := *blocks.NewHandler(queries, logger)

	return &Application{
		Config:         cfg,
		Logger:         logger,
		DB:             db,
		SessionManager: sm,
		AuthHandler:    authHandler,
		MediaHandler:   mediaHandler,
		ProjectHandler: projectHandler,
		PageHandler:    &pageHandler,
		BlockHandler:   &blockHandler,
	}
}

func (app *Application) GetUserID(r *http.Request) (uuid.UUID, error) {
	id, err := app.SessionManager.GetUserID(r)
	if err != nil || id == uuid.Nil {
		return uuid.Nil, errors.New("no session user")
	}
	return id, nil
}

func (app *Application) LoadUser(ctx context.Context, userID uuid.UUID) (any, error) {
	return app.AuthHandler.Service.GetUserByID(ctx, userID)
}

func (app *Application) EnsureAdminExists() error {
	var (
		fname    = "Admin"
		lname    = "User"
		email    = os.Getenv("ADMIN_EMAIL")
		password = os.Getenv("ADMIN_PASSWORD")
	)

	users, err := app.AuthHandler.Repo.ListUsers(context.Background())
	if err != nil {
		fmt.Printf("Error fetching Users: %v", err)
		return err
	}

	if len(users) == 0 {
		fmt.Printf("No Users found. Creating a default admin user...\n")
		user, err := app.AuthHandler.Service.Register(context.Background(), fname, lname, email, password)
		if err != nil {
			fmt.Printf("Error creating a user: %v\n", err)
			return err
		}

		fmt.Printf("Default Admin user created: %v.\n", user.Email)

	} else {
		fmt.Printf("✅ %d user(s) already exist. Skipping admin bootstrap.", len(users))
	}

	return nil
}

package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/iankencruz/threefive/backend/internal/application"
	"github.com/iankencruz/threefive/backend/internal/core/logger"
	"github.com/iankencruz/threefive/backend/internal/core/sessions"
	"github.com/iankencruz/threefive/backend/internal/database"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	ctx := context.Background()

	cfg := application.LoadConfig()
	// Logger setup
	log := logger.New()
	log.Info("💾 Starting ThreeFiveProject")

	// Database connection
	db, err := database.Connect(ctx, log, cfg.DB_URL)
	if err != nil {
		log.Error("❌ DB connect error", slog.String("Reason", err.Error()))
		os.Exit(1)
	}
	defer db.Close()

	// create session manager
	// Initialize the Session Manager
	sm := sessions.NewManager(db)

	app := application.New(
		ctx,
		cfg,
		db,
		sm,
		log,
	)

	// requireAuth := middleware.RequireAuth(app)

	r := application.Routes(app)

	log.Info("🚀 Server running at http://localhost:8080")
	srv := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  90 * time.Second,
	}
	if err := srv.ListenAndServe(); err != nil {
		log.Error("❌ Server error:", slog.Any("error", err))
	}

}

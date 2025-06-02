package main

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/alexedwards/scs/pgxstore"
	"github.com/alexedwards/scs/v2"
	"github.com/iankencruz/threefive/internal/application"
	"github.com/iankencruz/threefive/internal/database"
	"github.com/iankencruz/threefive/internal/logger"
	"github.com/iankencruz/threefive/internal/routes"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	ctx := context.Background()

	// Logger setup
	log := logger.New()
	log.Info("💾 Starting ThreeFiveProject")

	// Database connection
	db := database.Connect(ctx, log)
	defer db.Close()

	// create session manager
	sessionManager := scs.New()
	sessionManager.Store = pgxstore.New(db)
	sessionManager.Lifetime = 24 * time.Hour
	sessionManager.Cookie.SameSite = http.SameSiteLaxMode
	sessionManager.Cookie.Secure = false // Set to true if using HTTPS

	app := application.New(
		ctx,
		db,
		sessionManager,
		log,
	)

	// requireAuth := middleware.RequireAuth(app)

	r := routes.Routes(app)

	log.Info("🚀 Server running at http://localhost:8080")
	srv := &http.Server{
		Addr:         ":8080",
		Handler:      sessionManager.LoadAndSave(r),
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  90 * time.Second,
	}
	if err := srv.ListenAndServe(); err != nil {
		log.Error("❌ Server error:", slog.Any("error", err))
	}

}

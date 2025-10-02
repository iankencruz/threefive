package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/iankencruz/threefive/internal/auth"
	"github.com/iankencruz/threefive/internal/media"
	"github.com/iankencruz/threefive/internal/shared/middleware"
)

func (s *Server) setupRouter() http.Handler {
	r := chi.NewRouter()

	// Global middleware
	r.Use(middleware.RequestLogger)
	r.Use(middleware.Recovery)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{s.config.Frontend.URL},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "Cookie"},
		ExposedHeaders:   []string{"Set-Cookie"},
		MaxAge:           300,
		AllowCredentials: true,
	}))

	// Mount auth routes
	auth.RegisterRoutes(r, s.authHandler, s.sessionManager)

	// Basic routes
	r.Get("/", homeHandler)
	r.Get("/health", healthHandler(s.db))

	// API v1 routes
	r.Route("/api/v1", func(r chi.Router) {
		// Auth middleware
		r.Use(auth.NewMiddleware(s.sessionManager).RequireAuth)

		// Mount feature routes
		media.RegisterRoutes(r, s.mediaHandler)
		// user.RegisterRoutes(r, s.userHandler)
		// project.RegisterRoutes(r, s.projectHandler)

		r.Get("/status", statusHandler)
		r.Get("/db-test", dbTestHandler)
	})

	return r
}

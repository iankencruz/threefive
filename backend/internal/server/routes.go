package server

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/iankencruz/threefive/internal/auth"
	"github.com/iankencruz/threefive/internal/media"
	"github.com/iankencruz/threefive/internal/pages"
	"github.com/iankencruz/threefive/internal/shared/middleware"
)

func (s *Server) setupRouter() http.Handler {
	r := chi.NewRouter()

	// Global middleware
	r.Use(middleware.RequestLogger)
	r.Use(middleware.Recovery)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{s.Config.Frontend.URL},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "Cookie"},
		ExposedHeaders:   []string{"Set-Cookie"},
		MaxAge:           300,
		AllowCredentials: true,
	}))

	workDir, _ := os.Getwd()
	uploadsDir := http.Dir(filepath.Join(workDir, s.Config.Storage.LocalBasePath))
	FileServer(r, "/uploads", uploadsDir)

	// Mount auth routes
	auth.RegisterRoutes(r, s.AuthHandler, s.SessionManager)

	// Basic routes
	r.Get("/", homeHandler)
	r.Get("/health", healthHandler(s.DB))

	// API v1 routes
	r.Route("/api/v1", func(r chi.Router) {
		// Auth middleware
		r.Use(auth.NewMiddleware(s.SessionManager).RequireAuth)

		// Mount feature routes
		media.RegisterRoutes(r, s.MediaHandler)
		pages.RegisterRoutes(r, s.PageHandler, s.SessionManager)
		// user.RegisterRoutes(r, s.userHandler)
		// project.RegisterRoutes(r, s.projectHandler)

		r.Get("/status", statusHandler)
		r.Get("/db-test", dbTestHandler)
	})

	return r
}

// FileServer sets up a http.FileServer handler to serve static files from a http.FileSystem
func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", http.StatusMovedPermanently).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}

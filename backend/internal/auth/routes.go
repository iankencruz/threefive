package auth

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Routes creates and returns auth routes
func Routes(handler *Handler, middleware *Middleware) http.Handler {
	r := chi.NewRouter()

	// Public auth routes (no authentication required)
	r.Group(func(r chi.Router) {
		r.Use(middleware.RequireNoAuth) // Prevent already authenticated users

		r.Post("/register", handler.Register)
		r.Post("/login", handler.Login)
		r.Post("/request-password-reset", handler.RequestPasswordReset)
		r.Post("/reset-password", handler.ResetPassword)
	})

	// Protected auth routes (authentication required)
	r.Group(func(r chi.Router) {
		r.Use(middleware.RequireAuth)

		r.Get("/me", handler.Me)
		r.Post("/logout", handler.Logout)
		r.Post("/logout-all", handler.LogoutAll)
		r.Put("/change-password", handler.ChangePassword)
		r.Post("/refresh", handler.RefreshSession)

		// Session management routes
		r.Get("/sessions", handler.GetSessions)
		r.Delete("/sessions/{sessionID}", handler.RevokeSession)
	})

	return r
}

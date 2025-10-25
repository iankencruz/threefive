package auth

import (
	"github.com/go-chi/chi/v5"
	"github.com/iankencruz/threefive/internal/shared/session"
)

// RegisterRoutes creates and returns auth routes
func RegisterRoutes(r chi.Router, handler *Handler, sessionManager *session.Manager) {
	middleware := NewMiddleware(sessionManager)

	r.Route("/auth", func(r chi.Router) {
		// Public routes
		r.Group(func(r chi.Router) {
			r.Post("/register", handler.Register)
			r.Post("/login", handler.Login)
			r.Post("/request-password-reset", handler.RequestPasswordReset)
			r.Post("/reset-password", handler.ResetPassword)
		})

		// Protected routes
		r.Group(func(r chi.Router) {
			r.Use(middleware.RequireAuth)

			r.Get("/me", handler.Me)
			r.Post("/logout", handler.Logout)
			r.Post("/logout-all", handler.LogoutAll)
			r.Put("/change-password", handler.ChangePassword)
			r.Post("/refresh", handler.RefreshSession)

			r.Get("/sessions", handler.GetSessions)
			r.Delete("/sessions/{sessionID}", handler.RevokeSession)
		})
	})
}

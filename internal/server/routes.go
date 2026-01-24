package server

import (
	"github.com/iankencruz/threefive/internal/handler"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

// RegisterRoutes defines the API endpoints
func (s *Server) RegisterRoutes() {
	s.Log.Info("Registering routes...")

	// Initialize handlers
	authHandler := handler.NewAuthHandler(s.AuthService, s.SessionManager, s.Log)
	adminHandler := handler.NewAdminHandler(s.Log)
	mediaHandler := handler.NewMediaHandler(s.MediaService, s.Log)

	// Static assets
	s.Echo.Static("/assets", "assets")

	// Session Middleware
	s.Echo.Use(s.SessionMiddleware.Session)

	s.Echo.Pre(middleware.RemoveTrailingSlash())

	// handlers
	s.Echo.GET("/health", s.healthCheckHandler)

	// Public routes (no auth required)
	s.Echo.GET("/login", authHandler.ShowLoginPage)
	s.Echo.POST("/login", authHandler.HandleLogin)
	s.Echo.POST("/logout", authHandler.HandleLogout)

	// admin routes (require authentication)
	admin := s.Echo.Group("/admin")
	admin.Use(s.SessionMiddleware.RequireAuth)

	// redirect admin to dashboard
	admin.GET("", func(c *echo.Context) error {
		return c.Redirect(302, "/admin/dashboard")
	})

	admin.GET("/dashboard", adminHandler.ShowDashboard)

	admin.GET("/projects", adminHandler.ShowProjects)

	media := admin.Group("/media")
	media.GET("", mediaHandler.ShowMediaList)
	media.POST("/upload", mediaHandler.UploadMedia)
	media.DELETE("/:mediaID", mediaHandler.DeleteMedia)

	s.Log.Info("routes registered successfully")
}

// Simple handler methods on the server for basic routes
func (s *Server) healthCheckHandler(c *echo.Context) error {
	return c.JSON(200, map[string]string{
		"status": "ok",
	})
}

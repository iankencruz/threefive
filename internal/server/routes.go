package server

import (
	"github.com/iankencruz/threefive/internal/handler"
	"github.com/labstack/echo/v5"
)

// RegisterRoutes defines the API endpoints
func (s *Server) RegisterRoutes() {
	s.Log.Info("Registering routes...")

	// Initialize handlers
	authHandler := handler.NewAuthHandler(s.AuthService, s.SessionManager, s.Log)

	// Static assets
	s.Echo.Static("/assets", "assets")

	// Session Middleware
	s.Echo.Use(s.SessionMiddleware.Session)

	// handlers
	// s.Echo.GET("/", h.HealthCheckHandler)
	// s.Echo.GET("/hello", h.HelloWorldHandler)

	// Public routes (no auth required)
	s.Echo.GET("/login", authHandler.ShowLoginPage)
	s.Echo.POST("/login", authHandler.HandleLogin)
	s.Echo.POST("/logout", authHandler.HandleLogout)

	// Protected routes (require authentication)
	protected := s.Echo.Group("")
	protected.Use(s.SessionMiddleware.RequireAuth)

	protected.GET("/", s.healthCheckHandler)
	protected.GET("/hello", s.helloWorldHandler)

	s.Log.Info("routes registered successfully")
}

// Simple handler methods on the server for basic routes
func (s *Server) healthCheckHandler(c *echo.Context) error {
	return c.JSON(200, map[string]string{
		"status": "ok",
	})
}

func (s *Server) helloWorldHandler(c *echo.Context) error {
	return c.String(200, "Hello, World!\n")
}

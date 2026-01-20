package server

import (
	"github.com/iankencruz/threefive/internal/handler"
)

// RegisterRoutes defines the API endpoints
func (s *Server) RegisterRoutes() {
	h := handler.NewAuthHandler()

	s.Echo.Static("/assets", "assets")

	// handlers
	// s.Echo.GET("/", h.HealthCheckHandler)
	// s.Echo.GET("/hello", h.HelloWorldHandler)

	// Initialize public auth handler

	// Login routes
	s.Echo.GET("/login", h.ShowLoginPage)
	s.Echo.POST("/login", h.HandleLogin)

	// Placeholder admin dashboard
	// s.Echo.GET("/admin/dashboard", h.ShowDashboard)
	// s.Echo.GET("/admin", publicAuthHandler.ShowDashboard)
}

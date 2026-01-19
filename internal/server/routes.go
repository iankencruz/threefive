package server

import handlers "github.com/iankencruz/threefive/internal/handler/admin"

// RegisterRoutes defines the API endpoints
func (s *Server) RegisterRoutes() {
	h := handlers.NewHandler()

	s.echo.GET("/", h.HealthCheckHandler)
	s.echo.GET("/hello", h.HelloWorldHandler)
}

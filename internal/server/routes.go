package server

import (
	"os"
	"syscall"
	"time"

	handlers "github.com/iankencruz/threefive/internal/handler/admin"
	"github.com/labstack/echo/v5"
	"github.com/rs/zerolog/log"
)

// RegisterRoutes defines the API endpoints
func (s *Server) RegisterRoutes() {
	h := handlers.NewHandler()

	s.echo.GET("/", h.HealthCheckHandler)
	s.echo.GET("/hello", h.HelloWorldHandler)

	s.echo.GET("/slow", func(c *echo.Context) error {
		log.Info().Msg("Slow request started (7s wait)...")
		time.Sleep(7 * time.Second)
		log.Info().Msg("Slow request finished!")
		return c.String(200, "Slow response complete\n")
	})

	s.echo.GET("/shutdown", func(c *echo.Context) error {
		log.Warn().Msg("Shutdown route hit! Sending SIGTERM...")
		go func() {
			time.Sleep(500 * time.Millisecond)
			p, _ := os.FindProcess(os.Getpid())
			p.Signal(syscall.SIGTERM)
		}()
		return c.String(200, "Shutdown signal sent.\n")
	})
}

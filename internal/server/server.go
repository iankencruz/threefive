package server

import (
	"io"
	"log/slog"

	"github.com/iankencruz/threefive/internal/middleware"
	"github.com/labstack/echo/v5"
)

type Server struct {
	echo *echo.Echo
}

func New() *Server {
	// Initialize Echo with the config that uses the silent logger
	e := echo.NewWithConfig(echo.Config{
		Logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
	})

	e.Use(middleware.CustomRequestLogger())

	s := &Server{
		echo: e,
	}

	s.RegisterRoutes()

	return s

}

// Start runs the Echo server on a specific address
func (s *Server) Start(address string) error {
	return s.echo.Start(address)
}

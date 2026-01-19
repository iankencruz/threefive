package server

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/iankencruz/threefive/database"
	"github.com/iankencruz/threefive/internal/middleware"
	"github.com/iankencruz/threefive/internal/services"
	"github.com/labstack/echo/v5"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Server struct {
	Echo        *echo.Echo
	DB          database.Service
	UserService *services.UserService
}

func NewServer() *Server {

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	if os.Getenv("ENV") != "production" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	}

	// Initialize Echo with the config that uses the silent logger
	e := echo.NewWithConfig(echo.Config{
		Logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
	})

	e.Use(middleware.CustomRequestLogger())

	s := &Server{
		Echo: e,
		DB:   database.New(),
	}

	s.RegisterRoutes()

	return s

}

// Start runs the Echo server on a specific address
func (s *Server) Start(ctx context.Context, address string) error {

	if address == "" {
		address = ":8080" // fallback
	}

	srv := &http.Server{
		Addr:    address,
		Handler: s.Echo,
	}

	// Start server in background
	go func() {
		log.Info().Msgf("Starting server on %s", address)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error().Err(err).Msg("Server ListenAndServe failed")
		}
	}()

	// Wait for the signal context to be done
	<-ctx.Done()
	log.Warn().Msg("Signal received, beginning graceful shutdown...")

	// Create a deadline for the shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Error().Err(err).Msg("Graceful shutdown failed")
		return err
	}

	// Only close database pool after server has shutdown
	s.DB.Close()

	log.Info().Msg("Server exited cleanly")
	return nil
}

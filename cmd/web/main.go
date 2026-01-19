package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/iankencruz/threefive/internal/server"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

func main() {

	// Load environment variables from .env file if it exists
	err := godotenv.Load()
	if err != nil {
		log.Info().Msg("No .env file found, proceeding with environment variables")
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Create server instance
	s := server.NewServer()

	// Get port from environment
	port := fmt.Sprintf(":%s", os.Getenv("PORT"))

	// Start the server
	if err := s.Start(ctx, port); err != nil {
		log.Fatal().Err(err).Msg("Failed to start server")
	}

}

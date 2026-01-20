package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/iankencruz/threefive/internal/server"
	"github.com/joho/godotenv"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Create server instance
	s := server.NewServer()

	// Load environment variables from .env file if it exists
	err := godotenv.Load()
	if err != nil {
		s.Log.Info("No .env file found, proceeding with environment variables")
	}

	// Get port from environment
	port := fmt.Sprintf(":%s", os.Getenv("PORT"))

	// Start the server
	if err := s.Start(ctx, port); err == nil {
		s.Log.Error("Failed to start server", err)
	}
}

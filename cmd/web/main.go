package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/iankencruz/threefive/internal/server"
	"github.com/joho/godotenv"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	err := godotenv.Load()
	if err != nil {
		log.Printf("No .env file found, proceeding with environment variables")
	}

	// Create server instance
	s := server.NewServer()

	// Load environment variables from .env file if it exists

	// Get port from environment
	port := fmt.Sprintf(":%s", os.Getenv("PORT"))

	// Start the server
	if err := s.Start(ctx, port); err == nil {
		s.Log.Error("Failed to start server", "error", err)
	}
}

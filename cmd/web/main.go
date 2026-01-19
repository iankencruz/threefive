package main

import (
	"fmt"
	"os"
	"time"

	"github.com/iankencruz/threefive/internal/server"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Info().Msg("No .env file found, proceeding with environment variables")
	}

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	if os.Getenv("ENV") != "production" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	}

	s := server.New()

	// Get port from environment
	port := fmt.Sprintf(":%s", os.Getenv("PORT"))
	if port == "" {
		port = ":8080" // fallback
	}

	log.Info().Msgf("Starting %s server on port %s", os.Getenv("APP_NAME"), port)

	if err := s.Start(port); err != nil {

	}

}

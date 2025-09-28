package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Frontend FrontendConfig
}

type ServerConfig struct {
	Host string
	Port int
	Env  string
}

type DatabaseConfig struct {
	URL string
}

type FrontendConfig struct {
	URL string
}

// Load reads configuration from environment variables
func Load() *Config {
	// Load .env file - ignore error if file doesn't exist
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	cfg := &Config{
		Server: ServerConfig{
			Host: getEnv("SERVER_HOST", "localhost"),
			Port: getEnvAsInt("PORT", 8000),
			Env:  getEnv("ENV", "development"),
		},
		Database: DatabaseConfig{
			URL: getEnvRequired("DATABASE_URL"),
		},
		Frontend: FrontendConfig{
			URL: getEnvRequired("FRONTEND_URL"),
		},
	}

	// Validate required configuration
	if cfg.Database.URL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	// Log configuration in development
	if cfg.IsDevelopment() {
		log.Printf("🔧 Server: %s:%d (%s)", cfg.Server.Host, cfg.Server.Port, cfg.Server.Env)
		log.Printf("🔧 Database: Connected")
	}

	return cfg
}

// Helper functions
func getEnv(key, defaultVal string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultVal
}

func getEnvRequired(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Required environment variable %s is not set", key)
	}
	return value
}

func getEnvAsInt(key string, defaultVal int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
		log.Printf("Warning: Invalid integer for %s, using default: %d", key, defaultVal)
	}
	return defaultVal
}

// ServerAddress returns the full server address
func (c *Config) ServerAddress() string {
	return c.Server.Host + ":" + strconv.Itoa(c.Server.Port)
}

// IsDevelopment returns true if running in development mode
func (c *Config) IsDevelopment() bool {
	return c.Server.Env == "development"
}

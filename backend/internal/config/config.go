package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Frontend FrontendConfig
	Storage  StorageConfig
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

type StorageConfig struct {
	Type          string
	LocalBasePath string
	LocalBaseURL  string
	S3Bucket      string
	S3Region      string
	S3AccessKey   string
	S3SecretKey   string
	S3Endpoint    string
	S3PublicURL   string
	S3UseSSL      bool
}

// Load reads configuration from environment variables
func Load() *Config {
	// Load .env file - ignore error if file doesn't exist
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	// Determine storage type from env
	storageType := getEnv("STORAGE_TYPE", "local")

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
		Storage: StorageConfig{
			Type:          storageType,
			LocalBasePath: getEnv("LOCAL_STORAGE_PATH", "./uploads"),
			LocalBaseURL:  getEnv("LOCAL_STORAGE_URL", "http://localhost:8000/uploads"),
			S3Bucket:      getEnv("S3_BUCKET", ""),
			S3Region:      getEnv("S3_REGION", ""),
			S3AccessKey:   getEnv("S3_ACCESS_KEY", ""),
			S3SecretKey:   getEnv("S3_SECRET_KEY", ""),
			S3Endpoint:    getEnv("S3_ENDPOINT", ""),
			S3PublicURL:   getEnv("S3_PUBLIC_URL", ""),
			S3UseSSL:      getEnvAsBool("S3_USE_SSL", true),
		},
	}

	// Validate required configuration
	if cfg.Database.URL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	if cfg.Frontend.URL == "" {
		log.Fatal("FRONTEND_URL is required")
	}

	// Validate storage configuration
	if err := cfg.validateStorage(); err != nil {
		log.Fatalf("Storage configuration error: %v", err)
	}

	// Log configuration in development
	if cfg.IsDevelopment() {
		log.Printf("🔧 Server: %s:%d (%s)", cfg.Server.Host, cfg.Server.Port, cfg.Server.Env)
		log.Printf("🔧 Database: Connected")
	}

	return cfg
}

// validateStorage validates storage configuration based on type
func (c *Config) validateStorage() error {
	switch c.Storage.Type {
	case "local":
		// Local storage is always valid (will create dir if needed)
		return nil
	case "s3":
		// Validate S3 configuration
		if c.Storage.S3Bucket == "" {
			return fmt.Errorf("S3_BUCKET is required when STORAGE_TYPE=s3")
		}
		if c.Storage.S3Region == "" {
			return fmt.Errorf("S3_REGION is required when STORAGE_TYPE=s3")
		}
		if c.Storage.S3AccessKey == "" {
			return fmt.Errorf("S3_ACCESS_KEY is required when STORAGE_TYPE=s3")
		}
		if c.Storage.S3SecretKey == "" {
			return fmt.Errorf("S3_SECRET_KEY is required when STORAGE_TYPE=s3")
		}
		return nil
	default:
		return fmt.Errorf("invalid STORAGE_TYPE: %s (must be 'local' or 's3')", c.Storage.Type)
	}
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

func getEnvAsBool(key string, defaultVal bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolVal, err := strconv.ParseBool(value); err == nil {
			return boolVal
		}
		log.Printf("Warning: Invalid boolean for %s, using default: %t", key, defaultVal)
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

// IsProduction returns true if running in production mode
func (c *Config) IsProduction() bool {
	return c.Server.Env == "production"
}

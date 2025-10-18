// backend/internal/config/config.go
package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

type Config struct {
	// Server
	Server                 ServerConfig
	Database               DatabaseConfig
	Frontend               FrontendConfig
	Storage                StorageConfig
	SessionSecret          string
	AllowedOrigins         []string
	MediaUploadPath        string
	MediaMaxSize           int64 // in bytes
	AutoPurgeEnabled       bool
	AutoPurgeRetentionDays int
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

func Load() (*Config, error) {
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
			Type:          getEnv("STORAGE_TYPE", "local"),
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
		SessionSecret:          getEnv("SESSION_SECRET", ""),
		MediaUploadPath:        getEnv("MEDIA_UPLOAD_PATH", "./uploads"),
		MediaMaxSize:           getEnvAsInt64("MEDIA_MAX_SIZE", 10*1024*1024), // 10MB default
		AutoPurgeEnabled:       getEnvAsBool("AUTO_PURGE_ENABLED", true),
		AutoPurgeRetentionDays: getEnvAsInt("AUTO_PURGE_RETENTION_DAYS", 30),
	}

	// Parse allowed origins
	originsStr := getEnv("ALLOWED_ORIGINS", "http://localhost:3000")
	cfg.AllowedOrigins = parseOrigins(originsStr)

	// Validate required fields
	if cfg.SessionSecret == "" {
		return nil, fmt.Errorf("SESSION_SECRET is required")
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

	// Validate auto-purge retention days
	if cfg.AutoPurgeRetentionDays < 1 {
		return nil, fmt.Errorf("AUTO_PURGE_RETENTION_DAYS must be at least 1")
	}

	return cfg, nil
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

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvRequired(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Required environment variable %s is not set", key)
	}
	return value
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}

	return value
}

func getEnvAsInt64(key string, defaultValue int64) int64 {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.ParseInt(valueStr, 10, 64)
	if err != nil {
		return defaultValue
	}

	return value
}

func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.ParseBool(valueStr)
	if err != nil {
		return defaultValue
	}

	return value
}

func parseOrigins(originsStr string) []string {
	// Split by comma and trim spaces
	origins := []string{}
	for _, origin := range splitAndTrim(originsStr, ",") {
		if origin != "" {
			origins = append(origins, origin)
		}
	}
	return origins
}

func splitAndTrim(s, sep string) []string {
	parts := []string{}
	for _, part := range splitString(s, sep) {
		trimmed := trimSpace(part)
		if trimmed != "" {
			parts = append(parts, trimmed)
		}
	}
	return parts
}

func splitString(s, sep string) []string {
	if s == "" {
		return []string{}
	}
	result := []string{}
	current := ""
	for _, char := range s {
		if string(char) == sep {
			result = append(result, current)
			current = ""
		} else {
			current += string(char)
		}
	}
	result = append(result, current)
	return result
}

func trimSpace(s string) string {
	start := 0
	end := len(s)

	for start < end && (s[start] == ' ' || s[start] == '\t' || s[start] == '\n') {
		start++
	}

	for end > start && (s[end-1] == ' ' || s[end-1] == '\t' || s[end-1] == '\n') {
		end--
	}

	return s[start:end]
}

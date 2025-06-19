package application

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds the application configuration...
type Config struct {
	Env    string
	Port   string
	DB_URL string
	S3     S3Config
}

type S3Config struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	Bucket    string
	UseSSL    bool
	Region    string
	BaseURL   string // Base URL for S3 bucket, e.g., https://s3.yourdomain.com/media/
}

func LoadConfig() *Config {
	err := godotenv.Load(".env") // Load environment variables from .env file
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	cfg := &Config{
		Env:    getEnv("APP_ENV", "development"),
		Port:   getEnv("PORT", "8080"),
		DB_URL: getEnv("DATABASE_URL", "postgres://user:pass@localhost:5432/threefive?sslmode=disable&timezone=UTC"),
		S3: S3Config{
			Endpoint:  getEnv("S3_ENDPOINT", "localhost:9000"),
			AccessKey: getEnv("S3_ACCESS_KEY", ""),
			SecretKey: getEnv("S3_SECRET_KEY", ""),
			Bucket:    getEnv("S3_BUCKET", "threefive"),
			UseSSL:    getEnv("S3_USE_SSL", "false") == "true",
			Region:    getEnv("S3_REGION", "sgp1"),
			BaseURL:   getEnv("S3_BASE_URL", "https://sgp1.vultrobjects.com/"),
		},
	}

	if cfg.DB_URL == "" {
		log.Fatal("DATABASE_URL is required but not set")
	}

	return cfg
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

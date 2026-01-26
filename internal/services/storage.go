// internal/services/storage.go
package services

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

// StorageProvider defines the interface for file storage
type StorageProvider interface {
	Upload(ctx context.Context, file *multipart.FileHeader, key string) (string, error)
	Delete(ctx context.Context, key string) error
	GetURL(key string) string
}

// LocalStorage implements local file storage
type LocalStorage struct {
	uploadDir string
	baseURL   string
}

func NewLocalStorage(uploadDir, baseURL string) *LocalStorage {
	// Ensure upload directory exists
	os.MkdirAll(uploadDir, 0o755)

	return &LocalStorage{
		uploadDir: uploadDir,
		baseURL:   baseURL,
	}
}

func (s *LocalStorage) Upload(ctx context.Context, file *multipart.FileHeader, key string) (string, error) {
	uploadPath := filepath.Join(s.uploadDir, key)

	// Ensure directory exists
	dir := filepath.Dir(uploadPath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	// Open uploaded file
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	// Create destination file
	dst, err := os.Create(uploadPath)
	if err != nil {
		return "", fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dst.Close()

	// Copy file
	if _, err = io.Copy(dst, src); err != nil {
		return "", fmt.Errorf("failed to save file: %w", err)
	}

	return key, nil
}

func (s *LocalStorage) Delete(ctx context.Context, key string) error {
	filePath := filepath.Join(s.uploadDir, key)
	if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete file: %w", err)
	}
	return nil
}

func (s *LocalStorage) GetURL(key string) string {
	return fmt.Sprintf("%s/%s", s.baseURL, key)
}

// S3Storage implements S3-compatible storage (AWS S3, Vultr Object Storage, etc.)
type S3Storage struct {
	client   *s3.Client
	bucket   string
	region   string
	endpoint string // For S3-compatible services like Vultr
	baseURL  string // Public URL for accessing files
}

type S3Config struct {
	Bucket          string
	Region          string
	Endpoint        string // Optional - for S3-compatible services (e.g., Vultr)
	AccessKeyID     string
	SecretAccessKey string
	BaseURL         string // Optional - custom CDN URL
}

// internal/services/storage.go

func NewS3Storage(ctx context.Context, cfg S3Config) (*S3Storage, error) {
	// Validate required fields
	if cfg.Bucket == "" {
		return nil, fmt.Errorf("S3 bucket is required")
	}
	if cfg.Region == "" {
		return nil, fmt.Errorf("S3 region is required")
	}
	if cfg.AccessKeyID == "" {
		return nil, fmt.Errorf("S3 access key ID is required")
	}
	if cfg.SecretAccessKey == "" {
		return nil, fmt.Errorf("S3 secret access key is required")
	}

	// Create static credentials provider
	credsProvider := credentials.NewStaticCredentialsProvider(
		cfg.AccessKeyID,
		cfg.SecretAccessKey,
		"", // Session token (empty for MinIO/Vultr)
	)

	// Load AWS config with credentials
	awsCfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(cfg.Region),
		config.WithCredentialsProvider(credsProvider),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	// Create S3 client with custom options for S3-compatible services (MinIO, Vultr, etc.)
	var s3Client *s3.Client
	if cfg.Endpoint != "" {
		// S3-compatible service (MinIO, Vultr Object Storage, etc.)
		s3Client = s3.NewFromConfig(awsCfg, func(o *s3.Options) {
			o.BaseEndpoint = aws.String(cfg.Endpoint)
			o.UsePathStyle = true // Required for MinIO and some S3-compatible services
		})
	} else {
		// AWS S3
		s3Client = s3.NewFromConfig(awsCfg)
	}

	// Determine base URL for accessing files
	baseURL := cfg.BaseURL
	if baseURL == "" {
		if cfg.Endpoint != "" {
			// For MinIO/S3-compatible services
			baseURL = fmt.Sprintf("%s/%s", cfg.Endpoint, cfg.Bucket)
		} else {
			// For AWS S3
			baseURL = fmt.Sprintf("https://%s.s3.%s.amazonaws.com", cfg.Bucket, cfg.Region)
		}
	}

	return &S3Storage{
		client:   s3Client,
		bucket:   cfg.Bucket,
		region:   cfg.Region,
		endpoint: cfg.Endpoint,
		baseURL:  baseURL,
	}, nil
}

func (s *S3Storage) Upload(ctx context.Context, file *multipart.FileHeader, key string) (string, error) {
	// Open uploaded file
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	// Determine content type
	contentType := file.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	// Upload to S3/MinIO with public-read ACL
	_, err = s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(key),
		Body:        src,
		ContentType: aws.String(contentType),
		ACL:         "public-read", // Make file publicly accessible
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload to S3: %w", err)
	}

	return key, nil
}

func (s *S3Storage) Delete(ctx context.Context, key string) error {
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("failed to delete from S3: %w", err)
	}

	return nil
}

func (s *S3Storage) GetURL(key string) string {
	return fmt.Sprintf("%s/%s", s.baseURL, key)
}

// GenerateStorageKey generates a unique storage key for a file
func GenerateStorageKey(filename string) string {
	ext := filepath.Ext(filename)
	timestamp := time.Now().Format("2006/01/02") // Organize by date: 2024/01/22
	uniqueID := uuid.New().String()[:8]

	return fmt.Sprintf("media/%s/%s%s", timestamp, uniqueID, ext)
}

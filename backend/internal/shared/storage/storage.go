package storage

import (
	"context"
	"io"
)

// StorageType represents the type of storage backend
type StorageType string

const (
	StorageTypeLocal StorageType = "local"
	StorageTypeS3    StorageType = "s3"
)

// UploadInput contains the data needed to upload a file
type UploadInput struct {
	File         io.Reader
	Filename     string
	ContentType  string
	Size         int64
	GenerateName bool // If true, generate unique filename
}

// UploadResult contains the result of an upload operation
type UploadResult struct {
	Filename     string
	Path         string
	URL          string
	ThumbnailURL string
	Size         int64
	Width        int
	Height       int
	S3Bucket     string
	S3Key        string
	S3Region     string
}

// Storage is the interface that all storage backends must implement
type Storage interface {
	// Upload uploads a file and returns the result
	Upload(ctx context.Context, input UploadInput) (*UploadResult, error)

	// Delete removes a file from storage
	Delete(ctx context.Context, path string) error

	// GetURL returns the public URL for a file
	GetURL(path string) string

	// Type returns the storage type
	Type() StorageType
}

// Config holds storage configuration
type Config struct {
	Type StorageType

	// Local storage config
	LocalBasePath string
	LocalBaseURL  string

	// S3 config
	S3Bucket    string
	S3Region    string
	S3AccessKey string
	S3SecretKey string
	S3Endpoint  string // Optional, for S3-compatible services
	S3PublicURL string // Optional, for custom CDN URLs
}

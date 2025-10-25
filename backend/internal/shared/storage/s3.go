// backend/internal/shared/storage/s3.go
package storage

import (
	"context"
	"fmt"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// S3Storage implements Storage interface for S3-compatible storage (MinIO, AWS S3, etc)
type S3Storage struct {
	client    *minio.Client
	bucket    string
	region    string
	endpoint  string
	publicURL string
	useSSL    bool
}

// NewS3Storage creates a new S3 storage instance using MinIO client
func NewS3Storage(bucket, region, accessKey, secretKey, endpoint, publicURL string, useSSL bool) (*S3Storage, error) {
	// Initialize MinIO client
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
		Region: region,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create MinIO client: %w", err)
	}

	return &S3Storage{
		client:    client,
		bucket:    bucket,
		region:    region,
		endpoint:  endpoint,
		publicURL: publicURL,
		useSSL:    useSSL,
	}, nil
}

// Upload uploads a file to S3-compatible storage
func (s *S3Storage) Upload(ctx context.Context, input UploadInput) (*UploadResult, error) {
	// Generate filename if needed
	filename := input.Filename
	if input.GenerateName {
		filename = generateUniqueFilename(input.Filename)
	}

	// Upload to S3/MinIO
	_, err := s.client.PutObject(ctx, s.bucket, filename, input.File, input.Size, minio.PutObjectOptions{
		ContentType: input.ContentType,
		UserMetadata: map[string]string{
			"x-amz-acl": "public-read",
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to upload to S3: %w", err)
	}

	url := s.GetURL(filename)

	result := &UploadResult{
		Filename: filename,
		Path:     filename,
		URL:      url,
		Size:     input.Size,
		S3Bucket: s.bucket,
		S3Key:    filename,
		S3Region: s.region,
	}

	return result, nil
}

// Delete removes a file from S3-compatible storage
func (s *S3Storage) Delete(ctx context.Context, path string) error {
	err := s.client.RemoveObject(ctx, s.bucket, path, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete from S3: %w", err)
	}
	return nil
}

// GetURL returns the public URL for a file
func (s *S3Storage) GetURL(path string) string {
	if s.publicURL != "" {
		return fmt.Sprintf("%s/%s", s.publicURL, path)
	}

	protocol := "http"
	if s.useSSL {
		protocol = "https"
	}
	return fmt.Sprintf("%s://%s/%s/%s", protocol, s.endpoint, s.bucket, path)
}

// Type returns the storage type
func (s *S3Storage) Type() StorageType {
	return StorageTypeS3
}

// backend/internal/shared/storage/local.go
package storage

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// LocalStorage implements Storage interface for local filesystem
type LocalStorage struct {
	basePath string
	baseURL  string
}

// NewLocalStorage creates a new local storage instance
func NewLocalStorage(basePath, baseURL string) (*LocalStorage, error) {
	// Ensure base path exists
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create base path: %w", err)
	}

	return &LocalStorage{
		basePath: basePath,
		baseURL:  strings.TrimSuffix(baseURL, "/"),
	}, nil
}

// Upload uploads a file to local storage
func (ls *LocalStorage) Upload(ctx context.Context, input UploadInput) (*UploadResult, error) {
	// Generate filename if needed
	filename := input.Filename
	if input.GenerateName {
		filename = generateUniqueFilename(input.Filename)
	}

	// Create full path
	fullPath := filepath.Join(ls.basePath, filename)

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}

	// Create file
	file, err := os.Create(fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// Copy content
	written, err := io.Copy(file, input.File)
	if err != nil {
		os.Remove(fullPath) // Cleanup on error
		return nil, fmt.Errorf("failed to write file: %w", err)
	}

	// Get image dimensions if it's an image
	width, height := 0, 0
	if isImage(input.ContentType) {
		file.Seek(0, 0) // Reset to beginning
		if img, _, err := image.DecodeConfig(file); err == nil {
			width = img.Width
			height = img.Height
		}
	}

	result := &UploadResult{
		Filename: filename,
		Path:     filename,
		URL:      ls.GetURL(filename),
		Size:     written,
		Width:    width,
		Height:   height,
	}

	return result, nil
}

// Delete removes a file from local storage
func (ls *LocalStorage) Delete(ctx context.Context, path string) error {
	fullPath := filepath.Join(ls.basePath, path)
	if err := os.Remove(fullPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete file: %w", err)
	}
	return nil
}

// GetURL returns the public URL for a file
func (ls *LocalStorage) GetURL(path string) string {
	return fmt.Sprintf("%s/%s", ls.baseURL, path)
}

// Type returns the storage type
func (ls *LocalStorage) Type() StorageType {
	return StorageTypeLocal
}

// Helper functions

func generateUniqueFilename(original string) string {
	ext := filepath.Ext(original)
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes) + ext
}

func isImage(contentType string) bool {
	return strings.HasPrefix(contentType, "image/")
}

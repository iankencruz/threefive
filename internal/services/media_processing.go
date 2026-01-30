package services

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// ProcessVideoThumbnail generates a thumbnail and uploads it, returning the storage key
func (s *MediaService) ProcessVideoThumbnail(ctx context.Context, videoFile *multipart.FileHeader, originalKey string) (string, error) {
	// Create temp files
	tmpDir := os.TempDir()
	videoPath := filepath.Join(tmpDir, fmt.Sprintf("video_%d%s", os.Getpid(), filepath.Ext(videoFile.Filename)))
	thumbPath := filepath.Join(tmpDir, fmt.Sprintf("thumb_%d.jpg", os.Getpid()))

	defer os.Remove(videoPath)
	defer os.Remove(thumbPath)

	// Save video to temp file
	src, err := videoFile.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open video: %w", err)
	}
	defer src.Close()

	dst, err := os.Create(videoPath)
	if err != nil {
		return "", fmt.Errorf("failed to create temp video: %w", err)
	}
	if _, err = io.Copy(dst, src); err != nil {
		dst.Close()
		return "", fmt.Errorf("failed to save video: %w", err)
	}
	dst.Close()

	// Run FFmpeg to extract thumbnail
	cmd := exec.CommandContext(ctx, "ffmpeg",
		"-i", videoPath,
		"-ss", "00:00:01.000",
		"-vframes", "1",
		"-vf", "scale=640:-1",
		"-q:v", "2",
		"-y",
		thumbPath,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("ffmpeg failed: %w\nOutput: %s", err, string(output))
	}

	// Verify thumbnail was created
	if _, err := os.Stat(thumbPath); err != nil {
		return "", fmt.Errorf("thumbnail not created: %w", err)
	}

	// Generate storage key for thumbnail
	thumbnailKey := strings.TrimSuffix(originalKey, filepath.Ext(originalKey)) + "_thumb.jpg"

	// Upload based on storage type
	switch storage := s.storage.(type) {
	case *LocalStorage:
		// For local storage, copy the file
		dstPath := filepath.Join(storage.uploadDir, thumbnailKey)
		if err := os.MkdirAll(filepath.Dir(dstPath), 0o755); err != nil {
			return "", fmt.Errorf("failed to create thumbnail directory: %w", err)
		}
		if err := copyFile(thumbPath, dstPath); err != nil {
			return "", fmt.Errorf("failed to copy thumbnail: %w", err)
		}

	case *S3Storage:
		// For S3, upload the file
		thumbFile, err := os.Open(thumbPath)
		if err != nil {
			return "", fmt.Errorf("failed to open thumbnail: %w", err)
		}
		defer thumbFile.Close()

		if err := storage.uploadFile(ctx, thumbFile, thumbnailKey, "image/jpeg"); err != nil {
			return "", fmt.Errorf("failed to upload thumbnail to S3: %w", err)
		}
	default:
		return "", fmt.Errorf("unsupported storage type")
	}

	return thumbnailKey, nil
}

// copyFile is a helper to copy files
func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}

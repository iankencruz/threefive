// internal/services/media_response.go
package services

import (
	"strings"
	"time"

	"github.com/iankencruz/threefive/database/generated"
	"github.com/jackc/pgx/v5/pgtype"
)

// MediaResponse is a view model for media that includes the public URL
type MediaResponse struct {
	ID               pgtype.UUID
	Filename         string
	OriginalFilename string
	MimeType         string
	FileSize         int64
	Width            pgtype.Int4
	Height           pgtype.Int4
	Duration         pgtype.Int4
	StorageType      string
	AltText          string
	URL              string // Public URL for serving the file
	ThumbnailURL     string // URL for thumbnail (if available)
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// ToMediaResponse converts a generated.Media to MediaResponse with URLs
func (s *MediaService) ToMediaResponse(media *generated.Media) MediaResponse {
	resp := MediaResponse{
		ID:               media.ID,
		Filename:         media.Filename,
		OriginalFilename: media.OriginalFilename,
		MimeType:         media.MimeType,
		FileSize:         media.FileSize,
		StorageType:      media.StorageType,
		CreatedAt:        media.CreatedAt,
		UpdatedAt:        media.UpdatedAt,
	}

	// Handle optional fields
	if media.Width.Valid {
		width := media.Width.Int32
		resp.Width = pgtype.Int4{Int32: width, Valid: true}
	}
	if media.Height.Valid {
		height := media.Height.Int32
		resp.Height = pgtype.Int4{Int32: height, Valid: true}
	}
	if media.Duration.Valid {
		duration := media.Duration.Int32
		resp.Duration = pgtype.Int4{Int32: duration, Valid: true}

	}
	if media.AltText.Valid {
		resp.AltText = media.AltText.String
	}

	// Generate URLs
	resp.URL = s.GetMediaURL(media)
	resp.ThumbnailURL = s.GetThumbnailURL(media)

	return resp
}

// ToMediaResponses converts a slice of generated.Media to MediaResponses
func (s *MediaService) ToMediaResponses(mediaList []generated.Media) []MediaResponse {
	responses := make([]MediaResponse, len(mediaList))
	for i, media := range mediaList {
		responses[i] = s.ToMediaResponse(&media)
	}
	return responses
}

// GetThumbnailURL returns the thumbnail URL for a media file
func (s *MediaService) GetThumbnailURL(media *generated.Media) string {
	// If no thumbnail exists, return original URL
	if !media.ThumbnailKey.Valid || media.ThumbnailKey.String == "" {
		return s.GetMediaURL(media)
	}

	// Use the same URL generation logic as GetMediaURL, but with thumbnail key
	if media.StorageType == "s3" {
		// For S3-compatible storage, use GetURL from storage provider
		if storage, ok := s.storage.(*S3Storage); ok {
			return storage.GetURL(media.ThumbnailKey.String)
		}
	} else if media.StorageType == "local" {
		// For local storage
		if storage, ok := s.storage.(*LocalStorage); ok {
			return storage.GetURL(media.ThumbnailKey.String)
		}
	}

	// Fallback to original URL
	return s.GetMediaURL(media)
}

// Helper: Check if media is an image
func IsImage(mimeType string) bool {
	return strings.HasPrefix(mimeType, "image/")
}

// Helper: Check if media is a video
func IsVideo(mimeType string) bool {
	return strings.HasPrefix(mimeType, "video/")
}

// Helper: Get file extension from mime type
func GetExtensionFromMimeType(mimeType string) string {
	switch mimeType {
	case "image/jpeg", "image/jpg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/gif":
		return ".gif"
	case "image/webp":
		return ".webp"
	case "video/mp4":
		return ".mp4"
	case "application/pdf":
		return ".pdf"
	default:
		return ""
	}
}

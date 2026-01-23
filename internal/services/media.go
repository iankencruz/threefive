// internal/services/media.go
package services

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/iankencruz/threefive/database/generated"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MediaService struct {
	db      *pgxpool.Pool
	queries *generated.Queries
	config  MediaConfig
	logger  *slog.Logger
}

type MediaConfig struct {
	UploadDir    string // Local upload directory
	MaxFileSize  int64  // Max file size in bytes
	AllowedTypes []string
	S3Bucket     string // For future S3 support
	S3Region     string
	BaseURL      string // Base URL for serving files
}

func NewMediaService(db *pgxpool.Pool, queries *generated.Queries, config MediaConfig) *MediaService {
	// Set defaults
	if config.UploadDir == "" {
		config.UploadDir = "./uploads"
	}
	if config.MaxFileSize == 0 {
		config.MaxFileSize = 50 * 1024 * 1024 // 50MB default
	}
	if len(config.AllowedTypes) == 0 {
		config.AllowedTypes = []string{"image/jpeg", "image/png", "image/gif", "image/webp", "video/mp4", "application/pdf"}
	}
	if config.BaseURL == "" {
		config.BaseURL = "/uploads"
	}

	// Ensure upload directory exists
	os.MkdirAll(config.UploadDir, 0o755)

	return &MediaService{
		db:      db,
		queries: queries,
		config:  config,
	}
}

// UploadMedia handles file upload and creates media record
func (s *MediaService) UploadMedia(ctx context.Context, file *multipart.FileHeader, altText string, uploadedBy pgtype.UUID) (*generated.Media, error) {
	// Validate file size
	if file.Size > s.config.MaxFileSize {
		return nil, fmt.Errorf("file size exceeds maximum allowed size of %d bytes", s.config.MaxFileSize)
	}

	// Validate file type
	mimeType := file.Header.Get("Content-Type")
	if !s.isAllowedType(mimeType) {
		return nil, fmt.Errorf("file type %s is not allowed", mimeType)
	}

	// Generate unique filename
	ext := filepath.Ext(file.Filename)
	timestamp := time.Now().Format("20060102")
	uniqueID := uuid.New().String()[:8]
	filename := fmt.Sprintf("%s-%s%s", timestamp, uniqueID, ext)

	// Create upload path
	uploadPath := filepath.Join(s.config.UploadDir, filename)

	// Open uploaded file
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	// Create destination file
	dst, err := os.Create(uploadPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dst.Close()

	// Copy file
	_, err = io.Copy(dst, src)
	if err != nil {
		return nil, fmt.Errorf("failed to save file: %w", err)
	}

	// Get image dimensions if it's an image
	var width, height *int32
	if strings.HasPrefix(mimeType, "image/") {
		// TODO: Implement image dimension detection
		// You can use image.DecodeConfig for this
	}

	// Create media record in database
	mediaID := uuid.New()
	var pgMediaID pgtype.UUID
	pgMediaID.Scan(mediaID.String())

	var pgAltText pgtype.Text
	if altText != "" {
		pgAltText.Scan(altText)
	}

	media, err := s.queries.CreateMedia(ctx, generated.CreateMediaParams{
		ID:               pgMediaID,
		Filename:         &filename,
		OriginalFilename: &file.Filename,
		MimeType:         &mimeType,
		FileSize:         file.Size,
		Width:            width,
		Height:           &height,
		StorageType:      pgtype.Text{String: "local", Valid: true},
		OriginalKey:      pgtype.Text{String: filename, Valid: true},
		AltText:          pgAltText,
		UploadedBy:       uploadedBy,
	})
	if err != nil {
		// Clean up file if database insert fails
		os.Remove(uploadPath)
		return nil, fmt.Errorf("failed to create media record: %w", err)
	}

	return &media, nil
}

// GetMediaByID retrieves a media file by ID
func (s *MediaService) GetMediaByID(ctx context.Context, id pgtype.UUID) (*generated.Media, error) {
	media, err := s.queries.GetMediaByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get media: %w", err)
	}
	return &media, nil
}

// ListMedia lists all media files with pagination
func (s *MediaService) ListMedia(ctx context.Context, limit, offset int32) ([]generated.Media, error) {
	if limit <= 0 {
		limit = 20
	}

	media, err := s.queries.ListMedia(ctx, generated.ListMediaParams{
		LimitVal:  limit,
		OffsetVal: offset,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list media: %w", err)
	}

	return media, nil
}

// ListMediaByType lists media filtered by MIME type pattern
func (s *MediaService) ListMediaByType(ctx context.Context, mimeTypePattern string, limit, offset int32) ([]generated.Media, error) {
	if limit <= 0 {
		limit = 20
	}

	media, err := s.queries.ListMediaByType(ctx, generated.ListMediaByTypeParams{
		MimeTypePattern: mimeTypePattern,
		LimitVal:        limit,
		OffsetVal:       offset,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list media by type: %w", err)
	}

	return media, nil
}

// ListMediaByUploader lists media uploaded by a specific user
func (s *MediaService) ListMediaByUploader(ctx context.Context, uploadedBy pgtype.UUID, limit, offset int32) ([]generated.Media, error) {
	if limit <= 0 {
		limit = 20
	}

	media, err := s.queries.ListMediaByUploader(ctx, generated.ListMediaByUploaderParams{
		UploadedBy: uploadedBy,
		LimitVal:   limit,
		OffsetVal:  offset,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list media by uploader: %w", err)
	}

	return media, nil
}

// UpdateMediaAltText updates the alt text for a media file
func (s *MediaService) UpdateMediaAltText(ctx context.Context, id pgtype.UUID, altText string) (*generated.Media, error) {
	var pgAltText pgtype.Text
	pgAltText.Scan(altText)

	media, err := s.queries.UpdateMediaAltText(ctx, generated.UpdateMediaAltTextParams{
		ID:      id,
		AltText: pgAltText,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update media alt text: %w", err)
	}

	return &media, nil
}

// DeleteMedia soft deletes a media file
func (s *MediaService) DeleteMedia(ctx context.Context, id pgtype.UUID) error {
	err := s.queries.SoftDeleteMedia(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete media: %w", err)
	}

	// Note: We keep the physical file for recovery
	// Use PurgeOldDeletedMedia to permanently delete old files

	return nil
}

// RestoreMedia restores a soft-deleted media file
func (s *MediaService) RestoreMedia(ctx context.Context, id pgtype.UUID) error {
	err := s.queries.RestoreMedia(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to restore media: %w", err)
	}
	return nil
}

// HardDeleteMedia permanently deletes a media file and removes the physical file
func (s *MediaService) HardDeleteMedia(ctx context.Context, id pgtype.UUID) error {
	// Get media to find filename
	media, err := s.queries.GetMediaByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get media: %w", err)
	}

	// Delete from database
	err = s.queries.HardDeleteMedia(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete media from database: %w", err)
	}

	// Delete physical file
	filePath := filepath.Join(s.config.UploadDir, media.Filename)
	if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
		// Log error but don't fail - database record is already deleted
		fmt.Printf("Warning: failed to delete physical file %s: %v\n", filePath, err)
	}

	return nil
}

// PurgeOldDeletedMedia permanently deletes media files deleted more than 30 days ago
func (s *MediaService) PurgeOldDeletedMedia(ctx context.Context) error {
	// Get list of deleted media before purging
	deletedMedia, err := s.queries.GetDeletedMedia(ctx, generated.GetDeletedMediaParams{
		LimitVal:  1000,
		OffsetVal: 0,
	})
	if err != nil {
		return fmt.Errorf("failed to get deleted media: %w", err)
	}

	// Delete physical files
	for _, media := range deletedMedia {
		if media.DeletedAt.Valid && time.Since(media.DeletedAt.Time) > 30*24*time.Hour {
			filePath := filepath.Join(s.config.UploadDir, media.Filename)
			os.Remove(filePath) // Ignore errors
		}
	}

	// Purge from database
	err = s.queries.PurgeOldDeletedMedia(ctx)
	if err != nil {
		return fmt.Errorf("failed to purge old deleted media: %w", err)
	}

	return nil
}

// Media Relations

// LinkMediaToEntity creates a relationship between media and an entity
func (s *MediaService) LinkMediaToEntity(ctx context.Context, mediaID pgtype.UUID, entityType string, entityID pgtype.UUID, relationType string, sortOrder int32) error {
	relationID := uuid.New()
	var pgRelationID pgtype.UUID
	pgRelationID.Scan(relationID.String())

	var pgEntityID pgtype.UUID
	pgEntityID.Scan(entityID)

	_, err := s.queries.CreateMediaRelation(ctx, generated.CreateMediaRelationParams{
		ID:           pgRelationID,
		MediaID:      mediaID,
		EntityType:   entityType,
		EntityID:     pgEntityID,
		RelationType: pgtype.Text{String: relationType, Valid: true},
		SortOrder:    pgtype.Int4{Int32: sortOrder, Valid: true},
	})
	if err != nil {
		return fmt.Errorf("failed to link media to entity: %w", err)
	}

	return nil
}

// GetMediaForEntity retrieves all media linked to an entity
func (s *MediaService) GetMediaForEntity(ctx context.Context, entityType string, entityID pgtype.UUID) ([]generated.Media, error) {
	media, err := s.queries.GetMediaForEntity(ctx, generated.GetMediaForEntityParams{
		EntityType: entityType,
		EntityID:   entityID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get media for entity: %w", err)
	}

	return media, nil
}

// GetFeaturedMediaForEntity retrieves the featured media for an entity
func (s *MediaService) GetFeaturedMediaForEntity(ctx context.Context, entityType string, entityID pgtype.UUID) (*generated.Media, error) {
	media, err := s.queries.GetFeaturedMediaForEntity(ctx, generated.GetFeaturedMediaForEntityParams{
		EntityType: entityType,
		EntityID:   entityID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get featured media: %w", err)
	}

	return &media, nil
}

// GetGalleryMediaForEntity retrieves gallery media for an entity
func (s *MediaService) GetGalleryMediaForEntity(ctx context.Context, entityType string, entityID pgtype.UUID) ([]generated.Media, error) {
	media, err := s.queries.GetGalleryMediaForEntity(ctx, generated.GetGalleryMediaForEntityParams{
		EntityType: entityType,
		EntityID:   entityID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get gallery media: %w", err)
	}

	return media, nil
}

// UnlinkMediaFromEntity removes the relationship between media and an entity
func (s *MediaService) UnlinkMediaFromEntity(ctx context.Context, mediaID pgtype.UUID, entityType string, entityID pgtype.UUID) error {
	err := s.queries.DeleteMediaRelation(ctx, generated.DeleteMediaRelationParams{
		MediaID:    mediaID,
		EntityType: entityType,
		EntityID:   entityID,
	})
	if err != nil {
		return fmt.Errorf("failed to unlink media: %w", err)
	}

	return nil
}

// GetMediaStats retrieves media statistics
func (s *MediaService) GetMediaStats(ctx context.Context) (*generated.GetMediaStatsRow, error) {
	stats, err := s.queries.GetMediaStats(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get media stats: %w", err)
	}

	return &stats, nil
}

// CountMedia returns the total count of media files
func (s *MediaService) CountMedia(ctx context.Context) (int64, error) {
	count, err := s.queries.CountMedia(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to count media: %w", err)
	}
	return count, nil
}

// Helper methods

func (s *MediaService) isAllowedType(mimeType string) bool {
	for _, allowed := range s.config.AllowedTypes {
		if allowed == mimeType {
			return true
		}
	}
	return false
}

// GetMediaURL returns the public URL for a media file
func (s *MediaService) GetMediaURL(media *generated.Media) string {
	if media.StorageType.String == "s3" {
		// TODO: Construct S3 URL
		return ""
	}
	return fmt.Sprintf("%s/%s", s.config.BaseURL, media.Filename)
}

// internal/services/media.go
package services

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"slices"
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
	storage StorageProvider
	config  MediaConfig
}

type MediaConfig struct {
	MaxFileSize  int64
	AllowedTypes []string
}

func NewMediaService(db *pgxpool.Pool, queries *generated.Queries, storage StorageProvider, config MediaConfig) *MediaService {
	// Set defaults
	if config.MaxFileSize == 0 {
		config.MaxFileSize = 50 * 1024 * 1024 // 50MB
	}
	if len(config.AllowedTypes) == 0 {
		config.AllowedTypes = []string{
			"image/jpeg", "image/png", "image/gif", "image/webp",
			"video/mp4", "video/quicktime",
			"application/pdf",
		}
	}

	return &MediaService{
		db:      db,
		queries: queries,
		storage: storage,
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

	// Generate unique storage key
	storageKey := GenerateStorageKey(file.Filename)

	// Upload file using storage provider
	uploadedKey, err := s.storage.Upload(ctx, file, storageKey)
	if err != nil {
		return nil, fmt.Errorf("failed to upload file: %w", err)
	}

	// Get image dimensions if it's an image
	var width, height pgtype.Int4
	if strings.HasPrefix(mimeType, "image/") {
		// TODO: Implement image dimension detection using image.DecodeConfig
	}

	// Determine storage type
	var storageType string
	switch s.storage.(type) {
	case *S3Storage:
		storageType = "s3"
	case *LocalStorage:
		storageType = "local"
	default:
		storageType = "unknown"
	}

	// Get S3 details if applicable
	var s3Bucket, s3Region pgtype.Text
	if s3Storage, ok := s.storage.(*S3Storage); ok {
		s3Bucket = pgtype.Text{String: s3Storage.bucket, Valid: true}
		s3Region = pgtype.Text{String: s3Storage.region, Valid: true}
	}

	// Create media record in database
	mediaID := uuid.New()
	var pgMediaID pgtype.UUID
	if err := pgMediaID.Scan(mediaID.String()); err != nil {
		return nil, fmt.Errorf("failed to convert media ID: %w", err)
	}

	// Prepare optional fields
	var pgAltText pgtype.Text
	if altText != "" {
		pgAltText = pgtype.Text{String: altText, Valid: true}
	}

	var pgOriginalKey pgtype.Text
	pgOriginalKey = pgtype.Text{String: uploadedKey, Valid: true}

	filename := filepath.Base(uploadedKey)

	media, err := s.queries.CreateMedia(ctx, generated.CreateMediaParams{
		ID:               pgMediaID,
		Filename:         filename,
		OriginalFilename: file.Filename,
		MimeType:         mimeType,
		FileSize:         file.Size,
		Width:            width,
		Height:           height,
		StorageType:      storageType, // Fixed: SQLC expects string for NOT NULL TEXT
		S3Bucket:         s3Bucket,
		S3Region:         s3Region,
		OriginalKey:      pgOriginalKey,
		AltText:          pgAltText,
		UploadedBy:       uploadedBy,
	})
	if err != nil {
		// Clean up uploaded file if database insert fails
		s.storage.Delete(ctx, uploadedKey)
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
	if altText != "" {
		pgAltText = pgtype.Text{String: altText, Valid: true}
	}

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

// HardDeleteMedia permanently deletes a media file and removes from storage
func (s *MediaService) HardDeleteMedia(ctx context.Context, id pgtype.UUID) error {
	// Get media to find storage key
	media, err := s.GetMediaByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get media: %w", err)
	}

	// Delete from database
	err = s.queries.HardDeleteMedia(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete media from database: %w", err)
	}

	// Delete from storage
	if media.OriginalKey.Valid {
		if err := s.storage.Delete(ctx, media.OriginalKey.String); err != nil {
			fmt.Printf("Warning: failed to delete file from storage: %v\n", err)
		}
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

	// Delete physical files older than 30 days
	thirtyDaysAgo := time.Now().Add(-30 * 24 * time.Hour)
	for _, media := range deletedMedia {
		// Fixed: DeletedAt is a *time.Time pointer, not pgtype.Timestamp
		if media.DeletedAt != nil && media.DeletedAt.Before(thirtyDaysAgo) {
			if media.OriginalKey.Valid {
				s.storage.Delete(ctx, media.OriginalKey.String) // Ignore errors
			}
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
	if err := pgRelationID.Scan(relationID.String()); err != nil {
		return fmt.Errorf("failed to convert relation ID: %w", err)
	}

	var pgSortOrder pgtype.Int4
	pgSortOrder = pgtype.Int4{Int32: sortOrder, Valid: true}

	_, err := s.queries.CreateMediaRelation(ctx, generated.CreateMediaRelationParams{
		ID:           pgRelationID,
		MediaID:      mediaID,
		EntityType:   entityType,
		EntityID:     entityID,
		RelationType: relationType, // Fixed: SQLC expects string for NOT NULL TEXT
		SortOrder:    pgSortOrder,
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
	return slices.Contains(s.config.AllowedTypes, mimeType)
}

// GetMediaURL returns the public URL for a media file
func (s *MediaService) GetMediaURL(media *generated.Media) string {
	if !media.OriginalKey.Valid {
		return ""
	}
	return s.storage.GetURL(media.OriginalKey.String)
}

// UpdateMedia updates media metadata (alt text)
func (s *MediaService) UpdateMedia(ctx context.Context, mediaID pgtype.UUID, altText string) (*generated.Media, error) {
	// Convert altText to pgtype.Text
	var altTextPg pgtype.Text
	if altText != "" {
		altTextPg.String = altText
		altTextPg.Valid = true
	}

	// Update media in database
	media, err := s.queries.UpdateMediaAltText(ctx, generated.UpdateMediaAltTextParams{
		ID:      mediaID,
		AltText: altTextPg,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update media: %w", err)
	}

	return &media, nil
}

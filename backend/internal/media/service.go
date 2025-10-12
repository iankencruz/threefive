// backend/internal/media/service.go
package media

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/iankencruz/threefive/internal/shared/errors"
	"github.com/iankencruz/threefive/internal/shared/sqlc"
	"github.com/iankencruz/threefive/internal/shared/storage"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Service handles media business logic
type Service struct {
	db      *pgxpool.Pool
	queries *sqlc.Queries
	storage storage.Storage
}

// NewService creates a new media service
func NewService(db *pgxpool.Pool, queries *sqlc.Queries, storage storage.Storage) *Service {
	return &Service{
		db:      db,
		queries: queries,
		storage: storage,
	}
}

// UploadFile uploads a file to storage and creates a media record
func (s *Service) UploadFile(ctx context.Context, file *multipart.FileHeader, userID uuid.UUID) (*sqlc.Media, error) {
	// Validate file
	if err := s.validateFile(file); err != nil {
		return nil, err
	}

	// Open file
	src, err := file.Open()
	if err != nil {
		return nil, errors.Internal("Failed to open file", err)
	}
	defer src.Close()

	// Upload to storage
	uploadResult, err := s.storage.Upload(ctx, storage.UploadInput{
		File:         src,
		Filename:     file.Filename,
		ContentType:  file.Header.Get("Content-Type"),
		Size:         file.Size,
		GenerateName: true, // Generate unique filename
	})
	if err != nil {
		return nil, errors.Internal("Failed to upload file", err)
	}

	// Determine storage type values
	var s3Bucket, s3Key, s3Region *string
	if s.storage.Type() == storage.StorageTypeS3 {
		s3Bucket = &uploadResult.S3Bucket
		s3Key = &uploadResult.S3Key
		s3Region = &uploadResult.S3Region
	}

	// Create media record
	media, err := s.queries.CreateMedia(ctx, sqlc.CreateMediaParams{
		Filename:         uploadResult.Filename,
		OriginalFilename: file.Filename,
		MimeType:         file.Header.Get("Content-Type"),
		SizeBytes:        file.Size,
		Width:            intToNullInt32(uploadResult.Width),
		Height:           intToNullInt32(uploadResult.Height),
		StorageType:      sqlc.StorageType(s.storage.Type()),
		StoragePath:      uploadResult.Path,
		S3Bucket:         toNullString(s3Bucket),
		S3Key:            toNullString(s3Key),
		S3Region:         toNullString(s3Region),
		Url:              toNullString(&uploadResult.URL),
		ThumbnailUrl:     toNullString(&uploadResult.ThumbnailURL),
		UploadedBy:       userID,
	})
	if err != nil {
		// Cleanup uploaded file on database error
		s.storage.Delete(ctx, uploadResult.Path)
		return nil, errors.Internal("Failed to create media record", err)
	}

	return &media, nil
}

// GetMediaByID retrieves a media record by ID
func (s *Service) GetMediaByID(ctx context.Context, id uuid.UUID) (*sqlc.Media, error) {
	media, err := s.queries.GetMediaByID(ctx, id)
	if err != nil {
		return nil, errors.NotFound("Media not found", "media_not_found")
	}
	return &media, nil
}

// ListMedia retrieves paginated list of media
func (s *Service) ListMedia(ctx context.Context, limit, offset int32) ([]sqlc.Media, error) {
	media, err := s.queries.ListMedia(ctx, sqlc.ListMediaParams{
		LimitVal:  limit,
		OffsetVal: offset,
	})
	if err != nil {
		return nil, errors.Internal("Failed to list media", err)
	}
	return media, nil
}

// DeleteMedia soft deletes a media record and removes file from storage
func (s *Service) DeleteMedia(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	// Get media record
	media, err := s.queries.GetMediaByID(ctx, id)
	if err != nil {
		return errors.NotFound("Media not found", "media_not_found")
	}

	// Check ownership (admins can delete any, users can only delete their own)
	// TODO: Add proper role check when roles are implemented
	if media.UploadedBy != userID {
		return errors.Forbidden("Cannot delete media uploaded by another user", "forbidden")
	}

	// Delete from storage
	if err := s.storage.Delete(ctx, media.StoragePath); err != nil {
		// Log error but continue with database deletion
		// In production, you might want to queue this for retry
		fmt.Printf("Warning: failed to delete file from storage: %v\n", err)
	}

	// Soft delete from database
	if err := s.queries.SoftDeleteMedia(ctx, id); err != nil {
		return errors.Internal("Failed to delete media", err)
	}

	return nil
}

// HardDeleteMedia permanently deletes a media record
func (s *Service) HardDeleteMedia(ctx context.Context, id uuid.UUID) error {
	// Get media record
	media, err := s.queries.GetMediaByID(ctx, id)
	if err != nil {
		return errors.NotFound("Media not found", "media_not_found")
	}

	// Delete from storage
	if err := s.storage.Delete(ctx, media.StoragePath); err != nil {
		return errors.Internal("Failed to delete file from storage", err)
	}

	// Hard delete from database
	if err := s.queries.HardDeleteMedia(ctx, id); err != nil {
		return errors.Internal("Failed to delete media", err)
	}

	return nil
}

// LinkMediaToEntity links media to an entity (project, page, gallery, etc)
func (s *Service) LinkMediaToEntity(ctx context.Context, mediaID uuid.UUID, entityType string, entityID uuid.UUID, sortOrder int32) error {
	// Verify media exists
	if _, err := s.queries.GetMediaByID(ctx, mediaID); err != nil {
		return errors.NotFound("Media not found", "media_not_found")
	}

	// Create link
	_, err := s.queries.LinkMediaToEntity(ctx, sqlc.LinkMediaToEntityParams{
		MediaID:    mediaID,
		EntityType: entityType,
		EntityID:   entityID,
		SortOrder:  intToNullInt32(int(sortOrder)),
	})
	if err != nil {
		return errors.Internal("Failed to link media", err)
	}

	return nil
}

// UnlinkMediaFromEntity removes link between media and entity
func (s *Service) UnlinkMediaFromEntity(ctx context.Context, mediaID uuid.UUID, entityType string, entityID uuid.UUID) error {
	err := s.queries.UnlinkMediaFromEntity(ctx, sqlc.UnlinkMediaFromEntityParams{
		MediaID:    mediaID,
		EntityType: entityType,
		EntityID:   entityID,
	})
	if err != nil {
		return errors.Internal("Failed to unlink media", err)
	}
	return nil
}

// GetMediaForEntity retrieves all media linked to an entity
func (s *Service) GetMediaForEntity(ctx context.Context, entityType string, entityID uuid.UUID) ([]sqlc.Media, error) {
	media, err := s.queries.GetMediaForEntity(ctx, sqlc.GetMediaForEntityParams{
		EntityType: entityType,
		EntityID:   entityID,
	})
	if err != nil {
		return nil, errors.Internal("Failed to get media for entity", err)
	}
	return media, nil
}

// GetMediaStats returns statistics about media storage
func (s *Service) GetMediaStats(ctx context.Context) (*sqlc.GetMediaStatsRow, error) {
	stats, err := s.queries.GetMediaStats(ctx)
	if err != nil {
		return nil, errors.Internal("Failed to get media stats", err)
	}
	return &stats, nil
}

// Validation helpers

func (s *Service) validateFile(file *multipart.FileHeader) error {
	ext := strings.ToLower(filepath.Ext(file.Filename))

	// Different size limits based on file type
	var maxFileSize int64
	videoExtensions := map[string]bool{".mp4": true, ".mov": true, ".avi": true}

	if videoExtensions[ext] {
		maxFileSize = 200 * 1024 * 1024 // 200MB for videos
	} else {
		maxFileSize = 50 * 1024 * 1024 // 50MB for images/docs
	}

	if file.Size > maxFileSize {
		fileType := "file"
		if videoExtensions[ext] {
			fileType = "video"
		}
		return errors.BadRequest(
			fmt.Sprintf("File size exceeds %dMB limit for %ss", maxFileSize/(1024*1024), fileType),
			"file_too_large",
		)
	}

	// Validate file extension
	allowedExtensions := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".webp": true,
		".svg":  true,
		".pdf":  true,
		".mp4":  true,
		".mov":  true,
		".avi":  true,
	}

	if !allowedExtensions[ext] {
		return errors.BadRequest(fmt.Sprintf("File type %s not allowed", ext), "invalid_file_type")
	}

	return nil
}

// Helper functions for nullable types
// For nullable int from pointer
func toNullInt32(val *int) pgtype.Int4 {
	if val == nil {
		return pgtype.Int4{Valid: false}
	}
	return pgtype.Int4{Int32: int32(*val), Valid: true}
}

// For nullable int from direct value (use when value might be 0 = NULL)
func intToNullInt32(val int) pgtype.Int4 {
	if val == 0 {
		return pgtype.Int4{Valid: false}
	}
	return pgtype.Int4{Int32: int32(val), Valid: true}
}

// For nullable string from pointer
func toNullString(val *string) pgtype.Text {
	if val == nil || *val == "" {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{String: *val, Valid: true}
}

// For nullable string from direct value
func stringToNullText(val string) pgtype.Text {
	if val == "" {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{String: val, Valid: true}
}

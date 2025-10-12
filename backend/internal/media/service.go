// backend/internal/media/service.go
package media

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/iankencruz/threefive/internal/media/processing"
	"github.com/iankencruz/threefive/internal/shared/errors"
	"github.com/iankencruz/threefive/internal/shared/sqlc"
	"github.com/iankencruz/threefive/internal/shared/storage"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Service handles media business logic
type Service struct {
	db        *pgxpool.Pool
	queries   *sqlc.Queries
	storage   storage.Storage
	processor *processing.Processor
}

// NewService creates a new media service
func NewService(db *pgxpool.Pool, queries *sqlc.Queries, storage storage.Storage) *Service {
	// Initialize processor with default config
	processorConfig := processing.LoadConfigFromEnv()
	processor, err := processing.NewProcessor(processorConfig, "./tmp/processing")
	if err != nil {
		// Log warning but don't fail service creation
		fmt.Printf("Warning: Failed to initialize media processor: %v\n", err)
		fmt.Println("Media processing features will be disabled")
		processor = nil
	}

	return &Service{
		db:        db,
		queries:   queries,
		storage:   storage,
		processor: processor,
	}
}

// UploadFile uploads a file to storage and creates a media record
func (s *Service) UploadFile(ctx context.Context, file *multipart.FileHeader, userID uuid.UUID) (*sqlc.Media, error) {
	if err := s.validateFile(file); err != nil {
		return nil, err
	}

	src, err := file.Open()
	if err != nil {
		return nil, errors.Internal("Failed to open file", err)
	}
	defer src.Close()

	filename := file.Filename
	contentType := file.Header.Get("Content-Type")
	isImage := processing.IsImageFile(filename)
	isVideo := processing.IsVideoFile(filename)

	var uploadResult *storage.UploadResult

	// Process media if processor is available
	if s.processor != nil && (isImage || isVideo) {
		uploadResult, err = s.processAndUpload(ctx, src, filename, contentType, file.Size, isImage, isVideo)
		if err != nil {
			// Fall back to direct upload if processing fails
			fmt.Printf("Warning: Processing failed, falling back to direct upload: %v\n", err)
			src.Seek(0, io.SeekStart)
			uploadResult, err = s.directUpload(ctx, src, filename, contentType, file.Size)
			if err != nil {
				return nil, err
			}
		}
	} else {
		uploadResult, err = s.directUpload(ctx, src, filename, contentType, file.Size)
		if err != nil {
			return nil, err
		}
	}

	// Determine storage type values
	var s3Bucket, s3Key, s3Region *string
	if s.storage.Type() == storage.StorageTypeS3 {
		s3Bucket = &uploadResult.S3Bucket
		s3Key = &uploadResult.S3Key
		s3Region = &uploadResult.S3Region
	}

	// Create media record - NOTE: Capital U in Url and ThumbnailUrl
	media, err := s.queries.CreateMedia(ctx, sqlc.CreateMediaParams{
		Filename:         uploadResult.Filename,
		OriginalFilename: file.Filename,
		MimeType:         contentType,
		SizeBytes:        uploadResult.Size,
		Width:            intToNullInt32(uploadResult.Width),
		Height:           intToNullInt32(uploadResult.Height),
		StorageType:      sqlc.StorageType(s.storage.Type()),
		StoragePath:      uploadResult.Path,
		S3Bucket:         toNullString(s3Bucket),
		S3Key:            toNullString(s3Key),
		S3Region:         toNullString(s3Region),
		Url:              toNullString(&uploadResult.URL),          // Capital U
		ThumbnailUrl:     toNullString(&uploadResult.ThumbnailURL), // Capital U
		UploadedBy:       userID,
	})
	if err != nil {
		s.storage.Delete(ctx, uploadResult.Path)
		if uploadResult.ThumbnailURL != "" {
			thumbPath := strings.TrimPrefix(uploadResult.ThumbnailURL, s.storage.GetURL(""))
			s.storage.Delete(ctx, thumbPath)
		}
		return nil, errors.Internal("Failed to create media record", err)
	}

	return &media, nil
}

// processAndUpload processes media and uploads the results
func (s *Service) processAndUpload(ctx context.Context, src io.Reader, filename, contentType string, size int64, isImage, isVideo bool) (*storage.UploadResult, error) {
	var processResult *processing.ProcessResult
	var err error

	if isImage {
		processResult, err = s.processor.ProcessImage(ctx, src, filename)
		if err != nil {
			return nil, fmt.Errorf("image processing failed: %w", err)
		}
	} else if isVideo {
		// Save to temp file for video processing
		tempPath := filepath.Join(s.processor.WorkDir(), "temp_"+filename)
		tempFile, err := os.Create(tempPath)
		if err != nil {
			return nil, fmt.Errorf("failed to create temp file: %w", err)
		}

		if _, err := io.Copy(tempFile, src); err != nil {
			tempFile.Close()
			os.Remove(tempPath)
			return nil, fmt.Errorf("failed to save temp file: %w", err)
		}
		tempFile.Close()

		processResult, err = s.processor.ProcessVideo(ctx, tempPath, filename)
		if err != nil {
			os.Remove(tempPath)
			return nil, fmt.Errorf("video processing failed: %w", err)
		}

		defer os.Remove(tempPath)
	}

	// Upload processed file
	processedFile, err := os.Open(processResult.ProcessedPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open processed file: %w", err)
	}
	defer processedFile.Close()
	defer os.Remove(processResult.ProcessedPath)

	// Determine new content type
	newContentType := contentType
	if processResult.Format == "webp" {
		newContentType = "image/webp"
	} else if processResult.Format == "mp4" {
		newContentType = "video/mp4"
	}

	uploadResult, err := s.storage.Upload(ctx, storage.UploadInput{
		File:         processedFile,
		Filename:     filepath.Base(processResult.ProcessedPath),
		ContentType:  newContentType,
		Size:         processResult.Size,
		GenerateName: true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to upload processed file: %w", err)
	}

	uploadResult.Width = processResult.Width
	uploadResult.Height = processResult.Height

	// Upload thumbnail if exists
	if processResult.ThumbnailPath != "" {
		thumbFile, err := os.Open(processResult.ThumbnailPath)
		if err == nil {
			defer thumbFile.Close()
			defer os.Remove(processResult.ThumbnailPath)

			thumbInfo, _ := thumbFile.Stat()
			thumbResult, err := s.storage.Upload(ctx, storage.UploadInput{
				File:         thumbFile,
				Filename:     filepath.Base(processResult.ThumbnailPath),
				ContentType:  "image/webp",
				Size:         thumbInfo.Size(),
				GenerateName: true,
			})
			if err == nil {
				uploadResult.ThumbnailURL = thumbResult.URL
			}
		}
	}

	return uploadResult, nil
}

// directUpload uploads a file directly without processing
func (s *Service) directUpload(ctx context.Context, src io.Reader, filename, contentType string, size int64) (*storage.UploadResult, error) {
	return s.storage.Upload(ctx, storage.UploadInput{
		File:         src,
		Filename:     filename,
		ContentType:  contentType,
		Size:         size,
		GenerateName: true,
	})
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
	media, err := s.queries.GetMediaByID(ctx, id)
	if err != nil {
		return errors.NotFound("Media not found", "media_not_found")
	}

	if media.UploadedBy != userID {
		return errors.Forbidden("Cannot delete media uploaded by another user", "forbidden")
	}

	if err := s.storage.Delete(ctx, media.StoragePath); err != nil {
		fmt.Printf("Warning: failed to delete file from storage: %v\n", err)
	}

	// NOTE: Capital U in ThumbnailUrl
	if media.ThumbnailUrl.Valid && media.ThumbnailUrl.String != "" {
		thumbPath := strings.TrimPrefix(media.ThumbnailUrl.String, s.storage.GetURL(""))
		if err := s.storage.Delete(ctx, thumbPath); err != nil {
			fmt.Printf("Warning: failed to delete thumbnail from storage: %v\n", err)
		}
	}

	if err := s.queries.SoftDeleteMedia(ctx, id); err != nil {
		return errors.Internal("Failed to delete media", err)
	}

	return nil
}

// HardDeleteMedia permanently deletes a media record
func (s *Service) HardDeleteMedia(ctx context.Context, id uuid.UUID) error {
	media, err := s.queries.GetMediaByID(ctx, id)
	if err != nil {
		return errors.NotFound("Media not found", "media_not_found")
	}

	if err := s.storage.Delete(ctx, media.StoragePath); err != nil {
		fmt.Printf("Warning: failed to delete file from storage: %v\n", err)
	}

	// NOTE: Capital U in ThumbnailUrl
	if media.ThumbnailUrl.Valid && media.ThumbnailUrl.String != "" {
		thumbPath := strings.TrimPrefix(media.ThumbnailUrl.String, s.storage.GetURL(""))
		if err := s.storage.Delete(ctx, thumbPath); err != nil {
			fmt.Printf("Warning: failed to delete thumbnail from storage: %v\n", err)
		}
	}

	if err := s.queries.HardDeleteMedia(ctx, id); err != nil {
		return errors.Internal("Failed to hard delete media", err)
	}

	return nil
}

// LinkMediaToEntity links media to an entity
func (s *Service) LinkMediaToEntity(ctx context.Context, mediaID uuid.UUID, entityType string, entityID uuid.UUID, sortOrder int32) error {
	_, err := s.queries.LinkMediaToEntity(ctx, sqlc.LinkMediaToEntityParams{
		MediaID:    mediaID,
		EntityType: entityType,
		EntityID:   entityID,
		SortOrder:  intToNullInt32(int(sortOrder)),
	})
	return err
}

// UnlinkMediaFromEntity unlinks media from an entity
func (s *Service) UnlinkMediaFromEntity(ctx context.Context, mediaID uuid.UUID, entityType string, entityID uuid.UUID) error {
	return s.queries.UnlinkMediaFromEntity(ctx, sqlc.UnlinkMediaFromEntityParams{
		MediaID:    mediaID,
		EntityType: entityType,
		EntityID:   entityID,
	})
}

// GetMediaForEntity retrieves all media for an entity
func (s *Service) GetMediaForEntity(ctx context.Context, entityType string, entityID uuid.UUID) ([]sqlc.Media, error) {
	return s.queries.GetMediaForEntity(ctx, sqlc.GetMediaForEntityParams{
		EntityType: entityType,
		EntityID:   entityID,
	})
}

// GetMediaStats returns media statistics
func (s *Service) GetMediaStats(ctx context.Context) (*MediaStats, error) {
	stats, err := s.queries.GetMediaStats(ctx)
	if err != nil {
		return nil, errors.Internal("Failed to get media stats", err)
	}

	return &MediaStats{
		TotalCount:      stats.TotalFiles,
		TotalSize:       stats.TotalSizeBytes,
		UniqueUploaders: stats.UniqueUploaders,
	}, nil
}

// MediaStats contains media statistics
type MediaStats struct {
	TotalCount      int64 `json:"total_count"`
	TotalSize       int64 `json:"total_size"`
	UniqueUploaders int64 `json:"unique_uploaders"`
}

// validateFile validates uploaded file
func (s *Service) validateFile(file *multipart.FileHeader) error {
	maxSize := int64(50 << 20) // 50MB
	if processing.IsVideoFile(file.Filename) {
		maxSize = 200 << 20 // 200MB
	}

	if file.Size > maxSize {
		return errors.BadRequest(fmt.Sprintf("File too large (max %dMB)", maxSize>>20), "file_too_large")
	}

	contentType := file.Header.Get("Content-Type")
	if contentType == "" {
		return errors.BadRequest("Content-Type header required", "invalid_content_type")
	}

	return nil
}

// Helper functions
func intToNullInt32(v int) pgtype.Int4 {
	if v == 0 {
		return pgtype.Int4{Valid: false}
	}
	return pgtype.Int4{Int32: int32(v), Valid: true}
}

func toNullString(s *string) pgtype.Text {
	if s == nil || *s == "" {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{String: *s, Valid: true}
}

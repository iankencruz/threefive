// backend/internal/media/service.go
package media

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

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

	// Process media if processor is available
	if s.processor != nil && (isImage || isVideo) {
		if isImage {
			return s.processAndUploadImageVariants(ctx, src, filename, contentType, file.Filename, userID)
		} else if isVideo {
			return s.processAndUploadVideo(ctx, src, filename, contentType, file.Filename, file.Size, userID)
		}
	}

	// Fall back to direct upload
	fmt.Printf("Warning: Processing not available, falling back to direct upload\n")
	return s.directUploadMedia(ctx, src, filename, contentType, file.Filename, file.Size, userID)
}

// processAndUploadImageVariants processes an image into multiple variants and uploads them
func (s *Service) processAndUploadImageVariants(ctx context.Context, src io.Reader, filename, contentType, originalFilename string, userID uuid.UUID) (*sqlc.Media, error) {
	// Process image into variants
	variants, err := s.processor.ProcessImageVariants(ctx, src, filename)
	if err != nil {
		fmt.Printf("Warning: Image variant processing failed: %v\n", err)
		// Fall back to direct upload
		src.(io.Seeker).Seek(0, io.SeekStart)
		return s.directUploadMedia(ctx, src, filename, contentType, originalFilename, 0, userID)
	}

	// Upload all variants
	variantResult, err := s.uploadImageVariants(ctx, variants)
	if err != nil {
		return nil, fmt.Errorf("failed to upload image variants: %w", err)
	}

	// Clean up temporary files
	defer func() {
		os.Remove(variants.Original.ProcessedPath)
		os.Remove(variants.Large.ProcessedPath)
		os.Remove(variants.Medium.ProcessedPath)
		os.Remove(variants.Thumbnail.ProcessedPath)
	}()

	// Determine storage type values
	var s3Bucket, s3Key, s3Region *string
	if s.storage.Type() == storage.StorageTypeS3 {
		bucket := ""
		key := variantResult.LargePath
		region := ""
		s3Bucket = &bucket
		s3Key = &key
		s3Region = &region
	}

	// Create media record with all variant URLs
	media, err := s.queries.CreateMedia(ctx, sqlc.CreateMediaParams{
		Filename:         filepath.Base(variantResult.LargePath),
		OriginalFilename: originalFilename,
		MimeType:         "image/webp",
		SizeBytes:        variantResult.Size,
		Width:            intToNullInt32(variantResult.Width),
		Height:           intToNullInt32(variantResult.Height),
		StorageType:      sqlc.StorageType(s.storage.Type()),
		StoragePath:      variantResult.LargePath,
		S3Bucket:         toNullString(s3Bucket),
		S3Key:            toNullString(s3Key),
		S3Region:         toNullString(s3Region),
		Url:              toNullString(&variantResult.LargeURL), // Deprecated, kept for backwards compat
		OriginalUrl:      toNullString(&variantResult.OriginalURL),
		LargeUrl:         toNullString(&variantResult.LargeURL),
		MediumUrl:        toNullString(&variantResult.MediumURL),
		ThumbnailUrl:     toNullString(&variantResult.ThumbnailURL),
		OriginalPath:     toNullString(&variantResult.OriginalPath),
		LargePath:        toNullString(&variantResult.LargePath),
		MediumPath:       toNullString(&variantResult.MediumPath),
		ThumbnailPath:    toNullString(&variantResult.ThumbnailPath),
		UploadedBy:       userID,
	})
	if err != nil {
		// Clean up uploaded files on error
		s.storage.Delete(ctx, variantResult.OriginalPath)
		s.storage.Delete(ctx, variantResult.LargePath)
		s.storage.Delete(ctx, variantResult.MediumPath)
		s.storage.Delete(ctx, variantResult.ThumbnailPath)
		return nil, errors.Internal("Failed to create media record", err)
	}

	return &media, nil
}

// uploadImageVariants uploads all image variants to storage
func (s *Service) uploadImageVariants(ctx context.Context, variants *processing.ImageVariants) (*VariantUploadResult, error) {
	result := &VariantUploadResult{
		Width:  variants.Large.Width,
		Height: variants.Large.Height,
		Size:   variants.Large.Size,
	}

	// Upload original
	originalFile, err := os.Open(variants.Original.ProcessedPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open original: %w", err)
	}
	defer originalFile.Close()

	originalUpload, err := s.storage.Upload(ctx, storage.UploadInput{
		File:         originalFile,
		Filename:     filepath.Base(variants.Original.ProcessedPath),
		ContentType:  "image/webp",
		Size:         variants.Original.Size,
		GenerateName: true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to upload original: %w", err)
	}
	result.OriginalURL = originalUpload.URL
	result.OriginalPath = originalUpload.Path

	// Upload large
	largeFile, err := os.Open(variants.Large.ProcessedPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open large: %w", err)
	}
	defer largeFile.Close()

	largeUpload, err := s.storage.Upload(ctx, storage.UploadInput{
		File:         largeFile,
		Filename:     filepath.Base(variants.Large.ProcessedPath),
		ContentType:  "image/webp",
		Size:         variants.Large.Size,
		GenerateName: true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to upload large: %w", err)
	}
	result.LargeURL = largeUpload.URL
	result.LargePath = largeUpload.Path

	// Upload medium
	mediumFile, err := os.Open(variants.Medium.ProcessedPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open medium: %w", err)
	}
	defer mediumFile.Close()

	mediumUpload, err := s.storage.Upload(ctx, storage.UploadInput{
		File:         mediumFile,
		Filename:     filepath.Base(variants.Medium.ProcessedPath),
		ContentType:  "image/webp",
		Size:         variants.Medium.Size,
		GenerateName: true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to upload medium: %w", err)
	}
	result.MediumURL = mediumUpload.URL
	result.MediumPath = mediumUpload.Path

	// Upload thumbnail
	thumbnailFile, err := os.Open(variants.Thumbnail.ProcessedPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open thumbnail: %w", err)
	}
	defer thumbnailFile.Close()

	thumbnailUpload, err := s.storage.Upload(ctx, storage.UploadInput{
		File:         thumbnailFile,
		Filename:     filepath.Base(variants.Thumbnail.ProcessedPath),
		ContentType:  "image/webp",
		Size:         variants.Thumbnail.Size,
		GenerateName: true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to upload thumbnail: %w", err)
	}
	result.ThumbnailURL = thumbnailUpload.URL
	result.ThumbnailPath = thumbnailUpload.Path

	return result, nil
}

// processAndUploadVideo processes a video and uploads it
func (s *Service) processAndUploadVideo(ctx context.Context, src io.Reader, filename, contentType, originalFilename string, size int64, userID uuid.UUID) (*sqlc.Media, error) {
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
	defer os.Remove(tempPath)

	// Process video
	processResult, err := s.processor.ProcessVideo(ctx, tempPath, filename)
	if err != nil {
		fmt.Printf("Warning: Video processing failed: %v\n", err)
		// Fall back to direct upload
		srcFile, err := os.Open(tempPath)
		if err != nil {
			return nil, err
		}
		defer srcFile.Close()
		return s.directUploadMedia(ctx, srcFile, filename, contentType, originalFilename, size, userID)
	}

	// Upload processed video
	processedFile, err := os.Open(processResult.ProcessedPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open processed file: %w", err)
	}
	defer processedFile.Close()
	defer os.Remove(processResult.ProcessedPath)

	uploadResult, err := s.storage.Upload(ctx, storage.UploadInput{
		File:         processedFile,
		Filename:     filepath.Base(processResult.ProcessedPath),
		ContentType:  "video/mp4",
		Size:         processResult.Size,
		GenerateName: true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to upload processed video: %w", err)
	}

	// Upload video thumbnail if exists
	var thumbnailURL string
	var thumbnailPath string
	if processResult.ThumbnailPath != "" {
		thumbFile, err := os.Open(processResult.ThumbnailPath)
		if err == nil {
			defer thumbFile.Close()
			defer os.Remove(processResult.ThumbnailPath)

			thumbInfo, _ := thumbFile.Stat()
			thumbResult, err := s.storage.Upload(ctx, storage.UploadInput{
				File:         thumbFile,
				Filename:     filepath.Base(processResult.ThumbnailPath),
				ContentType:  "image/jpeg",
				Size:         thumbInfo.Size(),
				GenerateName: true,
			})
			if err == nil {
				thumbnailURL = thumbResult.URL
				thumbnailPath = thumbResult.Path
			}
		}
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
		OriginalFilename: originalFilename,
		MimeType:         "video/mp4",
		SizeBytes:        uploadResult.Size,
		Width:            intToNullInt32(processResult.Width),
		Height:           intToNullInt32(processResult.Height),
		StorageType:      sqlc.StorageType(s.storage.Type()),
		StoragePath:      uploadResult.Path,
		S3Bucket:         toNullString(s3Bucket),
		S3Key:            toNullString(s3Key),
		S3Region:         toNullString(s3Region),
		Url:              toNullString(&uploadResult.URL),
		LargeUrl:         toNullString(&uploadResult.URL), // Video uses same URL for large
		ThumbnailUrl:     toNullString(&thumbnailURL),
		ThumbnailPath:    toNullString(&thumbnailPath),
		UploadedBy:       userID,
	})
	if err != nil {
		s.storage.Delete(ctx, uploadResult.Path)
		if thumbnailURL != "" {
			s.storage.Delete(ctx, thumbnailPath)
		}
		return nil, errors.Internal("Failed to create media record", err)
	}

	return &media, nil
}

// directUploadMedia uploads a file directly without processing
func (s *Service) directUploadMedia(ctx context.Context, src io.Reader, filename, contentType, originalFilename string, size int64, userID uuid.UUID) (*sqlc.Media, error) {
	uploadResult, err := s.storage.Upload(ctx, storage.UploadInput{
		File:         src,
		Filename:     filename,
		ContentType:  contentType,
		Size:         size,
		GenerateName: true,
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

	media, err := s.queries.CreateMedia(ctx, sqlc.CreateMediaParams{
		Filename:         uploadResult.Filename,
		OriginalFilename: originalFilename,
		MimeType:         contentType,
		SizeBytes:        uploadResult.Size,
		StorageType:      sqlc.StorageType(s.storage.Type()),
		StoragePath:      uploadResult.Path,
		S3Bucket:         toNullString(s3Bucket),
		S3Key:            toNullString(s3Key),
		S3Region:         toNullString(s3Region),
		Url:              toNullString(&uploadResult.URL),
		LargeUrl:         toNullString(&uploadResult.URL),
		UploadedBy:       userID,
	})
	if err != nil {
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

// DeleteMedia soft deletes a media record and removes files from storage
func (s *Service) DeleteMedia(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	media, err := s.queries.GetMediaByID(ctx, id)
	if err != nil {
		return errors.NotFound("Media not found", "media_not_found")
	}

	if media.UploadedBy != userID {
		return errors.Forbidden("Cannot delete media uploaded by another user", "forbidden")
	}

	// Delete all variant files
	if media.OriginalPath.Valid && media.OriginalPath.String != "" {
		s.storage.Delete(ctx, media.OriginalPath.String)
	}
	if media.LargePath.Valid && media.LargePath.String != "" {
		s.storage.Delete(ctx, media.LargePath.String)
	}
	if media.MediumPath.Valid && media.MediumPath.String != "" {
		s.storage.Delete(ctx, media.MediumPath.String)
	}
	if media.ThumbnailPath.Valid && media.ThumbnailPath.String != "" {
		s.storage.Delete(ctx, media.ThumbnailPath.String)
	}

	// Delete from old storage_path field for backwards compatibility
	if media.StoragePath != "" {
		s.storage.Delete(ctx, media.StoragePath)
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

	// Delete all variant files
	if media.OriginalPath.Valid && media.OriginalPath.String != "" {
		s.storage.Delete(ctx, media.OriginalPath.String)
	}
	if media.LargePath.Valid && media.LargePath.String != "" {
		s.storage.Delete(ctx, media.LargePath.String)
	}
	if media.MediumPath.Valid && media.MediumPath.String != "" {
		s.storage.Delete(ctx, media.MediumPath.String)
	}
	if media.ThumbnailPath.Valid && media.ThumbnailPath.String != "" {
		s.storage.Delete(ctx, media.ThumbnailPath.String)
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

// validateFile validates uploaded file
func (s *Service) validateFile(file *multipart.FileHeader) error {
	maxSize := int64(100 << 20) // 100MB
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

// Search Media
func (s *Service) SearchMedia(ctx context.Context, params SearchMediaParams) ([]sqlc.Media, int64, error) {
	// Get total count with SQLC's generated params
	total, err := s.queries.CountSearchMedia(ctx, sqlc.CountSearchMediaParams{
		SearchQuery:    params.SearchQuery,
		MimeTypeFilter: params.MimeTypeFilter,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count media: %w", err)
	}

	// Get filtered media with SQLC's generated params
	media, err := s.queries.SearchMedia(ctx, sqlc.SearchMediaParams{
		SearchQuery:    params.SearchQuery,
		MimeTypeFilter: params.MimeTypeFilter,
		SortBy:         params.SortBy,
		SortOrder:      params.SortOrder,
		LimitVal:       params.Limit,
		OffsetVal:      params.Offset,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search media: %w", err)
	}

	return media, total, nil
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

package media

import (
	"context"
	"mime/multipart"

	"github.com/iankencruz/threefive/backend/internal/generated"
)

type S3Uploader interface {
	Upload(ctx context.Context, file multipart.File, filename, contentType string) (url string, thumb string, err error)
}

type Service struct {
	repo     Repository
	uploader S3Uploader
}

func NewService(repo Repository, uploader S3Uploader) *Service {
	return &Service{repo: repo, uploader: uploader}
}

func (s *Service) UploadMedia(
	ctx context.Context,
	file multipart.File,
	header *multipart.FileHeader,
	title string,
	altText string,
	sortOrder int32,
) (*generated.Media, error) {
	defer file.Close()

	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	// Upload to S3
	url, thumbURL, err := s.uploader.Upload(ctx, file, header.Filename, contentType)
	if err != nil {
		return nil, err
	}

	fileSize := int32(header.Size)

	// Insert into DB
	media, err := s.repo.Create(ctx, generated.CreateMediaParams{
		Url:          url,
		ThumbnailUrl: &thumbURL,
		Type:         "image", // placeholder, you can auto-detect later
		Title:        &title,
		AltText:      &altText,
		MimeType:     &contentType,
		FileSize:     &fileSize,
		SortOrder:    sortOrder,
	})
	if err != nil {
		return nil, err
	}

	return media, nil
}

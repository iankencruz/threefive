package media

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/iankencruz/threefive/internal/core/s3"
	"github.com/iankencruz/threefive/internal/generated"
)

type S3Uploader interface {
	PutObject(ctx context.Context, reader io.Reader, filename, contentType string) (string, error)
	JoinURL(filename string) string
	RemoveObject(ctx context.Context, key string) error
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

	buf, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	isJPG := strings.HasSuffix(strings.ToLower(header.Filename), ".jpg") || strings.HasSuffix(strings.ToLower(header.Filename), ".jpeg")
	isPNG := strings.HasSuffix(strings.ToLower(header.Filename), ".png")

	var mainFilename string
	var mainURL string
	var thumbURL, mediumURL *string
	mime := contentType

	if isJPG || isPNG {

		variants, err := s3.GenerateVariants(header.Filename, buf)
		if err != nil {
			return nil, err
		}

		mainS3Key, err := s.uploader.PutObject(ctx, bytes.NewReader(variants.WebP), variants.Filename, "image/webp")
		if err != nil {
			return nil, err
		}
		mainURL = s.uploader.JoinURL(mainS3Key)

		thumbS3, err := s.uploader.PutObject(ctx, bytes.NewReader(variants.Thumbnail), variants.ThumbName, "image/webp")
		if err == nil {
			tmp := s.uploader.JoinURL(thumbS3)
			thumbURL = &tmp
		}

		mediumS3, err := s.uploader.PutObject(ctx, bytes.NewReader(variants.Medium), variants.MediumName, "image/webp")
		if err == nil {
			tmp := s.uploader.JoinURL(mediumS3)
			mediumURL = &tmp
		}
		mime = "image/webp"
	} else {

		baseName := strings.TrimSuffix(header.Filename, filepath.Ext(header.Filename))
		mainFilename = fmt.Sprintf("%s-%s%s", uuid.New().String(), baseName, filepath.Ext(header.Filename))
		key, err := s.uploader.PutObject(ctx, bytes.NewReader(buf), mainFilename, contentType)
		if err != nil {
			return nil, err
		}
		mainURL = s.uploader.JoinURL(key)
	}

	fullURL := mainURL
	size := int32(len(buf))
	mediaType := inferMediaType(mime)

	media, err := s.repo.Create(ctx, generated.CreateMediaParams{
		Url:          fullURL,
		ThumbnailUrl: thumbURL,
		MediumUrl:    mediumURL,
		Type:         mediaType,
		Title:        &title,
		AltText:      &title,
		MimeType:     &mime,
		FileSize:     &size,
		SortOrder:    sortOrder,
	})
	if err != nil {
		return nil, err
	}

	return media, nil
}

func (s *Service) DeleteMediaWithVariants(ctx context.Context, media *generated.Media) error {
	// Remove from DB first
	if err := s.repo.Delete(ctx, media.ID); err != nil {
		return err
	}

	// Get object keys by trimming the base URL
	baseURL := s.uploader.JoinURL("") // ensures consistent formatting
	stripPrefix := strings.TrimSuffix(baseURL, "/")

	trim := func(fullURL *string) (key string, ok bool) {
		if fullURL != nil && strings.HasPrefix(*fullURL, stripPrefix) {
			return strings.TrimPrefix(*fullURL, stripPrefix+"/"), true
		}
		return "", false
	}

	keys := []string{}
	if key, ok := trim(&media.Url); ok {
		keys = append(keys, key)
	}
	if key, ok := trim(media.ThumbnailUrl); ok {
		keys = append(keys, key)
	}
	if key, ok := trim(media.MediumUrl); ok {
		keys = append(keys, key)
	}

	// Delete from S3
	for _, key := range keys {
		if err := s.uploader.RemoveObject(ctx, key); err != nil {
			// Optional: log warning instead of failing the whole process
			fmt.Printf("⚠️ failed to delete S3 object %s: %v\n", key, err)
		}
	}

	return nil
}

func inferMediaType(mimeType string) string {
	switch {
	case mimeType == "image/gif":
		return "image" // ✅ special case: let GIFs behave like images
	case strings.HasPrefix(mimeType, "image/"):
		return "image"
	case strings.HasPrefix(mimeType, "video/"):
		return "video"
	case mimeType == "video/webm":
		return "video"
	default:
		return "embed"
	}
}

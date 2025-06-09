package media

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/iankencruz/threefive/backend/internal/core/s3"
	"github.com/iankencruz/threefive/backend/internal/generated"
	"github.com/jackc/pgx/v5/pgtype"
)

type S3Uploader interface {
	PutObject(ctx context.Context, reader io.Reader, filename, contentType string) (string, error)
	JoinURL(filename string) string
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
	var webpBytes []byte
	var thumbURL, mediumURL *string
	mime := contentType

	if isJPG || isPNG {
		webpBytes, err = s3.ConvertToWebP(buf)
		if err != nil {
			return nil, err
		}

		baseName := strings.TrimSuffix(header.Filename, filepath.Ext(header.Filename))
		mainFilename = fmt.Sprintf("%s-%s.webp", uuid.New().String(), baseName)
		key, err := s.uploader.PutObject(ctx, bytes.NewReader(webpBytes), mainFilename, "image/webp")
		if err != nil {
			return nil, err
		}
		mainURL = s.uploader.JoinURL(key)

		// Thumbnail
		thumb, err := s3.ResizeImage(webpBytes, 0.25)
		if err == nil {
			thumbName := "thumb-" + mainFilename
			thumbS3Key, err := s.uploader.PutObject(ctx, bytes.NewReader(thumb), thumbName, "image/webp")
			if err == nil {
				thumb := s.uploader.JoinURL(thumbS3Key)
				thumbURL = &thumb
			}
		}

		// Medium
		medium, err := s3.ResizeImage(webpBytes, 0.5)
		if err == nil {
			mediumName := "medium-" + mainFilename
			mediumS3Key, err := s.uploader.PutObject(ctx, bytes.NewReader(medium), mediumName, "image/webp")
			if err == nil {
				medium := s.uploader.JoinURL(mediumS3Key)
				mediumURL = &medium
			}
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
		AltText:      &altText,
		MimeType:     &mime,
		FileSize:     &size,
		SortOrder:    sortOrder,
	})
	if err != nil {
		return nil, err
	}

	return media, nil
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

func parseUUIDParam(r *http.Request, id string) (pgtype.UUID, error) {
	u, err := uuid.Parse(id)
	if err != nil {
		return pgtype.UUID{}, err
	}
	return pgtype.UUID{
		Bytes: u,
		Valid: true,
	}, nil
}

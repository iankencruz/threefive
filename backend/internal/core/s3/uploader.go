package s3

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Uploader struct {
	Client     *minio.Client
	BucketName string
	BaseURL    string // e.g. https://s3.yourdomain.com/media/
}

func NewUploader(endpoint, accessKey, secretKey, bucket string, useSSL bool, baseURL string) (*Uploader, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, err
	}

	return &Uploader{
		Client:     client,
		BucketName: bucket,
		BaseURL:    baseURL,
	}, nil
}

// Upload uploads the file and returns (fileURL, thumbnailURL, error)
func (u *Uploader) Upload(ctx context.Context, file io.Reader, filename string, contentType string) (string, string, error) {
	objectName := fmt.Sprintf("%d-%s", time.Now().UnixNano(), filepath.Base(filename))

	_, err := u.Client.PutObject(ctx, u.BucketName, objectName, file, -1, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", "", err
	}

	fileURL := fmt.Sprintf("%s/%s", u.BaseURL, objectName)

	// optional: thumbnail logic
	thumbnailURL := "" // implement later if needed

	return fileURL, thumbnailURL, nil
}

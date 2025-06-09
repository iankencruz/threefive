package s3

import (
	"bytes"
	"context"
	"io"
	"log"
	"log/slog"
	"testing"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	tcminio "github.com/testcontainers/testcontainers-go/modules/minio"
)

func TestUploader_Upload(t *testing.T) {
	ctx := context.Background()
	logger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{})) // default to discard

	// Capture logs with `t.Log`
	logger = slog.New(slog.NewTextHandler(slogTestWriter{t}, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelInfo,
	}))

	minioContainer, err := tcminio.Run(ctx, "minio/minio:RELEASE.2024-01-16T16-07-38Z")

	if err != nil {
		log.Printf("failed to start container: %s", err)
		return
	}

	t.Cleanup(func() {
		_ = minioContainer.Terminate(ctx)
	})

	endpoint, err := minioContainer.Endpoint(ctx, "")
	require.NoError(t, err)

	const bucket = "media"
	baseURL := "http://" + endpoint + "/" + bucket

	logger.Info("MinIO started", "endpoint", endpoint)

	uploader, err := NewUploader(endpoint, "minioadmin", "minioadmin", bucket, "sgp1", false, baseURL)
	require.NoError(t, err)

	err = uploader.Client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{})
	require.NoError(t, err)
	logger.Info("Created test bucket", "bucket", bucket)

	tests := []struct {
		name        string
		fileContent string
		fileName    string
		contentType string
		wantErr     bool
	}{
		{"basic text file", "Hello, MinIO!", "hello.txt", "text/plain", false},
		{"empty file", "", "empty.txt", "text/plain", false},
		{"image content", "\xFF\xD8\xFF", "fake.jpg", "image/jpeg", false},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			logger.Info("Running test case", "name", tt.name)

			data := []byte(tt.fileContent)
			objectName := uuid.New().String() + "-" + tt.fileName
			reader := bytes.NewReader(data)

			logger.Info("Uploading file", "object", objectName)

			url, thumb, err := uploader.Upload(ctx, reader, objectName, tt.contentType)

			if tt.wantErr {
				assert.Error(t, err)
				logger.Info("Expected error occurred", "err", err)
				return
			}

			require.NoError(t, err)
			logger.Info("Upload complete", "url", url, "thumb", thumb)

			obj, err := uploader.Client.GetObject(ctx, bucket, objectName, minio.GetObjectOptions{})
			require.NoError(t, err)

			got, err := io.ReadAll(obj)
			require.NoError(t, err)
			assert.Equal(t, tt.fileContent, string(got))

			logger.Info("Verified object content matches")
		})
	}
}

// slogTestWriter allows slog logs to be routed through t.Log()
type slogTestWriter struct {
	t *testing.T
}

func (w slogTestWriter) Write(p []byte) (n int, err error) {
	w.t.Log(string(p))
	return len(p), nil
}

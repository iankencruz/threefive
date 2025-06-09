package s3

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"strings"

	"github.com/chai2010/webp"
	"github.com/disintegration/imaging"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Uploader struct {
	Client     *minio.Client
	BucketName string
	BaseURL    string // e.g. https://sgp1.vultrobjects.com/
}

// ✅ Uploads raw data to S3 and returns full public URL

func (u *Uploader) PutObject(ctx context.Context, data io.Reader, filename string, contentType string) (string, error) {
	objectName := filename

	// We need to read the data into a buffer to calculate its size
	buf, err := io.ReadAll(data)
	if err != nil {
		return "", err
	}

	// Upload using the minio client
	_, err = u.Client.PutObject(ctx, u.BucketName, objectName, bytes.NewReader(buf), int64(len(buf)), minio.PutObjectOptions{
		ContentType: contentType,
		UserMetadata: map[string]string{
			"x-amz-acl": "public-read",
		},
	})
	if err != nil {
		return "", err
	}

	// Return the public URL
	return objectName, nil
}

// ✅ Constructs full public URL
func (u *Uploader) JoinURL(filename string) string {
	base := strings.TrimSuffix(u.BaseURL, "/")
	bucket := strings.TrimSuffix(u.BucketName, "/")
	return fmt.Sprintf("%s/%s/%s", base, bucket, filename)
}

// ✅ Convert PNG or JPEG to WebP
func ConvertToWebP(data []byte) ([]byte, error) {
	img, format, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	if format != "jpeg" && format != "png" {
		return data, nil // skip conversion
	}

	var buf bytes.Buffer
	if err := webp.Encode(&buf, img, &webp.Options{Lossless: true}); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// ✅ Resize image to width in px (e.g., 300 for thumbnail)
func ResizeImage(input []byte, scale float64) ([]byte, error) {
	if scale <= 0 || scale >= 1 {
		return nil, errors.New("scale must be between 0 and 1")
	}

	img, _, err := image.Decode(bytes.NewReader(input))
	if err != nil {
		return nil, err
	}

	originalWidth := img.Bounds().Dx()
	targetWidth := int(float64(originalWidth) * scale)

	// 🛑 Prevent upscaling
	if targetWidth >= originalWidth {
		return input, nil // return original image
	}

	// Resize
	resized := imaging.Resize(img, targetWidth, 0, imaging.Lanczos)

	// Encode to WebP
	var buf bytes.Buffer
	err = webp.Encode(&buf, resized, &webp.Options{
		Lossless: false,
		Quality:  75,
	})
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// ✅ Optional: Create a new uploader instance
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

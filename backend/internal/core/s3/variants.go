package s3

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

type Variants struct {
	Filename   string
	WebP       []byte
	Thumbnail  []byte
	Medium     []byte
	ThumbName  string
	MediumName string
}

func GenerateVariants(originalFilename string, image []byte) (*Variants, error) {
	webpBytes, err := ConvertToWebP(image)
	if err != nil {
		return nil, err
	}

	ext := filepath.Ext(originalFilename)
	name := strings.TrimSuffix(originalFilename, ext)
	base := fmt.Sprintf("%s-%s", uuid.New().String(), name)

	thumb, err := ResizeImage(webpBytes, 0.25)
	if err != nil {
		return nil, err
	}

	medium, err := ResizeImage(webpBytes, 0.5)
	if err != nil {
		return nil, err
	}

	return &Variants{
		Filename:   base + ".webp",
		WebP:       webpBytes,
		Thumbnail:  thumb,
		Medium:     medium,
		ThumbName:  "thumb-" + base + ".webp",
		MediumName: "medium-" + base + ".webp",
	}, nil
}

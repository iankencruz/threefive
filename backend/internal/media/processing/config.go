// backend/internal/media/processing/config.go
package processing

import (
	"os"
	"strconv"
)

// LoadConfigFromEnv loads processor configuration from environment variables
func LoadConfigFromEnv() ProcessorConfig {
	config := DefaultConfig()

	if v := os.Getenv("MEDIA_WEBP_QUALITY"); v != "" {
		if quality, err := strconv.Atoi(v); err == nil && quality >= 1 && quality <= 100 {
			config.WebPQuality = quality
		}
	}

	if v := os.Getenv("MEDIA_MAX_IMAGE_WIDTH"); v != "" {
		if width, err := strconv.Atoi(v); err == nil && width > 0 {
			config.MaxImageWidth = width
		}
	}

	if v := os.Getenv("MEDIA_MAX_IMAGE_HEIGHT"); v != "" {
		if height, err := strconv.Atoi(v); err == nil && height > 0 {
			config.MaxImageHeight = height
		}
	}

	if v := os.Getenv("MEDIA_GENERATE_THUMBNAIL"); v != "" {
		if generate, err := strconv.ParseBool(v); err == nil {
			config.GenerateThumbnail = generate
		}
	}

	if v := os.Getenv("MEDIA_VIDEO_CODEC"); v != "" {
		config.VideoCodec = v
	}

	if v := os.Getenv("MEDIA_VIDEO_BITRATE"); v != "" {
		config.VideoBitrate = v
	}

	if v := os.Getenv("MEDIA_VIDEO_PRESET"); v != "" {
		config.VideoPreset = v
	}

	if v := os.Getenv("MEDIA_VIDEO_MAX_WIDTH"); v != "" {
		if width, err := strconv.Atoi(v); err == nil && width > 0 {
			config.VideoMaxWidth = width
		}
	}

	if v := os.Getenv("MEDIA_VIDEO_MAX_HEIGHT"); v != "" {
		if height, err := strconv.Atoi(v); err == nil && height > 0 {
			config.VideoMaxHeight = height
		}
	}

	if v := os.Getenv("MEDIA_THUMBNAIL_TIME"); v != "" {
		config.ThumbnailTime = v
	}

	if v := os.Getenv("MEDIA_THUMBNAIL_WIDTH"); v != "" {
		if width, err := strconv.Atoi(v); err == nil && width > 0 {
			config.ThumbnailWidth = width
		}
	}

	if v := os.Getenv("MEDIA_THUMBNAIL_HEIGHT"); v != "" {
		if height, err := strconv.Atoi(v); err == nil && height > 0 {
			config.ThumbnailHeight = height
		}
	}

	return config
}

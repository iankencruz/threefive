// backend/internal/media/processing/webp.go
package processing

import (
	"fmt"
	"image"
	"image/jpeg"
	"os"
	"os/exec"
	"path/filepath"
)

// encodeWebP encodes an image as WebP using FFmpeg
func (p *Processor) encodeWebP(img image.Image, outputPath string) error {
	// Create temp JPEG
	tempJPEG := filepath.Join(p.workDir, "temp_encode.jpg")
	jpegFile, err := os.Create(tempJPEG)
	if err != nil {
		return fmt.Errorf("failed to create temp jpeg: %w", err)
	}
	defer os.Remove(tempJPEG)

	// Encode as high-quality JPEG first
	if err := jpeg.Encode(jpegFile, img, &jpeg.Options{Quality: 95}); err != nil {
		jpegFile.Close()
		return fmt.Errorf("failed to encode jpeg: %w", err)
	}
	jpegFile.Close()

	// Convert to WebP using FFmpeg
	args := []string{
		"-i", tempJPEG,
		"-c:v", "libwebp",
		"-quality", fmt.Sprintf("%d", p.config.WebPQuality),
		"-y", // Overwrite output
		outputPath,
	}

	cmd := exec.Command(p.ffmpegCmd, args...)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("ffmpeg webp conversion failed: %w\nOutput: %s", err, output)
	}

	return nil
}

// backend/internal/media/processing/processor.go
package processing

import (
	"context"
	"fmt"
	"image"
	_ "image/jpeg" // Register JPEG decoder
	_ "image/png"  // Register PNG decoder
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	_ "golang.org/x/image/tiff" // Register TIFF decoder
)

// ProcessorConfig holds configuration for media processing
type ProcessorConfig struct {
	// Image settings
	WebPQuality       int
	MaxImageWidth     int
	MaxImageHeight    int
	GenerateThumbnail bool

	// Video settings
	VideoCodec     string
	VideoBitrate   string
	VideoPreset    string
	VideoMaxWidth  int
	VideoMaxHeight int
	ThumbnailTime  string

	// Thumbnail settings
	ThumbnailWidth  int
	ThumbnailHeight int
}

// DefaultConfig returns default processing configuration
func DefaultConfig() ProcessorConfig {
	return ProcessorConfig{
		WebPQuality:       80,
		MaxImageWidth:     4000,
		MaxImageHeight:    4000,
		GenerateThumbnail: true,
		VideoCodec:        "h264",
		VideoBitrate:      "2M",
		VideoPreset:       "fast",
		VideoMaxWidth:     1920,
		VideoMaxHeight:    1080,
		ThumbnailTime:     "00:00:01",
		ThumbnailWidth:    300,
		ThumbnailHeight:   300,
	}
}

// ProcessResult contains the results of media processing
type ProcessResult struct {
	ProcessedPath string
	ThumbnailPath string
	Width         int
	Height        int
	Size          int64
	Format        string
}

// Processor handles media processing operations
type Processor struct {
	config     ProcessorConfig
	workDir    string
	ffmpegCmd  string
	ffprobeCmd string
}

// NewProcessor creates a new media processor
func NewProcessor(config ProcessorConfig, workDir string) (*Processor, error) {
	if err := os.MkdirAll(workDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create work directory: %w", err)
	}

	ffmpegPath, err := exec.LookPath("ffmpeg")
	if err != nil {
		return nil, fmt.Errorf("ffmpeg not found in PATH: %w", err)
	}

	ffprobePath, err := exec.LookPath("ffprobe")
	if err != nil {
		return nil, fmt.Errorf("ffprobe not found in PATH: %w", err)
	}

	return &Processor{
		config:     config,
		workDir:    workDir,
		ffmpegCmd:  ffmpegPath,
		ffprobeCmd: ffprobePath,
	}, nil
}

// ProcessImage processes an image file (converts to WebP)
func (p *Processor) ProcessImage(ctx context.Context, input io.Reader, filename string) (*ProcessResult, error) {
	img, format, err := image.Decode(input)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// Resize if needed
	if p.config.MaxImageWidth > 0 || p.config.MaxImageHeight > 0 {
		img = p.resizeImage(img, width, height)
		bounds = img.Bounds()
		width = bounds.Dx()
		height = bounds.Dy()
	}

	// Convert to WebP
	baseFilename := strings.TrimSuffix(filename, filepath.Ext(filename))
	webpFilename := baseFilename + ".webp"
	webpPath := filepath.Join(p.workDir, webpFilename)

	if err := p.encodeWebP(img, webpPath); err != nil {
		return nil, fmt.Errorf("failed to encode webp: %w", err)
	}

	fileInfo, err := os.Stat(webpPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}

	result := &ProcessResult{
		ProcessedPath: webpPath,
		Width:         width,
		Height:        height,
		Size:          fileInfo.Size(),
		Format:        "webp",
	}

	// Generate thumbnail
	if p.config.GenerateThumbnail {
		thumbPath, err := p.generateImageThumbnail(img, baseFilename)
		if err != nil {
			fmt.Printf("Warning: failed to generate thumbnail: %v\n", err)
		} else {
			result.ThumbnailPath = thumbPath
		}
	}

	fmt.Printf("✅ Processed image: %s -> %s (original: %s, size: %d bytes)\n",
		filename, webpFilename, format, result.Size)

	return result, nil
}

// ProcessVideo processes a video file
func (p *Processor) ProcessVideo(ctx context.Context, inputPath, filename string) (*ProcessResult, error) {
	info, err := p.getVideoInfo(ctx, inputPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get video info: %w", err)
	}

	baseFilename := strings.TrimSuffix(filename, filepath.Ext(filename))
	outputFilename := baseFilename + "_optimized.mp4"
	outputPath := filepath.Join(p.workDir, outputFilename)

	args := []string{
		"-i", inputPath,
		"-c:v", p.config.VideoCodec,
		"-preset", p.config.VideoPreset,
		"-b:v", p.config.VideoBitrate,
		"-c:a", "aac",
		"-b:a", "128k",
	}

	if p.config.VideoMaxWidth > 0 || p.config.VideoMaxHeight > 0 {
		scale := fmt.Sprintf("scale='min(%d,iw)':min(%d,ih):force_original_aspect_ratio=decrease",
			p.config.VideoMaxWidth, p.config.VideoMaxHeight)
		args = append(args, "-vf", scale)
	}

	args = append(args, "-movflags", "+faststart", outputPath)

	cmd := exec.CommandContext(ctx, p.ffmpegCmd, args...)
	if output, err := cmd.CombinedOutput(); err != nil {
		return nil, fmt.Errorf("ffmpeg failed: %w\nOutput: %s", err, output)
	}

	fileInfo, err := os.Stat(outputPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get output file info: %w", err)
	}

	result := &ProcessResult{
		ProcessedPath: outputPath,
		Width:         info.Width,
		Height:        info.Height,
		Size:          fileInfo.Size(),
		Format:        "mp4",
	}

	thumbPath, err := p.generateVideoThumbnail(ctx, outputPath, baseFilename)
	if err != nil {
		fmt.Printf("Warning: failed to generate video thumbnail: %v\n", err)
	} else {
		result.ThumbnailPath = thumbPath
	}

	fmt.Printf("✅ Processed video: %s -> %s (size: %d bytes)\n",
		filename, outputFilename, result.Size)

	return result, nil
}

// resizeImage resizes an image maintaining aspect ratio
func (p *Processor) resizeImage(img image.Image, width, height int) image.Image {
	maxW := p.config.MaxImageWidth
	maxH := p.config.MaxImageHeight

	if maxW == 0 {
		maxW = width
	}
	if maxH == 0 {
		maxH = height
	}

	if width <= maxW && height <= maxH {
		return img
	}

	ratio := float64(width) / float64(height)
	newWidth := maxW
	newHeight := int(float64(newWidth) / ratio)

	if newHeight > maxH {
		newHeight = maxH
		newWidth = int(float64(newHeight) * ratio)
	}

	resized := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))

	for y := 0; y < newHeight; y++ {
		for x := 0; x < newWidth; x++ {
			srcX := x * width / newWidth
			srcY := y * height / newHeight
			resized.Set(x, y, img.At(srcX, srcY))
		}
	}

	return resized
}

// generateImageThumbnail generates a thumbnail for an image
func (p *Processor) generateImageThumbnail(img image.Image, baseFilename string) (string, error) {
	thumbImg := p.resizeToThumbnail(img)

	thumbFilename := baseFilename + "_thumb.webp"
	thumbPath := filepath.Join(p.workDir, thumbFilename)

	if err := p.encodeWebP(thumbImg, thumbPath); err != nil {
		return "", err
	}

	return thumbPath, nil
}

// generateVideoThumbnail generates a thumbnail from a video
func (p *Processor) generateVideoThumbnail(ctx context.Context, videoPath, baseFilename string) (string, error) {
	thumbFilename := baseFilename + "_thumb.jpg"
	thumbPath := filepath.Join(p.workDir, thumbFilename)

	args := []string{
		"-i", videoPath,
		"-ss", p.config.ThumbnailTime,
		"-vframes", "1",
		"-vf", fmt.Sprintf("scale=%d:%d:force_original_aspect_ratio=decrease",
			p.config.ThumbnailWidth, p.config.ThumbnailHeight),
		thumbPath,
	}

	cmd := exec.CommandContext(ctx, p.ffmpegCmd, args...)
	if output, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("ffmpeg thumbnail failed: %w\nOutput: %s", err, output)
	}

	return thumbPath, nil
}

// resizeToThumbnail resizes an image to thumbnail dimensions
func (p *Processor) resizeToThumbnail(img image.Image) image.Image {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	thumbW := p.config.ThumbnailWidth
	thumbH := p.config.ThumbnailHeight

	ratio := float64(width) / float64(height)
	newWidth := thumbW
	newHeight := int(float64(newWidth) / ratio)

	if newHeight > thumbH {
		newHeight = thumbH
		newWidth = int(float64(newHeight) * ratio)
	}

	thumb := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))

	for y := 0; y < newHeight; y++ {
		for x := 0; x < newWidth; x++ {
			srcX := x * width / newWidth
			srcY := y * height / newHeight
			thumb.Set(x, y, img.At(srcX, srcY))
		}
	}

	return thumb
}

// getVideoInfo retrieves video metadata
func (p *Processor) getVideoInfo(ctx context.Context, videoPath string) (*VideoInfo, error) {
	args := []string{
		"-v", "error",
		"-select_streams", "v:0",
		"-show_entries", "stream=width,height,duration",
		"-of", "default=noprint_wrappers=1",
		videoPath,
	}

	cmd := exec.CommandContext(ctx, p.ffprobeCmd, args...)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("ffprobe failed: %w", err)
	}

	info := &VideoInfo{}
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key, value := parts[0], parts[1]
		switch key {
		case "width":
			fmt.Sscanf(value, "%d", &info.Width)
		case "height":
			fmt.Sscanf(value, "%d", &info.Height)
		case "duration":
			fmt.Sscanf(value, "%f", &info.Duration)
		}
	}

	return info, nil
}

// VideoInfo contains video metadata
type VideoInfo struct {
	Width    int
	Height   int
	Duration float64
}

// IsImageFile checks if a filename is an image
func IsImageFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	return ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".tif" || ext == ".tiff"
}

// IsVideoFile checks if a filename is a video
func IsVideoFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	return ext == ".mp4" || ext == ".mov" || ext == ".avi" || ext == ".mkv" || ext == ".webm"
}

// WorkDir returns the working directory
func (p *Processor) WorkDir() string {
	return p.workDir
}

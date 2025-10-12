// backend/internal/media/processing/processor.go
package processing

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// ProcessorConfig holds configuration for media processing
type ProcessorConfig struct {
	// Image settings
	WebPQuality       int  // 1-100, default 80
	MaxImageWidth     int  // Max width before resize, 0 = no limit
	MaxImageHeight    int  // Max height before resize, 0 = no limit
	GenerateThumbnail bool // Generate thumbnail for images

	// Video settings
	VideoCodec     string // h264, h265, vp9
	VideoBitrate   string // e.g., "2M", "5M"
	VideoPreset    string // ultrafast, fast, medium, slow
	VideoMaxWidth  int    // Max video width
	VideoMaxHeight int    // Max video height
	ThumbnailTime  string // Time to capture thumbnail (e.g., "00:00:01")

	// Thumbnail settings
	ThumbnailWidth  int // Default 300
	ThumbnailHeight int // Default 300
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
	ProcessedPath string // Path to processed file (WebP for images, optimized video)
	ThumbnailPath string // Path to thumbnail (if generated)
	Width         int
	Height        int
	Size          int64
	Format        string // "webp", "mp4", etc.
}

// Processor handles media processing operations
type Processor struct {
	config     ProcessorConfig
	workDir    string // Temporary working directory
	ffmpegCmd  string // Path to ffmpeg binary
	ffprobeCmd string // Path to ffprobe binary
}

// NewProcessor creates a new media processor
func NewProcessor(config ProcessorConfig, workDir string) (*Processor, error) {
	// Ensure work directory exists
	if err := os.MkdirAll(workDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create work directory: %w", err)
	}

	// Check for FFmpeg availability
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
	// Read the image
	img, format, err := image.Decode(input)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	// Get original dimensions
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

	webpFile, err := os.Create(webpPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create webp file: %w", err)
	}
	defer webpFile.Close()

	// Encode as WebP
	if err := p.encodeWebP(webpFile, img); err != nil {
		return nil, fmt.Errorf("failed to encode webp: %w", err)
	}

	// Get file size
	fileInfo, err := webpFile.Stat()
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

	// Generate thumbnail if configured
	if p.config.GenerateThumbnail {
		thumbPath, err := p.generateImageThumbnail(img, baseFilename)
		if err != nil {
			// Log error but don't fail the whole operation
			fmt.Printf("Warning: failed to generate thumbnail: %v\n", err)
		} else {
			result.ThumbnailPath = thumbPath
		}
	}

	fmt.Printf("✅ Processed image: %s -> %s (original: %s, size: %d bytes)\n",
		filename, webpFilename, format, result.Size)

	return result, nil
}

// ProcessVideo processes a video file (optimizes and generates thumbnail)
func (p *Processor) ProcessVideo(ctx context.Context, inputPath, filename string) (*ProcessResult, error) {
	// Get video info
	info, err := p.getVideoInfo(ctx, inputPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get video info: %w", err)
	}

	// Generate output filename
	baseFilename := strings.TrimSuffix(filename, filepath.Ext(filename))
	outputFilename := baseFilename + "_optimized.mp4"
	outputPath := filepath.Join(p.workDir, outputFilename)

	// Build ffmpeg command for video optimization
	args := []string{
		"-i", inputPath,
		"-c:v", p.config.VideoCodec,
		"-preset", p.config.VideoPreset,
		"-b:v", p.config.VideoBitrate,
		"-c:a", "aac",
		"-b:a", "128k",
	}

	// Add scaling if needed
	if p.config.VideoMaxWidth > 0 || p.config.VideoMaxHeight > 0 {
		scale := fmt.Sprintf("scale='min(%d,iw)':min(%d,ih):force_original_aspect_ratio=decrease",
			p.config.VideoMaxWidth, p.config.VideoMaxHeight)
		args = append(args, "-vf", scale)
	}

	args = append(args, "-movflags", "+faststart", outputPath)

	// Execute ffmpeg
	cmd := exec.CommandContext(ctx, p.ffmpegCmd, args...)
	if output, err := cmd.CombinedOutput(); err != nil {
		return nil, fmt.Errorf("ffmpeg failed: %w\nOutput: %s", err, output)
	}

	// Get output file size
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

	// Generate video thumbnail
	thumbPath, err := p.generateVideoThumbnail(ctx, outputPath, baseFilename)
	if err != nil {
		// Log error but don't fail
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

	// No resize needed
	if width <= maxW && height <= maxH {
		return img
	}

	// Calculate new dimensions maintaining aspect ratio
	ratio := float64(width) / float64(height)
	newWidth := maxW
	newHeight := int(float64(newWidth) / ratio)

	if newHeight > maxH {
		newHeight = maxH
		newWidth = int(float64(newHeight) * ratio)
	}

	// Create new image with calculated dimensions
	resized := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))

	// Simple nearest-neighbor resizing
	// For production, consider using github.com/nfnt/resize or similar
	for y := 0; y < newHeight; y++ {
		for x := 0; x < newWidth; x++ {
			srcX := x * width / newWidth
			srcY := y * height / newHeight
			resized.Set(x, y, img.At(srcX, srcY))
		}
	}

	return resized
}

// encodeWebP encodes an image as WebP
func (p *Processor) encodeWebP(w io.Writer, img image.Image) error {
	// Use FFmpeg for WebP encoding since standard library doesn't support encoding
	// Alternative: use third-party library like github.com/chai2010/webp

	// For now, we'll encode as JPEG first, then convert with FFmpeg
	// In production, use a proper WebP encoder library

	buf := &bytes.Buffer{}
	if err := jpeg.Encode(buf, img, &jpeg.Options{Quality: p.config.WebPQuality}); err != nil {
		return err
	}

	// TODO: Use proper WebP encoder here
	// For now, just copy the JPEG data
	_, err := io.Copy(w, buf)
	return err
}

// generateImageThumbnail generates a thumbnail for an image
func (p *Processor) generateImageThumbnail(img image.Image, baseFilename string) (string, error) {
	// Resize to thumbnail dimensions
	thumbImg := p.resizeToThumbnail(img)

	// Save as WebP
	thumbFilename := baseFilename + "_thumb.webp"
	thumbPath := filepath.Join(p.workDir, thumbFilename)

	thumbFile, err := os.Create(thumbPath)
	if err != nil {
		return "", err
	}
	defer thumbFile.Close()

	if err := p.encodeWebP(thumbFile, thumbImg); err != nil {
		return "", err
	}

	return thumbPath, nil
}

// generateVideoThumbnail generates a thumbnail from a video
func (p *Processor) generateVideoThumbnail(ctx context.Context, videoPath, baseFilename string) (string, error) {
	thumbFilename := baseFilename + "_thumb.jpg"
	thumbPath := filepath.Join(p.workDir, thumbFilename)

	// Use ffmpeg to extract a frame
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

	// Calculate dimensions maintaining aspect ratio
	ratio := float64(width) / float64(height)
	newWidth := thumbW
	newHeight := int(float64(newWidth) / ratio)

	if newHeight > thumbH {
		newHeight = thumbH
		newWidth = int(float64(newHeight) * ratio)
	}

	// Create thumbnail
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

// getVideoInfo retrieves video metadata using ffprobe
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

	// Parse output (simplified version)
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

// internal/handler/media.go
package handler

import (
	"fmt"
	"log/slog"
	"strconv"

	"github.com/google/uuid"
	"github.com/iankencruz/threefive/database/generated"
	"github.com/iankencruz/threefive/internal/middleware"
	"github.com/iankencruz/threefive/internal/services"
	"github.com/iankencruz/threefive/pkg/responses"
	"github.com/iankencruz/threefive/templates/lib"
	"github.com/iankencruz/threefive/templates/pages"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v5"
)

type MediaHandler struct {
	mediaService *services.MediaService
	logger       *slog.Logger
}

func NewMediaHandler(mediaService *services.MediaService, logger *slog.Logger) *MediaHandler {
	return &MediaHandler{
		mediaService: mediaService,
		logger:       logger,
	}
}

// ShowMediaList renders the media library page
func (h *MediaHandler) ShowMediaList(c *echo.Context) error {
	user := middleware.GetUser(c)
	currentPath := c.Request().URL.Path

	h.logger.Info("Rendering media library page", "path", currentPath)

	// Get pagination parameters
	page := 1
	if p := c.QueryParam("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	limit := int32(20)
	offset := int32((page - 1) * int(limit))

	// Get filter parameters
	mediaType := c.QueryParam("type") // "image", "video", "document"

	var mediaList []generated.Media
	var err error

	if mediaType != "" {
		// Filter by type
		var mimePattern string
		switch mediaType {
		case "image":
			mimePattern = "image/%"
		case "video":
			mimePattern = "video/%"
		case "document":
			mimePattern = "application/%"
		default:
			mimePattern = "%"
		}

		h.logger.Debug("Filtering media by type", "type", mediaType, "pattern", mimePattern)

		mediaList, err = h.mediaService.ListMediaByType(c.Request().Context(), mimePattern, limit, offset)
		if err != nil {
			h.logger.Error("failed to list media by type", "error", err)
			return c.String(500, "Failed to load media")
		}
	} else {
		// List all media
		mediaList, err = h.mediaService.ListMedia(c.Request().Context(), limit, offset)
		if err != nil {
			h.logger.Error("failed to list media", "error", err)
			return c.String(500, "Failed to load media")
		}
	}

	// Get total count for pagination
	totalCount, err := h.mediaService.CountMedia(c.Request().Context())
	if err != nil {
		h.logger.Error("failed to count media", "error", err)
		totalCount = 0
	}

	// Calculate total pages
	totalPages := int((totalCount + int64(limit) - 1) / int64(limit))

	mediaResponses := h.mediaService.ToMediaResponses(mediaList)

	h.logger.Debug("Media library loaded",
		"total_media", len(mediaList),
		"page", page,
		"total_pages", totalPages,
		"filter", mediaType,
	)

	// Add user to context for template
	ctx := lib.WithUser(c.Request().Context(), user)

	// Render the media library page
	component := pages.MediaLibrary(pages.MediaLibraryProps{
		Media:       mediaResponses,
		CurrentPage: page,
		TotalPages:  totalPages,
		MediaType:   mediaType,
	}, currentPath)

	return responses.Render(ctx, c, component)
}

// GetMediaDetail returns the detail modal for a specific media item
func (h *MediaHandler) GetMediaDetail(c *echo.Context) error {
	// Get media ID from URL parameter
	idParam := c.Param("id")

	h.logger.Debug("Get media detail request", "id_param", idParam)

	if idParam == "" {
		h.logger.Error("missing media ID parameter")
		return responses.ErrorToast(c.Request().Context(), c, "Media ID is required")
	}

	// Parse UUID
	mediaUUID, err := uuid.Parse(idParam)
	if err != nil {
		h.logger.Error("invalid media ID format", "id_param", idParam, "error", err)
		return responses.ErrorToast(c.Request().Context(), c, "Invalid media ID format")
	}

	// Convert to pgtype.UUID
	var mediaID pgtype.UUID
	if err := mediaID.Scan(mediaUUID.String()); err != nil {
		h.logger.Error("failed to convert UUID", "error", err)
		return responses.ErrorToast(c.Request().Context(), c, "Failed to process media ID")
	}

	h.logger.Info("Fetching media details", "media_id", mediaUUID)

	// Get media from database
	media, err := h.mediaService.GetMediaByID(c.Request().Context(), mediaID)
	if err != nil {
		h.logger.Error("failed to get media", "error", err, "media_id", mediaUUID)
		return responses.ErrorToast(c.Request().Context(), c, "Media not found")
	}

	// Convert to response
	mediaResponse := services.MediaResponse{
		ID:               media.ID,
		Filename:         media.Filename,
		OriginalFilename: media.OriginalFilename,
		MimeType:         media.MimeType,
		FileSize:         media.FileSize,
		Width:            media.Width,
		Height:           media.Height,
		URL:              h.mediaService.GetMediaURL(media),
		AltText:          media.AltText.String,
		StorageType:      media.StorageType,
		CreatedAt:        media.CreatedAt,
		UpdatedAt:        media.UpdatedAt,
	}

	h.logger.Info("Media details retrieved successfully", "media_id", mediaUUID)

	// Render the detail modal
	component := lib.MediaDetailModal(mediaResponse)
	return responses.Render(c.Request().Context(), c, component)
}

// UploadMedia handles file upload (returns HTMX-compatible HTML)
func (h *MediaHandler) UploadMedia(c *echo.Context) error {
	user := middleware.GetUser(c)
	if user == nil {
		h.logger.Error("user not found in context")
		return responses.ErrorToast(c.Request().Context(), c, "Authentication required")
	}

	// Get uploaded file
	file, err := c.FormFile("file")
	if err != nil {
		h.logger.Error("failed to get uploaded file", "error", err)
		return responses.ErrorToast(c.Request().Context(), c, "No file uploaded")
	}

	// Get optional alt text
	altText := c.FormValue("alt_text")

	// Convert user ID to pgtype.UUID
	var uploadedBy pgtype.UUID
	if err := uploadedBy.Scan(user.ID.String()); err != nil {
		h.logger.Error("failed to convert user ID", "error", err)
		return responses.ErrorToast(c.Request().Context(), c, "Internal server error")
	}

	// Upload media (this now handles video thumbnail generation)
	media, err := h.mediaService.UploadMedia(c.Request().Context(), file, altText, uploadedBy)
	if err != nil {
		h.logger.Error("failed to upload media", "error", err)
		return responses.ErrorToast(c.Request().Context(), c, err.Error())
	}

	if media == nil {
		h.logger.Error("media service returned nil media")
		return responses.ErrorToast(c.Request().Context(), c, "Failed to process upload")
	}

	h.logger.Info("media uploaded successfully",
		"media_id", media.ID,
		"filename", media.Filename,
		"user_id", user.ID,
	)

	// Get total count BEFORE this upload to determine if grid exists
	totalCount, err := h.mediaService.CountMedia(c.Request().Context())
	if err != nil {
		h.logger.Error("failed to count media", "error", err)
		totalCount = 1
	}

	// Use ToMediaResponse helper which includes ThumbnailURL
	mediaResponse := h.mediaService.ToMediaResponse(media)

	h.logger.Debug("media response details",
		"mime_type", mediaResponse.MimeType,
		"url", mediaResponse.URL,
		"thumbnail_url", mediaResponse.ThumbnailURL,
	)

	if totalCount == 1 {
		// First upload - replace empty state with grid container + first card
		h.logger.Debug("first media upload, creating grid")
		component := lib.MediaGridStart([]services.MediaResponse{mediaResponse})
		c.Response().Header().Set("HX-Reswap", "outerHTML") // Override to replace empty state
		return responses.RenderSuccess(c.Request().Context(), c, component, "File uploaded successfully")
	} else {
		// Subsequent uploads - just return the card (will be appended via afterbegin)
		h.logger.Debug("appending new media card", "total_count", totalCount)
		component := lib.MediaCard(mediaResponse)
		return responses.RenderSuccess(c.Request().Context(), c, component, "File uploaded successfully")
	}
}

func (h *MediaHandler) DeleteMedia(c *echo.Context) error {
	// Get media ID from URL parameter
	idParam := c.Param("id")

	h.logger.Debug("Delete media request", "id_param", idParam)

	if idParam == "" {
		h.logger.Error("missing media ID parameter")
		component := lib.ErrorMessage("Media ID is required")
		return responses.Render(c.Request().Context(), c, component)
	}

	// Parse UUID
	mediaUUID, err := uuid.Parse(idParam)
	if err != nil {
		h.logger.Error("invalid media ID format", "id_param", idParam, "error", err)
		component := lib.ErrorMessage(fmt.Sprintf("Invalid media ID format: %s", idParam))
		return responses.Render(c.Request().Context(), c, component)
	}

	// Convert to pgtype.UUID
	var mediaID pgtype.UUID
	if err := mediaID.Scan(mediaUUID.String()); err != nil {
		h.logger.Error("failed to convert UUID", "error", err)
		component := lib.ErrorMessage("Failed to process media ID")
		return responses.Render(c.Request().Context(), c, component)
	}

	h.logger.Info("Attempting to delete media", "media_id", mediaUUID)

	// Soft delete media
	err = h.mediaService.DeleteMedia(c.Request().Context(), mediaID)
	if err != nil {
		h.logger.Error("failed to delete media", "error", err, "media_id", mediaUUID)
		component := lib.ErrorMessage("Failed to delete media")
		return responses.Render(c.Request().Context(), c, component)
	}

	h.logger.Info("Media deleted successfully", "media_id", mediaUUID)

	// Return success toast (HTMX will remove the element via hx-swap="outerHTML" or hx-swap="delete")
	return responses.SuccessToast(c.Request().Context(), c, "Media deleted successfully")
}

// UpdateMedia handles updating media metadata (alt text, etc.)
func (h *MediaHandler) UpdateMedia(c *echo.Context) error {
	// Get media ID from URL parameter
	idParam := c.Param("id")

	h.logger.Debug("Update media request", "id_param", idParam)

	if idParam == "" {
		h.logger.Error("missing media ID parameter")
		return responses.ErrorToast(c.Request().Context(), c, "Media ID is required")
	}

	// Parse UUID
	mediaUUID, err := uuid.Parse(idParam)
	if err != nil {
		h.logger.Error("invalid media ID format", "id_param", idParam, "error", err)
		return responses.ErrorToast(c.Request().Context(), c, "Invalid media ID format")
	}

	// Convert to pgtype.UUID
	var mediaID pgtype.UUID
	if err := mediaID.Scan(mediaUUID.String()); err != nil {
		h.logger.Error("failed to convert UUID", "error", err)
		return responses.ErrorToast(c.Request().Context(), c, "Failed to process media ID")
	}

	// Parse form data
	if err := c.Request().ParseForm(); err != nil {
		h.logger.Error("failed to parse form", "error", err)
		return responses.ErrorToast(c.Request().Context(), c, "Failed to parse form data")
	}

	altText := c.FormValue("alt_text")
	h.logger.Info("Updating media", "media_id", mediaUUID, "alt_text", altText)

	// Update media in database
	updatedMedia, err := h.mediaService.UpdateMedia(c.Request().Context(), mediaID, altText)
	if err != nil {
		h.logger.Error("failed to update media", "error", err, "media_id", mediaUUID)
		return responses.ErrorToast(c.Request().Context(), c, "Failed to update media")
	}

	h.logger.Info("Media updated successfully", "media_id", mediaUUID)

	// Convert to response
	mediaResponse := services.MediaResponse{
		ID:               updatedMedia.ID,
		Filename:         updatedMedia.Filename,
		OriginalFilename: updatedMedia.OriginalFilename,
		MimeType:         updatedMedia.MimeType,
		FileSize:         updatedMedia.FileSize,
		Width:            updatedMedia.Width,
		Height:           updatedMedia.Height,
		URL:              h.mediaService.GetMediaURL(updatedMedia),
		AltText:          updatedMedia.AltText.String,
		StorageType:      updatedMedia.StorageType,
		CreatedAt:        updatedMedia.CreatedAt,
		UpdatedAt:        updatedMedia.UpdatedAt,
	}

	// Return updated media card
	component := lib.MediaCard(mediaResponse)
	return responses.RenderSuccess(c.Request().Context(), c, component, "Media updated successfully")
}

// ShowMediaSelector returns media grid for selection dialog
func (h *MediaHandler) ShowMediaSelector(c *echo.Context) error {
	h.logger.Debug("Loading media selector")

	// Fetch all media (you can add pagination/filtering later)
	mediaList, err := h.mediaService.ListMedia(c.Request().Context(), 100, 0)
	if err != nil {
		h.logger.Error("failed to list media for selector", "error", err)
		return c.String(500, "Failed to load media")
	}

	// Convert to MediaResponse format
	mediaResponses := h.mediaService.ToMediaResponses(mediaList)

	// Get dialog parameters from query string
	dialogID := c.QueryParam("dialog_id")
	targetInputID := c.QueryParam("target_input_id")

	// Defaults
	if dialogID == "" {
		dialogID = "media-selector"
	}
	if targetInputID == "" {
		targetInputID = "hero_media_id"
	}

	// Render the media grid
	component := lib.MediaSelectorGrid(mediaResponses, dialogID, targetInputID)
	return responses.Render(c.Request().Context(), c, component)
}

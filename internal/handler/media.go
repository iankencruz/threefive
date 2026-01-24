// internal/handler/media.go
package handler

import (
	"log/slog"
	"strconv"

	"github.com/google/uuid"
	"github.com/iankencruz/threefive/components/toast"
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

	// Convert generated.Media to services.MediaResponse
	mediaResponses := make([]services.MediaResponse, len(mediaList))
	for i, m := range mediaList {
		mediaResponses[i] = services.MediaResponse{
			ID:               m.ID,
			Filename:         m.Filename,
			OriginalFilename: m.OriginalFilename,
			MimeType:         m.MimeType,
			FileSize:         m.FileSize,
			Width:            m.Width,
			Height:           m.Height,
			URL:              h.mediaService.GetMediaURL(&m),
			AltText:          m.AltText.String,
			CreatedAt:        m.CreatedAt,
			UpdatedAt:        m.UpdatedAt,
		}
	}

	h.logger.Info("Media library loaded",
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

// UploadMedia handles file upload (returns HTMX-compatible HTML)
func (h *MediaHandler) UploadMedia(c *echo.Context) error {
	user := middleware.GetUser(c)

	// Get uploaded file
	file, err := c.FormFile("file")
	if err != nil {
		h.logger.Error("failed to get uploaded file", "error", err)
		// Return error toast
		component := toast.Toast(toast.Props{
			Title:       "File Upload Failed",
			Description: "No file was uploaded.",
		})
		return responses.Render(c.Request().Context(), c, component)
	}

	// Get optional alt text
	altText := c.FormValue("alt_text")

	// Convert user ID to pgtype.UUID
	var uploadedBy pgtype.UUID
	if err := uploadedBy.Scan(user.ID.String()); err != nil {
		h.logger.Error("failed to convert user ID", "error", err)
		component := toast.Toast(toast.Props{
			Title:       "Internal server error",
			Description: "Could not process user information.",
		})
		return responses.Render(c.Request().Context(), c, component)
	}

	// Upload media
	media, err := h.mediaService.UploadMedia(c.Request().Context(), file, altText, uploadedBy)
	if err != nil {
		h.logger.Error("failed to upload media", "error", err)
		component := toast.Toast(toast.Props{
			Title:       "File Upload Failed",
			Description: err.Error(),
		})
		return responses.Render(c.Request().Context(), c, component)
	}

	h.logger.Info("media uploaded successfully",
		"media_id", media.ID,
		"filename", media.Filename,
		"user_id", user.ID,
	)

	mediaResponse := h.mediaService.ToMediaResponse(media)

	// Return success - render the new media card component
	component := lib.MediaCard(mediaResponse)
	return responses.RenderSuccess(c.Request().Context(), c, component, "File uploaded successfully")
}

// DeleteMedia handles media deletion
func (h *MediaHandler) DeleteMedia(c *echo.Context) error {
	// Get media ID from URL parameter
	idParam := c.Param("id")
	mediaUUID, err := uuid.Parse(idParam)
	if err != nil {
		component := toast.Toast(toast.Props{
			Title:       "Invalid Media ID",
			Description: "The provided media ID is not valid.",
		})
		return responses.Render(c.Request().Context(), c, component)
	}

	var mediaID pgtype.UUID
	if err := mediaID.Scan(mediaUUID.String()); err != nil {
		component := toast.Toast(toast.Props{
			Title:       "Internal server error",
			Description: "Could not process media ID.",
		})
		return responses.Render(c.Request().Context(), c, component)
	}

	// Soft delete media
	err = h.mediaService.DeleteMedia(c.Request().Context(), mediaID)
	if err != nil {
		h.logger.Error("failed to delete media", "error", err, "media_id", mediaUUID)
		component := toast.Toast(toast.Props{
			Title:       "Failed to Delete Media",
			Description: err.Error(),
		})
		return responses.Render(c.Request().Context(), c, component)
	}

	h.logger.Info("media deleted", "media_id", mediaUUID)

	// Return success toast (HTMX will remove the element via hx-swap="outerHTML")
	return responses.SuccessToast(c.Request().Context(), c, "Media deleted successfully")
}

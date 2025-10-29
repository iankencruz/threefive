// backend/internal/media/handler.go
package media

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/iankencruz/threefive/internal/auth"
	"github.com/iankencruz/threefive/internal/shared/sqlc"
	"github.com/iankencruz/threefive/internal/shared/storage"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Handler struct {
	service *Service
}

// NewHandler creates a new media handler with its own service
func NewHandler(db *pgxpool.Pool, queries *sqlc.Queries, storage storage.Storage) *Handler {
	service := NewService(db, queries, storage)
	return &Handler{
		service: service,
	}
}

// UploadHandler handles file uploads
// POST /api/v1/media/upload
func (h *Handler) UploadHandler(w http.ResponseWriter, r *http.Request) {
	// ✅ Get user from context using auth helper
	user := auth.MustGetUserFromContext(r.Context())

	// Parse multipart form (max 100MB)
	if err := r.ParseMultipartForm(100 << 20); err != nil {
		respondError(w, http.StatusBadRequest, "failed to parse form: "+err.Error())
		return
	}

	// Get the file from form
	file, header, err := r.FormFile("file")
	if err != nil {
		respondError(w, http.StatusBadRequest, "no file provided")
		return
	}
	defer file.Close()

	// Upload file
	media, err := h.service.UploadFile(r.Context(), header, user.ID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, media)
}

// GetMediaHandler retrieves a single media by ID
// GET /api/v1/media/{id}
func (h *Handler) GetMediaHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid media ID")
		return
	}

	media, err := h.service.GetMediaByID(r.Context(), id)
	if err != nil {
		respondError(w, http.StatusNotFound, "media not found")
		return
	}

	respondJSON(w, http.StatusOK, media)
}

// ListMediaHandler lists media with pagination
// GET /api/v1/media?page=1&limit=20
func (h *Handler) ListMediaHandler(w http.ResponseWriter, r *http.Request) {
	// Parse pagination params
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit

	media, err := h.service.ListMedia(r.Context(), int32(limit), int32(offset))
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := map[string]any{
		"data": media,
		"pagination": map[string]int{
			"page":  page,
			"limit": limit,
		},
	}

	respondJSON(w, http.StatusOK, response)
}

// DeleteMediaHandler soft deletes media
// DELETE /api/v1/media/{id}
func (h *Handler) DeleteMediaHandler(w http.ResponseWriter, r *http.Request) {
	// ✅ Get user from context using auth helper
	user := auth.MustGetUserFromContext(r.Context())

	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid media ID")
		return
	}

	// Check for hard delete query param
	hardDelete := r.URL.Query().Get("hard") == "true"

	var deleteErr error
	if hardDelete {
		deleteErr = h.service.HardDeleteMedia(r.Context(), id)
	} else {
		deleteErr = h.service.DeleteMedia(r.Context(), id, user.ID)
	}

	if deleteErr != nil {
		respondError(w, http.StatusInternalServerError, deleteErr.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{
		"message": "media deleted successfully",
	})
}

// LinkMediaHandler links media to an entity
// POST /api/v1/media/{id}/link
// Body: {"entity_type": "project", "entity_id": "uuid", "sort_order": 1}
func (h *Handler) LinkMediaHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	mediaID, err := uuid.Parse(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid media ID")
		return
	}

	var req struct {
		EntityType string    `json:"entity_type"`
		EntityID   uuid.UUID `json:"entity_id"`
		SortOrder  int       `json:"sort_order"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Validate entity type
	validTypes := map[string]bool{"project": true, "page": true, "gallery": true}
	if !validTypes[req.EntityType] {
		respondError(w, http.StatusBadRequest, "invalid entity_type (must be: project, page, or gallery)")
		return
	}

	if err := h.service.LinkMediaToEntity(r.Context(), mediaID, req.EntityType, req.EntityID, int32(req.SortOrder)); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{
		"message": "media linked successfully",
	})
}

// UnlinkMediaHandler unlinks media from an entity
// DELETE /api/v1/media/{id}/link
// Body: {"entity_type": "project", "entity_id": "uuid"}
func (h *Handler) UnlinkMediaHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	mediaID, err := uuid.Parse(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid media ID")
		return
	}

	var req struct {
		EntityType string    `json:"entity_type"`
		EntityID   uuid.UUID `json:"entity_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.service.UnlinkMediaFromEntity(r.Context(), mediaID, req.EntityType, req.EntityID); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{
		"message": "media unlinked successfully",
	})
}

// GetEntityMediaHandler gets all media for an entity
// GET /api/v1/media/entity/{type}/{id}
// Example: GET /api/v1/media/entity/project/uuid-here
func (h *Handler) GetEntityMediaHandler(w http.ResponseWriter, r *http.Request) {
	entityType := chi.URLParam(r, "type")
	entityIDStr := chi.URLParam(r, "id")

	entityID, err := uuid.Parse(entityIDStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid entity ID")
		return
	}

	// Validate entity type
	validTypes := map[string]bool{"project": true, "page": true, "gallery": true}
	if !validTypes[entityType] {
		respondError(w, http.StatusBadRequest, "invalid entity type (must be: project, page, or gallery)")
		return
	}

	media, err := h.service.GetMediaForEntity(r.Context(), entityType, entityID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]any{
		"entity_type": entityType,
		"entity_id":   entityID,
		"media":       media,
	})
}

// GetStatsHandler returns media storage statistics
// GET /api/v1/media/stats
func (h *Handler) GetStatsHandler(w http.ResponseWriter, r *http.Request) {
	stats, err := h.service.GetMediaStats(r.Context())
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, stats)
}

// Helper functions for JSON responses
func respondJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// helper function for JSON error responses
func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{
		"error": message,
	})
}

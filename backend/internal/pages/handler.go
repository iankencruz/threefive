// backend/internal/pages/handler.go
package pages

import (
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/iankencruz/threefive/internal/auth"
	"github.com/iankencruz/threefive/internal/blocks"
	"github.com/iankencruz/threefive/internal/config"
	"github.com/iankencruz/threefive/internal/shared/responses"
	"github.com/iankencruz/threefive/internal/shared/sqlc"
	"github.com/iankencruz/threefive/internal/shared/validation"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Handler struct {
	service *Service
}

// NewHandler creates a new pages handler with its own service
func NewHandler(db *pgxpool.Pool, queries *sqlc.Queries, cfg *config.Config) *Handler {
	// Create block service internally (only needs queries)
	blockService := blocks.NewService(queries)

	// Create service config
	serviceCfg := ServiceConfig{
		AutoPurgeRetentionDays: cfg.AutoPurgeRetentionDays,
	}

	// Create pages service
	service := NewService(db, queries, blockService, serviceCfg)

	return &Handler{
		service: service,
	}
}

// CreatePage handles page creation
// POST /api/v1/pages
func (h *Handler) CreatePage(w http.ResponseWriter, r *http.Request) {
	var req CreatePageRequest

	// Get current user from context
	user := auth.MustGetUserFromContext(r.Context())

	// Parse and validate request
	err := validation.ParseAndValidateJSON(r, &req, func(v *validation.Validator) {
		req.Validate(v)
	})
	if err != nil {
		responses.WriteErr(w, err)
		return
	}

	// Create page
	page, err := h.service.CreatePage(r.Context(), req, user.ID)
	if err != nil {
		responses.WriteErr(w, err)
		return
	}

	responses.WriteCreated(w, page)
}

// GetPage handles retrieving a single page by ID or slug
// GET /api/v1/pages/{idOrSlug}
func (h *Handler) GetPage(w http.ResponseWriter, r *http.Request) {
	idOrSlug := chi.URLParam(r, "idOrSlug")

	// Try to parse as UUID first
	id, err := uuid.Parse(idOrSlug)
	var page *PageResponse

	if err == nil {
		// It's a valid UUID, get by ID
		page, err = h.service.GetPageByID(r.Context(), id)
	} else {
		// It's a slug, get by slug
		page, err = h.service.GetPageBySlug(r.Context(), idOrSlug)
	}

	if err != nil {
		responses.WriteErr(w, err)
		return
	}

	responses.WriteOK(w, page)
}

// ListPages handles listing pages with pagination
// GET /api/v1/pages?page=1&limit=20
func (h *Handler) ListPages(w http.ResponseWriter, r *http.Request) {
	// Parse pagination params
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := int32((page - 1) * limit)

	// List pages
	result, err := h.service.ListPages(r.Context(), int32(limit), offset)
	if err != nil {
		responses.WriteErr(w, err)
		return
	}

	responses.WriteOK(w, result)
}

// UpdatePage handles updating a page
// PUT /api/v1/pages/{id}
func (h *Handler) UpdatePage(w http.ResponseWriter, r *http.Request) {
	var req UpdatePageRequest

	// Parse page ID
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		responses.WriteErr(w, err)
		return
	}

	// Parse and validate request
	err = validation.ParseAndValidateJSON(r, &req, func(v *validation.Validator) {
		req.Validate(v)
	})
	if err != nil {
		responses.WriteErr(w, err)
		return
	}

	// Update page
	page, err := h.service.UpdatePage(r.Context(), id, req)
	if err != nil {
		responses.WriteErr(w, err)
		return
	}

	responses.WriteOK(w, page)
}

// UpdatePageStatus handles updating page status
// PATCH /api/v1/pages/{id}/status
func (h *Handler) UpdatePageStatus(w http.ResponseWriter, r *http.Request) {
	var req UpdatePageStatusRequest

	// Parse page ID
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		responses.WriteErr(w, err)
		return
	}

	// Parse and validate request
	err = validation.ParseAndValidateJSON(r, &req, func(v *validation.Validator) {
		req.Validate(v)
	})
	if err != nil {
		responses.WriteErr(w, err)
		return
	}

	// Update status
	err = h.service.UpdatePageStatus(r.Context(), id, req.Status)
	if err != nil {
		responses.WriteErr(w, err)
		return
	}

	response := map[string]string{
		"message": "Page status updated successfully",
	}

	responses.WriteOK(w, response)
}

// DeletePage handles soft deleting a page
// DELETE /api/v1/pages/{id}
func (h *Handler) DeletePage(w http.ResponseWriter, r *http.Request) {
	// Parse page ID
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		responses.WriteErr(w, err)
		return
	}

	// Delete page
	err = h.service.DeletePage(r.Context(), id)
	if err != nil {
		responses.WriteErr(w, err)
		return
	}

	response := map[string]string{
		"message": "Page deleted successfully",
	}

	responses.WriteOK(w, response)
}

// PurgeDeletedPages manually triggers purge of old deleted pages
// POST /api/v1/pages/purge
func (h *Handler) PurgeDeletedPages(w http.ResponseWriter, r *http.Request) {
	log.Printf("[Pages] Manual purge triggered by user")

	// Call service to purge
	err := h.service.PurgeOldDeletedPages(r.Context())
	if err != nil {
		log.Printf("[Pages] Purge failed: %v", err)
		responses.WriteErr(w, err)
		return
	}

	log.Printf("[Pages] Manual purge completed successfully")

	response := map[string]string{
		"message": "Old deleted pages purged successfully",
	}

	responses.WriteOK(w, response)
}

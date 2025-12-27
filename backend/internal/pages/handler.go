// backend/internal/pages/handler.go
package pages

import (
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
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

	// Parse and validate request
	req, err := validation.ParseAndValidate[*CreatePageRequest](r)
	if err != nil {
		responses.WriteErr(w, err)
		return
	}

	// Create page
	page, err := h.service.CreatePage(r.Context(), req)
	if err != nil {
		responses.WriteErr(w, err)
		return
	}

	responses.WriteCreated(w, page)
}

// GetPageByID fetches a page by UUID (for admin editing)
func (h *Handler) GetPageByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	log.Printf("DEBUG GetPageByID: Captured ID parameter: '%s'", idStr)

	id, err := uuid.Parse(idStr)
	if err != nil {
		log.Printf("DEBUG GetPageByID: UUID parse failed for: '%s'", idStr)
		responses.WriteErr(w, err)
		return
	}

	page, err := h.service.GetPageByID(r.Context(), id)
	if err != nil {
		responses.WriteErr(w, err)
		return
	}

	responses.WriteOK(w, page)
}

// GetPage handles retrieving a single page by ID or slug
// GET /api/v1/pages/{slug}
func (h *Handler) GetPageBySlug(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	// log.Printf("DEBUG GetPageBySlug: Called with slug: '%s'", slug)
	// Try to parse as UUID first
	page, err := h.service.GetPageBySlug(r.Context(), slug)
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

	// Parse filter parameters
	var statusFilter *string
	if status := r.URL.Query().Get("status"); status != "" {
		// Validate status is a valid enum value
		if status == "draft" || status == "published" || status == "archived" {
			statusFilter = &status
		}
	}

	// Parse sort parameters
	sortBy := r.URL.Query().Get("sort")
	validSortFields := map[string]bool{
		"created_at":   true,
		"published_at": true,
		"title":        true,
	}
	if sortBy == "" || !validSortFields[sortBy] {
		sortBy = "created_at"
	}

	sortOrder := r.URL.Query().Get("order")
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "desc"
	}

	offset := int32((page - 1) * limit)

	// List pages with filters
	result, err := h.service.ListPages(r.Context(), ListPagesParams{
		StatusFilter: statusFilter,
		SortBy:       sortBy,
		SortOrder:    sortOrder,
		Limit:        int32(limit),
		Offset:       offset,
	})
	if err != nil {
		responses.WriteErr(w, err)
		return
	}

	// Build response with pagination and filters
	response := map[string]any{
		"data":       result.Pages,
		"pagination": result.Pagination,
		"filters": map[string]any{
			"status": statusFilter,
			"sort":   sortBy,
			"order":  sortOrder,
		},
	}

	responses.WriteJSON(w, http.StatusOK, response)
}

// ListPublishedPages handles listing only published pages (public route)
func (h *Handler) ListPublishedPages(w http.ResponseWriter, r *http.Request) {
	// Parse pagination params (no status filter needed)
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 || limit > 100 {
		limit = 20
	}

	// Parse sort params
	sortBy := r.URL.Query().Get("sort")
	if sortBy == "" {
		sortBy = "published_at"
	}

	sortOrder := r.URL.Query().Get("order")
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "desc"
	}

	offset := (page - 1) * limit

	result, err := h.service.ListPublishedPages(r.Context(), ListPagesParams{
		SortBy:    sortBy,
		SortOrder: sortOrder,
		Limit:     int32(limit),
		Offset:    int32(offset),
	})
	if err != nil {
		responses.WriteErr(w, err)
		return
	}

	data := map[string]any{
		"data":       result.Pages,
		"pagination": result.Pagination,
	}

	responses.WriteJSON(w, http.StatusOK, data)
}

// UpdatePage handles updating a page
// PUT /api/v1/pages/{id}
func (h *Handler) UpdatePage(w http.ResponseWriter, r *http.Request) {

	// Parse page ID
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		responses.WriteErr(w, err)
		return
	}

	// Parse and validate request
	// err = validation.ParseAndValidate(r, &req, func(v *validation.Validator) {
	// 	req.Validate(v)
	// })

	req, err := validation.ParseAndValidate[*UpdatePageRequest](r)
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

	// Parse page ID
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		responses.WriteErr(w, err)
		return
	}

	// Parse and validate request
	req, err := validation.ParseAndValidate[*UpdatePageStatusRequest](r)
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

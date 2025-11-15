package blogs

import (
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
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

// NewHandler creates a new blogs handler with its own service
func NewHandler(db *pgxpool.Pool, queries *sqlc.Queries, cfg *config.Config) *Handler {
	// Create block service internally (only needs queries)
	blockService := blocks.NewService(queries)

	// Create service config
	serviceCfg := ServiceConfig{
		AutoPurgeRetentionDays: cfg.AutoPurgeRetentionDays,
	}

	// Create blogs service
	service := NewService(db, queries, blockService, serviceCfg)

	return &Handler{
		service: service,
	}
}

// CreateBlog handles blog creation
// POST /api/v1/blogs
func (h *Handler) CreateBlog(w http.ResponseWriter, r *http.Request) {
	var req CreateBlogRequest

	// Parse and validate request
	err := validation.ParseAndValidateJSON(r, &req, func(v *validation.Validator) {
		req.Validate(v)
	})
	if err != nil {
		responses.WriteErr(w, err)
		return
	}

	// Create blog
	blog, err := h.service.CreateBlog(r.Context(), req)
	if err != nil {
		responses.WriteErr(w, err)
		return
	}

	responses.WriteCreated(w, blog)
}

// GetBlog handles retrieving a single blog by ID or slug
// GET /api/v1/blogs/{idOrSlug}
func (h *Handler) GetBlog(w http.ResponseWriter, r *http.Request) {
	idOrSlug := chi.URLParam(r, "idOrSlug")

	// Try parsing as UUID first
	if id, err := uuid.Parse(idOrSlug); err == nil {
		blog, err := h.service.GetBlogByID(r.Context(), id)
		if err != nil {
			responses.WriteErr(w, err)
			return
		}
		responses.WriteJSON(w, http.StatusOK, blog)
		return
	}

	// Otherwise, treat as slug
	blog, err := h.service.GetBlogBySlug(r.Context(), idOrSlug)
	if err != nil {
		responses.WriteErr(w, err)
		return
	}

	responses.WriteJSON(w, http.StatusOK, blog)
}

// ListBlogs handles retrieving a paginated list of blogs
// GET /api/v1/blogs
func (h *Handler) ListBlogs(w http.ResponseWriter, r *http.Request) {
	// Parse pagination parameters
	limit := int32(10)
	offset := int32(0)

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.ParseInt(limitStr, 10, 32); err == nil && l > 0 {
			limit = int32(l)
		}
	}

	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if o, err := strconv.ParseInt(offsetStr, 10, 32); err == nil && o >= 0 {
			offset = int32(o)
		}
	}

	// Limit max to 100
	if limit > 100 {
		limit = 100
	}

	blogs, err := h.service.ListBlogs(r.Context(), limit, offset)
	if err != nil {
		responses.WriteErr(w, err)
		return
	}

	responses.WriteJSON(w, http.StatusOK, blogs)
}

// UpdateBlog handles updating a blog
// PUT /api/v1/blogs/{id}
func (h *Handler) UpdateBlog(w http.ResponseWriter, r *http.Request) {
	// Get blog ID from URL
	idStr := chi.URLParam(r, "id")
	blogID, err := uuid.Parse(idStr)
	if err != nil {
		responses.WriteErr(w, err)
		return
	}

	var req UpdateBlogRequest

	// Parse and validate request
	err = validation.ParseAndValidateJSON(r, &req, func(v *validation.Validator) {
		req.Validate(v)
	})
	if err != nil {
		responses.WriteErr(w, err)
		return
	}

	// Update blog
	blog, err := h.service.UpdateBlog(r.Context(), blogID, req)
	if err != nil {
		responses.WriteErr(w, err)
		return
	}

	responses.WriteJSON(w, http.StatusOK, blog)
}

// UpdateBlogStatus handles updating only the blog status
// PATCH /api/v1/blogs/{id}/status
func (h *Handler) UpdateBlogStatus(w http.ResponseWriter, r *http.Request) {
	// Get blog ID from URL
	idStr := chi.URLParam(r, "id")
	blogID, err := uuid.Parse(idStr)
	if err != nil {
		responses.WriteErr(w, err)
		return
	}

	var req UpdateBlogStatusRequest

	// Parse and validate request
	err = validation.ParseAndValidateJSON(r, &req, func(v *validation.Validator) {
		req.Validate(v)
	})
	if err != nil {
		responses.WriteErr(w, err)
		return
	}

	// Update status
	blog, err := h.service.UpdateBlogStatus(r.Context(), blogID, req.Status)
	if err != nil {
		responses.WriteErr(w, err)
		return
	}

	responses.WriteJSON(w, http.StatusOK, blog)
}

// DeleteBlog handles soft-deleting a blog
// DELETE /api/v1/blogs/{id}
func (h *Handler) DeleteBlog(w http.ResponseWriter, r *http.Request) {
	// Get blog ID from URL
	idStr := chi.URLParam(r, "id")
	blogID, err := uuid.Parse(idStr)
	if err != nil {
		responses.WriteErr(w, err)
		return
	}

	// Soft delete blog
	err = h.service.SoftDeleteBlog(r.Context(), blogID)
	if err != nil {
		responses.WriteErr(w, err)
		return
	}

	responses.WriteNoContent(w)
}

// PurgeOldDeletedBlogs handles cleanup of old soft-deleted blogs
// POST /api/v1/blogs/purge
func (h *Handler) PurgeOldDeletedBlogs(w http.ResponseWriter, r *http.Request) {
	err := h.service.PurgeOldDeletedBlogs(r.Context())
	if err != nil {
		log.Printf("Error purging old deleted blogs: %v", err)
		responses.WriteErr(w, err)
		return
	}

	responses.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "Old deleted blogs purged successfully",
	})
}

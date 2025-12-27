package blogs

import (
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/iankencruz/threefive/internal/blocks"
	"github.com/iankencruz/threefive/internal/config"
	"github.com/iankencruz/threefive/internal/shared/errors"
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

	// Parse and validate request
	req, err := validation.ParseAndValidate[*CreateBlogRequest](r)
	if err != nil {
		responses.WriteErr(w, err)
		return
	}

	// Create blog
	blog, err := h.service.CreateBlog(r.Context(), req)
	if err != nil {
		responses.WriteErr(w, errors.BadRequest("Failed to create blog", "blog_creation_failed"))
		return
	}

	responses.WriteCreated(w, blog)
}

func (h *Handler) GetBlogByID(w http.ResponseWriter, r *http.Request) {
	// 👇 ADD THIS DEBUG BLOCK
	if rctx := chi.RouteContext(r.Context()); rctx != nil {
		log.Printf("🕵️ CHI DEBUG: Route Pattern: %s", rctx.RoutePattern())
		log.Printf("🕵️ CHI DEBUG: URL Params Keys: %v", rctx.URLParams.Keys)
		log.Printf("🕵️ CHI DEBUG: URL Params Values: %v", rctx.URLParams.Values)
	} else {
		log.Printf("🕵️ CHI DEBUG: RouteContext is NIL (Context was wiped!)")
	}
	// 👆 END DEBUG BLOCK

	// Try getting ALL possible parameter names
	idStr := chi.URLParam(r, "id")
	slugStr := chi.URLParam(r, "slug")

	log.Printf("🔵 URLParam 'id': '%s'", idStr)
	log.Printf("🔵 URLParam 'slug': '%s'", slugStr)

	// If id is empty but slug has the UUID, that's the problem
	if idStr == "" && slugStr != "" {
		log.Printf("⚠️ FOUND THE BUG: UUID is in 'slug' parameter, not 'id'!")
		idStr = slugStr
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		log.Printf("🔵 UUID parse error: %v", err)
		responses.WriteErr(w, errors.NotFound("Blog ID not found", "blog_not_found"))
		return
	}

	page, err := h.service.GetBlogByID(r.Context(), id)
	if err != nil {
		responses.WriteErr(w, errors.NotFound("Failed to retrieve blog", "blog_retrieval_failed"))
		return
	}

	responses.WriteOK(w, page)
}

// GetBlogBySlug handles retrieving a single blog by slug
// GET /api/v1/blogs/{slug}
func (h *Handler) GetBlogBySlug(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	blog, err := h.service.GetBlogBySlug(r.Context(), slug)
	if err != nil {
		responses.WriteErr(w, err)
		return
	}

	responses.WriteOK(w, blog)
}

// ListBlogs handles retrieving a paginated list of blogs with optional filters
// GET /api/v1/blogs?status=published&featured=true&page=1&limit=20&sort=created_at&order=desc
func (h *Handler) ListBlogs(w http.ResponseWriter, r *http.Request) {
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

	var featuredFilter *bool
	if featured := r.URL.Query().Get("featured"); featured != "" {
		if featured == "true" {
			t := true
			featuredFilter = &t
		} else if featured == "false" {
			f := false
			featuredFilter = &f
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

	offset := (page - 1) * limit

	// List blogs with filters
	result, err := h.service.ListBlogs(r.Context(), ListBlogsParams{
		StatusFilter:   statusFilter,
		FeaturedFilter: featuredFilter,
		SortBy:         sortBy,
		SortOrder:      sortOrder,
		Limit:          int32(limit),
		Offset:         int32(offset),
	})
	if err != nil {
		responses.WriteErr(w, err)
		return
	}

	// Build response with pagination and filters
	response := map[string]any{
		"data":       result.Blogs,
		"pagination": result.Pagination,
		"filters": map[string]any{
			"status":   statusFilter,
			"featured": featuredFilter,
			"sort":     sortBy,
			"order":    sortOrder,
		},
	}

	responses.WriteJSON(w, http.StatusOK, response)
}

// ListPublishedBlogs handles listing only published blogs (public route)
// GET /api/v1/blogs?page=1&limit=20&sort=published_at&order=desc
func (h *Handler) ListPublishedBlogs(w http.ResponseWriter, r *http.Request) {
	// Parse pagination params
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 || limit > 100 {
		limit = 20
	}

	// Parse sort parameters
	sortBy := r.URL.Query().Get("sort")
	validSortFields := map[string]bool{
		"created_at":   true,
		"published_at": true,
		"title":        true,
	}
	if sortBy == "" || !validSortFields[sortBy] {
		sortBy = "published_at" // Default to published_at for public
	}

	sortOrder := r.URL.Query().Get("order")
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "desc"
	}

	offset := (page - 1) * limit

	// List published blogs
	result, err := h.service.ListPublishedBlogs(r.Context(), ListBlogsParams{
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
		"data":       result.Blogs,
		"pagination": result.Pagination,
	}
	// Simple response without filters (no status filter needed for public)
	responses.WriteJSON(w, http.StatusOK, data)
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

	// Parse and validate request
	req, err := validation.ParseAndValidate[*UpdateBlogRequest](r)
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

	// Parse and validate request
	req, err := validation.ParseAndValidate[*UpdateBlogStatusRequest](r)
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

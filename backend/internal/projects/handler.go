// backend/internal/projects/handler.go
package projects

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/iankencruz/threefive/internal/auth"
	"github.com/iankencruz/threefive/internal/shared/responses"
	"github.com/iankencruz/threefive/internal/shared/sqlc"
	"github.com/iankencruz/threefive/internal/shared/validation"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Handler struct {
	service *Service
}

// NewHandler creates a new projects handler with its own service
func NewHandler(db *pgxpool.Pool, queries *sqlc.Queries) *Handler {
	// Create projects service (no longer needs block service)
	service := NewService(db, queries)

	return &Handler{
		service: service,
	}
}

// CreateProject handles project creation
// POST /api/v1/admin/projects
func (h *Handler) CreateProject(w http.ResponseWriter, r *http.Request) {
	// Get current user from context
	user := auth.MustGetUserFromContext(r.Context())

	// Parse and validate request
	req, err := validation.ParseAndValidate[*CreateProjectRequest](r)
	if err != nil {
		responses.WriteErr(w, err)
		return
	}

	// Create project
	project, err := h.service.CreateProject(r.Context(), req, user.ID)
	if err != nil {
		responses.WriteErr(w, err)
		return
	}

	responses.WriteCreated(w, project)
}

// GetProjectByID fetches a project by UUID (for admin editing)
// GET /api/v1/admin/projects/{id}
func (h *Handler) GetProjectByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		responses.WriteErr(w, err)
		return
	}

	project, err := h.service.GetProjectByID(r.Context(), id)
	if err != nil {
		responses.WriteErr(w, err)
		return
	}

	responses.WriteOK(w, project)
}

// GetProjectBySlug handles retrieving a single project by slug (public)
// GET /api/v1/projects/{slug}
func (h *Handler) GetProjectBySlug(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	project, err := h.service.GetProjectBySlug(r.Context(), slug)
	if err != nil {
		responses.WriteErr(w, err)
		return
	}

	responses.WriteOK(w, project)
}

// ListProjects handles listing projects with pagination and optional filters (admin)
// GET /api/v1/admin/projects?status=published&page=1&limit=20&sort=created_at&order=desc
func (h *Handler) ListProjects(w http.ResponseWriter, r *http.Request) {
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
		"project_date": true,
		"project_year": true,
	}
	if sortBy == "" || !validSortFields[sortBy] {
		sortBy = "created_at"
	}

	sortOrder := r.URL.Query().Get("order")
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "desc"
	}

	offset := int32((page - 1) * limit)

	// List projects with filters
	result, err := h.service.ListProjects(r.Context(), ListProjectsParams{
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
		"data":       result.Projects,
		"pagination": result.Pagination,
		"filters": map[string]any{
			"status": statusFilter,
			"sort":   sortBy,
			"order":  sortOrder,
		},
	}

	responses.WriteJSON(w, http.StatusOK, response)
}

// ListPublishedProjects handles listing published projects (public)
// GET /api/v1/projects?page=1&limit=20&sort=created_at&order=desc
func (h *Handler) ListPublishedProjects(w http.ResponseWriter, r *http.Request) {
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
		"project_date": true,
		"project_year": true,
	}
	if sortBy == "" || !validSortFields[sortBy] {
		sortBy = "created_at"
	}

	sortOrder := r.URL.Query().Get("order")
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "desc"
	}

	offset := int32((page - 1) * limit)

	// List published projects only
	publishedStatus := "published"
	result, err := h.service.ListProjects(r.Context(), ListProjectsParams{
		StatusFilter: &publishedStatus,
		SortBy:       sortBy,
		SortOrder:    sortOrder,
		Limit:        int32(limit),
		Offset:       offset,
	})
	if err != nil {
		responses.WriteErr(w, err)
		return
	}

	// Build response with pagination
	data := map[string]any{
		"data":       result.Projects,
		"pagination": result.Pagination,
	}

	responses.WriteJSON(w, http.StatusOK, data)
}

// UpdateProject handles updating a project
// PUT /api/v1/admin/projects/{id}
func (h *Handler) UpdateProject(w http.ResponseWriter, r *http.Request) {
	// Parse project ID
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		responses.WriteErr(w, err)
		return
	}

	// Parse and validate request
	req, err := validation.ParseAndValidate[*UpdateProjectRequest](r)
	if err != nil {
		responses.WriteErr(w, err)
		return
	}

	// Update project
	project, err := h.service.UpdateProject(r.Context(), id, req)
	if err != nil {
		responses.WriteErr(w, err)
		return
	}

	responses.WriteOK(w, project)
}

// UpdateProjectStatus handles updating project status
// PATCH /api/v1/admin/projects/{id}/status
func (h *Handler) UpdateProjectStatus(w http.ResponseWriter, r *http.Request) {
	// Parse project ID
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		responses.WriteErr(w, err)
		return
	}

	// Parse and validate request
	req, err := validation.ParseAndValidate[*UpdateProjectStatusRequest](r)
	if err != nil {
		responses.WriteErr(w, err)
		return
	}

	// Update status
	err = h.service.UpdateProjectStatus(r.Context(), id, req.Status)
	if err != nil {
		responses.WriteErr(w, err)
		return
	}

	response := map[string]string{
		"message": "Project status updated successfully",
	}

	responses.WriteOK(w, response)
}

// DeleteProject handles soft deleting a project
// DELETE /api/v1/admin/projects/{id}
func (h *Handler) DeleteProject(w http.ResponseWriter, r *http.Request) {
	// Parse project ID
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		responses.WriteErr(w, err)
		return
	}

	// Delete project
	err = h.service.DeleteProject(r.Context(), id)
	if err != nil {
		responses.WriteErr(w, err)
		return
	}

	response := map[string]string{
		"message": "Project deleted successfully",
	}

	responses.WriteOK(w, response)
}

// backend/internal/projects/handler.go
package projects

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/iankencruz/threefive/internal/auth"
	"github.com/iankencruz/threefive/internal/blocks"
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
	// Create block service internally (only needs queries)
	blockService := blocks.NewService(queries)

	// Create projects service
	service := NewService(db, queries, blockService)

	return &Handler{
		service: service,
	}
}

// CreateProject handles project creation
// POST /api/v1/projects
func (h *Handler) CreateProject(w http.ResponseWriter, r *http.Request) {
	var req CreateProjectRequest

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

	// Create project
	project, err := h.service.CreateProject(r.Context(), req, user.ID)
	if err != nil {
		responses.WriteErr(w, err)
		return
	}

	responses.WriteCreated(w, project)
}

// GetProject handles retrieving a single project by ID or slug
// GET /api/v1/projects/{idOrSlug}
func (h *Handler) GetProject(w http.ResponseWriter, r *http.Request) {
	idOrSlug := chi.URLParam(r, "idOrSlug")

	// Try to parse as UUID first
	id, err := uuid.Parse(idOrSlug)
	var project *ProjectResponse

	if err == nil {
		// It's a valid UUID, get by ID
		project, err = h.service.GetProjectByID(r.Context(), id)
	} else {
		// It's a slug, get by slug
		project, err = h.service.GetProjectBySlug(r.Context(), idOrSlug)
	}

	if err != nil {
		responses.WriteErr(w, err)
		return
	}

	responses.WriteOK(w, project)
}

// ListProjects handles listing projects with pagination
// GET /api/v1/projects?page=1&limit=20&status=published
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

	// Parse status filter (optional)
	status := r.URL.Query().Get("status")
	var statusPtr *string
	if status != "" {
		statusPtr = &status
	}

	offset := int32((page - 1) * limit)

	// List projects
	result, err := h.service.ListProjects(r.Context(), statusPtr, int32(limit), offset)
	if err != nil {
		responses.WriteErr(w, err)
		return
	}

	responses.WriteOK(w, result)
}

// UpdateProject handles updating a project
// PUT /api/v1/projects/{id}
func (h *Handler) UpdateProject(w http.ResponseWriter, r *http.Request) {
	var req UpdateProjectRequest

	// Parse project ID
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

	// Update project
	project, err := h.service.UpdateProject(r.Context(), id, req)
	if err != nil {
		responses.WriteErr(w, err)
		return
	}

	responses.WriteOK(w, project)
}

// UpdateProjectStatus handles updating project status
// PATCH /api/v1/projects/{id}/status
func (h *Handler) UpdateProjectStatus(w http.ResponseWriter, r *http.Request) {
	var req UpdateProjectStatusRequest

	// Parse project ID
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
// DELETE /api/v1/projects/{id}
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

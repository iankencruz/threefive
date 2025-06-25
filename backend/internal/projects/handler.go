package project

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/iankencruz/threefive/internal/core/response"
	"github.com/iankencruz/threefive/internal/core/validators"
	"github.com/iankencruz/threefive/internal/generated"
	"github.com/jackc/pgx/v5/pgtype"
)

type Handler struct {
	Repo    Repository
	Service *ProjectService
	Logger  *slog.Logger
}

func NewHandler(q *generated.Queries, logger *slog.Logger) *Handler {

	repo := NewRepository(q)
	service := NewProjectService(repo)

	return &Handler{
		Repo:    repo,
		Service: service,
		Logger:  logger,
	}
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	projects, err := h.Service.repo.ListProjects(r.Context())
	if err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, "❌ Failed to list projects", err)
		return
	}
	response.WriteJSON(w, http.StatusOK, "✅ Admin: List Projects success", projects)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req generated.CreateProjectParams
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteJSON(w, http.StatusBadRequest, "❌ Invalid request body", err)
		return
	}

	v := validators.New()

	v.Require("title", req.Title)
	v.Require("slug", req.Slug)
	v.MatchPattern("slug", req.Slug, validators.SlugRX, "❌ Slug must be lowercase and can only contain letters, numbers, and hyphens")

	if !v.Valid() {
		response.WriteJSON(w, http.StatusBadRequest, "❌ Validation failed", v.Errors)
		return
	}

	project, err := h.Service.Create(r.Context(), req)
	if err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, "❌ Failed to create project", nil)
		return
	}
	response.WriteJSON(w, http.StatusCreated, "✅ Create Project success", project)
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.WriteJSON(w, http.StatusBadRequest, "Invalid project ID", nil)
		return
	}

	project, err := h.Service.repo.GetProjectByID(r.Context(), id)
	if err != nil {
		response.WriteJSON(w, http.StatusNotFound, "Project not found", err)
		return
	}
	response.WriteJSON(w, http.StatusOK, "✅ Admin:Get Project success", project)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	var req generated.UpdateProjectParams
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteJSON(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	v := validators.New()
	v.MatchPattern("slug", req.Slug, validators.SlugRX, "❌ Slug must be lowercase and can only contain letters, numbers, and hyphens")

	if !v.Valid() {
		response.WriteJSON(w, http.StatusInternalServerError, "❌ Validation failed", v.Errors)
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.WriteJSON(w, http.StatusBadRequest, "Invalid project ID", nil)
		return
	}
	req.ID = id

	var ts pgtype.Timestamptz
	_ = ts.Scan(time.Now())
	req.UpdatedAt = ts

	project, err := h.Service.Update(r.Context(), req)
	if err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, "Failed to update project", nil)
		return
	}
	response.WriteJSON(w, http.StatusOK, "✅ Update Project success", project)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.WriteJSON(w, http.StatusBadRequest, "Invalid project ID", nil)
		return
	}

	if err := h.Service.repo.DeleteProject(r.Context(), id); err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, "Failed to delete project", err)
		return
	}
	response.WriteJSON(w, http.StatusNoContent, "", nil)
}

func (h *Handler) AddMedia(w http.ResponseWriter, r *http.Request) {
	var req generated.AddMediaToProjectParams
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteJSON(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.WriteJSON(w, http.StatusBadRequest, "Invalid project ID", err)
		return
	}
	req.ProjectID = id
	if err := h.Service.repo.AddMediaToProject(r.Context(), req); err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, "Failed to add media", err)
		return
	}
	response.WriteJSON(w, http.StatusCreated, "✅ Add Media success", nil)
}

func (h *Handler) RemoveMedia(w http.ResponseWriter, r *http.Request) {
	var req generated.RemoveMediaFromProjectParams
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteJSON(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.WriteJSON(w, http.StatusBadRequest, "Invalid project ID", err)
		return
	}
	req.ProjectID = id
	if err := h.Service.repo.RemoveMediaFromProject(r.Context(), req); err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, "Failed to remove media", err)
		return
	}
	response.WriteJSON(w, http.StatusNoContent, "✅ Remove Media sucess", nil)
}

func (h *Handler) UpdateSortOrder(w http.ResponseWriter, r *http.Request) {
	var req generated.UpdateProjectMediaSortOrderParams
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteJSON(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.WriteJSON(w, http.StatusBadRequest, "Invalid project ID", err)
		return
	}

	req.ProjectID = id
	if err := h.Service.repo.UpdateProjectMediaSortOrder(r.Context(), req); err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, "Failed to update sort order", err)
		return
	}
	response.WriteJSON(w, http.StatusOK, "✅ Update Sort Order success", nil)
}

func (h *Handler) GetPublicProjects(w http.ResponseWriter, r *http.Request) {
	projects, err := h.Repo.ListPublishedProjects(r.Context())
	if err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, "Failed to fetch projects", nil)
		return
	}
	response.WriteJSON(w, http.StatusOK, "✅ Public projects list", projects)
}

package gallery

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/iankencruz/threefive/internal/core/response"
	"github.com/iankencruz/threefive/internal/generated"
)

type Handler struct {
	Repo    Repository
	Service *GalleryService
	Logger  *slog.Logger
}

func NewHandler(q *generated.Queries, logger *slog.Logger) *Handler {
	repo := NewRepository(q)
	service := NewGalleryService(repo)

	return &Handler{
		Repo:    repo,
		Service: service,
		Logger:  logger,
	}
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	gallery, err := h.Service.repo.ListGalleries(r.Context())
	if err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, "❌ Failed to list galleries", err)
		return
	}
	response.WriteJSON(w, http.StatusOK, "✅ Admin: List Galleries success", gallery)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req generated.CreateGalleryParams
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		err = errors.New("❌ Invalid request body")
		response.WriteJSON(w, http.StatusBadRequest, "❌ Invalid request body", err)
		return
	}
	gallery, err := h.Service.Create(r.Context(), req)
	if err != nil {
		err = errors.New("Failed to create a Gallery.")
		response.WriteJSON(w, http.StatusBadRequest, "❌ Invalid request body", err)
		return
	}

	response.WriteJSON(w, http.StatusCreated, " Create Gallery Successful", gallery)
}

func (h *Handler) GetAdminGallery(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	if slug == "" {
		response.WriteJSON(w, http.StatusBadRequest, "Missing gallery slug", nil)
		return
	}

	gallery, err := h.Service.repo.GetGalleryBySlug(r.Context(), slug)
	if err != nil {
		response.WriteJSON(w, http.StatusNotFound, "Gallery not found", err)
		return
	}
	response.WriteJSON(w, http.StatusOK, "Admin: Get Gallery Success", gallery)
}

func (h *Handler) AddMedia(w http.ResponseWriter, r *http.Request) {
	var req generated.AddMediaToGalleryParams
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteJSON(w, http.StatusBadRequest, "Invalid Request Body", err)
		return
	}

	slug := chi.URLParam(r, "slug")

	gm, err := h.Repo.GetGalleryBySlug(r.Context(), slug)
	if err != nil {
		response.WriteJSON(w, http.StatusNotFound, "Project not found", nil)
		return
	}
	req.GalleryID = gm.ID
	if err := h.Service.repo.AddMediaToGallery(r.Context(), req); err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, "Failed to add media", err)
		return
	}
	response.WriteJSON(w, http.StatusCreated, "✅ Add Media success", nil)
}

func (h *Handler) RemoveMedia(w http.ResponseWriter, r *http.Request) {
	var req generated.RemoveMediaFromGalleryParams
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteJSON(w, http.StatusBadRequest, "Invalid request body", err)
		return

	}

	slug := chi.URLParam(r, "slug")

	pm, err := h.Service.repo.GetGalleryBySlug(r.Context(), slug)
	fmt.Printf("projectWithMedia: %v", pm)
	if err != nil {
		response.WriteJSON(w, http.StatusNotFound, "Project not found", nil)
		return
	}
	req.GalleryID = pm.ID
	if err := h.Service.repo.RemoveMediaFromGallery(r.Context(), req); err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, "Failed to remove media", err)
		return
	}
	response.WriteJSON(w, http.StatusNoContent, "✅ Remove Media sucess", nil)
}

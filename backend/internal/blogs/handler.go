package blogs

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/iankencruz/threefive/internal/core/response"
	"github.com/iankencruz/threefive/internal/generated"
)

type Handler struct {
	Repo    Repository
	Service *BlogService
}

func NewHandler(q *generated.Queries) *Handler {
	repo := NewRepository(q)
	service := NewBlogService(repo)
	return &Handler{
		Repo:    repo,
		Service: service,
	}
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	blog, err := h.Service.GetGalleryBySlug(r.Context(), slug)
	if err != nil {
		response.WriteJSON(w, http.StatusBadRequest, "Blog not found", err.Error())
		return
	}

	response.WriteJSON(w, http.StatusOK, "✅ Blog loaded", blog)
}

func (h *Handler) GetAdminGallery(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	if slug == "" {
		response.WriteJSON(w, http.StatusBadRequest, "Missing blog slug", nil)
		return
	}

	gallery, err := h.Service.repo.GetBySlug(r.Context(), slug)
	if err != nil {
		response.WriteJSON(w, http.StatusNotFound, "Blog not found", err)
		return
	}
	response.WriteJSON(w, http.StatusOK, "Admin: Get Blog Success", gallery)
}

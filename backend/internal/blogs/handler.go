package blogs

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/iankencruz/threefive/internal/blocks"
	"github.com/iankencruz/threefive/internal/core/response"
	"github.com/iankencruz/threefive/internal/generated"
)

type Handler struct {
	Repo    Repository
	Service *BlogService
}

func NewHandler(q *generated.Queries, blockRepo *blocks.Repository) *Handler {
	repo := NewRepository(q, blockRepo)
	service := NewBlogService(repo)
	return &Handler{
		Repo:    repo,
		Service: service,
	}
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	blog, err := h.Service.GetBySlug(r.Context(), slug)
	if err != nil {
		response.WriteJSON(w, http.StatusBadRequest, "Blog not found", err.Error())
		return
	}

	response.WriteJSON(w, http.StatusOK, "✅ Blog loaded", blog)
}

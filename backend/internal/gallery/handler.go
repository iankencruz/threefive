package gallery

import (
	"log/slog"
	"net/http"

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
	projects, err := h.Service.repo.ListGalleries(r.Context())
	if err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, "❌ Failed to list galleries", err)
		return
	}
	response.WriteJSON(w, http.StatusOK, "✅ Admin: List Galleries success", projects)
}

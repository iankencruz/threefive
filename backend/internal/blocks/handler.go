package blocks

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/iankencruz/threefive/internal/core/response"
	"github.com/iankencruz/threefive/internal/generated"
)

type Handler struct {
	Repo    *Repository
	Service *Service
	Logger  *slog.Logger
}

func NewHandler(q *generated.Queries, logger *slog.Logger) *Handler {
	repo := NewRepository(q)
	service := NewService(repo)

	return &Handler{
		Repo:    repo,
		Service: service,
		Logger:  logger,
	}
}

// POST /api/v1/admin/blocks
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Block generated.Block `json:"block"`
		Props any             `json:"props"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		response.WriteJSON(w, http.StatusBadRequest, "Invalid JSON", err)
		return
	}

	err := h.Service.CreateBlock(r.Context(), payload.Block, payload.Props)
	if err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, "Failed to create block", err)
		return
	}

	response.WriteJSON(w, http.StatusCreated, "✅ Block created", payload.Block)
}

// PUT /api/v1/admin/blocks/{id}
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.WriteJSON(w, http.StatusBadRequest, "Invalid block ID", err)
		return
	}

	var payload struct {
		Block generated.Block `json:"block"`
		Props any             `json:"props"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		response.WriteJSON(w, http.StatusBadRequest, "Invalid JSON", err)
		return
	}

	payload.Block.ID = id

	err = h.Service.UpdateBlock(r.Context(), payload.Block, payload.Props)
	if err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, "Failed to update block", err)
		return
	}

	response.WriteJSON(w, http.StatusOK, "✅ Block updated", nil)
}

// DELETE /api/v1/admin/blocks/{id}
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.WriteJSON(w, http.StatusBadRequest, "Invalid block ID", err)
		return
	}

	block, err := h.Service.Repo.GetBlockByID(r.Context(), id)
	if err != nil {
		response.WriteJSON(w, http.StatusNotFound, "Block not found", err)
		return
	}

	if err := h.Service.DeleteBlockByID(r.Context(), id, block.Type); err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, "Failed to delete block", err)
		return
	}

	response.WriteJSON(w, http.StatusOK, "✅ Block deleted", nil)
}

package pages

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/iankencruz/threefive/internal/blocks"
	"github.com/iankencruz/threefive/internal/core/response"
	"github.com/iankencruz/threefive/internal/core/validators"
	"github.com/iankencruz/threefive/internal/generated"
)

type Handler struct {
	Repo         Repository
	Service      *PageService
	BlockService *blocks.Service
	BlockRepo    *blocks.Repository
	Logger       *slog.Logger
}

func NewHandler(q *generated.Queries, blockRepo *blocks.Repository, blockService *blocks.Service, logger *slog.Logger) *Handler {
	repo := NewRepository(q)
	service := NewPageService(repo, blockRepo, *blockService)
	return &Handler{
		Repo:         repo,
		Service:      service,
		BlockService: blockService,
		BlockRepo:    blockRepo,
		Logger:       logger,
	}
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	pages, err := h.Service.Repo.ListPages(r.Context())
	if err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, "Failed to list pages", err)
		return
	}
	response.WriteJSON(w, http.StatusOK, "Admin: List Pages success", pages)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreatePageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteJSON(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	v := validators.New()
	v.Require("title", req.Page.Title)

	if !v.Valid() {
		response.WriteJSON(w, http.StatusBadRequest, fmt.Sprintf("Validation failed: %v\n", v.Errors), v.Errors)
		return
	}

	page, err := h.Service.Create(r.Context(), req.Page)
	if err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, "Failed to create page", nil)
		return
	}
	response.WriteJSON(w, http.StatusCreated, "Create Page success", page)
}

// Getter
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	if slug == "" {
		response.WriteJSON(w, http.StatusBadRequest, "Missing page slug", nil)
		return
	}

	page, err := h.Service.Repo.GetPageBySlug(r.Context(), slug)
	if err != nil {
		response.WriteJSON(w, http.StatusNotFound, "No valid Page found", err)
		return
	}

	blocks, err := h.Service.GetPageBlocks(r.Context(), page.ID)
	if err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, fmt.Sprintf("Failed to fetch Page Blocks: %v", err), err.Error())
		return
	}

	result := PageWithBlocks{
		Page:   *page,
		Blocks: blocks,
	}

	response.WriteJSON(w, http.StatusOK, "Admin:Get Page success", result)
}

// Update Page

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	var req UpdatePageWithBlocksRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteJSON(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	slug := chi.URLParam(r, "slug")
	page, err := h.Service.Repo.GetPageBySlug(r.Context(), slug)
	if err != nil {
		response.WriteJSON(w, http.StatusNotFound, "Page not found", err)
		return
	}
	req.Page.ID = page.ID

	updated, err := h.Service.UpdateWithBlocks(r.Context(), req)
	if err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, "Failed to update page.", err.Error())
		return
	}

	response.WriteJSON(w, http.StatusOK, "Update Page success", updated)
}

func (h *Handler) SortBlocks(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	if slug == "" {
		response.WriteJSON(w, http.StatusBadRequest, "Missing page slug", nil)
		return
	}

	var updates []struct {
		ID        uuid.UUID `json:"id"`
		SortOrder int       `json:"sort_order"`
	}
	if err := response.DecodeJSON(w, r, &updates); err != nil {
		response.WriteJSON(w, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	for _, u := range updates {
		args := generated.UpdateBlockSortOrderParams{
			ID:        u.ID,
			SortOrder: int32(u.SortOrder),
		}

		fmt.Printf("  - BlockID: %s → Order: %d\n", args.ID, args.SortOrder)

		if err := h.BlockRepo.UpdateBlockSortOrder(r.Context(), args); err != nil {
			response.WriteJSON(w, http.StatusInternalServerError, "Failed to update block", err.Error())
			return
		}
	}

	response.WriteJSON(w, http.StatusOK, "Block sort updated", nil)
}

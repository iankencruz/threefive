package pages

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/iankencruz/threefive/internal/core/response"
	"github.com/iankencruz/threefive/internal/core/validators"
	"github.com/iankencruz/threefive/internal/generated"
)

type Handler struct {
	Repo    Repository
	Service *PageService
	Logger  *slog.Logger
}

func NewHandler(q *generated.Queries, logger *slog.Logger) *Handler {
	repo := NewRepository(q)
	service := NewPageService(repo)
	return &Handler{
		Repo:    repo,
		Service: service,
		Logger:  logger,
	}
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	pages, err := h.Service.Repo.ListPages(r.Context())
	if err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, "❌ Failed to list pages", err)
		return
	}
	response.WriteJSON(w, http.StatusOK, "✅ Admin: List Pages success", pages)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req generated.CreatePageParams
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteJSON(w, http.StatusBadRequest, "❌ Invalid request body", err)
		return
	}

	v := validators.New()
	v.Require("title", req.Title)

	if !v.Valid() {
		response.WriteJSON(w, http.StatusBadRequest, "❌ Validation failed", v.Errors)
		return
	}

	page, err := h.Service.Create(r.Context(), req)
	if err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, "❌ Failed to create page", nil)
		return
	}
	response.WriteJSON(w, http.StatusCreated, "✅ Create Page success", page)
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	if slug == "" {
		response.WriteJSON(w, http.StatusBadRequest, "Missing page slug", nil)
		return
	}

	page, err := h.Service.Repo.GetPageBySlug(r.Context(), slug)
	fmt.Println("Page:", page)
	if err != nil {
		response.WriteJSON(w, http.StatusNotFound, "No valid Page found", err)
		return
	}
	response.WriteJSON(w, http.StatusOK, "✅ Admin:Get Page success", page)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	var req generated.UpdatePageParams
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fmt.Printf("Req err: %v\n", err)

		response.WriteJSON(w, http.StatusBadRequest, "❌ Invalid request body", err)
		return
	}

	fmt.Printf("Req: %v\n", req)

	slug := chi.URLParam(r, "slug")
	if slug == "" {
		response.WriteJSON(w, http.StatusBadRequest, "Missing page slug", nil)
		return
	}

	page, err := h.Service.Repo.GetPageBySlug(r.Context(), slug)
	if err != nil {
		response.WriteJSON(w, http.StatusNotFound, "Page not found", nil)
		return
	}

	req.ID = page.ID

	updated, err := h.Service.Update(r.Context(), req)
	if err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, "❌ Failed to update page", err)
		return
	}
	response.WriteJSON(w, http.StatusOK, "✅ Update Page success", updated)
}

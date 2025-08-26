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

func (h *Handler) HomePage(w http.ResponseWriter, r *http.Request) {
	// data := "Home Page"
	page, err := h.Service.Repo.GetPageBySlug(r.Context(), "home")
	if err != nil {
		h.Logger.Error("No valid Page found")
		response.WriteJSON(w, http.StatusNotFound, "No valid Page found", err)
		return
	}
	response.WriteJSON(w, http.StatusOK, "Home Page Public:", page)
}

func (h *Handler) ContactPage(w http.ResponseWriter, r *http.Request) {
	page, err := h.Service.Repo.GetPageBySlug(r.Context(), "contact")
	if err != nil {
		response.WriteJSON(w, http.StatusNotFound, "No valid Page found", err)
		return
	}

	response.WriteJSON(w, http.StatusOK, "Admin:Get Page success", page)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	sort := r.URL.Query().Get("sort") // optional query param
	pages, err := h.Repo.ListPages(r.Context(), sort)
	if err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, "Failed to fetch pages", nil)
		return
	}

	response.WriteJSON(w, http.StatusOK, "✅ Pages fetched", pages)
}

func (h *Handler) AboutPage(w http.ResponseWriter, r *http.Request) {
	page, err := h.Service.Repo.GetPageBySlug(r.Context(), "about")
	if err != nil {
		h.Logger.Error("No valid Page found")
		response.WriteJSON(w, http.StatusNotFound, "No valid Page found", err)
		return
	}

	response.WriteJSON(w, http.StatusOK, "Admin:Get Page success", page)
}

// Getter page by slug
func (h *Handler) GetAdminPages(w http.ResponseWriter, r *http.Request) {
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

	response.WriteJSON(w, http.StatusOK, "Admin:Get Page success", page)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req generated.CreatePageParams
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteJSON(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	v := validators.New()
	v.Require("title", req.Title)

	if !v.Valid() {
		response.WriteJSON(w, http.StatusBadRequest, fmt.Sprintf("Validation failed: %v\n", v.Errors), v.Errors)
		return
	}

	page, err := h.Service.Create(r.Context(), req)
	if err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, "Failed to create page", err.Error())
		return
	}
	response.WriteJSON(w, http.StatusCreated, "Create Page success", page)
}

// Update Page

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	var req generated.UpdatePageParams
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
	req.ID = page.ID

	page, err = h.Service.Update(r.Context(), req)
	if err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, "Failed to update page.", err.Error())
		return
	}

	response.WriteJSON(w, http.StatusOK, "Update Page success", page)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	if slug == "" {
		response.WriteJSON(w, http.StatusBadRequest, "Missing page slug", nil)
		return
	}

	page, err := h.Service.Repo.GetPageBySlug(r.Context(), slug)
	if err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, "failed to retrive page via slug", err.Error())
		return
	}

	if err := h.Service.Repo.DeletePage(r.Context(), page.ID); err != nil {
		response.WriteJSON(w, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}
	response.WriteJSON(w, http.StatusNoContent, "", nil)
}

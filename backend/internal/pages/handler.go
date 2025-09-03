package pages

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/iankencruz/threefive/internal/core/response"
	"github.com/iankencruz/threefive/internal/core/validators"
	"github.com/iankencruz/threefive/internal/gallery"
	"github.com/iankencruz/threefive/internal/generated"
)

type Handler struct {
	Repo           Repository
	Service        *PageService
	Logger         *slog.Logger
	GalleryService *gallery.GalleryService
}

func NewHandler(q *generated.Queries, galleryService *gallery.GalleryService, logger *slog.Logger) *Handler {
	repo := NewRepository(q)
	service := NewPageService(repo)
	return &Handler{
		Repo:           repo,
		Service:        service,
		Logger:         logger,
		GalleryService: galleryService,
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
	sortParam := r.URL.Query().Get("sort")

	pages, err := h.Service.ListPages(r.Context(), sortParam)
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

// GET /api/v1/admin/pages/{slug}/galleries
func (h *Handler) ListPageGalleries(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	page, err := h.Repo.GetPageBySlug(r.Context(), slug)
	if err != nil {
		response.WriteJSON(w, http.StatusNotFound, "Page not found", nil)
		return
	}

	galleries, err := h.GalleryService.ListByPage(r.Context(), page.ID)
	if err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, "Failed to fetch galleries", nil)
		return
	}

	// Always return [] instead of null
	if galleries == nil {
		galleries = []generated.Gallery{}
	}

	response.WriteJSON(w, http.StatusOK, "✅ Galleries fetched", galleries)
}

// POST /api/v1/admin/pages/{slug}/galleries
func (h *Handler) LinkGallery(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	page, err := h.Repo.GetPageBySlug(r.Context(), slug)
	if err != nil {
		response.WriteJSON(w, http.StatusNotFound, "Page not found", nil)
		return
	}

	var input struct {
		GalleryID uuid.UUID `json:"gallery_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.WriteJSON(w, http.StatusBadRequest, "Invalid payload", nil)
		return
	}

	if err := h.GalleryService.LinkToPage(r.Context(), input.GalleryID, page.ID); err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, "Failed to link gallery", nil)
		return
	}

	response.WriteJSON(w, http.StatusOK, "✅ Gallery linked", nil)
}

// DELETE /api/v1/admin/pages/{slug}/galleries/{galleryID}
func (h *Handler) UnlinkGallery(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	page, err := h.Repo.GetPageBySlug(r.Context(), slug)
	if err != nil {
		response.WriteJSON(w, http.StatusNotFound, "Page not found", nil)
		return
	}

	galleryIDStr := chi.URLParam(r, "galleryID")
	galleryID, err := uuid.Parse(galleryIDStr)
	if err != nil {
		response.WriteJSON(w, http.StatusBadRequest, "Invalid gallery ID", nil)
		return
	}

	if err := h.GalleryService.UnlinkFromPage(r.Context(), galleryID, page.ID); err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, "Failed to unlink gallery", nil)
		return
	}

	response.WriteJSON(w, http.StatusOK, "✅ Gallery unlinked", nil)
}

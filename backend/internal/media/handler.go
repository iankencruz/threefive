package media

import (
	"log/slog"
	"math"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/iankencruz/threefive/backend/internal/core/response"
	"github.com/iankencruz/threefive/backend/internal/core/s3"
	"github.com/iankencruz/threefive/backend/internal/generated"
)

type Handler struct {
	Repo     Repository
	Service  *Service
	Uploader *s3.Uploader
	Logger   *slog.Logger
}

func NewHandler(q *generated.Queries, logger *slog.Logger, uploader *s3.Uploader) *Handler {
	repo := NewRepository(q)
	service := NewService(repo, uploader)

	return &Handler{
		Repo:     repo,
		Service:  service,
		Uploader: uploader,
		Logger:   logger,
	}
}

// UploadMediaHandler handles POST /admin/media/upload
func (h *Handler) UploadMediaHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(20 << 20) // 20MB max
	if err != nil {
		h.Logger.Error("failed to parse multipart form", "err", err)
		response.WriteJSON(w, http.StatusBadRequest, "Invalid form data", err)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		h.Logger.Error("file missing in form", "err", err)
		response.WriteJSON(w, http.StatusBadRequest, "File is required", err)
		return
	}
	defer file.Close()

	title := r.FormValue("title")
	alt := r.FormValue("alt")
	sortStr := r.FormValue("sort")
	sortOrder := int32(0)
	if sortStr != "" {
		if n, err := strconv.Atoi(sortStr); err == nil {
			sortOrder = int32(n)
		}
	}

	media, err := h.Service.UploadMedia(r.Context(), file, header, title, alt, sortOrder)
	if err != nil {
		h.Logger.Error("media upload failed", "err", err)
		response.WriteJSON(w, http.StatusInternalServerError, "Upload failed", nil)
		return
	}

	response.WriteJSON(w, http.StatusOK, "✅ Upload Success", media)
}

func (h *Handler) ListMediaHandler(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	sizeStr := r.URL.Query().Get("limit") // fix inconsistent name
	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(sizeStr)

	if page < 1 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	// Fetch paginated media
	mediaList, err := h.Repo.ListPaginated(r.Context(), int32(pageSize), int32(offset))
	if err != nil {
		h.Logger.Error("failed to list media", "err", err)
		response.WriteJSON(w, http.StatusInternalServerError, "Failed to list media", nil)
		return
	}
	if mediaList == nil {
		mediaList = []generated.Media{}
	}

	// ✅ Add full S3 URL to each media item
	// bucket := strings.TrimSuffix(h.Uploader.BucketName, "/") + "/"
	for i := range mediaList {
		if mediaList[i].ThumbnailUrl != nil && *mediaList[i].ThumbnailUrl != "" {
			fullThumb := *mediaList[i].ThumbnailUrl
			mediaList[i].ThumbnailUrl = &fullThumb
		}
		if mediaList[i].MediumUrl != nil && *mediaList[i].MediumUrl != "" {
			fullMedium := *mediaList[i].MediumUrl
			mediaList[i].MediumUrl = &fullMedium
		}
	}

	// Fetch total count
	totalCount, err := h.Repo.CountMedia(r.Context())
	if err != nil {
		h.Logger.Error("failed to count media", "err", err)
		response.WriteJSON(w, http.StatusInternalServerError, "Failed to count media", nil)
		return
	}

	totalPages := int(math.Ceil(float64(totalCount) / float64(pageSize)))

	// Return full paginated result
	result := map[string]interface{}{
		"items":       mediaList,
		"total_count": totalCount,
		"total_pages": totalPages,
		"page":        page,
	}

	response.WriteJSON(w, http.StatusOK, "✅ Media list fetched", result)
}

func (h *Handler) DeleteMediaHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		response.WriteJSON(w, http.StatusBadRequest, "Missing media ID", nil)
		return
	}

	id, err := parseUUIDParam(r, idStr)
	if err != nil {
		http.Error(w, "invalid UUID", http.StatusBadRequest)
		return
	}

	err = h.Repo.Delete(r.Context(), id)
	if err != nil {
		h.Logger.Error("failed to delete media", "id", id, "err", err)
		response.WriteJSON(w, http.StatusInternalServerError, "Delete failed", nil)
		return
	}

	response.WriteJSON(w, http.StatusOK, "✅ Deleted", nil)
}

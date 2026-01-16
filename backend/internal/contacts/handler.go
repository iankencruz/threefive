// internal/contacts/handler.go
package contacts

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/iankencruz/threefive/internal/shared/email"
	"github.com/iankencruz/threefive/internal/shared/responses"
	"github.com/iankencruz/threefive/internal/shared/sqlc"
	"github.com/iankencruz/threefive/internal/shared/validation"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Handler struct {
	Service *Service
}

func NewHandler(db *pgxpool.Pool, queries *sqlc.Queries, emailService *email.Service) *Handler {
	return &Handler{
		Service: NewService(db, queries, emailService),
	}
}

// CreateContact handles POST /api/contact
func (h *Handler) CreateContact(w http.ResponseWriter, r *http.Request) {
	// Parse and validate request
	req, err := validation.ParseAndValidate[*ContactRequest](r)
	if err != nil {
		responses.WriteErr(w, err)
		return
	}

	// Get IP address and user agent
	ipAddress := getIPAddress(r)
	userAgent := r.UserAgent()

	contact, err := h.Service.CreateContact(r.Context(), req, ipAddress, userAgent)
	if err != nil {
		responses.WriteErr(w, err)
		return
	}

	response := map[string]interface{}{
		"message": "Contact form submitted successfully",
		"contact": contact,
	}

	responses.WriteCreated(w, response)
}

// GetContact handles GET /api/admin/contacts/{id}
func (h *Handler) GetContact(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		responses.WriteErr(w, err)
		return
	}

	contact, err := h.Service.GetContactByID(r.Context(), id)
	if err != nil {
		responses.WriteErr(w, err)
		return
	}

	responses.WriteOK(w, contact)
}

// ListContacts handles GET /api/admin/contacts
func (h *Handler) ListContacts(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	params := ListContactsParams{
		Status:  r.URL.Query().Get("status"),
		OrderBy: r.URL.Query().Get("order_by"),
	}

	// Parse limit
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err == nil {
			params.Limit = int32(limit)
		}
	}

	// Parse offset
	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		offset, err := strconv.Atoi(offsetStr)
		if err == nil {
			params.Offset = int32(offset)
		}
	}

	contacts, err := h.Service.ListContacts(r.Context(), params)
	if err != nil {
		responses.WriteErr(w, err)
		return
	}

	responses.WriteOK(w, contacts)
}

// UpdateContactStatus handles PATCH /api/admin/contacts/{id}/status
func (h *Handler) UpdateContactStatus(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		responses.WriteErr(w, err)
		return
	}

	// Parse and validate request
	req, err := validation.ParseAndValidate[*UpdateContactStatusRequest](r)
	if err != nil {
		responses.WriteErr(w, err)
		return
	}

	contact, err := h.Service.UpdateContactStatus(r.Context(), id, req.Status)
	if err != nil {
		responses.WriteErr(w, err)
		return
	}

	response := map[string]interface{}{
		"message": "Contact status updated successfully",
		"contact": contact,
	}

	responses.WriteOK(w, response)
}

// DeleteContact handles DELETE /api/admin/contacts/{id}
func (h *Handler) DeleteContact(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		responses.WriteErr(w, err)
		return
	}

	if err := h.Service.DeleteContact(r.Context(), id); err != nil {
		responses.WriteErr(w, err)
		return
	}

	response := map[string]string{
		"message": "Contact deleted successfully",
	}

	responses.WriteOK(w, response)
}

// getIPAddress extracts the real IP address from the request
func getIPAddress(r *http.Request) string {
	// Try X-Forwarded-For header first
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		ips := strings.Split(xff, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	// Try X-Real-IP header
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return strings.TrimSpace(xri)
	}

	// Fall back to RemoteAddr
	ip := r.RemoteAddr
	if idx := strings.LastIndex(ip, ":"); idx != -1 {
		ip = ip[:idx]
	}
	return ip
}

package auth

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/iankencruz/threefive/backend/internal/core/sessions"
)

type Handler struct {
	Service        Service
	SessionManager *sessions.Manager
	Logger         *slog.Logger
}

func (h *Handler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: implement handler logic
}

func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: implement handler logic
}

func (h *Handler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: implement handler logic
}

func (h *Handler) GetUserID(r *http.Request) (int32, error) {
	return h.SessionManager.GetUserID(r)
}

func (h *Handler) LoadUser(ctx context.Context, userID int32) (any, error) {
	return h.Service.GetUserByID(ctx, userID)
}

func (h *Handler) GetAuthenticatedUser(w http.ResponseWriter, r *http.Request) {
	userID, err := h.SessionManager.GetUserID(r)
	w.Header().Set("Content-Type", "application/json")

	if err != nil || userID == 0 {
		json.NewEncoder(w).Encode(map[string]interface{}{"user": nil})
		return
	}

	user, err := h.Service.GetUserByID(r.Context(), int32(userID))
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{"user": nil})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{"user": user})
}

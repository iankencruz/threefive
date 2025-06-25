package auth

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/iankencruz/threefive/internal/core/contextkeys"
	"github.com/iankencruz/threefive/internal/core/errors"
	"github.com/iankencruz/threefive/internal/core/response"
	"github.com/iankencruz/threefive/internal/core/sessions"
	"github.com/iankencruz/threefive/internal/core/validators"
	"github.com/iankencruz/threefive/internal/generated"
)

type Handler struct {
	Repo           Repository
	Service        Service
	SessionManager *sessions.Manager
	Logger         *slog.Logger
}

func NewHandler(q *generated.Queries, sessionManager *sessions.Manager, logger *slog.Logger) *Handler {

	repo := NewAuthRepository(q)
	service := NewAuthService(repo)

	return &Handler{
		Repo:           repo,
		Service:        service,
		SessionManager: sessionManager,
		Logger:         logger,
	}
}

func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		errResp := errors.BadRequest("Invalid login payload")
		response.WriteJSON(w, errResp.Code, errResp.Message, errResp)
		return
	}

	v := validators.New()
	v.Require("email", input.Email)
	v.MatchPattern("email", input.Email, validators.EmailRX, "Must be a valid email address")
	v.Require("password", input.Password)

	if !v.Valid() {
		errResp := errors.BadRequest("Validation failed")
		response.WriteJSON(w, errResp.Code, errResp.Message, v.Errors)
		return
	}

	// Use the service to Check the credentials
	user, err := h.Service.Login(r.Context(), input.Email, input.Password)
	if err != nil {
		errResp := errors.Internal(err.Error())
		response.WriteJSON(w, errResp.Code, errResp.Message, errResp)
		return
	}

	if err := h.SessionManager.SetUserID(w, r, user.ID); err != nil {
		errResp := errors.Internal("Failed to set session")
		response.WriteJSON(w, errResp.Code, errResp.Message, errResp)
		return
	}

	response.WriteJSON(w, http.StatusOK, "Logged in", map[string]any{"user": user})
}

func (h *Handler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		response.WriteJSON(w, http.StatusMethodNotAllowed, "Method not allowed", nil)
		return
	}

	var input struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
		Password  string `json:"password"`
	}

	if err := response.DecodeJSON(w, r, &input); err != nil {
		errResp := errors.BadRequest("Invalid JSON payload: " + err.Error())
		response.WriteJSON(w, errResp.Code, errResp.Message, nil)
		return
	}

	input.Email = strings.ToLower(strings.TrimSpace(input.Email))

	v := validators.New()
	v.Require("first_name", input.FirstName)
	v.Require("last_name", input.LastName)
	v.Require("email", input.Email)
	v.MatchPattern("email", input.Email, validators.EmailRX, "Invalid email address")
	v.Require("password", input.Password)
	v.MatchPattern("password", input.Password, validators.UppercaseRX, "Must include at least one uppercase letter")
	v.MatchPattern("password", input.Password, validators.NumberRX, "Must include at least one number")

	if !v.Valid() {
		errResp := errors.BadRequest("Validation failed")
		response.WriteJSON(w, errResp.Code, errResp.Message, v.Errors)
		return
	}

	user, err := h.Service.Register(r.Context(), input.FirstName, input.LastName, input.Email, input.Password)
	if err != nil {
		v.Errors["email"] = "Registration failed"
		response.WriteJSON(w, http.StatusBadRequest, "Registration failed", v.Errors)
		return
	}

	if err := h.SessionManager.SetUserID(w, r, user.ID); err != nil {
		errResp := errors.Internal("Session error")
		response.WriteJSON(w, errResp.Code, errResp.Message, nil)
		return
	}

	response.WriteJSON(w, http.StatusOK, "Registration successful", map[string]any{
		"user": map[string]any{
			"id":         user.ID,
			"firstName":  user.FirstName,
			"lastName":   user.LastName,
			"email":      user.Email,
			"created_at": user.CreatedAt,
		},
	})
}

func (h *Handler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: implement handler logic
	if err := h.SessionManager.Clear(w, r); err != nil {
		errResp := errors.Internal("Failed to clear session")
		response.WriteJSON(w, errResp.Code, errResp.Message, errResp)
		return
	}

	response.WriteJSON(w, http.StatusOK, "Logged out successfully", nil)

}

func (h *Handler) GetUserID(r *http.Request) (int32, error) {
	return h.SessionManager.GetUserID(r)
}

func (h *Handler) LoadUser(ctx context.Context, userID int32) (any, error) {
	return h.Service.GetUserByID(ctx, userID)
}

func (h *Handler) GetAuthenticatedUser(w http.ResponseWriter, r *http.Request) {
	userID, err := h.SessionManager.GetUserID(r)
	if err != nil || userID == 0 {
		_ = response.WriteJSON(w, http.StatusUnauthorized, "unauthorised", nil)
		return
	}

	user, err := h.Service.GetUserByID(r.Context(), int32(userID))
	if err != nil {
		_ = response.WriteJSON(w, http.StatusInternalServerError, "could not retrieve user", nil)
		return
	}

	_ = response.WriteJSON(w, http.StatusOK, "", user)
}

func (h *Handler) MeHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	sessionID := ctx.Value(contextkeys.SessionID).(string)
	if sessionID == "" {
		http.Error(w, "Unauthorised", http.StatusUnauthorized)
		return
	}

	userIDStr, err := h.SessionManager.GetString(ctx, sessionID, "userID")
	if err != nil || userIDStr == "" {
		http.Error(w, "Unauthorised", http.StatusUnauthorized)
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusInternalServerError)
		return
	}

	user, err := h.Service.GetUserByID(ctx, int32(userID))
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	response := map[string]any{"user": user}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response)
}

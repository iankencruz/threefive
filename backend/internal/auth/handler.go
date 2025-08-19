package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/google/uuid"
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
		v.Errors["email"] = fmt.Sprintf("Email Registration failed: %v", err)
		response.WriteJSON(w, http.StatusBadRequest, "Registration failed", v.Errors)
		return
	}

	fmt.Printf("User.ID: %v", user.ID)
	if err := h.SessionManager.SetUserID(w, r, user.ID); err != nil {
		errResp := errors.Internal(fmt.Sprintf("Session error: %v", err))
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

func (h *Handler) GetUserID(r *http.Request) (uuid.UUID, error) {
	return h.SessionManager.GetUserID(r)
}

func (h *Handler) LoadUser(ctx context.Context, userID uuid.UUID) (any, error) {
	return h.Service.GetUserByID(ctx, userID)
}

func (h *Handler) GetAuthenticatedUser(w http.ResponseWriter, r *http.Request) {
	userID, err := h.SessionManager.GetUserID(r)
	if err != nil || userID == uuid.Nil {
		_ = response.WriteJSON(w, http.StatusUnauthorized, "unauthorised", nil)
		return
	}

	user, err := h.Service.GetUserByID(r.Context(), userID)
	if err != nil {
		_ = response.WriteJSON(w, http.StatusInternalServerError, "could not retrieve user", nil)
		return
	}

	_ = response.WriteJSON(w, http.StatusOK, "", user)
}

func (h *Handler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.Repo.ListUsers(r.Context())
	if err != nil {
		_ = response.WriteJSON(w, http.StatusInternalServerError, "Coult not retrieve user", err)
	}

	_ = response.WriteJSON(w, http.StatusOK, "", users)
}

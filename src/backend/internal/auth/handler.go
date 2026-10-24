package auth

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/alexedwards/scs/v2"

	"github.com/zeeptobean/introse-portal67/src/backend/internal/httpx"
)

const sessionUserIDKey = "userID"

// Handler exposes HTTP endpoints for login, logout, and account identity.
type Handler struct {
	service  *Service
	sessions *scs.SessionManager
}

func NewHandler(service *Service, sessions *scs.SessionManager) *Handler {
	return &Handler{service: service, sessions: sessions}
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var in loginRequest
	if err := httpx.DecodeJSON(w, r, &in); err != nil {
		httpx.BadRequest(w, "invalid JSON body")
		return
	}
	if in.Email == "" || in.Password == "" {
		httpx.BadRequest(w, "email and password are required")
		return
	}

	ctx := r.Context()
	u, err := h.service.Authenticate(ctx, in.Email, in.Password)
	switch {
	case errors.Is(err, ErrInvalidCredentials):
		httpx.Error(w, http.StatusUnauthorized, err.Error())
		return
	case errors.Is(err, ErrInactive):
		httpx.Error(w, http.StatusForbidden, err.Error())
		return
	case err != nil:
		slog.ErrorContext(ctx, "log in", "error", err)
		httpx.InternalError(w)
		return
	}

	// A fresh token on privilege change prevents session fixation.
	if err := h.sessions.RenewToken(ctx); err != nil {
		slog.ErrorContext(ctx, "renew session token", "error", err)
		httpx.InternalError(w)
		return
	}
	h.sessions.Put(ctx, sessionUserIDKey, u.ID)
	httpx.WriteJSON(w, http.StatusOK, u)
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	if err := h.sessions.Destroy(r.Context()); err != nil {
		slog.ErrorContext(r.Context(), "destroy session", "error", err)
		httpx.InternalError(w)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, map[string]string{"status": "logged out"})
}

func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	u, ok := UserFromContext(r.Context())
	if !ok {
		httpx.NotLoggedIn(w)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, u)
}

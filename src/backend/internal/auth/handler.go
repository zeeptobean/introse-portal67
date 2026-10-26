package auth

import (
	"errors"
	"log/slog"
	"math"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/alexedwards/scs/v2"

	"github.com/zeeptobean/introse-portal67/src/backend/internal/httpx"
)

const sessionUserIDKey = "userID"

// Login attempts are limited per client IP and per email. Successful logins
// clear the email's counter but still count against the IP.
const (
	loginWindow           = 15 * time.Minute
	loginAttemptsPerIP    = 30
	loginAttemptsPerEmail = 5
)

// Handler exposes HTTP endpoints for login, logout, and account identity.
type Handler struct {
	service      *Service
	sessions     *scs.SessionManager
	loginByIP    *limiter
	loginByEmail *limiter
}

func NewHandler(service *Service, sessions *scs.SessionManager) *Handler {
	return &Handler{
		service:      service,
		sessions:     sessions,
		loginByIP:    newLimiter(loginAttemptsPerIP, loginWindow),
		loginByEmail: newLimiter(loginAttemptsPerEmail, loginWindow),
	}
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	if ok, retry := h.loginByIP.allow(clientIP(r)); !ok {
		tooManyAttempts(w, retry)
		return
	}
	var in loginRequest
	if err := httpx.DecodeJSON(w, r, &in); err != nil {
		httpx.BadRequest(w, "invalid JSON body")
		return
	}
	if in.Email == "" || in.Password == "" {
		httpx.BadRequest(w, "email and password are required")
		return
	}
	emailKey := NormalizeEmail(in.Email)
	if ok, retry := h.loginByEmail.allow(emailKey); !ok {
		tooManyAttempts(w, retry)
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
	h.loginByEmail.reset(emailKey)

	// A fresh token on privilege change prevents session fixation.
	if err := h.sessions.RenewToken(ctx); err != nil {
		slog.ErrorContext(ctx, "renew session token", "error", err)
		httpx.InternalError(w)
		return
	}
	h.sessions.Put(ctx, sessionUserIDKey, u.ID)
	httpx.WriteJSON(w, http.StatusOK, u)
}

// clientIP uses a verified proxy address when TrustClientIP supplied one;
// otherwise it uses the immediate connection. Raw forwarding headers are ignored.
func clientIP(r *http.Request) string {
	if address, ok := r.Context().Value(clientIPKey{}).(string); ok {
		return address
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

func tooManyAttempts(w http.ResponseWriter, retry time.Duration) {
	w.Header().Set("Retry-After", strconv.Itoa(int(math.Ceil(retry.Seconds()))))
	httpx.Error(w, http.StatusTooManyRequests, "too many login attempts; try again later")
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

// ListStudents lets staff look up students by the optional q query parameter.
func (h *Handler) ListStudents(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.Students(r.Context(), r.URL.Query().Get("q"))
	if err != nil {
		slog.ErrorContext(r.Context(), "list students", "error", err)
		httpx.InternalError(w)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, users)
}

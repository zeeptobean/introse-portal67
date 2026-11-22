package auth

import (
	"crypto/ed25519"
	"errors"
	"log/slog"
	"net/http"
	"slices"

	"github.com/zeeptobean/introse-portal67/src/backend/internal/httpx"
	"github.com/zeeptobean/introse-portal67/src/backend/internal/platform/authn"
)

func RequireServiceIdentity(keys []ed25519.PublicKey, audience string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			identity, err := authn.Verify(keys, r.Header.Get(authn.IdentityHeader), audience)
			role := Role(identity.Role)
			if err != nil || identity.UserID <= 0 || (role != RoleStaff && role != RoleLecturer && role != RoleStudent) {
				httpx.NotLoggedIn(w)
				return
			}
			user := User{ID: identity.UserID, Role: role, IsActive: true}
			next.ServeHTTP(w, r.WithContext(withUser(r.Context(), user)))
		})
	}
}

// RequireLogin loads the active user for the session and stores it in the
// request context, responding 401 when there is none.
func (h *Handler) RequireLogin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		id := h.sessions.GetInt64(ctx, sessionUserIDKey)
		if id == 0 {
			httpx.NotLoggedIn(w)
			return
		}
		u, err := h.service.UserByID(ctx, id)
		if errors.Is(err, ErrUserNotFound) {
			httpx.NotLoggedIn(w)
			return
		}
		if err != nil {
			slog.ErrorContext(ctx, "load session user", "error", err)
			httpx.InternalError(w)
			return
		}
		if !u.IsActive {
			httpx.NotLoggedIn(w)
			return
		}
		next.ServeHTTP(w, r.WithContext(withUser(ctx, u)))
	})
}

// RequireRole responds 403 unless the user has one of the given roles. It
// must run after RequireLogin.
func RequireRole(roles ...Role) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			u, ok := UserFromContext(r.Context())
			if !ok {
				httpx.NotLoggedIn(w)
				return
			}
			if !slices.Contains(roles, u.Role) {
				httpx.PermissionDenied(w)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

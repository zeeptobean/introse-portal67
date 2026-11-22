package app

import (
	"context"
	"crypto/ed25519"
	"net/http"
	"net/netip"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/zeeptobean/introse-portal67/src/backend/internal/auth"
	"github.com/zeeptobean/introse-portal67/src/backend/internal/httpx"
	"github.com/zeeptobean/introse-portal67/src/backend/internal/platform/authn"
	"github.com/zeeptobean/introse-portal67/src/backend/internal/platform/server"
	"github.com/zeeptobean/introse-portal67/src/backend/internal/routing"
)

func NewIdentityHandler(pool *pgxpool.Pool, trusted []netip.Prefix, keys []ed25519.PublicKey, limit auth.SharedLoginLimiter, checks ...func(context.Context) error) http.Handler {
	sessions := newSessionManager(pool)
	handler := auth.NewHandler(auth.NewService(auth.NewStore(pool)), sessions)
	handler.SharedLoginLimiter = limit
	mux := http.NewServeMux()
	server.HealthRoutes(mux, append([]func(context.Context) error{pool.Ping}, checks...)...)
	serviceRoutes{mux, handler.RequireLogin}.mountIdentity(handler)
	mux.HandleFunc("POST /internal/sessions/resolve", func(w http.ResponseWriter, r *http.Request) {
		identity, err := authn.Verify(keys, r.Header.Get(authn.IdentityHeader), routing.Identity)
		if err != nil || identity.UserID != 0 || identity.Role != "gateway" {
			httpx.PermissionDenied(w)
			return
		}
		handler.Resolve(w, r)
	})
	return auth.TrustClientIP(trusted, sessions.LoadAndSave(requestTrace(mux)))
}

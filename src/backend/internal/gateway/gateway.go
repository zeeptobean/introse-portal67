package gateway

import (
	"context"
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/netip"
	"net/url"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"

	"github.com/zeeptobean/introse-portal67/src/backend/internal/auth"
	"github.com/zeeptobean/introse-portal67/src/backend/internal/httpx"
	"github.com/zeeptobean/introse-portal67/src/backend/internal/platform/authn"
	"github.com/zeeptobean/introse-portal67/src/backend/internal/platform/server"
	"github.com/zeeptobean/introse-portal67/src/backend/internal/routing"
)

type Gateway struct {
	key      ed25519.PrivateKey
	identity *url.URL
	client   *http.Client
	cache    *redis.Client
	prefix   string
}

func New(key ed25519.PrivateKey, endpoints map[string]string, legacy string, trusted []netip.Prefix, cache *redis.Client, prefix string) (http.Handler, error) {
	identity, err := parseOrigin(endpoints[routing.Identity])
	if err != nil {
		return nil, fmt.Errorf("IDENTITY_URL: %w", err)
	}
	if prefix == "" {
		prefix = "p67"
	}
	g := &Gateway{key: key, identity: identity, cache: cache, prefix: prefix, client: &http.Client{
		Timeout:       3 * time.Second,
		CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse },
	}}
	mux := http.NewServeMux()
	upstreams := map[string]*url.URL{}
	for service, routes := range routing.Routes {
		upstream := endpoints[service]
		fallback := upstream == "" && service != routing.Notification && legacy != ""
		if fallback {
			upstream = legacy
		}
		if upstream == "" && service == routing.Notification {
			continue
		}
		target, err := parseOrigin(upstream)
		if err != nil {
			return nil, fmt.Errorf("%s upstream: %w", service, err)
		}
		upstreams[upstream] = target
		proxy := g.proxy(target, service == routing.Identity)
		var handler http.Handler = proxy
		if service != routing.Identity && !fallback {
			handler = g.authenticate(service, proxy)
		}
		for _, route := range routes {
			mux.Handle(route, handler)
		}
	}
	server.HealthRoutes(mux, func(ctx context.Context) error {
		for _, upstream := range upstreams {
			resp, err := g.request(ctx, http.MethodGet, upstream, "/api/ready", "", "")
			if err != nil {
				return err
			}
			_ = resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				return errors.New("upstream unavailable")
			}
		}
		if cache != nil {
			return cache.Ping(ctx).Err()
		}
		return nil
	})
	return auth.TrustClientIP(trusted, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Header.Del(authn.IdentityHeader)
		if len(r.Header.Get("X-Request-ID")) == 0 || len(r.Header.Get("X-Request-ID")) > 128 {
			r.Header.Set("X-Request-ID", uuid.NewString())
		}
		w.Header().Set("X-Request-ID", r.Header.Get("X-Request-ID"))
		mux.ServeHTTP(w, r)
	})), nil
}

func parseOrigin(raw string) (*url.URL, error) {
	target, err := url.Parse(raw)
	if err != nil || target.Host == "" || target.User != nil ||
		(target.Scheme != "http" && target.Scheme != "https") ||
		(target.Path != "" && target.Path != "/") || target.RawQuery != "" || target.Fragment != "" {
		return nil, errors.New("expected an HTTP(S) origin")
	}
	return target, nil
}

func (g *Gateway) proxy(target *url.URL, identity bool) *httputil.ReverseProxy {
	return &httputil.ReverseProxy{
		Rewrite: func(request *httputil.ProxyRequest) {
			request.SetURL(target)
			request.Out.Header.Set("X-Forwarded-For", auth.ClientIP(request.In))
		},
		ModifyResponse: func(response *http.Response) error {
			r := response.Request
			if identity && g.cache != nil && r.URL.Path == "/api/auth/logout" && response.StatusCode == http.StatusOK {
				if cookie, err := r.Cookie("session"); err == nil {
					ctx, cancel := context.WithTimeout(r.Context(), 300*time.Millisecond)
					defer cancel()
					// A tombstone also blocks cache fills racing with the logout.
					_ = g.cache.Set(ctx, g.sessionKey(cookie.Value), "revoked", time.Minute).Err()
				}
			}
			return nil
		},
		ErrorHandler: func(w http.ResponseWriter, _ *http.Request, _ error) {
			httpx.Error(w, http.StatusBadGateway, "backend unavailable")
		},
	}
}

func (g *Gateway) authenticate(service string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session")
		if err != nil {
			httpx.NotLoggedIn(w)
			return
		}
		session, err := g.resolve(r.Context(), cookie)
		if errors.Is(err, errNoSession) {
			httpx.NotLoggedIn(w)
			return
		}
		if err != nil {
			httpx.Error(w, http.StatusBadGateway, "identity unavailable")
			return
		}
		r.Header.Del("Cookie")
		r.Header.Set(authn.IdentityHeader, authn.Sign(g.key, session.Identity, service))
		next.ServeHTTP(w, r)
	})
}

var errNoSession = errors.New("session absent or expired")

func (g *Gateway) resolve(ctx context.Context, cookie *http.Cookie) (authn.Session, error) {
	var session authn.Session
	key := g.sessionKey(cookie.Value)
	if g.cache != nil {
		cacheCtx, cancel := context.WithTimeout(ctx, 300*time.Millisecond)
		value, err := g.cache.Get(cacheCtx, key).Result()
		cancel()
		if err == nil {
			if value == "revoked" {
				return session, errNoSession
			}
			if json.Unmarshal([]byte(value), &session) == nil && time.Now().Before(session.ExpiresAt) {
				return session, nil
			}
		}
	}
	proof := authn.Sign(g.key, authn.Identity{Role: "gateway"}, routing.Identity)
	response, err := g.request(ctx, http.MethodPost, g.identity, "/internal/sessions/resolve", cookie.String(), proof)
	if err != nil {
		return session, err
	}
	defer response.Body.Close()
	if response.StatusCode == http.StatusUnauthorized {
		return session, errNoSession
	}
	if response.StatusCode != http.StatusOK {
		return session, errors.New("identity did not resolve session")
	}
	if err := json.NewDecoder(io.LimitReader(response.Body, 4096)).Decode(&session); err != nil {
		return session, err
	}
	ttl := min(time.Until(session.ExpiresAt), 30*time.Second)
	if ttl <= 0 || session.UserID <= 0 {
		return session, errNoSession
	}
	if g.cache != nil {
		value, _ := json.Marshal(session)
		cacheCtx, cancel := context.WithTimeout(ctx, 300*time.Millisecond)
		defer cancel()
		_ = g.cache.SetNX(cacheCtx, key, value, ttl).Err()
	}
	return session, nil
}

func (g *Gateway) request(ctx context.Context, method string, base *url.URL, path, cookie, proof string) (*http.Response, error) {
	target := base.JoinPath(path)
	r, err := http.NewRequestWithContext(ctx, method, target.String(), nil)
	if err != nil {
		return nil, err
	}
	if cookie != "" {
		r.Header.Set("Cookie", cookie)
		r.Header.Set(authn.IdentityHeader, proof)
	}
	return g.client.Do(r)
}

func (g *Gateway) sessionKey(token string) string {
	return fmt.Sprintf("%s:sess:%x", g.prefix, sha256.Sum256([]byte(token)))
}

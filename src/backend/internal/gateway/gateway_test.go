package gateway

import (
	"crypto/ed25519"
	"crypto/rand"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/zeeptobean/introse-portal67/src/backend/internal/auth"
	"github.com/zeeptobean/introse-portal67/src/backend/internal/platform/authn"
	"github.com/zeeptobean/introse-portal67/src/backend/internal/routing"
)

func TestGatewayResolvesSessionAndSignsForTheRouteOwner(t *testing.T) {
	public, private, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	keys := []ed25519.PublicKey{public}
	identity := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claim, err := authn.Verify(keys, r.Header.Get(authn.IdentityHeader), routing.Identity)
		if err != nil || claim.Role != "gateway" || r.URL.Path != "/internal/sessions/resolve" {
			t.Error("unexpected identity request")
		}
		cookie, _ := r.Cookie("session")
		if cookie == nil || cookie.Value != "valid" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		fmt.Fprintf(w, `{"user_id":42,"role":"student","expires_at":%q}`, time.Now().Add(time.Hour).Format(time.RFC3339))
	}))
	defer identity.Close()
	academic := httptest.NewServer(auth.RequireServiceIdentity(keys, routing.Academic)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u, _ := auth.UserFromContext(r.Context())
		if u.ID != 42 || u.Role != auth.RoleStudent || r.Header.Get("Cookie") != "" {
			t.Error("gateway forwarded the wrong identity or session")
		}
		w.WriteHeader(http.StatusNoContent)
	})))
	defer academic.Close()
	legacy := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get(authn.IdentityHeader) != "" {
			t.Error("legacy received an untrusted service identity header")
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer legacy.Close()
	gateway, err := New(private, map[string]string{routing.Identity: identity.URL, routing.Academic: academic.URL}, legacy.URL, nil, nil, "")
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		path   string
		cookie string
		code   int
	}{
		{"/api/me/results", "valid", http.StatusNoContent},
		{"/api/me/results", "expired", http.StatusUnauthorized},
		{"/api/me/results", "", http.StatusUnauthorized},
		{"/api/sections", "", http.StatusNoContent},
		{"/internal/sessions/resolve", "valid", http.StatusNotFound},
	} {
		r := httptest.NewRequest(http.MethodGet, tc.path, nil)
		r.Header.Set(authn.IdentityHeader, authn.Sign(private, authn.Identity{UserID: 9, Role: "staff"}, routing.Academic))
		if tc.cookie != "" {
			r.AddCookie(&http.Cookie{Name: "session", Value: tc.cookie})
		}
		w := httptest.NewRecorder()
		gateway.ServeHTTP(w, r)
		if w.Code != tc.code {
			t.Errorf("%s with %q: status = %d, want %d", tc.path, tc.cookie, w.Code, tc.code)
		}
	}
}

func TestRejectsUpstreamURLsWithCredentials(t *testing.T) {
	raw := (&url.URL{Scheme: "http", Host: "localhost", User: url.UserPassword("user", "secret")}).String()
	if _, err := parseOrigin(raw); err == nil {
		t.Fatal("upstream credentials accepted")
	}
}

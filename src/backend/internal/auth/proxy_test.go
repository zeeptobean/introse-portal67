package auth

import (
	"net/http"
	"net/http/httptest"
	"net/netip"
	"testing"
)

func TestTrustedClientIP(t *testing.T) {
	for _, tc := range []struct {
		name, peer, forwarded, want string
		trusted                     bool
	}{
		{"default ignores forwarding", "127.0.0.1:54321", "203.0.113.7", "127.0.0.1", false},
		{"trusted immediate proxy", "127.0.0.1:54321", "203.0.113.7", "203.0.113.7", true},
		{"untrusted caller cannot spoof", "198.51.100.8:54321", "203.0.113.7", "198.51.100.8", true},
		{"forwarded chain rejected", "127.0.0.1:54321", "203.0.113.7, 198.51.100.8", "127.0.0.1", true},
		{"malformed address rejected", "127.0.0.1:54321", "not-an-ip", "127.0.0.1", true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var proxies []netip.Prefix
			if tc.trusted {
				proxies = []netip.Prefix{netip.MustParsePrefix("127.0.0.1/32")}
			}
			handler := TrustClientIP(proxies, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if got := clientIP(r); got != tc.want {
					t.Errorf("clientIP = %q, want %q", got, tc.want)
				}
			}))
			request := httptest.NewRequest(http.MethodPost, "/api/auth/login", nil)
			request.RemoteAddr = tc.peer
			request.Header.Set("X-Forwarded-For", tc.forwarded)
			handler.ServeHTTP(httptest.NewRecorder(), request)
		})
	}
}

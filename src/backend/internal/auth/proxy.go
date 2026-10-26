package auth

import (
	"context"
	"net"
	"net/http"
	"net/netip"
)

type clientIPKey struct{}

// TrustClientIP accepts a single client address only from an explicitly trusted
// immediate proxy. SvelteKit overwrites X-Forwarded-For instead of appending it.
// No configured proxies preserves the original direct-connection behavior.
func TrustClientIP(proxies []netip.Prefix, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		host, _, err := net.SplitHostPort(r.RemoteAddr)
		peer, parseErr := netip.ParseAddr(host)
		if err == nil && parseErr == nil {
			for _, prefix := range proxies {
				if !prefix.Contains(peer.Unmap()) {
					continue
				}
				values := r.Header.Values("X-Forwarded-For")
				if len(values) == 1 {
					if address, err := netip.ParseAddr(values[0]); err == nil {
						r = r.WithContext(context.WithValue(r.Context(), clientIPKey{}, address.Unmap().String()))
					}
				}
				break
			}
		}
		next.ServeHTTP(w, r)
	})
}

package middleware

import (
	"context"
	"net"
	"net/http"
	"strings"
)

func TrustedProxy(networks []*net.IPNet) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			client := ClientIP(r)
			peer := net.ParseIP(client)

			if peer != nil && trusted(peer, networks) {
				forwarded := strings.TrimSpace(strings.Split(r.Header.Get("X-Forwarded-For"), ",")[0])
				if net.ParseIP(forwarded) != nil {
					client = forwarded
				}
			}

			ctx := context.WithValue(r.Context(), clientIPContextKey, client)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func ClientIP(r *http.Request) string {
	if ip, ok := r.Context().Value(clientIPContextKey).(string); ok {
		return ip
	}

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil {
		return host
	}

	return r.RemoteAddr
}

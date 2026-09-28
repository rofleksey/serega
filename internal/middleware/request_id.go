package middleware

import (
	"context"
	"net/http"

	"github.com/rofleksey/meg/httpx"
)

func RequestID(next http.Handler) http.Handler {
	return httpx.RequestID(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Downstream protocol adapters read the validated ID from the header too.
		r.Header.Set("X-Request-ID", httpx.RequestIDFrom(r.Context()))
		next.ServeHTTP(w, r)
	}))
}

func RequestIDFrom(ctx context.Context) string {
	return httpx.RequestIDFrom(ctx)
}

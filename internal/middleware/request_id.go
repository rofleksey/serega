package middleware

import (
	"context"
	"net/http"

	"github.com/rofleksey/meg/httpx"
)

func RequestID(next http.Handler) http.Handler {
	return httpx.RequestID(next)
}

func RequestIDFrom(ctx context.Context) string {
	return httpx.RequestIDFrom(ctx)
}

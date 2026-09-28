package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/rofleksey/meg/ratelimit"
)

// RateLimit creates a client-IP keyed rate limiter for one router-owned
// operation or route group.
func RateLimit(limit int, window time.Duration) func(http.Handler) http.Handler {
	limiter := ratelimit.NewLimiter(limit, window)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if allowed, retryAfter := limiter.Allow(ClientIP(r)); !allowed {
				seconds := max(int(retryAfter.Seconds()), 1)
				w.Header().Set("Retry-After", strconv.Itoa(seconds))
				writeError(w, http.StatusTooManyRequests, "rate_limited", "too many requests; try again later", RequestIDFrom(r.Context()))

				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

package middleware

import (
	"net/http"
)

const MaxAPIRequestBytes int64 = 320 << 10

func RequestBodyLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !unsafeMethod(r.Method) {
			next.ServeHTTP(w, r)

			return
		}

		if r.ContentLength > MaxAPIRequestBytes {
			writeError(w, http.StatusRequestEntityTooLarge, "request_too_large", "request body is too large", RequestIDFrom(r.Context()))

			return
		}

		r.Body = http.MaxBytesReader(w, r.Body, MaxAPIRequestBytes)
		next.ServeHTTP(w, r)
	})
}

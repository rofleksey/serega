package middleware

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

// RequestJSON validates JSON syntax and media type when a mutation carries a
// body. Required-body and schema rules remain owned by the generated contract
// and concrete handler, so bodyless operations need no route-name exceptions.
func RequestJSON(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !unsafeMethod(r.Method) || r.Body == nil || r.Body == http.NoBody || r.ContentLength == 0 {
			next.ServeHTTP(w, r)

			return
		}

		if !jsonContentType(r.Header.Get("Content-Type")) {
			writeError(w, http.StatusUnsupportedMediaType, "unsupported_media_type", "expected application/json", RequestIDFrom(r.Context()))

			return
		}

		body, err := io.ReadAll(r.Body)

		var maxErr *http.MaxBytesError
		if errors.As(err, &maxErr) {
			writeError(w, http.StatusRequestEntityTooLarge, "request_too_large", "request body is too large", RequestIDFrom(r.Context()))

			return
		}

		if err != nil || len(bytes.TrimSpace(body)) > 0 && !json.Valid(body) {
			writeError(w, http.StatusBadRequest, "invalid_request", "request body must be valid JSON", RequestIDFrom(r.Context()))

			return
		}

		r.Body = io.NopCloser(bytes.NewReader(body))
		next.ServeHTTP(w, r)
	})
}

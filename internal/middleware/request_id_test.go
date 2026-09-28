package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rofleksey/meg/httpx"
)

func TestRequestIDReplacesInvalidInput(t *testing.T) {
	handler := RequestID(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if RequestIDFrom(r.Context()) == "bad" {
			t.Fatal("invalid request ID reached the handler")
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	request := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/", nil)
	request.Header.Set("X-Request-ID", "bad")

	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)

	if !httpx.ValidRequestID(response.Header().Get("X-Request-ID")) {
		t.Fatalf("generated request ID = %q", response.Header().Get("X-Request-ID"))
	}
}

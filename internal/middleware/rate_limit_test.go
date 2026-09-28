package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"testing/synctest"
	"time"
)

func TestRateLimitThrottlesByClient(t *testing.T) {
	handler := RateLimit(12, time.Minute)(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))

	var response *httptest.ResponseRecorder

	for range 13 {
		request := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/", nil)
		request.RemoteAddr = "192.0.2.1:1234"
		response = httptest.NewRecorder()
		handler.ServeHTTP(response, request)
	}

	if response.Code != http.StatusTooManyRequests || response.Header().Get("Retry-After") == "" {
		t.Fatalf("status = %d, headers = %v", response.Code, response.Header())
	}
}

func TestRateLimitRoundsRetryAfterUp(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		handler := RateLimit(1, 1500*time.Millisecond)(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		}))
		request := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/", nil)

		handler.ServeHTTP(httptest.NewRecorder(), request)

		response := httptest.NewRecorder()
		handler.ServeHTTP(response, request)

		if response.Code != http.StatusTooManyRequests || response.Header().Get("Retry-After") != "2" {
			t.Fatalf("status = %d, headers = %v; want a two-second retry", response.Code, response.Header())
		}

		time.Sleep(2 * time.Second)

		response = httptest.NewRecorder()
		handler.ServeHTTP(response, request)

		if response.Code != http.StatusNoContent {
			t.Fatalf("retry status = %d, want %d", response.Code, http.StatusNoContent)
		}
	})
}

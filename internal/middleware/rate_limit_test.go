package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
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

package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRequestJSONValidatesOnlyPresentBodies(t *testing.T) {
	handler := RequestJSON(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))

	bodyless := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/", nil)
	bodylessResponse := httptest.NewRecorder()
	handler.ServeHTTP(bodylessResponse, bodyless)

	if bodylessResponse.Code != http.StatusNoContent {
		t.Fatalf("bodyless status = %d", bodylessResponse.Code)
	}

	invalid := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/", strings.NewReader("not-json"))
	invalid.Header.Set("Content-Type", "application/json")

	invalidResponse := httptest.NewRecorder()
	handler.ServeHTTP(invalidResponse, invalid)

	if invalidResponse.Code != http.StatusBadRequest {
		t.Fatalf("invalid JSON status = %d", invalidResponse.Code)
	}
}

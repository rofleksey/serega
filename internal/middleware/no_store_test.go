package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNoStoreAppliesToErrorResponses(t *testing.T) {
	handler := NoStore(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/", nil))

	if response.Header().Get("Cache-Control") != "no-store" {
		t.Error("account response may be cached")
	}
}

package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCSRFRejectsUnauthenticatedMutation(t *testing.T) {
	handler := CSRF(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusNoContent) }), true)
	request := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/", nil)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)

	if response.Code != http.StatusForbidden {
		t.Fatalf("status = %d", response.Code)
	}
}

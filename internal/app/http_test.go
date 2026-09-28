package app

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

type unavailableDatabase struct{}

func (unavailableDatabase) Ping(context.Context) error { return errors.New("unavailable") }
func (unavailableDatabase) Close()                     {}

func TestHealthAndReadinessDistinguishDatabaseFailure(t *testing.T) {
	handler := NewHTTPHandler(HTTPDependencies{Database: unavailableDatabase{}})

	for path, status := range map[string]int{"/healthz": http.StatusNoContent, "/readyz": http.StatusServiceUnavailable, "/api/v1/missing": http.StatusNotFound} {
		response := httptest.NewRecorder()
		handler.ServeHTTP(response, httptest.NewRequestWithContext(t.Context(), http.MethodGet, path, nil))

		if response.Code != status {
			t.Errorf("%s status = %d", path, response.Code)
		}
	}
}

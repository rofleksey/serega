package frontend

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandlerServesEmbeddedIndexForSPARoutes(t *testing.T) {
	handler := NewHandler()

	for _, route := range []string{"/", "/settings"} {
		response := httptest.NewRecorder()
		request := httptest.NewRequestWithContext(context.Background(), http.MethodGet, route, nil)
		handler.ServeHTTP(response, request)

		if response.Code != http.StatusOK || !strings.Contains(response.Body.String(), `<div id="root"></div>`) {
			t.Fatalf("GET %s = (%d, %q)", route, response.Code, response.Body.String())
		}
	}
}

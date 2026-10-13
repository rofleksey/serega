package frontend

import (
	"context"
	"net/http"
	"net/http/httptest"
	"regexp"
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

		if cache := response.Header().Get("Cache-Control"); cache != "no-cache" {
			t.Fatalf("GET %s cache = %q, want no-cache", route, cache)
		}
	}
}

func TestHandlerServesAssetsImmutablyAndRejectsMissingFiles(t *testing.T) {
	handler := NewHandler()

	index := httptest.NewRecorder()

	handler.ServeHTTP(index, httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/", nil))

	match := regexp.MustCompile(`src="([^"]+\.js)"`).FindStringSubmatch(index.Body.String())

	if match == nil {
		t.Fatal("index does not reference a built asset")
	}

	asset := httptest.NewRecorder()

	handler.ServeHTTP(asset, httptest.NewRequestWithContext(context.Background(), http.MethodGet, match[1], nil))

	if asset.Code != http.StatusOK || !strings.Contains(asset.Header().Get("Cache-Control"), "immutable") {
		t.Fatalf("GET %s = (%d, %q)", match[1], asset.Code, asset.Header().Get("Cache-Control"))
	}

	missing := httptest.NewRecorder()

	handler.ServeHTTP(missing, httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/assets/missing-chunk.js", nil))

	if missing.Code != http.StatusNotFound {
		t.Fatalf("missing asset status = %d, want %d", missing.Code, http.StatusNotFound)
	}
}

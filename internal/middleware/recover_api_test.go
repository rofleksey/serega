package middleware

import (
	"bytes"
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRecoverAPIUsesStableAPIEnvelope(t *testing.T) {
	var logs bytes.Buffer

	logger := slog.New(slog.NewJSONHandler(&logs, nil))
	handler := RecoverAPI(logger)(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		panic("boom")
	}))
	request := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/api/v1/servers", nil)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)

	if response.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d", response.Code)
	}

	if strings.Contains(response.Body.String(), "boom") {
		t.Fatalf("response leaked panic details: %s", response.Body.String())
	}

	if output := logs.String(); !strings.Contains(output, `"error_code":"handler_panic"`) || !strings.Contains(output, `"panic":"boom"`) || !strings.Contains(output, `"stack":`) {
		t.Fatalf("panic log = %s", output)
	}
}

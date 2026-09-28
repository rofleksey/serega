package middleware

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"

	"github.com/rofleksey/serega/internal/observability"
)

func TestRequestLoggerRecordsResponseMetadata(t *testing.T) {
	var output bytes.Buffer

	handler := RequestLogger(observability.NewEventRuntime(observability.NewLogger("json", &output)), nil)(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusAccepted)
	}))
	request := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/test", nil)
	handler.ServeHTTP(httptest.NewRecorder(), request)

	if !strings.Contains(output.String(), `"http.status":202`) || !strings.Contains(output.String(), `"http.route":"unmatched"`) {
		t.Fatalf("log = %s", output.String())
	}
}

func TestStatusWriterSupportsStreaming(t *testing.T) {
	recorder := httptest.NewRecorder()
	writer := &statusWriter{ResponseWriter: recorder}

	writer.Flush()

	if !recorder.Flushed || writer.status != http.StatusOK || writer.Unwrap() != recorder {
		t.Fatalf("streaming writer = %#v, flushed = %t", writer, recorder.Flushed)
	}
}

func TestStatusWriterSupportsHijacking(t *testing.T) {
	underlying := newHijackingWriter(t)
	writer := &statusWriter{ResponseWriter: underlying}

	connection, _, err := writer.Hijack()
	if err != nil {
		t.Fatal(err)
	}

	_ = connection.Close()

	underlying.closePeer()
}

func TestRequestLoggerOmitsRawPaths(t *testing.T) {
	for _, test := range []struct {
		name  string
		path  string
		route string
	}{
		{name: "matched", path: "/items/private-resource?token=private-query", route: "/items/{itemID}"},
		{name: "unmatched", path: "/unknown/private-resource?token=private-query", route: "unmatched"},
	} {
		t.Run(test.name, func(t *testing.T) {
			var output bytes.Buffer

			router := chi.NewRouter()
			router.Use(RequestLogger(observability.NewEventRuntime(observability.NewLogger("json", &output)), nil))
			router.Get("/items/{itemID}", func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusNoContent)
			})

			request := httptest.NewRequestWithContext(t.Context(), http.MethodGet, test.path, nil)
			router.ServeHTTP(httptest.NewRecorder(), request)

			event := output.String()
			if !strings.Contains(event, `"http.route":"`+test.route+`"`) {
				t.Fatalf("event omitted normalized route %q: %s", test.route, event)
			}

			if strings.Contains(event, "private-resource") || strings.Contains(event, "private-query") {
				t.Fatalf("event contains private URL values: %s", event)
			}
		})
	}
}

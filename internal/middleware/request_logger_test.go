package middleware

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

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

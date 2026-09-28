package observability

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"testing"

	"github.com/happytoolin/unolog"
	"github.com/rofleksey/serega/internal/entity"
)

func TestLoggerRedactsSensitiveFields(t *testing.T) {
	var output bytes.Buffer

	logger := NewLogger("json", &output)
	logger.InfoContext(context.Background(), "safe event",
		"token", "do-not-log",
		"nested", slog.GroupValue(slog.String("password", "also-secret"), slog.String("result", "ok")),
	)

	var record map[string]any
	if err := json.Unmarshal(output.Bytes(), &record); err != nil {
		t.Fatal(err)
	}

	if record["token"] != "[REDACTED]" {
		t.Fatalf("token = %v", record["token"])
	}

	nested, ok := record["nested"].(map[string]any)
	if !ok || nested["password"] != "[REDACTED]" || nested["result"] != "ok" {
		t.Fatalf("nested fields = %#v", record["nested"])
	}
}

func TestWideEventHasStableCompletionSchema(t *testing.T) {
	sink := unolog.NewTestSink()
	runtime := unolog.MustCompile(unolog.Config{Sink: sink, SamplingRate: 1})
	event := StartOperation(context.Background(), runtime, unolog.DomainHTTP, entity.EventHTTPRequest, "",
		entity.FieldOperation, entity.OperationAuthenticate,
	)
	CompleteOperation(event, entity.OutcomeSucceeded, nil)

	events := sink.Events()
	if len(events) != 1 {
		t.Fatalf("events = %d", len(events))
	}

	name, _ := events[0].Lookup(entity.FieldEvent)
	outcome, _ := events[0].Lookup(unolog.KeyOpOutcome)

	duration, found := events[0].Lookup(unolog.KeyDurationMS)
	if name != entity.EventHTTPRequest || outcome != string(unolog.OutcomeSuccess) || !found || duration == nil {
		t.Fatalf("wide event = %#v", events[0].Fields())
	}
}

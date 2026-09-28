package observability

import (
	"context"
	"testing"
	"time"

	"github.com/rofleksey/serega/internal/entity"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

func TestMetricsRecordHTTPAndDatabaseDimensions(t *testing.T) {
	reader := sdkmetric.NewManualReader()
	provider := sdkmetric.NewMeterProvider(sdkmetric.WithReader(reader))

	t.Cleanup(func() { _ = provider.Shutdown(context.Background()) })

	metrics, err := NewMetrics(provider, "test")
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	metrics.HTTPStarted(ctx, "GET")
	metrics.HTTPCompleted(ctx, "GET", "/api/v1/cards/{cardId}", 200, 25*time.Millisecond)

	if err := metrics.RegisterDatabasePool(func() DatabasePoolSnapshot {
		return DatabasePoolSnapshot{Acquired: 2, Idle: 3, Maximum: 10}
	}); err != nil {
		t.Fatal(err)
	}

	var collected metricdata.ResourceMetrics
	if err := reader.Collect(ctx, &collected); err != nil {
		t.Fatal(err)
	}

	names := make(map[string]bool)

	for _, scope := range collected.ScopeMetrics {
		for _, item := range scope.Metrics {
			names[item.Name] = true
		}
	}

	for _, name := range []string{entity.MetricHTTPRequests, entity.MetricHTTPDuration, entity.MetricHTTPActive, entity.MetricDBConnections, entity.MetricDBMaxConnections} {
		if !names[name] {
			t.Fatalf("metric %q was not collected: %#v", name, names)
		}
	}
}

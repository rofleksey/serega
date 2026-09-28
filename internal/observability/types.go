package observability

import (
	"github.com/rofleksey/meg/otelmetrics"
	"go.opentelemetry.io/otel/metric"
)

type Provider struct {
	Metrics *Metrics
	runtime *otelmetrics.Provider
}

type Metrics struct {
	meter            metric.Meter
	httpRequests     metric.Int64Counter
	httpDuration     metric.Float64Histogram
	httpActive       metric.Int64UpDownCounter
	dbConnections    metric.Int64ObservableGauge
	dbMaxConnections metric.Int64ObservableGauge
}

type DatabasePoolSnapshot struct {
	Acquired int64
	Idle     int64
	Maximum  int64
}

type DatabasePoolSnapshotFunc func() DatabasePoolSnapshot

package observability

import (
	"context"
	"time"

	"github.com/rofleksey/meg/otelmetrics"
	"github.com/rofleksey/serega/internal/entity"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

func Open(ctx context.Context, serviceName string) (*Provider, error) {
	runtime, err := otelmetrics.Open(ctx, serviceName)
	if err != nil {
		return nil, err
	}

	metrics, err := NewMetrics(runtime.MeterProvider(), serviceName)
	if err != nil {
		_ = runtime.Shutdown(ctx)
		return nil, err
	}

	return &Provider{Metrics: metrics, runtime: runtime}, nil
}

func Disabled(serviceName string) *Provider {
	runtime := otelmetrics.Disabled()

	metrics, err := NewMetrics(runtime.MeterProvider(), serviceName)
	if err != nil {
		panic(err)
	}

	return &Provider{Metrics: metrics, runtime: runtime}
}

func NewMetrics(provider metric.MeterProvider, instrumentationName string) (*Metrics, error) {
	meter := provider.Meter(instrumentationName)

	httpRequests, err := meter.Int64Counter(entity.MetricHTTPRequests, metric.WithUnit("{request}"))
	if err != nil {
		return nil, err
	}

	httpDuration, err := meter.Float64Histogram(entity.MetricHTTPDuration, metric.WithUnit("ms"))
	if err != nil {
		return nil, err
	}

	httpActive, err := meter.Int64UpDownCounter(entity.MetricHTTPActive, metric.WithUnit("{request}"))
	if err != nil {
		return nil, err
	}

	return &Metrics{
		meter:        meter,
		httpRequests: httpRequests,
		httpDuration: httpDuration,
		httpActive:   httpActive,
	}, nil
}

func (metrics *Metrics) RegisterDatabasePool(snapshot DatabasePoolSnapshotFunc) error {
	if metrics == nil || snapshot == nil {
		return nil
	}

	connections, err := metrics.meter.Int64ObservableGauge(entity.MetricDBConnections,
		metric.WithUnit("{connection}"),
		metric.WithInt64Callback(func(_ context.Context, observer metric.Int64Observer) error {
			current := snapshot()
			observer.Observe(current.Acquired, metric.WithAttributes(attribute.String(entity.MetricAttrState, entity.MetricStateAcquired)))
			observer.Observe(current.Idle, metric.WithAttributes(attribute.String(entity.MetricAttrState, entity.MetricStateIdle)))

			return nil
		}),
	)
	if err != nil {
		return err
	}

	maximum, err := metrics.meter.Int64ObservableGauge(entity.MetricDBMaxConnections,
		metric.WithUnit("{connection}"),
		metric.WithInt64Callback(func(_ context.Context, observer metric.Int64Observer) error {
			observer.Observe(snapshot().Maximum)

			return nil
		}),
	)
	if err != nil {
		return err
	}

	metrics.dbConnections = connections
	metrics.dbMaxConnections = maximum

	return nil
}

func (provider *Provider) Shutdown(ctx context.Context) error {
	if provider == nil || provider.runtime == nil {
		return nil
	}

	return provider.runtime.Shutdown(ctx)
}

func (metrics *Metrics) HTTPStarted(ctx context.Context, method string) {
	if metrics == nil {
		return
	}

	metrics.httpActive.Add(ctx, 1, metric.WithAttributes(attribute.String(entity.MetricAttrMethod, boundedHTTPMethod(method))))
}

func (metrics *Metrics) HTTPCompleted(ctx context.Context, method, route string, status int, duration time.Duration) {
	if metrics == nil {
		return
	}

	method = boundedHTTPMethod(method)
	active := metric.WithAttributes(attribute.String(entity.MetricAttrMethod, method))
	metrics.httpActive.Add(ctx, -1, active)

	attributes := metric.WithAttributes(
		attribute.String(entity.MetricAttrMethod, method),
		attribute.String(entity.MetricAttrRoute, route),
		attribute.Int(entity.MetricAttrStatusCode, status),
	)
	metrics.httpRequests.Add(ctx, 1, attributes)
	metrics.httpDuration.Record(ctx, durationMilliseconds(duration), attributes)
}

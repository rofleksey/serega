# Serega observability collector

Serega emits JSON wide-event logs to standard error and OTLP/HTTP metrics to a
collector. It does not emit distributed traces.

Run a pinned `otel/opentelemetry-collector-contrib:0.161.0` image with
`collector.yaml`. Supply `SEREGA_OTLP_EXPORT_ENDPOINT` and
`SEREGA_OTLP_EXPORT_AUTHORIZATION` through the deployment secret mechanism. Do
not commit either value. The sample pipeline has an explicit memory limit,
bounded sending queue, bounded retry period, batching, and a health endpoint.
Configure retention on the selected metrics backend; the Collector forwards
telemetry and is not a durable store.

Configure Serega with:

```text
OTEL_METRICS_EXPORTER=otlp
OTEL_EXPORTER_OTLP_METRICS_ENDPOINT=http://collector:4318/v1/metrics
OTEL_EXPORTER_OTLP_METRICS_PROTOCOL=http/protobuf
```

Use `OTEL_METRICS_EXPORTER=none` for local runs. Export failures do not block
board operations. Shutdown attempts a bounded final flush. Collect JSON process
logs with your platform's log shipper and index `event.name`, `outcome`, and
`request_id`. Never record passwords, cookies, authorization headers, or card
bodies in telemetry. Request metrics use bounded methods and route patterns.

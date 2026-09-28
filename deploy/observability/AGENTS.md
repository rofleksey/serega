# Optional standard telemetry pipeline

- Use OpenTelemetry Collector/OTLP; do not invent exporters or metric transports.
- Keep batching, bounded queues/retries/timeouts, health checks, and retention
  explicit for enabled components.
- Export failure must not break normal application work.
- Keep endpoints/credentials in deployment configuration; never commit secrets.
- Match service/metric names to the current Serega vocabulary.
- Do not expose raw card content, authentication material, or unbounded IDs as
  metric attributes.
- Document the optional setup independently from the minimal app/database path.

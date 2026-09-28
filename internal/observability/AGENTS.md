# Observability bridges

`logger.go` bridges `meg/logging` and `meg/operation` to redacting slog and
unolog. `metrics.go` uses OpenTelemetry through `meg/otelmetrics`. Keep generic
redaction, event lifecycle, and export mechanics in those libraries.

- Start one rich completion event at the owning boundary, enrich through
  context, and finish once from one goroutine. Independent concurrent work owns
  separate operations; do not add a parallel accumulator or lifecycle registry.
- Expected decisions enrich the event. Standalone logs describe process
  lifecycle or exceptional diagnostics. Handler failures attach a safe code and
  error type rather than raw dependency errors containing user input.
- Card titles/descriptions and unrestricted payloads are not telemetry fields.
- Metrics cover HTTP and database-pool behavior. Labels use bounded method,
  route pattern, status, and outcome; never usernames, card IDs, raw paths/errors.
- Balance active/completed accounting. Initialization/export fails open;
  queues, timeouts, and shutdown flushing remain bounded.
- Use library in-memory sinks for precise assertions. Public integration tests
  verify safe logs, correlation, OTLP delivery/failure, and shutdown behavior.
  The optional Collector example is in `deploy/observability`.

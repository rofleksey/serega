# Library bridges and telemetry

Keep this a narrow bridge from Serega's domain vocabulary to meg, unolog, slog,
and OpenTelemetry, following Akio's infrastructure boundary.

- Application logging uses the redacting slog logger.
- unolog via the existing operation bridge owns event lifecycle. Do not create
  a parallel field accumulator, exporter, or lifecycle registry.
- Start at the operation boundary, propagate context, safely enrich, and finish
  once from one goroutine. Independent work owns separate operations.
- Expected decisions enrich a completion event; standalone logs are for process
  lifecycle or exceptional diagnostics.
- Never record card titles/descriptions, passwords/hashes, cookies/tokens,
  headers, connection strings, or unrestricted payloads.
- Metric vocabulary lives in entity. Labels must be bounded: route patterns,
  method, status, outcome; never user/card IDs or raw errors.
- Keep active/completed accounting balanced and initialization/export fail-open.
  Shutdown flush, queues, and timeouts must be bounded.
- Use library in-memory test sinks for precise assertions. External tests verify
  safe emitted logs, correlation, OTLP delivery/failure, and shutdown behavior.

# Deployment

`production/README.md` and `observability/README.md` contain operator commands.
Keep them aligned with the example configuration.

## Production

- The Dockerfile builds Vite with pinned Node and generates bindings with pinned
  Go before building the embedded executable. Keep toolchain versions, lockfile,
  binary path, image command, and runtime user aligned.
- The minimal runtime runs as a non-root user. Preserve the read-only app
  filesystem, dropped capabilities, no-new-privileges, and bounded shutdown.
- Production Compose runs only the app and requires an external `DATABASE_URL`.
  Operators own PostgreSQL availability, durable storage, backups, migration,
  and user provisioning. Ensure database readiness and successful explicit
  migration before serving; do not add automatic serve-time migrations.
- HTTP binds to loopback for an HTTPS reverse proxy. Keep Secure cookies enabled
  and configure trusted proxy CIDRs for the actual proxy networks.
- Keep image publication and deployment separate from local verification.

## Optional telemetry

- `observability/collector.yaml` is a standard OTLP/HTTP metrics pipeline. Logs
  go to stderr for a platform log shipper; the application emits no traces.
- Keep Collector memory limits, batching, sending queues, retries/timeouts,
  and health checks explicit and bounded. Match service/metric names to the app.
- The Collector forwards metrics; retention belongs to the selected backend.
  Keep the app/database deployment usable without this optional pipeline.

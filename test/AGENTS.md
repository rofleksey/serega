# Integration tests

Run `make test-integration` to build the current executable and run both suites.
Keep the `integration` build tag and Makefile targets aligned. Testcontainers
owns dependency lifecycle by default; pin container images and bound startup
and teardown. Use deterministic assertions with deadlines instead of arbitrary
sleeps; keep load tests separate and fixtures minimal and synthetic.

## Application boundary

- `integration/application` interacts only through CLI, HTTP, process, logs, and
  OTLP boundaries; do not import internal packages or instantiate services.
- Its harness starts disposable PostgreSQL, runs explicit migration, starts the
  executable, captures logs, and supplies a fake OTLP receiver. Provision users
  through the CLI. `SEREGA_BINARY` overrides the default `bin/serega` artifact;
  ensure any override was built from the current source and frontend.
- Use independent cookie jars and CSRF flows to prove shared board visibility
  with separate authentication state. Cover CRUD/status movement, attribution,
  validation, missing cards, stale edits/deletes, and session behavior.
- Verify health/readiness, embedded UI and API routing on one origin, persistent
  data, request-event correlation, exporter failure tolerance, and signal-driven
  shutdown with a final telemetry flush.
- Guarantee process, receiver, and container cleanup after partial startup or
  failed assertions. Bound readiness polling and shutdown waits.

## PostgreSQL boundary

- `integration/postgres` applies real embedded Goose migrations and uses the
  production adapter/SQLC queries. Test uniqueness, constraints, ordering,
  attribution, persistence, transactions, and read-only schema readiness.
- Exercise concurrent expected-version mutations and distinguish stale from
  missing records. Fake-store tests cannot establish database atomicity.
- Verify account creation preserves other users and session changes are scoped.
  Schema changes require fresh migration and reapplication/upgrade coverage.
- The harness creates and drops a separate database per test. Its optional
  `SEREGA_INTEGRATION_DATABASE_URL` must point to a disposable PostgreSQL instance
  with database-creation privileges; it is not a development-database shortcut.

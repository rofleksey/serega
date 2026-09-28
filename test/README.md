# Integration tests

Unit tests remain beside their packages. This directory tests real boundaries:

- `integration/application` launches the built Serega executable with disposable
  PostgreSQL and uses only public HTTP, CLI, and process boundaries. It verifies
  separate user sessions, shared card CRUD, stale-write protection, CSRF,
  authentication, structured request events, OTLP failure tolerance, and shutdown.
- `integration/postgres` tests embedded migrations, read-only readiness, account
  uniqueness, persistence, and concurrent version checks against PostgreSQL.

Run `make test-integration` with Docker available. Testcontainers owns the pinned
PostgreSQL lifecycle. These suites never use a development or production database.
`SEREGA_BINARY` can override the built artifact for the application suite.

# Composition and lifecycle

This is where concrete PostgreSQL adapters, account/board services, handlers,
configuration, and telemetry are connected.

- Inject the store into services and services into handlers. Do not pass stores
  into transport or CLI code.
- Open and ping dependencies before accepting traffic; clean up partially
  opened resources on startup failure.
- Serving never migrates. The migration command owns explicit schema changes.
- Compose request ID, proxy trust, completion events, recovery, limits,
  authentication, sessions, and CSRF in intentional router-owned order.
- Keep dependency readiness separate from liveness and return safe diagnostics.
- UI and API share one origin. Unknown API paths return API errors, not SPA HTML.
- Bound server timeouts and shutdown: drain HTTP, stop owned work, close the pool,
  and flush telemetry with a deadline.
- Telemetry errors must not block application work.
- Lifecycle interfaces are appropriate for focused tests; global registries and
  speculative fake production adapters are not.
- Test routing/lifecycle here and built-process behavior in the application
  integration suite.

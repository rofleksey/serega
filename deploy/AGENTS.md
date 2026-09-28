# Deployment examples

Keep deployment simple: one immutable Go/frontend application plus PostgreSQL,
with optional standard telemetry infrastructure.

- Production examples must not require Vite, mutable asset mounts, source mounts,
  default credentials, or automatic serve-time migrations.
- Keep secrets in untracked deployment environment/secret storage.
- Build versions and names must agree with module/runtime/docs.
- Provide explicit migration order, readiness, bounded shutdown, least-privilege
  containers, persistent database volume, and clear TLS expectations.
- Examples are not permission to deploy. Deployment requires user authorization.
- Update operational docs with changes; validate configuration without contacting
  production services.

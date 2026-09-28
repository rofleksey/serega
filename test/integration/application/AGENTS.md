# Built-process application tests

Verify the artifact users run, using public boundaries only.

- Do not import internal application packages or instantiate handlers/services.
- Start disposable PostgreSQL with Testcontainers, run explicit migration/user
  commands through the executable, then launch serve.
- Use independent browser cookie jars/CSRF flows to prove multi-user access and
  isolation of authentication state.
- Cover shared visibility, CRUD/status movement, stale edit/delete conflicts,
  invalid input, missing targets, unauthenticated requests, and session behavior.
- Verify readiness, same-origin UI/API behavior, process signals/shutdown, safe
  wide-event correlation, and telemetry where relevant.
- Capture process output safely and guarantee cleanup on partial startup failure.
- Keep waits bounded and teardown reliable even after failed assertions.
- Build current frontend/backend before running; stale bin output is not evidence.

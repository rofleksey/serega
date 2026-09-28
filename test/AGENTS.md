# Integration-only test tree

Focused unit tests stay beside Go/frontend packages. This tree verifies real
application and PostgreSQL boundaries.

- application runs the built executable and interacts through public HTTP,
  process, CLI, logs, and OTLP boundaries. It does not import internal packages.
- postgres tests schema versions, transactions, and adapter guarantees against
  a real isolated database.
- Use Testcontainers with pinned images, deterministic setup, bounded readiness/
  polling, and cleanup. Never point tests at production or shared persistent state.
- Keep build tags/Makefile targets aligned. A missing Docker daemon is a reported
  environmental blocker, never a passing test.
- Prefer deterministic assertions and eventually-with-deadline over arbitrary sleeps.
- Keep fixtures minimal and synthetic; never copy credentials or user data.

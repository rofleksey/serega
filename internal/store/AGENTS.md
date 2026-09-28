# Outbound persistence

Concrete adapters satisfy ports owned by consuming use cases.

- Transport and CLI must not depend on this subtree.
- Map generated rows to domain values here; do not leak SQLC into application
  services or HTTP contracts.
- Own transaction mechanics and defensive consistency, while services own policy.
- Prefer domain/behavior files over a generic repository framework.
- Follow postgres/AGENTS.md for the current implementation.

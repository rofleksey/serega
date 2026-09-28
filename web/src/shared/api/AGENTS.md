# Browser transport boundary

api/openapi.yaml is canonical. generated.ts is ignored generator output;
change the contract and regenerate rather than editing generated definitions.

- Keep generated openapi-fetch calls centralized in client/transport modules.
- Use same-origin credentials and the shared CSRF header helper for mutations.
- CSRF tokens stay in memory; never put cookies/passwords/session tokens in
  browser persistent storage.
- Keep ApiRequestError and field-error extraction as the shared error model.
  Preserve request IDs for useful diagnostics without exposing server internals.
- Notify session expiry centrally and consistently with server status semantics.
- Domain aliases derive from generated schemas, not duplicate interfaces.
- Treat no-content success correctly and do not parse empty bodies as JSON.
- Never silently retry a version conflict as an unconditional overwrite.
- Unit-test request method/path/header/body construction, envelope failures,
  and session-expiry behavior when changing this boundary.

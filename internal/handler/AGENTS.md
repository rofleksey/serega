# HTTP handlers

One endpoint implementation per descriptive source file. Shared dependency
contracts belong in types.go; routing and response mechanics have clear owners.

1. Apply the configured authentication/CSRF/body validation boundaries.
2. Decode generated request types and validated route parameters.
3. Resolve the authenticated user as required.
4. Invoke an account/board service with the request context.
5. Map known domain outcomes into the common error envelope.
6. Explicitly convert domain values into generated responses.

- Never call a store, SQLC query, pgx pool, or persistence interface directly.
- Keep validation policy and authorization in services; transport parses and maps.
- Preserve safe stable error codes, request IDs, and useful field errors.
- Missing cards and stale versions must remain distinguishable.
- All authenticated users share cards; identity provides attribution, not a
  per-user visibility filter.
- Preserve login token renewal, logout destruction, and current-session lookup.
- Session expiry handling must agree with the frontend and remain distinct from
  a rejected login attempt.
- Compose middleware through router ownership, not raw path-segment inspection.
- Keep handwritten behavior outside generated packages.
- Focus tests on HTTP mapping/auth/validation; public integration tests prove
  the complete path through services and persistence.

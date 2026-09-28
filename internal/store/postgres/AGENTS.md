# PostgreSQL adapter

Preserve Akio's behavior-owned account/card files, explicit conversions, and
small shared construction/types files.

- Static SQL lives in sqlc/query. Call generated bindings and map rows explicitly.
- Keep generated parameter/row types inside this boundary.
- Bind generated queries to transactions; roll back on early exits and commit
  only after the full operation succeeds.
- Compare expected version and increment it in the SQL mutation itself.
  Read-then-unconditional-write is not concurrency protection.
- Distinguish absent card from stale existing version without exposing pgx errors.
- Actor IDs come from authenticated service input. Do not add per-user filters
  to the shared board.
- User creation never deletes others. Password/session changes affect only the
  intended user; never restore sole-administrator replacement.
- Constraints defend service rules; the adapter must not be their only owner.
- Parameterize runtime values. Handwritten dynamic SQL requires documented
  necessity and constrained identifiers.
- Private mapping helpers stay beside behavior; helper.go is genuinely shared.
- Real transaction/schema/concurrency guarantees belong in PostgreSQL integration
  tests rather than mocks that merely repeat pgx method calls.

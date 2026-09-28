# PostgreSQL adapter and queries

Keep behavior-specific account/card files with explicit row conversions.
`store.go` binds queries to the pool or the transaction carried by its private
context key. Nested `WithinTransaction` calls reuse that transaction; all
participating operations must receive the callback's context. Roll back early
exits and commit only after the complete operation succeeds.

## Card concurrency

- Lock the card row while deciding missing versus stale, then perform the
  version-predicated mutation in the same transaction. An unconditional write
  after a version read does not protect against concurrent changes.
- Updates atomically increment version and persist actor/time; deletion also
  compares expected version. Return distinct domain missing/conflict errors.
- Read persisted values in the transaction or use `RETURNING`; do not guess
  timestamps or versions in Go. Attribution comes from service input.
- Verify transactions, schema constraints, and concurrency against real
  PostgreSQL rather than mocks that merely repeat driver method calls.

## SQLC sources

- Use SQLC for fixed-shape queries in `sqlc/query/auth.sql` and `cards.sql`.
  Dynamic SQL needs documented necessity, constrained identifiers, and
  parameterized runtime values.
- Choose accurate query cardinality/result shapes. `SELECT *` means schema
  changes also change the generated row contract. Order lists deterministically.
- Account lookup is identity-scoped; board listing is shared across users.
- `sqlc/sqlc.yaml` reads the Goose migrations and uses pgx/v5 with explicit UUID,
  nullable, and time mappings. Preserve NULL/empty distinctions during mapping;
  adjust adapter conversions when these contracts change.
- Index changes belong with migrations. Regenerate and run
  `go tool sqlc vet -f internal/store/postgres/sqlc/sqlc.yaml` from the repo root
  after query/config changes.

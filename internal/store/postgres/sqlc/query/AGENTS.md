# Static SQL query sources

Group named SQLC queries by account/board behavior.

- Choose the correct cardinality annotation and result shape. With SELECT *,
  remember that a schema change also changes the generated row contract.
- Parameterize all runtime values and order results deterministically where
  callers require stable ordering.
- Account lookup is identity-scoped; board listing is shared across users.
- Card updates/deletes compare expected version in the mutation predicate.
  Updates increment version and persist actor/time atomically.
- Prefer RETURNING persisted values to guessing timestamps/versions in Go.
- Distinguish NULL/empty deliberately and follow sqlc.yaml mappings.
- Indexes belong in migrations. Changed consistency guarantees require real
  PostgreSQL integration coverage.
- Generated bindings are disposable output, never the source of truth.

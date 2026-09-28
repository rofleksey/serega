# PostgreSQL lifecycle and migrations

- `migrate_pgx.go` uses an instance-based Goose provider with embedded migration
  files. Keep provider state local; do not mutate a global migration filesystem
  or introduce another schema versioning system.
- Serving checks schema readiness read-only. `readiness.go` deliberately queries
  Goose metadata directly: Goose version helpers can initialize tables. This
  infrastructure exception does not permit handwritten application queries.
- Close partially opened connections and honor cancellation.

## Schema changes

- Follow the numbered filenames and Goose Up/Down/transaction conventions in
  `migrations/`. After deployment, append migrations; a fresh-project reset must
  deliberately account for any existing data and deployments.
- Defend unique identities, allowed statuses, required attribution, positive
  card versions, and referential integrity with database constraints.
- Keep the sessions schema compatible with SCS storage and the card schema
  suitable for atomic expected-version updates/deletes.
- Add indexes for real query needs. Consider upgrade and rollback behavior;
  destructive down migrations are not ordinary deployment steps.
- Do not seed demo credentials/cards or environment-specific values.
- Verify migration application/reapplication and schema/store guarantees in the
  isolated PostgreSQL integration suite, including SQLC generation/vet after
  schema changes.

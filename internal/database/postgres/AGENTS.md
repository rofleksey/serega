# PostgreSQL lifecycle

Goose owns migrations. SQLC reads the migration directory to infer query types
but does not execute or own schema changes.

- Keep migration embedding and explicit migrate-up entry points here.
- Serving must not migrate as a side effect.
- Reuse maintained Goose/pgx lifecycle; do not add another version table/system.
- Close partially opened pools and honor cancellation.
- Never log the database URL.
- Verify migration application/reapplication and schema guarantees in the
  isolated PostgreSQL integration suite.
- Read migrations/AGENTS.md before changing schema history.

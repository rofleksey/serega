# SQLC generation boundary

This directory contains handwritten generator source/configuration and queries.
The separate generated/ child contains ignored generator-owned output only.

- Edit sqlc.yaml, generate.go, and query/*.sql as needed.
- Never edit generated/*.go or add handwritten helpers/docs to generated/.
- Schema input is database/postgres/migrations; do not add a second schema.
- Preserve pgx/v5 and explicit UUID/nullable/time mappings unless an intentional
  contract change requires corresponding adapter updates.
- Generation must work from a clean checkout.
- Regenerate and run go tool sqlc vet -f internal/store/postgres/sqlc/sqlc.yaml
  after changes.
- Read query/AGENTS.md for query correctness.

# Goose migration sources

This directory is the canonical ordered PostgreSQL schema and SQLC schema input.

- Follow existing numbered filenames and Goose Up/Down/transaction conventions.
- After deployment, append migrations; never rewrite applied history. A fresh
  project reset must deliberately account for existing data/deployments.
- Defend invariants with unique identity, valid statuses, required attribution,
  positive versions, and referential integrity.
- Database constraints complement use-case validation and helpful errors.
- Session schema must remain compatible with persistent SCS storage.
- Card schema must support atomic expected-version mutations.
- Do not seed demo credentials/cards or environment-specific values.
- Consider upgrades, real query indexes, and rollback safety; destructive down
  operations are never ordinary deployment.
- Run SQLC generation/vet and isolated schema/store tests after changes.

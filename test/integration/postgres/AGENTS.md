# Real PostgreSQL contracts

Use a disposable isolated database to prove schema/store guarantees.

- Apply real Goose migrations and use the production adapter/SQLC queries.
- Verify uniqueness, constraints, attribution, ordering, and transaction behavior.
- Prove expected-version mutations are atomic and distinguish stale/missing
  outcomes; sequential fake-store tests cannot establish concurrency safety.
- Check account creation preserves other users and session changes are scoped.
- Keep tests deterministic and independent; reset through owned database/schema
  lifecycle, never shared production state.
- An override URL is acceptable only when the documented harness supports it
  and it targets isolated throwaway infrastructure.
- Cover fresh migration and reapplication/upgrade behavior when schema changes.

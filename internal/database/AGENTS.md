# Database lifecycle

Keep implementation-specific lifecycle under its database package.

- PostgreSQL migration/connection lifecycle belongs here; ordinary account/card
  queries belong in internal/store/postgres.
- An additional database needs a deliberate implementation and migration
  strategy, not speculative abstraction in the template.
- Use contexts, clear resource ownership, and errors that exclude credentials.
- Read the PostgreSQL and migration subtree guides before changing schema behavior.

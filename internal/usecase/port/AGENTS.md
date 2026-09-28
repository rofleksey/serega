# Shared application ports

Keep only contracts genuinely shared by application responsibilities.

- Transaction describes atomic execution without exposing SQL, pgx, or a concrete
  adapter. Implementations own commit/rollback and context binding.
- Domain-specific Store interfaces remain in their consuming use-case package.
- Do not grow this package into a generic repository/service registry.
- Shared types must not import handlers, generated contracts, or infrastructure.
- Retain a port only while a concrete consumer or adapter contract uses it.

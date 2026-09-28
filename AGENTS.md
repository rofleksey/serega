# Serega agent guide

Serega is a working reference application: one shared Kanban board, multiple
independently authenticated users, and a Go executable serving its embedded
React application. It was created by cloning Akio Main and adapting the
checkout. Preserve the demonstrated engineering boundaries when molding this
template into another product.

Read README.md for commands, ARCHITECTURE.md for ownership, and
docs/AKIO_PATTERNS.md for concrete upstream examples and adaptation rationale.
Subtree AGENTS.md files add focused instructions. Generated output directories
deliberately have no handwritten instructions.

## Scope and product policy

- Complete requested changes without introducing unrelated platforms, workers,
  frameworks, role systems, or infrastructure.
- Every authenticated user can read and change every card. User identity
  provides authentication and attribution; it is not a tenant boundary.
- Keep useful rationale comments. Delete obsolete files and dependencies when
  replacing features; do not preserve dead code as an example.
- Preserve practical usability, keyboard access, labels, and useful
  pending/error/empty states. Expand accessibility or internationalization when
  the task calls for it.
- Do not mutate production unless requested. Ordinary authorized development,
  generation, verification, and fixes do not need extra approval ceremonies.
- Future products can change domain scope deliberately. Update the affected
  architecture and subtree instructions along with changed policy.

## Dependency direction and ownership

Runtime calls follow cmd → cli → app → handler → usecase → store/postgres.
Use cases declare consumer-owned ports; they do not import the concrete store.
The app package constructs and connects concrete dependencies.

- Handlers and CLI commands call application services, never store/database
  methods, including through disguised persistence interfaces.
- Use cases own validation and decisions. Stores own SQL, row mapping,
  transactions, and defensive constraints.
- Generated HTTP and SQL models stay at their boundaries and are explicitly
  converted to entity/use-case values.
- Shared domain values and stable vocabulary belong in internal/entity.
- Inbound protocol handling belongs in handler or a deliberately introduced
  transport package. Adapter means a concrete outbound port implementation.
- Avoid pass-through orchestration layers around transport serialization.

## Go structure and libraries

- Keep one concrete handler per source file, and one middleware constructor or
  value per source file with focused behavioral coverage.
- Shared exported contracts belong in types.go or descriptive types_<domain>.go.
  Private types/helpers stay beside their behavior.
- helper.go is for genuinely shared mechanics, not unrelated leftovers.
- Stable statuses, actions, telemetry keys, and metric names belong in entity;
  prose, HTTP headers, SQL fragments, and local literals normally stay local.
- Preserve one HTTP error-envelope boundary and consistent domain-error mapping.
  Wrap errors with context and use errors.Is/As.
- Use contexts for outbound work and explicit lifecycle ownership.
- Reuse maintained infrastructure libraries. Application logging uses slog;
  unolog owns wide-event lifecycle; the existing meg bridges own generic
  redaction, operation, rate-limit, and metrics mechanics.
- Do not invent custom AST checks or cosmetic linters. Use the standard
  approved toolchain and document limitations.

## Source of truth and generation

- api/openapi.yaml owns HTTP shapes.
- internal/database/postgres/migrations owns Goose schema versions.
- internal/store/postgres/sqlc/query owns static SQL.
- Edit sources and regenerate; never patch generated Go or TypeScript.
- Do not put handwritten helpers or docs inside internal/api/generated or
  frontend dist. SQLC source/configuration lives in sqlc/, and its ignored
  output lives in the separate sqlc/generated/ directory.
- Build Vite assets before compiling Go. Production assets are embedded and
  immutable; do not introduce a mutable runtime asset directory.
- OpenAPI supports generation and validation. Do not expose documentation or
  raw schema endpoints unless the product requires them.

## Persistence, accounts, and collaboration

- Use SQLC for fixed-shape application queries. Handwritten dynamic SQL requires
  a documented necessity, constrained identifiers, and parameterized values.
- Compare expected card versions atomically in SQL. Stale updates and deletes
  must return a conflict rather than overwrite another person's work.
- Account creation and password/session changes affect the intended user.
  Never restore sole-administrator replacement or global session deletion.
- Keep credentials, hashes, cookies, tokens, unrestricted bodies, and connection
  strings out of responses, telemetry, fixtures, and committed configuration.
- Preserve secure cookie defaults, persistent hashed sessions, CSRF, request
  limits, and trusted-proxy handling. Local HTTP allowances must be explicit.

## Frontend

Retain React/TypeScript, MUI, TanStack Query, React Hook Form, Zod, and the
generated OpenAPI client unless the requested task changes the stack.

- Layer direction: app → pages → widgets → features → entities → shared.
  Lower layers never import higher layers; cross-folder imports use aliases.
- Query owns server state; features own form state. Invalidate affected queries
  after successful mutations and clear private query state on logout/expiry.
- Keep drafts on recoverable errors and conflicts. A refreshed list does not
  grant permission to overwrite a user's open form.
- Five-second polling is sufficient for this board. Introduce push transports
  only for a requirement that warrants their lifecycle complexity.
- Server validation remains authoritative; frontend schemas improve feedback.

## Observability

- Emit one rich completion event per owned operation using the existing unolog
  lifecycle and redacting slog sink.
- Start/finish at the owning boundary and enrich via context. One goroutine owns
  completion; independent concurrent work gets its own operation.
- Keep standalone logs for process lifecycle and exceptional diagnostics.
- Metrics use OpenTelemetry/OTLP, bounded labels, bounded timeouts/queues, and a
  bounded shutdown flush. Telemetry failure must not fail application work.
- Never label metrics with usernames, card IDs, raw dynamic paths, or errors.
  Never record card titles/descriptions or authentication material.

## Verification and completion

Use actual Makefile and web/package.json targets. The complete normal gate is
generation, formatting, explicit lint, vet, standalone staticcheck, Go/frontend
tests, and a production build.

- Keep golangci-lint v2 with linters.default: none, all enabled linters explicitly
  listed, and gofmt/goimports formatters. Preserve the approved list and finish
  without findings; exclusions must stay narrow and explained.
- Keep meaningful unit tests beside packages/features. Do not create tests that
  merely duplicate trivial implementation.
- test/ is reserved for integration tests. Application tests run the built
  executable through public boundaries; PostgreSQL tests verify schema/store
  contracts directly.
- Use disposable isolated infrastructure and Testcontainers with pinned images,
  bounded waits, and teardown. Never use production to verify changes.
- Run integration coverage for authentication, persistence, migration,
  concurrency, runtime, and deployment changes.
- A clean checkout must regenerate and build. Do not rely on local ignored
  artifacts. Report unavailable tools/infrastructure precisely; skipped is not
  passed.
- Inspect the diff before committing; exclude secrets, binaries, node_modules,
  generated output, and unrelated edits.

## Adapting into the next project

Follow a complete vertical slice: contract → entity → use-case port/service →
schema/query → store mapping → handler → client/feature/page → meaningful tests.
Rename module imports, executable/environment/cookie/browser-event names, image
names, and docs consistently. Remove unused domain vocabulary. Retain a small
working example of each engineering mechanism the new product still uses.

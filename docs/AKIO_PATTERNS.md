# Akio Main pattern study and Serega adaptation

This reference records concrete patterns observed in Akio Main before editing
the cloned checkout. It is intended to help an agent understand why files live
where they do, then apply the same boundaries to a different product.

## Provenance and method

Upstream: `github.com/rofleksey/akio-main`, commit
`6bb22aadc48f9e739cc5ea2e9887de5379a5ab9c`.

The source checkout was cloned into a separate Serega directory before product
changes. Serega was not reconstructed in an unrelated framework, and the prior
discarded Serega implementation was not used as its source.

The study read the upstream root AGENTS.md and ARCHITECTURE.md, then followed
actual code through command/composition, account authentication, generated HTTP
routing and validation, SQLC persistence, migrations, middleware, observability,
frontend session/forms/layers, build/deployment, and integration harnesses.
This document distinguishes verified source patterns from product-specific
features intentionally removed.

Upstream paths below refer to that fixed commit. They can be inspected using
`git show 6bb22aadc48f9e739cc5ea2e9887de5379a5ab9c:<path>` when inherited history is
available, or the [upstream commit](https://github.com/rofleksey/akio-main/tree/6bb22aadc48f9e739cc5ea2e9887de5379a5ab9c)
when access permits. Serega's own paths refer to the working template.

## 1. Small process entry and explicit composition

Observed upstream:

- `cmd/main.go` delegates to the CLI.
- `internal/cli/root.go`, `serve.go`, `migrate.go`, and administrator
  command files separate Cobra command construction by responsibility.
- `internal/app/runtime.go` opens/pings the pool and owns serving/shutdown.
- `internal/app/helper.go` performs concrete service/adapter wiring.
- `internal/config/config.go` exposes `Load` and a testable `FromLookup`.

Pattern: parsing and lifecycle are explicit boundaries. Neither HTTP nor CLI
should acquire direct persistence responsibilities. Configuration is parsed
once through a declared environment surface.

Serega keeps this layout and serves through the same kind of runtime. Its
`user create` command calls the application/account service path. Akio fleet
background loops were removed because the board needs no fleet controller.

## 2. Domain data, use-case contracts, and small services

Observed upstream:

- `internal/entity/account.go` defines shared account values.
- `internal/entity/vocabulary.go` and `observability.go` hold stable shared
  operation/telemetry vocabulary.
- `internal/usecase/account/types.go` declares the narrow Store port and
  Service structure, with compatibility aliases to entity types.
- `internal/usecase/account/service.go` owns authentication/password/token
  decisions; `password.go` owns Argon2id mechanics.
- Other use-case packages follow the same responsibility-specific layout.

Pattern: interfaces belong to consumers, services own validation/orchestration,
and the concrete store is selected by composition. The project does not need a
generic CRUD framework or DTO hierarchy to preserve these boundaries.

Serega uses `entity/card.go`, `usecase/board/types.go`, and small
`create.go`, `update.go`, and `delete.go` files. The account/password
implementation is retained/adapted. Multi-user semantics replace the upstream
sole-administrator policy.

Critical domain difference: upstream `SetSoleAdministrator` deletes other
users, and its password-change path deletes all sessions. Those policies are
inappropriate for a shared multi-user board and must not be copied back.
Serega provisions independent users and does not expose password-reset or
administrator-management flows.

## 3. Handwritten handlers around generated contracts

Observed upstream:

- `api/openapi.yaml` is canonical.
- `api/oapi-codegen*.yaml` and `internal/api/generate.go` generate models
  and stdlib HTTP server bindings.
- `internal/api/validation.go` and `response.go` wrap generated code from
  a normal handwritten package.
- `internal/handler/login.go` and `create_api_token.go` each implement one
  endpoint, convert input, call a service, and map outcomes.
- `internal/handler/helper.go` centralizes decoding/error response mechanics.
- `internal/handler/router.go` assembles the generated router.

Pattern: generation owns routine transport shape/routing, while handwritten
code owns authentication, application decisions, and error semantics.
Generated directories are not extension points for handwritten helpers.

Serega retains this arrangement with card operations and session endpoints.
`internal/handler/responses.go` maps domain values/errors; individual card
handler files remain thin. The common envelope carries a stable code, message,
request ID, and optional field errors. The API schema is not hosted.

The upstream router also contained product-specific stream and configuration
exceptions. Removing those routes removes their exceptions; there is no reason
to recreate Akio's live-log/config streaming machinery for a polling board.

## 4. SQLC for queries, handwritten transaction composition

Observed upstream:

- `internal/store/postgres/sqlc/sqlc.yaml` reads Goose migrations and query
  source, emits pgx/v5 bindings, and explicitly maps UUID/time/nullable values.
- `internal/store/postgres/account.go` calls generated queries and converts
  rows to domain values.
- Behavior-specific adapter files own their mapping and transaction logic.
- The store's transaction boundary binds generated query objects to the active
  pool/transaction; shared transaction contracts live under `usecase/port`.
- Exceptional dynamic monitoring SQL was documented separately upstream.

Pattern: generated query methods replace repetitive scanning and binding.
Transactions and concurrency semantics stay handwritten because they are part
of the application contract. SQLC does not own migrations.

Serega keeps account/card SQL in `sqlc/query/auth.sql` and `cards.sql`.
Its adapter locks the card, distinguishes missing/stale versions, applies a
version-predicated mutation, and reads the persisted representation in the
same transaction. This is a compact demonstration of the transaction pattern
without Akio's job queues, advisory locks, monitoring ranges, or bulk paths.

Upstream mixed SQLC source/configuration and ignored output in one directory.
Serega retains `sqlc/generate.go`, `sqlc/sqlc.yaml`, and `sqlc/query` as source,
with generated bindings moved into `sqlc/generated/` to make the generator-only
boundary explicit. Never edit that child directory or put handwritten files
inside it.

## 5. Schema lifecycle is separate from serving

Observed upstream:

- `internal/database/postgres/migrations` owns versioned Goose SQL.
- `internal/database/postgres/migrate.go` is an explicit lifecycle entry.
- CLI migration is separate from serving.
- The schema and persistent SCS session storage share PostgreSQL.

Pattern: deploying schema and starting traffic are separate operations.
Migrations stay with their database implementation instead of a generic root
directory detached from lifecycle ownership.

Serega has a fresh users/sessions/cards schema. It uses an instance-based Goose
provider and embedded migrations. Startup performs a read-only schema check.
The one metadata query in `readiness.go` is an explicit infrastructure
exception: Goose helpers that initialize tables cannot run during ordinary
serving. It is not permission to scatter handwritten application SQL.

## 6. Security mechanics are shared boundary components

Observed upstream:

- SCS session manager: hashed tokens, persistent pgx storage, idle/lifetime
  bounds, Secure/HttpOnly/SameSite cookies.
- Login renews the token before storing the authenticated identity.
- nosurf provides CSRF protection.
- Request IDs, proxy trust, body limits, JSON/media checks, recovery, request
  logging, and rate limiting live in separate middleware files.
- `rate_limit.go` delegates to `meg/ratelimit` instead of hand-rolled
  time-window arithmetic.

Pattern: preserve proven generic mechanics when changing domains. Compose
middleware by route ownership, not by parsing raw path segments.

Serega retains these components, adds the small private-response no-store
boundary, and makes cookie security explicit for local HTTP development while
keeping Secure on by default. Bearer-token and Agent enrollment paths were
removed because no remaining feature consumes them.

## 7. Logging is a library lifecycle, not arbitrary log calls

Observed upstream:

- `internal/observability/logger.go` delegates to `meg/logging` and
  `meg/operation`, exposing slog and unolog lifecycle through narrow wrappers.
- `internal/middleware/request_logger.go` starts/finishes one request event.
- Services enrich its context with operation/decision fields.
- `internal/observability/metrics.go` uses OTel instruments and
  `meg/otelmetrics`, not a custom registry/exporter.
- The runtime provides fail-open initialization and shutdown flushing.

Pattern: a safe, schema-stable completion event describes the operation across
layers. IDs can correlate safe event fields, while metrics use bounded labels.
Credentials and unrestricted data never belong in telemetry.

Serega keeps HTTP/database telemetry and removes fleet/webhook/job metrics.
Handler failures attach a safe code and error type to the owning event rather
than logging a raw database error that might contain card/user input.
The optional Collector remains a standard deployment example.

A subtle upstream documentation mismatch was resolved by reading code:
generic pointer helpers and logging/operation primitives already use the
shared `meg` module. Do not recreate older application-owned helpers simply
because an older architecture paragraph described them.

## 8. Feature layers are enforced, not just folder names

Observed upstream:

- `web/src/app/App.tsx` wires Query, MUI theme, feedback, browser routing,
  session provider, and protected routes.
- `web/eslint.config.js` forbids lower layers importing higher layers.
- `web/tsconfig.json` and Vite share stable layer aliases.
- Route pages compose widgets/features/entities/shared modules.
- `web/src/features/administration/form-schemas.ts` and
  `AdministrationDialogs.tsx` pair Zod with React Hook Form.
- Mutations invalidate Query data and translate server field errors.

Pattern: app → pages → widgets → features → entities → shared. File placement
corresponds to responsibility, not arbitrary component size. Shared code cannot
reach back into a feature or route to access a convenient value.

Serega's equivalent slice:

- `pages/board/BoardPage.tsx`: thin route composition.
- `widgets/kanban-board/KanbanBoard.tsx`: shared query, columns, move/delete.
- `features/card-management/CardDialog.tsx`: form/mutation lifecycle.
- `features/card-management/form-schemas.ts`: validation/defaults.
- `entities/card/CardItem.tsx`, `status.ts`: presentation and metadata.
- `entities/session/auth.tsx`: current-user context and session lifecycle.
- `shared/api`, `shared/ui`, `shared/lib`: common infrastructure.

The existing theme, feedback provider, time formatting, required-text helper,
loading/empty/error states, confirmations, and header/action primitives are
retained. Akio monitoring charts, editors, fleet pages, job streams, and their
dependencies are removed.

## 9. Typed browser transport and session invalidation

Observed upstream:

- `web/src/shared/api/transport.ts` configures generated openapi-fetch,
  same-origin credentials, an in-memory CSRF cache, and ApiRequestError.
- `apiFormErrors` provides structured field feedback.
- `session-expiry.ts` centralizes the 419 notification.
- `entities/session/auth.tsx` clears CSRF and Query state on session expiry
  or sign-out and navigates to login.

Pattern: components call a typed client, not ad hoc fetch. Authenticated cache
lifecycle belongs to session handling. Private cached data cannot survive
across identities merely because a page stays mounted.

Serega adapts the provider from administrator to ordinary user semantics,
retains the error/CSRF/session mechanisms, and uses generated card shapes.
Version conflict handling preserves the user's draft; a background polling
response cannot silently become a new edit precondition.

## 10. Immutable frontend and clean-checkout generation

Observed upstream:

- Vite emits `internal/frontend/dist`.
- `internal/frontend/assets.go` embeds the production build.
- Handwritten SPA serving lives in `handler.go`, `helper.go`, and `types.go`.
- Makefile, CI, and Docker build generated bindings and frontend before Go.
- Generated sources and output assets are gitignored.

Pattern: the executable/image is the deployable unit. A mutable runtime assets
directory and production Vite process would weaken that boundary.

Serega retains this build chain and SPA implementation. Names/configuration are
adapted to the new app; generators remain source-driven. A clean clone must work
without copying ignored build output from the original development machine.

## 11. Unit tests and public-boundary integration are different contracts

Observed upstream:

- Focused Go tests sit beside packages; frontend tests sit beside behavior.
- `test/integration/application/harness_test.go` owns a built executable,
  disposable PostgreSQL, a fake OTLP receiver, bounded readiness, log capture,
  signal shutdown, and cleanup.
- Application tests do not import internal application packages.
- `test/integration/postgres` verifies migrations/schema/store behavior with
  a real PostgreSQL harness and isolated test data.
- Testcontainers owns dependency lifecycle.

Pattern: a store mock cannot prove SQL atomicity, and directly calling a
handler cannot prove the release executable starts, embeds its UI, or flushes
telemetry. Keep these suites complementary.

Serega adapts the harnesses for independent accounts and card operations.
Important scenarios are two cookie jars seeing one board, CSRF/auth rejection,
validation, attribution, stale updates/deletes, persistence, schema lifecycle,
and bounded process cleanup. Deterministic tests stay separate from load tests.

## 12. Explicit quality gates and reusable deployment

Observed upstream:

- golangci-lint v2 uses `linters.default: none` and an explicit suite.
- gofmt/goimports, vet, standalone staticcheck, SQLC vet, frontend ESLint/type
  checks, unit/race/integration tests, and builds are separate gates.
- GitHub Actions pins third-party actions.
- Production Docker is multistage, includes generated/embedded assets, and
  runs non-root.
- Compose orders PostgreSQL health, explicit migration, then serving.
- Collector configuration is an optional standard telemetry pipeline.

Pattern: retain executable verification and operational defaults while removing
product-specific external dependencies. A reference template should not need a
private peer checkout or a previous owner's secret to pass CI.

Serega removes the Akio-Agent fixture checkout/protocol secret and old product
publishing assumptions. The board does not preserve unused wire fixtures,
WebSocket hubs, webhooks, monitoring, repository deployment, volume management,
or Agent commands simply to increase resemblance. Their general engineering
patterns remain represented by a smaller working application.

## Concrete source-retention map

The following files were inspected in both the upstream checkout and adapted
checkout. “Retained/adapted” means the same responsibility and source lineage,
not a promise of byte-for-byte identity after renaming and verification fixes.

| Upstream source | Serega source or outcome |
| --- | --- |
| cmd/main.go; internal/cli/serve.go | Same files, renamed application boundary |
| internal/config/config.go | Same loader pattern, smaller SEREGA_* surface |
| internal/app/http.go, runtime.go, helper.go | Same routing/runtime/composition roles, fleet loops removed |
| internal/api/generate.go; api/oapi-codegen*.yaml | Retained generator arrangement |
| internal/api/validation.go, response.go | Retained/adapted validation/error boundary |
| internal/handler/login.go, logout.go, helper.go | Retained/adapted session and response mechanics |
| internal/usecase/account/password.go | Retained Argon2id implementation |
| internal/middleware/request_id.go, trusted_proxy.go, rate_limit.go | Retained generic middleware |
| internal/middleware/request_logger.go, session_manager.go, csrf.go | Retained/adapted to reduced domain and cookie configuration |
| internal/observability/logger.go, metrics.go | Retained bridges; irrelevant domain instruments removed |
| internal/store/postgres/account.go; sqlc/sqlc.yaml | Adapted accounts; retained SQLC generation/type mapping pattern |
| internal/database/postgres | Retained lifecycle location, fresh schema and provider/readiness adaptation |
| internal/frontend/{assets,handler,helper,types}.go | Retained embedded-SPA structure |
| web/src/app/theme.ts; providers/feedback.tsx | Retained presentation foundation |
| web/src/shared/ui; shared/lib | Retained generic primitives/utilities and applicable tests |
| web/src/shared/api/transport.ts; entities/session/auth.tsx | Retained/adapted transport and user session lifecycle |
| web/eslint.config.js; web/tsconfig.json | Retained layer rules and aliases |
| test/integration/application/harness_test.go | Adapted built-process/Testcontainers/telemetry harness |
| test/integration/postgres/harness_test.go | Adapted real-database harness |
| deploy/production; deploy/observability | Retained layout and operational patterns, renamed/reduced configuration |

## How a future agent should use this study

Read the root guide and the guides along the target file's path. Follow the
existing card or account slice end to end before creating a new one. Keep a
small real example of each still-useful mechanism: typed transport, service
policy, persistence, authentication, UI form/query lifecycle, safe observability,
and public-boundary tests.

Replace domain vocabulary deliberately. Update contracts, migrations, frontend
schemas, tests, and the affected documentation together. Remove an unused
mechanism completely instead of leaving a misleading skeleton. Add complexity
only when a new product requirement needs it.

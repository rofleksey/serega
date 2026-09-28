# Serega architecture

Serega is a single shared Kanban board with independently authenticated users.
Its purpose is to demonstrate the Akio Main engineering structure in a small
application that can be cloned and reshaped into another product.

The implementation began with an actual clone of Akio Main at
`6bb22aadc48f9e739cc5ea2e9887de5379a5ab9c`. The upstream pattern study and
concrete retained/adapted file map are in [docs/AKIO_PATTERNS.md](docs/AKIO_PATTERNS.md).
Root and subtree AGENTS.md files turn those findings into local editing rules.

## Product boundary

There is one board and three fixed statuses: `todo`, `doing`, and `done`.
Every authenticated user can list, create, edit, move, and delete cards.
Cards contain a title and optional description, a version, creation/update
timestamps, and creator/updater identity. Creation starts in Todo.

Accounts are provisioned through `serega user create --username ... --password-stdin`.
There is no public registration, browser user administration, password-reset
flow, role hierarchy, tenant model, custom board configuration, or Akio Agent
dependency. These are intentional scope boundaries, not implied unfinished
features.

The board refreshes every five seconds and on window focus. This is ordinary
HTTP polling; no WebSocket coordination service is required. Other clients see
successful changes on their next refresh. Expected-version checks prevent
polling delays from silently losing another user's changes.

## Runtime and dependency map

```text
cmd/main.go
    └── internal/cli                  command parsing, stdin, signals
          └── internal/app           concrete dependency composition/lifecycle
                ├── config           validated process environment
                ├── database/postgres explicit migration/readiness lifecycle
                ├── frontend         embedded production SPA
                ├── observability    slog/unolog/OTel library bridges
                └── handler + middleware
                      └── usecase/account, usecase/board
                            └── consumer-owned Store ports
                                  └── store/postgres
                                        └── SQLC → PostgreSQL
```

This describes runtime calls. Compile-time dependency inversion remains:
use cases define their interfaces and do not import the concrete store.
`app` injects adapters into those ports. `entity` contains shared domain
models and stable vocabulary. Generated HTTP and SQLC models are converted
explicitly at their respective boundaries.

Transport and CLI must not bypass services, even through a store-shaped
interface. Services own validation and decisions; persistence owns transaction
mechanics and defensive constraints.

## Backend vertical slice: editing a card

1. The generated OpenAPI router binds the update operation to
   `internal/handler/update_card.go`.
2. Middleware/handler assembly validates authentication, CSRF, body limits,
   media type, and contract shape. The handler maps generated request data into
   a board use-case input and supplies the authenticated actor.
3. `internal/usecase/board/update.go` validates domain values and calls its
   persistence port with the request context.
4. `internal/store/postgres/board.go` owns the database transaction. It locks
   the card row to distinguish absence from a stale version, then calls the
   SQLC mutation that also compares the expected version.
5. The mutation increments the version and updates actor/time atomically.
   The adapter maps the persisted representation into an entity value.
6. The handler maps known outcomes into a common error envelope or the generated
   card response.

A stale version is a conflict, not a successful overwrite. Deletion uses the
same expected-version protection. A missing card remains distinguishable
from a conflict. The frontend retains the edit draft and offers explicit
review of the latest saved version before resubmission.

## Contracts and generation

| Source of truth | Generated consumer | Owner |
| --- | --- | --- |
| api/openapi.yaml | internal/api/generated/*.gen.go | go:generate + oapi-codegen |
| api/openapi.yaml | web/src/shared/api/generated.ts | openapi-typescript |
| PostgreSQL migrations + sqlc/query/*.sql | internal/store/postgres/sqlc/generated | SQLC |
| web/src + frontend config | internal/frontend/dist | Vite |
| embedded migration SQL | Go migration filesystem | go:embed |

Generated Go/TypeScript/assets are ignored and reproducible. Run
`make generate` before validating a clean checkout. The API contract supports
generation and validation; the application does not host raw OpenAPI or Swagger.

`internal/api` is the handwritten wrapper around a generator-only
`generated` subdirectory. SQLC source stays in `sqlc/generate.go`,
`sqlc/sqlc.yaml`, and `sqlc/query/*.sql`; its ignored bindings live in the
separate `sqlc/generated/` child directory. This small refinement makes the
generator-only boundary explicit without changing query or adapter ownership.

## Database ownership and concurrency

Goose migrations live under `internal/database/postgres/migrations`.
The migration command creates an instance-based Goose provider using the
embedded migration filesystem. It does not mutate a package-global migration
filesystem. Ordinary serving never applies migrations.

Startup checks schema readiness read-only. The query against Goose's metadata
in `internal/database/postgres/readiness.go` is infrastructure lifecycle SQL,
not a second application-query layer; its inline comment explains why Goose
helpers that initialize metadata cannot be used during serving.

Static account/card queries live under `internal/store/postgres/sqlc/query`.
SQLC generates their pgx/v5 bindings. UUIDs, nullable values, and timestamps use
explicit mappings in `sqlc.yaml`.

The PostgreSQL adapter's transaction helper carries the active transaction in
a private context key. Generated queries bind to the pool or that transaction.
The shared transaction port describes atomic execution without exposing SQL
to consumers. Row locks and version predicates keep a card's decision and
mutation in one transaction; use-case validation remains outside the adapter.

## Authentication and HTTP mechanics

Passwords use the inherited Argon2id helper with a random salt, fixed accepted
parameters, and constant-time comparison. Password hashes never enter public
user responses. Account creation adds an independent user and never deletes
existing accounts.

SCS stores hashed session tokens in PostgreSQL. Login renews the session token;
logout destroys it. Session and CSRF cookies are HttpOnly/SameSite and secure
by default. `SEREGA_COOKIE_SECURE=false` is an explicit local HTTP setting,
not a production default. nosurf protects cookie-authenticated mutations.

Request IDs, trusted-proxy parsing, rate limiting, body/media limits, panic
recovery, no-store policy, and consistent error mapping are separate middleware
or boundary responsibilities. Rate limiting reuses the existing maintained
library abstraction; it does not implement custom window arithmetic.

The frontend and API share an origin. Unknown API routes return an API error
instead of the SPA document. Health and readiness routes have distinct roles.
The runtime owns HTTP deadlines, cleanup of persistent-session background
work, the database pool, and a bounded telemetry shutdown flush.

## Frontend layers

| Layer | Concrete examples | Owns |
| --- | --- | --- |
| app | App.tsx, AppLayout.tsx, RequireUser.tsx | Providers, routes, protected shell, theme |
| pages | board/BoardPage.tsx, login/LoginPage.tsx | Route experience |
| widgets | kanban-board/KanbanBoard.tsx | Cards query, columns, move/delete composition |
| features | card-management/CardDialog.tsx, form-schemas.ts | Create/edit form and mutation feedback |
| entities | card/CardItem.tsx, card/status.ts, session/auth.tsx | Domain presentation and session identity |
| shared | api, ui, lib | Typed transport, generic UI, small utilities |

ESLint enforces downward dependency direction. Stable aliases avoid fragile
cross-folder relative imports. The existing MUI theme, feedback provider,
async states, confirmation dialog, and related primitives are reused.

TanStack Query owns server data under the `cards` query key. Mutations
invalidate it on success. React Hook Form and Zod own form state and immediate
validation. The API remains authoritative.

The editing card is captured as a draft base; background refetches do not
replace its expected version. A conflict keeps the draft and displays the
latest saved card for review. Keyboard-operable status controls provide
movement without requiring drag-and-drop.

The shared OpenAPI transport owns same-origin credentials, an in-memory CSRF
token cache, structured request errors, field-error mapping, and session-expiry
signaling. The session provider clears private query/CSRF state on logout or
expiry. Credentials are not stored in localStorage.

## Observability

Application code uses `log/slog`. The narrow adapters in
`internal/observability` reuse `meg/logging`, `meg/operation`, unolog, and
OpenTelemetry rather than reimplementing generic infrastructure.

The request boundary owns one completion event. Handlers/services contribute
safe fields through the context. Stable event/decision/metric vocabulary lives
in `entity`. Card text and authentication material must never enter telemetry.

Metrics cover HTTP and database pool behavior using bounded method/route/status
attributes. IDs and raw paths are not metric labels. OTLP initialization/export
fails open and shutdown flushing is bounded. The optional Collector example
lives in `deploy/observability`.

## Verification and distribution

Focused Go tests live beside packages; frontend Vitest/Testing Library tests
live beside components. `test/integration/application` runs the actual
executable through public CLI/HTTP/process/logging/OTLP boundaries.
`test/integration/postgres` verifies migrations, schema, transactions, and
adapter behavior against a real isolated PostgreSQL database.

The Makefile and CI generate from source, run explicit lint/static checks,
test, and build. Containerized integration infrastructure is disposable and
owned by Testcontainers. Missing infrastructure must be reported, not treated
as a passing gate.

Vite writes `internal/frontend/dist`; Go embeds it into the executable.
The production Dockerfile builds frontend and generated Go code before the
non-root runtime image. Compose orders database readiness, explicit migration,
then serving. Deployments provide secrets externally and use HTTPS.

## Evolving the template

Keep the architecture while replacing product policy. For a new domain, follow
the card/account slice through contract, entity, service port, schema/query,
adapter, handler, frontend layers, and tests. Add a layer or service only when
it has a real responsibility. Remove unused Serega vocabulary and revise
nearby agent guides as the product changes.

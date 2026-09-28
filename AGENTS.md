# Serega agent guide

Serega is a standalone project template with a working shared Kanban board:
independently authenticated users, PostgreSQL persistence, and a Go executable
serving an embedded React application. Every authenticated user can change every
card; identity provides authentication and attribution, not tenancy.

## Where guidance lives

[README.md](README.md) owns setup, commands, and configuration. This guide owns
repository-wide rules; read the scoped guides along the path you are editing.
Child guides add local constraints without restating their ancestors. Add a new
one only when a subtree needs substantial, distinct guidance.

| Area | Guide |
| --- | --- |
| Go package ownership and runtime | [internal/AGENTS.md](internal/AGENTS.md) |
| HTTP contract | [api/AGENTS.md](api/AGENTS.md) |
| Frontend layers and interaction | [web/AGENTS.md](web/AGENTS.md) |
| Integration testing | [test/AGENTS.md](test/AGENTS.md) |
| Containers and telemetry deployment | [deploy/AGENTS.md](deploy/AGENTS.md) |
| CI | [.github/workflows/AGENTS.md](.github/workflows/AGENTS.md) |

## Architecture and shared conventions

Runtime calls follow `cmd → cli → app → handler → usecase → store/postgres`.
Use cases declare consumer-owned ports; they never import concrete stores.
`app` constructs and connects dependencies. CLI commands also invoke services
through `app`, without going through HTTP handlers.

- Keep `cmd/main.go` limited to CLI delegation, process output, and exit status.
  Handlers and commands must not bypass services, even through a store-shaped
  interface. Services own validation/decisions; stores own SQL, row mapping,
  transactions, and defensive constraints.
- Explicitly convert generated HTTP/SQL models at their boundaries. Shared
  domain values, statuses, actions, and telemetry vocabulary belong in
  `internal/entity`; local prose, headers, SQL, and UI text stay local.
- Inbound protocol work belongs in handlers or a transport package; an adapter
  implements an outbound port. Add layers only when they own a real responsibility.
- Keep one concrete handler and one middleware constructor/value per source
  file. Put shared exported contracts in `types.go` or `types_<domain>.go`, and
  private types/helpers beside their behavior. Reserve `helper.go` for shared
  mechanics within that package.
- Wrap errors with context and use `errors.Is/As`. Keep one HTTP error-envelope
  boundary and consistent domain-error mapping. Propagate contexts for outbound
  work and give resources explicit lifecycle owners.
- Reuse maintained infrastructure libraries and existing helpers; application
  logging uses `slog`. Avoid custom CRUD frameworks, global service locators,
  utility frameworks, AST checks, and cosmetic linters.
- Keep credentials, hashes, cookies, tokens, unrestricted bodies, and connection
  strings out of public responses, telemetry, fixtures, and committed config.
- Preserve usable keyboard controls, labels, and pending/error/empty states.
  Expand accessibility or internationalization when the task calls for it.
- Keep useful rationale comments; remove obsolete code and dependencies when
  replacing features. Do not add unrelated platforms or infrastructure.
- Do not mutate production unless requested. Authorized development, generation,
  verification, and fixes need no extra approval steps.

## Sources and generated output

| Source of truth | Generated output |
| --- | --- |
| `api/openapi.yaml` | `internal/api/generated/*.gen.go`, `web/src/shared/api/generated.ts` |
| `internal/database/postgres/migrations` + `internal/store/postgres/sqlc/query` | `internal/store/postgres/sqlc/generated/` |
| `web/src` and frontend configuration | `internal/frontend/dist/` |

Edit sources and regenerate with `make generate`; never patch generated files.
Keep generated output ignored and free of handwritten helpers or guidance.
SQLC configuration and generation directives remain in its parent `sqlc/` directory.
Go embeds the versioned Goose migrations and the Vite production assets; build
the frontend before compiling Go. The executable/image is the deployable unit,
with immutable assets and no production Vite process or mutable asset directory.

## Verification

Use the actual Makefile and `web/package.json` targets. The normal code-change
gate is generation, formatting, explicit lint, vet, standalone staticcheck,
Go/frontend tests, and a production build. See the README for commands.

- Keep golangci-lint v2 with `linters.default: none`, explicitly enabled linters,
  and gofmt/goimports. Preserve the configured suite; keep exclusions narrow and
  explained, and finish without findings.
- Keep meaningful unit tests beside the behavior they cover; avoid tests that
  duplicate trivial implementation. `test/` owns integration coverage.
- Run integration coverage for authentication, persistence, migration,
  concurrency, runtime, and deployment changes using isolated disposable
  infrastructure. Never use production to verify changes.
- A clean checkout must regenerate and build without copied ignored artifacts.
  Report unavailable tools/infrastructure accurately; skipped checks are not passes.
- For documentation-only edits, check accuracy, local links, stale references,
  and the diff; runtime suites are unnecessary unless behavior also changes.
- Inspect the diff before committing. Exclude secrets, binaries, `node_modules`,
  generated output, and unrelated edits.

## Adapting into the next project

Follow a complete vertical slice: contract → entity → use-case port/service →
schema/query → store mapping → handler → client/feature/page → meaningful tests.
Keep a small working example of each mechanism the new product still uses.

Rename module/import prefixes, package/executable names, environment variables,
cookies, browser events, metrics, images, and product copy consistently. Replace
domain policy deliberately, remove unused vocabulary/mechanisms, and update the
contracts, schemas, tests, and affected guides together.

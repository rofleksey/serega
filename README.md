# Serega

A standalone project template with a small shared Kanban board as its working
example. One board, three columns, individual accounts.

Create cards with a title and description, edit them, move them between **To do**,
**Doing**, and **Done**, and delete them with confirmation. Every signed-in user
can collaborate on every card. Changes refresh across browsers every five seconds
and on window focus. Version checks reject stale edits and deletes, so concurrent
work is never silently overwritten. PostgreSQL persists accounts, sessions, and
cards.

## Start locally

Requirements: Go **1.27.1**, Node.js **24.21.0** (`.nvmrc`), Docker with Compose,
and `make`. PostgreSQL **17.6** runs as a local dependency. The dependency lockfiles
and Go tool declarations pin application libraries and generators.

```sh
docker compose up -d --wait postgres
export DATABASE_URL='postgres://serega:serega@127.0.0.1:5432/serega?sslmode=disable'
export SEREGA_HTTP_ADDR=127.0.0.1:8080
export SEREGA_COOKIE_SECURE=false
export OTEL_METRICS_EXPORTER=none
npm --prefix web ci
make build
./bin/serega migrate up
./bin/serega user create --username alice
./bin/serega user create --username bob
./bin/serega serve
```

Account creation prompts for a password twice with terminal echo disabled
(minimum 12 characters). `--password-stdin` accepts one password from standard
input for automation. Passwords are Argon2id hashed. Creating another account
never replaces existing accounts or signs other users out. There is no public
signup, default password, account role hierarchy, or separate board membership.

Open [localhost:8080](http://localhost:8080), sign in, and use a second browser or
private window for the other account. The Go executable serves the embedded
production frontend. The serving process checks schema readiness and never runs
migrations itself. Re-running `migrate up` is safe.

`.env.example` is a reference; the executable does not automatically load dotenv
files. `SEREGA_COOKIE_SECURE=false` is for local HTTP development. Keep the default
`true` behind HTTPS when deploying. `docker compose down` preserves the board;
`docker compose down -v` destroys the local database.

## Use as a template

Choose **Use this template** on GitHub, or clone the repository into a new project.
Start with [AGENTS.md](AGENTS.md) for architecture and repository-wide rules,
then read the scoped guides for the areas you change. Each guide adds only
the ownership, examples, and constraints specific to its subtree.

1. Rename the module/import prefix, package name, executable, environment prefix,
   cookie names, metric prefix, and product copy for the new project.
2. Replace the card domain from its OpenAPI contract through use cases, SQLC
   queries, migrations, frontend feature, and public-boundary tests.
3. Keep applicable foundations: explicit composition, generated contracts,
   secure sessions/CSRF, immutable embedded UI, structured logs and metrics,
   database migrations, and reproducible verification.
4. Remove components only when the new product has no corresponding concern.
   Update the affected `AGENTS.md` guides with the actual decisions.

This is intentionally one service and one board. It omits organizations, projects,
attachments, comments, due dates, drag ordering, background workers, and external
integrations. Polling keeps collaboration simple and works across server
processes without an in-memory broadcast service.

## Verify

Install **golangci-lint v2.14.0**, then:

```sh
make generate          # Go/TypeScript API, SQLC, production frontend
make format            # gofmt + goimports from the lint config
make check             # strict Go lint, vet, staticcheck, frontend checks
make test              # Go units + Vitest/React Testing Library
make test-integration  # built-process and database tests, disposable PostgreSQL
```

CI runs the same generation, lint, static analysis, unit tests, race checks, build,
and isolated integration tests. No private repository checkout, deploy key, or
production service is required. Generated files and embedded build assets are
ignored and regenerated from tracked sources. Go tools require the frontend
build first; `make generate`, `make build`, `make check`, and `make test` arrange it.

## Configuration and operations

| Variable | Default / purpose |
| --- | --- |
| `DATABASE_URL` | Required PostgreSQL connection URL |
| `SEREGA_HTTP_ADDR` | `:8080` |
| `SEREGA_COOKIE_SECURE` | `true`; disable explicitly for local HTTP |
| `SEREGA_LOG_FORMAT` | `json`; also accepts `text` |
| `SEREGA_SHUTDOWN_TIMEOUT` | `15s` |
| `SEREGA_TRUSTED_PROXY_CIDRS` | Empty; trust no forwarding headers |
| `OTEL_METRICS_EXPORTER` | Set `none` locally or `otlp` with a collector |

`GET /healthz` checks process liveness and `GET /readyz` checks database
availability. Request completion events use redacted `slog`/`unolog`; metrics use
OpenTelemetry with bounded dimensions and optional OTLP/HTTP export. Export
failure does not block requests. See [container guidance](deploy/production/README.md)
and the [collector example](deploy/observability/README.md). No production target
or automatic deployment is configured.

The authoritative API is [api/openapi.yaml](api/openapi.yaml). Session/CSRF
endpoints live at `/api/v1/auth`; card endpoints live at `/api/v1/cards`. All
mutations require a CSRF token. The raw API contract and Swagger UI are not
served by the application.

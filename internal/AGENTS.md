# Backend guide

Read the root guide first. Focused guides add rules for services, middleware,
PostgreSQL lifecycle/persistence, and observability; other packages use this guide.

| Package | Responsibility |
| --- | --- |
| app | Concrete wiring, routing, startup, shutdown |
| cli | Cobra command construction, input and presentation |
| config | Declared process environment and startup validation |
| entity | Shared domain models and vocabulary |
| api | Handwritten contract validation and response encoding |
| handler | Endpoint decoding, session handling, service/response mapping |
| middleware | HTTP security and request mechanics |
| usecase | Account/board policy and consumer-owned ports |
| database/postgres | Explicit Goose lifecycle and schema readiness |
| store/postgres | SQLC adapter, transactions, row conversion |
| frontend | Embedded SPA serving |
| observability | Logging/operation bridges and HTTP/database metrics |

## Composition and lifecycle

- `app/helper.go` wires adapters, services, session storage, and handlers;
  `app/http.go` owns route/middleware composition; `app/runtime.go` owns serving.
- Open, ping, and check schema readiness before accepting traffic. Release
  partially opened resources if startup fails. Schema changes belong to the
  explicit migration command; see `database/postgres/AGENTS.md`.
- Keep liveness separate from dependency readiness, with safe diagnostics.
- UI and API share one origin. Unknown API routes return the API error envelope
  rather than the SPA document.
- Bound HTTP deadlines and shutdown. Drain requests, stop session cleanup,
  close the pool, and flush telemetry with a deadline. Lifecycle interfaces
  can support focused tests without adding global registries.

## Commands and configuration

- Keep descriptive Cobra command files and constructors. Validate flags and
  arguments with Cobra, use its input/output streams where practical, and
  return errors to the process entry point without duplicate printing.
- `serve` uses a signal-cancellable context; `migrate up` explicitly changes
  schema; `user create` provisions an independent account.
- Password input uses a hidden terminal prompt with confirmation or the bounded
  `--password-stdin` path. Do not introduce credential flags or environment input.
- Keep serving environment parsing in `config.Load` → `FromLookup`; test the
  latter without mutating process environment. Migration/provisioning commands
  may load only the configuration they require.
- Validate addresses, durations, booleans, log modes, and proxy CIDRs at startup.
  Update README and deployment examples when the supported variables change.
- Cover changed input, flags, malformed/default configuration, and exit behavior.

## HTTP boundary

- `api/validation.go` owns contract validation and schema loading;
  `api/response.go` owns the envelope with code, safe message, request ID,
  and optional field errors. Do not expose decoder/database errors to clients.
- Handlers decode generated types, resolve the authenticated actor, call the
  service with request context, and map results in `handler/responses.go`.
  Schema validation alone does not establish identity, existence, or version.
- Keep missing-card and stale-version responses distinct. Session expiry must
  agree with the browser's expiry flow and differ from rejected login.
- Login renews the session token; logout destroys it; current-user lookup uses
  the current session. Apply security middleware through router composition.
- Focus HTTP tests on mapping, authentication, and malformed/unsupported input;
  application integration tests cover the complete public path.

## Embedded SPA

UI source lives in `web/`. Keep `frontend` focused on safe GET/HEAD asset serving,
path handling, SPA fallback, and appropriate asset/HTML caching. The app router
owns API fallbacks. Test handwritten serving here and UI behavior under `web/src`.

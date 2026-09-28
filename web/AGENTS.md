# Frontend guide

Use React/TypeScript, MUI, TanStack Query, React Hook Form, Zod, and the
OpenAPI client. This guide owns frontend-specific rules; repository-wide
policy, generation, and verification live in the root AGENTS.md.

## Ownership and imports

Dependencies flow app → pages → widgets → features → entities → shared.
ESLint enforces layer boundaries; pages also cannot import sibling pages.
Use @app, @pages, @widgets, @features, @entities, and @shared aliases across
folders, with relative imports for local siblings. Keep aliases aligned in
TypeScript and Vite. Prefer named exports, PascalCase component filenames,
and kebab-case utility/schema filenames.

| Layer | Responsibility and examples |
| --- | --- |
| app | Providers, router, protected shell, theme, generic test setup: App.tsx, layout/AppLayout.tsx, routing/RequireUser.tsx |
| pages | Thin route composition: board/BoardPage.tsx, login/LoginPage.tsx |
| widgets | Screen sections and their server state: kanban-board/KanbanBoard.tsx |
| features | User actions, form schemas/defaults, mutation feedback: card-management/CardDialog.tsx, form-schemas.ts |
| entities | Reusable domain presentation/metadata and identity: card/CardItem.tsx, status.ts, session/auth.tsx |
| shared | api for typed transport, ui for domain-independent primitives, lib for small reusable mechanics |

Wire QueryClient, theme/CssBaseline, feedback, routing, and session providers
once in app. Keep route protection there; the server owns authorization.
Shared components and utilities must not depend on domain policy or routes.
Keep feature mocks/scenarios beside their tests, outside generic test setup.

## Server state and drafts

- TanStack Query owns remote data and mutations. Keep form drafts local rather
  than duplicating live server data in another global store.
- The board widget owns the `cards` query, column composition, and move/delete
  actions. Poll every five seconds and refresh on focus; preserve existing data
  while showing a recoverable refresh error. Add push transport only when a
  product requirement warrants its lifecycle complexity.
- Features own substantial forms through React Hook Form and Zod, with schemas
  and defaults beside the behavior. Mirror API limits for immediate feedback;
  server validation remains authoritative. Show field and general errors.
- Invalidate affected queries after successful mutations. A conflict may also
  refresh the displayed list, but must never retry an overwrite automatically.
- Capture the expected card version when opening a draft. Background refetches
  cannot change it. On conflict, preserve input and require explicit review of
  the latest saved version before resubmission. Preserve a deleted card's draft
  so the user can copy it. Reset only on success or intentional cancellation.
- Disable repeated or contradictory actions while pending. Deletion confirms
  the target and effect and submits the version the user reviewed.
- Card presentation receives callbacks instead of owning mutation/navigation
  workflows. Separate todo/doing/done transport values from labels/colors, and
  keep ID, version, attribution, and timestamps distinct. Render user text as
  text with readable wrapping, never as injected HTML. Requests do not select
  arbitrary actor identities.

## Transport and session lifecycle

- Components use shared/api's typed client and generated schema aliases, not
  duplicated transport interfaces or ad hoc fetch calls.
- client/transport modules own openapi-fetch, same-origin credentials, mutation
  CSRF headers, and the in-memory CSRF cache. Credentials and session/CSRF tokens
  must not enter browser persistent storage; login passwords stay in form memory.
- ApiRequestError and apiFormErrors provide common envelope/field handling.
  Preserve request IDs for diagnostics without exposing server internals, and
  handle no-content success without parsing an empty body as JSON.
- shared/api centralizes session-expiry signaling. The entities/session provider
  obtains the current user through Query and handles expiry once, cancelling
  outstanding queries and clearing private Query and CSRF state on logout or
  expiry so late responses cannot restore another identity's cache.
- After login, refresh the current user and honor a safe intended destination.
  Distinguish session expiry from invalid credentials; authentication errors
  must not reveal whether an account exists.

## UI and focused coverage

- Reuse the MUI theme and responsive sx layouts, feedback provider, async states,
  confirmations, headers, and action menus when their semantics fit. Generic
  utilities such as time formatting and required-text validation belong in lib;
  handle missing/invalid inputs where their caller contract permits them.
- Keep empty columns and the create path clear. Status changes must remain
  keyboard-operable without dragging. Confirmation and dialog controls need
  useful labels, focus behavior, and pending states; show retry only when supported.
- Use Vitest and Testing Library roles/labels for observable behavior, not
  implementation snapshots or CSS-class assertions. Cover form validation,
  pending/success/error states, draft-safe conflicts, mutation refresh, and
  routing/session transitions where changed. Transport changes need coverage
  for method/path/header/body construction, error envelopes, empty responses,
  and expiry signaling; test reusable behavior at its owning layer.

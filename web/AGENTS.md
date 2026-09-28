# Frontend engineering guide

This is Akio's React/TypeScript/MUI frontend structure adapted to one shared
Kanban board. Preserve existing libraries and layer ownership.

## Structure and imports

- app wires providers, routing, layout, and theme.
- pages compose route-level experiences.
- widgets compose substantial reusable screen sections.
- features own user actions and forms.
- entities own domain presentation/session identity.
- shared owns transport, generic UI, and small utilities.
- Dependency direction is downward in that list. ESLint enforces the boundary.
- Use @app, @pages, @widgets, @features, @entities, and @shared aliases across
  folders. Local sibling imports may be relative.
- Keep PascalCase component files and kebab-case utility/schema files consistent
  with neighboring Akio code. Prefer named exports.

## State, forms, and HTTP

- TanStack Query owns remote state and mutation lifecycle. Do not build another
  global store or copy live server data into unrelated component state.
- React Hook Form + Zod own substantial create/edit forms. Defaults and schemas
  stay beside the feature. The server still validates authoritatively.
- Use generated OpenAPI types via shared/api. Do not invent duplicate transport
  interfaces or sprinkle fetch calls through components.
- shared/api owns cookie credentials, CSRF headers, common errors, and session
  expiry notifications. Never store sessions/passwords in localStorage.
- Invalidate affected query keys after success. Keep forms/drafts open after
  recoverable failures; map field errors where possible.
- Refetch/polling must not silently replace an open edit's expected version.
  A user reviews/reloads the latest card before submitting after a conflict.

## UI and verification

- Preserve MUI theme tokens, responsive sx layouts, semantic headings, labels,
  and clear loading/empty/error/pending states.
- Reuse async-states, confirmation-dialog, page-header, and other existing
  shared primitives where their semantics fit.
- Status changes need a keyboard-usable control. Do not make dragging the only
  way to use the board.
- Show destructive actions clearly and keep them disabled while pending.
- Tests live beside the behavior; use Testing Library roles/labels and Vitest,
  not implementation-detail snapshots or class-name assertions.
- Run generate, check, test, and build scripts. Vite writes ../internal/frontend/dist,
  which Go embeds. A standalone browser bundle is not the final deployable app.
- Keep dependency versions/lockfile aligned and remove abandoned Akio-only
  dependencies when removing their features.

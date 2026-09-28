# Shared UI primitives

Preserve Akio's generic loading, empty, recoverable-error, confirmation, header,
and action-menu components where they fit.

- No domain services, feature imports, route assumptions, or card-specific policy.
- Keep labels, headings, focus/dialog behavior, and keyboard controls meaningful.
- Use MUI/theme tokens and responsive sx patterns rather than parallel styling.
- Confirmations identify target, effect, action, and pending state.
- Error components show a useful safe message and retry only when supported.
- Test reusable behavior through roles/labels; avoid implementation snapshots.

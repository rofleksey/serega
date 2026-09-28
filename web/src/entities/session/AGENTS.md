# Session identity and cache lifecycle

Adapt Akio's session provider to ordinary independent users.

- Current user comes from shared/api through Query; do not infer it from a
  client-stored token or remember a password.
- Clear CSRF cache and private Query state on logout/session expiry.
- Cancel outstanding requests where appropriate so a late response cannot
  repopulate another session's cache.
- The expiry event is centralized in shared/api and handled once by the provider.
- Use safe intended navigation after login; distinguish session expiry from
  invalid credentials.
- Do not retain sole-administrator naming/assumptions or invent client-only roles.
- The server remains the authorization boundary.

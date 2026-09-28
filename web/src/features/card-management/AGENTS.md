# Card forms and mutations

Own card create/edit form schemas, defaults, dialog behavior, and corresponding
mutation feedback.

- Title, description, and status validation mirrors the API for early feedback;
  server validation remains authoritative.
- Use generated request types. Never send arbitrary actor IDs.
- Capture the card version when opening/editing a draft and submit that version.
- On conflict, retain the draft and explain that the card changed. Reload/review
  the latest version explicitly before submitting again; do not silently
  substitute a background-refetched version.
- Do not close or reset on failed mutations. Reset intentionally on successful
  completion or confirmed cancellation.
- Map common API field errors to form fields and show any general error.
- Invalidate the cards query on success, not before it.
- Destructive confirmation must name the intended target and preserve the
  expected-version check.
- Cover required/bounded fields, pending states, success reset, and draft-safe
  conflict handling in focused tests.

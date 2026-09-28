# Application services

- Declare each narrow Store port with its consumer in `types.go` and inject
  dependencies through constructors. Use entity/standard Go values, without
  HTTP statuses, generated models, or driver types.
- Own input normalization/validation even when the contract and UI validate it.
  Expose stable domain errors for decisions the HTTP boundary must distinguish.
- Preserve dependency error causes and request context. Enrich the active
  operation with safe decision fields instead of logging user input.
- `port.Transaction` describes atomic execution without exposing SQL; its
  implementation owns commit/rollback and context binding. Keep only contracts
  with concrete consumers or adapter use in this shared package.
- Use small fake ports to prove validation, decisions, and dependency failures,
  including that rejected input never reaches persistence.

## Accounts

`account` owns provisioning, username/password validation, hashing, and
authentication. Preserve Argon2id, random salts, fixed accepted parameters and
lengths, and constant-time comparison. Store only hashes and keep them out of
public user models. Login failures must not disclose whether a username exists.
Creation adds an independent user; password/session changes affect only that
account. Do not infer privileges from account order or invalidate all sessions.
Accounts are CLI-provisioned; public registration, password reset, and browser
user administration are outside the current scope. Test credential failures,
duplicate/invalid identities, and user independence.

## Shared board

`board` validates title, description, status, identifiers, and expected version.
Creation starts in `todo`; supported statuses are `todo`, `doing`, and `done`.
Changing statuses requires aligned SQL constraints, service rules, OpenAPI, and
UI metadata; display labels belong in the UI. The authenticated actor supplies
attribution, never payload creator/updater IDs. Version is the edit precondition
and remains separate from identity, timestamps, and attribution.

Keep success, missing, and conflict outcomes distinct for updates and deletes:
a stale confirmation cannot delete another user's revision. PostgreSQL tests
prove atomicity; HTTP integration tests prove collaboration across independent
sessions.

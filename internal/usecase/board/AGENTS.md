# Shared board policy

One board, statuses todo/doing/done, equal authenticated access.

- Validate title, description, status, identifiers, and expected version here.
  Keep server rules consistent with contract and form feedback.
- The authenticated actor supplies attribution; never trust payload creator or
  updater IDs.
- No implicit tenants, personal lists, board memberships, role system, or
  custom workflows.
- Keep atomic mutation success/missing/conflict outcomes distinct so clients
  can recover from concurrent edits.
- Deletion also uses expected version; a stale confirmation must not remove a
  card another user has changed.
- Ports use domain types, not SQLC/pgx or HTTP statuses.
- Enrich the active operation safely; no card titles/descriptions in telemetry.
- Service tests prove validation/decisions, PostgreSQL tests prove atomicity,
  and HTTP integration proves collaboration across independent sessions.

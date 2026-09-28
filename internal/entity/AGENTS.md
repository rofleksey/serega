# Domain data and vocabulary

Keep domain values independent from HTTP, SQLC, and database-driver types.

- Shared account/card models live here. Password hashes are internal data and
  must never be passed through generated transport models.
- Status values are todo, doing, done. Display labels belong in frontend code.
  Status changes require aligned SQL constraints, service rules, OpenAPI, and UI.
- Card identity, version, timestamps, and attribution have separate semantics.
  Version enforces collaboration; it is not merely a display field.
- Stable shared operation names, decisions, metric names, and telemetry keys
  belong here. Local prose, headers, SQL, and UI text should remain local.
- Avoid generic base entities or speculative common models.
- Reuse existing maintained helpers rather than growing a utility framework.

# Board policy

The template's working example is one shared Kanban board. Replace this document
together with the domain when adapting the template.

## Identity and collaboration

- Accounts are provisioned by operators through the CLI; there is no public
  signup, default password, role hierarchy, or per-board membership. Every
  authenticated account can read and change every card.
- Identity provides authentication and attribution, not tenancy: the acting
  account supplies `createdBy` and `updatedBy`; requests never select an actor.
- Collaboration is deliberately poll-based: browsers refresh every five seconds
  and on window focus. Do not add a push transport without a product
  requirement; polling keeps collaboration correct across server processes
  without an in-memory broadcast service.

## Card lifecycle and versioning

- Statuses are `todo`, `doing`, and `done`; new cards start in `todo`. Changing
  the supported set requires aligned OpenAPI, SQL constraints, service
  validation, and UI metadata; display labels belong to the UI.
- `version` is the edit precondition for updates and deletes and stays separate
  from identity, timestamps, and attribution.
- Updates and deletes are version-predicated mutations executed in one
  transaction with the row locked while deciding missing versus stale. Stale
  edits and deletions return conflict (`409`), missing cards return `404`, and
  the two outcomes stay distinct.
- The service validates title, description, status, identifiers, and expected
  version even though the contract and UI validate too; persistence enforces
  its constraints defensively.
- A conflicting edit is never retried automatically: the client preserves the
  draft and requires explicit review of the latest saved version.

## Presentation

- The UI derives labels and colors from the stable status values; keep those
  values out of presentation and SQL. Card text renders as text with readable
  wrapping.
- Empty columns, pending actions, and recoverable errors have explicit states;
  confirmations name the target and effect.

## Telemetry

- Card titles, descriptions, identifiers, and unrestricted payloads never enter
  logs, metrics, or events. Request events carry bounded fields only.

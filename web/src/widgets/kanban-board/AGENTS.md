# Shared Kanban widget

Own the board's server query and column/card composition.

- Use the shared cards query key and five-second polling plus focus refresh.
  Polling shows other users' changes without extra server protocol machinery.
- Preserve background data while showing recoverable refresh errors.
- Group cards using the canonical todo/doing/done statuses and stable keys.
- Show meaningful empty columns and a clear create path.
- Create/edit forms belong in card-management; card presentation/status metadata
  belongs in entities/card.
- Successful mutations invalidate cards. Conflict handling must not silently
  retry with a newly fetched version and overwrite another user's edit.
- Editing a draft remains independent from background refetches.
- Keep actions usable by keyboard and on narrow screens; dragging is optional,
  never the only status-change mechanism.
- Verify multi-column rendering, mutation refresh, failure/conflict feedback,
  and available actions through user-facing tests.

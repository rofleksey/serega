# Board route

This page represents the single shared board.

- Compose the board widget and its actions; keep server-state/mutation ownership
  where the widget/feature guides specify.
- Do not introduce a board selector, project hierarchy, or implicit personal board.
- Ensure the heading, shared nature, refresh/error state, and create action are clear.
- Keep useful content readable on narrow screens and preserve keyboard controls.
- Domain interactions belong in card-management and kanban-board, not duplicated here.

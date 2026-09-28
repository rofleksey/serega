# Route pages

Pages compose lower layers for a route; they do not depend on app or sibling pages.

- Use shared headers/states and widgets/features/entities where ownership fits.
- Put reusable actions/forms in features and reusable substantial compositions
  in widgets rather than building a second implementation in another page.
- Login owns the login route experience; board owns the one-board route.
- Respect URL/navigation state and route protection without introducing hidden
  alternate data stores.
- Keep loading, recoverable failure, and empty states distinct.
- No raw fetch, generated-type copies, database assumptions, or client-only
  authorization policies.
- Test user-visible route behavior, not incidental component nesting.

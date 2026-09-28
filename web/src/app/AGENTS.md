# Application composition

Own providers, routes, layout, theme, and test setup.

- Wire QueryClient, theme/CssBaseline, feedback, router, and user/session provider
  once in a clear order. Lower layers must not import app.
- Keep protected routes behind the session guard; client guards improve UX but
  never replace server authorization.
- Route pages are composition points. Layout should not own card persistence or
  form validation.
- Keep session expiry/logout clearing all private QueryClient data so another
  user cannot see cached data from the previous session.
- Use theme values and existing responsive layout patterns.
- Add pages/providers only for real product needs; the template is a small app.
- Keep test setup generic; feature mocks and scenarios belong beside their tests.
- Cover routing/session transitions with behavioral tests.

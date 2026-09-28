# Immutable embedded frontend

Edit UI source in web/. This package serves its Vite build embedded in Go.

- Build internal/frontend/dist before compiling the executable.
- dist/ is ignored generator output: no edits, commits, docs, or reliance on
  pre-existing local files.
- Preserve go:embed. Production has no mutable runtime asset directory or
  required Vite server.
- Separate asset serving from SPA fallback; unknown API routes are handled by
  the application router.
- Preserve safe path/method handling and sensible static asset versus HTML caching.
- Test handwritten serving behavior here; UI behavior is tested under web/src.

# HTTP middleware

One middleware constructor/value per file with focused behavioral coverage.
Related context accessors may stay together; helper.go is only for shared mechanics.

- Scope middleware through router composition, not path parsing.
- Keep request IDs bounded and safe, with consistent error/event correlation.
- Trust forwarded client identity only from configured proxies before rate limiting.
- Reuse meg/ratelimit and maintained primitives. Keep per-client storage bounded;
  do not implement custom time-window arithmetic.
- Preserve request body/media limits, panic recovery, and the common API envelope.
- SCS stores hashed session tokens; cookies remain HttpOnly and SameSite with
  Secure enabled by default.
- nosurf CSRF uses the same cookie-security configuration. Cookie-authenticated
  unsafe methods must not bypass it.
- Request logging owns exactly one completion event. Preserve status/byte
  recording and required optional response-writer interfaces.
- Metric route attributes use router patterns, never raw paths with IDs.
- Test changed limits, proxy trust, session security, CSRF rejection, rate refill
  and storage bounds, and panic-to-response behavior.

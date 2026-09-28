# HTTP middleware

- Scope middleware through route composition. Bound and sanitize request IDs
  for consistent error/event correlation.
- Resolve forwarded client identity only through configured trusted proxies
  before rate limiting. Reuse `meg/ratelimit` with bounded per-client storage;
  avoid custom time-window arithmetic.
- Preserve body/media limits, panic recovery, and no-store headers for private
  responses, using the shared API envelope for failures.
- SCS persists hashed session tokens with bounded idle/lifetime settings.
  Cookies remain HttpOnly/SameSite and Secure by default.
- nosurf and SCS must receive the same cookie-security setting.
  `SEREGA_COOKIE_SECURE=false` is an explicit local HTTP allowance.
  Cookie-authenticated unsafe methods require CSRF protection.
- Request logging owns the completion event. Preserve status/byte accounting
  and optional response-writer interfaces. Use router patterns for metric
  routes and follow `../observability/AGENTS.md` for telemetry behavior.
- Cover changes to limits, proxy trust, session security, CSRF rejection,
  rate refill/storage bounds, and panic-to-response behavior in focused tests.

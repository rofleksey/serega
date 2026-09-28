# Explicit environment configuration

Follow Akio's Load → FromLookup pattern to test parsing without mutating the
process environment.

- Supported variables belong here; do not scatter os.Getenv through services.
- Validate addresses, durations, booleans, log modes, and proxy CIDRs at startup.
- Require the database for serving; migration configuration may intentionally
  have fewer requirements.
- Secure cookies default on. An explicit local HTTP allowance must be supplied
  consistently to both session and CSRF middleware.
- Forwarded identity requires configured trusted proxies.
- Never include secret-bearing values or database URLs in errors or config dumps.
- Keep shared config contracts in types.go and private parsers near their use.
- Add meaningful cases for defaults, malformed values, and security-relevant
  combinations. Update README/deployment examples with configuration changes.

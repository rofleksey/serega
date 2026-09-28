# Shared frontend infrastructure

This is the lowest frontend layer; it cannot import app/pages/widgets/features/entities.

- api owns generated transport types, typed request construction, common errors,
  CSRF caching, same-origin credentials, and session-expiry signaling.
- ui owns domain-independent presentational primitives.
- lib owns small reusable mechanics with clear semantics, not a miscellaneous
  utility framework.
- Share code only when behavior and ownership are actually shared.
- Keep names/types explicit and preserve useful failure states.
- Test reusable behavior here so feature tests can focus on domain outcomes.

# User actions and forms

Features own cohesive actions, not routes or application-wide composition.

- Import entities/shared and related feature-local code, never widgets/pages/app.
- Substantial forms use React Hook Form + Zod with schemas/defaults beside the
  feature. Keep field and general server errors visible.
- Use Query mutations and invalidate owned data on success.
- Keep local draft state local. Avoid copying all server data into another store.
- Disable repeated/contradictory actions while pending and preserve user input
  after recoverable failures.
- Reuse shared confirmation and async/error primitives.
- Test real form validation/submission/error behaviors.

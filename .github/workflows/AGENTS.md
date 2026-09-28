# CI for a reusable public template

Preserve reproducible Akio checks while removing dependencies on its product.

- Work from a clean checkout: install locked frontend dependencies, generate
  Go/SQLC/TypeScript/assets, then run explicit checks/tests/build.
- Keep Go/Node/linter versions aligned with repository declarations.
- Pin third-party actions to reviewed commits and grant minimum permissions.
- Include meaningful unit, race, and disposable integration coverage where the
  workflow supports their infrastructure.
- No private Akio-Agent checkout, protocol secret, production credential, or
  external account is required to validate this template.
- Keep image publication/deployment separate from verification and deliberate.
  Do not make a new template consumer publish to a previous owner's registry.
- Never mask failed checks with broad continue-on-error or excluded packages.
- Generated files remain ignored; validate source regeneration rather than
  committing artifacts to satisfy CI.

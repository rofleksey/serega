# Canonical HTTP contract

openapi.yaml owns HTTP transport shapes. This directory contains contracts and
generator configuration, not a Go package.

- Define methods, operation IDs, parameters, required fields, successful
  responses, and error outcomes in the contract.
- Stable operation IDs generate Go handler methods and typed client paths.
- Reuse user/card/error schemas. Make optional, nullable, bounded text, enums,
  and concurrency preconditions intentional.
- Never return passwords or password hashes in user schemas.
- Include applicable authentication, CSRF, validation, missing-target, stale
  version, and internal-failure outcomes.
- OpenAPI security declarations document policy; handwritten middleware and
  handlers still implement it.
- Regenerate through make generate after contract changes. Go directives live
  in internal/api/generate.go; web's generate script owns TypeScript output.
- Fix the source contract, not generated output. Keep validation, service input,
  frontend client, and tests aligned.
- Do not host the raw schema or Swagger UI without a product requirement.

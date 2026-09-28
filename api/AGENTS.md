# HTTP contract

`openapi.yaml` and the generator configuration define the public HTTP contract;
this directory is not a Go package.

- Define methods, stable operation IDs, parameters, required fields, successful
  responses, and error outcomes. Operation IDs generate Go handler methods and
  TypeScript operation types.
- Reuse user/card/error schemas. Make optional/nullable values, text bounds,
  enums, and expected-version preconditions intentional.
- Declare authentication, CSRF, validation, missing-target, stale-version, and
  internal-failure outcomes where applicable. Security declarations document
  policy; middleware and handlers enforce it.
- Keep service inputs, handwritten validation, client behavior, and tests aligned
  with contract changes. Go generator directives live in `internal/api/generate.go`;
  `web/package.json` owns TypeScript generation.
- Do not serve raw OpenAPI or Swagger unless a product requirement calls for it.

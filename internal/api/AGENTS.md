# Handwritten contract boundary

This package owns generation directives, request validation, and shared
response encoding around generated bindings.

- Change ../../api/openapi.yaml for transport shapes and regenerate.
- generated/ contains generator-owned files only. Put wrappers, policy, tests,
  and documentation in this handwritten package.
- Preserve one error envelope with a stable code, safe readable message,
  request correlation, and optional field-level errors.
- Never expose raw database/decoder errors as client messages.
- Contract validation does not establish authentication, authorization,
  existence, or current version; services and handlers still own those.
- Keep schema-loading glue out of individual endpoints.
- Test malformed/unsupported request shapes and envelope consistency when
  changing adapters. Regenerate both Go and TypeScript consumers.

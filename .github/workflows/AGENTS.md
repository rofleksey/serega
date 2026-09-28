# CI workflows

- Read Go/Node versions from repository declarations and keep the linter version
  aligned with local tooling. Install frontend dependencies with `npm ci`.
- Pin third-party actions to reviewed commits, use minimal permissions, and
  disable persisted checkout credentials. Keep job timeouts and concurrency
  cancellation explicit.
- CI must validate this repository without private peer checkouts or external
  accounts. Integration jobs need a working Docker daemon for Testcontainers.
- Preserve explicit integration-tag linting and race coverage in addition to
  the root verification gate. Never mask failures with broad `continue-on-error`
  settings or package exclusions.
- Publishing requires a separately configured destination and workflow; routine
  verification must not push images or deploy.

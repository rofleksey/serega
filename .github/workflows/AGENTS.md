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
- The verify job builds the production image and runs the container smoke test
  (`test/container/smoke.py`) against it before any publish. Keep the smoke data
  synthetic and disposable.
- Publishing is a gated job in this workflow: it depends on every verification
  job, runs only for `main` pushes or explicit dispatch, and publishes only
  verified images. Pull requests and verification steps never push images or deploy.

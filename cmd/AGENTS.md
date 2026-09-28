# Executable boundary

Keep cmd/main.go as the tiny process entry point inherited from Akio.
Root rules apply.

- Delegate command parsing/execution to internal/cli.
- Keep exit status and process-level output here; do not open databases, wire
  routes, hash passwords, or implement policy.
- New commands belong in cli; new application behavior belongs in a service
  wired by app.
- Build via the root target so the executable contains current frontend assets.
- Verify arguments, exit codes, standard input, and signals through the built
  executable when those public contracts change.

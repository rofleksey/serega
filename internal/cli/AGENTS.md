# Cobra commands

Preserve Akio's small descriptive command files and constructor pattern.

- Cobra owns argument/flag validation and command help.
- Parse input, call app/application services, and present results. No direct
  store methods, SQLC calls, or SQL.
- serve loads serving config and uses a signal-cancellable context.
- migrate up explicitly invokes schema lifecycle; serve does not migrate.
- user create provisions an independent account. It must not delete or reset
  other users or invalidate their sessions.
- Read passwords from standard input as documented. Never add password flags,
  print supplied credentials, or include them in logs/errors.
- Use Cobra's input/output streams where practical to keep behavior testable.
- Return errors to one process-level boundary; avoid nested os.Exit and
  duplicate usage/error printing.
- Cover changed flags, standard input, validation, and exit behavior.

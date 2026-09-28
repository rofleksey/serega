# Production image and Compose

- Preserve the multistage build: pinned Node builds Vite, pinned Go generates
  bindings and embeds those assets, a minimal runtime runs as a non-root user.
- Keep production assets immutable and build from a clean checkout.
- Align module/toolchain, lockfile, binary path, image name, user, and command.
- Compose uses explicit required secret variables with no real credentials.
- Database health precedes explicit migration; serving follows migration success.
- Keep database data persistent, app filesystem read-only where practical,
  capabilities/privileges limited, and exposure intentional.
- Cookies default Secure; document TLS termination and trusted proxy CIDRs.
  Do not use local HTTP allowances as production defaults.
- Never add deployment/push side effects to routine local verification.

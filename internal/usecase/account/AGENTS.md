# Independent accounts

Adapt Akio's account service/password design to independent users.

- Own username/password validation, hashing/checking, authentication, provisioning.
- Preserve Argon2id, random salts, strict parameter bounds, and constant-time
  comparison. Store only hashes.
- Never log or return supplied credentials/hashes. Login errors must not
  disclose whether a username exists.
- Creating a user does not replace others. Never introduce SetSoleAdministrator,
  delete-other-users, or all-users session invalidation.
- Password/session operations affect the intended account only.
- All board users have equal permissions; do not infer an admin role from
  account order or inherit administrator naming accidentally.
- Keep the Store port narrow; hashing/validation stay outside PostgreSQL.
- Test credential failures, password rules, invalid/duplicate identities, and
  multi-user independence at the appropriate boundary.

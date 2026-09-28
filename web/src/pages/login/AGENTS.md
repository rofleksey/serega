# Login route

Preserve the existing form/validation/session-provider pattern.

- Submit through shared/api, with the CSRF/session machinery it owns.
- Passwords stay in form memory and are never logged, cached, or stored locally.
- Show safe authentication errors without revealing account existence.
- Keep pending submission and field labels clear; support normal keyboard submit.
- Refresh current-user state after login and honor safe intended navigation.
- Explain session expiry separately from invalid credentials.
- Do not add public self-registration to a CLI-provisioned application by accident.

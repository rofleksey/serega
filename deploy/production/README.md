# Container build

Build from the repository root:

```sh
docker build -f deploy/production/Dockerfile -t serega .
```

The multi-stage build generates the API and SQLC bindings and embeds the browser
assets into one non-root Go image. No production destination or automated image
publishing is configured in this template.

The example compose file binds HTTP to loopback for an HTTPS reverse proxy.
Supply your own PostgreSQL URL, explicitly run `serega migrate up`, provision
users with `serega user create`, and keep secure cookies enabled. Configure
`SEREGA_TRUSTED_PROXY_CIDRS` only for trusted proxy networks. Back up PostgreSQL;
the application image is stateless. Never use the development database password
on an exposed database. Review network, TLS, backup, and secret management for
the actual deployment before running it.

# syntax=docker/dockerfile:1
# The UI is built into the same immutable image that serves it.  Serega
# deliberately has no separate development web server at runtime.
FROM node:24.21.0-bookworm-slim AS web-build
WORKDIR /src
COPY api/openapi.yaml ./api/openapi.yaml
COPY web/package.json web/package-lock.json ./web/
RUN npm --prefix web ci
COPY web/ ./web/
RUN npm --prefix web run generate && npm --prefix web run build

FROM golang:1.27.1-bookworm AS go-build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN go generate ./internal/api ./internal/store/postgres/sqlc
COPY --from=web-build /src/internal/frontend/dist ./internal/frontend/dist
RUN CGO_ENABLED=0 go build -trimpath -ldflags='-s -w' -o /out/serega ./cmd

FROM debian:bookworm-slim
RUN apt-get update \
    && apt-get install -y --no-install-recommends ca-certificates \
    && rm -rf /var/lib/apt/lists/*
RUN useradd --system --uid 10001 --home-dir /nonexistent --shell /usr/sbin/nologin serega
WORKDIR /opt/serega
COPY --from=go-build /out/serega /usr/local/bin/serega
USER 10001:10001
EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/serega"]
CMD ["serve"]

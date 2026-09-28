.PHONY: format frontend generate lint check test test-integration test-integration-postgres test-integration-application build release

GO_FILES := $(shell find cmd internal test -name '*.go' -not -path '*/generated/*' -not -name '*.gen.go')
GO_PACKAGES := ./cmd/... ./internal/...

format:
	golangci-lint fmt

frontend:
	cd web && npm run generate && npm run build

generate:
	go generate ./internal/api ./internal/store/postgres/sqlc
	go tool sqlc vet -f internal/store/postgres/sqlc/sqlc.yaml
	$(MAKE) frontend

lint: generate
	golangci-lint run
	golangci-lint run --build-tags=integration ./test/integration/...

check: lint
	gofmt -d $(GO_FILES) | (! grep .)
	go tool sqlc vet -f internal/store/postgres/sqlc/sqlc.yaml
	go vet $(GO_PACKAGES)
	go list $(GO_PACKAGES) | grep -v '/internal/api/generated$$' | sed 's|^github.com/rofleksey/serega/|./|' | xargs go tool staticcheck
	cd web && npm run check

test: generate
	go test $(GO_PACKAGES)
	cd web && npm test

test-integration: build
	go test -tags=integration -count=1 ./test/integration/...

test-integration-postgres: generate
	go test -tags=integration -count=1 ./test/integration/postgres

test-integration-application: build
	go test -tags=integration -count=1 ./test/integration/application

build: generate
	go build -trimpath -o bin/serega ./cmd

release: check test build

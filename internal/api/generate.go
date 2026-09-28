package api

//go:generate go tool oapi-codegen --config ../../api/oapi-codegen.yaml -o generated/models.gen.go ../../api/openapi.yaml
//go:generate go tool oapi-codegen --config ../../api/oapi-codegen-server.yaml -o generated/server.gen.go ../../api/openapi.yaml

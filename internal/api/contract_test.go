package api

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestOpenAPIReferencesAndOperationIDs(t *testing.T) {
	raw, err := os.ReadFile(filepath.Join("..", "..", "api", "openapi.yaml"))
	if err != nil {
		t.Fatal(err)
	}

	var document map[string]any
	if err := yaml.Unmarshal(raw, &document); err != nil {
		t.Fatalf("parse OpenAPI YAML: %v", err)
	}

	operations := make(map[string]string)

	walkContract(document, func(reference string) {
		if !strings.HasPrefix(reference, "#/") {
			return
		}

		current := any(document)
		for _, part := range strings.Split(strings.TrimPrefix(reference, "#/"), "/") {
			object, ok := current.(map[string]any)
			if !ok {
				t.Errorf("reference %s traverses a non-object", reference)
				return
			}

			current, ok = object[part]
			if !ok {
				t.Errorf("unresolved reference %s", reference)
				return
			}
		}
	})

	for path, rawPath := range contractObject(t, document["paths"], "paths") {
		for method, rawOperation := range contractObject(t, rawPath, path) {
			operation, ok := rawOperation.(map[string]any)
			if !ok {
				continue
			}

			operationID, ok := operation["operationId"].(string)
			if !ok {
				continue
			}

			if previous, duplicate := operations[operationID]; duplicate {
				t.Errorf("duplicate operationId %s at %s and %s %s", operationID, previous, method, path)
			}

			operations[operationID] = method + " " + path
		}
	}
}

func TestGeneratedBindingsContainCurrentContract(t *testing.T) {
	root := filepath.Join("..", "..")
	for path, required := range map[string][]string{
		filepath.Join("generated", "server.gen.go"): {
			`http.MethodPatch+" "+options.BaseURL+"/v1/cards/{cardId}"`,
			"func GetSpec()",
		},
		filepath.Join("generated", "models.gen.go"): {
			`json:"createdBy"`, `json:"version"`,
		},
		filepath.Join(root, "web", "src", "shared", "api", "generated.ts"): {
			"createdBy:", "version: number",
		},
	} {
		raw, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("read generated binding %s: %v; run make generate", path, err)
		}

		for _, field := range required {
			if !strings.Contains(string(raw), field) {
				t.Errorf("%s is missing %q", path, field)
			}
		}
	}
}

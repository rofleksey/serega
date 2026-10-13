package usecase_test

import (
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"
	"strings"
	"testing"
)

// Use cases own application policy and declare consumer-owned ports; they must
// not depend on transports, concrete adapters, or persistence mechanics.
func TestUseCasesDoNotDependOnOuterLayers(t *testing.T) {
	err := filepath.WalkDir(".", func(path string, entry fs.DirEntry, err error) error {
		if err != nil || entry.IsDir() || !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return err
		}

		file, err := parser.ParseFile(token.NewFileSet(), path, nil, parser.ImportsOnly)
		if err != nil {
			return err
		}

		for _, spec := range file.Imports {
			dependency := strings.Trim(spec.Path.Value, `"`)
			if strings.Contains(dependency, "/internal/app") ||
				strings.Contains(dependency, "/internal/cli") ||
				strings.Contains(dependency, "/internal/handler") ||
				strings.Contains(dependency, "/internal/frontend") ||
				strings.Contains(dependency, "/internal/middleware") ||
				strings.Contains(dependency, "/internal/api/generated") ||
				strings.Contains(dependency, "/internal/store/") ||
				strings.Contains(dependency, "/internal/database/") ||
				strings.Contains(dependency, "github.com/jackc/pgx") {
				t.Errorf("%s imports outer layer %s", path, dependency)
			}
		}

		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

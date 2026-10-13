package handler

import (
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"
	"strings"
	"testing"
)

// Handlers decode transport input and call application services; persistence
// stays behind consumer-owned use-case ports and the composition root.
func TestHandlersDoNotDependOnPersistence(t *testing.T) {
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
			if strings.Contains(dependency, "/internal/store/") ||
				strings.Contains(dependency, "/internal/database/") ||
				strings.Contains(dependency, "github.com/jackc/pgx") {
				t.Errorf("%s imports persistence layer %s", path, dependency)
			}
		}

		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

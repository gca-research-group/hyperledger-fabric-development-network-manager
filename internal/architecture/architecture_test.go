package architecture

import (
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

func TestDependencyDirection(t *testing.T) {
	root := filepath.Clean(filepath.Join("..", ".."))
	err := filepath.WalkDir(filepath.Join(root, "internal"), func(path string, entry os.DirEntry, err error) error {
		if err != nil || entry.IsDir() || !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return err
		}
		file, err := parser.ParseFile(token.NewFileSet(), path, nil, parser.ImportsOnly)
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(filepath.Join(root, "internal"), path)
		top := strings.Split(filepath.ToSlash(rel), "/")[0]
		for _, imported := range file.Imports {
			name, _ := strconv.Unquote(imported.Path.Value)
			if top != "cli" && strings.Contains(name, "/internal/cli") {
				t.Errorf("%s imports CLI package %s", rel, name)
			}
			if top != "application" && top != "cli" && strings.Contains(name, "/internal/application") {
				t.Errorf("infrastructure package %s imports application package %s", rel, name)
			}
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

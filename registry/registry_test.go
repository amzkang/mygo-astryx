package registry_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"strings"
	"testing"

	"github.com/amzkang/mygo-astryx/ds"
	"github.com/amzkang/mygo-astryx/registry"
)

func TestRegistryMatchesExportedDSComponents(t *testing.T) {
	exported := map[string]bool{}
	fset := token.NewFileSet()
	err := fs.WalkDir(ds.SourceFS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return err
		}
		src, err := ds.SourceFS.ReadFile(path)
		if err != nil {
			return err
		}
		f, err := parser.ParseFile(fset, path, src, 0)
		if err != nil {
			return err
		}
		// Ensure swizzle-safe imports only.
		for _, imp := range f.Imports {
			p := strings.Trim(imp.Path.Value, `"`)
			if strings.HasPrefix(p, "github.com/") &&
				p != "github.com/egoist/mygo/ui" &&
				p != "github.com/amzkang/mygo-astryx/theme" &&
				p != "github.com/amzkang/mygo-astryx/icons" {
				t.Errorf("%s imports disallowed package %q (breaks swizzle contract)", path, p)
			}
		}
		for _, decl := range f.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || fn.Recv != nil || !fn.Name.IsExported() {
				continue
			}
			exported[fn.Name.Name] = true
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	documented := map[string]bool{}
	for _, c := range registry.Components {
		documented[c.Name] = true
		if len(c.SourceFiles) == 0 {
			t.Errorf("component %q has no SourceFiles", c.Name)
		}
		for _, sf := range c.SourceFiles {
			if _, err := fs.Stat(ds.SourceFS, sf); err != nil {
				t.Errorf("component %q references missing source file %q: %v", c.Name, sf, err)
			}
		}
		if c.Example == "" || c.Summary == "" || c.A11y == "" {
			t.Errorf("component %q is missing Example, Summary, or A11y", c.Name)
		}
	}

	for name := range exported {
		if !documented[name] {
			t.Errorf("exported ds.%s is missing from registry.Components", name)
		}
	}
	for name := range documented {
		if !exported[name] {
			t.Errorf("registry.Components lists %s which is not an exported ds function", name)
		}
	}
}

func TestSearchAndPresetsAndTemplates(t *testing.T) {
	if len(registry.Presets) != 7 {
		t.Fatalf("expected 7 presets, got %d", len(registry.Presets))
	}
	if len(registry.Templates) < 2 {
		t.Fatalf("expected >=2 templates, got %d", len(registry.Templates))
	}
	hits := registry.Search("button", "")
	if len(hits) == 0 || hits[0].Name != "Button" {
		t.Fatalf("expected Button as top hit for 'button', got %+v", hits)
	}
	if _, ok := registry.FindComponent("CommandPalette"); !ok {
		t.Fatal("expected FindComponent(CommandPalette) to succeed")
	}
}

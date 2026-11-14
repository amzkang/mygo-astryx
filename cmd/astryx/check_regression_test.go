package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestStripSourceFSEmbedRemovesEmptyImports(t *testing.T) {
	for _, source := range []string{
		"package sample\nimport (\"embed\")\nvar SourceFS embed.FS\n",
		"package sample\nimport \"embed\"\nvar (SourceFS embed.FS; Keep = 1)\n",
		"package sample\nimport (\"embed\"; \"fmt\")\nvar SourceFS embed.FS\nvar Keep = fmt.Sprint(1)\n",
	} {
		fset := token.NewFileSet()
		file, err := parser.ParseFile(fset, "sample.go", source, 0)
		if err != nil {
			t.Fatal(err)
		}
		stripSourceFSEmbed(file)
		for _, decl := range file.Decls {
			if gen, ok := decl.(*ast.GenDecl); ok && len(gen.Specs) == 0 {
				t.Fatal("ejection retained an empty declaration")
			}
		}
		for _, imp := range file.Imports {
			if imp.Path.Value == `"embed"` {
				t.Fatal("ejection retained the embed import in AST metadata")
			}
		}
		var output bytes.Buffer
		if err := format.Node(&output, fset, file); err != nil {
			t.Fatal(err)
		}
		if strings.Contains(output.String(), "SourceFS") {
			t.Fatal("ejection retained SourceFS")
		}
		if strings.Contains(source, "fmt") && !strings.Contains(output.String(), `"fmt"`) {
			t.Fatal("ejection removed an unrelated import")
		}
		if _, err := parser.ParseFile(token.NewFileSet(), "ejected.go", output.Bytes(), 0); err != nil {
			t.Fatalf("ejected source is invalid: %v\n%s", err, output.String())
		}
	}
}

func TestCheckDirectoryKeysSegmentedAndOverlayControls(t *testing.T) {
	for _, call := range []string{
		`ds.Link(%s, "Docs", "https://example.com")`,
		`ui.Link(%s, "Docs", "https://example.com")`,
		"ds.SideNav(%s, &nav, items)",
		"ds.MobileNav(%s, &open, func() {})",
		"ds.AppShell(%s, &shell, slots)",
		"ds.ScrollableArea(%s, func() {})",
		"ds.MoreMenu(%s, &open, nil)",
		"ds.SegmentedControl(%s, &selected, labels)",
		"ui.SegmentedBase(%s, &selected, len(labels))",
		"ui.DateInput(%s, &date)",
		"ui.PopoverBase(%s, anchor, &open, func(panel ui.Element) {})",
		"ui.TooltipBase(%s, anchor, func(tip ui.Element) {})",
		"ui.ToastViewportBase(%s, func(viewport ui.Element, toasts []ui.Toast) {})",
		"ds.ContextMenu(%s, anchor, &open, nil)",
		"ds.ModalCard(%s, &open, \"Dialog\", 320, func() {})",
		"ds.Drawer(%s, &open, func() {})",
		"ds.Popover(%s, anchor, &open, func() {})",
		"ds.Tooltip(%s, anchor, \"Help\")",
		"ds.ToastViewport(%s)",
	} {
		t.Run(strings.Split(call, "(")[0], func(t *testing.T) {
			for _, keyed := range []bool{false, true} {
				context := "c"
				if keyed {
					context = "c.Key(i)"
				}
				source := `package sample
import (
    "time"
    "github.com/egoist/mygo/ui"
    "github.com/amzkang/mygo-astryx/ds"
)
func View(c *ui.Context, anchor ui.Element) {
    var selected int
    var date time.Time
    open := false
    labels := []string{"A", "B"}
    for i := range labels {
        ` + fmt.Sprintf(call, context) + `
    }
}
`
				dir := t.TempDir()
				if err := os.WriteFile(filepath.Join(dir, "sample.go"), []byte(source), 0o644); err != nil {
					t.Fatal(err)
				}
				issues, err := CheckDirectory(dir)
				if err != nil {
					t.Fatal(err)
				}
				if keyed && len(issues) != 0 {
					t.Fatalf("keyed control reported issues: %+v", issues)
				}
				if !keyed && (len(issues) != 1 || issues[0].Rule != "missing-key-in-loop") {
					t.Fatalf("unkeyed control must report one missing-key issue: %+v", issues)
				}
			}
		})
	}
}

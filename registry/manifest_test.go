package registry

import (
	"bytes"
	"encoding/json"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"
)

type inventoryFixture struct {
	Repository string            `json:"repository"`
	Commit     string            `json:"commit"`
	Paths      map[string]string `json:"paths"`
	Entries    []struct {
		Kind       string `json:"kind"`
		SourcePath string `json:"sourcePath"`
		TreeSHA    string `json:"treeSHA"`
	} `json:"entries"`
}

func readInventory(t *testing.T) inventoryFixture {
	t.Helper()
	data, err := os.ReadFile("testdata/astryx-inventory.json")
	if err != nil {
		t.Fatal(err)
	}
	var fixture inventoryFixture
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&fixture); err != nil {
		t.Fatal(err)
	}
	return fixture
}

func matchInventory(catalog parityCatalog, fixture inventoryFixture) bool {
	if catalog.Baseline.Repository != fixture.Repository || catalog.Baseline.Commit != fixture.Commit || !reflect.DeepEqual(catalog.Baseline.Paths, fixture.Paths) {
		return false
	}
	actual := map[string]bool{}
	for _, entry := range catalog.Entries {
		actual[entry.Kind+":"+entry.SourcePath] = true
	}
	expected := map[string]bool{}
	for _, entry := range fixture.Entries {
		expected[entry.Kind+":"+entry.SourcePath] = true
	}
	return len(actual) == len(catalog.Entries) && len(expected) == len(fixture.Entries) && reflect.DeepEqual(actual, expected)
}

func TestParityMatchesFixedInventory(t *testing.T) {
	catalog, err := decodeParity(parityJSON)
	if err != nil {
		t.Fatal(err)
	}
	fixture := readInventory(t)
	if !matchInventory(catalog, fixture) {
		t.Fatal("parity mapping differs from the fixed upstream directory inventory")
	}
	counts := map[string]int{}
	for _, entry := range fixture.Entries {
		counts[entry.Kind]++
		if len(entry.TreeSHA) != 40 {
			t.Errorf("%s has no fixed directory tree SHA", entry.SourcePath)
		}
	}
	if !reflect.DeepEqual(counts, map[string]int{"core": 124, "block": 172, "page": 58}) {
		t.Fatalf("unexpected inventory counts %v", counts)
	}
	for _, entry := range catalog.Entries {
		for _, evidence := range entry.Evidence {
			if strings.HasPrefix(evidence, "https://") {
				if !strings.Contains(evidence, upstreamCommit) {
					t.Errorf("%s cites an unpinned upstream source %s", entry.SourcePath, evidence)
				}
				continue
			}
			if _, err := os.Stat(filepath.Join("..", evidence)); err != nil {
				t.Errorf("%s cites missing local evidence %s", entry.SourcePath, evidence)
			}
		}
	}
	categories := map[string]int{}
	for _, entry := range catalog.Entries {
		if entry.Kind == "core" {
			categories[entry.Category]++
		}
		if entry.SourcePath == "packages/core/src/Stepper" {
			for _, ref := range entry.LocalRefs {
				if ref.Kind == "component" && ref.Name == "Stepper" {
					t.Fatal("workflow Stepper must not map to numeric Stepper")
				}
			}
		}
	}
	if !reflect.DeepEqual(categories, map[string]int{"component": 107, "support": 12, "internal": 5}) {
		t.Fatalf("unexpected core categories %v", categories)
	}
}

func TestParityRejectsInvalidMappings(t *testing.T) {
	for _, tc := range []struct {
		name   string
		mutate func(*parityCatalog)
	}{
		{"missing entry", func(c *parityCatalog) { c.Entries = c.Entries[1:] }},
		{"duplicate entry", func(c *parityCatalog) { c.Entries = append(c.Entries, c.Entries[0]) }},
		{"unknown status", func(c *parityCatalog) { c.Entries[0].Status = "done" }},
		{"unknown core category", func(c *parityCatalog) {
			for i := range c.Entries {
				if c.Entries[i].Kind == "core" {
					c.Entries[i].Category = "unknown"
					break
				}
			}
		}},
		{"missing reason", func(c *parityCatalog) { c.Entries[0].Reason = " " }},
		{"missing evidence", func(c *parityCatalog) { c.Entries[0].Evidence = nil }},
		{"dangling reference", func(c *parityCatalog) {
			c.Entries[0].LocalRefs = []ResourceRef{{Kind: "component", Name: "Unknown"}}
		}},
		{"implemented without reference", func(c *parityCatalog) {
			c.Entries[0].Status = "implemented"
			c.Entries[0].LocalRefs = []ResourceRef{}
		}},
		{"wrong commit", func(c *parityCatalog) { c.Baseline.Commit = strings.Repeat("0", 40) }},
		{"wrong MyGo version", func(c *parityCatalog) { c.Baseline.MyGoVersion = "v0.0.0" }},
		{"wrong source path", func(c *parityCatalog) { c.Baseline.Paths["core"] = "wrong" }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			catalog, err := decodeParity(parityJSON)
			if err != nil {
				t.Fatal(err)
			}
			tc.mutate(&catalog)
			if validateParity(catalog) == nil && matchInventory(catalog, readInventory(t)) {
				t.Fatal("invalid mapping passed validation")
			}
		})
	}
	bad := bytes.Replace(parityJSON, []byte(`"baseline"`), []byte(`"basline"`), 1)
	if _, err := decodeParity(bad); err == nil {
		t.Fatal("unknown JSON fields must fail strict decoding")
	}
	if _, err := decodeParity(append(append([]byte{}, parityJSON...), []byte(`{}`)...)); err == nil {
		t.Fatal("trailing JSON must fail decoding")
	}
}

func TestParityUsesActualSegmentAndTagPrimitivesAndIconCatalog(t *testing.T) {
	m := BuildManifest()
	entries := map[string]ParityEntry{}
	for _, entry := range m.Parity {
		entries[entry.SourcePath] = entry
	}
	segmented := entries["packages/core/src/SegmentedControl"]
	if !reflect.DeepEqual(segmented.Primitives, []PrimitiveRef{{Package: "github.com/egoist/mygo/ui", Symbol: "SegmentedBase"}}) {
		t.Fatal("SegmentedControl must identify the SegmentedBase used by its implementation")
	}
	tag := entries["packages/core/src/Token"]
	if !reflect.DeepEqual(tag.Primitives, []PrimitiveRef{
		{Package: "github.com/egoist/mygo/ui", Symbol: "Row"},
		{Package: "github.com/egoist/mygo/ui", Symbol: "Text"},
		{Package: "github.com/egoist/mygo/ui", Symbol: "ButtonBase"},
	}) {
		t.Fatal("Token must identify display tag primitives")
	}
	refs := []ResourceRef{}
	for _, icon := range m.Icons {
		refs = append(refs, ResourceRef{Kind: "icon", Name: icon.Name})
	}
	if !reflect.DeepEqual(entries["packages/core/src/Icon"].LocalRefs, refs) {
		t.Fatal("Icon coverage must enumerate the actual local icon catalog")
	}
}

func TestManifestDeterministicAndSorted(t *testing.T) {
	before, err := json.Marshal(Components)
	if err != nil {
		t.Fatal(err)
	}
	m := BuildManifest()
	first, err := json.Marshal(m)
	if err != nil {
		t.Fatal(err)
	}
	second, err := json.Marshal(BuildManifest())
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(first, second) || m.SchemaVersion != 1 {
		t.Fatal("manifest must be deterministic and use schema version 1")
	}
	after, _ := json.Marshal(Components)
	if !bytes.Equal(before, after) {
		t.Fatal("manifest sorting changed registry components")
	}
	var object map[string]json.RawMessage
	if err := json.Unmarshal(first, &object); err != nil {
		t.Fatal(err)
	}
	for field, expected := range map[string]int{
		"components": len(Components), "templates": len(Templates), "tokens": len(Tokens),
		"presets": len(Presets), "guides": len(Guides), "icons": len(IconNames()), "parity": 354,
	} {
		var rows []map[string]json.RawMessage
		if err := json.Unmarshal(object[field], &rows); err != nil || rows == nil {
			t.Fatalf("%s must be a JSON array: %v", field, err)
		}
		if len(rows) != expected {
			t.Fatalf("%s has %d records, want %d", field, len(rows), expected)
		}
		var keys []string
		for _, row := range rows {
			key := row["name"]
			if field == "tokens" {
				key = row["path"]
			}
			if field == "parity" {
				key = append(append([]byte{}, row["kind"]...), row["sourcePath"]...)
			}
			keys = append(keys, string(key))
		}
		if !sort.StringsAreSorted(keys) {
			t.Errorf("%s is not sorted", field)
		}
	}
	for _, icon := range m.Icons {
		if icon.Name == "" || icon.Symbol == "" {
			t.Fatal("manifest icons require names and Go symbols")
		}
	}
}

func TestParityPrimitiveSymbolsExist(t *testing.T) {
	cmd := exec.Command(filepath.Join("..", "bin", "go"), "list", "-m", "-json", "github.com/egoist/mygo")
	output, err := cmd.Output()
	if err != nil {
		t.Fatal(err)
	}
	var module struct{ Dir, Version string }
	if err := json.Unmarshal(output, &module); err != nil {
		t.Fatal(err)
	}
	if module.Version != "v0.8.3" {
		t.Fatalf("parity primitives require MyGo v0.8.3, resolved %s", module.Version)
	}
	packages := map[string]map[string]bool{}
	for _, entry := range BuildManifest().Parity {
		for _, primitive := range entry.Primitives {
			if _, ok := packages[primitive.Package]; !ok {
				dir := filepath.Join(module.Dir, strings.TrimPrefix(primitive.Package, "github.com/egoist/mygo"))
				parsed, err := parser.ParseDir(token.NewFileSet(), dir, func(info os.FileInfo) bool {
					return !strings.HasSuffix(info.Name(), "_test.go")
				}, 0)
				if err != nil {
					t.Fatal(err)
				}
				symbols := map[string]bool{}
				for _, pkg := range parsed {
					for _, file := range pkg.Files {
						for _, decl := range file.Decls {
							if fn, ok := decl.(*ast.FuncDecl); ok && fn.Name.IsExported() {
								name := fn.Name.Name
								if fn.Recv != nil {
									receiver := fn.Recv.List[0].Type
									if pointer, ok := receiver.(*ast.StarExpr); ok {
										receiver = pointer.X
									}
									if ident, ok := receiver.(*ast.Ident); ok {
										name = ident.Name + "." + name
									}
								}
								symbols[name] = true
							}
						}
						for name := range file.Scope.Objects {
							symbols[name] = true
						}
					}
				}
				packages[primitive.Package] = symbols
			}
			if !packages[primitive.Package][primitive.Symbol] {
				t.Errorf("%s cites missing MyGo primitive %s.%s", entry.SourcePath, primitive.Package, primitive.Symbol)
			}
		}
	}
}

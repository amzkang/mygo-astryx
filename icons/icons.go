package icons

import (
	"embed"
	"io/fs"
	"path"
	"sort"
	"strings"

	"github.com/egoist/mygo/ui"
)

//go:embed svg/*.svg
var svgFS embed.FS

// IconEntry pairs an icon's kebab-case slug with its parsed *ui.SVG and Go symbol name.
type IconEntry struct {
	Name   string
	Symbol string
	SVG    *ui.SVG
}

var (
	registry = map[string]*ui.SVG{}
	entries  []IconEntry
)

func mustLoad(name, symbol string) *ui.SVG {
	data, err := fs.ReadFile(svgFS, path.Join("svg", name+".svg"))
	if err != nil {
		panic("icons: missing embedded svg " + name + ": " + err.Error())
	}
	svg := ui.MustParseSVG(data)
	registry[name] = svg
	entries = append(entries, IconEntry{Name: name, Symbol: symbol, SVG: svg})
	return svg
}

// ByName looks up a built-in icon by kebab-case name (e.g. "chevron-right").
func ByName(name string) (*ui.SVG, bool) {
	svg, ok := registry[strings.ToLower(strings.TrimSpace(name))]
	return svg, ok
}

// All returns all registered icons sorted alphabetically by name.
func All() []IconEntry {
	out := make([]IconEntry, len(entries))
	copy(out, entries)
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}

// Find returns icons whose name or Go symbol contains query (case-insensitive).
func Find(query string) []IconEntry {
	q := strings.ToLower(strings.TrimSpace(query))
	all := All()
	if q == "" {
		return all
	}
	var out []IconEntry
	for _, e := range all {
		if strings.Contains(e.Name, q) || strings.Contains(strings.ToLower(e.Symbol), q) {
			out = append(out, e)
		}
	}
	return out
}

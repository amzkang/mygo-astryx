package icons

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

func TestAllIconsParsedAndSearchable(t *testing.T) {
	all := All()
	if len(all) < 50 {
		t.Fatalf("expected at least 50 icons, got %d", len(all))
	}
	for _, entry := range all {
		if entry.SVG == nil {
			t.Fatalf("icon %q has nil SVG", entry.Name)
		}
		got, ok := ByName(entry.Name)
		if !ok || got != entry.SVG {
			t.Fatalf("ByName(%q) lookup mismatch", entry.Name)
		}
	}
	hits := Find("chevron")
	if len(hits) != 4 {
		t.Fatalf("expected 4 chevron hits, got %d", len(hits))
	}
	img := ui.Render(func(c *ui.Context) {
		ui.Row(c).Gap(8).Padding(8).Children(func() {
			ui.Icon(c, Check).Size(20, 20)
			ui.Icon(c, Sparkles).Size(20, 20)
		})
	}, 100, 48, 1)
	if img.Bounds().Dx() != 100 || img.Bounds().Dy() != 48 {
		t.Fatalf("unexpected render bounds: %v", img.Bounds())
	}
}

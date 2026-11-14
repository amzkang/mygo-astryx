package theme

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

func TestStockThemesContrastAndToUI(t *testing.T) {
	for _, ts := range StockThemes() {
		for _, dark := range []bool{false, true} {
			th := ts.For(dark)
			if ratio := ContrastRatio(th.Colors.Text, th.Colors.Background); ratio < 7.0 {
				t.Errorf("%s (dark=%v): text/bg contrast %.2f < 7.0", th.Name, dark, ratio)
			}
			for _, name := range []string{"accent", "danger", "warning", "success", "info"} {
				tone := th.Tone(name)
				if ratio := ContrastRatio(tone.OnColor, tone.Base); ratio < 4.5 {
					t.Errorf("%s (dark=%v) tone %s: OnColor/Base contrast %.2f < 4.5", th.Name, dark, name, ratio)
				}
				if ratio := ContrastRatio(tone.Text, th.Colors.Background); ratio < 3.0 {
					t.Errorf("%s (dark=%v) tone %s: Text/Bg contrast %.2f < 3.0", th.Name, dark, name, ratio)
				}
			}
			uit := th.ToUI()
			if uit.Dark != dark || uit.Background != th.Colors.Background || uit.Accent != th.Colors.Accent.Base {
				t.Errorf("%s (dark=%v): ToUI mismatch", th.Name, dark)
			}
		}
	}
}

func TestWithThemeScopeIsolation(t *testing.T) {
	var seenOuter, seenInner, seenAfter string
	view := func(c *ui.Context) {
		SetRoot(c, Neutral.Light)
		seenOuter = Of(c).Name
		WithTheme(c, Matcha.Dark, func() {
			seenInner = Of(c).Name
			ui.Text(c, "inside:"+Of(c).Name)
		})
		seenAfter = Of(c).Name
		ui.Text(c, "outside:"+Of(c).Name)
	}
	tt := ui.NewTester(view, 400, 200)
	if seenOuter != "neutral" || seenInner != "matcha" || seenAfter != "neutral" {
		t.Fatalf("unexpected theme scope values: outer=%q inner=%q after=%q", seenOuter, seenInner, seenAfter)
	}
	if !tt.HasText("inside:matcha") || !tt.HasText("outside:neutral") {
		t.Fatalf("missing expected scoped text: %q", tt.Texts())
	}
}

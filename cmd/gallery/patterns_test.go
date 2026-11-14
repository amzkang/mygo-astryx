package main

import (
	"testing"

	"github.com/amzkang/mygo-astryx/registry"
	"github.com/amzkang/mygo-astryx/theme"
	"github.com/egoist/mygo/ui"
)

func TestP1aGalleryPatternsAndPlaygrounds(t *testing.T) {
	for _, name := range patternNames {
		t.Run(name, func(t *testing.T) {
			app := newGalleryState()
			app.pattern = name
			for _, rtl := range []bool{false, true} {
				app.rtl = rtl
				tester := ui.NewTester(func(c *ui.Context) {
					theme.WithTheme(c, theme.Matcha.Light, func() {
						root := ui.Column(c).Fill()
						if rtl {
							root.Direction(ui.RightToLeft)
						}
						root.Children(func() { renderPatternsPage(c, app, theme.Matcha.Light) })
					}).Fill()
				}, 1100, 700)
				assertNonBlankImage(t, tester.Image(), 1100, 700, name)
				if name == "app-shell-narrow" {
					if err := tester.Click("Open navigation"); err != nil || !app.patternShell.MobileNavOpen {
						t.Fatalf("narrow scene must open its real drawer rtl=%v: %v", rtl, err)
					}
					if err := tester.Click("Design system"); err != nil || app.patternShell.MobileNavOpen || app.patternShell.Nav.SelectedID != "design" {
						t.Fatalf("narrow scene destination: %v, %+v", err, app.patternShell)
					}
				}
				if name == "action-overflow" {
					if err := tester.Click("More actions"); err != nil {
						t.Fatal(err)
					}
					if err := tester.Click("Export workspace"); err != nil || app.patternAction != "Export workspace" {
						t.Fatalf("scene action: %v", err)
					}
				}
			}
		})
	}
	for _, name := range []string{"Layout", "TopNav", "SideNav", "MobileNav", "AppShell", "ScrollableArea", "MoreMenu", "Grid", "AspectRatio", "Link"} {
		t.Run("playground-"+name, func(t *testing.T) {
			app := newGalleryState()
			comp, ok := registry.FindComponent(name)
			if !ok || generateSnippet(comp, app) == "" {
				t.Fatal("component needs metadata and a usable snippet")
			}
			img := ui.Render(func(c *ui.Context) {
				theme.WithTheme(c, theme.Neutral.Light, func() {
					renderP1aWidget(c, app, theme.Neutral.Light, name)
				})
			}, 800, 500, 1)
			assertNonBlankImage(t, img, 800, 500, name)
		})
	}
	app := newGalleryState()
	tester := ui.NewTester(func(c *ui.Context) { renderGallery(c, app) }, 1200, 820)
	if err := tester.Click("Patterns"); err != nil || app.page != pagePatterns || !tester.HasText("Application patterns") {
		t.Fatalf("gallery entry must navigate to patterns: %v", err)
	}
}

package main

import (
	"image"
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"
	"github.com/amzkang/mygo-astryx/registry"
	"github.com/amzkang/mygo-astryx/theme"
)

func TestGalleryNavigationThemeDarkDensityAndRTL(t *testing.T) {
	app := newGalleryState()
	tester := ui.NewTester(func(c *ui.Context) {
		renderGallery(c, app)
	}, 1200, 820)

	// Switch theme via command palette
	if err := tester.Click("Command Palette"); err != nil {
		t.Fatalf("failed to open command palette: %v", err)
	}
	tester.Type("gothic")
	tester.Key(0, ui.KeyEnter)
	if app.themeName != "gothic" {
		t.Fatalf("expected themeName=gothic, got %q", app.themeName)
	}

	// Toggle Dark and RTL switches
	if err := tester.Click("Dark"); err != nil || !app.dark {
		t.Fatalf("expected dark=true, err=%v", err)
	}
	if err := tester.Click("RTL"); err != nil || !app.rtl {
		t.Fatalf("expected rtl=true, err=%v", err)
	}

	// Navigate through all pages
	for _, nav := range []struct {
		label string
		want  galleryPage
	}{
		{"Components", pageComponents},
		{"Tokens & WCAG", pageTokens},
		{"Icons (56)", pageIcons},
		{"Overview", pageOverview},
	} {
		if err := tester.Click(nav.label); err != nil {
			t.Fatalf("click %q failed: %v", nav.label, err)
		}
		if app.page != nav.want {
			t.Fatalf("expected page %q, got %q", nav.want, app.page)
		}
	}

	// Render headless PNG snapshot sanity check
	img := ui.Render(func(c *ui.Context) {
		renderGallery(c, app)
	}, 1200, 820, 1)
	assertNonBlankImage(t, img, 1200, 820, "overview gothic dark")
}

func TestGalleryPlaygroundPropsAndSnippetGenerator(t *testing.T) {
	app := newGalleryState()
	app.page = pageComponents
	app.selectedCmp = "Button"

	tester := ui.NewTester(func(c *ui.Context) {
		renderGallery(c, app)
	}, 1280, 920)

	// Toggle playground controls via UI
	if err := tester.Click("Soft"); err != nil || app.pgVariantIdx != 1 {
		t.Fatalf("expected pgVariantIdx=1 after clicking Soft, err=%v got=%d", err, app.pgVariantIdx)
	}
	if err := tester.Click("LG"); err != nil || app.pgSizeIdx != 3 {
		t.Fatalf("expected pgSizeIdx=3 after clicking LG, err=%v got=%d", err, app.pgSizeIdx)
	}
	if err := tester.Click("Success"); err != nil || app.pgToneIdx != 2 {
		t.Fatalf("expected pgToneIdx=2 after clicking Success, err=%v got=%d", err, app.pgToneIdx)
	}
	if err := tester.Click("Loading"); err != nil || !app.pgLoading {
		t.Fatalf("expected pgLoading=true after clicking Loading, err=%v", err)
	}

	comp, ok := registry.FindComponent("Button")
	if !ok {
		t.Fatal("Button missing from registry")
	}
	snippet := generateSnippet(comp, app)
	for _, want := range []string{"ds.VariantSoft", "ds.ToneSuccess", "ds.SizeLG", "Loading: true"} {
		if !strings.Contains(snippet, want) {
			t.Fatalf("expected live snippet to contain %q, got:\n%s", want, snippet)
		}
	}

	// Focus SearchInput and type "Dropdown" so DropdownMenu scrolls into viewport, then click it
	if err := tester.Click("Search components (e.g. button, modal, input)..."); err != nil {
		t.Fatalf("failed to focus component search: %v", err)
	}
	tester.Type("Dropdown")
	if err := tester.Click("DropdownMenu"); err != nil || app.selectedCmp != "DropdownMenu" {
		t.Fatalf("expected selectedCmp=DropdownMenu, err=%v got=%q", err, app.selectedCmp)
	}
	ddComp, _ := registry.FindComponent("DropdownMenu")
	ddSnippet := generateSnippet(ddComp, app)
	if !strings.Contains(ddSnippet, "ds.DropdownMenu") || !strings.Contains(ddSnippet, "ds.VariantSoft") {
		t.Fatalf("expected live DropdownMenu snippet, got:\n%s", ddSnippet)
	}
}

func TestGalleryTokensContrastMatrixAndDesktopShowcase(t *testing.T) {
	app := newGalleryState()
	tester := ui.NewTester(func(c *ui.Context) {
		renderGallery(c, app)
	}, 1280, 920)

	// Switch to Desktop Controls tab on Overview
	if err := tester.Click("Desktop Controls"); err != nil || app.playTab != 2 {
		t.Fatalf("expected playTab=2 for Desktop Controls, err=%v got=%d", err, app.playTab)
	}
	if err := tester.Click("Actions"); err != nil || !app.pgMenuOpen {
		t.Fatalf("expected DropdownMenu pgMenuOpen=true, err=%v", err)
	}
	if err := tester.Click("Confirm Dialog"); err != nil || !app.pgAlertOpen {
		t.Fatalf("expected AlertDialog pgAlertOpen=true, err=%v", err)
	}
	if err := tester.Click("Cancel"); err != nil || app.pgAlertOpen {
		t.Fatalf("expected AlertDialog closed on Cancel, err=%v open=%v", err, app.pgAlertOpen)
	}

	// Navigate to Tokens & WCAG -> 7-Theme Matrix and switch theme to Y2K
	if err := tester.Click("Tokens & WCAG"); err != nil {
		t.Fatalf("failed to open Tokens & WCAG: %v", err)
	}
	if err := tester.Click("7-Theme Matrix"); err != nil || app.tokensTab != 1 {
		t.Fatalf("expected tokensTab=1, err=%v got=%d", err, app.tokensTab)
	}
	if err := tester.Click("Use Y2k"); err != nil || app.themeName != "y2k" {
		t.Fatalf("expected themeName=y2k after matrix click, err=%v got=%q", err, app.themeName)
	}
}

func TestGalleryRenderRegressionAllPagesAndThemes(t *testing.T) {
	pages := []galleryPage{pageOverview, pageComponents, pageTokens, pageIcons}
	themes := []string{"neutral", "matcha", "gothic", "y2k"}
	for _, thName := range themes {
		for _, dark := range []bool{false, true} {
			for _, pg := range pages {
				app := newGalleryState()
				app.themeName = thName
				app.dark = dark
				app.page = pg
				if pg == pageTokens {
					app.tokensTab = 1
				}
				if pg == pageOverview {
					app.playTab = 2
				}
				img := ui.Render(func(c *ui.Context) {
					renderGallery(c, app)
				}, 1180, 780, 1)
				assertNonBlankImage(t, img, 1180, 780, thName+" "+string(pg))
				// Ensure theme exists cleanly
				_ = theme.ByName(thName).For(dark)
			}
		}
	}
}

func assertNonBlankImage(t *testing.T, img *image.RGBA, w, h int, label string) {
	t.Helper()
	if img == nil || img.Bounds().Dx() != w || img.Bounds().Dy() != h {
		t.Fatalf("%s: unexpected bounds %v", label, img)
	}
	first := img.RGBAAt(0, 0)
	varied := false
	for y := 0; y < h; y += 16 {
		for x := 0; x < w; x += 16 {
			if img.RGBAAt(x, y) != first {
				varied = true
				break
			}
		}
		if varied {
			break
		}
	}
	if !varied {
		t.Fatalf("%s: rendered image appears blank/uniform (%v)", label, first)
	}
}

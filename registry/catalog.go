package registry

import (
	"fmt"

	"github.com/egoist/mygo/ui"
	"github.com/amzkang/mygo-astryx/icons"
	"github.com/amzkang/mygo-astryx/theme"
)

// Tokens documents all theme token paths in theme.Theme.
var Tokens = []TokenDoc{
	{Path: "Colors.Background", Group: "surface", Type: "ui.Color", Description: "Base window canvas color"},
	{Path: "Colors.Surface", Group: "surface", Type: "ui.Color", Description: "Default card and control surface"},
	{Path: "Colors.SurfaceSubtle", Group: "surface", Type: "ui.Color", Description: "Recessed wells, inputs, soft pills"},
	{Path: "Colors.SurfaceElevated", Group: "surface", Type: "ui.Color", Description: "Modals, popovers, drawers, elevated cards"},
	{Path: "Colors.SurfaceHover", Group: "surface", Type: "ui.Color", Description: "Hover background state"},
	{Path: "Colors.SurfacePressed", Group: "surface", Type: "ui.Color", Description: "Active pressed background state"},
	{Path: "Colors.Border", Group: "border", Type: "ui.Color", Description: "Default hairline border"},
	{Path: "Colors.BorderSubtle", Group: "border", Type: "ui.Color", Description: "Soft dividers and section boundaries"},
	{Path: "Colors.BorderStrong", Group: "border", Type: "ui.Color", Description: "High-emphasis control borders"},
	{Path: "Colors.Text", Group: "text", Type: "ui.Color", Description: "Primary text (>= 7:1 contrast on Background)"},
	{Path: "Colors.TextSecondary", Group: "text", Type: "ui.Color", Description: "Secondary body text"},
	{Path: "Colors.TextMuted", Group: "text", Type: "ui.Color", Description: "Captions, placeholders, muted icons"},
	{Path: "Colors.Accent", Group: "semantic", Type: "theme.SemanticTone", Description: "Primary brand accent (Base, Hover, Pressed, Subtle, Border, Text, OnColor)"},
	{Path: "Colors.Danger", Group: "semantic", Type: "theme.SemanticTone", Description: "Destructive and error intent"},
	{Path: "Colors.Warning", Group: "semantic", Type: "theme.SemanticTone", Description: "Caution and alert intent"},
	{Path: "Colors.Success", Group: "semantic", Type: "theme.SemanticTone", Description: "Positive completion intent"},
	{Path: "Colors.Info", Group: "semantic", Type: "theme.SemanticTone", Description: "Informational callout intent"},
	{Path: "Colors.Focus", Group: "chrome", Type: "ui.Color", Description: "Keyboard focus ring color"},
	{Path: "Colors.Selection", Group: "chrome", Type: "ui.Color", Description: "Text and row selection highlight"},
	{Path: "Colors.Inverse", Group: "chrome", Type: "ui.Color", Description: "Inverse tooltip/badge background"},
	{Path: "Colors.InverseText", Group: "chrome", Type: "ui.Color", Description: "Foreground on Inverse surface"},
	{Path: "Typography.Display", Group: "typography", Type: "theme.TextStyle", Description: "Hero display text (28px, weight 700)"},
	{Path: "Typography.H1", Group: "typography", Type: "theme.TextStyle", Description: "Page title (22px, weight 700)"},
	{Path: "Typography.H2", Group: "typography", Type: "theme.TextStyle", Description: "Section heading (17px, weight 600)"},
	{Path: "Typography.H3", Group: "typography", Type: "theme.TextStyle", Description: "Subsection heading (14.5px, weight 600)"},
	{Path: "Typography.Body", Group: "typography", Type: "theme.TextStyle", Description: "Standard body copy"},
	{Path: "Typography.BodySM", Group: "typography", Type: "theme.TextStyle", Description: "Compact secondary copy"},
	{Path: "Typography.Caption", Group: "typography", Type: "theme.TextStyle", Description: "Small metadata and badges"},
	{Path: "Typography.Code", Group: "typography", Type: "theme.TextStyle", Description: "Monospace code snippets"},
	{Path: "Radius.XS", Group: "radius", Type: "float32", Description: "Small checkboxes and badges"},
	{Path: "Radius.SM", Group: "radius", Type: "float32", Description: "Menu items and tags"},
	{Path: "Radius.MD", Group: "radius", Type: "float32", Description: "Buttons and inputs"},
	{Path: "Radius.LG", Group: "radius", Type: "float32", Description: "Cards and popovers"},
	{Path: "Radius.XL", Group: "radius", Type: "float32", Description: "Modals and feature panels"},
	{Path: "Radius.Full", Group: "radius", Type: "float32", Description: "Pill badges and switches"},
	{Path: "Shadows.SM", Group: "shadow", Type: "theme.ShadowSpec", Description: "Subtle card lift"},
	{Path: "Shadows.MD", Group: "shadow", Type: "theme.ShadowSpec", Description: "Elevated card and dropdown shadow"},
	{Path: "Shadows.LG", Group: "shadow", Type: "theme.ShadowSpec", Description: "Modal and drawer elevation"},
	{Path: "Motion.Fast", Group: "motion", Type: "ui.Animation", Description: "110ms micro-interaction transition"},
	{Path: "Motion.Normal", Group: "motion", Type: "ui.Animation", Description: "180ms standard transition"},
	{Path: "Motion.SpringSnappy", Group: "motion", Type: "ui.Animation", Description: "Snappy spring curve for toggles"},
}

// Presets summarizes all 7 built-in themes derived from theme.StockThemes().
var Presets = buildPresets()

func buildPresets() []PresetDoc {
	descs := map[string]string{
		"neutral":   "Clean crisp zinc surfaces with balanced cobalt accent.",
		"matcha":    "Warm botanical matcha greens on calm stone paper.",
		"stone":     "Editorial warm sandstone neutrals with terracotta accent.",
		"butter":    "Soft sunlit cream canvas with rich amber-gold accent.",
		"chocolate": "Warm espresso cocoa surfaces with caramel highlight.",
		"gothic":    "High-contrast obsidian ink with crimson velvet accent.",
		"y2k":       "Glossy cyber-silver chrome with electric cyan-blue accent.",
	}
	var out []PresetDoc
	for _, ts := range theme.StockThemes() {
		name := ts.Light.Name
		out = append(out, PresetDoc{
			Name:        name,
			Description: descs[name],
			LightBg:     hex(ts.Light.Colors.Background),
			LightAccent: hex(ts.Light.Colors.Accent.Base),
			DarkBg:      hex(ts.Dark.Colors.Background),
			DarkAccent:  hex(ts.Dark.Colors.Accent.Base),
		})
	}
	return out
}

func hex(c ui.Color) string {
	return fmt.Sprintf("#%02X%02X%02X", c.R, c.G, c.B)
}

// IconNames returns all available embedded icon names.
func IconNames() []string {
	all := icons.All()
	out := make([]string, len(all))
	for i, ic := range all {
		out[i] = ic.Name
	}
	return out
}

// Guides captures essential developer and agent guidance for mygo-astryx.
var Guides = []GuideDoc{
	{
		Name:    "mygo-rules",
		Title:   "MyGo Native UI & mygo vet Rules",
		Summary: "Frame-lifetime, concurrency, identity, and theme-scoping rules enforced by mygo vet.",
		Tags:    []string{"mygo", "vet", "lifetime", "rules"},
		Body: `# MyGo Native UI Rules
1. Elements do not outlive their frame: never store ui.Element, *ui.Context, or *Parts in structs, globals, or goroutine closures.
2. Never block the UI thread: run I/O in goroutines and apply results on the main thread via Window.Update / Window.Invalidate.
3. Scope themes per frame with theme.WithTheme(c, th, fn) and read tokens via theme.Of(c).
4. Always key repeated/dynamic stateful controls with c.Key(...) before constructing them.
5. Verify every change with: ./bin/go build ./... && ./bin/go vet ./... && ./bin/go tool mygo vet . && ./bin/go test ./...`,
	},
	{
		Name:    "theming",
		Title:   "Multi-Theme Token Architecture",
		Summary: "How DefineTheme, ColorPair, Density, SemanticTone, and WithTheme work together.",
		Tags:    []string{"theme", "tokens", "wcag", "dark"},
		Body: `# Theming Guide
Use one of the 7 stock ThemeSets (Neutral, Matcha, Stone, Butter, Chocolate, Gothic, Y2K) or create your own via theme.DefineTheme(theme.Spec{...}).
- Call th := theme.Matcha.For(dark).WithDensity(theme.DensityCompact)
- Wrap your root view in theme.WithTheme(c, th, func() { ... })
- Use theme.WithTheme(c, th2, ...) on any nested subtree; overlays opened inside inherit both ui.Theme and Astryx tokens.`,
	},
	{
		Name:    "swizzle",
		Title:   "Swizzling Components Locally",
		Summary: "How astryx swizzle ejects self-contained ds components into your own project.",
		Tags:    []string{"swizzle", "eject", "cli"},
		Body: `# Swizzle Guide
Run 'astryx swizzle Button --out ./ui/ds --pkg ds' to eject a component's Go source files into your local project.
Swizzle automatically pulls transitive internal helper files (such as options.go), rewrites the package declaration, and validates compilation.`,
	},
}

// Templates provides ready-to-scaffold native MyGo + Astryx starter applications.
var Templates = []TemplateDoc{
	{
		Name:        "starter",
		Title:       "Astryx Desktop Starter",
		Description: "Single-window native MyGo app with theme/dark toggles, FormField, Button, and ToastViewport.",
		Tags:        []string{"app", "minimal", "desktop"},
		Files: map[string]string{
			"main.go": `package main

import (
	"log"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
	"github.com/amzkang/mygo-astryx/ds"
	"github.com/amzkang/mygo-astryx/icons"
	"github.com/amzkang/mygo-astryx/theme"
)

func main() {
	dark := false
	themeName := "matcha"
	name := ""

	mygo.App.WhenReady(func() {
		mygo.NewWindow(mygo.WindowOptions{
			Title:  "Astryx Starter",
			Width:  720,
			Height: 480,
			Content: ui.View(func(c *ui.Context) {
				th := theme.ByName(themeName).For(dark)
				theme.WithTheme(c, th, func() {
					ui.Column(c).Fill().Background(th.Colors.Background).Padding(th.Space(6)).Gap(th.Space(4)).Children(func() {
						ds.HStack(c, 2, func() {
							ds.Heading(c, 1, "Astryx Starter").Grow(1)
							ds.Select(c, &themeName, []string{"neutral", "matcha", "stone", "butter", "chocolate", "gothic", "y2k"}, nil)
							ds.Switch(c, &dark, "Dark")
						})
						ds.Card(c, func() {
							ds.FormField(c, "Project Name", "Lowercase identifier", "", func() {
								ds.TextInput(c, &name, ds.InputOpts{Placeholder: "my-native-app"})
							})
							ds.HStack(c, 2, func() {
								if ds.PrimaryButton(c, "Create", ds.ButtonOpts{LeadingIcon: icons.Sparkles}).Clicked() {
									c.Toast("Created " + name)
								}
							})
						})
					})
					ds.ToastViewport(c)
				})
			}),
		})
	})
	if err := mygo.App.Run(); err != nil {
		log.Fatal(err)
	}
}
`,
		},
	},
	{
		Name:        "settings-shell",
		Title:       "Astryx Settings Shell",
		Description: "Two-pane settings layout with Tabs, Accordion, Switch controls, and Modal confirmation.",
		Tags:        []string{"settings", "layout", "tabs"},
		Files: map[string]string{
			"main.go": `package main

import (
	"log"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
	"github.com/amzkang/mygo-astryx/ds"
	"github.com/amzkang/mygo-astryx/theme"
)

func main() {
	tab := 0
	telemetry := false
	autoSave := true
	confirmReset := false

	mygo.App.WhenReady(func() {
		mygo.NewWindow(mygo.WindowOptions{
			Title:  "Settings",
			Width:  820,
			Height: 540,
			Content: ui.View(func(c *ui.Context) {
				th := theme.Neutral.For(c.Theme().Dark)
				theme.WithTheme(c, th, func() {
					ui.Column(c).Fill().Background(th.Colors.Background).Padding(th.Space(6)).Gap(th.Space(4)).Children(func() {
						ds.Heading(c, 1, "Preferences")
						ds.Tabs(c, &tab, []string{"General", "Privacy"}, ds.TabsOpts{Pills: true})
						ds.Card(c, func() {
							ds.Switch(c, &autoSave, "Auto-save changes")
							ds.Divider(c)
							ds.Switch(c, &telemetry, "Share crash diagnostics")
							ds.HStack(c, 2, func() {
								if ds.Button(c, "Reset Defaults", ds.ButtonOpts{Variant: ds.VariantOutline, Tone: ds.ToneDanger}).Clicked() {
									confirmReset = true
								}
							})
						})
					})
					ds.ModalCard(c, &confirmReset, "Reset Preferences?", 380, func() {
						ds.Body(c, "All custom preferences will revert to defaults.")
						ds.HStack(c, 2, func() {
							if ds.Button(c, "Cancel", ds.ButtonOpts{Variant: ds.VariantOutline}).Clicked() {
								confirmReset = false
							}
							if ds.Button(c, "Reset", ds.ButtonOpts{Variant: ds.VariantSolid, Tone: ds.ToneDanger}).Clicked() {
								autoSave = true
								telemetry = false
								confirmReset = false
								c.Toast("Preferences reset")
							}
						})
					})
					ds.ToastViewport(c)
				})
			}),
		})
	})
	if err := mygo.App.Run(); err != nil {
		log.Fatal(err)
	}
}
`,
		},
	},
}

package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
	"github.com/amzkang/mygo-astryx/ds"
	"github.com/amzkang/mygo-astryx/icons"
	"github.com/amzkang/mygo-astryx/registry"
	"github.com/amzkang/mygo-astryx/theme"
)

type galleryPage string

const (
	pageOverview   galleryPage = "overview"
	pageComponents galleryPage = "components"
	pageTokens     galleryPage = "tokens"
	pageIcons      galleryPage = "icons"
)

type galleryState struct {
	page        galleryPage
	themeName   string
	dark        bool
	densityIdx  int // 0=comfortable, 1=compact, 2=spacious
	rtl         bool
	sidebarSel  string
	compSearch  string
	iconSearch  string
	selectedCmp string

	// Playground interactive controls
	playTab      int
	playPage     int
	playInput    string
	playNotes    string
	playChecked  bool
	playSwitch   bool
	playRadio    string
	playSeg      int
	playSelect   string
	playCombo    string
	playSlider   float64
	playAccOpen1 bool
	playAccOpen2 bool
	playModal    bool
	playDrawer   bool
	playPopover  bool
	paletteOpen  bool
	paletteState ds.CommandPaletteState

	// Component Playground & Desktop Showcase state
	pgVariantIdx int
	pgSizeIdx    int
	pgToneIdx    int
	pgDisabled   bool
	pgLoading    bool
	pgDot        bool
	pgLabel      string
	pgMenuOpen   bool
	pgAlertOpen  bool
	pgCollOpen   bool
	pgTreeOpen   bool
	pgTreeSel    string
	pgTableState ui.ListState
	pgNumVal     float64
	pgStepVal    float64
	pgRangeLo    float64
	pgRangeHi    float64
	pgTokens     []string
	pgDate       time.Time
	pgSplit      float32

	// Tokens & WCAG mode (0=active theme swatches, 1=7-theme contrast matrix)
	tokensTab int
}

func newGalleryState() *galleryState {
	return &galleryState{
		page:         pageOverview,
		themeName:    "matcha",
		sidebarSel:   "overview",
		selectedCmp:  "Button",
		playPage:     2,
		playInput:    "Astryx Native App",
		playChecked:  true,
		playSwitch:   true,
		playRadio:    "solid",
		playSelect:   "Matcha",
		playSlider:   72,
		playAccOpen1: true,
		pgSizeIdx:    2, // MD
		pgToneIdx:    1, // Accent
		pgLabel:      "Deploy Release",
		pgTreeOpen:   true,
		pgTreeSel:    "ds",
		pgNumVal:     12,
		pgStepVal:    3,
		pgRangeLo:    20,
		pgRangeHi:    80,
		pgTokens:     []string{"native", "wcag"},
		pgDate:       time.Date(2026, 11, 14, 0, 0, 0, 0, time.UTC),
		pgSplit:      210,
	}
}

func main() {
	app := newGalleryState()
	mygo.App.WhenReady(func() {
		mygo.NewWindow(mygo.WindowOptions{
			Title:  "Astryx Gallery — Native MyGo Design System",
			Width:  1180,
			Height: 780,
			Content: ui.View(func(c *ui.Context) {
				renderGallery(c, app)
			}),
		})
	})
	if err := mygo.App.Run(); err != nil {
		log.Fatal(err)
	}
}

func currentDensity(idx int) theme.Density {
	switch idx {
	case 1:
		return theme.DensityCompact
	case 2:
		return theme.DensitySpacious
	default:
		return theme.DensityComfortable
	}
}

func renderGallery(c *ui.Context, app *galleryState) {
	if c.Shortcut(ui.Cmd, ui.KeyK) {
		app.paletteOpen = !app.paletteOpen
	}

	ts := theme.ByName(app.themeName)
	th := ts.For(app.dark).WithDensity(currentDensity(app.densityIdx))

	theme.WithTheme(c, th, func() {
		root := ui.Column(c).Fill().Background(th.Colors.Background)
		if app.rtl {
			root.Direction(ui.RightToLeft)
		} else {
			root.Direction(ui.LeftToRight)
		}
		root.Children(func() {
			renderTopBar(c, app, th)
			ds.Divider(c)
			ui.Row(c).Grow(1).MinHeight(0).Children(func() {
				renderSidebar(c, app, th)
				ui.Box(c).Width(1).FillHeight().Background(th.Colors.BorderSubtle)
				ui.Scroll(c).Grow(1).MinHeight(0).Padding(th.Space(6)).Children(func() {
					ds.Container(c, 980, func() {
						switch app.page {
						case pageComponents:
							renderComponentsPage(c, app, th)
						case pageTokens:
							renderTokensPage(c, app, th)
						case pageIcons:
							renderIconsPage(c, app, th)
						default:
							renderOverviewPage(c, app, th)
						}
					})
				})
			})
		})

		ds.ModalCard(c, &app.playModal, "Theme Confirmation", 420, func() {
			ds.Body(c, fmt.Sprintf("Active theme is %q (%s density, dark=%v).", th.Name, th.Density, th.Dark))
			ds.HStack(c, 2, func() {
				if ds.Button(c, "Cancel", ds.ButtonOpts{Variant: ds.VariantOutline}).Clicked() {
					app.playModal = false
				}
				if ds.PrimaryButton(c, "Apply & Toast").Clicked() {
					app.playModal = false
					c.Toast("Applied theme " + th.Name)
				}
			})
		})

		ds.Drawer(c, &app.playDrawer, func() {
			ds.VStack(c, 3, func() {
				ds.Body(c, "Inspect live WCAG ratios and layout tokens for the current theme.")
				ds.KeyValueRow(c, "Theme", th.Name)
				ds.KeyValueRow(c, "Density", th.Density.String())
				ds.KeyValueRow(c, "Text/Bg Contrast", fmt.Sprintf("%.2f:1", theme.ContrastRatio(th.Colors.Text, th.Colors.Background)))
				ds.KeyValueRow(c, "Accent OnColor/Base", fmt.Sprintf("%.2f:1", theme.ContrastRatio(th.Colors.Accent.OnColor, th.Colors.Accent.Base)))
			})
		}, ds.DrawerOpts{Title: "Token Inspector", Side: ds.DrawerRight})

		ds.CommandPalette(c, &app.paletteOpen, &app.paletteState, buildPaletteCommands(app, c))
		ds.ToastViewport(c)
	})
}

func renderTopBar(c *ui.Context, app *galleryState, th *theme.Theme) {
	themeNames := []string{"neutral", "matcha", "stone", "butter", "chocolate", "gothic", "y2k"}
	ui.Row(c).
		FillWidth().
		AlignItems(ui.Center).
		Justify(ui.SpaceBetween).
		Gap(th.Space(3)).
		Padding(th.Space(2.5), th.Space(4)).
		Background(th.Colors.Surface).
		Children(func() {
			ds.HStack(c, 2, func() {
				ui.Box(c).Size(28, 28).Center().Radius(th.Radius.MD).Background(th.Colors.Accent.Base).Children(func() {
					ui.Icon(c, icons.Sparkles).Size(16, 16).TextColor(th.Colors.Accent.OnColor)
				})
				ds.Heading(c, 3, "Astryx Gallery")
				ds.Badge(c, "MyGo Native", ds.BadgeOpts{Tone: ds.ToneAccent})
			})
			ds.HStack(c, 2.5, func() {
				if ds.Button(c, "Command Palette", ds.ButtonOpts{Variant: ds.VariantOutline, Size: ds.SizeSM, LeadingIcon: icons.Command}).Clicked() {
					app.paletteOpen = true
				}
				ds.Select(c, &app.themeName, themeNames, func(s string) string {
					return strings.ToUpper(s[:1]) + s[1:]
				})
				ds.SegmentedControl(c, &app.densityIdx, []string{"Comfy", "Compact", "Spacious"})
				ds.Switch(c, &app.dark, "Dark")
				ds.Switch(c, &app.rtl, "RTL")
			})
		})
}

func renderSidebar(c *ui.Context, app *galleryState, th *theme.Theme) {
	ui.Column(c).Width(220).FillHeight().Background(th.Colors.SurfaceSubtle).Padding(th.Space(3)).Gap(th.Space(1)).Children(func() {
		navItem := func(id galleryPage, icon *ui.SVG, label string) {
			active := app.page == id
			btn := ui.ButtonBase(c.Key("nav-"+string(id))).
				FillWidth().
				AlignItems(ui.Center).
				Justify(ui.Start).
				Gap(th.Space(2.5)).
				Padding(th.Space(2), th.Space(2.5)).
				Radius(th.Radius.MD).
				Cursor(ui.CursorPointer)
			if active {
				btn.Background(th.Colors.Accent.Subtle)
			} else if btn.Hovered() {
				btn.Background(th.Colors.SurfaceHover)
			}
			fg := th.Colors.TextSecondary
			weight := 500
			if active {
				fg = th.Colors.Accent.Text
				weight = 600
			}
			btn.Children(func() {
				ui.Icon(c, icon).Size(16, 16).TextColor(fg)
				ui.Text(c, label).FontSize(th.Typography.BodySM.Size).FontWeight(weight).TextColor(fg).SingleLine()
			})
			if btn.Clicked() {
				app.page = id
				app.sidebarSel = string(id)
			}
		}

		ds.Caption(c, "EXPLORE").Padding(th.Space(1), th.Space(2))
		navItem(pageOverview, icons.Home, "Overview")
		navItem(pageComponents, icons.Layers, "Components")
		navItem(pageTokens, icons.Palette, "Tokens & WCAG")
		navItem(pageIcons, icons.Box, "Icons (56)")

		ui.Spacer(c)
		ds.Card(c, func() {
			ds.Caption(c, "Active Theme")
			ds.HStack(c, 1.5, func() {
				ui.Box(c).Size(12, 12).Radius(6).Background(th.Colors.Accent.Base)
				ui.Text(c, th.Name).FontSize(th.Typography.BodySM.Size).FontWeight(600).TextColor(th.Colors.Text)
			})
			if ds.Button(c, "Inspect Tokens", ds.ButtonOpts{Variant: ds.VariantOutline, Size: ds.SizeXS, FullWidth: true}).Clicked() {
				app.playDrawer = true
			}
		}, ds.CardOpts{Padding: th.Space(3)})
	})
}

func renderOverviewPage(c *ui.Context, app *galleryState, th *theme.Theme) {
	ds.VStack(c, 5, func() {
		ds.Section(c, "Native Go Design System", "Built on MyGo v0.8.3 GPU UI with Astryx multi-theme tokens, WCAG contrast validation, swizzle ejection, and MCP.", func() {
			ui.Row(c).FillWidth().Gap(th.Space(3)).Children(func() {
				ds.StatTile(c, "Components", fmt.Sprint(len(registry.Components)), "ds package", ds.ToneAccent).Grow(1)
				ds.StatTile(c, "Stock Themes", fmt.Sprint(len(registry.Presets)), "light + dark", ds.ToneSuccess).Grow(1)
				ds.StatTile(c, "Embedded Icons", fmt.Sprint(len(icons.All())), "24x24 SVG", ds.ToneInfo).Grow(1)
				ds.StatTile(c, "Text/Bg Contrast", fmt.Sprintf("%.1f:1", theme.ContrastRatio(th.Colors.Text, th.Colors.Background)), "WCAG AAA", ds.ToneSuccess).Grow(1)
			})
		})

		ds.Card(c, func() {
			ds.Heading(c, 2, "Interactive Component Showcase")
			ds.Tabs(c, &app.playTab, []string{"Controls", "Feedback & Overlays", "Desktop Controls", "Nested Theme Scope"}, ds.TabsOpts{Pills: true})
			switch app.playTab {
			case 0:
				ds.VStack(c, 3, func() {
					ds.HStack(c, 2, func() {
						if ds.PrimaryButton(c, "Primary", ds.ButtonOpts{LeadingIcon: icons.Sparkles}).Clicked() {
							c.Toast("Primary clicked")
						}
						ds.Button(c, "Soft Success", ds.ButtonOpts{Variant: ds.VariantSoft, Tone: ds.ToneSuccess, LeadingIcon: icons.Check})
						ds.Button(c, "Outline", ds.ButtonOpts{Variant: ds.VariantOutline})
						ds.Button(c, "Danger", ds.ButtonOpts{Variant: ds.VariantSolid, Tone: ds.ToneDanger})
						ds.CopyButton(c, "Copy CLI", "go run ./cmd/astryx search button")
					})
					ds.FormField(c, "Workspace Name", "Live bound to Go state", "", func() {
						ds.TextInput(c, &app.playInput, ds.InputOpts{LeadingIcon: icons.Edit, Placeholder: "Enter name"})
					})
					ds.HStack(c, 4, func() {
						ds.Checkbox(c, &app.playChecked, "Strict mygo vet")
						ds.Radio(c.Key("rad-solid"), &app.playRadio, "solid", "Solid")
						ds.Radio(c.Key("rad-soft"), &app.playRadio, "soft", "Soft")
						ds.Select(c, &app.playSelect, []string{"Matcha", "Neutral", "Gothic", "Y2K"}, nil)
					})
					ds.Slider(c, &app.playSlider, 0, 100, "Progress Preview")
					ds.Progress(c, app.playSlider/100, ds.ToneAccent)
				})
			case 1:
				ds.VStack(c, 3, func() {
					ds.Alert(c, "All 7 themes satisfy >=7:1 primary text contrast and >=4.5:1 button foreground contrast.", ds.AlertOpts{
						Tone:   ds.ToneSuccess,
						Title:  "WCAG Verified",
						Action: "Inspect",
						OnClick: func() {
							app.playDrawer = true
						},
					})
					ds.HStack(c, 2, func() {
						if ds.Button(c, "Open Modal").Clicked() {
							app.playModal = true
						}
						if ds.Button(c, "Open Drawer").Clicked() {
							app.playDrawer = true
						}
						popBtn := ds.Button(c, "Toggle Popover")
						if popBtn.Clicked() {
							app.playPopover = !app.playPopover
						}
						ds.Tooltip(c, popBtn, "Anchored using ui.PopoverBase")
						ds.Popover(c, popBtn, &app.playPopover, func() {
							ds.VStack(c, 1.5, func() {
								ds.Heading(c, 3, "Popover Panel")
								ds.Caption(c, "Inherits theme.WithTheme automatically.")
							})
						})
						if ds.Button(c, "Fire Toast", ds.ButtonOpts{Tone: ds.ToneInfo, Variant: ds.VariantSoft}).Clicked() {
							c.AddToast(ui.Toast{
								Title:       "Astryx Toast",
								Description: "Rendered via ds.ToastViewport",
								Type:        "success",
							})
						}
					})
					ds.Pagination(c, &app.playPage, 9)
					ds.Accordion(c, []ds.AccordionItemSpec{
						{Title: "How Swizzle Works", Open: &app.playAccOpen1, Body: func() {
							ds.Body(c, "Run 'astryx swizzle ModalCard --out ./ui/ds' to eject self-contained Go files.")
						}},
						{Title: "Why Frame-Scoped Tokens", Open: &app.playAccOpen2, Body: func() {
							ds.Body(c, "theme.WithTheme pushes onto a per-frame stack via ui.Local so no pointers outlive the frame.")
						}},
					})
				})
			case 2:
				renderDesktopControlsShowcase(c, app, th)
			case 3:
				ds.HStack(c, 3, func() {
					for _, previewName := range []string{"gothic", "butter", "y2k"} {
						subTh := theme.ByName(previewName).For(app.dark)
						theme.WithTheme(c.Key("sub-"+previewName), subTh, func() {
							ds.Card(c, func() {
								ds.HStack(c, 2, func() {
									ds.Avatar(c, strings.ToUpper(previewName[:2]), 32, nil)
									ds.Heading(c, 3, strings.ToUpper(previewName[:1])+previewName[1:])
								})
								ds.Caption(c, "Nested theme.WithTheme subtree")
								ds.PrimaryButton(c, "Scoped CTA", ds.ButtonOpts{Size: ds.SizeSM})
							}).Grow(1)
						})
					}
				})
			}
		}, ds.CardOpts{Elevated: true})
	})
}

func renderComponentsPage(c *ui.Context, app *galleryState, th *theme.Theme) {
	ds.VStack(c, 4, func() {
		ds.Section(c, "Component Catalog", "Browse all exported ds components, props, accessibility notes, and copyable Go snippets.", func() {
			ds.SearchInput(c, &app.compSearch, "Search components (e.g. button, modal, input)...")
		})
		hits := registry.Search(app.compSearch, "component")
		ui.Row(c).FillWidth().Gap(th.Space(4)).AlignItems(ui.Start).Children(func() {
			ds.Card(c, func() {
				ds.Caption(c, fmt.Sprintf("%d components", len(hits)))
				ui.Column(c).Gap(th.Space(1)).Children(func() {
					for _, h := range hits {
						active := strings.EqualFold(h.Name, app.selectedCmp)
						btn := ui.ButtonBase(c.Key("cmp-"+h.Name)).
							FillWidth().
							Justify(ui.Start).
							Padding(th.Space(1.5), th.Space(2.5)).
							Radius(th.Radius.SM).
							Cursor(ui.CursorPointer)
						if active {
							btn.Background(th.Colors.Accent.Subtle)
						} else if btn.Hovered() {
							btn.Background(th.Colors.SurfaceHover)
						}
						fg := th.Colors.Text
						if active {
							fg = th.Colors.Accent.Text
						}
						btn.Children(func() {
							ui.Text(c, h.Name).FontSize(th.Typography.BodySM.Size).FontWeight(500).TextColor(fg)
						})
						if btn.Clicked() {
							app.selectedCmp = h.Name
						}
					}
				})
			}, ds.CardOpts{Padding: th.Space(3)}).Width(240)

			if comp, ok := registry.FindComponent(app.selectedCmp); ok {
				ds.VStack(c, 3, func() {
					renderComponentPlayground(c, app, th, comp)
					ds.Card(c, func() {
						ds.HStack(c, 2, func() {
							ds.Heading(c, 2, comp.Name).Grow(1)
							ds.Badge(c, comp.Category, ds.BadgeOpts{Tone: ds.ToneAccent})
							ds.CopyButton(c, "Copy Snippet", comp.Example, ds.ButtonOpts{Size: ds.SizeSM})
						})
						ds.Body(c, comp.Summary)
						ds.KeyValueRow(c, "Signature", comp.Signature)
						ds.KeyValueRow(c, "Swizzle Files", strings.Join(comp.SourceFiles, ", "))
						ds.Alert(c, comp.A11y, ds.AlertOpts{Tone: ds.ToneInfo, Title: "Accessibility"})
						ds.Heading(c, 3, "Reference Example")
						ui.Box(c).FillWidth().Padding(th.Space(3)).Radius(th.Radius.MD).Background(th.Colors.SurfaceSubtle).Border(1, th.Colors.BorderSubtle).Children(func() {
							ui.Text(c, comp.Example).Font(th.MonoFont).FontSize(th.Typography.Code.Size).TextColor(th.Colors.Text)
						})
					})
				}).Grow(1)
			}
		})
	})
}

func renderTokensPage(c *ui.Context, app *galleryState, th *theme.Theme) {
	ds.VStack(c, 4, func() {
		ds.Section(c, "Theme Tokens & WCAG Inspector", fmt.Sprintf("Inspecting %s (dark=%v, density=%s)", th.Name, th.Dark, th.Density), func() {
			ds.Tabs(c, &app.tokensTab, []string{"Current Theme", "7-Theme Matrix"}, ds.TabsOpts{Pills: true})
			if app.tokensTab == 1 {
				renderThemeContrastMatrix(c, app, th)
				return
			}
			swatches := []struct {
				name string
				bg   ui.Color
				fg   ui.Color
			}{
				{"Background / Text", th.Colors.Background, th.Colors.Text},
				{"Surface / Text", th.Colors.Surface, th.Colors.Text},
				{"Accent / OnColor", th.Colors.Accent.Base, th.Colors.Accent.OnColor},
				{"Success / OnColor", th.Colors.Success.Base, th.Colors.Success.OnColor},
				{"Warning / OnColor", th.Colors.Warning.Base, th.Colors.Warning.OnColor},
				{"Danger / OnColor", th.Colors.Danger.Base, th.Colors.Danger.OnColor},
				{"Info / OnColor", th.Colors.Info.Base, th.Colors.Info.OnColor},
			}
			ds.VStack(c, 2, func() {
				for _, sw := range swatches {
					ratio := theme.ContrastRatio(sw.fg, sw.bg)
					grade := "AA"
					tone := ds.ToneSuccess
					if ratio >= 7.0 {
						grade = "AAA"
					} else if ratio < 4.5 {
						grade = "Large"
						tone = ds.ToneWarning
					}
					ds.Card(c, func() {
						ds.HStack(c, 3, func() {
							ui.Box(c).Size(38, 38).Center().Radius(th.Radius.MD).Background(sw.bg).Border(1, th.Colors.Border).Children(func() {
								ui.Text(c, "Aa").FontWeight(700).TextColor(sw.fg)
							})
							ui.Column(c).Grow(1).Gap(2).Children(func() {
								ui.Text(c, sw.name).FontWeight(600).TextColor(th.Colors.Text)
								ds.Caption(c, fmt.Sprintf("bg=#%02X%02X%02X  fg=#%02X%02X%02X", sw.bg.R, sw.bg.G, sw.bg.B, sw.fg.R, sw.fg.G, sw.fg.B))
							})
							ds.Badge(c, fmt.Sprintf("%.2f:1 (%s)", ratio, grade), ds.BadgeOpts{Tone: tone, Dot: true})
						})
					}, ds.CardOpts{Padding: th.Space(3)})
				}
			})
		})
	})
}

func renderIconsPage(c *ui.Context, app *galleryState, th *theme.Theme) {
	ds.VStack(c, 4, func() {
		ds.Section(c, "Embedded Icon Pack", "56 crisp 24x24 Lucide-style SVG icons compiled directly into the binary.", func() {
			ds.SearchInput(c, &app.iconSearch, "Filter icons by name (e.g. arrow, check, sun)...")
		})
		matches := icons.Find(app.iconSearch)
		ds.Caption(c, fmt.Sprintf("Showing %d icons", len(matches)))
		if len(matches) == 0 {
			ds.EmptyState(c, "No matching icons", "Try another keyword such as chevron or alert.")
			return
		}
		const cols = 4
		ds.VStack(c, 2, func() {
			for i := 0; i < len(matches); i += cols {
				ui.Row(c).FillWidth().Gap(th.Space(2)).Children(func() {
					for j := 0; j < cols; j++ {
						idx := i + j
						if idx >= len(matches) {
							ui.Spacer(c)
							continue
						}
						ic := matches[idx]
						ds.Card(c, func() {
							ds.HStack(c, 2, func() {
								ui.Icon(c, ic.SVG).Size(20, 20).TextColor(th.Colors.Text)
								ui.Column(c).Grow(1).MinWidth(0).Gap(1).Children(func() {
									ui.Text(c, ic.Symbol).FontSize(th.Typography.BodySM.Size).FontWeight(600).TextColor(th.Colors.Text).SingleLine()
									ds.Caption(c, "icons."+ic.Symbol)
								})
							})
						}, ds.CardOpts{Padding: th.Space(2.5)}).Grow(1)
					}
				})
			}
		})
	})
}

func buildPaletteCommands(app *galleryState, c *ui.Context) []ds.CommandItem {
	items := []ds.CommandItem{
		{ID: "go-overview", Title: "Go to Overview", Category: "Navigate", Icon: icons.Home, OnRun: func() { app.page = pageOverview }},
		{ID: "go-components", Title: "Browse Components", Category: "Navigate", Icon: icons.Layers, OnRun: func() { app.page = pageComponents }},
		{ID: "go-tokens", Title: "Inspect Tokens & WCAG", Category: "Navigate", Icon: icons.Palette, OnRun: func() { app.page = pageTokens }},
		{ID: "go-icons", Title: "Browse Icons", Category: "Navigate", Icon: icons.Box, OnRun: func() { app.page = pageIcons }},
		{ID: "toggle-dark", Title: "Toggle Dark / Light Mode", Category: "Appearance", Icon: icons.Moon, OnRun: func() { app.dark = !app.dark }},
		{ID: "toggle-rtl", Title: "Toggle RTL Direction", Category: "Appearance", Icon: icons.Sliders, OnRun: func() { app.rtl = !app.rtl }},
	}
	for _, ts := range theme.StockThemes() {
		name := ts.Light.Name
		items = append(items, ds.CommandItem{
			ID:       "theme-" + name,
			Title:    "Switch Theme: " + strings.ToUpper(name[:1]) + name[1:],
			Category: "Theme",
			Icon:     icons.Palette,
			OnRun:    func() { app.themeName = name },
		})
	}
	return items
}

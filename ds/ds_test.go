package ds_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/amzkang/mygo-astryx/ds"
	"github.com/amzkang/mygo-astryx/icons"
	"github.com/amzkang/mygo-astryx/theme"
)

func TestButtonsAndCopyButton(t *testing.T) {
	clicked := 0
	tester := ui.NewTester(func(c *ui.Context) {
		th := theme.Matcha.For(false)
		theme.WithTheme(c, th, func() {
			ds.VStack(c, 2, func() {
				if ds.PrimaryButton(c, "Create").Clicked() {
					clicked++
				}
				ds.Button(c, "Soft", ds.ButtonOpts{Variant: ds.VariantSoft, Tone: ds.ToneSuccess, LeadingIcon: icons.Check})
				ds.Button(c, "Loading", ds.ButtonOpts{Loading: true})
				ds.IconButton(c, "Settings", icons.Settings)
				ds.CopyButton(c, "Copy Token", "astyx_123")
			})
		})
	}, 640, 480)

	if err := tester.Click("Create"); err != nil {
		t.Fatalf("failed to click PrimaryButton Create: %v", err)
	}
	if clicked != 1 {
		t.Fatalf("expected clicked=1, got %d", clicked)
	}
	if err := tester.Click("Copy Token"); err != nil {
		t.Fatalf("failed to click CopyButton: %v", err)
	}
	if _, ok := tester.AccessByRole(ui.RoleButton, "Settings"); !ok {
		t.Fatal("expected Settings IconButton in accessibility tree")
	}
	if node, ok := tester.AccessByRole(ui.RoleButton, "Loading"); !ok || !node.Disabled {
		t.Fatalf("expected Loading button to be disabled in access tree, ok=%v node=%+v", ok, node)
	}
}

func TestFormsControlsAndBindings(t *testing.T) {
	text := ""
	area := ""
	query := "matcha"
	agreed := false
	enabled := false
	choice := "a"
	seg := 0
	sel := "Matcha"
	combo := ""
	sliderVal := 40.0

	tester := ui.NewTester(func(c *ui.Context) {
		th := theme.Neutral.For(false)
		theme.WithTheme(c, th, func() {
			ds.VStack(c, 2, func() {
				ds.FormField(c, "Project Name", "Lowercase slug", "", func() {
					ds.TextInput(c, &text, ds.InputOpts{Placeholder: "Enter name"})
				})
				ds.TextArea(c, &area, "Notes", 80)
				ds.SearchInput(c, &query, "Filter")
				ds.Checkbox(c, &agreed, "Accept terms")
				ds.Switch(c, &enabled, "Enable telemetry")
				ds.HStack(c, 2, func() {
					ds.Radio(c.Key("r1"), &choice, "a", "Alpha")
					ds.Radio(c.Key("r2"), &choice, "b", "Beta")
				})
				ds.SegmentedControl(c, &seg, []string{"List", "Grid"})
				ds.Select(c, &sel, []string{"Neutral", "Matcha", "Gothic"}, nil)
				ds.Combobox(c, &combo, []string{"Button", "Badge", "Drawer"}, "Pick widget")
				ds.Slider(c, &sliderVal, 0, 100, "Opacity")
				ds.Progress(c, 0.65, ds.ToneAccent)
			})
		})
	}, 800, 900)

	if err := tester.Click("Accept terms"); err != nil || !agreed {
		t.Fatalf("Checkbox click did not toggle agreed: err=%v agreed=%v", err, agreed)
	}
	if err := tester.Click("Enable telemetry"); err != nil || !enabled {
		t.Fatalf("Switch click did not toggle enabled: err=%v enabled=%v", err, enabled)
	}
	if err := tester.Click("Beta"); err != nil || choice != "b" {
		t.Fatalf("Radio click did not select Beta: err=%v choice=%q", err, choice)
	}
	if err := tester.Click("Grid"); err != nil || seg != 1 {
		t.Fatalf("SegmentedControl click did not select Grid: err=%v seg=%d", err, seg)
	}
	if err := tester.Click("Clear search"); err != nil || query != "" {
		t.Fatalf("SearchInput clear button did not reset query: err=%v query=%q", err, query)
	}
}

func TestNavigationOverlaysAndDisplay(t *testing.T) {
	tab := 0
	page := 2
	crumb := -1
	accOpen := true
	modalOpen := false
	drawerOpen := false
	popOpen := false
	paletteOpen := false
	var palState ds.CommandPaletteState
	ranCmd := ""

	tester := ui.NewTester(func(c *ui.Context) {
		th := theme.Gothic.For(true)
		theme.WithTheme(c, th, func() {
			ds.VStack(c, 2, func() {
				ds.Heading(c, 1, "Astryx Showcase")
				ds.Breadcrumbs(c, []string{"Home", "Tokens", "Colors"}, &crumb)
				ds.Tabs(c, &tab, []string{"Overview", "Themes", "Code"})
				ds.Pagination(c, &page, 10)

				ds.HStack(c, 2, func() {
					if ds.Button(c, "Open Modal").Clicked() {
						modalOpen = true
					}
					if ds.Button(c, "Open Drawer").Clicked() {
						drawerOpen = true
					}
					popAnchor := ds.Button(c, "Open Popover")
					if popAnchor.Clicked() {
						popOpen = !popOpen
					}
					ds.Tooltip(c, popAnchor, "Show details")
					ds.Popover(c, popAnchor, &popOpen, func() {
						ds.Body(c, "Popover content")
					})
					if ds.Button(c, "Open Palette").Clicked() {
						paletteOpen = true
					}
				})

				ds.Accordion(c, []ds.AccordionItemSpec{
					{
						Title: "Design Principles",
						Open:  &accOpen,
						Body:  func() { ds.Body(c, "High contrast semantic layers.") },
					},
				})

				ds.HStack(c, 2, func() {
					ds.StatTile(c, "Themes", "7", "+100%", ds.ToneSuccess)
					ds.Avatar(c, "Ada Lovelace", 40, nil)
					ds.Badge(c, "Stable", ds.BadgeOpts{Tone: ds.ToneSuccess, Dot: true})
					ds.Tag(c, "Go 1.27", ds.TagOpts{Tone: ds.ToneInfo})
				})

				ds.Alert(c, "All contrast checks passed", ds.AlertOpts{Tone: ds.ToneSuccess, Title: "Verified"})
				ds.EmptyState(c, "No items", "Create one to begin.")
				ds.Skeleton(c, 120, 16, 4)
				ds.KeyValueRow(c, "Theme", th.Name)
			})

			ds.ModalCard(c, &modalOpen, "Confirm Action", 400, func() {
				ds.Body(c, "Inside modal")
				if ds.PrimaryButton(c, "Close Modal").Clicked() {
					modalOpen = false
				}
			})
			ds.Drawer(c, &drawerOpen, func() {
				ds.Body(c, "Inside drawer")
			}, ds.DrawerOpts{Title: "Inspector", Side: ds.DrawerRight})
			ds.CommandPalette(c, &paletteOpen, &palState, []ds.CommandItem{
				{ID: "switch-matcha", Title: "Switch to Matcha", Category: "Theme", OnRun: func() { ranCmd = "matcha" }},
				{ID: "open-docs", Title: "Open Docs", Category: "Help", OnRun: func() { ranCmd = "docs" }},
			})
			ds.ToastViewport(c)
		})
	}, 960, 1000)

	if err := tester.Click("Code"); err != nil || tab != 2 {
		t.Fatalf("expected tab=2 after clicking Code: err=%v tab=%d", err, tab)
	}
	if err := tester.Click("Next page"); err != nil || page != 3 {
		t.Fatalf("expected page=3: err=%v page=%d", err, page)
	}

	if err := tester.Click("Open Modal"); err != nil || !modalOpen {
		t.Fatalf("expected modalOpen=true: err=%v", err)
	}
	tester.Key(0, ui.KeyEscape)
	if modalOpen {
		t.Fatal("expected Escape to close ModalCard")
	}

	if err := tester.Click("Open Drawer"); err != nil || !drawerOpen {
		t.Fatalf("expected drawerOpen=true: err=%v", err)
	}
	tester.Key(0, ui.KeyEscape)
	if drawerOpen {
		t.Fatal("expected Escape to close Drawer")
	}

	if err := tester.Click("Open Palette"); err != nil || !paletteOpen {
		t.Fatalf("expected paletteOpen=true: err=%v", err)
	}
	tester.Type("matcha")
	tester.Key(0, ui.KeyEnter)
	if paletteOpen || ranCmd != "matcha" {
		t.Fatalf("expected command palette to run matcha and close, open=%v ran=%q", paletteOpen, ranCmd)
	}
}

func TestRenderSmokeAcrossThemes(t *testing.T) {
	for _, ts := range theme.StockThemes() {
		for _, dark := range []bool{false, true} {
			th := ts.For(dark)
			img := ui.Render(func(c *ui.Context) {
				theme.WithTheme(c, th, func() {
					ds.Card(c, func() {
						ds.Heading(c, 2, th.Name)
						ds.HStack(c, 2, func() {
							ds.PrimaryButton(c, "Save")
							ds.Badge(c, "OK", ds.BadgeOpts{Tone: ds.ToneSuccess, Dot: true})
						})
					})
				})
			}, 400, 220, 1)
			if img == nil || img.Bounds().Dx() != 400 || img.Bounds().Dy() != 220 {
				t.Fatalf("unexpected rendered image for %s dark=%v", th.Name, dark)
			}
		}
	}
}

func TestMenuAlertDataAndAdvancedForms(t *testing.T) {
	menuOpen := false
	ctxOpen := false
	alertOpen := false
	alertConfirmed := false
	menuAction := ""
	collOpen := false
	treeOpen := true
	selectedNode := "root"
	splitWidth := float32(220)
	numVal := 12.0
	stepVal := 3.0
	lowVal, highVal := 20.0, 80.0
	tokens := []string{"native"}
	due := time.Date(2026, 11, 14, 0, 0, 0, 0, time.UTC)
	var tableState ui.ListState
	rows := [][]string{{"Neutral", "AA"}, {"Matcha", "AAA"}}
	cols := []ui.TableColumn{{Title: "Theme", Width: 140}, {Title: "Grade", Width: 80}}

	tester := ui.NewTester(func(c *ui.Context) {
		th := theme.Matcha.For(false)
		theme.WithTheme(c, th, func() {
			ds.SplitPane(c, &splitWidth, false, func() {
				ds.TreeView(c, []ds.TreeNode{
					{
						ID:   "root",
						Label: "Themes",
						Icon: icons.Palette,
						Open: &treeOpen,
						Children: []ds.TreeNode{
							{ID: "matcha", Label: "Matcha", Icon: icons.Sparkles},
						},
					},
				}, &selectedNode)
			}, func() {
				ds.VStack(c, 2, func() {
					ds.Toolbar(c, func() {
						ds.DropdownMenu(c, "Menu Actions", &menuOpen, []ds.MenuItemSpec{
							{Label: "Duplicate", Icon: icons.Copy, OnClick: func() { menuAction = "duplicate" }},
							{Separator: true},
							{Label: "Delete", Danger: true, OnClick: func() { menuAction = "delete" }},
						})
						if ds.Button(c, "Trigger Alert").Clicked() {
							alertOpen = true
						}
					})

					target := ds.Card(c, func() {
						ds.Body(c, "Context target")
					})
					ds.ContextMenu(c, target, &ctxOpen, []ds.MenuItemSpec{
						{Label: "Inspect", OnClick: func() { menuAction = "inspect" }},
					})

					ds.Collapsible(c, "Advanced Settings", &collOpen, func() {
						ds.Body(c, "Inside collapsible")
					})

					ds.DataTable(c, &tableState, cols, len(rows), "Empty", func(row, col int) {
						ds.Body(c, rows[row][col])
					}).Height(140)

					ds.HStack(c, 2, func() {
						ds.Spinner(c, "Loading tokens", ds.SizeSM)
						ds.Meter(c, "Usage", 75, 0, 100, &ui.MeterLevels{Warning: 70, Critical: 90})
					})

					ds.HStack(c, 2, func() {
						ds.NumberInput(c, &numVal, 0, 100, 1)
						ds.Stepper(c, &stepVal, 0, 10, 1, "Retries")
						ds.DateInput(c, &due)
					})
					ds.RangeSlider(c, &lowVal, &highVal, 0, 100, 5, "Contrast Range")
					ds.TokenField(c, &tokens, []string{"wcag", "oklab"})
				})
			})

			if ds.AlertDialog(c, &alertOpen, "Confirm Reset?", "This resets values.", "Confirm", "Cancel", ds.ToneDanger) {
				alertConfirmed = true
			}
		})
	}, 980, 860)

	// DropdownMenu open + click item
	if err := tester.Click("Menu Actions"); err != nil {
		t.Fatalf("failed to open DropdownMenu: %v", err)
	}
	if err := tester.Click("Duplicate"); err != nil || menuAction != "duplicate" {
		t.Fatalf("expected menuAction=duplicate, err=%v action=%q", err, menuAction)
	}

	// Collapsible toggle
	if err := tester.Click("Advanced Settings"); err != nil || !collOpen {
		t.Fatalf("expected collOpen=true: err=%v", err)
	}

	// AlertDialog confirm
	if err := tester.Click("Trigger Alert"); err != nil || !alertOpen {
		t.Fatalf("expected alertOpen=true: err=%v", err)
	}
	if err := tester.Click("Confirm"); err != nil || !alertConfirmed || alertOpen {
		t.Fatalf("expected AlertDialog confirm to return true and close: err=%v confirmed=%v open=%v", err, alertConfirmed, alertOpen)
	}
}


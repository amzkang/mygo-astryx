package main

import (
	"fmt"
	"strings"

	"github.com/egoist/mygo/ui"
	"github.com/amzkang/mygo-astryx/ds"
	"github.com/amzkang/mygo-astryx/icons"
	"github.com/amzkang/mygo-astryx/registry"
	"github.com/amzkang/mygo-astryx/theme"
)

func playgroundVariant(idx int) (ds.Variant, string) {
	switch idx {
	case 1:
		return ds.VariantSoft, "ds.VariantSoft"
	case 2:
		return ds.VariantOutline, "ds.VariantOutline"
	case 3:
		return ds.VariantGhost, "ds.VariantGhost"
	default:
		return ds.VariantSolid, "ds.VariantSolid"
	}
}

func playgroundSize(idx int) (ds.Size, string) {
	switch idx {
	case 0:
		return ds.SizeXS, "ds.SizeXS"
	case 1:
		return ds.SizeSM, "ds.SizeSM"
	case 3:
		return ds.SizeLG, "ds.SizeLG"
	default:
		return ds.SizeMD, "ds.SizeMD"
	}
}

func playgroundTone(idx int) (ds.Tone, string) {
	switch idx {
	case 0:
		return ds.ToneNeutral, "ds.ToneNeutral"
	case 2:
		return ds.ToneSuccess, "ds.ToneSuccess"
	case 3:
		return ds.ToneWarning, "ds.ToneWarning"
	case 4:
		return ds.ToneDanger, "ds.ToneDanger"
	case 5:
		return ds.ToneInfo, "ds.ToneInfo"
	default:
		return ds.ToneAccent, "ds.ToneAccent"
	}
}

func generateSnippet(comp registry.ComponentDoc, app *galleryState) string {
	_, variantSym := playgroundVariant(app.pgVariantIdx)
	_, sizeSym := playgroundSize(app.pgSizeIdx)
	_, toneSym := playgroundTone(app.pgToneIdx)
	label := strings.TrimSpace(app.pgLabel)
	if label == "" {
		label = "Deploy Release"
	}

	switch comp.Name {
	case "Button", "PrimaryButton":
		fn := "ds.Button"
		if comp.Name == "PrimaryButton" {
			fn = "ds.PrimaryButton"
		}
		var opts []string
		if comp.Name == "Button" {
			opts = append(opts, "Variant: "+variantSym)
		}
		opts = append(opts, "Tone: "+toneSym, "Size: "+sizeSym)
		if app.pgDisabled {
			opts = append(opts, "Disabled: true")
		}
		if app.pgLoading {
			opts = append(opts, "Loading: true")
		}
		return fmt.Sprintf(`if %s(c, %q, ds.ButtonOpts{%s}).Clicked() {
	c.Toast("Triggered %s")
}`, fn, label, strings.Join(opts, ", "), label)

	case "IconButton":
		return fmt.Sprintf(`ds.IconButton(c, %q, icons.Settings, ds.IconButtonOpts{
	Variant:  %s,
	Tone:     %s,
	Size:     %s,
	Disabled: %v,
})`, label, variantSym, toneSym, sizeSym, app.pgDisabled)

	case "CopyButton":
		return fmt.Sprintf(`ds.CopyButton(c, %q, "go run ./cmd/astryx", ds.ButtonOpts{
	Variant: %s,
	Tone:    %s,
	Size:    %s,
})`, label, variantSym, toneSym, sizeSym)

	case "Badge":
		return fmt.Sprintf(`ds.Badge(c, %q, ds.BadgeOpts{
	Variant: %s,
	Tone:    %s,
	Dot:     %v,
})`, label, variantSym, toneSym, app.pgDot)

	case "Tag":
		return fmt.Sprintf(`ds.Tag(c, %q, ds.TagOpts{
	Tone: %s,
})`, label, toneSym)

	case "Alert":
		return fmt.Sprintf(`ds.Alert(c, "Contrast and tokens verified for active theme.", ds.AlertOpts{
	Title: %q,
	Tone:  %s,
})`, label, toneSym)

	case "Spinner":
		return fmt.Sprintf(`ds.Spinner(c, %q, %s)`, label, sizeSym)

	case "Progress":
		return fmt.Sprintf(`ds.Progress(c, %.2f, %s)`, app.playSlider/100.0, toneSym)

	case "TextInput":
		return fmt.Sprintf(`ds.TextInput(c, &value, ds.InputOpts{
	Placeholder: %q,
	Size:        %s,
	Disabled:    %v,
})`, label, sizeSym, app.pgDisabled)

	case "StatTile":
		return fmt.Sprintf(`ds.StatTile(c, %q, "99.8%%", "+4.2%%", %s)`, label, toneSym)

	case "AlertDialog":
		return fmt.Sprintf(`if ds.AlertDialog(c, &open, %q, "This action cannot be undone.", "Confirm", "Cancel", %s) {
	// confirmed
}`, label, toneSym)

	case "DropdownMenu":
		return fmt.Sprintf(`ds.DropdownMenu(c, %q, &menuOpen, []ds.MenuItemSpec{
	{Label: "Inspect", Icon: icons.Eye, OnClick: func() {}},
	{Separator: true},
	{Label: "Delete", Danger: true, OnClick: func() {}},
}, ds.ButtonOpts{Variant: %s, Tone: %s, Size: %s})`, label, variantSym, toneSym, sizeSym)

	default:
		return comp.Example
	}
}

func renderComponentPlayground(c *ui.Context, app *galleryState, th *theme.Theme, comp registry.ComponentDoc) {
	variant, _ := playgroundVariant(app.pgVariantIdx)
	size, _ := playgroundSize(app.pgSizeIdx)
	tone, _ := playgroundTone(app.pgToneIdx)
	label := strings.TrimSpace(app.pgLabel)
	if label == "" {
		label = "Deploy Release"
	}

	ds.Card(c, func() {
		ds.HStack(c, 2, func() {
			ds.Heading(c, 3, "Interactive Playground").Grow(1)
			ds.Badge(c, comp.Category, ds.BadgeOpts{Tone: tone, Dot: app.pgDot})
		})

		// Prop Controls
		ds.VStack(c, 2.5, func() {
			ds.HStack(c, 3, func() {
				ui.Column(c).Gap(th.Space(1)).Children(func() {
					ds.Caption(c, "Variant")
					ds.SegmentedControl(c, &app.pgVariantIdx, []string{"Solid", "Soft", "Outline", "Ghost"})
				})
				ui.Column(c).Gap(th.Space(1)).Children(func() {
					ds.Caption(c, "Size")
					ds.SegmentedControl(c, &app.pgSizeIdx, []string{"XS", "SM", "MD", "LG"})
				})
			})
			ds.HStack(c, 3, func() {
				ui.Column(c).Gap(th.Space(1)).Children(func() {
					ds.Caption(c, "Tone")
					ds.SegmentedControl(c, &app.pgToneIdx, []string{"Neutral", "Accent", "Success", "Warning", "Danger", "Info"})
				})
			})
			ds.HStack(c, 3, func() {
				ds.TextInput(c, &app.pgLabel, ds.InputOpts{Placeholder: "Preview label", Size: ds.SizeSM}).Width(200)
				ds.Switch(c, &app.pgDisabled, "Disabled")
				ds.Switch(c, &app.pgLoading, "Loading")
				ds.Switch(c, &app.pgDot, "Dot")
			})
		})

		// Live Preview Canvas
		ui.Box(c).
			FillWidth().
			Padding(th.Space(5)).
			Radius(th.Radius.MD).
			Background(th.Colors.SurfaceSubtle).
			Border(1, th.Colors.BorderSubtle).
			Center().
			Children(func() {
				renderPlaygroundWidget(c, app, th, comp, label, variant, size, tone)
			})

		// Dynamic Snippet Output
		snippet := generateSnippet(comp, app)
		ds.HStack(c, 2, func() {
			ds.Caption(c, "Live Go Snippet (updates with props above)").Grow(1)
			ds.CopyButton(c, "Copy Live Snippet", snippet, ds.ButtonOpts{Size: ds.SizeXS})
		})
		ui.Box(c).
			FillWidth().
			Padding(th.Space(3)).
			Radius(th.Radius.MD).
			Background(th.Colors.SurfaceSubtle).
			Border(1, th.Colors.BorderSubtle).
			Children(func() {
				ui.Text(c, snippet).Font(th.MonoFont).FontSize(th.Typography.Code.Size).TextColor(th.Colors.Text)
			})
	}, ds.CardOpts{Elevated: true})
}

func renderPlaygroundWidget(c *ui.Context, app *galleryState, th *theme.Theme, comp registry.ComponentDoc, label string, variant ds.Variant, size ds.Size, tone ds.Tone) {
	switch comp.Name {
	case "Button":
		if ds.Button(c, label, ds.ButtonOpts{
			Variant:     variant,
			Tone:        tone,
			Size:        size,
			Disabled:    app.pgDisabled,
			Loading:     app.pgLoading,
			LeadingIcon: icons.Sparkles,
		}).Clicked() {
			c.Toast("Clicked " + label)
		}
	case "PrimaryButton":
		if ds.PrimaryButton(c, label, ds.ButtonOpts{
			Tone:        tone,
			Size:        size,
			Disabled:    app.pgDisabled,
			Loading:     app.pgLoading,
			LeadingIcon: icons.Check,
		}).Clicked() {
			c.Toast("PrimaryButton clicked")
		}
	case "IconButton":
		if ds.IconButton(c, label, icons.Settings, ds.IconButtonOpts{
			Variant:  variant,
			Tone:     tone,
			Size:     size,
			Disabled: app.pgDisabled,
		}).Clicked() {
			c.Toast("IconButton clicked")
		}
	case "CopyButton":
		ds.CopyButton(c, label, "go run ./cmd/astryx", ds.ButtonOpts{
			Variant:  variant,
			Tone:     tone,
			Size:     size,
			Disabled: app.pgDisabled,
		})
	case "Badge":
		ds.Badge(c, label, ds.BadgeOpts{Variant: variant, Tone: tone, Dot: app.pgDot})
	case "Tag":
		ds.Tag(c, label, ds.TagOpts{Tone: tone})
	case "Alert":
		ds.Alert(c, "Contrast and tokens verified for active theme.", ds.AlertOpts{Title: label, Tone: tone})
	case "EmptyState":
		ds.EmptyState(c, label, "Try adjusting filters or creating a new item.", ds.EmptyStateOpts{Icon: icons.Inbox})
	case "Skeleton":
		ds.HStack(c, 2, func() {
			ds.Skeleton(c, 40, 40, th.Radius.Full)
			ds.VStack(c, 1.5, func() {
				ds.Skeleton(c, 180, 14, th.Radius.SM)
				ds.Skeleton(c, 120, 12, th.Radius.SM)
			})
		})
	case "Spinner":
		ds.HStack(c, 2, func() {
			ds.Spinner(c, label, size)
			ds.Body(c, label)
		})
	case "Meter":
		ds.Meter(c, label, app.playSlider, 0, 100, &ui.MeterLevels{Warning: 70, Critical: 90})
	case "Progress":
		ds.Progress(c, app.playSlider/100.0, tone)
	case "StatTile":
		ds.StatTile(c, label, "99.8%", "+4.2%", tone).Width(240)
	case "Avatar":
		ds.Avatar(c, label, 44, nil)
	case "KeyValueRow":
		ds.KeyValueRow(c, label, th.Name).Width(320)
	case "TextInput":
		ds.TextInput(c, &app.playInput, ds.InputOpts{Placeholder: label, Size: size, Disabled: app.pgDisabled, LeadingIcon: icons.Edit}).Width(280)
	case "TextArea":
		ds.TextArea(c, &app.playNotes, label, 84).Width(320)
	case "SearchInput":
		ds.SearchInput(c, &app.compSearch, label).Width(280)
	case "Checkbox":
		ds.Checkbox(c, &app.playChecked, label)
	case "Switch":
		ds.Switch(c, &app.playSwitch, label)
	case "Radio":
		ds.HStack(c, 3, func() {
			ds.Radio(c.Key("pg-rad-a"), &app.playRadio, "solid", label+" A")
			ds.Radio(c.Key("pg-rad-b"), &app.playRadio, "soft", label+" B")
		})
	case "SegmentedControl":
		ds.SegmentedControl(c, &app.playSeg, []string{"Overview", "Tokens", "Code"})
	case "Select":
		ds.Select(c, &app.playSelect, []string{"Matcha", "Neutral", "Gothic", "Y2K"}, nil)
	case "Combobox":
		ds.Combobox(c, &app.playCombo, []string{"Button", "DropdownMenu", "AlertDialog", "DataTable"}, label)
	case "Slider":
		ds.Slider(c, &app.playSlider, 0, 100, label).Width(320)
	case "NumberInput":
		ds.NumberInput(c, &app.pgNumVal, 0, 100, 1)
	case "Stepper":
		ds.Stepper(c, &app.pgStepVal, 0, 20, 1, label)
	case "RangeSlider":
		ds.RangeSlider(c, &app.pgRangeLo, &app.pgRangeHi, 0, 100, 5, label).Width(320)
	case "TokenField":
		ds.TokenField(c, &app.pgTokens, []string{"wcag", "oklab", "native", "swizzle"}).Width(360)
	case "DateInput":
		ds.DateInput(c, &app.pgDate)
	case "DropdownMenu":
		ds.DropdownMenu(c, label, &app.pgMenuOpen, []ds.MenuItemSpec{
			{Label: "Inspect", Icon: icons.Eye, OnClick: func() { c.Toast("Inspect clicked") }},
			{Label: "Duplicate", Icon: icons.Copy, OnClick: func() { c.Toast("Duplicate clicked") }},
			{Separator: true},
			{Label: "Delete", Danger: true, OnClick: func() { c.Toast("Delete clicked") }},
		}, ds.ButtonOpts{Variant: variant, Tone: tone, Size: size, Disabled: app.pgDisabled})
	case "AlertDialog":
		if ds.Button(c, "Open "+label, ds.ButtonOpts{Variant: variant, Tone: tone, Size: size}).Clicked() {
			app.pgAlertOpen = true
		}
		if ds.AlertDialog(c, &app.pgAlertOpen, label, "This demonstrates AlertDialog with your active tone and theme.", "Confirm", "Cancel", tone) {
			c.Toast("Confirmed " + label)
		}
	case "Collapsible":
		ds.Collapsible(c, label, &app.pgCollOpen, func() {
			ds.Body(c, "Collapsible content expands and animates cleanly within the frame.")
		}).Width(360)
	case "DataTable":
		cols := []ui.TableColumn{{Title: "Theme", Width: 140}, {Title: "Text/Bg", Width: 110}, {Title: "Status", Width: 100}}
		rows := [][]string{{"Neutral", "15.4:1", "AAA"}, {"Matcha", "12.8:1", "AAA"}, {"Gothic", "14.2:1", "AAA"}}
		ds.DataTable(c, &app.pgTableState, cols, len(rows), "No rows", func(row, col int) {
			ds.Body(c, rows[row][col])
		}).Width(380).Height(150)
	case "TreeView":
		ds.TreeView(c, []ds.TreeNode{
			{
				ID:    "themes",
				Label: "Stock Themes",
				Icon:  icons.Palette,
				Open:  &app.pgTreeOpen,
				Children: []ds.TreeNode{
					{ID: "matcha", Label: "Matcha", Icon: icons.Sparkles},
					{ID: "gothic", Label: "Gothic", Icon: icons.Moon},
				},
			},
		}, &app.pgTreeSel).Width(280)
	case "Toolbar":
		ds.Toolbar(c, func() {
			ds.Button(c, label, ds.ButtonOpts{Variant: variant, Tone: tone, Size: ds.SizeSM, LeadingIcon: icons.Plus})
			ds.Button(c, "Export", ds.ButtonOpts{Variant: ds.VariantOutline, Size: ds.SizeSM, LeadingIcon: icons.Download})
		}).Width(360)
	case "Tabs":
		ds.Tabs(c, &app.playTab, []string{"Overview", "Props", "Code"}, ds.TabsOpts{Pills: true})
	case "Breadcrumbs":
		crumb := -1
		ds.Breadcrumbs(c, []string{"Home", "Components", label}, &crumb)
	case "Pagination":
		ds.Pagination(c, &app.playPage, 7)
	case "Accordion":
		ds.Accordion(c, []ds.AccordionItemSpec{
			{Title: label, Open: &app.playAccOpen1, Body: func() { ds.Body(c, "Expanded inside playground.") }},
		}).Width(360)
	default:
		ds.HStack(c, 2, func() {
			ds.Badge(c, comp.Name, ds.BadgeOpts{Variant: variant, Tone: tone, Dot: app.pgDot})
			ds.Button(c, label, ds.ButtonOpts{Variant: variant, Tone: tone, Size: size, Disabled: app.pgDisabled, Loading: app.pgLoading})
		})
	}
}

func renderDesktopControlsShowcase(c *ui.Context, app *galleryState, th *theme.Theme) {
	cols := []ui.TableColumn{
		{Title: "Theme", Width: 130},
		{Title: "Accent", Width: 120},
		{Title: "WCAG", Width: 90},
	}
	rows := [][]string{
		{"Matcha", "Botanical", "AAA"},
		{"Gothic", "Obsidian", "AAA"},
		{"Y2K", "Chrome", "AAA"},
		{"Butter", "Amber", "AAA"},
	}

	ds.VStack(c, 3, func() {
		ds.Toolbar(c, func() {
			ds.DropdownMenu(c, "Actions", &app.pgMenuOpen, []ds.MenuItemSpec{
				{Label: "Inspect Theme", Icon: icons.Palette, OnClick: func() { app.page = pageTokens }},
				{Label: "Browse Icons", Icon: icons.Box, OnClick: func() { app.page = pageIcons }},
				{Separator: true},
				{Label: "Reset Demo", Danger: true, OnClick: func() { app.pgAlertOpen = true }},
			}, ds.ButtonOpts{Size: ds.SizeSM})
			if ds.Button(c, "Confirm Dialog", ds.ButtonOpts{Variant: ds.VariantSoft, Tone: ds.ToneDanger, Size: ds.SizeSM}).Clicked() {
				app.pgAlertOpen = true
			}
			ui.Spacer(c)
			ds.Spinner(c, "Syncing", ds.SizeSM)
			ds.Badge(c, "v0.8.3", ds.BadgeOpts{Tone: ds.ToneInfo})
		})

		ds.SplitPane(c, &app.pgSplit, false, func() {
			ds.TreeView(c, []ds.TreeNode{
				{
					ID:    "root",
					Label: "mygo-astryx",
					Icon:  icons.Layers,
					Open:  &app.pgTreeOpen,
					Children: []ds.TreeNode{
						{ID: "ds", Label: "ds (61)", Icon: icons.Box},
						{ID: "theme", Label: "theme (7)", Icon: icons.Palette},
						{ID: "icons", Label: "icons (56)", Icon: icons.Sparkles},
					},
				},
			}, &app.pgTreeSel)
		}, func() {
			ds.VStack(c, 2.5, func() {
				ds.DataTable(c, &app.pgTableState, cols, len(rows), "Empty", func(row, col int) {
					ds.Body(c, rows[row][col])
				}).Height(140)
				ds.Meter(c, "Token Budget", app.playSlider, 0, 100, &ui.MeterLevels{Warning: 70, Critical: 90})
				ds.Collapsible(c, "Advanced Inputs (Number / Stepper / Range / Tokens / Date)", &app.pgCollOpen, func() {
					ds.VStack(c, 2, func() {
						ds.HStack(c, 2, func() {
							ds.NumberInput(c, &app.pgNumVal, 0, 100, 1)
							ds.Stepper(c, &app.pgStepVal, 0, 20, 1, "Retries")
							ds.DateInput(c, &app.pgDate)
						})
						ds.RangeSlider(c, &app.pgRangeLo, &app.pgRangeHi, 0, 100, 5, "Contrast Window")
						ds.TokenField(c, &app.pgTokens, []string{"wcag", "oklab", "native", "swizzle"})
					})
				})
			}).Padding(0, 0, 0, th.Space(3))
		}).Height(320)

		if ds.AlertDialog(c, &app.pgAlertOpen, "Reset Playground State?", "This restores default playground sample values.", "Reset", "Cancel", ds.ToneDanger) {
			app.pgNumVal = 12
			app.pgStepVal = 3
			app.pgRangeLo, app.pgRangeHi = 20, 80
			c.Toast("Playground defaults restored")
		}
	})
}

func renderThemeContrastMatrix(c *ui.Context, app *galleryState, th *theme.Theme) {
	ds.VStack(c, 2.5, func() {
		ds.Caption(c, fmt.Sprintf("Comparing all 7 stock themes in %s mode against WCAG AA (>=4.5:1) and AAA (>=7.0:1).", map[bool]string{false: "light", true: "dark"}[app.dark]))
		for _, ts := range theme.StockThemes() {
			target := ts.For(app.dark)
			name := target.Name
			textBg := theme.ContrastRatio(target.Colors.Text, target.Colors.Background)
			accentBtn := theme.ContrastRatio(target.Colors.Accent.OnColor, target.Colors.Accent.Base)
			accentText := theme.ContrastRatio(target.Colors.Accent.Text, target.Colors.Background)
			dangerBtn := theme.ContrastRatio(target.Colors.Danger.OnColor, target.Colors.Danger.Base)

			isCurrent := strings.EqualFold(app.themeName, name)
			ds.Card(c.Key("matrix-"+name), func() {
				ds.HStack(c, 3, func() {
					ui.Row(c).Gap(th.Space(1)).AlignItems(ui.Center).Children(func() {
						ui.Box(c).Size(24, 24).Radius(th.Radius.SM).Background(target.Colors.Background).Border(1, target.Colors.Border)
						ui.Box(c).Size(24, 24).Radius(th.Radius.SM).Background(target.Colors.Accent.Base)
						ui.Box(c).Size(24, 24).Radius(th.Radius.SM).Background(target.Colors.Text)
					})
					ui.Column(c).Width(110).Gap(1).Children(func() {
						ui.Text(c, strings.ToUpper(name[:1])+name[1:]).FontWeight(600).TextColor(th.Colors.Text)
						if isCurrent {
							ds.Caption(c, "Active theme")
						} else {
							ds.Caption(c, map[bool]string{false: "Light", true: "Dark"}[app.dark])
						}
					})
					contrastPill(c, "Text/Bg", textBg).Grow(1)
					contrastPill(c, "Accent Btn", accentBtn).Grow(1)
					contrastPill(c, "Accent Text", accentText).Grow(1)
					contrastPill(c, "Danger Btn", dangerBtn).Grow(1)
					btnLabel := "Use " + strings.ToUpper(name[:1]) + name[1:]
					variant := ds.VariantOutline
					if isCurrent {
						variant = ds.VariantSoft
					}
					if ds.Button(c.Key("use-theme-"+name), btnLabel, ds.ButtonOpts{Variant: variant, Size: ds.SizeSM}).Clicked() {
						app.themeName = name
					}
				})
			}, ds.CardOpts{Padding: th.Space(3)})
		}
	})
}

func contrastPill(c *ui.Context, label string, ratio float64) ui.Element {
	grade := "AA"
	tone := ds.ToneSuccess
	if ratio >= 7.0 {
		grade = "AAA"
	} else if ratio < 4.5 {
		grade = "Fail"
		tone = ds.ToneDanger
	}
	return ui.Column(c).Gap(2).Children(func() {
		ds.Caption(c, label)
		ds.Badge(c, fmt.Sprintf("%.2f:1 %s", ratio, grade), ds.BadgeOpts{Tone: tone, Dot: true})
	})
}

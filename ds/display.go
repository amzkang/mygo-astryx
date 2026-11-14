package ds

import (
	"fmt"

	"github.com/egoist/mygo/ui"
	"github.com/amzkang/mygo-astryx/icons"
	"github.com/amzkang/mygo-astryx/theme"
)

// AccordionItemSpec defines one section inside Accordion.
type AccordionItemSpec struct {
	Title string
	Open  *bool
	Body  func()
}

// Accordion renders a themed collapsible stack built on ui.AccordionBase.
func Accordion(c *ui.Context, items []AccordionItemSpec) ui.Element {
	th := theme.Of(c)
	acc := ui.AccordionBase(c)
	acc.Root.
		FillWidth().
		Radius(th.Radius.LG).
		Background(th.Colors.Surface).
		Border(1, th.Colors.Border).
		Clip()
	acc.Root.Children(func() {
		for i, item := range items {
			if i > 0 {
				Divider(c.Key(fmt.Sprintf("acc-div-%d", i)))
			}
			sec := acc.Item(item.Open)
			sec.Trigger.
				FillWidth().
				AlignItems(ui.Center).
				Justify(ui.SpaceBetween).
				Padding(th.Space(3), th.Space(4)).
				Cursor(ui.CursorPointer)
			if sec.Trigger.Hovered() {
				sec.Trigger.Background(th.Colors.SurfaceHover)
			}
			sec.Trigger.Children(func() {
				ui.Text(c, item.Title).
					FontSize(th.Typography.Body.Size).
					FontWeight(600).
					TextColor(th.Colors.Text)
				chev := icons.ChevronRight
				if *item.Open {
					chev = icons.ChevronDown
				}
				ui.Icon(c, chev).Size(16, 16).TextColor(th.Colors.TextMuted)
			})
			panel := sec.Panel(func() {
				if item.Body != nil {
					item.Body()
				}
			})
			if panel != (ui.Element{}) {
				panel.
					FillWidth().
					Padding(0, th.Space(4), th.Space(3.5), th.Space(4)).
					Gap(th.Space(2))
			}
		}
	})
	return acc.Root
}

// Popover renders an anchored floating panel built on ui.PopoverBase.
func Popover(c *ui.Context, anchor ui.Element, open *bool, fn func()) ui.Element {
	th := theme.Of(c)
	return ui.PopoverBase(c, anchor, open, func(panel ui.Element) {
		panel.
			Margin(th.Space(1.5), 0, 0, 0).
			Padding(th.Space(3.5)).
			Gap(th.Space(2.5)).
			Radius(th.Radius.LG).
			Background(th.Colors.SurfaceElevated).
			Border(1, th.Colors.Border)
		th.Shadows.LG.Apply(panel)
		if fn != nil {
			fn()
		}
	})
}

// Tooltip attaches a themed tooltip to anchor built on ui.TooltipBase.
func Tooltip(c *ui.Context, anchor ui.Element, text string) ui.Element {
	if text == "" {
		return anchor
	}
	th := theme.Of(c)
	anchor.Description(text)
	ui.TooltipBase(c, anchor, func(tip ui.Element) {
		tip.
			MaxWidth(280).
			Padding(th.Space(1.25), th.Space(2.25)).
			Radius(th.Radius.SM).
			Background(th.Colors.Inverse).
			TextColor(th.Colors.InverseText).
			AttachTo(anchor, ui.AnchorTop, ui.AnchorBottom).
			Margin(0, 0, th.Space(1.5), 0)
		th.Shadows.MD.Apply(tip)
		tip.Children(func() {
			ui.Text(c, text).
				FontSize(th.Typography.Caption.Size).
				TextColor(th.Colors.InverseText)
		})
	})
	return anchor
}

// StatTile renders a compact KPI / metric card with optional delta badge.
func StatTile(c *ui.Context, label, value, delta string, tone ...Tone) ui.Element {
	th := theme.Of(c)
	t := ToneAccent
	if len(tone) > 0 {
		t = tone[0]
	}
	return Card(c, func() {
		ui.Row(c).FillWidth().AlignItems(ui.Center).Justify(ui.SpaceBetween).Gap(th.Space(2)).Children(func() {
			Caption(c, label)
			if delta != "" {
				Badge(c, delta, BadgeOpts{Variant: VariantSoft, Tone: t})
			}
		})
		ui.Text(c, value).
			FontSize(th.Typography.H1.Size).
			FontWeight(700).
			TextColor(th.Colors.Text)
	}, CardOpts{Padding: th.Space(3.5)})
}

// Avatar renders a circular user/brand badge using initials when bitmap is nil.
func Avatar(c *ui.Context, name string, size float32, img *ui.Bitmap) ui.Element {
	th := theme.Of(c)
	if size <= 0 {
		size = 36
	}
	box := ui.Box(c).
		Size(size, size).
		Center().
		Radius(size / 2).
		Background(th.Colors.Accent.Subtle).
		Border(1, th.Colors.Accent.Border).
		Clip().
		Shrink(0)
	if img != nil {
		return box.Children(func() {
			ui.Image(c, img).Size(size, size).Fit(ui.Cover)
		})
	}
	initials := "?"
	if name != "" {
		runes := []rune(name)
		initials = string(runes[0])
		for i := 1; i < len(runes); i++ {
			if runes[i-1] == ' ' && runes[i] != ' ' {
				initials += string(runes[i])
				break
			}
		}
	}
	return box.Children(func() {
		ui.Text(c, initials).
			FontSize(max(11, size*0.38)).
			FontWeight(700).
			TextColor(th.Colors.Accent.Text).
			SingleLine()
	})
}

// KeyValueRow renders a two-column label/value inspection row.
func KeyValueRow(c *ui.Context, key, value string) ui.Element {
	th := theme.Of(c)
	return ui.Row(c).
		FillWidth().
		AlignItems(ui.Center).
		Justify(ui.SpaceBetween).
		Gap(th.Space(3)).
		Padding(th.Space(1.75), 0).
		Children(func() {
			Caption(c, key)
			CodeText(c, value)
		})
}

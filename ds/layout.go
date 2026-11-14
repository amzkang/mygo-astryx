package ds

import (
	"github.com/egoist/mygo/ui"
	"github.com/amzkang/mygo-astryx/theme"
)

// VStack lays children out vertically with theme-scaled gap.
func VStack(c *ui.Context, gapUnits float32, fn func()) ui.Element {
	th := theme.Of(c)
	return ui.Column(c).Gap(th.Space(gapUnits)).Children(fn)
}

// HStack lays children out horizontally, vertically centered, with theme-scaled gap.
func HStack(c *ui.Context, gapUnits float32, fn func()) ui.Element {
	th := theme.Of(c)
	return ui.Row(c).AlignItems(ui.Center).Gap(th.Space(gapUnits)).Children(fn)
}

// Center centers its children inside available space.
func Center(c *ui.Context, fn func()) ui.Element {
	return ui.Column(c).Fill().Center().Children(fn)
}

// Container constrains content width to maxWidth DIPs and centers it horizontally.
func Container(c *ui.Context, maxWidth float32, fn func()) ui.Element {
	th := theme.Of(c)
	return ui.Column(c).FillWidth().AlignItems(ui.Center).PaddingX(th.Space(4)).Children(func() {
		ui.Column(c).FillWidth().MaxWidth(maxWidth).Children(fn)
	})
}

// CardOpts configures Card surface elevation and padding.
type CardOpts struct {
	Elevated bool
	Padding  float32 // 0 uses theme default (4 space units)
	Radius   float32 // 0 uses theme.Radius.LG
}

// Card creates a bordered, rounded surface container.
func Card(c *ui.Context, fn func(), opts ...CardOpts) ui.Element {
	th := theme.Of(c)
	var o CardOpts
	if len(opts) > 0 {
		o = opts[0]
	}
	pad := o.Padding
	if pad <= 0 {
		pad = th.Space(4)
	}
	rad := o.Radius
	if rad <= 0 {
		rad = th.Radius.LG
	}
	bg := th.Colors.Surface
	if o.Elevated {
		bg = th.Colors.SurfaceElevated
	}
	el := ui.Column(c).
		Padding(pad).
		Gap(th.Space(3)).
		Radius(rad).
		Background(bg).
		Border(1, th.Colors.Border)
	if o.Elevated {
		el = th.Shadows.MD.Apply(el)
	} else {
		el = th.Shadows.SM.Apply(el)
	}
	return el.Children(fn)
}

// Section creates a semantic region with a heading and optional subtitle.
func Section(c *ui.Context, title, subtitle string, fn func()) ui.Element {
	th := theme.Of(c)
	return ui.Column(c).Gap(th.Space(3)).Role(ui.RoleRegion).Label(title).Children(func() {
		ui.Column(c).Gap(th.Space(1)).Children(func() {
			if title != "" {
				Heading(c, 2, title)
			}
			if subtitle != "" {
				ui.Text(c, subtitle).
					FontSize(th.Typography.BodySM.Size).
					TextColor(th.Colors.TextMuted)
			}
		})
		if fn != nil {
			fn()
		}
	})
}

// Heading renders a semantic heading with level 1..3 mapped to theme.TypographyScale.
func Heading(c *ui.Context, level int, text string) ui.Element {
	th := theme.Of(c)
	var ts theme.TextStyle
	switch level {
	case 1:
		ts = th.Typography.H1
	case 3:
		ts = th.Typography.H3
	default:
		level = 2
		ts = th.Typography.H2
	}
	el := ui.Text(c, text).
		FontSize(ts.Size).
		FontWeight(ts.Weight).
		TextColor(th.Colors.Text).
		Role(ui.RoleHeading).
		Level(level)
	if ts.Font != "" {
		el.Font(ts.Font)
	}
	return el
}

// Body renders primary body copy using theme.Typography.Body.
func Body(c *ui.Context, text string) ui.Element {
	th := theme.Of(c)
	ts := th.Typography.Body
	return ui.Text(c, text).FontSize(ts.Size).FontWeight(ts.Weight).TextColor(th.Colors.Text)
}

// Caption renders muted secondary metadata text using theme.Typography.Caption.
func Caption(c *ui.Context, text string) ui.Element {
	th := theme.Of(c)
	ts := th.Typography.Caption
	return ui.Text(c, text).FontSize(ts.Size).FontWeight(ts.Weight).TextColor(th.Colors.TextMuted)
}

// CodeText renders inline monospace code with subtle surface pill styling.
func CodeText(c *ui.Context, code string) ui.Element {
	th := theme.Of(c)
	ts := th.Typography.Code
	return ui.Box(c).
		Padding(th.Space(0.5), th.Space(1.5)).
		Radius(th.Radius.XS).
		Background(th.Colors.SurfaceSubtle).
		Border(1, th.Colors.BorderSubtle).
		Children(func() {
			ui.Text(c, code).Font(ts.Font).FontSize(ts.Size).TextColor(th.Colors.Text)
		})
}

// Kbd renders a keyboard shortcut cap (e.g. "Cmd+K").
func Kbd(c *ui.Context, shortcut string) ui.Element {
	th := theme.Of(c)
	return ui.Box(c).
		Padding(th.Space(0.5), th.Space(1.5)).
		Radius(th.Radius.XS).
		Background(th.Colors.SurfaceElevated).
		Border(1, th.Colors.BorderStrong).
		Children(func() {
			ui.Text(c, shortcut).
				Font(th.MonoFont).
				FontSize(th.Typography.Caption.Size).
				FontWeight(600).
				TextColor(th.Colors.TextSecondary).
				SingleLine()
		})
}

// SplitPane renders a resizable horizontal or vertical split view built on ui.Split / ui.SplitVertical.
func SplitPane(c *ui.Context, size *float32, vertical bool, first, second func()) ui.Element {
	if vertical {
		return ui.SplitVertical(c, size, first, second).Fill()
	}
	return ui.Split(c, size, first, second).Fill()
}

// Toolbar renders a horizontal action bar with subtle surface background and bottom border.
func Toolbar(c *ui.Context, fn func()) ui.Element {
	th := theme.Of(c)
	return ui.Toolbar(c, fn).
		FillWidth().
		AlignItems(ui.Center).
		Gap(th.Space(2)).
		Padding(th.Space(2), th.Space(3)).
		Radius(th.Radius.MD).
		Background(th.Colors.SurfaceSubtle).
		Border(1, th.Colors.BorderSubtle)
}


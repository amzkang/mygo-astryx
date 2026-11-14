package ds

import (
	"math"

	"github.com/amzkang/mygo-astryx/theme"
	"github.com/egoist/mygo/ui"
)

// LayoutSlots builds optional regions in the current frame. Start and End
// follow the inherited writing direction. Empty slots occupy no space.
type LayoutSlots struct {
	Header, Start, Content, End, Footer func()
}

// LayoutOpts sizes a five-region layout. Dimensions are DIPs, zero uses defaults.
// AutoHeight lets the layout size to content instead of filling its parent.
type LayoutOpts struct {
	StartWidth, EndWidth float32
	Gap, Padding         float32
	AutoHeight           bool
}

// Layout composes fixed-height headers/footers around a flexible content row.
func Layout(c *ui.Context, slots LayoutSlots, opts ...LayoutOpts) ui.Element {
	th := theme.Of(c)
	var o LayoutOpts
	if len(opts) > 0 {
		o = opts[0]
	}
	startWidth, endWidth := o.StartWidth, o.EndWidth
	if startWidth <= 0 {
		startWidth = th.Space(60)
	}
	if endWidth <= 0 {
		endWidth = th.Space(80)
	}
	root := ui.Column(c).MinWidth(0).MinHeight(0).Gap(max(o.Gap, 0)).Padding(max(o.Padding, 0))
	if o.AutoHeight {
		root.FillWidth()
	} else {
		root.Fill()
	}
	return root.Children(func() {
		if slots.Header != nil {
			ui.Column(c.Key("header")).FillWidth().Shrink(0).Children(slots.Header)
		}
		if slots.Start != nil || slots.Content != nil || slots.End != nil {
			row := ui.Row(c.Key("body")).FillWidth().MinHeight(0).AlignItems(ui.Stretch).Gap(max(o.Gap, 0))
			if !o.AutoHeight {
				row.Grow(1)
			}
			row.Children(func() {
				if slots.Start != nil {
					ui.Column(c.Key("start")).Width(startWidth).MinHeight(0).Shrink(0).Children(slots.Start)
				}
				if slots.Content != nil {
					ui.Column(c.Key("content")).Grow(1).MinWidth(0).MinHeight(0).Children(slots.Content)
				}
				if slots.End != nil {
					ui.Column(c.Key("end")).Width(endWidth).MinHeight(0).Shrink(0).Children(slots.End)
				}
			})
		}
		if slots.Footer != nil {
			ui.Column(c.Key("footer")).FillWidth().Shrink(0).Children(slots.Footer)
		}
	})
}

// GridOpts sets equal columns or explicit tracks, and density-scaled default gaps.
// ColumnTracks takes precedence over Columns. Gap is in DIPs, zero is four units.
type GridOpts struct {
	Columns      int
	ColumnTracks []ui.Track
	Gap          float32
}

// Grid lays out content row by row using MyGo's native grid and inherited RTL.
func Grid(c *ui.Context, fn func(), opts ...GridOpts) ui.Element {
	var o GridOpts
	if len(opts) > 0 {
		o = opts[0]
	}
	gap := o.Gap
	if gap <= 0 {
		gap = theme.Of(c).Space(4)
	}
	if o.Columns <= 0 {
		o.Columns = 2
	}
	e := ui.Grid(c).FillWidth().MinWidth(0).Columns(o.Columns).Gap(gap)
	if len(o.ColumnTracks) > 0 {
		e.ColumnTracks(o.ColumnTracks...)
	}
	return e.Children(fn)
}

// AspectRatio constrains a full-width media container to width/height ratio.
// Invalid or nonpositive ratios use a square. Children inherit its clipping.
func AspectRatio(c *ui.Context, ratio float32, fn func()) ui.Element {
	if ratio <= 0 || math.IsNaN(float64(ratio)) || math.IsInf(float64(ratio), 0) {
		ratio = 1
	}
	return ui.Box(c).FillWidth().MinWidth(0).AspectRatio(ratio).Clip().Children(fn)
}

// Link renders a themed, keyboard-accessible native URL link. The host opens URLs.
func Link(c *ui.Context, label, url string) ui.Element {
	th := theme.Of(c)
	return ui.Link(c, label, url).TextColor(th.Colors.Accent.Text).FontSize(th.Typography.Body.Size)
}

// ScrollOpts configures ScrollableArea. State optionally persists the scroll
// position in caller state. Both takes precedence over Horizontal.
type ScrollOpts struct {
	Label       string
	Horizontal  bool
	Both        bool
	State       *ui.ScrollState
	HideShadows bool
}

// ScrollableArea adds a named keyboard-scrollable region and overflow edge fades.
func ScrollableArea(c *ui.Context, fn func(), opts ...ScrollOpts) ui.Element {
	th := theme.Of(c)
	var o ScrollOpts
	if len(opts) > 0 {
		o = opts[0]
	}
	var e ui.Element
	if o.Both {
		e = ui.ScrollBoth(c)
	} else if o.Horizontal {
		e = ui.ScrollHorizontal(c)
	} else {
		e = ui.Scroll(c)
	}
	label := o.Label
	if label == "" {
		label = "Scrollable content"
	}
	e.Fill().MinWidth(0).MinHeight(0).Role(ui.RoleRegion).Label(label).Focusable()
	state := o.State
	if state == nil {
		state = ui.Local(e, "ds_scroll", func() ui.ScrollState { return ui.ScrollState{} })
	}
	e.TrackScroll(state)
	// Only consume keys while the region itself is focused. Inputs inside it
	// retain their arrow, Home/End and PageUp/PageDown behavior.
	if e.Focused() {
		step := th.Space(8)
		if e.Shortcut(0, ui.KeyDown) {
			state.Y += step
		}
		if e.Shortcut(0, ui.KeyUp) {
			state.Y -= step
		}
		if e.Shortcut(0, ui.KeyPageDown) {
			state.Y += max(step, e.Bounds().H-step)
		}
		if e.Shortcut(0, ui.KeyPageUp) {
			state.Y -= max(step, e.Bounds().H-step)
		}
		if e.Shortcut(0, ui.KeyEnd) {
			state.Y = state.MaxY
			state.X = state.MaxX
		}
		if e.Shortcut(0, ui.KeyHome) {
			state.X, state.Y = 0, 0
		}
		if c.Direction() == ui.RightToLeft {
			step = -step
		}
		if e.Shortcut(0, ui.KeyRight) {
			state.X += step
		}
		if e.Shortcut(0, ui.KeyLeft) {
			state.X -= step
		}
	}
	// Apply requests before laying out this frame's children.
	e.TrackScroll(state).Children(fn)
	if !o.HideShadows {
		color := th.Colors.Text.Alpha(0.12)
		rtl := c.Direction() == ui.RightToLeft
		e.DrawOver(func(p *ui.Painter, r ui.Rect) {
			fade := min(th.Space(3), min(r.W, r.H)/2)
			if state.Y > 0 {
				p.FillGradient(ui.Rect{X: r.X, Y: r.Y, W: r.W, H: fade}, ui.LinearGradient{From: color, To: color.Alpha(0), Angle: 180}, 0)
			}
			if state.Y < state.MaxY {
				p.FillGradient(ui.Rect{X: r.X, Y: r.Y + r.H - fade, W: r.W, H: fade}, ui.LinearGradient{From: color.Alpha(0), To: color, Angle: 180}, 0)
			}
			left, right := state.X > 0, state.X < state.MaxX
			if rtl {
				left, right = right, left
			}
			if left {
				p.FillGradient(ui.Rect{X: r.X, Y: r.Y, W: fade, H: r.H}, ui.LinearGradient{From: color, To: color.Alpha(0), Angle: 90}, 0)
			}
			if right {
				p.FillGradient(ui.Rect{X: r.X + r.W - fade, Y: r.Y, W: fade, H: r.H}, ui.LinearGradient{From: color.Alpha(0), To: color, Angle: 90}, 0)
			}
		})
	}
	return e
}

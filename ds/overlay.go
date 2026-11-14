package ds

import (
	"github.com/egoist/mygo/ui"
	"github.com/amzkang/mygo-astryx/icons"
	"github.com/amzkang/mygo-astryx/theme"
)

// DrawerSide controls which edge Drawer slides out from.
type DrawerSide uint8

const (
	DrawerRight DrawerSide = iota
	DrawerLeft
)

// DrawerOpts configures Drawer dimensions and behavior.
type DrawerOpts struct {
	Side  DrawerSide
	Width float32 // Default 380
	Title string
}

// Drawer renders an accessible slide-over panel using ui.Overlay, Modal focus trapping,
// Escape dismissal, outside-click dismissal, and ElementTransition motion.
func Drawer(c *ui.Context, open *bool, fn func(), opts ...DrawerOpts) ui.Element {
	th := theme.Of(c)
	var o DrawerOpts
	if len(opts) > 0 {
		o = opts[0]
	}
	w := o.Width
	if w <= 0 {
		w = 380
	}
	enterX := float32(28)
	if o.Side == DrawerLeft {
		enterX = -28
	}

	var out ui.Element
	ui.Overlay(c.Key("ds_drawer_overlay"), func() {
		if open == nil || !*open {
			return
		}
		winW, winH := c.Size()
		panelW := min(w, max(0, winW-24))

		backdrop := ui.Box(c.Key("drawer_backdrop")).
			Absolute().Top(0).Left(0).Right(0).Bottom(0).
			Background(th.Colors.Backdrop).
			Transition(ui.ElementTransition{
				Animation: &th.Motion.Normal,
				Enter:     &ui.Motion{},
				Exit:      &ui.Motion{},
			})
		if backdrop.Clicked() {
			*open = false
		}

		panel := ui.Column(c.Key("drawer_panel")).
			Absolute().Top(0).Height(winH).Width(panelW).
			Background(th.Colors.SurfaceElevated).
			Role(ui.RoleDialog).
			Modal().
			Label(o.Title).
			Transition(ui.ElementTransition{
				Animation: &th.Motion.Normal,
				Enter:     &ui.Motion{X: enterX},
				Exit:      &ui.Motion{X: enterX},
			})
		if o.Side == DrawerLeft {
			panel.Left(0).BorderWidth(0, 1, 0, 0).BorderColor(th.Colors.Border)
		} else {
			panel.Right(0).BorderWidth(0, 0, 0, 1).BorderColor(th.Colors.Border)
		}
		th.Shadows.LG.Apply(panel)

		if panel.OverlayShortcut(0, ui.KeyEscape) || panel.PressedOutside() {
			*open = false
		}

		panel.Children(func() {
			ui.Row(c).
				FillWidth().
				AlignItems(ui.Center).
				Padding(th.Space(3.5), th.Space(4)).
				BorderWidth(0, 0, 1, 0).
				BorderColor(th.Colors.BorderSubtle).
				Children(func() {
					Heading(c, 3, o.Title).Grow(1)
					if IconButton(c, "Close drawer", icons.X, IconButtonOpts{Variant: VariantGhost, Size: SizeSM}).Clicked() {
						*open = false
					}
				})
			ui.Scroll(c).Grow(1).MinHeight(0).Padding(th.Space(4)).Children(fn)
		})
		out = panel
	})
	return out
}

// ModalCard renders a themed centered modal dialog built on ui.DialogBase.
func ModalCard(c *ui.Context, open *bool, title string, width float32, fn func()) ui.Element {
	th := theme.Of(c)
	if width <= 0 {
		width = 440
	}
	return ui.DialogBase(c, open, func(backdrop, panel ui.Element) {
		backdrop.Background(th.Colors.Backdrop)
		panel.
			Width(width).
			MaxWidthPercent(92).
			Padding(th.Space(5)).
			Gap(th.Space(4)).
			Radius(th.Radius.XL).
			Background(th.Colors.SurfaceElevated).
			Border(1, th.Colors.Border).
			Label(title)
		th.Shadows.LG.Apply(panel)

		ui.Row(c).FillWidth().AlignItems(ui.Center).Children(func() {
			Heading(c, 2, title).Grow(1)
			if IconButton(c, "Close dialog", icons.X, IconButtonOpts{Variant: VariantGhost, Size: SizeSM}).Clicked() {
				*open = false
			}
		})
		if fn != nil {
			fn()
		}
	})
}

package ds

import (
	"fmt"

	"github.com/egoist/mygo/ui"
	"github.com/amzkang/mygo-astryx/icons"
	"github.com/amzkang/mygo-astryx/theme"
)

// MenuItemSpec describes a single entry inside DropdownMenu or ContextMenu.
type MenuItemSpec struct {
	ID        string
	Label     string
	Icon      *ui.SVG
	Shortcut  string
	Checked   bool
	Checkable bool
	Disabled  bool
	Danger    bool
	Separator bool
	OnClick   func()
}

// DropdownMenu renders a themed menu button built on ui.MenuBase.
func DropdownMenu(c *ui.Context, triggerLabel string, open *bool, items []MenuItemSpec, opts ...ButtonOpts) ui.Element {
	th := theme.Of(c)
	var o ButtonOpts
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Variant == 0 && o.Tone == 0 {
		o.Variant = VariantOutline
	}
	h, padX, _, fontSize, iconSize := sizeMetrics(th, o.Size)
	m := ui.MenuBase(c, open)
	m.Trigger.
		Height(h).
		Padding(0, padX).
		Gap(th.Space(1.5)).
		Radius(th.Radius.MD).
		Cursor(ui.CursorPointer)
	if o.Disabled {
		m.Trigger.Disabled(true)
	}
	bg, hover, pressed, fg, border := buttonColors(th, o.Variant, o.Tone)
	applyInteractiveFill(m.Trigger, bg, hover, pressed)
	m.Trigger.TextColor(fg)
	if border.A > 0 {
		m.Trigger.Border(1, border)
	}
	m.Trigger.Children(func() {
		if o.LeadingIcon != nil {
			ui.Icon(c, o.LeadingIcon).Size(iconSize, iconSize).TextColor(fg)
		}
		if triggerLabel != "" {
			ui.Text(c, triggerLabel).FontSize(fontSize).FontWeight(500).TextColor(fg).SingleLine()
		}
		ui.Icon(c, icons.ChevronDown).Size(14, 14).TextColor(fg)
	})

	renderMenuPopup(c, th, m, items)
	return m.Trigger
}

// ContextMenu attaches a themed context menu (secondary-click / Shift+F10) to target via ui.ContextMenuBase.
func ContextMenu(c *ui.Context, target ui.Element, open *bool, items []MenuItemSpec) ui.Element {
	th := theme.Of(c)
	m := ui.ContextMenuBase(c, target, open)
	renderMenuPopup(c, th, m, items)
	return target
}

func renderMenuPopup(c *ui.Context, th *theme.Theme, m ui.MenuParts, items []MenuItemSpec) {
	m.Popup(func(panel ui.Element) {
		panel.
			MinWidth(190).
			Margin(th.Space(1), 0, 0, 0).
			Padding(th.Space(1)).
			Gap(th.Space(0.5)).
			Radius(th.Radius.LG).
			Background(th.Colors.SurfaceElevated).
			Border(1, th.Colors.Border)
		th.Shadows.LG.Apply(panel)

		for i, spec := range items {
			if spec.Separator {
				ui.Box(c.Key(fmt.Sprintf("menu-sep-%d", i))).
					FillWidth().
					Height(1).
					Margin(th.Space(1), th.Space(1.5)).
					Background(th.Colors.BorderSubtle)
				continue
			}
			var item ui.Element
			if spec.Checkable {
				item = m.CheckItem(spec.Checked)
			} else {
				item = m.Item()
			}
			item.
				FillWidth().
				AlignItems(ui.Center).
				Justify(ui.SpaceBetween).
				Gap(th.Space(3)).
				Padding(th.Space(1.75), th.Space(2.5)).
				Radius(th.Radius.SM).
				Cursor(ui.CursorPointer)
			if spec.Disabled {
				item.Disabled(true)
			}

			fg := th.Colors.Text
			if spec.Danger {
				fg = th.Colors.Danger.Text
			}
			if spec.Disabled {
				fg = th.Colors.TextMuted
			} else if item.Highlighted() {
				if spec.Danger {
					item.Background(th.Colors.Danger.Subtle)
				} else {
					item.Background(th.Colors.Accent.Subtle)
					fg = th.Colors.Accent.Text
				}
			}

			item.Children(func() {
				ui.Row(c).AlignItems(ui.Center).Gap(th.Space(2)).Children(func() {
					if spec.Checkable {
						if spec.Checked {
							ui.Icon(c, icons.Check).Size(14, 14).TextColor(fg)
						} else {
							ui.Box(c).Size(14, 14)
						}
					} else if spec.Icon != nil {
						ui.Icon(c, spec.Icon).Size(14, 14).TextColor(fg)
					}
					ui.Text(c, spec.Label).FontSize(th.Typography.BodySM.Size).FontWeight(500).TextColor(fg).SingleLine()
				})
				if spec.Shortcut != "" {
					ui.Text(c, spec.Shortcut).Font(th.MonoFont).FontSize(th.Typography.Caption.Size).TextColor(th.Colors.TextMuted).SingleLine()
				}
			})

			if item.Clicked() && spec.OnClick != nil {
				spec.OnClick()
			}
		}
	})
}

// AlertDialog renders a themed confirmation dialog with Cancel and Confirm buttons built on ui.DialogBase.
// It returns true in the frame the user confirms the action.
func AlertDialog(c *ui.Context, open *bool, title, message, confirmLabel, cancelLabel string, tone Tone) bool {
	th := theme.Of(c)
	if confirmLabel == "" {
		confirmLabel = "Confirm"
	}
	if cancelLabel == "" {
		cancelLabel = "Cancel"
	}
	if tone == ToneNeutral {
		tone = ToneAccent
	}
	confirmed := false
	ui.DialogBase(c, open, func(backdrop, panel ui.Element) {
		backdrop.Background(th.Colors.Backdrop)
		panel.
			Width(420).
			MaxWidthPercent(92).
			Padding(th.Space(5)).
			Gap(th.Space(3.5)).
			Radius(th.Radius.XL).
			Background(th.Colors.SurfaceElevated).
			Border(1, th.Colors.Border).
			Role(ui.RoleAlertDialog).
			Label(title).
			Description(message)
		th.Shadows.LG.Apply(panel)

		ui.Row(c).FillWidth().AlignItems(ui.Start).Gap(th.Space(3)).Children(func() {
			sem := resolveTone(th, tone)
			icon := icons.Info
			if tone == ToneDanger {
				icon = icons.AlertCircle
			} else if tone == ToneWarning {
				icon = icons.AlertTriangle
			}
			ui.Box(c).
				Size(34, 34).
				Center().
				Radius(th.Radius.Full).
				Background(sem.Subtle).
				Shrink(0).
				Children(func() {
					ui.Icon(c, icon).Size(18, 18).TextColor(sem.Text)
				})
			ui.Column(c).Grow(1).MinWidth(0).Gap(th.Space(1.5)).Children(func() {
				Heading(c, 3, title)
				if message != "" {
					ui.Text(c, message).FontSize(th.Typography.BodySM.Size).TextColor(th.Colors.TextSecondary)
				}
			})
		})

		ui.Row(c).FillWidth().Justify(ui.End).Gap(th.Space(2)).Children(func() {
			if Button(c, cancelLabel, ButtonOpts{Variant: VariantOutline}).Clicked() {
				*open = false
			}
			if Button(c, confirmLabel, ButtonOpts{Variant: VariantSolid, Tone: tone}).Clicked() {
				*open = false
				confirmed = true
			}
		})
	})
	return confirmed
}

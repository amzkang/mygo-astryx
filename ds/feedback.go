package ds

import (
	"github.com/egoist/mygo/ui"
	"github.com/amzkang/mygo-astryx/icons"
	"github.com/amzkang/mygo-astryx/theme"
)

// BadgeOpts configures Badge styling.
type BadgeOpts struct {
	Variant Variant
	Tone    Tone
	Dot     bool
}

// Badge renders a semantic pill indicator.
func Badge(c *ui.Context, text string, opts ...BadgeOpts) ui.Element {
	th := theme.Of(c)
	var o BadgeOpts
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Variant == 0 {
		o.Variant = VariantSoft
	}
	sem := resolveTone(th, o.Tone)
	bg, fg, bdr := sem.Subtle, sem.Text, ui.Transparent
	switch o.Variant {
	case VariantSolid:
		bg, fg = sem.Base, sem.OnColor
		if o.Tone == ToneNeutral {
			bdr = th.Colors.Border
		}
	case VariantOutline:
		bg, fg, bdr = ui.Transparent, sem.Text, sem.Border
		if o.Tone == ToneNeutral {
			bdr = th.Colors.BorderStrong
		}
	default: // VariantSoft
		if o.Tone == ToneNeutral {
			bg, fg = th.Colors.SurfaceSubtle, th.Colors.TextSecondary
		}
	}
	el := ui.Row(c).
		AlignItems(ui.Center).
		Gap(th.Space(1.25)).
		Padding(th.Space(0.5), th.Space(2.25)).
		Radius(th.Radius.Full).
		Background(bg).
		Shrink(0)
	if bdr.A > 0 {
		el.Border(1, bdr)
	}
	return el.Children(func() {
		if o.Dot {
			dotColor := sem.Base
			if o.Variant == VariantSolid && o.Tone != ToneNeutral {
				dotColor = sem.OnColor
			} else if o.Tone == ToneNeutral {
				dotColor = th.Colors.TextMuted
			}
			ui.Box(c).Size(6, 6).Radius(3).Background(dotColor)
		}
		ui.Text(c, text).
			FontSize(th.Typography.Caption.Size).
			FontWeight(600).
			TextColor(fg).
			SingleLine()
	})
}

// TagOpts configures an interactive or removable Tag chip.
type TagOpts struct {
	Tone      Tone
	Removable bool
	OnRemove  func()
}

// Tag renders a compact label chip with an optional dismiss button.
func Tag(c *ui.Context, label string, opts ...TagOpts) ui.Element {
	th := theme.Of(c)
	var o TagOpts
	if len(opts) > 0 {
		o = opts[0]
	}
	sem := resolveTone(th, o.Tone)
	bg := th.Colors.SurfaceSubtle
	fg := th.Colors.Text
	bdr := th.Colors.Border
	if o.Tone != ToneNeutral {
		bg, fg, bdr = sem.Subtle, sem.Text, sem.Border
	}
	return ui.Row(c).
		AlignItems(ui.Center).
		Gap(th.Space(1)).
		Padding(th.Space(0.75), th.Space(2)).
		Radius(th.Radius.SM).
		Background(bg).
		Border(1, bdr).
		Shrink(0).
		Children(func() {
			ui.Text(c, label).FontSize(th.Typography.BodySM.Size).FontWeight(500).TextColor(fg).SingleLine()
			if o.Removable {
				rm := ui.ButtonBase(c).
					Label("Remove " + label).
					Size(16, 16).
					Center().
					Radius(th.Radius.Full).
					Cursor(ui.CursorPointer)
				if rm.Hovered() {
					rm.Background(th.Colors.SurfaceHover)
				}
				rm.Children(func() {
					ui.Icon(c, icons.X).Size(11, 11).TextColor(fg)
				})
				if rm.Clicked() && o.OnRemove != nil {
					o.OnRemove()
				}
			}
		})
}

// AlertOpts configures inline Alert / Callout banners.
type AlertOpts struct {
	Tone    Tone
	Title   string
	Icon    *ui.SVG
	Action  string
	OnClick func()
}

// Alert renders a semantic callout banner with appropriate ARIA role (RoleAlert for Warning/Danger, RoleStatus otherwise).
func Alert(c *ui.Context, message string, opts ...AlertOpts) ui.Element {
	th := theme.Of(c)
	var o AlertOpts
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Tone == ToneNeutral {
		o.Tone = ToneInfo
	}
	sem := resolveTone(th, o.Tone)
	icon := o.Icon
	if icon == nil {
		switch o.Tone {
		case ToneSuccess:
			icon = icons.CheckCircle
		case ToneWarning:
			icon = icons.AlertTriangle
		case ToneDanger:
			icon = icons.AlertCircle
		default:
			icon = icons.Info
		}
	}
	role := ui.RoleStatus
	if o.Tone == ToneWarning || o.Tone == ToneDanger {
		role = ui.RoleAlert
	}
	return ui.Row(c).
		FillWidth().
		AlignItems(ui.Start).
		Gap(th.Space(3)).
		Padding(th.Space(3), th.Space(3.5)).
		Radius(th.Radius.MD).
		Background(sem.Subtle).
		Border(1, sem.Border).
		Role(role).
		Children(func() {
			ui.Box(c).Padding(2, 0, 0, 0).Shrink(0).Children(func() {
				ui.Icon(c, icon).Size(17, 17).TextColor(sem.Text)
			})
			ui.Column(c).Grow(1).MinWidth(0).Gap(th.Space(1)).Children(func() {
				if o.Title != "" {
					ui.Text(c, o.Title).
						FontSize(th.Typography.Body.Size).
						FontWeight(600).
						TextColor(th.Colors.Text)
				}
				if message != "" {
					ui.Text(c, message).
						FontSize(th.Typography.BodySM.Size).
						TextColor(th.Colors.TextSecondary)
				}
			})
			if o.Action != "" {
				btn := Button(c, o.Action, ButtonOpts{Variant: VariantOutline, Size: SizeSM, Tone: o.Tone})
				if btn.Clicked() && o.OnClick != nil {
					o.OnClick()
				}
			}
		})
}

// EmptyStateOpts configures EmptyState placeholders.
type EmptyStateOpts struct {
	Icon        *ui.SVG
	ActionLabel string
	OnAction    func()
}

// EmptyState renders a centered placeholder with icon, title, description, and optional action button.
func EmptyState(c *ui.Context, title, description string, opts ...EmptyStateOpts) ui.Element {
	th := theme.Of(c)
	var o EmptyStateOpts
	if len(opts) > 0 {
		o = opts[0]
	}
	icon := o.Icon
	if icon == nil {
		icon = icons.Inbox
	}
	return ui.Column(c).
		FillWidth().
		MinHeight(180).
		Center().
		Padding(th.Space(6)).
		Gap(th.Space(2.5)).
		Radius(th.Radius.LG).
		Border(1, th.Colors.BorderSubtle).
		Background(th.Colors.SurfaceSubtle).
		Children(func() {
			ui.Box(c).
				Size(44, 44).
				Center().
				Radius(th.Radius.Full).
				Background(th.Colors.Surface).
				Border(1, th.Colors.Border).
				Children(func() {
					ui.Icon(c, icon).Size(22, 22).TextColor(th.Colors.TextMuted)
				})
			ui.Text(c, title).FontSize(th.Typography.H3.Size).FontWeight(600).TextColor(th.Colors.Text).TextAlign(ui.Center)
			if description != "" {
				ui.Text(c, description).
					MaxWidth(360).
					FontSize(th.Typography.BodySM.Size).
					TextColor(th.Colors.TextMuted).
					TextAlign(ui.Center)
			}
			if o.ActionLabel != "" {
				if PrimaryButton(c, o.ActionLabel, ButtonOpts{Size: SizeSM}).Clicked() && o.OnAction != nil {
					o.OnAction()
				}
			}
		})
}

// Skeleton renders a pulsing placeholder rectangle that respects ReduceMotion.
func Skeleton(c *ui.Context, width, height, radius float32) ui.Element {
	th := theme.Of(c)
	if radius <= 0 {
		radius = th.Radius.SM
	}
	el := ui.Box(c).
		Size(width, height).
		Radius(radius).
		Role(ui.RoleStatus).
		Busy(true).
		Label("Loading")
	bg := th.Colors.Surface
	if !c.Preferences().ReduceMotion {
		pulse := el.Loop("skeleton-pulse", 1400000000, ui.Bounce(ui.EaseInOut))
		bg = th.Colors.Surface.MixOklab(th.Colors.SurfaceHover, pulse)
	}
	return el.Background(bg)
}

// ToastViewport renders window toasts using ui.ToastViewportBase and semantic theme tones.
func ToastViewport(c *ui.Context) ui.Element {
	th := theme.Of(c)
	return ui.ToastViewportBase(c, func(viewport ui.Element, toasts []ui.Toast) {
		viewport.Padding(th.Space(4)).AlignItems(ui.End).Gap(th.Space(2))
		for _, t := range toasts {
			toast := ui.ToastBase(c.Key(t.ID), t)
			tone := ToneNeutral
			var icon *ui.SVG
			switch t.Type {
			case "error", "danger":
				tone, icon = ToneDanger, icons.AlertCircle
			case "warning":
				tone, icon = ToneWarning, icons.AlertTriangle
			case "success":
				tone, icon = ToneSuccess, icons.CheckCircle
			case "info":
				tone, icon = ToneInfo, icons.Info
			}
			sem := resolveTone(th, tone)
			bdr := th.Colors.Border
			if tone != ToneNeutral {
				bdr = sem.Border
			}
			root := toast.Root.
				Row().
				AlignItems(ui.Center).
				Gap(th.Space(2.5)).
				MinWidth(260).
				MaxWidth(400).
				Padding(th.Space(2.5), th.Space(3.5)).
				Radius(th.Radius.MD).
				Background(th.Colors.SurfaceElevated).
				Border(1, bdr).
				Transition(ui.ElementTransition{
					Enter: &ui.Motion{Y: 10},
					Exit:  &ui.Motion{},
				})
			th.Shadows.MD.Apply(root)
			root.Children(func() {
				if icon != nil {
					ui.Icon(c, icon).Size(16, 16).TextColor(sem.Text)
				}
				ui.Column(c).Grow(1).MinWidth(0).Gap(2).Children(func() {
					if t.Title != "" {
						ui.Text(c, t.Title).FontSize(th.Typography.BodySM.Size).FontWeight(600).TextColor(th.Colors.Text)
					}
					if t.Description != "" {
						ui.Text(c, t.Description).FontSize(th.Typography.Caption.Size).TextColor(th.Colors.TextMuted)
					}
				})
				if t.Action != "" {
					act := toast.ActionButton().
						Padding(th.Space(1), th.Space(2)).
						Radius(th.Radius.SM).
						Background(th.Colors.SurfaceSubtle).
						Cursor(ui.CursorPointer)
					act.Children(func() {
						ui.Text(c, t.Action).FontSize(th.Typography.Caption.Size).FontWeight(600).TextColor(th.Colors.Accent.Text)
					})
				}
				cls := toast.CloseButton().
					Label("Close").
					Size(22, 22).
					Center().
					Radius(th.Radius.SM).
					Cursor(ui.CursorPointer)
				if cls.Hovered() {
					cls.Background(th.Colors.SurfaceHover)
				}
				cls.Children(func() {
					ui.Icon(c, icons.X).Size(13, 13).TextColor(th.Colors.TextMuted)
				})
			})
		}
	})
}

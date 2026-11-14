package ds

import (
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/amzkang/mygo-astryx/icons"
	"github.com/amzkang/mygo-astryx/theme"
)

// ButtonOpts configures Button appearance and behavior.
type ButtonOpts struct {
	Variant      Variant
	Size         Size
	Tone         Tone
	LeadingIcon  *ui.SVG
	TrailingIcon *ui.SVG
	Loading      bool
	Disabled     bool
	FullWidth    bool
	Radius       float32
}

// Button renders an accessible themed push button built on ui.ButtonBase.
func Button(c *ui.Context, label string, opts ...ButtonOpts) ui.Element {
	th := theme.Of(c)
	var o ButtonOpts
	if len(opts) > 0 {
		o = opts[0]
	}
	variant := o.Variant
	tone := o.Tone
	if variant == VariantPrimary {
		variant = VariantSolid
		if tone == ToneNeutral {
			tone = ToneAccent
		}
	}
	h, padX, _, fontSize, iconSize := sizeMetrics(th, o.Size)
	rad := o.Radius
	if rad <= 0 {
		rad = th.Radius.MD
	}

	b := ui.ButtonBase(c).
		Height(h).
		Padding(0, padX).
		Gap(th.Space(1.5)).
		Radius(rad).
		Cursor(ui.CursorPointer)
	if o.FullWidth {
		b.FillWidth()
	}
	if o.Disabled || o.Loading {
		b.Disabled(true)
	}

	bg, hover, pressed, fg, border := buttonColors(th, variant, tone)
	applyInteractiveFill(b, bg, hover, pressed)
	b.TextColor(fg)
	if border.A > 0 {
		b.Border(1, border)
	}
	b.Children(func() {
		if o.Loading {
			ui.Spinner(c).Size(iconSize, iconSize).TextColor(fg)
		} else if o.LeadingIcon != nil {
			ui.Icon(c, o.LeadingIcon).Size(iconSize, iconSize).TextColor(fg)
		}
		if label != "" {
			ui.Text(c, label).FontSize(fontSize).FontWeight(500).TextColor(fg).SingleLine()
		}
		if o.TrailingIcon != nil && !o.Loading {
			ui.Icon(c, o.TrailingIcon).Size(iconSize, iconSize).TextColor(fg)
		}
	})
	return b
}

// PrimaryButton is a shorthand for an accent-filled primary Button.
func PrimaryButton(c *ui.Context, label string, opts ...ButtonOpts) ui.Element {
	var o ButtonOpts
	if len(opts) > 0 {
		o = opts[0]
	}
	o.Variant = VariantSolid
	o.Tone = ToneAccent
	return Button(c, label, o)
}

// IconButtonOpts configures IconButton.
type IconButtonOpts struct {
	Variant  Variant
	Size     Size
	Tone     Tone
	Disabled bool
	Tooltip  string
}

// IconButton renders a square icon-only button with required accessible Label and optional Tooltip.
func IconButton(c *ui.Context, label string, icon *ui.SVG, opts ...IconButtonOpts) ui.Element {
	th := theme.Of(c)
	var o IconButtonOpts
	if len(opts) > 0 {
		o = opts[0]
	}
	h, _, _, _, iconSize := sizeMetrics(th, o.Size)
	b := ui.ButtonBase(c).
		Label(label).
		Size(h, h).
		Center().
		Radius(th.Radius.MD).
		Cursor(ui.CursorPointer)
	if o.Disabled {
		b.Disabled(true)
	}
	tip := o.Tooltip
	if tip == "" {
		tip = label
	}
	if tip != "" {
		b.Tooltip(tip)
	}
	bg, hover, pressed, fg, border := buttonColors(th, o.Variant, o.Tone)
	applyInteractiveFill(b, bg, hover, pressed)
	b.TextColor(fg)
	if border.A > 0 {
		b.Border(1, border)
	}
	b.Children(func() {
		if icon != nil {
			ui.Icon(c, icon).Size(iconSize, iconSize).TextColor(fg)
		}
	})
	return b
}

type copyButtonState struct {
	copiedUntil time.Time
}

// CopyButton copies value to the system clipboard and temporarily flips to a "Copied" confirmation state.
func CopyButton(c *ui.Context, label, value string, opts ...ButtonOpts) ui.Element {
	var o ButtonOpts
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Variant == 0 && o.Tone == 0 {
		o.Variant = VariantOutline
	}
	o.LeadingIcon = icons.Copy
	btn := Button(c, label, o)
	st := ui.Local(btn, "copy-state", func() copyButtonState { return copyButtonState{} })
	copied := !st.copiedUntil.IsZero() && c.Now().Before(st.copiedUntil)
	if copied {
		c.After(st.copiedUntil.Sub(c.Now()))
	}
	if btn.Clicked() {
		c.WriteClipboardThen(value, func(err error) {
			if err == nil {
				c.Announce("Copied to clipboard")
			}
		})
		st.copiedUntil = c.Now().Add(1800 * time.Millisecond)
		c.Invalidate()
	}
	if copied {
		o.LeadingIcon = icons.ClipboardCheck
		o.Tone = ToneSuccess
		if o.Variant == VariantOutline {
			o.Variant = VariantSoft
		}
		copiedLabel := "Copied"
		if label == "" {
			copiedLabel = ""
		}
		return Button(c.Key("copied-badge"), copiedLabel, o).Label(label)
	}
	return btn
}

func buttonColors(th *theme.Theme, variant Variant, tone Tone) (base, hover, pressed, fg, border ui.Color) {
	sem := resolveTone(th, tone)
	switch variant {
	case VariantSoft:
		if tone == ToneNeutral {
			return th.Colors.SurfaceSubtle, th.Colors.SurfaceHover, th.Colors.SurfacePressed, th.Colors.Text, ui.Transparent
		}
		return sem.Subtle, sem.Subtle.MixOklab(sem.Base, 0.14), sem.Subtle.MixOklab(sem.Base, 0.24), sem.Text, ui.Transparent
	case VariantOutline:
		bdr := th.Colors.Border
		if tone != ToneNeutral {
			bdr = sem.Border
		}
		return ui.Transparent, th.Colors.SurfaceHover, th.Colors.SurfacePressed, sem.Text, bdr
	case VariantGhost:
		return ui.Transparent, th.Colors.SurfaceHover, th.Colors.SurfacePressed, sem.Text, ui.Transparent
	default: // VariantSolid
		if tone == ToneNeutral {
			return th.Colors.Surface, th.Colors.SurfaceHover, th.Colors.SurfacePressed, th.Colors.Text, th.Colors.Border
		}
		return sem.Base, sem.Hover, sem.Pressed, sem.OnColor, ui.Transparent
	}
}

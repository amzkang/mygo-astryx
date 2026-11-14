package ds

import (
	"fmt"
	"strings"
	"time"

	"github.com/amzkang/mygo-astryx/icons"
	"github.com/amzkang/mygo-astryx/theme"
	"github.com/egoist/mygo/ui"
)

// InputOpts configures TextInput.
type InputOpts struct {
	Placeholder string
	Error       string
	Disabled    bool
	LeadingIcon *ui.SVG
	Size        Size
	Password    bool
}

// TextInput renders a themed single-line text input built on ui.TextInputBase.
func TextInput(c *ui.Context, value *string, opts ...InputOpts) ui.Element {
	th := theme.Of(c)
	var o InputOpts
	if len(opts) > 0 {
		o = opts[0]
	}
	h, padX, _, fontSize, iconSize := sizeMetrics(th, o.Size)
	borderColor := th.Colors.Border
	if o.Error != "" {
		borderColor = th.Colors.Danger.Base
	}
	wrap := ui.Row(c).
		AlignItems(ui.Center).
		Gap(th.Space(2)).
		Height(h).
		Padding(0, padX).
		Radius(th.Radius.MD).
		Background(th.Colors.Surface).
		Border(1, borderColor)

	var inner ui.Element
	wrap.Children(func() {
		if o.LeadingIcon != nil {
			ui.Icon(c, o.LeadingIcon).Size(iconSize, iconSize).TextColor(th.Colors.TextMuted)
		}
		inp := ui.TextInputBase(c, value).
			Grow(1).
			FontSize(fontSize).
			TextColor(th.Colors.Text)
		if o.Placeholder != "" {
			inp.Placeholder(o.Placeholder)
		}
		if o.Disabled {
			inp.Disabled(true)
		}
		if o.Password {
			inp.Password()
		}
		if inp.Focused() {
			if o.Error != "" {
				wrap.Border(1.5, th.Colors.Danger.Base)
			} else {
				wrap.Border(1.5, th.Colors.Focus)
			}
		}
		inner = inp
	})
	_ = inner
	return wrap
}

// TextArea renders a themed multi-line text input built on ui.TextAreaBase.
func TextArea(c *ui.Context, value *string, placeholder string, height float32) ui.Element {
	th := theme.Of(c)
	if height <= 0 {
		height = 96
	}
	ta := ui.TextAreaBase(c, value).
		FillWidth().
		Height(height).
		Padding(th.Space(2.5), th.Space(3)).
		Radius(th.Radius.MD).
		Background(th.Colors.Surface).
		Border(1, th.Colors.Border).
		FontSize(th.Typography.Body.Size).
		TextColor(th.Colors.Text)
	if placeholder != "" {
		ta.Placeholder(placeholder)
	}
	if ta.Focused() {
		ta.Border(1.5, th.Colors.Focus)
	}
	return ta
}

// SearchInput renders a search field with leading search icon and a clear button when non-empty.
func SearchInput(c *ui.Context, query *string, placeholder string) ui.Element {
	th := theme.Of(c)
	if placeholder == "" {
		placeholder = "Search…"
	}
	h, padX, _, fontSize, iconSize := sizeMetrics(th, SizeMD)
	wrap := ui.Row(c).
		AlignItems(ui.Center).
		Gap(th.Space(2)).
		Height(h).
		Padding(0, padX).
		Radius(th.Radius.MD).
		Background(th.Colors.Surface).
		Border(1, th.Colors.Border)
	wrap.Children(func() {
		ui.Icon(c, icons.Search).Size(iconSize, iconSize).TextColor(th.Colors.TextMuted)
		inp := ui.TextInputBase(c, query).
			Placeholder(placeholder).
			Grow(1).
			FontSize(fontSize).
			TextColor(th.Colors.Text).
			Label(placeholder)
		if inp.Focused() {
			wrap.Border(1.5, th.Colors.Focus)
		}
		if query != nil && *query != "" {
			clr := ui.ButtonBase(c).
				Label("Clear search").
				Size(20, 20).
				Center().
				Radius(th.Radius.Full).
				Cursor(ui.CursorPointer)
			if clr.Hovered() {
				clr.Background(th.Colors.SurfaceHover)
			}
			clr.Children(func() {
				ui.Icon(c, icons.X).Size(12, 12).TextColor(th.Colors.TextMuted)
			})
			if clr.Clicked() {
				*query = ""
			}
		}
	})
	return wrap
}

// FormField wraps any control with label, optional hint, and validation error text.
func FormField(c *ui.Context, label, hint, errMsg string, fn func()) ui.Element {
	th := theme.Of(c)
	return ui.Column(c).FillWidth().Gap(th.Space(1.5)).Children(func() {
		if label != "" {
			ui.Text(c, label).
				FontSize(th.Typography.BodySM.Size).
				FontWeight(600).
				TextColor(th.Colors.Text)
		}
		if fn != nil {
			fn()
		}
		if errMsg != "" {
			ui.Row(c).AlignItems(ui.Center).Gap(th.Space(1)).Role(ui.RoleAlert).Children(func() {
				ui.Icon(c, icons.AlertCircle).Size(13, 13).TextColor(th.Colors.Danger.Text)
				ui.Text(c, errMsg).FontSize(th.Typography.Caption.Size).TextColor(th.Colors.Danger.Text)
			})
		} else if hint != "" {
			ui.Text(c, hint).FontSize(th.Typography.Caption.Size).TextColor(th.Colors.TextMuted)
		}
	})
}

// Checkbox renders a themed checkbox built on ui.CheckboxBase.
func Checkbox(c *ui.Context, checked *bool, label string) ui.Element {
	th := theme.Of(c)
	on := checked != nil && *checked
	row := ui.CheckboxBase(c, checked).
		AlignItems(ui.Center).
		Gap(th.Space(2)).
		Cursor(ui.CursorPointer).
		FocusRing(false)
	row.Children(func() {
		box := ui.Box(c).Size(18, 18).Center().Radius(th.Radius.XS).Shrink(0)
		if on {
			box.Background(th.Colors.Accent.Base)
		} else {
			bdr := th.Colors.BorderStrong
			if row.Hovered() {
				bdr = th.Colors.Accent.Base
			}
			box.Background(th.Colors.Surface).Border(1.5, bdr)
		}
		if row.FocusVisible() {
			box.Border(2, th.Colors.Focus)
		}
		box.Children(func() {
			if on {
				ui.Icon(c, icons.Check).Size(13, 13).TextColor(th.Colors.Accent.OnColor)
			}
		})
		if label != "" {
			ui.Text(c, label).FontSize(th.Typography.Body.Size).TextColor(th.Colors.Text)
		}
	})
	return row
}

// Switch renders a themed toggle switch built on ui.SwitchBase with optional label.
func Switch(c *ui.Context, on *bool, label string) ui.Element {
	th := theme.Of(c)
	active := on != nil && *on
	sw := ui.SwitchBase(c, on).
		AlignItems(ui.Center).
		Gap(th.Space(2.5)).
		Cursor(ui.CursorPointer)
	sw.Children(func() {
		track := ui.Row(c).
			Size(38, 22).
			Padding(2).
			Radius(th.Radius.Full).
			AlignItems(ui.Center)
		if active {
			track.Background(th.Colors.Accent.Base).Justify(ui.End)
		} else {
			track.Background(th.Colors.BorderStrong).Justify(ui.Start)
		}
		track.Children(func() {
			thumbBg := th.Colors.SurfaceElevated
			if th.Dark {
				thumbBg = th.Colors.Text
			}
			if active {
				thumbBg = th.Colors.Accent.OnColor
			}
			thumb := ui.Box(c).
				Size(18, 18).
				Radius(9).
				Background(thumbBg)
			th.Shadows.SM.Apply(thumb)
		})
		if label != "" {
			ui.Text(c, label).FontSize(th.Typography.Body.Size).TextColor(th.Colors.Text)
		}
	})
	return sw
}

// Radio renders a themed radio option built on ui.RadioBase.
func Radio[T comparable](c *ui.Context, selected *T, value T, label string) ui.Element {
	th := theme.Of(c)
	on := selected != nil && *selected == value
	row := ui.RadioBase(c, selected, value).
		AlignItems(ui.Center).
		Gap(th.Space(2)).
		Cursor(ui.CursorPointer)
	row.Children(func() {
		outer := ui.Box(c).Size(18, 18).Center().Radius(9).Shrink(0)
		if on {
			outer.Background(th.Colors.Accent.Base)
		} else {
			bdr := th.Colors.BorderStrong
			if row.Hovered() {
				bdr = th.Colors.Accent.Base
			}
			outer.Background(th.Colors.Surface).Border(1.5, bdr)
		}
		outer.Children(func() {
			if on {
				ui.Box(c).Size(7, 7).Radius(3.5).Background(th.Colors.Accent.OnColor)
			}
		})
		if label != "" {
			ui.Text(c, label).FontSize(th.Typography.Body.Size).TextColor(th.Colors.Text)
		}
	})
	return row
}

// SegmentedControl renders a themed segmented radio group built on ui.SegmentedBase.
func SegmentedControl(c *ui.Context, selected *int, labels []string) ui.Element {
	th := theme.Of(c)
	parts := ui.SegmentedBase(c, selected, len(labels))
	parts.Track.
		Padding(th.Space(0.75)).
		Gap(th.Space(0.75)).
		Radius(th.Radius.MD).
		Background(th.Colors.SurfaceSubtle).
		Border(1, th.Colors.Border)
	parts.Track.Children(func() {
		for i, lbl := range labels {
			on := i == *selected
			seg := parts.Segment(i).
				Padding(th.Space(1.25), th.Space(3)).
				Radius(max(0, th.Radius.MD-2)).
				Cursor(ui.CursorPointer)
			if on {
				seg.Background(th.Colors.SurfaceElevated)
				th.Shadows.SM.Apply(seg)
			} else if seg.Hovered() {
				seg.Background(th.Colors.SurfaceHover)
			}
			fg := th.Colors.TextMuted
			weight := 500
			if on {
				fg = th.Colors.Text
				weight = 600
			}
			seg.Children(func() {
				ui.Text(c, lbl).FontSize(th.Typography.BodySM.Size).FontWeight(weight).TextColor(fg).SingleLine()
			})
		}
	})
	return parts.Track
}

// Select renders a themed dropdown selector built on ui.SelectBase.
func Select[T comparable](c *ui.Context, selected *T, options []T, format func(T) string) ui.Element {
	th := theme.Of(c)
	labelOf := func(v T) string {
		if format != nil {
			return format(v)
		}
		return fmt.Sprint(v)
	}
	sel := ui.SelectBase(c, selected)
	h, padX, _, fontSize, _ := sizeMetrics(th, SizeMD)
	sel.Trigger.
		Height(h).
		MinWidth(140).
		Padding(0, padX).
		Gap(th.Space(2)).
		Justify(ui.SpaceBetween).
		Radius(th.Radius.MD).
		Background(th.Colors.Surface).
		Border(1, th.Colors.Border).
		Cursor(ui.CursorPointer)
	if sel.Trigger.Hovered() {
		sel.Trigger.Background(th.Colors.SurfaceHover)
	}
	sel.Trigger.Children(func() {
		ui.Text(c, labelOf(*selected)).FontSize(fontSize).TextColor(th.Colors.Text).SingleLine()
		ui.Icon(c, icons.ChevronDown).Size(15, 15).TextColor(th.Colors.TextMuted)
	})
	sel.Popup(func(panel ui.Element) {
		panel.
			Margin(th.Space(1), 0, 0, 0).
			Padding(th.Space(1)).
			Gap(th.Space(0.5)).
			Radius(th.Radius.MD).
			Background(th.Colors.SurfaceElevated).
			Border(1, th.Colors.Border)
		th.Shadows.MD.Apply(panel)
		for _, opt := range options {
			isSelected := opt == *selected
			item := sel.Item(opt).
				FillWidth().
				AlignItems(ui.Center).
				Justify(ui.SpaceBetween).
				Gap(th.Space(2)).
				Padding(th.Space(1.75), th.Space(2.5)).
				Radius(th.Radius.SM).
				Cursor(ui.CursorPointer)
			if item.Highlighted() {
				item.Background(th.Colors.Accent.Subtle)
			}
			item.Children(func() {
				fg := th.Colors.Text
				if item.Highlighted() {
					fg = th.Colors.Accent.Text
				}
				ui.Text(c, labelOf(opt)).FontSize(th.Typography.BodySM.Size).TextColor(fg).SingleLine()
				if isSelected {
					ui.Icon(c, icons.Check).Size(14, 14).TextColor(th.Colors.Accent.Text)
				}
			})
		}
	})
	return sel.Trigger
}

// Combobox renders a filterable text dropdown built on ui.ComboboxBase.
func Combobox(c *ui.Context, text *string, options []string, placeholder string) ui.Element {
	th := theme.Of(c)
	cb := ui.ComboboxBase(c, text)
	h, padX, _, fontSize, _ := sizeMetrics(th, SizeMD)
	cb.Input.
		Height(h).
		MinWidth(180).
		Padding(0, padX).
		Radius(th.Radius.MD).
		Background(th.Colors.Surface).
		Border(1, th.Colors.Border).
		FontSize(fontSize).
		TextColor(th.Colors.Text)
	if placeholder != "" {
		cb.Input.Placeholder(placeholder)
	}
	if cb.Input.Focused() {
		cb.Input.Border(1.5, th.Colors.Focus)
	}
	cb.Popup(func(panel ui.Element) {
		panel.
			Margin(th.Space(1), 0, 0, 0).
			Padding(th.Space(1)).
			Gap(th.Space(0.5)).
			Radius(th.Radius.MD).
			Background(th.Colors.SurfaceElevated).
			Border(1, th.Colors.Border)
		th.Shadows.MD.Apply(panel)
		q := strings.ToLower(strings.TrimSpace(*text))
		matches := 0
		for _, opt := range options {
			if cb.Filtering() && q != "" && !strings.Contains(strings.ToLower(opt), q) {
				continue
			}
			matches++
			item := cb.Item(opt).
				FillWidth().
				AlignItems(ui.Center).
				Padding(th.Space(1.75), th.Space(2.5)).
				Radius(th.Radius.SM).
				Cursor(ui.CursorPointer)
			if item.Highlighted() {
				item.Background(th.Colors.Accent.Subtle)
			}
			item.Children(func() {
				fg := th.Colors.Text
				if item.Highlighted() {
					fg = th.Colors.Accent.Text
				}
				ui.Text(c, opt).FontSize(th.Typography.BodySM.Size).TextColor(fg).SingleLine()
			})
		}
		if matches == 0 {
			ui.Box(c).Padding(th.Space(2), th.Space(3)).Children(func() {
				Caption(c, "No matching options")
			})
		}
	})
	if chosen, ok := cb.Chosen(); ok {
		*text = chosen
	}
	return cb.Input
}

// Slider wraps ui.Slider with an optional label and live formatted readout.
func Slider(c *ui.Context, value *float64, lo, hi float64, label string) ui.Element {
	th := theme.Of(c)
	return ui.Column(c).FillWidth().Gap(th.Space(1.5)).Children(func() {
		if label != "" {
			ui.Row(c).FillWidth().Justify(ui.SpaceBetween).AlignItems(ui.Center).Children(func() {
				ui.Text(c, label).FontSize(th.Typography.BodySM.Size).FontWeight(500).TextColor(th.Colors.Text)
				ui.Text(c, fmt.Sprintf("%.0f", *value)).Font(th.MonoFont).FontSize(th.Typography.Caption.Size).TextColor(th.Colors.TextMuted)
			})
		}
		ui.Slider(c, value, lo, hi).FillWidth().Label(label)
	})
}

// Progress renders a semantic-toned horizontal progress bar (0..1 or negative for indeterminate).
func Progress(c *ui.Context, value float64, tone ...Tone) ui.Element {
	th := theme.Of(c)
	t := ToneAccent
	if len(tone) > 0 {
		t = tone[0]
	}
	sem := resolveTone(th, t)
	if value < 0 {
		return ui.Progress(c, value).FillWidth()
	}
	clamped := max(0.0, min(1.0, value))
	fill := sem.Base
	if t == ToneNeutral {
		fill = th.Colors.TextSecondary
	}
	track := ui.Box(c).
		FillWidth().
		Height(th.Space(2)).
		Radius(th.Radius.Full).
		Background(th.Colors.SurfaceSubtle).
		Border(1, th.Colors.BorderSubtle).
		Clip().
		Role(ui.RoleProgress).
		Range(0, 1, clamped)
	track.Children(func() {
		ui.Box(c).
			WidthPercent(float32(clamped * 100)).
			FillHeight().
			Radius(th.Radius.Full).
			Background(fill)
	})
	return track
}

// NumberInput renders a themed numeric field with increment/decrement support built on ui.NumberInput.
func NumberInput(c *ui.Context, value *float64, lo, hi, step float64) ui.Element {
	th := theme.Of(c)
	return ui.NumberInput(c, value, lo, hi, step).
		Radius(th.Radius.MD).
		Background(th.Colors.Surface).
		Border(1, th.Colors.Border)
}

// Stepper renders a labeled spin-button stepper control built on ui.Stepper.
func Stepper(c *ui.Context, value *float64, lo, hi, step float64, label string) ui.Element {
	th := theme.Of(c)
	return ui.Row(c).AlignItems(ui.Center).Gap(th.Space(2)).Children(func() {
		if label != "" {
			ui.Text(c, label).FontSize(th.Typography.BodySM.Size).FontWeight(500).TextColor(th.Colors.Text)
		}
		if value != nil {
			ui.Text(c, fmt.Sprintf("%g", *value)).Font(th.MonoFont).FontSize(th.Typography.BodySM.Size).TextColor(th.Colors.TextSecondary)
		}
		ui.Stepper(c, value, lo, hi, step).Label(label)
	})
}

// RangeSlider renders a dual-thumb range selector with live low..high readout built on ui.RangeSlider.
func RangeSlider(c *ui.Context, low, high *float64, lo, hi, step float64, label string) ui.Element {
	th := theme.Of(c)
	return ui.Column(c).FillWidth().Gap(th.Space(1.5)).Children(func() {
		if label != "" && low != nil && high != nil {
			ui.Row(c).FillWidth().Justify(ui.SpaceBetween).AlignItems(ui.Center).Children(func() {
				ui.Text(c, label).FontSize(th.Typography.BodySM.Size).FontWeight(500).TextColor(th.Colors.Text)
				ui.Text(c, fmt.Sprintf("%.0f – %.0f", *low, *high)).Font(th.MonoFont).FontSize(th.Typography.Caption.Size).TextColor(th.Colors.TextMuted)
			})
		}
		ui.RangeSlider(c, low, high, lo, hi, step).FillWidth().Label(label)
	})
}

// TokenField renders a multi-tag autocomplete input built on ui.TokenField.
func TokenField(c *ui.Context, tokens *[]string, suggestions []string) ui.Element {
	th := theme.Of(c)
	return ui.TokenField(c, tokens, suggestions).
		FillWidth().
		Radius(th.Radius.MD).
		Background(th.Colors.Surface).
		Border(1, th.Colors.Border)
}

// DateInput renders a themed calendar date picker input built on ui.DateInput.
func DateInput(c *ui.Context, date *time.Time) ui.Element {
	th := theme.Of(c)
	return ui.DateInput(c, date).
		Radius(th.Radius.MD).
		Background(th.Colors.Surface).
		Border(1, th.Colors.Border)
}

package theme

import (
	"runtime"
	"time"

	"github.com/egoist/mygo/ui"
)

// Spec defines a theme declaratively with light/dark pairs, mirroring Astryx's defineTheme().
// Any omitted derived colors (hover, pressed, subtle, border, on-color, typography, motion)
// are automatically computed from the primary palette and mode.
type Spec struct {
	Name string

	Background      ColorPair
	Surface         ColorPair
	SurfaceSubtle   ColorPair
	SurfaceElevated ColorPair
	SurfaceHover    ColorPair
	SurfacePressed  ColorPair

	Border       ColorPair
	BorderSubtle ColorPair
	BorderStrong ColorPair

	Text          ColorPair
	TextSecondary ColorPair
	TextMuted     ColorPair

	Accent  ColorPair
	Danger  ColorPair
	Warning ColorPair
	Success ColorPair
	Info    ColorPair

	Inverse     ColorPair
	InverseText ColorPair

	Radius          float32
	Font            string
	MonoFont        string
	FontSize        float32
	Density         Density
	FocusRingWidth  float32
	FocusRingOffset float32
	TooltipDelay    time.Duration
}

// ThemeSet holds a compiled theme pair (Light and Dark) produced by DefineTheme.
type ThemeSet struct {
	Spec  Spec
	Light *Theme
	Dark  *Theme
}

// For returns the Light or Dark theme from the set.
func (s *ThemeSet) For(dark bool) *Theme {
	if dark {
		return s.Dark
	}
	return s.Light
}

// DefineTheme compiles a declarative Spec into resolved Light and Dark *Theme instances.
func DefineTheme(spec Spec) *ThemeSet {
	if spec.Name == "" {
		spec.Name = "custom"
	}
	return &ThemeSet{
		Spec:  spec,
		Light: resolveSpec(spec, false),
		Dark:  resolveSpec(spec, true),
	}
}

// WithDensity returns a shallow copy of t with adjusted spacing density.
func (t *Theme) WithDensity(d Density) *Theme {
	cp := *t
	cp.Density = d
	cp.Spacing = d.SpacingUnit()
	return &cp
}

// WithAccent returns a copy of t with a custom accent color (re-deriving hover/pressed/subtle/focus/selection).
func (t *Theme) WithAccent(accent ui.Color) *Theme {
	cp := *t
	cp.Colors.Accent = buildTone(accent, ui.Color{}, cp.Colors.Background, cp.Dark)
	cp.Colors.Selection = accent.Alpha(func() float32 {
		if cp.Dark {
			return 0.36
		}
		return 0.22
	}())
	cp.Colors.Focus = accent.Alpha(0.55)
	return &cp
}

// ToUI converts the design system Theme into a *ui.Theme so built-in mygo/ui widgets
// (Table, TextInput, Select, DateInput, etc.) automatically inherit consistent tokens.
func (t *Theme) ToUI() *ui.Theme {
	invFill, invText := t.Colors.Inverse, t.Colors.InverseText
	if invFill.A == 0 {
		invFill = t.Colors.Text
	}
	if invText.A == 0 {
		invText = t.Colors.Background
	}
	return &ui.Theme{
		Dark:            t.Dark,
		Background:      t.Colors.Background,
		Surface:         t.Colors.Surface,
		SurfaceHover:    t.Colors.SurfaceHover,
		SurfacePressed:  t.Colors.SurfacePressed,
		Border:          t.Colors.Border,
		Text:            t.Colors.Text,
		TextMuted:       t.Colors.TextMuted,
		Accent:          t.Colors.Accent.Base,
		AccentHover:     t.Colors.Accent.Hover,
		AccentPressed:   t.Colors.Accent.Pressed,
		AccentText:      t.Colors.Accent.OnColor,
		InactiveAccent:  t.Colors.InactiveAccent,
		Danger:          t.Colors.Danger.Base,
		Warning:         t.Colors.Warning.Base,
		Success:         t.Colors.Success.Base,
		Selection:       t.Colors.Selection,
		Focus:           t.Colors.Focus,
		FocusRingWidth:  t.FocusRingWidth,
		FocusRingOffset: t.FocusRingOffset,
		Inverse:         invFill,
		InverseText:     invText,
		Scrollbar:       t.Colors.Scrollbar,
		ScrollbarWidth:  t.ScrollbarWidth,
		Radius:          t.Radius.MD,
		Spacing:         t.Spacing,
		FontSize:        t.FontSize,
		Font:            t.Font,
		TooltipDelay:    t.TooltipDelay,
	}
}

func resolveSpec(spec Spec, dark bool) *Theme {
	pick := func(p ColorPair, fallbackLight, fallbackDark ui.Color) ui.Color {
		c := p.Resolve(dark)
		if c.A == 0 {
			if dark {
				return fallbackDark
			}
			return fallbackLight
		}
		return c
	}

	bg := pick(spec.Background, ui.Hex("#ffffff"), ui.Hex("#141417"))
	surface := pick(spec.Surface, ui.Hex("#f4f4f6"), ui.Hex("#1e1e23"))
	surfaceSubtle := pick(spec.SurfaceSubtle, bg.MixOklab(surface, 0.55), bg.MixOklab(surface, 0.50))
	surfaceElevated := pick(spec.SurfaceElevated, ui.Hex("#ffffff"), surface.MixOklab(ui.RGB(255, 255, 255), 0.05))
	surfaceHover := pick(spec.SurfaceHover, DeriveHover(surface, dark), DeriveHover(surface, dark))
	surfacePressed := pick(spec.SurfacePressed, DerivePressed(surface, dark), DerivePressed(surface, dark))

	border := pick(spec.Border, ui.Hex("#dcdce2"), ui.Hex("#32323a"))
	borderSubtle := pick(spec.BorderSubtle, border.MixOklab(bg, 0.45), border.MixOklab(bg, 0.45))
	borderStrong := pick(spec.BorderStrong, border.MixOklab(pick(spec.Text, ui.Hex("#18181b"), ui.Hex("#f4f4f5")), 0.28), border.MixOklab(pick(spec.Text, ui.Hex("#18181b"), ui.Hex("#f4f4f5")), 0.28))

	text := pick(spec.Text, ui.Hex("#18181b"), ui.Hex("#f4f4f6"))
	textSecondary := pick(spec.TextSecondary, text.MixOklab(bg, 0.25), text.MixOklab(bg, 0.25))
	textMuted := pick(spec.TextMuted, ui.Hex("#71717a"), ui.Hex("#a1a1aa"))

	accentBase := pick(spec.Accent, ui.Hex("#2563eb"), ui.Hex("#3b82f6"))
	dangerBase := pick(spec.Danger, ui.Hex("#dc2626"), ui.Hex("#ef4444"))
	warningBase := pick(spec.Warning, ui.Hex("#d97706"), ui.Hex("#f59e0b"))
	successBase := pick(spec.Success, ui.Hex("#16a34a"), ui.Hex("#22c55e"))
	infoBase := pick(spec.Info, ui.Hex("#0284c7"), ui.Hex("#38bdf8"))

	invFill := pick(spec.Inverse, ui.Hex("#18181b"), ui.Hex("#f4f4f6"))
	invText := pick(spec.InverseText, ui.Hex("#fafafa"), ui.Hex("#18181b"))

	radiusBase := spec.Radius
	if radiusBase <= 0 {
		radiusBase = 6
	}
	fontSize := spec.FontSize
	if fontSize <= 0 {
		if runtime.GOOS == "darwin" {
			fontSize = 13
		} else {
			fontSize = 14
		}
	}
	monoFont := spec.MonoFont
	if monoFont == "" {
		monoFont = "monospace"
	}
	focusW := spec.FocusRingWidth
	if focusW <= 0 {
		focusW = 2
	}
	focusOff := spec.FocusRingOffset
	if focusOff == 0 {
		focusOff = 1
	}
	spacing := spec.Density.SpacingUnit()

	scrollbar := ui.RGBA(0, 0, 0, 0.30)
	backdrop := ui.RGBA(0, 0, 0, 0.42)
	selAlpha := float32(0.22)
	if dark {
		scrollbar = ui.RGBA(255, 255, 255, 0.30)
		backdrop = ui.RGBA(0, 0, 0, 0.62)
		selAlpha = 0.36
	}

	shadowSM := ShadowSpec{X: 0, Y: 1, Blur: 3, Spread: 0, Color: ui.RGBA(0, 0, 0, 0.08)}
	shadowMD := ShadowSpec{X: 0, Y: 4, Blur: 12, Spread: 0, Color: ui.RGBA(0, 0, 0, 0.12)}
	shadowLG := ShadowSpec{X: 0, Y: 12, Blur: 32, Spread: 0, Color: ui.RGBA(0, 0, 0, 0.20)}
	if dark {
		shadowSM.Color = ui.RGBA(0, 0, 0, 0.32)
		shadowMD.Color = ui.RGBA(0, 0, 0, 0.45)
		shadowLG.Color = ui.RGBA(0, 0, 0, 0.60)
	}

	return &Theme{
		Name:    spec.Name,
		Dark:    dark,
		Density: spec.Density,
		Colors: Colors{
			Background:      bg,
			Surface:         surface,
			SurfaceSubtle:   surfaceSubtle,
			SurfaceElevated: surfaceElevated,
			SurfaceHover:    surfaceHover,
			SurfacePressed:  surfacePressed,
			Border:          border,
			BorderSubtle:    borderSubtle,
			BorderStrong:    borderStrong,
			Text:            text,
			TextSecondary:   textSecondary,
			TextMuted:       textMuted,
			Accent:          buildTone(accentBase, ui.Color{}, bg, dark),
			Danger:          buildTone(dangerBase, ui.Color{}, bg, dark),
			Warning:         buildTone(warningBase, ui.Color{}, bg, dark),
			Success:         buildTone(successBase, ui.Color{}, bg, dark),
			Info:            buildTone(infoBase, ui.Color{}, bg, dark),
			Selection:       accentBase.Alpha(selAlpha),
			Focus:           accentBase.Alpha(0.55),
			Inverse:         invFill,
			InverseText:     invText,
			Scrollbar:       scrollbar,
			Backdrop:        backdrop,
		},
		Typography: TypographyScale{
			Display: TextStyle{Size: fontSize * 2.0, Weight: 700, LineHeight: 1.15, Font: spec.Font},
			H1:      TextStyle{Size: fontSize * 1.55, Weight: 700, LineHeight: 1.25, Font: spec.Font},
			H2:      TextStyle{Size: fontSize * 1.28, Weight: 600, LineHeight: 1.30, Font: spec.Font},
			H3:      TextStyle{Size: fontSize * 1.08, Weight: 600, LineHeight: 1.35, Font: spec.Font},
			Body:    TextStyle{Size: fontSize, Weight: 400, LineHeight: 1.45, Font: spec.Font},
			BodySM:  TextStyle{Size: max(11, fontSize-1.5), Weight: 400, LineHeight: 1.40, Font: spec.Font},
			Caption: TextStyle{Size: max(10, fontSize-2.5), Weight: 500, LineHeight: 1.35, Font: spec.Font},
			Code:    TextStyle{Size: max(11, fontSize-1), Weight: 400, LineHeight: 1.45, Font: monoFont},
		},
		Radius: RadiusScale{
			None: 0,
			XS:   max(2, radiusBase*0.5),
			SM:   max(3, radiusBase*0.75),
			MD:   radiusBase,
			LG:   radiusBase * 1.5,
			XL:   radiusBase * 2.2,
			Full: 9999,
		},
		Shadows: ShadowScale{
			SM: shadowSM,
			MD: shadowMD,
			LG: shadowLG,
		},
		Motion: MotionTokens{
			Fast:         ui.Animation{Duration: 120 * time.Millisecond, Ease: ui.EaseOut, Oklab: true},
			Normal:       ui.Animation{Duration: 200 * time.Millisecond, Ease: ui.EaseInOut, Oklab: true},
			Slow:         ui.Animation{Duration: 320 * time.Millisecond, Ease: ui.EaseInOut, Oklab: true},
			SpringSnappy: ui.Snappy,
			SpringSmooth: ui.Smooth,
			SpringBouncy: ui.Bouncy,
		},
		Spacing:         spacing,
		FontSize:        fontSize,
		Font:            spec.Font,
		MonoFont:        monoFont,
		FocusRingWidth:  focusW,
		FocusRingOffset: focusOff,
		ScrollbarWidth:  6,
		TooltipDelay:    spec.TooltipDelay,
	}
}

func buildTone(base, onColor, bg ui.Color, dark bool) SemanticTone {
	if onColor.A == 0 {
		onColor = OnColor(base)
	}
	textColor := base
	if ContrastRatio(textColor, bg) < 4.5 {
		if dark {
			textColor = base.MixOklab(ui.RGB(255, 255, 255), 0.25)
		} else {
			textColor = base.MixOklab(ui.RGB(0, 0, 0), 0.20)
		}
	}
	return SemanticTone{
		Base:    base,
		Hover:   DeriveHover(base, dark),
		Pressed: DerivePressed(base, dark),
		Subtle:  DeriveSubtle(base, bg, dark),
		Border:  DeriveBorder(base, bg, dark),
		Text:    textColor,
		OnColor: onColor,
	}
}

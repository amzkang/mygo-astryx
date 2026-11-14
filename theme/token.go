package theme

import (
	"time"

	"github.com/egoist/mygo/ui"
)

// ColorPair stores explicit light and dark color values, inspired by Astryx's light-dark() tokens.
type ColorPair struct {
	Light ui.Color
	Dark  ui.Color
}

// Pair creates a ColorPair from light and dark colors.
func Pair(light, dark ui.Color) ColorPair {
	return ColorPair{Light: light, Dark: dark}
}

// Resolve returns the concrete color for the requested appearance.
func (p ColorPair) Resolve(dark bool) ui.Color {
	if dark {
		return p.Dark
	}
	return p.Light
}

// Density controls spacing scale and component compactness.
type Density uint8

const (
	DensityComfortable Density = iota // Spacing = 4 (default)
	DensityCompact                    // Spacing = 3
	DensitySpacious                   // Spacing = 5
)

// SpacingUnit returns the base DIP spacing unit for d.
func (d Density) SpacingUnit() float32 {
	switch d {
	case DensityCompact:
		return 3
	case DensitySpacious:
		return 5
	default:
		return 4
	}
}

// String returns the human-readable density name.
func (d Density) String() string {
	switch d {
	case DensityCompact:
		return "compact"
	case DensitySpacious:
		return "spacious"
	default:
		return "comfortable"
	}
}

// SemanticTone holds the complete state & surface family for a semantic intent.
type SemanticTone struct {
	Base    ui.Color
	Hover   ui.Color
	Pressed ui.Color
	Subtle  ui.Color
	Border  ui.Color
	Text    ui.Color // Foreground text when drawn directly on canvas/subtle surface
	OnColor ui.Color // High-contrast foreground text on top of Base
}

// TextStyle captures typography scale token properties.
type TextStyle struct {
	Size       float32
	Weight     int
	LineHeight float32
	Font       string
}

// TypographyScale defines Astryx-style hierarchical text tokens.
type TypographyScale struct {
	Display TextStyle
	H1      TextStyle
	H2      TextStyle
	H3      TextStyle
	Body    TextStyle
	BodySM  TextStyle
	Caption TextStyle
	Code    TextStyle
}

// RadiusScale defines corner radius tokens.
type RadiusScale struct {
	None float32
	XS   float32
	SM   float32
	MD   float32
	LG   float32
	XL   float32
	Full float32
}

// ShadowSpec defines a single box shadow.
type ShadowSpec struct {
	X, Y, Blur, Spread float32
	Color              ui.Color
}

// Apply attaches the shadow to an element when alpha > 0.
func (s ShadowSpec) Apply(e ui.Element) ui.Element {
	if s.Color.A == 0 {
		return e
	}
	return e.Shadow(s.X, s.Y, s.Blur, s.Spread, s.Color)
}

// ShadowScale defines elevation shadow tokens.
type ShadowScale struct {
	SM ShadowSpec
	MD ShadowSpec
	LG ShadowSpec
}

// MotionTokens centralizes animation curves and springs across the design system.
type MotionTokens struct {
	Fast         ui.Animation
	Normal       ui.Animation
	Slow         ui.Animation
	SpringSnappy ui.Animation
	SpringSmooth ui.Animation
	SpringBouncy ui.Animation
}

// Colors holds resolved concrete colors for the current light/dark mode.
type Colors struct {
	// Layered surfaces
	Background      ui.Color
	Surface         ui.Color
	SurfaceSubtle   ui.Color
	SurfaceElevated ui.Color
	SurfaceHover    ui.Color
	SurfacePressed  ui.Color

	// Borders
	Border       ui.Color
	BorderSubtle ui.Color
	BorderStrong ui.Color

	// Text hierarchy
	Text          ui.Color
	TextSecondary ui.Color
	TextMuted     ui.Color

	// Semantic families
	Accent  SemanticTone
	Danger  SemanticTone
	Warning SemanticTone
	Success SemanticTone
	Info    SemanticTone

	// Chrome & Focus
	Selection      ui.Color
	Focus          ui.Color
	Inverse        ui.Color
	InverseText    ui.Color
	Scrollbar      ui.Color
	Backdrop       ui.Color
	InactiveAccent ui.Color
}

// Theme is the full resolved mygo-astryx design system theme.
type Theme struct {
	Name       string
	Dark       bool
	Density    Density
	Colors     Colors
	Typography TypographyScale
	Radius     RadiusScale
	Shadows    ShadowScale
	Motion     MotionTokens

	// Base metrics synchronized with ui.Theme
	Spacing         float32
	FontSize        float32
	Font            string
	MonoFont        string
	FocusRingWidth  float32
	FocusRingOffset float32
	ScrollbarWidth  float32
	TooltipDelay    time.Duration
}

// Space returns n multiples of the theme's Spacing unit.
func (t *Theme) Space(n float32) float32 {
	s := t.Spacing
	if s <= 0 {
		s = t.Density.SpacingUnit()
	}
	return n * s
}

// Tone resolves one of the 5 built-in semantic families by name ("accent", "danger", "warning", "success", "info").
func (t *Theme) Tone(name string) SemanticTone {
	switch name {
	case "danger", "error":
		return t.Colors.Danger
	case "warning", "warn":
		return t.Colors.Warning
	case "success":
		return t.Colors.Success
	case "info":
		return t.Colors.Info
	default:
		return t.Colors.Accent
	}
}

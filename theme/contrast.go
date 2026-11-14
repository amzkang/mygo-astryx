package theme

import (
	"math"

	"github.com/egoist/mygo/ui"
)

// RelativeLuminance returns the WCAG 2.x relative luminance of c in [0, 1].
func RelativeLuminance(c ui.Color) float64 {
	linearize := func(v uint8) float64 {
		s := float64(v) / 255.0
		if s <= 0.04045 {
			return s / 12.92
		}
		return math.Pow((s+0.055)/1.055, 2.4)
	}
	return 0.2126*linearize(c.R) + 0.7152*linearize(c.G) + 0.0722*linearize(c.B)
}

// ContrastRatio returns the WCAG 2.x contrast ratio between a and b in [1, 21].
// Semi-transparent foregrounds are composited over bg first if b is opaque.
func ContrastRatio(fg, bg ui.Color) float64 {
	if fg.A < 255 {
		fg = fg.Over(bg)
	}
	l1 := RelativeLuminance(fg)
	l2 := RelativeLuminance(bg)
	if l1 < l2 {
		l1, l2 = l2, l1
	}
	return (l1 + 0.05) / (l2 + 0.05)
}

// OnColor picks dark or light foreground color that achieves higher contrast against bg.
func OnColor(bg ui.Color) ui.Color {
	dark := ui.Hex("#111113")
	light := ui.Hex("#ffffff")
	if ContrastRatio(light, bg) >= ContrastRatio(dark, bg) {
		return light
	}
	return dark
}

// DeriveHover shifts a color slightly toward contrast in Oklab space.
func DeriveHover(base ui.Color, dark bool) ui.Color {
	if dark {
		return base.MixOklab(ui.RGB(255, 255, 255), 0.09)
	}
	return base.MixOklab(ui.RGB(0, 0, 0), 0.08)
}

// DerivePressed shifts a color more noticeably toward contrast in Oklab space.
func DerivePressed(base ui.Color, dark bool) ui.Color {
	if dark {
		return base.MixOklab(ui.RGB(255, 255, 255), 0.16)
	}
	return base.MixOklab(ui.RGB(0, 0, 0), 0.15)
}

// DeriveSubtle blends a semantic base color with the canvas background for soft badges/alerts.
func DeriveSubtle(base, bg ui.Color, dark bool) ui.Color {
	if dark {
		return bg.MixOklab(base, 0.18)
	}
	return bg.MixOklab(base, 0.11)
}

// DeriveBorder creates a subtle semantic border from a semantic base and canvas background.
func DeriveBorder(base, bg ui.Color, dark bool) ui.Color {
	if dark {
		return bg.MixOklab(base, 0.40)
	}
	return bg.MixOklab(base, 0.30)
}

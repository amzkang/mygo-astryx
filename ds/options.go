package ds

import (
	"embed"

	"github.com/egoist/mygo/ui"
	"github.com/amzkang/mygo-astryx/theme"
)

//go:embed *.go
var SourceFS embed.FS

// Variant controls visual emphasis of buttons, badges, and tags.
type Variant uint8

const (
	VariantSolid   Variant = iota // High-contrast filled surface
	VariantSoft                   // Tinted subtle fill
	VariantOutline                // Bordered transparent surface
	VariantGhost                  // Borderless until hovered
	VariantPrimary                // Alias for VariantSolid with ToneAccent when Tone is unset
)

// Size controls component height, padding, and font scale.
type Size uint8

const (
	SizeMD Size = iota // Default medium control size
	SizeXS
	SizeSM
	SizeLG
)

// Tone controls semantic intent coloring.
type Tone uint8

const (
	ToneNeutral Tone = iota
	ToneAccent
	ToneInfo
	ToneSuccess
	ToneWarning
	ToneDanger
)

func resolveTone(th *theme.Theme, tone Tone) theme.SemanticTone {
	switch tone {
	case ToneAccent:
		return th.Colors.Accent
	case ToneInfo:
		return th.Colors.Info
	case ToneSuccess:
		return th.Colors.Success
	case ToneWarning:
		return th.Colors.Warning
	case ToneDanger:
		return th.Colors.Danger
	default:
		return theme.SemanticTone{
			Base:    th.Colors.Surface,
			Hover:   th.Colors.SurfaceHover,
			Pressed: th.Colors.SurfacePressed,
			Subtle:  th.Colors.SurfaceSubtle,
			Border:  th.Colors.Border,
			Text:    th.Colors.Text,
			OnColor: th.Colors.Text,
		}
	}
}

func sizeMetrics(th *theme.Theme, sz Size) (height, padX, padY, fontSize, iconSize float32) {
	u := th.Space(1)
	base := th.FontSize
	switch sz {
	case SizeXS:
		return u * 6.5, u * 2, u * 0.75, max(11, base-2), 13
	case SizeSM:
		return u * 7.5, u * 2.5, u * 1, max(12, base-1), 14
	case SizeLG:
		return u * 10.5, u * 4.5, u * 2, base + 1.5, 18
	default: // SizeMD
		return u * 9, u * 3.5, u * 1.5, base, 16
	}
}

func applyInteractiveFill(el ui.Element, base, hover, pressed ui.Color) {
	if el.IsDisabled() {
		el.Background(base)
		return
	}
	if el.Pressed() {
		el.Background(pressed)
	} else if el.Hovered() {
		el.Background(hover)
	} else {
		el.Background(base)
	}
}

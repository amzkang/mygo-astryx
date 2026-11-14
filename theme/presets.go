package theme

import "github.com/egoist/mygo/ui"

// Seven stock themes inspired by Meta Astryx's public theme packages:
// neutral, matcha, stone, butter, chocolate, gothic, y2k.
var (
	// Neutral is the crisp default grayscale + blue accent theme.
	Neutral = DefineTheme(Spec{
		Name:            "neutral",
		Background:      Pair(ui.Hex("#ffffff"), ui.Hex("#141417")),
		Surface:         Pair(ui.Hex("#f4f4f6"), ui.Hex("#1e1e24")),
		SurfaceElevated: Pair(ui.Hex("#ffffff"), ui.Hex("#26262e")),
		Border:          Pair(ui.Hex("#e0e0e6"), ui.Hex("#32323c")),
		Text:            Pair(ui.Hex("#18181b"), ui.Hex("#f4f4f6")),
		TextMuted:       Pair(ui.Hex("#71717a"), ui.Hex("#a1a1aa")),
		Accent:          Pair(ui.Hex("#2563eb"), ui.Hex("#3b82f6")),
		Radius:          6,
	})

	// Matcha pairs warm paper neutrals with botanical matcha green accents.
	Matcha = DefineTheme(Spec{
		Name:            "matcha",
		Background:      Pair(ui.Hex("#fbfcf8"), ui.Hex("#121714")),
		Surface:         Pair(ui.Hex("#f1f5ea"), ui.Hex("#1a221d")),
		SurfaceElevated: Pair(ui.Hex("#ffffff"), ui.Hex("#212b25")),
		Border:          Pair(ui.Hex("#dce4d2"), ui.Hex("#2d3b32")),
		Text:            Pair(ui.Hex("#1a231c"), ui.Hex("#eef3eb")),
		TextMuted:       Pair(ui.Hex("#5d6b5f"), ui.Hex("#99a89b")),
		Accent:          Pair(ui.Hex("#4d7c0f"), ui.Hex("#84cc16")),
		Success:         Pair(ui.Hex("#15803d"), ui.Hex("#4ade80")),
		Radius:          8,
	})

	// Stone delivers a calm architectural mineral palette.
	Stone = DefineTheme(Spec{
		Name:            "stone",
		Background:      Pair(ui.Hex("#fafaf9"), ui.Hex("#161514")),
		Surface:         Pair(ui.Hex("#f1f0ed"), ui.Hex("#21201e")),
		SurfaceElevated: Pair(ui.Hex("#ffffff"), ui.Hex("#292724")),
		Border:          Pair(ui.Hex("#e2e0da"), ui.Hex("#383530")),
		Text:            Pair(ui.Hex("#1c1917"), ui.Hex("#f5f5f4")),
		TextMuted:       Pair(ui.Hex("#78716c"), ui.Hex("#a8a29e")),
		Accent:          Pair(ui.Hex("#44403c"), ui.Hex("#d6d3d1")),
		Radius:          4,
	})

	// Butter brings warm cream surfaces and honey-gold highlights.
	Butter = DefineTheme(Spec{
		Name:            "butter",
		Background:      Pair(ui.Hex("#fffdf5"), ui.Hex("#1a1710")),
		Surface:         Pair(ui.Hex("#fbf5de"), ui.Hex("#252016")),
		SurfaceElevated: Pair(ui.Hex("#ffffff"), ui.Hex("#2f291c")),
		Border:          Pair(ui.Hex("#eadfbc"), ui.Hex("#413825")),
		Text:            Pair(ui.Hex("#272011"), ui.Hex("#fbf7eb")),
		TextMuted:       Pair(ui.Hex("#7c6f52"), ui.Hex("#b7a98a")),
		Accent:          Pair(ui.Hex("#b45309"), ui.Hex("#f59e0b")),
		Radius:          10,
	})

	// Chocolate provides rich espresso surfaces and warm terracotta accents.
	Chocolate = DefineTheme(Spec{
		Name:            "chocolate",
		Background:      Pair(ui.Hex("#fdf8f5"), ui.Hex("#181210")),
		Surface:         Pair(ui.Hex("#f5ebe5"), ui.Hex("#241b18")),
		SurfaceElevated: Pair(ui.Hex("#ffffff"), ui.Hex("#2f231f")),
		Border:          Pair(ui.Hex("#e4d3c9"), ui.Hex("#43332d")),
		Text:            Pair(ui.Hex("#291b16"), ui.Hex("#f7efe9")),
		TextMuted:       Pair(ui.Hex("#7a6258"), ui.Hex("#b49a8e")),
		Accent:          Pair(ui.Hex("#9a3412"), ui.Hex("#ea580c")),
		Radius:          6,
	})

	// Gothic emphasizes stark high-contrast monochrome with crimson accent.
	Gothic = DefineTheme(Spec{
		Name:            "gothic",
		Background:      Pair(ui.Hex("#f8f7fa"), ui.Hex("#0c0b0f")),
		Surface:         Pair(ui.Hex("#eceaf1"), ui.Hex("#16141c")),
		SurfaceElevated: Pair(ui.Hex("#ffffff"), ui.Hex("#1f1c28")),
		Border:          Pair(ui.Hex("#d3cfe0"), ui.Hex("#2f2a3d")),
		Text:            Pair(ui.Hex("#121017"), ui.Hex("#f2f0f7")),
		TextMuted:       Pair(ui.Hex("#635d73"), ui.Hex("#9b93ad")),
		Accent:          Pair(ui.Hex("#9f1239"), ui.Hex("#f43f5e")),
		Radius:          2,
	})

	// Y2K brings retro cyber-silver surfaces and electric cyan-violet accents.
	Y2K = DefineTheme(Spec{
		Name:            "y2k",
		Background:      Pair(ui.Hex("#f3f6fb"), ui.Hex("#0d111a")),
		Surface:         Pair(ui.Hex("#e5ebf5"), ui.Hex("#161d2b")),
		SurfaceElevated: Pair(ui.Hex("#ffffff"), ui.Hex("#1e283b")),
		Border:          Pair(ui.Hex("#c5d2e8"), ui.Hex("#2d3d59")),
		Text:            Pair(ui.Hex("#0f172a"), ui.Hex("#eef4ff")),
		TextMuted:       Pair(ui.Hex("#566680"), ui.Hex("#91a3c2")),
		Accent:          Pair(ui.Hex("#6d28d9"), ui.Hex("#22d3ee")),
		Radius:          12,
	})
)

// StockThemes returns all built-in ThemeSets in display order.
func StockThemes() []*ThemeSet {
	return []*ThemeSet{Neutral, Matcha, Stone, Butter, Chocolate, Gothic, Y2K}
}

// ByName finds a built-in ThemeSet by case-exact name, falling back to Neutral.
func ByName(name string) *ThemeSet {
	for _, ts := range StockThemes() {
		if ts.Spec.Name == name {
			return ts
		}
	}
	return Neutral
}

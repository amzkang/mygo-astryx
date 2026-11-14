package ds

import (
	"fmt"
	"strings"

	"github.com/egoist/mygo/ui"
	"github.com/amzkang/mygo-astryx/icons"
	"github.com/amzkang/mygo-astryx/theme"
)

// Pagination renders accessible page navigation controls for 1-based *page in [1, totalPages].
func Pagination(c *ui.Context, page *int, totalPages int) ui.Element {
	th := theme.Of(c)
	if totalPages < 1 {
		totalPages = 1
	}
	if *page < 1 {
		*page = 1
	}
	if *page > totalPages {
		*page = totalPages
	}

	return ui.Row(c).
		AlignItems(ui.Center).
		Gap(th.Space(1)).
		Role(ui.RoleNavigation).
		Label("Pagination").
		Children(func() {
			prev := IconButton(c.Key("prev"), "Previous page", icons.ChevronLeft, IconButtonOpts{
				Variant:  VariantOutline,
				Size:     SizeSM,
				Disabled: *page <= 1,
			})
			if prev.Clicked() && *page > 1 {
				*page--
			}

			for idx, p := range pageWindow(*page, totalPages) {
				if p == -1 {
					ui.Text(c.Key(fmt.Sprintf("gap-%d", idx)), "…").
						Padding(0, th.Space(1.5)).
						TextColor(th.Colors.TextMuted)
					continue
				}
				active := p == *page
				v := VariantGhost
				tone := ToneNeutral
				if active {
					v = VariantSolid
					tone = ToneAccent
				}
				btn := Button(c.Key(fmt.Sprintf("page-%d", p)), fmt.Sprint(p), ButtonOpts{
					Variant: v,
					Tone:    tone,
					Size:    SizeSM,
				})
				if active {
					btn.Current(ui.CurrentPage)
				}
				if btn.Clicked() {
					*page = p
				}
			}

			next := IconButton(c.Key("next"), "Next page", icons.ChevronRight, IconButtonOpts{
				Variant:  VariantOutline,
				Size:     SizeSM,
				Disabled: *page >= totalPages,
			})
			if next.Clicked() && *page < totalPages {
				*page++
			}
		})
}

func pageWindow(cur, total int) []int {
	if total <= 7 {
		out := make([]int, total)
		for i := range total {
			out[i] = i + 1
		}
		return out
	}
	out := []int{1}
	start := max(2, cur-1)
	end := min(total-1, cur+1)
	if start > 2 {
		out = append(out, -1)
	}
	for p := start; p <= end; p++ {
		out = append(out, p)
	}
	if end < total-1 {
		out = append(out, -1)
	}
	out = append(out, total)
	return out
}

// CommandItem describes an action inside CommandPalette.
type CommandItem struct {
	ID       string
	Title    string
	Subtitle string
	Category string
	Shortcut string
	Icon     *ui.SVG
	OnRun    func()
}

// CommandPaletteState keeps the search query and list scroll/selection state across frames.
type CommandPaletteState struct {
	Query    string
	Selected int
	List     ui.ListState
}

// CommandPalette renders a Cmd+K command launcher dialog using ui.DialogBase and virtualized ui.List.
func CommandPalette(c *ui.Context, open *bool, state *CommandPaletteState, items []CommandItem) ui.Element {
	th := theme.Of(c)
	return ui.DialogBase(c, open, func(backdrop, panel ui.Element) {
		backdrop.Background(th.Colors.Backdrop)
		panel.
			Width(540).
			MaxWidthPercent(92).
			Radius(th.Radius.XL).
			Background(th.Colors.SurfaceElevated).
			Border(1, th.Colors.Border).
			Clip().
			Label("Command Palette")
		th.Shadows.LG.Apply(panel)

		q := strings.ToLower(strings.TrimSpace(state.Query))
		var filtered []CommandItem
		for _, it := range items {
			if q == "" ||
				strings.Contains(strings.ToLower(it.Title), q) ||
				strings.Contains(strings.ToLower(it.Subtitle), q) ||
				strings.Contains(strings.ToLower(it.Category), q) {
				filtered = append(filtered, it)
			}
		}
		if len(filtered) == 0 {
			state.Selected = -1
		} else if state.Selected < 0 || state.Selected >= len(filtered) {
			state.Selected = 0
		}

		ui.Row(c).
			FillWidth().
			AlignItems(ui.Center).
			Gap(th.Space(2.5)).
			Padding(th.Space(3), th.Space(4)).
			BorderWidth(0, 0, 1, 0).
			BorderColor(th.Colors.BorderSubtle).
			Children(func() {
				ui.Icon(c, icons.Search).Size(17, 17).TextColor(th.Colors.TextMuted)
				inp := ui.TextInputBase(c, &state.Query).
					Placeholder("Type a command or search…").
					AutoFocus().
					Grow(1).
					TextColor(th.Colors.Text).
					FontSize(th.Typography.Body.Size).
					Label("Command search")
				if inp.Changed() {
					state.Selected = 0
				}
				if inp.Submitted() && len(filtered) > 0 {
					idx := max(0, min(state.Selected, len(filtered)-1))
					chosen := filtered[idx]
					*open = false
					state.Query = ""
					if chosen.OnRun != nil {
						chosen.OnRun()
					}
				}
				Kbd(c, "Esc")
			})

		if len(filtered) == 0 {
			ui.Column(c).FillWidth().Padding(th.Space(8)).Center().Gap(th.Space(2)).Children(func() {
				ui.Icon(c, icons.Search).Size(24, 24).TextColor(th.Colors.TextMuted)
				Caption(c, "No matching commands")
			})
			return
		}

		state.List.Selected = &state.Selected
		state.List.Key = func(i int) any { return filtered[i].ID }
		listEl := ui.List(c, &state.List, len(filtered), func(i int) {
			it := filtered[i]
			selected := i == state.Selected
			row := ui.Row(c).
				FillWidth().
				AlignItems(ui.Center).
				Gap(th.Space(2.5)).
				Padding(th.Space(2), th.Space(3.5)).
				Radius(th.Radius.SM).
				Cursor(ui.CursorPointer)
			if selected {
				row.Background(th.Colors.Accent.Subtle)
			}
			row.Children(func() {
				icon := it.Icon
				if icon == nil {
					icon = icons.Command
				}
				fg := th.Colors.TextMuted
				if selected {
					fg = th.Colors.Accent.Text
				}
				ui.Icon(c, icon).Size(16, 16).TextColor(fg)
				ui.Column(c).Grow(1).MinWidth(0).Gap(1).Children(func() {
					ui.Text(c, it.Title).FontSize(th.Typography.Body.Size).FontWeight(500).TextColor(th.Colors.Text).SingleLine()
					if it.Subtitle != "" {
						ui.Text(c, it.Subtitle).FontSize(th.Typography.Caption.Size).TextColor(th.Colors.TextMuted).SingleLine()
					}
				})
				if it.Category != "" {
					Badge(c, it.Category, BadgeOpts{Variant: VariantSoft, Tone: ToneNeutral})
				}
				if it.Shortcut != "" {
					Kbd(c, it.Shortcut)
				}
			})
		}).Height(min(float32(len(filtered))*52+12, 320)).Padding(th.Space(1.5))

		if listEl.Submitted() && state.Selected >= 0 && state.Selected < len(filtered) {
			chosen := filtered[state.Selected]
			*open = false
			state.Query = ""
			if chosen.OnRun != nil {
				chosen.OnRun()
			}
		}
	})
}

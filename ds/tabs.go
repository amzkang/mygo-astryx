package ds

import (
	"fmt"

	"github.com/egoist/mygo/ui"
	"github.com/amzkang/mygo-astryx/icons"
	"github.com/amzkang/mygo-astryx/theme"
)

// TabsOpts configures Tabs appearance.
type TabsOpts struct {
	Pills bool
}

// Tabs renders accessible horizontal tabs built on ui.TabsBase with theme-styled indicator/pills.
func Tabs(c *ui.Context, selected *int, labels []string, opts ...TabsOpts) ui.Element {
	th := theme.Of(c)
	var o TabsOpts
	if len(opts) > 0 {
		o = opts[0]
	}
	parts := ui.TabsBase(c, selected, len(labels))
	list := parts.List.Gap(th.Space(1))
	if o.Pills {
		list.Padding(th.Space(0.75)).Radius(th.Radius.MD).Background(th.Colors.SurfaceSubtle).Border(1, th.Colors.BorderSubtle)
	} else {
		list.BorderWidth(0, 0, 1, 0).BorderColor(th.Colors.BorderSubtle)
	}
	list.Children(func() {
		for i, lbl := range labels {
			on := i == *selected
			tab := parts.Tab(i).
				Center().
				Gap(th.Space(1.5)).
				Padding(th.Space(1.75), th.Space(3)).
				Cursor(ui.CursorPointer)
			if o.Pills {
				tab.Radius(max(0, th.Radius.MD-2))
				if on {
					tab.Background(th.Colors.SurfaceElevated).Border(1, th.Colors.Border)
					th.Shadows.SM.Apply(tab)
				} else if tab.Hovered() {
					tab.Background(th.Colors.SurfaceHover)
				}
			} else {
				tab.Radius(th.Radius.SM)
				if on {
					tab.BorderWidth(0, 0, 2, 0).BorderColor(th.Colors.Accent.Base)
				} else if tab.Hovered() {
					tab.Background(th.Colors.SurfaceHover)
				}
			}
			fg := th.Colors.TextMuted
			weight := 500
			if on {
				fg = th.Colors.Text
				weight = 600
			}
			tab.Children(func() {
				ui.Text(c, lbl).
					FontSize(th.Typography.BodySM.Size).
					FontWeight(weight).
					TextColor(fg).
					SingleLine()
			})
		}
	})
	return list
}

// Breadcrumbs renders a semantic navigation trail and writes the clicked index to *chosen (if non-nil).
func Breadcrumbs(c *ui.Context, items []string, chosen *int) ui.Element {
	th := theme.Of(c)
	return ui.Row(c).
		AlignItems(ui.Center).
		Gap(th.Space(1)).
		Role(ui.RoleNavigation).
		Label("Breadcrumbs").
		Children(func() {
			for i, item := range items {
				isLast := i == len(items)-1
				if i > 0 {
					ui.Icon(c.Key(fmt.Sprintf("sep-%d", i)), icons.ChevronRight).
						Size(13, 13).
						TextColor(th.Colors.TextMuted)
				}
				if isLast {
					ui.Text(c.Key(fmt.Sprintf("crumb-%d", i)), item).
						FontSize(th.Typography.BodySM.Size).
						FontWeight(600).
						TextColor(th.Colors.Text).
						Current(ui.CurrentPage).
						SingleLine()
					continue
				}
				btn := ui.ButtonBase(c.Key(fmt.Sprintf("crumb-%d", i))).
					Padding(th.Space(0.5), th.Space(1.5)).
					Radius(th.Radius.SM).
					Cursor(ui.CursorPointer)
				if btn.Hovered() {
					btn.Background(th.Colors.SurfaceHover)
				}
				btn.Children(func() {
					ui.Text(c, item).
						FontSize(th.Typography.BodySM.Size).
						TextColor(th.Colors.TextMuted).
						SingleLine()
				})
				if btn.Clicked() && chosen != nil {
					*chosen = i
				}
			}
		})
}

package ds

import (
	"fmt"

	"github.com/egoist/mygo/ui"
	"github.com/amzkang/mygo-astryx/icons"
	"github.com/amzkang/mygo-astryx/theme"
)

// DataTable renders a themed sortable/scrollable table backed by ui.Table, showing EmptyState when rowCount == 0.
func DataTable(c *ui.Context, state *ui.ListState, columns []ui.TableColumn, rowCount int, emptyTitle string, cell func(row, col int)) ui.Element {
	th := theme.Of(c)
	if rowCount == 0 {
		if emptyTitle == "" {
			emptyTitle = "No rows to display"
		}
		return EmptyState(c, emptyTitle, "", EmptyStateOpts{Icon: icons.Inbox})
	}
	return ui.Box(c).
		FillWidth().
		Radius(th.Radius.LG).
		Background(th.Colors.Surface).
		Border(1, th.Colors.Border).
		Clip().
		Children(func() {
			ui.Table(c, state, columns, rowCount, cell).FillWidth()
		})
}

// TreeNode describes a hierarchical item rendered by TreeView.
type TreeNode struct {
	ID       string
	Label    string
	Icon     *ui.SVG
	Open     *bool
	Children []TreeNode
}

// TreeView renders a hierarchical expandable tree built on ui.TreeBase.
func TreeView(c *ui.Context, roots []TreeNode, selected *string) ui.Element {
	th := theme.Of(c)
	tree := ui.TreeBase(c)
	tree.Root.
		FillWidth().
		Padding(th.Space(1)).
		Gap(th.Space(0.5)).
		Radius(th.Radius.LG).
		Background(th.Colors.Surface).
		Border(1, th.Colors.Border)

	var renderNode func(n TreeNode)
	renderNode = func(n TreeNode) {
		item := tree.Item(n.Open)
		isSelected := selected != nil && *selected == n.ID
		indent := float32(item.Depth()) * th.Space(4)
		row := item.Row.
			FillWidth().
			AlignItems(ui.Center).
			Gap(th.Space(1.5)).
			Padding(th.Space(1.5), th.Space(2.5), th.Space(1.5), th.Space(2.5)+indent).
			Radius(th.Radius.SM).
			Cursor(ui.CursorPointer)
		if isSelected {
			row.Background(th.Colors.Accent.Subtle)
		} else if row.Hovered() {
			row.Background(th.Colors.SurfaceHover)
		}

		fg := th.Colors.Text
		if isSelected {
			fg = th.Colors.Accent.Text
		}
		row.Children(func() {
			if n.Open != nil && len(n.Children) > 0 {
				chev := icons.ChevronRight
				if *n.Open {
					chev = icons.ChevronDown
				}
				toggle := ui.ButtonBase(c.Key("toggle-"+n.ID)).
					Label("Toggle " + n.Label).
					Size(18, 18).
					Center().
					Radius(th.Radius.XS)
				toggle.Children(func() {
					ui.Icon(c, chev).Size(13, 13).TextColor(th.Colors.TextMuted)
				})
				if toggle.Clicked() {
					*n.Open = !*n.Open
				}
			} else {
				ui.Box(c).Size(18, 18)
			}
			if n.Icon != nil {
				ui.Icon(c, n.Icon).Size(15, 15).TextColor(fg)
			}
			ui.Text(c, n.Label).
				FontSize(th.Typography.BodySM.Size).
				FontWeight(500).
				TextColor(fg).
				SingleLine()
		})
		if row.Clicked() && selected != nil {
			*selected = n.ID
		}
		if len(n.Children) > 0 {
			item.Group(func() {
				for _, child := range n.Children {
					renderNode(child)
				}
			})
		}
	}

	tree.Root.Children(func() {
		for _, r := range roots {
			renderNode(r)
		}
	})
	return tree.Root
}

// Spinner renders a size-scaled activity indicator with accessible Label.
func Spinner(c *ui.Context, label string, size Size) ui.Element {
	th := theme.Of(c)
	_, _, _, _, iconSize := sizeMetrics(th, size)
	if label == "" {
		label = "Loading"
	}
	return ui.Spinner(c).
		Size(iconSize, iconSize).
		TextColor(th.Colors.Accent.Base).
		Label(label)
}

// Meter renders a labeled level indicator with threshold coloring built on ui.Meter.
func Meter(c *ui.Context, label string, value, lo, hi float64, levels *ui.MeterLevels) ui.Element {
	th := theme.Of(c)
	return ui.Column(c).FillWidth().Gap(th.Space(1.5)).Children(func() {
		if label != "" {
			ui.Row(c).FillWidth().Justify(ui.SpaceBetween).AlignItems(ui.Center).Children(func() {
				ui.Text(c, label).FontSize(th.Typography.BodySM.Size).FontWeight(500).TextColor(th.Colors.Text)
				ui.Text(c, fmt.Sprintf("%.0f / %.0f", value, hi)).Font(th.MonoFont).FontSize(th.Typography.Caption.Size).TextColor(th.Colors.TextMuted)
			})
		}
		ui.Meter(c, value, lo, hi, levels).FillWidth().Label(label)
	})
}

// Collapsible renders a single expandable disclosure card built on ui.CollapsibleBase.
func Collapsible(c *ui.Context, title string, open *bool, fn func()) ui.Element {
	th := theme.Of(c)
	parts := ui.CollapsibleBase(c, open)
	return ui.Column(c).
		FillWidth().
		Radius(th.Radius.LG).
		Background(th.Colors.Surface).
		Border(1, th.Colors.Border).
		Clip().
		Children(func() {
			tr := parts.Trigger.
				FillWidth().
				AlignItems(ui.Center).
				Justify(ui.SpaceBetween).
				Padding(th.Space(3), th.Space(4)).
				Cursor(ui.CursorPointer)
			if tr.Hovered() {
				tr.Background(th.Colors.SurfaceHover)
			}
			tr.Children(func() {
				ui.Text(c, title).FontSize(th.Typography.Body.Size).FontWeight(600).TextColor(th.Colors.Text)
				ui.Icon(c, icons.ChevronRight).
					Size(16, 16).
					Rotate(90 * parts.Progress()).
					TextColor(th.Colors.TextMuted)
			})
			if panel := parts.Panel(func() {
				if fn != nil {
					fn()
				}
			}); panel != (ui.Element{}) {
				panel.FillWidth().Padding(0, th.Space(4), th.Space(3.5), th.Space(4)).Gap(th.Space(2))
			}
		})
}

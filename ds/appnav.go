package ds

import (
	"unicode/utf8"

	"github.com/amzkang/mygo-astryx/icons"
	"github.com/amzkang/mygo-astryx/theme"
	"github.com/egoist/mygo/ui"
)

// TopNavSlots builds the title and optional start, center and end controls in-frame.
type TopNavSlots struct {
	Heading, Start, Center, End func()
}

// TopNavOpts names and sizes the application header. Zero height sizes to content.
type TopNavOpts struct {
	Label  string
	Height float32
}

// TopNav composes a themed navigation header with flexible central content.
func TopNav(c *ui.Context, slots TopNavSlots, opts ...TopNavOpts) ui.Element {
	th := theme.Of(c)
	var o TopNavOpts
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Label == "" {
		o.Label = "Top navigation"
	}
	e := ui.Row(c).FillWidth().MinWidth(0).Shrink(0).AlignItems(ui.Center).
		Gap(th.Space(3)).Padding(th.Space(3)).Background(th.Colors.Surface).
		BorderWidth(0, 0, 1, 0).BorderColor(th.Colors.BorderSubtle).
		Role(ui.RoleNavigation).Label(o.Label)
	if o.Height > 0 {
		e.Height(o.Height)
	}
	return e.Children(func() {
		if slots.Start != nil {
			ui.Row(c.Key("start")).AlignItems(ui.Center).Gap(th.Space(2)).Children(slots.Start)
		}
		if slots.Heading != nil {
			ui.Column(c.Key("heading")).MinWidth(0).Children(slots.Heading)
		}
		ui.Row(c.Key("center")).Grow(1).MinWidth(0).AlignItems(ui.Center).Justify(ui.Center).Children(slots.Center)
		if slots.End != nil {
			ui.Row(c.Key("end")).AlignItems(ui.Center).Gap(th.Space(2)).Children(slots.End)
		}
	})
}

// NavItem describes a selectable destination or a collapsible section.
// IDs must be nonempty and unique across the complete navigation tree.
type NavItem struct {
	ID, Label string
	Icon      *ui.SVG
	Disabled  bool
	Children  []NavItem
}

// NavState holds only durable navigation data. Missing OpenSections keys default
// to open. Collapsed is the caller-controlled icon-rail presentation.
type NavState struct {
	SelectedID   string
	Collapsed    bool
	OpenSections map[string]bool
	activation   uint64
}

// NavOpts names and sizes navigation. Width defaults to 60 theme units and
// CollapsedWidth to 16. Title is the mobile drawer's accessible dialog name.
type NavOpts struct {
	Label, Title          string
	Width, CollapsedWidth float32
}

// SideNav renders a native sidebar with collapsible, keyboard-activatable sections.
// The sidebar's arrows/Home/End choose enabled destinations. Tab reaches section
// buttons, whose Enter/Space toggles their children. Disabled items are skipped.
func SideNav(c *ui.Context, state *NavState, items []NavItem, opts ...NavOpts) ui.Element {
	th := theme.Of(c)
	var o NavOpts
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Label == "" {
		o.Label = "Side navigation"
	}
	if o.Width <= 0 {
		o.Width = th.Space(60)
	}
	if o.CollapsedWidth <= 0 {
		o.CollapsedWidth = th.Space(16)
	}
	if state.OpenSections == nil {
		state.OpenSections = make(map[string]bool)
	}
	width := o.Width
	if state.Collapsed {
		width = o.CollapsedWidth
	}
	return ui.Column(c).Width(width).FillHeight().MinHeight(0).Shrink(0).
		Background(th.Colors.Surface).Role(ui.RoleNavigation).Label(o.Label).Children(func() {
		ui.Sidebar(c.Key("destinations"), &state.SelectedID, func() {
			renderNavItems(c, state, items, 1, false)
		}).Fill().MinWidth(0).Background(th.Colors.Surface)
	})
}

func renderNavItems(c *ui.Context, state *NavState, items []NavItem, level int, disabled bool) {
	th := theme.Of(c)
	for _, item := range items {
		key := c.Key(item.ID)
		off := disabled || item.Disabled
		if len(item.Children) > 0 {
			open, exists := state.OpenSections[item.ID]
			if !exists {
				open = true
			}
			var button ui.Element
			if state.Collapsed {
				icon := item.Icon
				if icon == nil {
					icon = icons.Menu
				}
				button = IconButton(key, item.Label, icon, IconButtonOpts{Variant: VariantGhost, Disabled: off})
			} else {
				arrow := icons.ChevronRight
				if c.Direction() == ui.RightToLeft {
					arrow = icons.ChevronLeft
				}
				if open {
					arrow = icons.ChevronDown
				}
				button = Button(key, item.Label, ButtonOpts{Variant: VariantGhost, LeadingIcon: item.Icon, TrailingIcon: arrow, Disabled: off, FullWidth: true}).Justify(ui.Start)
			}
			button.Expanded(open)
			if button.Clicked() {
				open = !open
				state.OpenSections[item.ID] = open
			}
			if open {
				ui.Column(key.Key("children")).FillWidth().PaddingStart(th.Space(2)).Children(func() {
					renderNavItems(c, state, item.Children, level+1, off)
				})
			}
			continue
		}
		label := item.Label
		if state.Collapsed {
			label = ""
		}
		var row ui.Element
		if off {
			// A disabled destination stays in the accessibility tree, but does
			// not enter Sidebar's native keyboard selection list.
			row = ui.Row(key).Disabled(true).Role(ui.RoleTreeItem).Gap(th.Space(2)).Padding(th.Space(1.5), th.Space(2)).Children(func() {
				if item.Icon != nil {
					ui.Icon(c, item.Icon).Size(16, 16).TextColor(th.Colors.TextMuted)
				}
				ui.Text(c, label).SingleLine().TextColor(th.Colors.TextMuted)
			})
		} else {
			row = ui.SidebarItem(key, item.ID, item.Icon, label)
			if row.Clicked() {
				state.activation++
			}
		}
		row.Label(item.Label).Level(level).MinWidth(0)
		if state.SelectedID == item.ID {
			row.Current(ui.CurrentPage)
		}
		if state.Collapsed {
			row.Tooltip(item.Label)
			if item.Icon == nil {
				r, _ := utf8.DecodeRuneInString(item.Label)
				row.Children(func() { ui.Text(c, string(r)).TextColor(th.Colors.Text) })
			}
		}
	}
}

// MobileNav hosts frame-local content in a modal drawer at the writing-direction
// start edge. Drawer owns Escape/outside dismissal and modal focus restoration.
func MobileNav(c *ui.Context, open *bool, fn func(), opts ...NavOpts) ui.Element {
	th := theme.Of(c)
	var o NavOpts
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.Title == "" {
		o.Title = "Navigation"
	}
	if o.Width <= 0 {
		o.Width = th.Space(80)
	}
	side := DrawerLeft
	if c.Direction() == ui.RightToLeft {
		side = DrawerRight
	}
	return Drawer(c, open, fn, DrawerOpts{Side: side, Width: o.Width, Title: o.Title})
}

// AppShellState persists navigation and narrow-window drawer state.
type AppShellState struct {
	Nav           NavState
	MobileNavOpen bool
}

// AppShellSlots builds the application regions in the current frame.
type AppShellSlots struct {
	TopNav, SideNav, Banner, Content func()
}

// AppShellOpts controls the DIP breakpoint and region names. Zero dimensions use
// 60 theme units of navigation, 640 DIPs of content and six units of gap.
type AppShellOpts struct {
	NavWidth, ContentMinWidth, Gap float32
	NavigationLabel, ContentLabel  string
}

// AppShell builds a window-root layout, switching its sidebar to a modal drawer
// below NavWidth+ContentMinWidth+Gap. F6 moves between navigation and main content.
func AppShell(c *ui.Context, state *AppShellState, slots AppShellSlots, opts ...AppShellOpts) ui.Element {
	th := theme.Of(c)
	var o AppShellOpts
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.NavWidth <= 0 {
		o.NavWidth = th.Space(60)
	}
	if o.ContentMinWidth <= 0 {
		o.ContentMinWidth = 640
	}
	if o.Gap <= 0 {
		o.Gap = th.Space(6)
	}
	if o.NavigationLabel == "" {
		o.NavigationLabel = "Application navigation"
	}
	if o.ContentLabel == "" {
		o.ContentLabel = "Main content"
	}
	width, _ := c.Size()
	narrow := slots.SideNav != nil && width < o.NavWidth+o.ContentMinWidth+o.Gap
	if !narrow {
		state.MobileNavOpen = false
	}
	root := ui.Column(c).Fill().MinWidth(0).MinHeight(0).Background(th.Colors.Background)
	focus := ui.Local(root, "region_focus", func() string { return "" })
	if !state.MobileNavOpen && c.Shortcut(0, ui.KeyF6) {
		*focus = "content"
		if ui.FocusedValue(c, focus) == "content" && slots.SideNav != nil {
			*focus = "navigation"
		}
	}
	root.Children(func() {
		ui.Row(c.Key("header")).FillWidth().AlignItems(ui.Center).Children(func() {
			if narrow {
				trigger := IconButton(c.Key("open_navigation"), "Open navigation", icons.Menu, IconButtonOpts{Variant: VariantGhost})
				trigger.FocusBind(focus, "navigation").Expanded(state.MobileNavOpen)
				if trigger.Clicked() {
					state.MobileNavOpen = true
				}
			}
			if slots.TopNav != nil {
				ui.Column(c.Key("topnav")).Grow(1).MinWidth(0).Children(slots.TopNav)
			}
		})
		if slots.Banner != nil {
			ui.Column(c.Key("banner")).FillWidth().Shrink(0).Children(slots.Banner)
		}
		ui.Row(c.Key("body")).Grow(1).MinHeight(0).AlignItems(ui.Stretch).Gap(o.Gap).Children(func() {
			if !narrow && slots.SideNav != nil {
				navWidth := o.NavWidth
				if state.Nav.Collapsed {
					navWidth = th.Space(16)
				}
				ui.Column(c.Key("navigation")).Width(navWidth).Shrink(0).MinHeight(0).
					Role(ui.RoleNavigation).Label(o.NavigationLabel).TabStop(false).FocusBind(focus, "navigation").Children(slots.SideNav)
			}
			ui.Column(c.Key("content")).Grow(1).MinWidth(0).MinHeight(0).
				Role(ui.RoleRegion).Label(o.ContentLabel).TabStop(false).FocusBind(focus, "content").Children(slots.Content)
		})
	})
	if narrow {
		selected := state.Nav.SelectedID
		activation := state.Nav.activation
		MobileNav(c.Key("mobile"), &state.MobileNavOpen, slots.SideNav, NavOpts{Title: o.NavigationLabel, Width: o.NavWidth + th.Space(8)})
		if state.Nav.SelectedID != selected || state.Nav.activation != activation {
			state.MobileNavOpen = false
		}
	}
	return root
}

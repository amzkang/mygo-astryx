package main

import (
	"fmt"

	"github.com/amzkang/mygo-astryx/ds"
	"github.com/amzkang/mygo-astryx/icons"
	"github.com/amzkang/mygo-astryx/theme"
	"github.com/egoist/mygo/ui"
)

var patternNames = []string{
	"app-shell-basic", "app-shell-narrow", "nav-hierarchy", "nav-collapsed",
	"scroll-long-content", "media-grid", "action-overflow", "responsive-media-detail",
}

func patternNavItems() []ds.NavItem {
	return []ds.NavItem{
		{ID: "home", Label: "Workspace home", Icon: icons.Home},
		{ID: "projects", Label: "Projects", Icon: icons.Layers, Children: []ds.NavItem{
			{ID: "design", Label: "Design system", Icon: icons.Palette},
			{ID: "releases", Label: "Releases", Icon: icons.Box},
		}},
		{ID: "settings", Label: "Workspace settings", Icon: icons.Settings},
		{ID: "locked", Label: "Billing unavailable", Icon: icons.AlertCircle, Disabled: true},
	}
}

func renderPatternsPage(c *ui.Context, app *galleryState, th *theme.Theme) {
	ds.Section(c, "Application patterns", "Layout and navigation compositions", func() {
		ds.Select(c.Key("pattern"), &app.pattern, patternNames, nil).Label("Pattern").Width(300)
		if app.lastPattern != app.pattern {
			app.patternShell.Nav.Collapsed = app.pattern == "nav-collapsed"
			app.lastPattern = app.pattern
		}
		ui.Box(c).FillWidth().Height(460).Clip().Border(1, th.Colors.Border).Children(func() {
			renderPattern(c, app, th)
		})
	})
}

func renderPattern(c *ui.Context, app *galleryState, th *theme.Theme) {
	content := func() {
		ds.ScrollableArea(c.Key("pattern-content"), func() {
			ds.VStack(c, 3, func() {
				ds.Heading(c, 2, "Workspace overview")
				ds.Body(c, "Selected destination: "+app.patternShell.Nav.SelectedID)
				renderPatternMedia(c, th)
				ds.Link(c, "MyGo documentation", "https://github.com/egoist/mygo")
			})
		}, ds.ScrollOpts{Label: "Workspace activity", State: &app.patternScroll})
	}
	top := func() {
		ds.TopNav(c, ds.TopNavSlots{
			Heading: func() { ds.Heading(c, 3, "Astryx workspace") },
			End: func() {
				if ds.IconButton(c.Key("collapse"), "Toggle side navigation", icons.PanelRight, ds.IconButtonOpts{Variant: ds.VariantGhost}).Clicked() {
					app.patternShell.Nav.Collapsed = !app.patternShell.Nav.Collapsed
				}
				renderPatternMoreMenu(c, app)
			},
		})
	}
	side := func() { ds.SideNav(c.Key("pattern-nav"), &app.patternShell.Nav, patternNavItems()) }
	switch app.pattern {
	case "app-shell-basic", "app-shell-narrow", "nav-collapsed":
		o := ds.AppShellOpts{}
		if app.pattern == "app-shell-narrow" {
			// Force the mobile presentation within this gallery example.
			o.ContentMinWidth = 2000
		}
		ds.AppShell(c, &app.patternShell, ds.AppShellSlots{TopNav: top, SideNav: side, Content: func() {
			ds.Layout(c, ds.LayoutSlots{Content: content})
		}}, o)
	case "nav-hierarchy":
		ds.Layout(c, ds.LayoutSlots{Header: top, Start: side, Content: func() {
			ds.Heading(c, 2, "Navigation hierarchy")
			ds.Link(c, "MyGo documentation", "https://github.com/egoist/mygo")
		}})
	case "scroll-long-content":
		ds.Layout(c, ds.LayoutSlots{Content: func() {
			ds.ScrollableArea(c, func() {
				ds.Grid(c, func() {
					for i := range 40 {
						ds.Card(c.Key(fmt.Sprintf("activity-%d", i)), func() {
							ds.Body(c, fmt.Sprintf("Activity %02d", i+1))
						})
					}
				})
			}, ds.ScrollOpts{Label: "Activity feed", State: &app.patternScroll})
		}})
	case "media-grid":
		ds.ScrollableArea(c, func() { renderPatternMedia(c, th) })
	case "action-overflow":
		ds.Layout(c, ds.LayoutSlots{Header: top, Content: func() {
			ds.Body(c, "Last action: "+app.patternAction)
			ds.Link(c, "MyGo documentation", "https://github.com/egoist/mygo")
		}})
	case "responsive-media-detail":
		ds.Layout(c, ds.LayoutSlots{Header: func() {
			if ds.Button(c, "Show media navigation").Clicked() {
				app.patternMobile = true
			}
		}, Content: func() {
			ds.AspectRatio(c, 16.0/9, func() { ds.Body(c, "Media detail") }).Background(th.Colors.Accent.Subtle)
		}})
		ds.MobileNav(c, &app.patternMobile, func() { ds.Body(c, "Media collection"); side() })
	}
}

func renderPatternMedia(c *ui.Context, th *theme.Theme) {
	ds.Grid(c, func() {
		for _, title := range []string{"Design library", "Release notes", "Workspace settings", "Team activity"} {
			ds.AspectRatio(c.Key(title), 16.0/9, func() {
				ui.Column(c).Fill().Center().Gap(th.Space(2)).Children(func() {
					ui.Icon(c, icons.Layers).Size(28, 28).TextColor(th.Colors.Accent.Text)
					ds.Link(c, title, "https://github.com/egoist/mygo")
				})
			}).Background(th.Colors.Accent.Subtle).Radius(th.Radius.MD)
		}
	}, ds.GridOpts{Columns: 2})
}

func renderPatternMoreMenu(c *ui.Context, app *galleryState) {
	ds.MoreMenu(c.Key("pattern-menu"), &app.pgMenuOpen, []ds.MenuItemSpec{
		{ID: "export", Label: "Export workspace", Icon: icons.Download, OnClick: func() { app.patternAction = "Export workspace" }},
		{ID: "duplicate", Label: "Duplicate workspace", Icon: icons.Copy, OnClick: func() { app.patternAction = "Duplicate workspace" }},
		{ID: "delete", Label: "Delete workspace", Disabled: true},
	})
}

func renderP1aWidget(c *ui.Context, app *galleryState, th *theme.Theme, name string) {
	switch name {
	case "Layout":
		ui.Box(c).Width(460).Height(180).Children(func() {
			ds.Layout(c, ds.LayoutSlots{
				Header:  func() { ds.Heading(c, 3, "Workspace") },
				Start:   func() { ds.Body(c, "Navigation") },
				Content: func() { ds.Body(c, "Main content") },
				Footer:  func() { ds.Caption(c, "Saved") },
			}, ds.LayoutOpts{StartWidth: 120, Gap: th.Space(3)})
		})
	case "TopNav":
		ds.TopNav(c, ds.TopNavSlots{Heading: func() { ds.Heading(c, 3, "Workspace") }, End: func() { renderPatternMoreMenu(c, app) }}).Width(460)
	case "SideNav":
		ui.Box(c).Height(280).Children(func() { ds.SideNav(c, &app.patternShell.Nav, patternNavItems()) })
	case "MobileNav":
		if ds.Button(c, "Open mobile navigation").Clicked() {
			app.patternMobile = true
		}
		ds.MobileNav(c, &app.patternMobile, func() { ds.SideNav(c, &app.patternShell.Nav, patternNavItems()) })
	case "AppShell":
		ui.Box(c).Width(540).Height(260).Children(func() {
			ds.AppShell(c, &app.patternShell, ds.AppShellSlots{
				TopNav:  func() { ds.TopNav(c, ds.TopNavSlots{Heading: func() { ds.Heading(c, 3, "Workspace") }}) },
				SideNav: func() { ds.SideNav(c, &app.patternShell.Nav, patternNavItems(), ds.NavOpts{Width: 180}) },
				Content: func() { ds.Body(c, "Selected: "+app.patternShell.Nav.SelectedID) },
			}, ds.AppShellOpts{NavWidth: 180, ContentMinWidth: 280, Gap: 16})
		})
	case "ScrollableArea":
		ui.Box(c).Width(400).Height(180).Children(func() {
			ds.ScrollableArea(c, func() {
				for i := range 20 {
					ds.Body(c.Key(i), fmt.Sprintf("Activity %02d", i+1))
				}
			}, ds.ScrollOpts{Label: "Activity", State: &app.patternScroll})
		})
	case "MoreMenu":
		renderPatternMoreMenu(c, app)
	case "Grid":
		ds.Grid(c, func() {
			ds.Card(c, func() { ds.Body(c, "First column") })
			ds.Card(c, func() { ds.Body(c, "Second column") })
		}).Width(400)
	case "AspectRatio":
		ds.AspectRatio(c, 16.0/9, func() { ds.Body(c, "16:9 media") }).Width(320).Background(th.Colors.Accent.Subtle)
	case "Link":
		ds.Link(c, "MyGo documentation", "https://github.com/egoist/mygo")
	}
}

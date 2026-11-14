package ds_test

import (
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/amzkang/mygo-astryx/ds"
	"github.com/amzkang/mygo-astryx/icons"
	"github.com/amzkang/mygo-astryx/theme"
	"github.com/egoist/mygo/ui"
)

func navigationItems() []ds.NavItem {
	return []ds.NavItem{
		{ID: "home", Label: "Home", Icon: icons.Home},
		{ID: "disabled", Label: "Unavailable", Disabled: true},
		{ID: "workspace", Label: "Workspace", Children: []ds.NavItem{
			{ID: "projects", Label: "Projects", Icon: icons.Layers},
			{ID: "settings", Label: "Settings", Icon: icons.Settings},
		}},
	}
}

func TestSideNavSelectionSectionsDisabledAndIdentity(t *testing.T) {
	var state ds.NavState
	items := navigationItems()
	tester := ui.NewTester(func(c *ui.Context) {
		theme.WithTheme(c, theme.Matcha.Light, func() {
			ds.SideNav(c.Key("nav"), &state, items)
		}).Fill()
	}, 400, 500)
	if state.SelectedID != "" {
		t.Fatal("empty selection must remain empty")
	}
	if err := tester.Click("Home"); err != nil || state.SelectedID != "home" {
		t.Fatalf("home selection: %v, %+v", err, state)
	}
	if info, ok := tester.Access("Unavailable"); !ok || !info.Disabled {
		t.Fatalf("disabled destination semantics: %+v, %v", info, ok)
	}
	if err := tester.Click("Unavailable"); err == nil && state.SelectedID != "home" {
		t.Fatal("disabled destination changed selection")
	}
	if err := tester.Click("Home"); err != nil {
		t.Fatal(err)
	}
	tester.Key(0, ui.KeyDown)
	if state.SelectedID != "projects" {
		t.Fatalf("Down must skip disabled destination: %q", state.SelectedID)
	}
	tester.Key(0, ui.KeyEnd)
	if state.SelectedID != "settings" {
		t.Fatalf("End must choose last destination: %q", state.SelectedID)
	}
	if info, ok := tester.Access("Settings"); !ok || info.Current != ui.CurrentPage || !info.Selected {
		t.Fatalf("selected destination semantics: %+v, %v", info, ok)
	}
	if err := tester.Click("Workspace"); err != nil || state.OpenSections["workspace"] {
		t.Fatalf("section should close: %v, %+v", err, state)
	}
	tester.Key(0, ui.KeySpace)
	if !state.OpenSections["workspace"] {
		t.Fatal("Space must reopen a focused section")
	}
	// Reorder ordinary data, retaining IDs. Selection and section state survive.
	items[0], items[1] = items[1], items[0]
	state.Collapsed = true
	tester.Frame()
	if state.SelectedID != "settings" || !state.OpenSections["workspace"] {
		t.Fatalf("state lost after reorder/collapse: %+v", state)
	}
	if info, ok := tester.Access("Settings"); !ok || info.Name != "Settings" {
		t.Fatal("icon rail lost its accessible destination name")
	}
	if rect, ok := tester.Find("Side navigation"); !ok || rect.W != theme.Matcha.Light.Space(16) {
		t.Fatalf("collapsed width: %+v, %v", rect, ok)
	}
}

func TestMobileNavDismissalAndFocusRestoration(t *testing.T) {
	for _, rtl := range []bool{false, true} {
		t.Run(fmt.Sprint(rtl), func(t *testing.T) {
			open := false
			tester := ui.NewTester(func(c *ui.Context) {
				theme.WithTheme(c, theme.Neutral.Light, func() {
					root := ui.Column(c).Fill()
					if rtl {
						root.Direction(ui.RightToLeft)
					}
					root.Children(func() {
						if ds.Button(c.Key("trigger"), "Show navigation").Clicked() {
							open = true
						}
						ds.MobileNav(c.Key("mobile"), &open, func() {
							ds.Button(c, "Destination")
						}, ds.NavOpts{Width: 240})
					})
				}).Fill()
			}, 600, 400)
			if err := tester.Click("Show navigation"); err != nil || !open {
				t.Fatalf("open: %v", err)
			}
			tester.Advance(time.Second)
			if !tester.Focused("Close drawer") {
				t.Fatal("drawer must initially focus its close control")
			}
			panel, ok := tester.AccessByRole(ui.RoleDialog, "Navigation")
			if !ok || (!rtl && panel.Bounds.X != 0) || (rtl && panel.Bounds.X <= 0) {
				t.Fatalf("drawer start edge: %+v, %v", panel, ok)
			}
			tester.Key(ui.Shift, ui.KeyTab)
			if !tester.Focused("Destination") {
				t.Fatal("modal must wrap focus inside the drawer")
			}
			tester.Key(0, ui.KeyEscape)
			if open || !tester.Focused("Show navigation") {
				t.Fatalf("Escape must close and restore focus: open=%v", open)
			}
			if err := tester.Click("Show navigation"); err != nil {
				t.Fatal(err)
			}
			x := float32(590)
			if rtl {
				x = 10
			}
			tester.ClickAt(x, 200)
			if open || !tester.Focused("Show navigation") {
				t.Fatal("outside click must close and restore focus")
			}
		})
	}
}

func TestAppShellBreakpointSelectionAndRegionFocus(t *testing.T) {
	for _, density := range []theme.Density{theme.DensityComfortable, theme.DensityCompact, theme.DensitySpacious} {
		t.Run(density.String(), func(t *testing.T) {
			th := theme.Matcha.Light.WithDensity(density)
			threshold := int(th.Space(60) + 640 + th.Space(6))
			state := ds.AppShellState{Nav: ds.NavState{SelectedID: "home", OpenSections: map[string]bool{"workspace": true}}}
			tester := ui.NewTester(func(c *ui.Context) {
				theme.WithTheme(c, th, func() {
					ds.AppShell(c.Key("shell"), &state, ds.AppShellSlots{
						TopNav:  func() { ds.TopNav(c, ds.TopNavSlots{Heading: func() { ds.Heading(c, 2, "App") }}) },
						SideNav: func() { ds.SideNav(c, &state.Nav, navigationItems()) },
						Content: func() { ds.Body(c, "Page content") },
					})
				}).Fill()
			}, threshold, 500)
			if _, ok := tester.Find("Open navigation"); ok {
				t.Fatal("exact breakpoint must use wide mode")
			}
			tester.Key(0, ui.KeyF6)
			if !tester.Focused("Main content") {
				t.Fatal("first F6 must focus main content")
			}
			tester.Key(0, ui.KeyF6)
			if !tester.Focused("Application navigation") {
				t.Fatal("next F6 must focus navigation")
			}
			tester.SetSize(threshold-1, 500)
			if err := tester.Click("Open navigation"); err != nil || !state.MobileNavOpen {
				t.Fatalf("narrow drawer open: %v", err)
			}
			if err := tester.Click("Projects"); err != nil || state.MobileNavOpen || state.Nav.SelectedID != "projects" {
				t.Fatalf("destination must select and close drawer: %v, %+v", err, state)
			}
			if !tester.Focused("Open navigation") {
				t.Fatal("selection dismissal must restore trigger focus")
			}
			if err := tester.Click("Open navigation"); err != nil {
				t.Fatal(err)
			}
			if err := tester.Click("Projects"); err != nil || state.MobileNavOpen {
				t.Fatalf("choosing the current destination must also close: %v", err)
			}
			state.Nav.Collapsed = true
			tester.Frame()
			if err := tester.Click("Open navigation"); err != nil {
				t.Fatal(err)
			}
			tester.SetSize(threshold+1, 500)
			if state.MobileNavOpen || state.Nav.SelectedID != "projects" || !state.Nav.Collapsed || !state.Nav.OpenSections["workspace"] {
				t.Fatalf("widening lost durable state or retained drawer: %+v", state)
			}
			for width := threshold - 3; width <= threshold+3; width++ {
				tester.SetSize(width, 500)
			}
		})
	}
}

func TestP1aLayoutGridRatioAndScroll(t *testing.T) {
	var scroll ui.ScrollState
	tester := ui.NewTester(func(c *ui.Context) {
		theme.WithTheme(c, theme.Neutral.Light, func() {
			ds.Layout(c, ds.LayoutSlots{
				Header: func() { ui.Box(c).FillWidth().Height(40).Label("Header") },
				Start:  func() { ui.Box(c).Fill().Label("Start") },
				End:    func() { ui.Box(c).Fill().Label("End") },
				Footer: func() { ui.Box(c).FillWidth().Height(30).Label("Footer") },
				Content: func() {
					ds.ScrollableArea(c.Key("scroll"), func() {
						ds.Grid(c, func() {
							ds.AspectRatio(c, 2, func() { ds.Link(c, "Docs", "https://github.com/egoist/mygo") }).Label("Media")
							ui.Box(c).Height(40).Label("Second cell")
						})
						ui.Box(c).Height(1200)
					}, ds.ScrollOpts{Label: "Activity", State: &scroll})
				},
			}, ds.LayoutOpts{StartWidth: 100, EndWidth: 80})
		}).Fill()
	}, 800, 500)
	start, _ := tester.Find("Start")
	end, _ := tester.Find("End")
	activity, _ := tester.Find("Activity")
	if start.W != 100 || end.W != 80 || activity.H != 430 || activity.W != 620 {
		h, _ := tester.Find("Header")
		f, _ := tester.Find("Footer")
		t.Fatalf("five-region geometry: start=%+v end=%+v content=%+v header=%+v footer=%+v", start, end, activity, h, f)
	}
	media, _ := tester.Find("Media")
	second, _ := tester.Find("Second cell")
	if math.Abs(float64(media.W/media.H-2)) > .01 || second.X <= media.X {
		t.Fatalf("grid or aspect ratio: media=%+v second=%+v", media, second)
	}
	link, ok := tester.AccessByRole(ui.RoleLink, "Docs")
	if !ok || !link.Focusable {
		t.Fatalf("native link semantics: %+v, %v", link, ok)
	}
	// Click empty space inside the named region to focus its keyboard scroller.
	tester.ClickAt(activity.X+activity.W-5, activity.Y+activity.H-5)
	tester.Key(0, ui.KeyPageDown)
	if scroll.Y <= 0 {
		t.Fatalf("PageDown did not scroll: %+v", scroll)
	}
	tester.Key(0, ui.KeyEnd)
	if scroll.Y != scroll.MaxY || scroll.MaxY == 0 {
		t.Fatalf("End did not reach overflow end: %+v", scroll)
	}
	tester.Key(0, ui.KeyHome)
	if scroll.Y != 0 {
		t.Fatal("Home did not return to start")
	}
	tester.SetDirection(ui.RightToLeft)
	start, _ = tester.Find("Start")
	end, _ = tester.Find("End")
	if start.X <= end.X {
		t.Fatal("layout start/end must follow RTL")
	}
}

func TestMoreMenuKeyboardAndDisabledActions(t *testing.T) {
	open, ran := false, ""
	tester := ui.NewTester(func(c *ui.Context) {
		theme.WithTheme(c, theme.Gothic.Dark, func() {
			ds.MoreMenu(c, &open, []ds.MenuItemSpec{
				{ID: "disabled", Label: "Disabled action", Disabled: true, OnClick: func() { ran = "disabled" }},
				{ID: "export", Label: "Export", OnClick: func() { ran = "export" }},
			})
		}).Fill()
	}, 400, 300)
	if err := tester.Click("More actions"); err != nil || !open {
		t.Fatalf("open overflow menu: %v", err)
	}
	tester.Key(0, ui.KeyDown)
	tester.Key(0, ui.KeyEnter)
	if ran != "export" || open {
		t.Fatalf("keyboard selection must skip disabled and close: ran=%q open=%v", ran, open)
	}
	if !tester.Focused("More actions") {
		t.Fatal("menu must restore trigger focus")
	}
	if err := tester.Click("More actions"); err != nil {
		t.Fatal(err)
	}
	tester.Key(0, ui.KeyEscape)
	if open {
		t.Fatal("Escape must close overflow menu")
	}
}

func TestLayoutAutoHeightEmptySlotsAndInvalidRatios(t *testing.T) {
	for _, ratio := range []float32{0, -2, float32(math.NaN()), float32(math.Inf(1)), 2} {
		tester := ui.NewTester(func(c *ui.Context) {
			ds.Layout(c, ds.LayoutSlots{Content: func() {
				ds.Layout(c, ds.LayoutSlots{Content: func() {
					ds.AspectRatio(c, ratio, nil).Label("Ratio")
				}}, ds.LayoutOpts{AutoHeight: true})
			}}, ds.LayoutOpts{AutoHeight: true})
			ui.Text(c, "After layout")
		}, 200, 500)
		r, ok := tester.Find("Ratio")
		after, _ := tester.Find("After layout")
		wantHeight := float32(200)
		if ratio == 2 {
			wantHeight = 100
		}
		if !ok || r.W != 200 || r.H != wantHeight || after.Y != wantHeight {
			t.Fatalf("ratio=%v, rect=%+v, following=%+v", ratio, r, after)
		}
	}
	calls := 0
	ui.Render(func(c *ui.Context) { ds.Layout(c, ds.LayoutSlots{}); calls++ }, 1, 1, 1)
	if calls == 0 {
		t.Fatal("empty layout did not render")
	}
}

func TestScrollableAreaAxesAndChildKeys(t *testing.T) {
	for _, rtl := range []bool{false, true} {
		var state ui.ScrollState
		text := ""
		tester := ui.NewTester(func(c *ui.Context) {
			root := ui.Column(c).Fill()
			if rtl {
				root.Direction(ui.RightToLeft)
			}
			root.Children(func() {
				ds.ScrollableArea(c, func() {
					ui.Box(c).Size(1000, 800)
				}, ds.ScrollOpts{Both: true, State: &state, Label: "Canvas"})
			})
		}, 400, 300)
		tester.ClickAt(200, 150)
		next := ui.KeyRight
		if rtl {
			next = ui.KeyLeft
		}
		tester.Key(0, next)
		if state.X <= 0 || state.MaxX <= 0 || state.MaxY <= 0 {
			t.Fatalf("logical horizontal scroll rtl=%v: %+v", rtl, state)
		}
		tester.Key(0, ui.KeyEnd)
		if state.X != state.MaxX || state.Y != state.MaxY {
			t.Fatalf("both-axis End: %+v", state)
		}
		var inputState ui.ScrollState
		inputTester := ui.NewTester(func(c *ui.Context) {
			ds.ScrollableArea(c, func() {
				ds.TextInput(c, &text, ds.InputOpts{Placeholder: "Edit inside scroll"}).Label("Input container")
				ui.Box(c).Height(1000)
			}, ds.ScrollOpts{State: &inputState})
		}, 400, 300)
		inputState.X, inputState.Y = 0, 0
		inputTester.Frame()
		input, ok := inputTester.Find("Input container")
		if !ok || input.W <= 0 || input.H <= 0 {
			t.Fatal("scroll content must contain a visible input")
		}
		inputTester.ClickAt(input.X+input.W/2, input.Y+input.H/2)
		inputTester.Type("abc")
		inputTester.Key(0, ui.KeyHome)
		inputTester.Type("z")
		if text != "zabc" || inputState.Y != 0 {
			t.Fatalf("child input lost Home key: text=%q scroll=%+v", text, inputState)
		}
	}
}

func TestTopNavSlotsAndNativeLinkActivation(t *testing.T) {
	url := "https://github.com/egoist/mygo"
	tester := ui.NewTester(func(c *ui.Context) {
		ds.TopNav(c, ds.TopNavSlots{
			Start:   func() { ds.Button(c, "Back") },
			Heading: func() { ds.Heading(c, 2, "Workspace") },
			Center:  func() { ds.Link(c, "Documentation", url) },
			End:     func() { ds.Button(c, "Account") },
		}, ds.TopNavOpts{Label: "Workspace navigation"})
	}, 800, 300)
	if info, ok := tester.AccessByRole(ui.RoleNavigation, "Workspace navigation"); !ok || info.Bounds.W != 800 {
		t.Fatalf("navigation semantics and width: %+v %v", info, ok)
	}
	if err := tester.Click("Documentation"); err != nil {
		t.Fatal(err)
	}
	tester.Key(0, ui.KeyEnter)
	if urls := tester.OpenedURLs(); len(urls) != 2 || urls[0] != url || urls[1] != url {
		t.Fatalf("native pointer/keyboard URL activation: %v", urls)
	}
}

func TestAppShellThemesDensityRTLAndCustomBreakpoint(t *testing.T) {
	for _, stock := range theme.StockThemes() {
		for _, dark := range []bool{false, true} {
			for _, density := range []theme.Density{theme.DensityComfortable, theme.DensityCompact, theme.DensitySpacious} {
				th := stock.For(dark).WithDensity(density)
				for _, rtl := range []bool{false, true} {
					var state ds.AppShellState
					menuOpen := false
					tester := ui.NewTester(func(c *ui.Context) {
						theme.WithTheme(c, th, func() {
							root := ui.Column(c).Fill()
							if rtl {
								root.Direction(ui.RightToLeft)
							}
							root.Children(func() {
								ds.AppShell(c, &state, ds.AppShellSlots{
									TopNav: func() {
										ds.TopNav(c, ds.TopNavSlots{
											Heading: func() { ds.Heading(c, 3, "Workspace with a longer localized title") },
											End:     func() { ds.MoreMenu(c, &menuOpen, nil) },
										})
									},
									SideNav: func() { ds.SideNav(c, &state.Nav, navigationItems(), ds.NavOpts{Width: 180}) },
									Content: func() {
										ds.Layout(c, ds.LayoutSlots{Content: func() {
											ds.ScrollableArea(c, func() {
												ds.Grid(c, func() {
													ds.AspectRatio(c, 2, func() { ds.Link(c, "Docs", "https://github.com/egoist/mygo") })
													ds.Body(c, "Content")
												})
											})
										}})
									},
								}, ds.AppShellOpts{NavWidth: 180, ContentMinWidth: 400, Gap: 20})
							})
						}).Fill()
					}, 600, 350)
					content, ok := tester.AccessByRole(ui.RoleRegion, "Main content")
					if !ok || content.Bounds.W != 400 || content.Bounds.H <= 0 {
						t.Fatalf("theme=%s dark=%v density=%s rtl=%v content=%+v", th.Name, dark, density, rtl, content)
					}
					header, _ := tester.AccessByRole(ui.RoleNavigation, "Top navigation")
					if content.Bounds.Y < header.Bounds.Y+header.Bounds.H {
						t.Fatal("scrolling content must stay below the header")
					}
					tester.SetSize(599, 350)
					if _, ok := tester.Find("Open navigation"); !ok {
						t.Fatal("custom breakpoint ignored")
					}
					tester.SetSize(0, 0)
					tester.SetSize(600, 350)
				}
			}
		}
	}
}

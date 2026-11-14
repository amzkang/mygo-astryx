# mygo-astryx

A pure-Go native UI design system, component library, multi-theme token engine, and agent-ready CLI/MCP toolkit built on [MyGo](https://github.com/egoist/mygo) v0.8.3 (inspired by [Meta Astryx](https://github.com/meta/astryx)).

## Highlights

- **71 Native Components (`ds/`)**: Zero-WebView GPU-rendered controls spanning layout, actions, forms, feedback, navigation, overlays, menus, and data display—built directly on MyGo primitives (`ui.ButtonBase`, `ui.MenuBase`, `ui.DialogBase`, `ui.TreeBase`, `ui.Table`, `ui.ToastViewportBase`, etc.).
- **Multi-Theme Token Engine (`theme/`)**: 7 built-in light + dark theme presets (`Neutral`, `Matcha`, `Stone`, `Butter`, `Chocolate`, `Gothic`, `Y2K`), OKLab color mixing, density scaling (`Compact`, `Comfortable`, `Spacious`), and programmatic WCAG contrast validation (`>=7:1` body text, `>=4.5:1` interactive foregrounds).
- **56 Embedded Icons (`icons/`)**: Crisp 24x24 Lucide-style SVG icons compiled directly into the binary via `//go:embed`.
- **Interactive Native Gallery (`cmd/gallery`)**: Live component catalog with interactive prop playground, real-time Go snippet generator, 7-theme WCAG contrast matrix, RTL/Dark/Density switching, and `Cmd+K` Command Palette.
- **Agent-Ready CLI & MCP Server (`cmd/astryx`)**: Search, dense component docs, static hygiene linter (`astryx check --vet`), template scaffolding (`starter`, `settings-shell`), `.mcp.json` + `AGENTS.md`/`CLAUDE.md` bootstrapping (`astryx init`), swizzle ejection (`astryx swizzle`), and stdio MCP server (`astryx mcp`).

## Project Layout

```
mygo-astryx/
├── theme/        # Design tokens, DefineTheme, 7 stock presets, WCAG contrast, WithTheme
├── icons/        # 56 embedded 24x24 SVG icons
├── ds/           # 71 swizzle-safe exported UI components + embedded SourceFS
├── registry/     # Structured component, token, preset, guide, and template metadata
├── cmd/astryx/   # Agent-ready CLI + MCP stdio server
├── cmd/gallery/  # Native desktop showcase & interactive playground
├── scripts/      # Full verification script (verify.sh)
└── bin/go        # Repo Go toolchain wrapper (Go 1.27.3)
```

## Quick Start

### Run the Native Gallery

```bash
./bin/go run ./cmd/gallery
```

Press `Cmd+K` (or click **Command Palette**) to switch pages, themes, dark mode, or RTL layout.

### Use in Your Own MyGo App

```go
package main

import (
	"log"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
	"github.com/amzkang/mygo-astryx/ds"
	"github.com/amzkang/mygo-astryx/icons"
	"github.com/amzkang/mygo-astryx/theme"
)

func main() {
	dark := false
	name := ""

	mygo.App.WhenReady(func() {
		mygo.NewWindow(mygo.WindowOptions{
			Title:  "My Astryx App",
			Width:  720,
			Height: 480,
			Content: ui.View(func(c *ui.Context) {
				th := theme.Matcha.For(dark)
				theme.WithTheme(c, th, func() {
					ui.Column(c).Fill().Background(th.Colors.Background).Padding(th.Space(6)).Gap(th.Space(4)).Children(func() {
						ds.HStack(c, 2, func() {
							ds.Heading(c, 1, "Hello Astryx").Grow(1)
							ds.Switch(c, &dark, "Dark")
						})
						ds.Card(c, func() {
							ds.FormField(c, "Project", "Lowercase slug", "", func() {
								ds.TextInput(c, &name, ds.InputOpts{Placeholder: "my-app"})
							})
							if ds.PrimaryButton(c, "Create", ds.ButtonOpts{LeadingIcon: icons.Sparkles}).Clicked() {
								c.Toast("Created " + name)
							}
						})
					})
					ds.ToastViewport(c)
				})
			}),
		})
	})
	if err := mygo.App.Run(); err != nil {
		log.Fatal(err)
	}
}
```

## CLI & MCP Toolkit (`astryx`)

```bash
# Bootstrap AGENTS.md / CLAUDE.md and .mcp.json in a project
./bin/go run ./cmd/astryx init --dir . --target both

# Search components, tokens, presets, guides, templates, or icons
./bin/go run ./cmd/astryx search dropdown --kind component

# Export the versioned offline capability manifest and pinned upstream parity map
./bin/go run ./cmd/astryx manifest

# Print compact token-efficient component docs for agents
./bin/go run ./cmd/astryx component AlertDialog --dense

# Inspect theme presets & WCAG contrast
./bin/go run ./cmd/astryx theme gothic --dark

# Scaffold a starter app from built-in templates (starter, settings-shell)
./bin/go run ./cmd/astryx template starter --out ./my-app

# Eject self-contained component source files into your project
./bin/go run ./cmd/astryx swizzle ModalCard --out ./ui/ds --pkg ds --dry-run
./bin/go run ./cmd/astryx swizzle ModalCard --out ./ui/ds --pkg ds

# Run Astryx hygiene rules + go vet + mygo vet
./bin/go run ./cmd/astryx check --dir . --vet

# Start the stdio MCP server (exposes `search` and `get`)
./bin/go run ./cmd/astryx mcp
```

### MCP Tools

When running `astryx mcp` (or configured via `.mcp.json` created by `astryx init`), coding agents get two tools:

- `search`: Filter across `component`, `token`, `preset`, `guide`, `template`, and `icon`.
- `get`: Retrieve full Go signatures, props, accessibility notes, token paths, icon usage, dark/light preset hex palettes, or template sources.

## MyGo Native UI Rules

`astryx check` and `go tool mygo vet .` enforce MyGo's frame-lifetime and theming invariants:

1. **No stored frame handles (`no-stored-frame-handle`)**: Never store `ui.Element`, `*ui.Context`, or `ui.*Parts` in structs, package globals, or goroutine closures.
2. **Scoped themes (`use-with-theme`)**: Always use `theme.WithTheme(c, th, fn)` instead of calling `ui.Themed` directly.
3. **Semantic colors (`hardcoded-color-in-ui`)**: Read colors from `theme.Of(c)` rather than hardcoding `ui.RGB` / `ui.RGBA` / `ui.Hex` in UI packages.
4. **Stable loop identity (`missing-key-in-loop`)**: Pass `c.Key(...)` when rendering stateful controls inside loops.
5. **Swizzle-safe `ds/` imports**: Exported files in `ds/` only import stdlib, `github.com/egoist/mygo/ui`, `github.com/amzkang/mygo-astryx/theme`, and `github.com/amzkang/mygo-astryx/icons`.

### Application layout and navigation

The P1a components provide application shells, desktop and mobile navigation, five-region layouts, grids, aspect ratios, links, scrolling, and overflow menus. Run the gallery and open **Patterns** for eight interactive compositions. See the [delivery record](docs/astryx-delivery-plan.md#p1a-本地交付与验收) for supported behavior and remaining upstream variants.

```go
var shell ds.AppShellState

func view(c *ui.Context) {
    theme.WithTheme(c, theme.Matcha.Light, func() {
        ds.AppShell(c, &shell, ds.AppShellSlots{
            TopNav: func() {
                ds.TopNav(c, ds.TopNavSlots{
                    Heading: func() { ds.Heading(c, 2, "Workspace") },
                })
            },
            SideNav: func() {
                ds.SideNav(c, &shell.Nav, []ds.NavItem{
                    {ID: "home", Label: "Home"},
                    {ID: "settings", Label: "Settings"},
                })
            },
            Content: func() { ds.Body(c, "Workspace content") },
        })
    }).Fill()
}
```

## Development

Start with [AGENTS.md](AGENTS.md) and [CONTRIBUTING.md](CONTRIBUTING.md) for repository rules, component acceptance criteria, and development workflows. See the [architecture guide](docs/architecture.md) for module boundaries and frame ownership, the [glossary](CONTEXT.md) for project terms, and the [roadmap](docs/astryx-parity-roadmap.md) for planned capabilities. Maintainer documentation is primarily in Chinese.

`astryx manifest` emits schema version 1 with components, templates, tokens, presets, guides, icon metadata, and 354 parity assessments against a fixed Astryx commit. Existing `search --json` arrays and `component --json` objects retain their shapes. Each parity record states the assessed coverage, remaining work, reusable MyGo primitives, and evidence. The [delivery plan](docs/astryx-delivery-plan.md) records the next component interfaces and all 58 page milestones.

## Verification

Run the full verification suite:

```bash
./scripts/verify.sh
```

Or run commands individually:

```bash
./bin/go build ./... && ./bin/go vet ./... && ./bin/go tool mygo vet .
./bin/go run ./cmd/astryx check --dir . --vet
./bin/go test -count=1 ./...
```

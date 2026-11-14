package main

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/amzkang/mygo-astryx/icons"
	"github.com/amzkang/mygo-astryx/registry"
	"github.com/amzkang/mygo-astryx/theme"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type searchInput struct {
	Query string `json:"query" jsonschema:"Search keywords (component name, token, theme, icon, guide, or template)"`
	Kind  string `json:"kind,omitempty" jsonschema:"Optional filter: component, token, preset, guide, template, or icon"`
}

type getInput struct {
	Kind  string `json:"kind" jsonschema:"Resource kind: component, guide, preset, template, token, or icon"`
	Name  string `json:"name" jsonschema:"Name of the resource to retrieve"`
	Dense bool   `json:"dense,omitempty" jsonschema:"For components, return compact token-efficient markdown"`
	Dark  bool   `json:"dark,omitempty" jsonschema:"For presets, return dark palette tokens"`
}

type iconDoc struct {
	Name   string `json:"name"`
	Symbol string `json:"symbol"`
	Usage  string `json:"usage"`
}

type presetOutput struct {
	Name       string  `json:"name"`
	Dark       bool    `json:"dark"`
	Density    string  `json:"density"`
	Background string  `json:"background"`
	Surface    string  `json:"surface"`
	Text       string  `json:"text"`
	Accent     string  `json:"accent"`
	OnAccent   string  `json:"onAccent"`
	Contrast   float64 `json:"contrast"`
}

func runMCPServer() error {
	return newMCPServer().Run(context.Background(), &mcp.StdioTransport{})
}

func newMCPServer() *mcp.Server {
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "mygo-astryx",
		Version: "0.1.0",
	}, nil)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "search",
		Description: "Search mygo-astryx native UI components, design tokens, themes, guides, templates, and icons.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, in searchInput) (*mcp.CallToolResult, any, error) {
		hits := registry.Search(in.Query, in.Kind)
		if len(hits) > 30 {
			hits = hits[:30]
		}
		var b strings.Builder
		if len(hits) == 0 {
			b.WriteString("No matches found.")
		} else {
			for _, h := range hits {
				fmt.Fprintf(&b, "- [%s] **%s**: %s\n", h.Kind, h.Name, h.Summary)
			}
		}
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: b.String()}},
		}, map[string]any{"hits": hits}, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get",
		Description: "Get full documentation, Go signature, accessibility notes, and example code for a component, guide, preset, template, token, or icon.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, in getInput) (*mcp.CallToolResult, any, error) {
		kind := strings.ToLower(strings.TrimSpace(in.Kind))
		if kind == "" {
			kind = "component"
		}
		switch kind {
		case "component":
			c, ok := registry.FindComponent(in.Name)
			if !ok {
				return nil, nil, fmt.Errorf("unknown component %q", in.Name)
			}
			txt := c.FormatDense()
			if !in.Dense && len(c.Tags) > 0 {
				txt += fmt.Sprintf("\nTags: %s\n", strings.Join(c.Tags, ", "))
			}
			return &mcp.CallToolResult{
				Content: []mcp.Content{&mcp.TextContent{Text: txt}},
			}, c, nil
		case "guide":
			g, ok := registry.FindGuide(in.Name)
			if !ok {
				return nil, nil, fmt.Errorf("unknown guide %q", in.Name)
			}
			return &mcp.CallToolResult{
				Content: []mcp.Content{&mcp.TextContent{Text: g.Body}},
			}, g, nil
		case "preset", "theme":
			ts := theme.ByName(in.Name)
			th := ts.For(in.Dark)
			out := presetOutput{
				Name:       th.Name,
				Dark:       th.Dark,
				Density:    fmt.Sprint(th.Density),
				Background: fmt.Sprintf("#%02X%02X%02X", th.Colors.Background.R, th.Colors.Background.G, th.Colors.Background.B),
				Surface:    fmt.Sprintf("#%02X%02X%02X", th.Colors.Surface.R, th.Colors.Surface.G, th.Colors.Surface.B),
				Text:       fmt.Sprintf("#%02X%02X%02X", th.Colors.Text.R, th.Colors.Text.G, th.Colors.Text.B),
				Accent:     fmt.Sprintf("#%02X%02X%02X", th.Colors.Accent.Base.R, th.Colors.Accent.Base.G, th.Colors.Accent.Base.B),
				OnAccent:   fmt.Sprintf("#%02X%02X%02X", th.Colors.Accent.OnColor.R, th.Colors.Accent.OnColor.G, th.Colors.Accent.OnColor.B),
				Contrast:   theme.ContrastRatio(th.Colors.Accent.OnColor, th.Colors.Accent.Base),
			}
			txt := fmt.Sprintf("# Theme %s (dark=%v)\nBackground: %s\nSurface: %s\nText: %s\nAccent: %s (OnColor: %s, Contrast: %.2f:1)\n",
				out.Name, out.Dark, out.Background, out.Surface, out.Text, out.Accent, out.OnAccent, out.Contrast)
			return &mcp.CallToolResult{
				Content: []mcp.Content{&mcp.TextContent{Text: txt}},
			}, out, nil
		case "template":
			tpl, ok := registry.FindTemplate(in.Name)
			if !ok {
				return nil, nil, fmt.Errorf("unknown template %q", in.Name)
			}
			var b strings.Builder
			fmt.Fprintf(&b, "# %s\n%s\n\n", tpl.Title, tpl.Description)
			paths := make([]string, 0, len(tpl.Files))
			for p := range tpl.Files {
				paths = append(paths, p)
			}
			sort.Strings(paths)
			for _, path := range paths {
				fmt.Fprintf(&b, "## %s\n```go\n%s\n```\n", path, tpl.Files[path])
			}
			return &mcp.CallToolResult{
				Content: []mcp.Content{&mcp.TextContent{Text: b.String()}},
			}, tpl, nil
		case "token":
			tok, ok := findToken(in.Name)
			if !ok {
				return nil, nil, fmt.Errorf("unknown token %q", in.Name)
			}
			txt := fmt.Sprintf("# Token `%s`\nGroup: %s\nType: `%s`\nDescription: %s\nUsage: `th := theme.Of(c); _ = th.%s`\n",
				tok.Path, tok.Group, tok.Type, tok.Description, tok.Path)
			return &mcp.CallToolResult{
				Content: []mcp.Content{&mcp.TextContent{Text: txt}},
			}, tok, nil
		case "icon":
			ic, ok := findIcon(in.Name)
			if !ok {
				return nil, nil, fmt.Errorf("unknown icon %q", in.Name)
			}
			doc := iconDoc{
				Name:   ic.Name,
				Symbol: ic.Symbol,
				Usage:  fmt.Sprintf("ui.Icon(c, icons.%s).Size(16, 16).TextColor(theme.Of(c).Colors.Text)", ic.Symbol),
			}
			txt := fmt.Sprintf("# Icon `%s` (`icons.%s`)\nUsage:\n```go\n%s\n```\n", doc.Name, doc.Symbol, doc.Usage)
			return &mcp.CallToolResult{
				Content: []mcp.Content{&mcp.TextContent{Text: txt}},
			}, doc, nil
		default:
			return nil, nil, fmt.Errorf("unsupported kind %q", in.Kind)
		}
	})

	return server
}

func findToken(name string) (registry.TokenDoc, bool) {
	q := strings.ToLower(strings.TrimSpace(name))
	if q == "" {
		return registry.TokenDoc{}, false
	}
	for _, t := range registry.Tokens {
		if strings.ToLower(t.Path) == q {
			return t, true
		}
	}
	for _, t := range registry.Tokens {
		if strings.HasSuffix(strings.ToLower(t.Path), "."+q) {
			return t, true
		}
	}
	return registry.TokenDoc{}, false
}

func findIcon(name string) (icons.IconEntry, bool) {
	q := strings.ToLower(strings.TrimSpace(name))
	if q == "" {
		return icons.IconEntry{}, false
	}
	for _, ic := range icons.All() {
		if strings.ToLower(ic.Name) == q || strings.ToLower(ic.Symbol) == q {
			return ic, true
		}
	}
	return icons.IconEntry{}, false
}

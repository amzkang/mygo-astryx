package main

import (
	"context"
	"strings"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestMCPServerSearchAndGet(t *testing.T) {
	ctx := context.Background()
	serverTransport, clientTransport := mcp.NewInMemoryTransports()

	server := newMCPServer()
	serverSession, err := server.Connect(ctx, serverTransport, nil)
	if err != nil {
		t.Fatalf("server.Connect failed: %v", err)
	}
	defer serverSession.Close()

	client := mcp.NewClient(&mcp.Implementation{
		Name:    "test-client",
		Version: "0.1.0",
	}, nil)
	session, err := client.Connect(ctx, clientTransport, nil)
	if err != nil {
		t.Fatalf("client.Connect failed: %v", err)
	}
	defer session.Close()

	toolsRes, err := session.ListTools(ctx, &mcp.ListToolsParams{})
	if err != nil {
		t.Fatalf("ListTools failed: %v", err)
	}
	if len(toolsRes.Tools) != 2 {
		t.Fatalf("expected 2 MCP tools, got %d", len(toolsRes.Tools))
	}

	// Search for components
	searchRes, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name: "search",
		Arguments: map[string]any{
			"query": "DropdownMenu",
			"kind":  "component",
		},
	})
	if err != nil {
		t.Fatalf("CallTool search failed: %v", err)
	}
	if searchRes.IsError || !strings.Contains(firstText(t, searchRes), "DropdownMenu") {
		t.Fatalf("expected DropdownMenu in search output, got %+v", searchRes)
	}

	// Get component (dense & non-dense)
	compRes, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name: "get",
		Arguments: map[string]any{
			"kind":  "component",
			"name":  "AlertDialog",
			"dense": true,
		},
	})
	if err != nil {
		t.Fatalf("CallTool get component failed: %v", err)
	}
	if compRes.IsError || !strings.Contains(firstText(t, compRes), "AlertDialog") {
		t.Fatalf("expected AlertDialog doc, got %+v", compRes)
	}

	// Get token by exact path and suffix
	for _, tokName := range []string{"Colors.Accent", "Accent"} {
		tokRes, err := session.CallTool(ctx, &mcp.CallToolParams{
			Name: "get",
			Arguments: map[string]any{
				"kind": "token",
				"name": tokName,
			},
		})
		if err != nil {
			t.Fatalf("CallTool get token %q failed: %v", tokName, err)
		}
		if tokRes.IsError || !strings.Contains(firstText(t, tokRes), "Colors.Accent") {
			t.Fatalf("expected Colors.Accent token doc, got %+v", tokRes)
		}
	}

	// Get icon by slug and Go symbol
	for _, iconName := range []string{"chevron-right", "ChevronRight"} {
		iconRes, err := session.CallTool(ctx, &mcp.CallToolParams{
			Name: "get",
			Arguments: map[string]any{
				"kind": "icon",
				"name": iconName,
			},
		})
		if err != nil {
			t.Fatalf("CallTool get icon %q failed: %v", iconName, err)
		}
		txt := firstText(t, iconRes)
		if iconRes.IsError || !strings.Contains(txt, "icons.ChevronRight") {
			t.Fatalf("expected icons.ChevronRight in icon doc, got %q", txt)
		}
	}

	// Get preset dark mode & guide & template
	presetRes, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name: "get",
		Arguments: map[string]any{
			"kind": "preset",
			"name": "gothic",
			"dark": true,
		},
	})
	if err != nil || presetRes.IsError || !strings.Contains(firstText(t, presetRes), "dark=true") {
		t.Fatalf("unexpected preset result: err=%v res=%+v", err, presetRes)
	}

	guideRes, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name: "get",
		Arguments: map[string]any{
			"kind": "guide",
			"name": "mygo-rules",
		},
	})
	if err != nil || guideRes.IsError || !strings.Contains(firstText(t, guideRes), "WithTheme") {
		t.Fatalf("unexpected guide result: err=%v res=%+v", err, guideRes)
	}

	tplRes, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name: "get",
		Arguments: map[string]any{
			"kind": "template",
			"name": "starter",
		},
	})
	if err != nil || tplRes.IsError || !strings.Contains(firstText(t, tplRes), "mygo.App.WhenReady") {
		t.Fatalf("unexpected template result: err=%v res=%+v", err, tplRes)
	}

	// Unknown resource returns tool error
	badRes, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name: "get",
		Arguments: map[string]any{
			"kind": "token",
			"name": "Does.Not.Exist",
		},
	})
	if err == nil && !badRes.IsError {
		t.Fatalf("expected error for unknown token, got %+v", badRes)
	}
}

func firstText(t *testing.T, res *mcp.CallToolResult) string {
	t.Helper()
	for _, c := range res.Content {
		if tc, ok := c.(*mcp.TextContent); ok {
			return tc.Text
		}
	}
	t.Fatalf("no TextContent in CallToolResult: %+v", res)
	return ""
}

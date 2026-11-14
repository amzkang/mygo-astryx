package main

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/amzkang/mygo-astryx/registry"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func requireJSONKeys(t *testing.T, object map[string]json.RawMessage, expected []string) {
	t.Helper()
	var keys []string
	for key := range object {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	sort.Strings(expected)
	if !reflect.DeepEqual(keys, expected) {
		t.Fatalf("JSON fields = %v, want %v", keys, expected)
	}
}

func TestCLIJSONAndMCPContracts(t *testing.T) {
	repoRoot, err := filepath.Abs("../..")
	if err != nil {
		t.Fatal(err)
	}
	binary := filepath.Join(t.TempDir(), "astryx")
	build := exec.Command(filepath.Join(repoRoot, "bin", "go"), "build", "-o", binary, ".")
	if output, err := build.CombinedOutput(); err != nil {
		t.Fatalf("build CLI: %v\n%s", err, output)
	}
	runCLI := func(t *testing.T, args ...string) []byte {
		t.Helper()
		cmd := exec.Command(binary, args...)
		cmd.Dir = t.TempDir()
		cmd.Env = append(os.Environ(), "GOPROXY=off", "GOSUMDB=off")
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		output, err := cmd.Output()
		if err != nil || stderr.Len() != 0 {
			t.Fatalf("CLI %v failed: %v\n%s", args, err, stderr.String())
		}
		if !json.Valid(output) {
			t.Fatalf("CLI %v must emit JSON only", args)
		}
		return output
	}

	t.Run("manifest offline and stable", func(t *testing.T) {
		first := runCLI(t, "manifest")
		if !bytes.Equal(first, runCLI(t, "manifest")) {
			t.Fatal("manifest output changed between offline invocations")
		}
		var object map[string]json.RawMessage
		if err := json.Unmarshal(first, &object); err != nil {
			t.Fatal(err)
		}
		requireJSONKeys(t, object, []string{"schemaVersion", "baseline", "components", "templates", "tokens", "presets", "guides", "icons", "parity"})
		var manifest registry.CapabilityManifest
		if err := json.Unmarshal(first, &manifest); err != nil {
			t.Fatal(err)
		}
		if manifest.SchemaVersion != 1 || len(manifest.Components) != len(registry.Components) || len(manifest.Parity) != 354 {
			t.Fatal("manifest omitted capabilities or parity records")
		}
		var icons []map[string]json.RawMessage
		if err := json.Unmarshal(object["icons"], &icons); err != nil {
			t.Fatal(err)
		}
		for _, icon := range icons {
			requireJSONKeys(t, icon, []string{"name", "symbol"})
		}
	})

	componentJSON := runCLI(t, "component", "primarybutton", "--json")
	var component map[string]json.RawMessage
	if err := json.Unmarshal(componentJSON, &component); err != nil {
		t.Fatal(err)
	}
	requireJSONKeys(t, component, []string{"name", "category", "summary", "signature", "sourceFiles", "tags", "a11y", "example"})
	if string(component["name"]) != `"PrimaryButton"` {
		t.Fatal("component lookup must retain its canonical name")
	}
	if output := strings.TrimSpace(string(runCLI(t, "search", "__astryx_no_match__", "--json"))); output != "null" {
		t.Fatalf("no-match search changed its existing null representation: %s", output)
	}
	for _, args := range [][]string{{"component", "Unknown", "--json"}, {"manifest", "--unknown"}} {
		cmd := exec.Command(binary, args...)
		var stdout, stderr bytes.Buffer
		cmd.Stdout, cmd.Stderr = &stdout, &stderr
		if err := cmd.Run(); err == nil || stdout.Len() != 0 || stderr.Len() == 0 {
			t.Fatalf("CLI %v must fail with stderr only", args)
		}
	}
	for _, help := range []string{"-h", "--help"} {
		cmd := exec.Command(binary, "manifest", help)
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		output, err := cmd.Output()
		if err != nil || stderr.Len() != 0 || !bytes.Contains(output, []byte("astryx manifest")) || !bytes.Contains(output, []byte("JSON")) {
			t.Fatalf("manifest %s must print help to stdout: %v\n%s\n%s", help, err, output, stderr.String())
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	serverTransport, clientTransport := mcp.NewInMemoryTransports()
	serverSession, err := newMCPServer().Connect(ctx, serverTransport, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer serverSession.Close()
	client := mcp.NewClient(&mcp.Implementation{Name: "contract-test", Version: "1"}, nil)
	session, err := client.Connect(ctx, clientTransport, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer session.Close()
	for _, dense := range []bool{false, true} {
		result, err := session.CallTool(ctx, &mcp.CallToolParams{
			Name: "get", Arguments: map[string]any{"kind": "component", "name": "primarybutton", "dense": dense},
		})
		if err != nil || result.IsError {
			t.Fatalf("MCP get failed: %v", err)
		}
		encoded, err := json.Marshal(result.StructuredContent)
		if err != nil {
			t.Fatal(err)
		}
		var cliValue, mcpValue any
		if err := json.Unmarshal(componentJSON, &cliValue); err != nil {
			t.Fatal(err)
		}
		if err := json.Unmarshal(encoded, &mcpValue); err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(cliValue, mcpValue) {
			t.Fatal("CLI component JSON and MCP get diverged")
		}
	}
	for _, query := range []struct{ text, kind string }{{"", ""}, {"button", "component"}, {"__astryx_no_match__", ""}} {
		var cliHits []registry.SearchHit
		if err := json.Unmarshal(runCLI(t, "search", query.text, "--kind", query.kind, "--json"), &cliHits); err != nil {
			t.Fatal(err)
		}
		result, err := session.CallTool(ctx, &mcp.CallToolParams{
			Name: "search", Arguments: map[string]any{"query": query.text, "kind": query.kind},
		})
		if err != nil || result.IsError {
			t.Fatalf("MCP search failed: %v", err)
		}
		encoded, err := json.Marshal(result.StructuredContent)
		if err != nil {
			t.Fatal(err)
		}
		var structured struct {
			Hits []registry.SearchHit `json:"hits"`
		}
		if err := json.Unmarshal(encoded, &structured); err != nil {
			t.Fatal(err)
		}
		if len(cliHits) > 30 {
			cliHits = cliHits[:30]
		}
		if !reflect.DeepEqual(cliHits, structured.Hits) {
			t.Fatal("CLI search and MCP search diverged after applying the existing MCP limit")
		}
	}
}

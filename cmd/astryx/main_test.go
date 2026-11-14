package main

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/amzkang/mygo-astryx/ds"
	"github.com/amzkang/mygo-astryx/registry"
)

func TestInitIdempotency(t *testing.T) {
	for _, target := range []string{"agents", "claude", "both"} {
		t.Run(target, func(t *testing.T) {
			tmp := t.TempDir()
			original := "# Project rules\n\nKeep this prefix.\n" + managedBegin + "\nOld managed instructions.\n" + managedEnd + "\nKeep this suffix.\n"
			for _, name := range []string{"AGENTS.md", "CLAUDE.md"} {
				if err := os.WriteFile(filepath.Join(tmp, name), []byte(original), 0o644); err != nil {
					t.Fatal(err)
				}
			}
			for range 2 {
				if err := cmdInit([]string{"--dir", tmp, "--target", target}); err != nil {
					t.Fatalf("init failed: %v", err)
				}
			}
			for _, name := range []string{"AGENTS.md", "CLAUDE.md"} {
				b, err := os.ReadFile(filepath.Join(tmp, name))
				if err != nil {
					t.Fatal(err)
				}
				body := string(b)
				selected := target == "both" || target == "agents" && name == "AGENTS.md" || target == "claude" && name == "CLAUDE.md"
				if !selected {
					if body != original {
						t.Fatalf("init changed unselected file %s", name)
					}
					continue
				}
				if strings.Count(body, managedBegin) != 1 || strings.Count(body, managedEnd) != 1 {
					t.Fatalf("%s expected exactly one managed block", name)
				}
				if !strings.HasPrefix(body, "# Project rules\n\nKeep this prefix.\n") || !strings.HasSuffix(body, "\nKeep this suffix.\n") || strings.Contains(body, "Old managed instructions.") {
					t.Fatalf("%s did not preserve user content and replace managed instructions", name)
				}
				for _, want := range []string{"go build ./...", "go vet ./...", "go tool mygo vet .", "go test -count=1 ./...", "Window.Update", "Window.Invalidate", "c.Key(...)"} {
					if !strings.Contains(body, want) {
						t.Errorf("%s missing generated instruction %q", name, want)
					}
				}
			}
		})
	}
}

func TestInitVerificationUsesTargetProject(t *testing.T) {
	for _, tc := range []struct {
		name  string
		files map[string]os.FileMode
		want  string
	}{
		{"plain", nil, "go build ./... && go vet ./... && go tool mygo vet . && go test -count=1 ./..."},
		{"wrapper", map[string]os.FileMode{"bin/go": 0o755}, "./bin/go build ./... && ./bin/go vet ./... && ./bin/go tool mygo vet . && ./bin/go test -count=1 ./..."},
		{"non-executable wrapper", map[string]os.FileMode{"bin/go": 0o644}, "go build ./... && go vet ./... && go tool mygo vet . && go test -count=1 ./..."},
		{"non-executable script", map[string]os.FileMode{"scripts/verify.sh": 0o644}, "bash ./scripts/verify.sh"},
		{"script", map[string]os.FileMode{"bin/go": 0o755, "scripts/verify.sh": 0o755}, "./scripts/verify.sh"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			tmp := t.TempDir()
			for name, mode := range tc.files {
				path := filepath.Join(tmp, name)
				if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(path, []byte("#!/bin/sh\n"), mode); err != nil {
					t.Fatal(err)
				}
				info, err := os.Stat(path)
				if err != nil {
					t.Fatal(err)
				}
				if mode&0o111 != 0 && info.Mode()&0o111 == 0 {
					t.Skip("filesystem does not preserve executable permission bits")
				}
			}
			// Re-initializing a fresh project must produce byte-identical instructions.
			var first []byte
			for range 2 {
				if err := cmdInit([]string{"--dir", tmp, "--target", "both", "--mcp=false"}); err != nil {
					t.Fatal(err)
				}
				for _, name := range []string{"AGENTS.md", "CLAUDE.md"} {
					body, err := os.ReadFile(filepath.Join(tmp, name))
					if err != nil {
						t.Fatal(err)
					}
					if !strings.Contains(string(body), "`"+tc.want+"`") {
						t.Fatalf("%s missing target verification %q", name, tc.want)
					}
					if first == nil {
						first = body
					}
					if string(body) != string(first) {
						t.Fatalf("%s generated instructions differ between files or runs", name)
					}
				}
			}
		})
	}
}

func TestSwizzleCompilesAndPassesMyGoVet(t *testing.T) {
	repoRoot, err := filepath.Abs("../..")
	if err != nil {
		t.Fatal(err)
	}
	selections := []string{"--all"}
	for _, component := range registry.Components {
		selections = append(selections, component.Name)
	}
	for _, selection := range selections {
		t.Run(selection, func(t *testing.T) {
			tmp := t.TempDir()
			ejectedDir := filepath.Join(tmp, "localds")
			if err := cmdSwizzle([]string{selection, "--out", ejectedDir, "--pkg", "localds"}); err != nil {
				t.Fatalf("swizzle failed: %v", err)
			}

			goMod := `module example.com/swizzledapp

go 1.27.3

require (
	github.com/egoist/mygo v0.8.3
	github.com/amzkang/mygo-astryx v0.0.0
)

replace github.com/amzkang/mygo-astryx => ` + repoRoot + `

tool github.com/egoist/mygo/cmd/mygo
`
			if err := os.WriteFile(filepath.Join(tmp, "go.mod"), []byte(goMod), 0o644); err != nil {
				t.Fatal(err)
			}
			symbol := selection
			if selection == "--all" {
				symbol = "ModalCard"
			}
			if selection == "Radio" || selection == "Select" {
				symbol += "[string]"
			}
			mainGo := "package main\n\nimport \"example.com/swizzledapp/localds\"\n\nfunc main() { _ = localds." + symbol + " }\n"
			if selection == "ModalCard" {
				mainGo = `package main

import (
	"example.com/swizzledapp/localds"
	"github.com/egoist/mygo/ui"
	"github.com/amzkang/mygo-astryx/theme"
)

func main() {
	open := true
	img := ui.Render(func(c *ui.Context) {
		theme.WithTheme(c, theme.Matcha.Light, func() {
			localds.ModalCard(c, &open, "Swizzled", 360, func() {
				localds.PrimaryButton(c, "OK")
			})
		})
	}, 400, 300, 1)
	if img == nil || img.Bounds().Dx() != 400 || img.Bounds().Dy() != 300 {
		panic("swizzled ModalCard did not render a frame")
	}
}
`
			}

			if selection == "AppShell" {
				mainGo = `package main

import (
    "example.com/swizzledapp/localds"
    "github.com/egoist/mygo/ui"
    "github.com/amzkang/mygo-astryx/theme"
)

func main() {
    for _, width := range []int{600, 1000} {
        var state localds.AppShellState
        img := ui.Render(func(c *ui.Context) {
            theme.WithTheme(c, theme.Matcha.Light, func() {
                localds.AppShell(c, &state, localds.AppShellSlots{
                    TopNav: func() { localds.TopNav(c, localds.TopNavSlots{Heading: func() { ui.Text(c, "App") }}) },
                    SideNav: func() { localds.SideNav(c, &state.Nav, []localds.NavItem{{ID: "home", Label: "Home"}}) },
                    Content: func() { ui.Text(c, "Content") },
                })
            }).Fill()
        }, width, 400, 1)
        if img == nil || img.Bounds().Dx() != width || img.Bounds().Dy() != 400 {
            panic("swizzled AppShell did not render")
        }
    }
}
`
			}

			if err := os.WriteFile(filepath.Join(tmp, "main.go"), []byte(mainGo), 0o644); err != nil {
				t.Fatal(err)
			}
			repoSum, err := os.ReadFile(filepath.Join(repoRoot, "go.sum"))
			if err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(filepath.Join(tmp, "go.sum"), repoSum, 0o644); err != nil {
				t.Fatal(err)
			}

			goBin := filepath.Join(repoRoot, "bin", "go")
			for _, cmdArgs := range [][]string{
				{"build", "-mod=mod", "./..."},
				{"vet", "-mod=mod", "./..."},
				{"tool", "mygo", "vet", "."},
			} {
				cmd := exec.Command(goBin, cmdArgs...)
				cmd.Dir = tmp
				out, err := cmd.CombinedOutput()
				if err != nil {
					t.Fatalf("go %v failed in swizzled temp module: %v\n%s", cmdArgs, err, string(out))
				}
			}
			if selection == "ModalCard" || selection == "AppShell" {
				cmd := exec.Command(goBin, "run", ".")
				cmd.Dir = tmp
				if output, err := cmd.CombinedOutput(); err != nil {
					t.Fatalf("swizzled render smoke failed: %v\n%s", err, output)
				}
			}
		})
	}
}

func TestCheckDirectoryCatchesViolationsAndPassesRepo(t *testing.T) {
	repoRoot, err := filepath.Abs("../..")
	if err != nil {
		t.Fatal(err)
	}
	issues, err := CheckDirectory(repoRoot)
	if err != nil {
		t.Fatal(err)
	}
	if len(issues) != 0 {
		t.Fatalf("expected 0 check issues in repo, got %+v", issues)
	}

	tmp := t.TempDir()
	badFile := `package bad
import "github.com/egoist/mygo/ui"
type Leaky struct {
	El ui.Element
}
`
	if err := os.WriteFile(filepath.Join(tmp, "bad.go"), []byte(badFile), 0o644); err != nil {
		t.Fatal(err)
	}
	badIssues, err := CheckDirectory(tmp)
	if err != nil {
		t.Fatal(err)
	}
	if len(badIssues) != 1 || badIssues[0].Rule != "no-stored-frame-handle" {
		t.Fatalf("expected 1 no-stored-frame-handle issue, got %+v", badIssues)
	}
}

func TestInitMCPMerge(t *testing.T) {
	tmp := t.TempDir()
	existing := `{"mcpServers":{"other":{"command":"echo","args":["hi"]}},"customFlag":true}`
	mcpPath := filepath.Join(tmp, ".mcp.json")
	if err := os.WriteFile(mcpPath, []byte(existing), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := cmdInit([]string{"--dir", tmp, "--target", "agents"}); err != nil {
		t.Fatalf("cmdInit failed: %v", err)
	}
	if err := cmdInit([]string{"--dir", tmp, "--target", "agents"}); err != nil {
		t.Fatalf("second cmdInit failed: %v", err)
	}

	raw, err := os.ReadFile(mcpPath)
	if err != nil {
		t.Fatal(err)
	}
	var parsed map[string]any
	if err := json.Unmarshal(raw, &parsed); err != nil {
		t.Fatal(err)
	}
	if parsed["customFlag"] != true {
		t.Fatalf("expected customFlag preserved, got %v", parsed["customFlag"])
	}
	servers, ok := parsed["mcpServers"].(map[string]any)
	if !ok {
		t.Fatalf("expected mcpServers map, got %T", parsed["mcpServers"])
	}
	if _, ok := servers["other"]; !ok {
		t.Fatalf("expected existing 'other' MCP server preserved, got %+v", servers)
	}
	if _, ok := servers["mygo-astryx"]; !ok {
		t.Fatalf("expected 'mygo-astryx' MCP server added, got %+v", servers)
	}

	// Verify --mcp=false skips creating .mcp.json
	tmpNoMCP := t.TempDir()
	if err := cmdInit([]string{"--dir", tmpNoMCP, "--target", "claude", "--mcp=false"}); err != nil {
		t.Fatalf("cmdInit --mcp=false failed: %v", err)
	}
	if _, err := os.Stat(filepath.Join(tmpNoMCP, ".mcp.json")); !os.IsNotExist(err) {
		t.Fatalf("expected no .mcp.json when --mcp=false, stat err=%v", err)
	}
}

func TestSwizzleDryRunForceAndAll(t *testing.T) {
	tmp := t.TempDir()
	outDir := filepath.Join(tmp, "ui", "ds")

	// --dry-run must not write files
	dryRes, err := Swizzle(SwizzleOptions{
		Component: "Button",
		OutDir:    outDir,
		PkgName:   "localds",
		DryRun:    true,
	})
	if err != nil {
		t.Fatalf("dry-run failed: %v", err)
	}
	if len(dryRes) == 0 || dryRes[0].Action != "would-create" {
		t.Fatalf("unexpected dry-run results: %+v", dryRes)
	}
	if _, err := os.Stat(outDir); !os.IsNotExist(err) {
		t.Fatalf("expected dry-run not to create %s", outDir)
	}

	// Real swizzle creates files
	if err := cmdSwizzle([]string{"Button", "--out", outDir, "--pkg", "localds"}); err != nil {
		t.Fatalf("cmdSwizzle Button failed: %v", err)
	}
	if _, err := os.Stat(filepath.Join(outDir, "button.go")); err != nil {
		t.Fatalf("expected button.go to exist: %v", err)
	}

	// Without --force, swizzling over existing files fails
	if err := cmdSwizzle([]string{"Button", "--out", outDir, "--pkg", "localds"}); err == nil {
		t.Fatal("expected error overwriting without --force")
	}

	// With --force and --all, all ds source files are ejected cleanly
	if err := cmdSwizzle([]string{"--all", "--force", "--out", outDir, "--pkg", "localds"}); err != nil {
		t.Fatalf("cmdSwizzle --all --force failed: %v", err)
	}
	entries, err := ds.SourceFS.ReadDir(".")
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".go") || strings.HasSuffix(entry.Name(), "_test.go") {
			continue
		}
		if _, err := os.Stat(filepath.Join(outDir, entry.Name())); err != nil {
			t.Fatalf("expected %s after --all swizzle: %v", entry.Name(), err)
		}
	}
}

func TestCheckDirectoryNewRules(t *testing.T) {
	tmp := t.TempDir()
	badCode := `package sample

import (
	"github.com/egoist/mygo/ui"
	"github.com/amzkang/mygo-astryx/ds"
	"github.com/amzkang/mygo-astryx/theme"
)

var globalParts ui.MenuParts

type BadState struct {
	Ctx *ui.Context
}

func View(c *ui.Context, items []string) {
	ui.Themed(ui.LightTheme(), func() {
		ui.Box(c).Background(ui.RGB(12, 34, 56))
		for _, item := range items {
			ds.Button(c, item)
		}
		for _, item := range items {
			ds.Button(c.Key("ok-"+item), item)
			theme.WithTheme(c.Key("th-"+item), theme.Matcha.Light, func() {
				ds.PrimaryButton(c, item)
			})
		}
	})
}
`
	if err := os.WriteFile(filepath.Join(tmp, "bad.go"), []byte(badCode), 0o644); err != nil {
		t.Fatal(err)
	}

	issues, err := CheckDirectory(tmp)
	if err != nil {
		t.Fatalf("CheckDirectory failed: %v", err)
	}
	gotRules := map[string]int{}
	for _, is := range issues {
		gotRules[is.Rule]++
	}
	if gotRules["no-stored-frame-handle"] != 2 {
		t.Fatalf("expected 2 no-stored-frame-handle issues (global + struct), got %+v (%+v)", gotRules, issues)
	}
	if gotRules["use-with-theme"] != 1 {
		t.Fatalf("expected 1 use-with-theme issue, got %+v", gotRules)
	}
	if gotRules["hardcoded-color-in-ui"] != 1 {
		t.Fatalf("expected 1 hardcoded-color-in-ui issue, got %+v", gotRules)
	}
	if gotRules["missing-key-in-loop"] != 1 {
		t.Fatalf("expected 1 missing-key-in-loop issue, got %+v (%+v)", gotRules, issues)
	}
}

func TestTemplatesE2EBuildVetAndCheck(t *testing.T) {
	repoRoot, err := filepath.Abs("../..")
	if err != nil {
		t.Fatal(err)
	}
	goBin := filepath.Join(repoRoot, "bin", "go")
	repoSum, _ := os.ReadFile(filepath.Join(repoRoot, "go.sum"))

	for _, tpl := range registry.Templates {
		t.Run(tpl.Name, func(t *testing.T) {
			appDir := t.TempDir()
			if err := cmdTemplate([]string{tpl.Name, "--out", appDir}); err != nil {
				t.Fatalf("cmdTemplate(%q) failed: %v", tpl.Name, err)
			}

			goMod := `module example.com/tpltest/` + strings.ReplaceAll(tpl.Name, "-", "") + `

go 1.27.3

require (
	github.com/egoist/mygo v0.8.3
	github.com/amzkang/mygo-astryx v0.0.0
)

replace github.com/amzkang/mygo-astryx => ` + repoRoot + `

tool github.com/egoist/mygo/cmd/mygo
`
			if err := os.WriteFile(filepath.Join(appDir, "go.mod"), []byte(goMod), 0o644); err != nil {
				t.Fatal(err)
			}
			if len(repoSum) > 0 {
				if err := os.WriteFile(filepath.Join(appDir, "go.sum"), repoSum, 0o644); err != nil {
					t.Fatal(err)
				}
			}

			issues, err := CheckDirectory(appDir)
			if err != nil {
				t.Fatalf("CheckDirectory failed: %v", err)
			}
			if len(issues) != 0 {
				t.Fatalf("expected 0 Astryx check issues for template %s, got %+v", tpl.Name, issues)
			}

			for _, cmdArgs := range [][]string{
				{"build", "-mod=mod", "./..."},
				{"vet", "-mod=mod", "./..."},
				{"tool", "mygo", "vet", "."},
			} {
				cmd := exec.Command(goBin, cmdArgs...)
				cmd.Dir = appDir
				out, err := cmd.CombinedOutput()
				if err != nil {
					t.Fatalf("template %s: go %v failed: %v\n%s", tpl.Name, cmdArgs, err, string(out))
				}
			}
		})
	}
}

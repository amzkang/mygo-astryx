package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/amzkang/mygo-astryx/ds"
	"github.com/amzkang/mygo-astryx/icons"
	"github.com/amzkang/mygo-astryx/registry"
	"github.com/amzkang/mygo-astryx/theme"
)

const (
	managedBegin = "<!-- BEGIN MYGO-ASTRYX -->"
	managedEnd   = "<!-- END MYGO-ASTRYX -->"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "astryx: %v\n", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) == 0 || args[0] == "-h" || args[0] == "--help" || args[0] == "help" {
		printUsage()
		return nil
	}
	sub := args[0]
	rest := args[1:]
	switch sub {
	case "init":
		return cmdInit(rest)
	case "search":
		return cmdSearch(rest)
	case "manifest":
		if len(rest) == 1 && (rest[0] == "-h" || rest[0] == "--help") {
			fmt.Print("Export the versioned capability manifest as JSON using embedded offline data.\n\nUsage\n  astryx manifest\n")
			return nil
		}
		if len(rest) != 0 {
			return fmt.Errorf("usage: astryx manifest")
		}
		return writeJSON(registry.BuildManifest())
	case "component":
		return cmdComponent(rest)
	case "docs":
		return cmdDocs(rest)
	case "theme":
		return cmdTheme(rest)
	case "swizzle":
		return cmdSwizzle(rest)
	case "template":
		return cmdTemplate(rest)
	case "check":
		return cmdCheck(rest)
	case "mcp":
		return runMCPServer()
	default:
		return fmt.Errorf("unknown subcommand %q (run 'astryx --help')", sub)
	}
}

func printUsage() {
	fmt.Print(`astryx — Agent-ready CLI & MCP toolkit for mygo-astryx

Usage:
  astryx init [--dir .] [--target agents|claude|both] [--mcp=true]
  astryx search <query> [--kind component|token|guide|preset|template|icon] [--json]
  astryx manifest
  astryx component <name> [--dense] [--json]
  astryx docs [guide] [--json]
  astryx theme [name] [--dark] [--json]
  astryx swizzle <component> [--out ./ui/ds] [--pkg ds] [--dry-run] [--force] [--all]
  astryx template <name> [--out .] [--list]
  astryx check [--dir .] [--vet] [--json]
  astryx mcp
`)
}

func cmdInit(args []string) error {
	fsFlags := flag.NewFlagSet("init", flag.ContinueOnError)
	dir := fsFlags.String("dir", ".", "Target project directory")
	target := fsFlags.String("target", "both", "Which instruction files to update: agents, claude, or both")
	withMCP := fsFlags.Bool("mcp", true, "Create or update .mcp.json with the mygo-astryx MCP server")
	if err := fsFlags.Parse(args); err != nil {
		return err
	}

	var targets []string
	switch strings.ToLower(*target) {
	case "agents":
		targets = []string{"AGENTS.md"}
	case "claude":
		targets = []string{"CLAUDE.md"}
	case "both":
		targets = []string{"AGENTS.md", "CLAUDE.md"}
	default:
		return fmt.Errorf("invalid --target %q", *target)
	}

	goCommand := "go"
	if info, err := os.Stat(filepath.Join(*dir, "bin", "go")); err == nil && info.Mode().IsRegular() && info.Mode()&0o111 != 0 {
		goCommand = "./bin/go"
	}
	verification := fmt.Sprintf("%s build ./... && %s vet ./... && %s tool mygo vet . && %s test -count=1 ./...", goCommand, goCommand, goCommand, goCommand)
	if info, err := os.Stat(filepath.Join(*dir, "scripts", "verify.sh")); err == nil && info.Mode().IsRegular() {
		verification = "bash ./scripts/verify.sh"
		if info.Mode()&0o111 != 0 {
			verification = "./scripts/verify.sh"
		}
	}

	block := strings.TrimSpace(`
` + managedBegin + `
## mygo-astryx UI Design System

- **Theme scope**. Wrap your root view in ` + "`theme.WithTheme(c, th, fn)`" + ` and read tokens with ` + "`theme.Of(c)`" + `.
- **Stock themes**. ` + "`theme.Neutral`, `theme.Matcha`, `theme.Stone`, `theme.Butter`, `theme.Chocolate`, `theme.Gothic`, `theme.Y2K`" + `.
- **MyGo lifetime rules**. Never store ` + "`ui.Element`" + `, ` + "`*ui.Context`" + `, or ` + "`*Parts`" + ` across frames or send them to goroutines.
- **UI thread**. Run I/O in goroutines and publish results via ` + "`Window.Update` / `Window.Invalidate`" + `. Keep frame handles on the UI thread.
- **Stable identity**. Use ` + "`c.Key(...)`" + ` before constructing repeated or dynamic stateful controls.
- **Discovery**. Run ` + "`astryx search <query>`" + ` or ` + "`astryx component <Name> --dense`" + ` before building custom widgets.
- **Verification**. After each batch of changes, run ` + "`" + verification + "`" + ` from the project root. Include build, Go vet, MyGo vet, and all tests in the project's verification script.
` + managedEnd)

	if err := os.MkdirAll(*dir, 0o755); err != nil {
		return err
	}
	for _, name := range targets {
		p := filepath.Join(*dir, name)
		if err := upsertManagedBlock(p, block); err != nil {
			return err
		}
		fmt.Printf("updated %s\n", p)
	}
	if *withMCP {
		mcpPath := filepath.Join(*dir, ".mcp.json")
		if err := upsertMCPConfig(mcpPath); err != nil {
			return err
		}
		fmt.Printf("updated %s\n", mcpPath)
	}
	return nil
}

func upsertMCPConfig(path string) error {
	root := map[string]any{}
	raw, err := os.ReadFile(path)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	if len(bytes.TrimSpace(raw)) > 0 {
		if err := json.Unmarshal(raw, &root); err != nil {
			return fmt.Errorf("parse %s: %w", path, err)
		}
	}
	servers, ok := root["mcpServers"].(map[string]any)
	if !ok || servers == nil {
		servers = map[string]any{}
	}
	servers["mygo-astryx"] = map[string]any{
		"command": "astryx",
		"args":    []string{"mcp"},
	}
	root["mcpServers"] = servers
	out, err := json.MarshalIndent(root, "", "  ")
	if err != nil {
		return err
	}
	out = append(out, '\n')
	return os.WriteFile(path, out, 0o644)
}

func upsertManagedBlock(path, block string) error {
	raw, err := os.ReadFile(path)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	content := string(raw)
	if start := strings.Index(content, managedBegin); start >= 0 {
		if end := strings.Index(content[start:], managedEnd); end >= 0 {
			endIdx := start + end + len(managedEnd)
			updated := content[:start] + block + content[endIdx:]
			return os.WriteFile(path, []byte(updated), 0o644)
		}
	}
	if len(content) > 0 && !strings.HasSuffix(content, "\n") {
		content += "\n"
	}
	if len(content) > 0 {
		content += "\n"
	}
	content += block + "\n"
	return os.WriteFile(path, []byte(content), 0o644)
}

func cmdSearch(args []string) error {
	fsFlags := flag.NewFlagSet("search", flag.ContinueOnError)
	kind := fsFlags.String("kind", "", "Filter by kind")
	asJSON := fsFlags.Bool("json", false, "Output JSON")
	query, err := parseFlagsMixed(fsFlags, args)
	if err != nil {
		return err
	}
	hits := registry.Search(query, *kind)
	if *asJSON {
		return writeJSON(hits)
	}
	if len(hits) == 0 {
		fmt.Println("No matches found.")
		return nil
	}
	for _, h := range hits {
		fmt.Printf("[%-9s] %-22s %s\n", h.Kind, h.Name, h.Summary)
	}
	return nil
}

func cmdComponent(args []string) error {
	fsFlags := flag.NewFlagSet("component", flag.ContinueOnError)
	dense := fsFlags.Bool("dense", false, "Emit compact token-efficient agent output")
	asJSON := fsFlags.Bool("json", false, "Output JSON")
	name, err := parseFlagsMixed(fsFlags, args)
	if err != nil {
		return err
	}
	if name == "" {
		return fmt.Errorf("usage: astryx component <name> [--dense] [--json]")
	}
	comp, ok := registry.FindComponent(name)
	if !ok {
		return fmt.Errorf("unknown component %q", name)
	}
	if *asJSON {
		return writeJSON(comp)
	}
	if *dense {
		fmt.Print(comp.FormatDense())
		return nil
	}
	fmt.Printf("%s (%s)\n%s\n\nSignature:\n  %s\n\nSource files:\n  %s\n\nAccessibility:\n  %s\n\nExample:\n%s\n",
		comp.Name, comp.Category, comp.Summary, comp.Signature, strings.Join(comp.SourceFiles, ", "), comp.A11y, comp.Example)
	return nil
}

func cmdDocs(args []string) error {
	fsFlags := flag.NewFlagSet("docs", flag.ContinueOnError)
	asJSON := fsFlags.Bool("json", false, "Output JSON")
	name, err := parseFlagsMixed(fsFlags, args)
	if err != nil {
		return err
	}
	if name == "" {
		if *asJSON {
			return writeJSON(registry.Guides)
		}
		for _, g := range registry.Guides {
			fmt.Printf("%-14s %s\n", g.Name, g.Summary)
		}
		return nil
	}
	g, ok := registry.FindGuide(name)
	if !ok {
		return fmt.Errorf("unknown guide %q", name)
	}
	if *asJSON {
		return writeJSON(g)
	}
	fmt.Println(g.Body)
	return nil
}

func cmdTheme(args []string) error {
	fsFlags := flag.NewFlagSet("theme", flag.ContinueOnError)
	dark := fsFlags.Bool("dark", false, "Inspect dark mode palette")
	asJSON := fsFlags.Bool("json", false, "Output JSON")
	name, err := parseFlagsMixed(fsFlags, args)
	if err != nil {
		return err
	}
	if name == "" {
		if *asJSON {
			return writeJSON(registry.Presets)
		}
		for _, p := range registry.Presets {
			fmt.Printf("%-11s light=(bg %s, accent %s) dark=(bg %s, accent %s) — %s\n",
				p.Name, p.LightBg, p.LightAccent, p.DarkBg, p.DarkAccent, p.Description)
		}
		return nil
	}
	ts := theme.ByName(name)
	th := ts.For(*dark)
	if *asJSON {
		return writeJSON(map[string]any{
			"name":       th.Name,
			"dark":       th.Dark,
			"density":    fmt.Sprint(th.Density),
			"background": fmt.Sprintf("#%02X%02X%02X", th.Colors.Background.R, th.Colors.Background.G, th.Colors.Background.B),
			"surface":    fmt.Sprintf("#%02X%02X%02X", th.Colors.Surface.R, th.Colors.Surface.G, th.Colors.Surface.B),
			"text":       fmt.Sprintf("#%02X%02X%02X", th.Colors.Text.R, th.Colors.Text.G, th.Colors.Text.B),
			"accent":     fmt.Sprintf("#%02X%02X%02X", th.Colors.Accent.Base.R, th.Colors.Accent.Base.G, th.Colors.Accent.Base.B),
			"onAccent":   fmt.Sprintf("#%02X%02X%02X", th.Colors.Accent.OnColor.R, th.Colors.Accent.OnColor.G, th.Colors.Accent.OnColor.B),
			"contrast":   theme.ContrastRatio(th.Colors.Accent.OnColor, th.Colors.Accent.Base),
		})
	}
	fmt.Printf("Theme: %s (dark=%v, density=%s)\n", th.Name, th.Dark, th.Density)
	fmt.Printf("Background: #%02X%02X%02X  Surface: #%02X%02X%02X  Text: #%02X%02X%02X\n",
		th.Colors.Background.R, th.Colors.Background.G, th.Colors.Background.B,
		th.Colors.Surface.R, th.Colors.Surface.G, th.Colors.Surface.B,
		th.Colors.Text.R, th.Colors.Text.G, th.Colors.Text.B)
	fmt.Printf("Accent: #%02X%02X%02X  OnColor: #%02X%02X%02X  Contrast: %.2f:1\n",
		th.Colors.Accent.Base.R, th.Colors.Accent.Base.G, th.Colors.Accent.Base.B,
		th.Colors.Accent.OnColor.R, th.Colors.Accent.OnColor.G, th.Colors.Accent.OnColor.B,
		theme.ContrastRatio(th.Colors.Accent.OnColor, th.Colors.Accent.Base))
	return nil
}

// SwizzleOptions configures ejecting ds source files into a local package.
type SwizzleOptions struct {
	Component string
	All       bool
	OutDir    string
	PkgName   string
	DryRun    bool
	Force     bool
}

// SwizzleResult describes a file created, overwritten, or previewed by Swizzle.
type SwizzleResult struct {
	File   string `json:"file"`
	Action string `json:"action"`
}

func cmdSwizzle(args []string) error {
	fsFlags := flag.NewFlagSet("swizzle", flag.ContinueOnError)
	outDir := fsFlags.String("out", "./ui/ds", "Destination directory for ejected files")
	pkgName := fsFlags.String("pkg", "ds", "Package name for ejected files")
	dryRun := fsFlags.Bool("dry-run", false, "Preview ejected files without writing to disk")
	force := fsFlags.Bool("force", false, "Overwrite existing files in --out")
	all := fsFlags.Bool("all", false, "Swizzle all ds component files")
	name, err := parseFlagsMixed(fsFlags, args)
	if err != nil {
		return err
	}
	if !*all && name == "" {
		return fmt.Errorf("usage: astryx swizzle <component> [--out ./ui/ds] [--pkg ds] [--dry-run] [--force] [--all]")
	}
	results, err := Swizzle(SwizzleOptions{
		Component: name,
		All:       *all,
		OutDir:    *outDir,
		PkgName:   *pkgName,
		DryRun:    *dryRun,
		Force:     *force,
	})
	if err != nil {
		return err
	}
	for _, r := range results {
		fmt.Printf("%s %s\n", r.Action, r.File)
	}
	if *dryRun {
		fmt.Println("dry-run: no files were modified")
	}
	return nil
}

// SwizzleComponent ejects a component's source files and transitive internal dependencies
// into outDir under pkgName, stripping embed.go boilerplate.
func SwizzleComponent(name, outDir, pkgName string) ([]string, error) {
	results, err := Swizzle(SwizzleOptions{
		Component: name,
		OutDir:    outDir,
		PkgName:   pkgName,
		Force:     true,
	})
	if err != nil {
		return nil, err
	}
	written := make([]string, len(results))
	for i, r := range results {
		written[i] = r.File
	}
	return written, nil
}

// Swizzle ejects component source files according to opts.
func Swizzle(opts SwizzleOptions) ([]SwizzleResult, error) {
	pkgName := opts.PkgName
	if pkgName == "" {
		pkgName = "ds"
	}
	if opts.OutDir == "" {
		opts.OutDir = "./ui/ds"
	}

	var files []string
	if opts.All {
		entries, err := fs.ReadDir(ds.SourceFS, ".")
		if err != nil {
			return nil, fmt.Errorf("list embedded ds files: %w", err)
		}
		for _, e := range entries {
			if e.IsDir() || !strings.HasSuffix(e.Name(), ".go") || strings.HasSuffix(e.Name(), "_test.go") {
				continue
			}
			files = append(files, e.Name())
		}
	} else {
		comp, ok := registry.FindComponent(opts.Component)
		if !ok {
			return nil, fmt.Errorf("unknown component %q", opts.Component)
		}
		files = append([]string(nil), comp.SourceFiles...)
	}
	sort.Strings(files)

	type preparedFile struct {
		dest   string
		body   []byte
		exists bool
	}

	fset := token.NewFileSet()
	prepared := make([]preparedFile, 0, len(files))
	for _, fileName := range files {
		raw, err := ds.SourceFS.ReadFile(fileName)
		if err != nil {
			return nil, fmt.Errorf("read embedded %s: %w", fileName, err)
		}
		f, err := parser.ParseFile(fset, fileName, raw, parser.ParseComments)
		if err != nil {
			return nil, fmt.Errorf("parse %s: %w", fileName, err)
		}
		f.Name.Name = pkgName
		stripSourceFSEmbed(f)

		for _, cg := range f.Comments {
			for _, cm := range cg.List {
				if strings.HasPrefix(cm.Text, "//go:embed ") {
					cm.Text = "// swizzled from mygo-astryx"
				}
			}
		}

		var buf bytes.Buffer
		if err := format.Node(&buf, fset, f); err != nil {
			return nil, fmt.Errorf("format %s: %w", fileName, err)
		}
		dest := filepath.Join(opts.OutDir, fileName)
		_, statErr := os.Stat(dest)
		exists := statErr == nil
		if exists && !opts.Force && !opts.DryRun {
			return nil, fmt.Errorf("%s already exists (pass --force to overwrite)", dest)
		}
		prepared = append(prepared, preparedFile{
			dest:   dest,
			body:   buf.Bytes(),
			exists: exists,
		})
	}

	if !opts.DryRun {
		if err := os.MkdirAll(opts.OutDir, 0o755); err != nil {
			return nil, err
		}
	}

	results := make([]SwizzleResult, 0, len(prepared))
	for _, pf := range prepared {
		action := "swizzled"
		if opts.DryRun {
			if pf.exists {
				if opts.Force {
					action = "would-overwrite"
				} else {
					action = "would-conflict"
				}
			} else {
				action = "would-create"
			}
		} else {
			if pf.exists {
				action = "overwrote"
			}
			if err := os.WriteFile(pf.dest, pf.body, 0o644); err != nil {
				return nil, err
			}
		}
		results = append(results, SwizzleResult{
			File:   pf.dest,
			Action: action,
		})
	}
	return results, nil
}

func stripSourceFSEmbed(f *ast.File) {
	// Remove SourceFS var and "embed" import from options.go when swizzled into user codebases.
	filtered := f.Decls[:0]
	f.Imports = nil
	for _, d := range f.Decls {
		gen, ok := d.(*ast.GenDecl)
		if ok && gen.Tok == token.VAR {
			keptSpecs := gen.Specs[:0]
			for _, s := range gen.Specs {
				vs, isVS := s.(*ast.ValueSpec)
				if isVS && len(vs.Names) == 1 && vs.Names[0].Name == "SourceFS" {
					continue
				}
				keptSpecs = append(keptSpecs, s)
			}
			gen.Specs = keptSpecs
			if len(gen.Specs) == 0 {
				continue
			}
		}
		if ok && gen.Tok == token.IMPORT {
			kept := gen.Specs[:0]
			for _, s := range gen.Specs {
				imp, ok := s.(*ast.ImportSpec)
				if ok && imp.Path != nil && strings.Trim(imp.Path.Value, "\"`") == "embed" {
					continue
				}
				kept = append(kept, s)
				if ok {
					f.Imports = append(f.Imports, imp)
				}
			}
			gen.Specs = kept
			if len(gen.Specs) == 0 {
				continue
			}
		}
		filtered = append(filtered, d)
	}
	f.Decls = filtered
}

func cmdTemplate(args []string) error {
	fsFlags := flag.NewFlagSet("template", flag.ContinueOnError)
	outDir := fsFlags.String("out", ".", "Output directory")
	list := fsFlags.Bool("list", false, "List templates")
	name, err := parseFlagsMixed(fsFlags, args)
	if err != nil {
		return err
	}
	if *list || name == "" {
		for _, tpl := range registry.Templates {
			fmt.Printf("%-16s %s — %s\n", tpl.Name, tpl.Title, tpl.Description)
		}
		return nil
	}
	tpl, ok := registry.FindTemplate(name)
	if !ok {
		return fmt.Errorf("unknown template %q", name)
	}
	if err := os.MkdirAll(*outDir, 0o755); err != nil {
		return err
	}
	for rel, body := range tpl.Files {
		dest := filepath.Join(*outDir, rel)
		if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(dest, []byte(body), 0o644); err != nil {
			return err
		}
		fmt.Printf("wrote %s\n", dest)
	}
	return nil
}

// CheckIssue represents a single static hygiene finding in a Go file.
type CheckIssue struct {
	File    string `json:"file"`
	Line    int    `json:"line"`
	Rule    string `json:"rule"`
	Message string `json:"message"`
}

type checkReport struct {
	Issues    []CheckIssue `json:"issues"`
	VetOK     *bool        `json:"vet_ok,omitempty"`
	VetOutput string       `json:"vet_output,omitempty"`
}

func cmdCheck(args []string) error {
	fsFlags := flag.NewFlagSet("check", flag.ContinueOnError)
	dir := fsFlags.String("dir", ".", "Directory to inspect")
	runVet := fsFlags.Bool("vet", false, "Also run go vet ./... and go tool mygo vet .")
	asJSON := fsFlags.Bool("json", false, "Output JSON")
	if err := fsFlags.Parse(args); err != nil {
		return err
	}
	issues, err := CheckDirectory(*dir)
	if err != nil {
		return err
	}

	var vetOK *bool
	var vetOut string
	if *runVet {
		ok, out := runGoAndMyGoVet(*dir)
		vetOK = &ok
		vetOut = out
	}

	if *asJSON {
		return writeJSON(checkReport{
			Issues:    issues,
			VetOK:     vetOK,
			VetOutput: vetOut,
		})
	}

	for _, is := range issues {
		fmt.Printf("%s:%d: [%s] %s\n", is.File, is.Line, is.Rule, is.Message)
	}
	if vetOK != nil && !*vetOK {
		fmt.Print(vetOut)
	}
	if len(issues) > 0 {
		return fmt.Errorf("found %d issue(s)", len(issues))
	}
	if vetOK != nil && !*vetOK {
		return fmt.Errorf("go vet / mygo vet failed")
	}
	vetNote := ""
	if vetOK != nil {
		vetNote = ", go vet + mygo vet passed"
	}
	fmt.Printf("ok: no Astryx hygiene issues found in %s (%d icons, %d components, %d themes%s)\n",
		*dir, len(icons.All()), len(registry.Components), len(theme.StockThemes()), vetNote)
	return nil
}

func runGoAndMyGoVet(dir string) (bool, string) {
	goBin := findGoBinary(dir)
	var outBuf strings.Builder
	cmds := [][]string{
		{"vet", "./..."},
		{"tool", "mygo", "vet", "."},
	}
	for _, args := range cmds {
		cmd := exec.Command(goBin, args...)
		cmd.Dir = dir
		b, err := cmd.CombinedOutput()
		if len(b) > 0 {
			outBuf.Write(b)
			if !strings.HasSuffix(string(b), "\n") {
				outBuf.WriteByte('\n')
			}
		}
		if err != nil {
			fmt.Fprintf(&outBuf, "%s %s: %v\n", filepath.Base(goBin), strings.Join(args, " "), err)
			return false, outBuf.String()
		}
	}
	return true, outBuf.String()
}

func findGoBinary(dir string) string {
	candidates := []string{
		filepath.Join(dir, "bin", "go"),
		"./bin/go",
	}
	for _, cand := range candidates {
		if st, err := os.Stat(cand); err == nil && !st.IsDir() {
			if abs, err := filepath.Abs(cand); err == nil {
				return abs
			}
			return cand
		}
	}
	return "go"
}

var statefulLoopControls = map[string]bool{
	"ds.Link":              true,
	"ui.Link":              true,
	"ds.SideNav":           true,
	"ds.MobileNav":         true,
	"ds.AppShell":          true,
	"ds.ScrollableArea":    true,
	"ds.MoreMenu":          true,
	"ui.ButtonBase":        true,
	"ui.CheckboxBase":      true,
	"ui.SwitchBase":        true,
	"ui.RadioBase":         true,
	"ui.TextInputBase":     true,
	"ui.TextAreaBase":      true,
	"ui.Slider":            true,
	"ui.RangeSlider":       true,
	"ui.NumberInput":       true,
	"ui.Stepper":           true,
	"ui.SelectBase":        true,
	"ui.ComboboxBase":      true,
	"ui.SegmentedBase":     true,
	"ui.DateInput":         true,
	"ui.PopoverBase":       true,
	"ui.TooltipBase":       true,
	"ui.ToastViewportBase": true,
	"ui.CollapsibleBase":   true,
	"ui.MenuBase":          true,
	"ui.ContextMenuBase":   true,
	"ui.DialogBase":        true,
	"ui.TabsBase":          true,
	"ui.AccordionBase":     true,
	"ui.TreeBase":          true,
	"ui.ToastBase":         true,
	"ds.Button":            true,
	"ds.PrimaryButton":     true,
	"ds.IconButton":        true,
	"ds.CopyButton":        true,
	"ds.Checkbox":          true,
	"ds.Switch":            true,
	"ds.Radio":             true,
	"ds.SegmentedControl":  true,
	"ds.TextInput":         true,
	"ds.TextArea":          true,
	"ds.SearchInput":       true,
	"ds.Select":            true,
	"ds.Combobox":          true,
	"ds.Slider":            true,
	"ds.RangeSlider":       true,
	"ds.NumberInput":       true,
	"ds.Stepper":           true,
	"ds.TokenField":        true,
	"ds.DateInput":         true,
	"ds.DropdownMenu":      true,
	"ds.Collapsible":       true,
	"ds.Accordion":         true,
	"ds.Tabs":              true,
	"ds.TreeView":          true,
	"ds.Pagination":        true,
	"ds.ContextMenu":       true,
	"ds.ModalCard":         true,
	"ds.Drawer":            true,
	"ds.Popover":           true,
	"ds.Tooltip":           true,
	"ds.ToastViewport":     true,
}

// CheckDirectory scans Go files under dir for common mygo-astryx mistakes:
//   - Storing ui.Element, *ui.Context, or ui.*Parts fields in structs or package globals.
//   - Using ui.Themed directly outside theme/scope.go instead of theme.WithTheme.
//   - Hardcoding ui.RGB / ui.RGBA / ui.Hex colors in UI packages outside theme/.
//   - Constructing stateful controls inside loops without c.Key(...).
func CheckDirectory(dir string) ([]CheckIssue, error) {
	var issues []CheckIssue
	fset := token.NewFileSet()
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			base := d.Name()
			if base == ".git" || base == ".gocache" || base == ".gopath" || base == "build" || base == "vendor" {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}
		f, perr := parser.ParseFile(fset, path, nil, 0)
		if perr != nil {
			return perr
		}
		rel, _ := filepath.Rel(dir, path)
		rel = filepath.ToSlash(rel)

		// Package-level var checks for frame-scoped handles.
		for _, decl := range f.Decls {
			gen, ok := decl.(*ast.GenDecl)
			if !ok || gen.Tok != token.VAR {
				continue
			}
			for _, spec := range gen.Specs {
				vs, ok := spec.(*ast.ValueSpec)
				if !ok || vs.Type == nil {
					continue
				}
				typeStr := exprString(vs.Type)
				if isForbiddenLifetimeType(typeStr) {
					pos := fset.Position(vs.Pos())
					issues = append(issues, CheckIssue{
						File:    rel,
						Line:    pos.Line,
						Rule:    "no-stored-frame-handle",
						Message: fmt.Sprintf("package variable has frame-scoped type %s", typeStr),
					})
				}
			}
		}

		isThemeFile := rel == "theme" || strings.HasPrefix(rel, "theme/")

		ast.Inspect(f, func(n ast.Node) bool {
			if st, ok := n.(*ast.StructType); ok && st.Fields != nil {
				for _, field := range st.Fields.List {
					typeStr := exprString(field.Type)
					if isForbiddenLifetimeType(typeStr) {
						pos := fset.Position(field.Pos())
						issues = append(issues, CheckIssue{
							File:    rel,
							Line:    pos.Line,
							Rule:    "no-stored-frame-handle",
							Message: fmt.Sprintf("struct field has frame-scoped type %s", typeStr),
						})
					}
				}
			}
			if call, ok := n.(*ast.CallExpr); ok {
				if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
					if ident, ok := sel.X.(*ast.Ident); ok && ident.Name == "ui" {
						if sel.Sel.Name == "Themed" && rel != "theme/scope.go" {
							pos := fset.Position(call.Pos())
							issues = append(issues, CheckIssue{
								File:    rel,
								Line:    pos.Line,
								Rule:    "use-with-theme",
								Message: "use theme.WithTheme(c, th, fn) instead of raw ui.Themed so Astryx tokens stay synchronized",
							})
						}
						if !isThemeFile && (sel.Sel.Name == "RGB" || sel.Sel.Name == "RGBA" || sel.Sel.Name == "Hex") {
							pos := fset.Position(call.Pos())
							issues = append(issues, CheckIssue{
								File:    rel,
								Line:    pos.Line,
								Rule:    "hardcoded-color-in-ui",
								Message: fmt.Sprintf("avoid hardcoded ui.%s(...) in UI code; read semantic tokens via theme.Of(c)", sel.Sel.Name),
							})
						}
					}
				}
			}
			switch loop := n.(type) {
			case *ast.ForStmt:
				checkLoopBodyKeys(fset, rel, loop.Body, &issues)
			case *ast.RangeStmt:
				checkLoopBodyKeys(fset, rel, loop.Body, &issues)
			}
			return true
		})
		return nil
	})
	return issues, err
}

func checkLoopBodyKeys(fset *token.FileSet, rel string, body *ast.BlockStmt, issues *[]CheckIssue) {
	if body == nil {
		return
	}
	ast.Inspect(body, func(n ast.Node) bool {
		// Nested loops will be inspected by the outer ast.Inspect walk, avoiding duplicate reports.
		if _, ok := n.(*ast.ForStmt); ok {
			return false
		}
		if _, ok := n.(*ast.RangeStmt); ok {
			return false
		}
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}
		callee := callName(call.Fun)
		// Virtualized lists/tables manage item identity via ListState.Key or row/col indexing.
		if callee == "ui.List" || callee == "ui.Table" || callee == "ds.DataTable" {
			return false
		}
		// Any call whose context argument is already keyed (c.Key(...)) scopes its subtree.
		if len(call.Args) > 0 && isKeyCall(call.Args[0]) {
			return false
		}
		if statefulLoopControls[callee] && len(call.Args) > 0 {
			if _, isBareIdent := call.Args[0].(*ast.Ident); isBareIdent {
				pos := fset.Position(call.Pos())
				*issues = append(*issues, CheckIssue{
					File:    rel,
					Line:    pos.Line,
					Rule:    "missing-key-in-loop",
					Message: fmt.Sprintf("%s inside loop uses unkeyed context; pass c.Key(...) for stable identity", callee),
				})
			}
		}
		return true
	})
}

func callName(fn ast.Expr) string {
	sel, ok := fn.(*ast.SelectorExpr)
	if !ok {
		return ""
	}
	ident, ok := sel.X.(*ast.Ident)
	if !ok {
		return ""
	}
	return ident.Name + "." + sel.Sel.Name
}

func isKeyCall(expr ast.Expr) bool {
	call, ok := expr.(*ast.CallExpr)
	if !ok {
		return false
	}
	sel, ok := call.Fun.(*ast.SelectorExpr)
	return ok && sel.Sel.Name == "Key"
}

func exprString(e ast.Expr) string {
	var buf bytes.Buffer
	_ = format.Node(&buf, token.NewFileSet(), e)
	return buf.String()
}

func isForbiddenLifetimeType(t string) bool {
	if t == "ui.Element" || t == "*ui.Context" {
		return true
	}
	return strings.HasPrefix(t, "ui.") && strings.HasSuffix(t, "Parts")
}

func parseFlagsMixed(fsFlags *flag.FlagSet, args []string) (string, error) {
	var flagArgs []string
	var posArgs []string
	for i := 0; i < len(args); i++ {
		a := args[i]
		if strings.HasPrefix(a, "-") {
			flagArgs = append(flagArgs, a)
			if !strings.Contains(a, "=") && i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
				// Peek if this flag expects a value.
				name := strings.TrimLeft(a, "-")
				if fl := fsFlags.Lookup(name); fl != nil {
					if bf, ok := fl.Value.(interface{ IsBoolFlag() bool }); !ok || !bf.IsBoolFlag() {
						i++
						flagArgs = append(flagArgs, args[i])
					}
				}
			}
			continue
		}
		posArgs = append(posArgs, a)
	}
	if err := fsFlags.Parse(flagArgs); err != nil {
		return "", err
	}
	posArgs = append(posArgs, fsFlags.Args()...)
	if len(posArgs) > 0 {
		return posArgs[0], nil
	}
	return "", nil
}

func writeJSON(v any) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}

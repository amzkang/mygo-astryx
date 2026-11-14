package registry

import (
	"fmt"
	"sort"
	"strings"
)

// PropDoc describes a single parameter or option field on a component.
type PropDoc struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Default     string `json:"default,omitempty"`
	Description string `json:"description"`
	Required    bool   `json:"required,omitempty"`
}

// VariantDoc describes an enumerated variant, size, or tone.
type VariantDoc struct {
	Name        string `json:"name"`
	Values      string `json:"values"`
	Description string `json:"description"`
}

// ComponentDoc captures Astryx-style structured metadata for both humans and coding agents.
type ComponentDoc struct {
	Name        string       `json:"name"`
	Category    string       `json:"category"`
	Summary     string       `json:"summary"`
	Signature   string       `json:"signature"`
	SourceFiles []string     `json:"sourceFiles"`
	Tags        []string     `json:"tags,omitempty"`
	Props       []PropDoc    `json:"props,omitempty"`
	Variants    []VariantDoc `json:"variants,omitempty"`
	A11y        string       `json:"a11y"`
	Rules       []string     `json:"rules,omitempty"`
	Example     string       `json:"example"`
}

// TokenDoc describes a design token path exposed by theme.Theme.
type TokenDoc struct {
	Path        string `json:"path"`
	Group       string `json:"group"`
	Type        string `json:"type"`
	Description string `json:"description"`
}

// GuideDoc captures architecture and convention guides for agents.
type GuideDoc struct {
	Name    string   `json:"name"`
	Title   string   `json:"title"`
	Summary string   `json:"summary"`
	Body    string   `json:"body"`
	Tags    []string `json:"tags,omitempty"`
}

// TemplateDoc provides a ready-to-run starter app template.
type TemplateDoc struct {
	Name        string   `json:"name"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Tags        []string `json:"tags,omitempty"`
	Files       map[string]string `json:"files"`
}

// PresetDoc summarizes one of the 7 stock Astryx themes.
type PresetDoc struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	LightBg     string `json:"lightBg"`
	LightAccent string `json:"lightAccent"`
	DarkBg      string `json:"darkBg"`
	DarkAccent  string `json:"darkAccent"`
}

// SearchHit is a unified search result returned by Search and MCP search.
type SearchHit struct {
	Kind    string `json:"kind"` // component, token, guide, preset, template, icon
	Name    string `json:"name"`
	Summary string `json:"summary"`
	Score   int    `json:"score"`
}

// FormatDense renders a compact token-efficient specification for LLM coding agents.
func (c ComponentDoc) FormatDense() string {
	var b strings.Builder
	fmt.Fprintf(&b, "# %s (%s)\n", c.Name, c.Category)
	fmt.Fprintf(&b, "%s\n\n", c.Summary)
	fmt.Fprintf(&b, "Signature: `%s`\n", c.Signature)
	fmt.Fprintf(&b, "Files: %s\n", strings.Join(c.SourceFiles, ", "))
	if len(c.Props) > 0 {
		b.WriteString("\nProps:\n")
		for _, p := range c.Props {
			req := ""
			if p.Required {
				req = " [required]"
			}
			def := ""
			if p.Default != "" {
				def = " (default: " + p.Default + ")"
			}
			fmt.Fprintf(&b, "- %s: %s%s%s — %s\n", p.Name, p.Type, req, def, p.Description)
		}
	}
	if len(c.Variants) > 0 {
		b.WriteString("\nVariants:\n")
		for _, v := range c.Variants {
			fmt.Fprintf(&b, "- %s: %s — %s\n", v.Name, v.Values, v.Description)
		}
	}
	if c.A11y != "" {
		fmt.Fprintf(&b, "\nA11y: %s\n", c.A11y)
	}
	if len(c.Rules) > 0 {
		b.WriteString("\nRules:\n")
		for _, r := range c.Rules {
			fmt.Fprintf(&b, "- %s\n", r)
		}
	}
	fmt.Fprintf(&b, "\nExample:\n```go\n%s\n```\n", strings.TrimSpace(c.Example))
	return b.String()
}

// FindComponent looks up a component by case-insensitive name.
func FindComponent(name string) (ComponentDoc, bool) {
	key := strings.ToLower(strings.TrimSpace(name))
	for _, c := range Components {
		if strings.ToLower(c.Name) == key {
			return c, true
		}
	}
	return ComponentDoc{}, false
}

// FindGuide looks up a guide by case-insensitive name.
func FindGuide(name string) (GuideDoc, bool) {
	key := strings.ToLower(strings.TrimSpace(name))
	for _, g := range Guides {
		if strings.ToLower(g.Name) == key {
			return g, true
		}
	}
	return GuideDoc{}, false
}

// FindTemplate looks up a template by case-insensitive name.
func FindTemplate(name string) (TemplateDoc, bool) {
	key := strings.ToLower(strings.TrimSpace(name))
	for _, t := range Templates {
		if strings.ToLower(t.Name) == key {
			return t, true
		}
	}
	return TemplateDoc{}, false
}

// Search performs a ranked search across components, tokens, presets, guides, templates, and icons.
func Search(query string, kind string) []SearchHit {
	q := strings.ToLower(strings.TrimSpace(query))
	k := strings.ToLower(strings.TrimSpace(kind))
	var hits []SearchHit

	add := func(hitKind, name, summary string, tags []string) {
		if k != "" && k != hitKind {
			return
		}
		score := matchScore(q, name, summary, tags)
		if q == "" || score > 0 {
			hits = append(hits, SearchHit{Kind: hitKind, Name: name, Summary: summary, Score: score})
		}
	}

	for _, c := range Components {
		add("component", c.Name, c.Summary, append([]string{c.Category}, c.Tags...))
	}
	for _, t := range Tokens {
		add("token", t.Path, t.Description, []string{t.Group, t.Type})
	}
	for _, p := range Presets {
		add("preset", p.Name, p.Description, []string{"theme", "palette"})
	}
	for _, g := range Guides {
		add("guide", g.Name, g.Summary, g.Tags)
	}
	for _, tpl := range Templates {
		add("template", tpl.Name, tpl.Description, tpl.Tags)
	}
	for _, ic := range IconNames() {
		add("icon", ic, "Embedded 24x24 Lucide-style SVG icon", []string{"svg", "icon"})
	}

	sort.Slice(hits, func(i, j int) bool {
		if hits[i].Score != hits[j].Score {
			return hits[i].Score > hits[j].Score
		}
		if hits[i].Kind != hits[j].Kind {
			return hits[i].Kind < hits[j].Kind
		}
		return hits[i].Name < hits[j].Name
	})
	return hits
}

func matchScore(q, name, summary string, tags []string) int {
	if q == "" {
		return 1
	}
	ln := strings.ToLower(name)
	ls := strings.ToLower(summary)
	if ln == q {
		return 100
	}
	score := 0
	if strings.HasPrefix(ln, q) {
		score += 60
	} else if strings.Contains(ln, q) {
		score += 40
	}
	for _, t := range tags {
		lt := strings.ToLower(t)
		if lt == q {
			score += 30
		} else if strings.Contains(lt, q) {
			score += 15
		}
	}
	if strings.Contains(ls, q) {
		score += 10
	}
	return score
}

package registry

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"path"
	"sort"
	"strings"

	"github.com/amzkang/mygo-astryx/icons"
)

const upstreamCommit = "e364afc958609c990d42514b8927a6ef54b86b00"

// Baseline identifies the immutable upstream inventory used for parity assessment.
type Baseline struct {
	Repository  string            `json:"repository"`
	Commit      string            `json:"commit"`
	MyGoVersion string            `json:"mygoVersion"`
	Paths       map[string]string `json:"paths"`
}

// ResourceRef identifies an implemented resource by its canonical registry name.
type ResourceRef struct {
	Kind string `json:"kind"`
	Name string `json:"name"`
}

// PrimitiveRef identifies an upstream MyGo API that can support a capability.
type PrimitiveRef struct {
	Package string `json:"package"`
	Symbol  string `json:"symbol"`
}

// ParityEntry records a source directory's assessed coverage and remaining work.
type ParityEntry struct {
	Kind       string         `json:"kind"`
	Category   string         `json:"category,omitempty"`
	SourcePath string         `json:"sourcePath"`
	Status     string         `json:"status"`
	LocalRefs  []ResourceRef  `json:"localRefs"`
	Primitives []PrimitiveRef `json:"primitives"`
	Reason     string         `json:"reason"`
	Evidence   []string       `json:"evidence"`
}

// ManifestIcon describes an icon without embedding parsed rendering objects.
type ManifestIcon struct {
	Name   string `json:"name"`
	Symbol string `json:"symbol"`
}

// CapabilityManifest is the versioned, offline inventory exposed by astryx manifest.
type CapabilityManifest struct {
	SchemaVersion int            `json:"schemaVersion"`
	Baseline      Baseline       `json:"baseline"`
	Components    []ComponentDoc `json:"components"`
	Templates     []TemplateDoc  `json:"templates"`
	Tokens        []TokenDoc     `json:"tokens"`
	Presets       []PresetDoc    `json:"presets"`
	Guides        []GuideDoc     `json:"guides"`
	Icons         []ManifestIcon `json:"icons"`
	Parity        []ParityEntry  `json:"parity"`
}

type parityCatalog struct {
	Baseline Baseline      `json:"baseline"`
	Entries  []ParityEntry `json:"entries"`
}

//go:embed parity.json
var parityJSON []byte

// BuildManifest returns a deterministic snapshot without reordering registry data.
func BuildManifest() CapabilityManifest {
	catalog, err := decodeParity(parityJSON)
	if err != nil {
		panic("registry: invalid embedded parity catalog: " + err.Error())
	}
	m := CapabilityManifest{
		SchemaVersion: 1,
		Baseline:      catalog.Baseline,
		Components:    append([]ComponentDoc{}, Components...),
		Templates:     append([]TemplateDoc{}, Templates...),
		Tokens:        append([]TokenDoc{}, Tokens...),
		Presets:       append([]PresetDoc{}, Presets...),
		Guides:        append([]GuideDoc{}, Guides...),
		Icons:         []ManifestIcon{},
		Parity:        catalog.Entries,
	}
	for _, icon := range icons.All() {
		m.Icons = append(m.Icons, ManifestIcon{Name: icon.Name, Symbol: icon.Symbol})
	}
	sort.Slice(m.Components, func(i, j int) bool { return m.Components[i].Name < m.Components[j].Name })
	sort.Slice(m.Templates, func(i, j int) bool { return m.Templates[i].Name < m.Templates[j].Name })
	sort.Slice(m.Tokens, func(i, j int) bool { return m.Tokens[i].Path < m.Tokens[j].Path })
	sort.Slice(m.Presets, func(i, j int) bool { return m.Presets[i].Name < m.Presets[j].Name })
	sort.Slice(m.Guides, func(i, j int) bool { return m.Guides[i].Name < m.Guides[j].Name })
	sort.Slice(m.Parity, func(i, j int) bool {
		if m.Parity[i].Kind != m.Parity[j].Kind {
			return m.Parity[i].Kind < m.Parity[j].Kind
		}
		return m.Parity[i].SourcePath < m.Parity[j].SourcePath
	})
	return m
}

func decodeParity(data []byte) (parityCatalog, error) {
	var catalog parityCatalog
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&catalog); err != nil {
		return catalog, err
	}
	if err := decoder.Decode(new(any)); err != io.EOF {
		return catalog, fmt.Errorf("expected a single parity document")
	}
	if err := validateParity(catalog); err != nil {
		return catalog, err
	}
	return catalog, nil
}

func validateParity(catalog parityCatalog) error {
	b := catalog.Baseline
	if b.Repository != "meta/astryx" || b.Commit != upstreamCommit || b.MyGoVersion != "v0.8.3" {
		return fmt.Errorf("unexpected parity baseline")
	}
	expectedPaths := map[string]string{
		"core":  "packages/core/src",
		"block": "packages/cli/assets/templates/blocks/components",
		"page":  "packages/cli/assets/templates/pages",
	}
	if len(b.Paths) != len(expectedPaths) {
		return fmt.Errorf("unexpected baseline paths")
	}
	for kind, expected := range expectedPaths {
		if b.Paths[kind] != expected {
			return fmt.Errorf("unexpected baseline path for %s", kind)
		}
	}
	known := map[ResourceRef]bool{}
	for _, hit := range Search("", "") {
		known[ResourceRef{Kind: hit.Kind, Name: hit.Name}] = true
	}
	seen := map[string]bool{}
	counts := map[string]int{}
	for _, entry := range catalog.Entries {
		key := entry.Kind + ":" + entry.SourcePath
		root, ok := expectedPaths[entry.Kind]
		if !ok || path.Dir(entry.SourcePath) != root || seen[key] {
			return fmt.Errorf("invalid or duplicate parity entry %s", key)
		}
		seen[key] = true
		counts[entry.Kind]++
		if entry.Kind == "core" {
			switch entry.Category {
			case "component", "support", "internal":
			default:
				return fmt.Errorf("%s has unknown core category %q", key, entry.Category)
			}
		} else if entry.Category != "" {
			return fmt.Errorf("%s has a core-only category", key)
		}
		switch entry.Status {
		case "implemented", "partial":
			if len(entry.LocalRefs) == 0 {
				return fmt.Errorf("%s requires local references", key)
			}
		case "missing", "deferred", "notApplicable":
		default:
			return fmt.Errorf("%s has unknown status %q", key, entry.Status)
		}
		if strings.TrimSpace(entry.Reason) == "" || len(entry.Evidence) == 0 || entry.LocalRefs == nil || entry.Primitives == nil {
			return fmt.Errorf("%s requires reason, evidence and reference arrays", key)
		}
		for _, ref := range entry.LocalRefs {
			if !known[ref] {
				return fmt.Errorf("%s references unknown resource %s:%s", key, ref.Kind, ref.Name)
			}
		}
		for _, ref := range entry.Primitives {
			if ref.Package == "" || ref.Symbol == "" {
				return fmt.Errorf("%s has an incomplete primitive reference", key)
			}
		}
		for _, evidence := range entry.Evidence {
			if strings.TrimSpace(evidence) == "" {
				return fmt.Errorf("%s has empty evidence", key)
			}
		}
	}
	if counts["core"] != 124 || counts["block"] != 172 || counts["page"] != 58 {
		return fmt.Errorf("incomplete parity inventory: %v", counts)
	}
	return nil
}

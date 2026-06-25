package compilehtml

import (
	"path/filepath"
	"strings"

	"github.com/Dcrispim/psrt.core/compilesvg"
	"github.com/Dcrispim/psrt.core/psrt"
	"github.com/Dcrispim/psrt.core/psrt/editor"
)

// Variant is one PSRT document bundled into a compiled HTML file.
type Variant struct {
	Label string
	Doc   psrt.Document
}

// VariantPSRT is a variant supplied as raw PSRT text (e.g. from a browser file picker).
type VariantPSRT struct {
	Label   string
	Content string
}

// LoadVariantsFromPaths parses each path and resolves constants; primary path should be first.
func LoadVariantsFromPaths(paths []string) ([]Variant, error) {
	if len(paths) == 0 {
		return nil, nil
	}
	out := make([]Variant, 0, len(paths))
	for _, p := range paths {
		doc, err := editor.LoadDocument(p)
		if err != nil {
			return nil, err
		}
		doc = compilesvg.ResolveDocument(doc)
		label := filepath.Base(p)
		out = append(out, Variant{Label: label, Doc: doc})
	}
	return out, nil
}

// LoadVariantsFromPSRT parses each body and resolves constants.
func LoadVariantsFromPSRT(items []VariantPSRT) ([]Variant, error) {
	if len(items) == 0 {
		return nil, nil
	}
	out := make([]Variant, 0, len(items))
	for _, item := range items {
		doc, err := psrt.Parse(strings.NewReader(item.Content))
		if err != nil {
			return nil, err
		}
		doc = compilesvg.ResolveDocument(doc)
		label := strings.TrimSpace(item.Label)
		if label == "" {
			label = "variant.psrt"
		}
		out = append(out, Variant{Label: label, Doc: doc})
	}
	return out, nil
}

func pageByName(doc psrt.Document, name string) *psrt.Page {
	for i := range doc.Pages {
		if doc.Pages[i].Name == name {
			return &doc.Pages[i]
		}
	}
	return nil
}

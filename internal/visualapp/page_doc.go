package visualapp

import (
	"os"
	"path/filepath"
	"strings"

	"psrt/psrt"
	"psrt/psrt/editor"
)

// extractPageDocument returns a copy of doc containing only the named page (fonts/consts preserved).
func extractPageDocument(doc psrt.Document, pageName string) (psrt.Document, error) {
	p, err := editor.FindPage(&doc, pageName)
	if err != nil {
		return psrt.Document{}, err
	}
	return psrt.Document{
		Pages:  []psrt.Page{*p},
		Fonts:  append([]string(nil), doc.Fonts...),
		Consts: cloneConsts(doc.Consts),
	}, nil
}

func cloneConsts(m map[string]string) map[string]string {
	if len(m) == 0 {
		return nil
	}
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}

// writeIntermediatePagePSRT saves a single-page PSRT next to the asset cache for preview compiles.
func writeIntermediatePagePSRT(storeRoot, pageName string, doc psrt.Document) (string, error) {
	dir := filepath.Join(storeRoot, "preview")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	safe := sanitizePageFilename(pageName)
	path := filepath.Join(dir, safe+".psrt")
	data, err := editor.FormatDocument(&doc)
	if err != nil {
		return "", err
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return "", err
	}
	return path, nil
}

func sanitizePageFilename(name string) string {
	name = strings.TrimSpace(name)
	if name == "" {
		return "page"
	}
	var b strings.Builder
	for _, r := range name {
		switch {
		case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z', r >= '0' && r <= '9', r == '-', r == '_':
			b.WriteRune(r)
		default:
			b.WriteRune('_')
		}
	}
	if b.Len() == 0 {
		return "page"
	}
	return b.String()
}

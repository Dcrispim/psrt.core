package visualapp

import (
	"context"
	"os"
	"path/filepath"

	"psrt/compilehtml"
	"psrt/compileopts"
	"psrt/compilesvg"
)

// ExportSVGFromDocument writes one SVG per page into parentDir/baseName/.
func (a *App) ExportSVGFromDocument(docJSON, parentDir, baseName string) (CompileSVGResult, error) {
	doc, err := parseDocumentJSON(docJSON)
	if err != nil {
		return CompileSVGResult{}, err
	}
	if baseName == "" {
		baseName = "document"
	}
	outDir := filepath.Join(parentDir, baseName)

	a.mu.Lock()
	store := a.store
	path := a.filePath
	client := a.client
	a.mu.Unlock()

	if store == nil && path != "" {
		store, _ = openStoreForPath(path)
	}
	ctx := context.Background()
	if store != nil {
		_ = store.EnsureDocument(ctx, client, doc)
	}
	batch, err := compilesvg.CompileWithCache(ctx, doc, client, outDir, store)
	if err != nil {
		return CompileSVGResult{}, err
	}
	return CompileSVGResult{
		URI:                outDir,
		UsedGoTextFallback: batch.UsedGoTextFallback,
	}, nil
}

// ExportHTMLFromDocument writes a single HTML file into dir as baseName.html.
// morePaths are additional .psrt files; morePSRT are variants read in the UI when paths are unavailable.
func (a *App) ExportHTMLFromDocument(docJSON, dir, baseName string, morePaths []string, morePSRT []VariantPSRT) (string, error) {
	doc, err := parseDocumentJSON(docJSON)
	if err != nil {
		return "", err
	}
	if baseName == "" {
		baseName = "document"
	}
	outPath := filepath.Join(dir, baseName+".html")

	a.mu.Lock()
	store := a.store
	path := a.filePath
	client := a.client
	a.mu.Unlock()

	if store == nil && path != "" {
		store, _ = openStoreForPath(path)
	}
	ctx := context.Background()
	if store != nil {
		_ = store.EnsureDocument(ctx, client, doc)
	}
	extras := make([]compilehtml.VariantPSRT, len(morePSRT))
	for i, v := range morePSRT {
		extras[i] = compilehtml.VariantPSRT{Label: v.Label, Content: v.Content}
	}
	html, err := compilehtml.CompileWithCacheFrom(ctx, doc, path, morePaths, extras, client, store, compileopts.Options{})
	if err != nil {
		return "", err
	}
	if err := os.WriteFile(outPath, html, 0o644); err != nil {
		return "", err
	}
	return outPath, nil
}

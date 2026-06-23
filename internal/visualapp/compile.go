package visualapp

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"psrt/compilehtml"
	"psrt/compileopts"
	"psrt/compilesvg"
)

var compileDebounce = 500 * time.Millisecond

func (a *App) scheduleAutoCompile() {
	a.compileMu.Lock()
	defer a.compileMu.Unlock()
	if a.compileTimer != nil {
		a.compileTimer.Stop()
	}
	page := a.activePage
	a.compileTimer = time.AfterFunc(compileDebounce, func() {
		res, err := a.CompilePageSVG(page)
		if err == nil && a.emit != nil {
			a.emit("compile:done", map[string]interface{}{
				"type":               "svg",
				"data":               res.URI,
				"usedGoTextFallback": res.UsedGoTextFallback,
			})
		}
	})
}

func (a *App) maybeAutoCompile() {
	a.mu.Lock()
	auto := a.autoCompile
	a.mu.Unlock()
	if auto {
		a.scheduleAutoCompile()
	}
}

// CompileSVGResult is returned when compiling page SVG for preview or export.
type CompileSVGResult struct {
	URI                string
	UsedGoTextFallback bool
}

// CompilePageSVG returns a data URI for the current page SVG.
func (a *App) CompilePageSVG(pageName string) (CompileSVGResult, error) {
	a.mu.Lock()
	doc := a.doc
	store := a.store
	client := a.client
	if pageName == "" {
		pageName = a.activePage
	}
	a.mu.Unlock()

	ctx := context.Background()
	if store != nil {
		_ = store.EnsureDocument(ctx, client, doc)
	}
	pageRes, err := compilesvg.CompilePageSVG(ctx, doc, pageName, client, store)
	if err != nil {
		return CompileSVGResult{}, err
	}
	return CompileSVGResult{
		URI:                encodePreview(pageRes.Data, "image/svg+xml"),
		UsedGoTextFallback: pageRes.UsedGoTextFallback,
	}, nil
}

// CompileDocumentHTML returns data URI for full HTML (all pages).
func (a *App) CompileDocumentHTML() (string, error) {
	a.mu.Lock()
	doc := a.doc
	store := a.store
	client := a.client
	filePath := a.filePath
	a.mu.Unlock()

	ctx := context.Background()
	if store != nil {
		_ = store.EnsureDocument(ctx, client, doc)
	}
	html, err := compilehtml.CompileWithCacheFrom(ctx, doc, filePath, nil, nil, client, store, compileopts.Options{})
	if err != nil {
		return "", err
	}
	return encodePreview(html, "text/html"), nil
}

// CompilePageHTML writes an intermediate single-page .psrt under the cache preview/ dir
// and returns a data URI for HTML compiled from that page only.
func (a *App) CompilePageHTML(pageName string) (string, error) {
	a.mu.Lock()
	doc := a.doc
	store := a.store
	client := a.client
	if pageName == "" {
		pageName = a.activePage
	}
	a.mu.Unlock()

	if pageName == "" {
		return "", fmt.Errorf("no active page")
	}

	pageDoc, err := extractPageDocument(doc, pageName)
	if err != nil {
		return "", err
	}

	ctx := context.Background()
	if store != nil {
		if _, err := writeIntermediatePagePSRT(store.RootDir, pageName, pageDoc); err != nil {
			return "", fmt.Errorf("write intermediate psrt: %w", err)
		}
		_ = store.EnsureDocument(ctx, client, pageDoc)
	}
	html, err := compilehtml.CompileWithCache(ctx, pageDoc, client, store)
	if err != nil {
		return "", err
	}
	return encodePreview(html, "text/html"), nil
}

// ExportSVG writes all page SVGs to dir.
func (a *App) ExportSVG(dir string) (CompileSVGResult, error) {
	a.mu.Lock()
	doc := a.doc
	store := a.store
	client := a.client
	a.mu.Unlock()
	batch, err := compilesvg.CompileWithCache(context.Background(), doc, client, dir, store)
	if err != nil {
		return CompileSVGResult{}, err
	}
	return CompileSVGResult{UsedGoTextFallback: batch.UsedGoTextFallback}, nil
}

// ExportHTML writes HTML next to psrt file.
func (a *App) ExportHTML() (string, error) {
	a.mu.Lock()
	path := a.filePath
	doc := a.doc
	store := a.store
	client := a.client
	a.mu.Unlock()
	if path == "" {
		return "", errNoFile()
	}
	base := filepath.Base(path)
	out := filepath.Join(filepath.Dir(path), base[:len(base)-len(filepath.Ext(base))]+".html")
	html, err := compilehtml.CompileWithCacheFrom(context.Background(), doc, path, nil, nil, client, store, compileopts.Options{})
	if err != nil {
		return "", err
	}
	if err := os.WriteFile(out, html, 0o644); err != nil {
		return "", err
	}
	return out, nil
}

func errNoFile() error {
	return &os.PathError{Op: "save", Path: "", Err: os.ErrInvalid}
}

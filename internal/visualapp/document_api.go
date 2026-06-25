package visualapp

import (
	"context"
	"fmt"

	"github.com/Dcrispim/psrt.core/compileasset/cache"
	"github.com/Dcrispim/psrt.core/compilehtml"
	"github.com/Dcrispim/psrt.core/compileopts"
	"github.com/Dcrispim/psrt.core/compilesvg"
	"github.com/Dcrispim/psrt.core/psrt"
	"github.com/Dcrispim/psrt.core/psrt/editor"
)

// GetDocumentJSON returns the in-memory document as JSON.
func (a *App) GetDocumentJSON() (string, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	return marshalDocument(&a.doc)
}

// SaveDocumentJSON parses doc JSON, updates memory, and writes the PSRT file.
func (a *App) SaveDocumentJSON(docJSON string) error {
	a.mu.Lock()
	path := a.filePath
	a.mu.Unlock()
	if path == "" {
		return fmt.Errorf("no file open")
	}
	return a.SaveDocumentJSONTo(docJSON, path)
}

// SaveDocumentJSONTo saves doc JSON to path and updates editor cache.
func (a *App) SaveDocumentJSONTo(docJSON, path string) error {
	doc, err := parseDocumentJSON(docJSON)
	if err != nil {
		return err
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	a.doc = doc
	a.filePath = path
	store, err := cache.NewStore("", path)
	if err != nil {
		return err
	}
	a.store = store
	return a.saveTo(path)
}

// ParseDocumentPSRT parses PSRT text and returns document JSON (does not save).
func (a *App) ParseDocumentPSRT(text string) (string, error) {
	doc, err := psrt.ParseString(text)
	if err != nil {
		return "", fmt.Errorf("parse: %w", err)
	}
	return marshalDocument(&doc)
}

// FormatDocumentJSON formats document JSON as PSRT text.
func (a *App) FormatDocumentJSON(docJSON string) (string, error) {
	return FormatDocumentJSON(docJSON)
}

// FormatPageDocumentJSON formats one page plus document fonts and constants as PSRT.
func FormatPageDocumentJSON(docJSON, pageName string) (string, error) {
	doc, err := parseDocumentJSON(docJSON)
	if err != nil {
		return "", err
	}
	pageDoc, err := extractPageDocument(doc, pageName)
	if err != nil {
		return "", err
	}
	data, err := editor.FormatDocument(&pageDoc)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// MergePageDocumentPSRT applies a page PSRT fragment (page, fonts, consts) into a full document JSON.
func MergePageDocumentPSRT(fullDocJSON, pageName, psrtText string) (string, error) {
	full, err := parseDocumentJSON(fullDocJSON)
	if err != nil {
		return "", err
	}
	parsed, err := psrt.ParseString(psrtText)
	if err != nil {
		return "", fmt.Errorf("parse: %w", err)
	}
	if len(parsed.Pages) == 0 {
		return "", fmt.Errorf("PSRT has no $START page block")
	}
	parsedPage := &parsed.Pages[0]
	for i := range parsed.Pages {
		if parsed.Pages[i].Name == pageName {
			parsedPage = &parsed.Pages[i]
			break
		}
	}
	idx := -1
	for i := range full.Pages {
		if full.Pages[i].Name == pageName {
			idx = i
			break
		}
	}
	if idx < 0 {
		return "", fmt.Errorf("page %q not found", pageName)
	}
	full.Pages[idx] = *parsedPage
	if len(parsed.Fonts) > 0 {
		full.Fonts = append([]string(nil), parsed.Fonts...)
	}
	if len(parsed.Consts) > 0 {
		full.Consts = cloneConsts(parsed.Consts)
	}
	return marshalDocument(&full)
}

// FormatDocumentJSON formats document JSON as PSRT text.
func FormatDocumentJSON(docJSON string) (string, error) {
	doc, err := parseDocumentJSON(docJSON)
	if err != nil {
		return "", err
	}
	data, err := editor.FormatDocument(&doc)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (a *App) FormatPageDocumentJSON(docJSON, pageName string) (string, error) {
	return FormatPageDocumentJSON(docJSON, pageName)
}

func (a *App) MergePageDocumentPSRT(fullDocJSON, pageName, psrtText string) (string, error) {
	return MergePageDocumentPSRT(fullDocJSON, pageName, psrtText)
}

// CompilePageSVGFromDocument compiles one page SVG from document JSON (editor state).
func (a *App) CompilePageSVGFromDocument(docJSON, pageName string) (CompileSVGResult, error) {
	doc, err := parseDocumentJSON(docJSON)
	if err != nil {
		return CompileSVGResult{}, err
	}
	a.mu.Lock()
	store := a.store
	client := a.client
	path := a.filePath
	if pageName == "" {
		pageName = a.activePage
	}
	a.mu.Unlock()
	if pageName == "" {
		return CompileSVGResult{}, fmt.Errorf("no page")
	}
	if store == nil && path != "" {
		store, _ = openStoreForPath(path)
	}
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

// CompilePageHTMLFromDocument compiles one page HTML from document JSON (editor state).
func (a *App) CompilePageHTMLFromDocument(docJSON, pageName string) (string, error) {
	doc, err := parseDocumentJSON(docJSON)
	if err != nil {
		return "", err
	}
	a.mu.Lock()
	store := a.store
	client := a.client
	path := a.filePath
	if pageName == "" {
		pageName = a.activePage
	}
	a.mu.Unlock()
	if pageName == "" {
		return "", fmt.Errorf("no page")
	}
	pageDoc, err := extractPageDocument(doc, pageName)
	if err != nil {
		return "", err
	}
	if store == nil && path != "" {
		store, _ = openStoreForPath(path)
	}
	ctx := context.Background()
	if store != nil {
		if _, err := writeIntermediatePagePSRT(store.RootDir, pageName, pageDoc); err != nil {
			return "", fmt.Errorf("write intermediate psrt: %w", err)
		}
		_ = store.EnsureDocument(ctx, client, pageDoc)
	}
	html, err := compilehtml.CompileWithCacheFrom(ctx, pageDoc, path, nil, nil, client, store, compileopts.Options{})
	if err != nil {
		return "", err
	}
	return encodePreview(html, "text/html"), nil
}

func parseDocumentJSON(docJSON string) (psrt.Document, error) {
	if docJSON == "" {
		return psrt.Document{}, fmt.Errorf("empty document")
	}
	return psrt.ParseJSON([]byte(docJSON))
}

func marshalDocument(doc *psrt.Document) (string, error) {
	b, err := psrt.ToJSON(*doc)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func openStoreForPath(path string) (*cache.Store, error) {
	return cache.NewStore("", path)
}

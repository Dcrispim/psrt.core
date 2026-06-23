package compilesvg

import (
	"context"
	"fmt"
	"net/http"

	"psrt/compileasset"
	"psrt/compileasset/cache"
	"psrt/compileopts"
	"psrt/psrt"
)

// CompilePageSVG returns SVG bytes for a single page by name.
func CompilePageSVG(ctx context.Context, doc psrt.Document, pageName string, client *http.Client, store *cache.Store) (PageSVGResult, error) {
	return CompilePageSVGWithOptions(ctx, doc, pageName, client, store, compileopts.Options{})
}

// CompilePageSVGWithOptions compiles one page to SVG with compile flags.
func CompilePageSVGWithOptions(ctx context.Context, doc psrt.Document, pageName string, client *http.Client, store *cache.Store, opts compileopts.Options) (PageSVGResult, error) {
	psrt.CleanEmptyTextBlockStyles(&doc)
	resolved, err := ResolveDocumentStrict(doc)
	if err != nil {
		return PageSVGResult{}, err
	}
	var p *psrt.Page
	for i := range resolved.Pages {
		if resolved.Pages[i].Name == pageName {
			p = &resolved.Pages[i]
			break
		}
	}
	if p == nil {
		return PageSVGResult{}, fmt.Errorf("page %q not found", pageName)
	}
	urls := compileasset.CollectAssetURLs(psrt.Document{Pages: []psrt.Page{*p}, Fonts: resolved.Fonts, Consts: resolved.Consts})
	pageURLs, fontURLs := compileasset.PartitionAssetURLs(resolved.Fonts, urls)
	var assets map[string]compileasset.Asset
	if store != nil {
		assets, err = cache.FetchDocumentURLsWithCache(ctx, client, store, resolved, pageURLs)
	} else {
		assets, err = compileasset.FetchURLs(client, pageURLs)
	}
	if err != nil {
		return PageSVGResult{}, err
	}
	if len(fontURLs) > 0 {
		fontAssets, ferr := compileasset.FetchFontAssets(ctx, client, fontURLs)
		if ferr != nil {
			return PageSVGResult{}, ferr
		}
		compileasset.MergeFontAssets(assets, fontAssets)
	}
	slug := UniqueSlugs([]string{p.Name})[0]
	return RenderPageSVGWithContext(ctx, p, slug, resolved.Fonts, assets, opts)
}

package compilesvg

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Dcrispim/psrt.core/compileasset"
	"github.com/Dcrispim/psrt.core/compileasset/cache"
	"github.com/Dcrispim/psrt.core/compileopts"
	"github.com/Dcrispim/psrt.core/compilesvg/textoutline"
	"github.com/Dcrispim/psrt.core/psrt"
)

// CompileResult holds batch SVG compile output metadata.
type CompileResult struct {
	UsedGoTextFallback bool
}

// Compile resolves constants, fetches assets, and writes one SVG per page into outDir.
func Compile(doc psrt.Document, client *http.Client, outDir string) (CompileResult, error) {
	return CompileWithOptions(context.Background(), doc, client, outDir, nil, compileopts.Options{})
}

// CompileWithOptions compiles all pages to SVG files with shared compile flags.
func CompileWithOptions(ctx context.Context, doc psrt.Document, client *http.Client, outDir string, store *cache.Store, opts compileopts.Options) (CompileResult, error) {
	return compileWithCache(ctx, doc, client, outDir, store, opts)
}

// CompileWithCache uses local asset cache when store is non-nil.
func CompileWithCache(ctx context.Context, doc psrt.Document, client *http.Client, outDir string, store *cache.Store) (CompileResult, error) {
	return compileWithCache(ctx, doc, client, outDir, store, compileopts.Options{})
}

func compileWithCache(ctx context.Context, doc psrt.Document, client *http.Client, outDir string, store *cache.Store, opts compileopts.Options) (CompileResult, error) {
	defer textoutline.CloseBrowser()

	psrt.CleanEmptyTextBlockStyles(&doc)
	resolved, err := ResolveDocumentStrict(doc)
	if err != nil {
		return CompileResult{}, err
	}
	urls := compileasset.CollectAssetURLs(resolved)
	pageURLs, fontURLs := compileasset.PartitionAssetURLs(resolved.Fonts, urls)
	var assets map[string]compileasset.Asset
	if store != nil {
		assets, err = cache.FetchDocumentURLsWithCache(ctx, client, store, resolved, pageURLs)
	} else {
		assets, err = compileasset.FetchURLs(client, pageURLs)
	}
	if err != nil {
		return CompileResult{}, err
	}
	if len(fontURLs) > 0 {
		fontAssets, ferr := compileasset.FetchFontAssets(ctx, client, fontURLs)
		if ferr != nil {
			return CompileResult{}, ferr
		}
		compileasset.MergeFontAssets(assets, fontAssets)
	}
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return CompileResult{}, fmt.Errorf("mkdir %q: %w", outDir, err)
	}

	names := make([]string, len(resolved.Pages))
	for i := range resolved.Pages {
		names[i] = resolved.Pages[i].Name
	}
	slugs := UniqueSlugs(names)

	var usedGoText bool
	for i := range resolved.Pages {
		p := &resolved.Pages[i]
		pageRes, err := RenderPageSVGWithContext(ctx, p, slugs[i], resolved.Fonts, assets, opts)
		if err != nil {
			return CompileResult{}, fmt.Errorf("page %q: %w", p.Name, err)
		}
		if pageRes.UsedGoTextFallback {
			usedGoText = true
		}
		path := filepath.Join(outDir, slugs[i]+".svg")
		if err := os.WriteFile(path, pageRes.Data, 0o644); err != nil {
			return CompileResult{}, fmt.Errorf("write %q: %w", path, err)
		}
	}
	return CompileResult{UsedGoTextFallback: usedGoText}, nil
}

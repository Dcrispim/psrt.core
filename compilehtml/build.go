package compilehtml

import (
	"context"
	"net/http"
	"path/filepath"
	"strings"

	"psrt/compileasset"
	"psrt/compileasset/cache"
	"psrt/compileopts"
	"psrt/compilesvg"
	"psrt/psrt"
)

// Compile downloads all assets referenced by doc and returns a self-contained HTML document.
func Compile(doc psrt.Document, client *http.Client) ([]byte, error) {
	return CompileWithOptions(context.Background(), doc, client, nil, compileopts.Options{})
}

// CompileWithOptions compiles with shared compile flags.
func CompileWithOptions(ctx context.Context, doc psrt.Document, client *http.Client, store *cache.Store, opts compileopts.Options) ([]byte, error) {
	return CompileWithCacheFrom(ctx, doc, "", nil, nil, client, store, opts)
}

// CompileWithCache uses local asset cache when store is non-nil.
func CompileWithCache(ctx context.Context, doc psrt.Document, client *http.Client, store *cache.Store) ([]byte, error) {
	return CompileWithCacheFrom(ctx, doc, "", nil, nil, client, store, compileopts.Options{})
}

// CompileWithCacheFrom is like CompileWithCache. Variants come only from sourcePath plus morePaths
// (explicit compile inputs), never from scanning the source directory.
func CompileWithCacheFrom(ctx context.Context, doc psrt.Document, sourcePath string, morePaths []string, morePSRT []VariantPSRT, client *http.Client, store *cache.Store, opts compileopts.Options) ([]byte, error) {
	psrt.CleanEmptyTextBlockStyles(&doc)
	resolved := compilesvg.ResolveDocument(doc)
	variants, err := buildVariants(sourcePath, resolved, morePaths, morePSRT)
	if err != nil {
		return nil, err
	}
	for i := range variants {
		psrt.CleanEmptyTextBlockStyles(&variants[i].Doc)
	}
	urls := collectAssetURLs(variants)
	pageURLs, fontURLs := compileasset.PartitionAssetURLs(resolved.Fonts, urls)
	var assets map[string]compileasset.Asset
	if store != nil {
		assets, err = cache.FetchDocumentURLsWithCache(ctx, client, store, resolved, pageURLs)
	} else {
		assets, err = compileasset.FetchURLs(client, pageURLs)
	}
	if err != nil {
		return nil, err
	}
	if len(fontURLs) > 0 {
		fontAssets, ferr := compileasset.FetchFontAssets(ctx, client, fontURLs)
		if ferr != nil {
			return nil, ferr
		}
		compileasset.MergeFontAssets(assets, fontAssets)
	}
	return RenderHTMLBundle(variants, assets, opts)
}

func buildVariants(sourcePath string, primary psrt.Document, morePaths []string, morePSRT []VariantPSRT) ([]Variant, error) {
	var extra []string
	for _, p := range morePaths {
		p = strings.TrimSpace(p)
		if p == "" || p == "-" {
			continue
		}
		extra = append(extra, p)
	}
	var variants []Variant
	if len(extra) == 0 && len(morePSRT) == 0 {
		label := "PSRT"
		if p := strings.TrimSpace(sourcePath); p != "" && p != "-" {
			label = filepath.Base(p)
		}
		return []Variant{{Label: label, Doc: primary}}, nil
	}
	if len(extra) > 0 {
		paths := make([]string, 0, 1+len(extra))
		if p := strings.TrimSpace(sourcePath); p != "" && p != "-" {
			paths = append(paths, p)
		}
		paths = append(paths, extra...)
		loaded, err := LoadVariantsFromPaths(paths)
		if err != nil {
			return nil, err
		}
		if len(loaded) == 0 {
			variants = []Variant{{Label: "PSRT", Doc: primary}}
		} else {
			loaded[0].Doc = primary
			variants = loaded
		}
	} else {
		label := "PSRT"
		if p := strings.TrimSpace(sourcePath); p != "" && p != "-" {
			label = filepath.Base(p)
		}
		variants = []Variant{{Label: label, Doc: primary}}
	}
	if len(morePSRT) > 0 {
		loaded, err := LoadVariantsFromPSRT(morePSRT)
		if err != nil {
			return nil, err
		}
		variants = append(variants, loaded...)
	}
	return variants, nil
}

func collectAssetURLs(variants []Variant) []string {
	seen := make(map[string]struct{})
	var out []string
	add := func(urls []string) {
		for _, u := range urls {
			if _, ok := seen[u]; ok {
				continue
			}
			seen[u] = struct{}{}
			out = append(out, u)
		}
	}
	for i := range variants {
		add(compileasset.CollectAssetURLs(variants[i].Doc))
	}
	return out
}

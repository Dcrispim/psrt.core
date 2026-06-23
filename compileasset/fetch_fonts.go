package compileasset

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

// PartitionAssetURLs splits collected URLs into page/image refs vs document font entries.
func PartitionAssetURLs(fonts, urls []string) (pageURLs, fontURLs []string) {
	fontSet := make(map[string]struct{}, len(fonts))
	for _, f := range fonts {
		if s := strings.TrimSpace(f); s != "" {
			fontSet[s] = struct{}{}
		}
	}
	for _, u := range urls {
		u = strings.TrimSpace(u)
		if u == "" {
			continue
		}
		if _, ok := fontSet[u]; ok {
			fontURLs = append(fontURLs, u)
		} else {
			pageURLs = append(pageURLs, u)
		}
	}
	return pageURLs, fontURLs
}

// FetchFontAssets downloads font files; map keys are the original $FONTS URLs from the document.
func FetchFontAssets(ctx context.Context, client *http.Client, fontURLs []string) (map[string]Asset, error) {
	out := make(map[string]Asset, len(fontURLs))
	for _, fontURL := range fontURLs {
		original := strings.TrimSpace(fontURL)
		if original == "" {
			continue
		}
		fetchURL, err := FontFetchURL(ctx, client, original)
		if err != nil {
			return nil, fmt.Errorf("font %q: %w", original, err)
		}
		fetched, err := FetchURLs(client, []string{fetchURL})
		if err != nil {
			return nil, fmt.Errorf("font %q: %w", original, err)
		}
		asset, ok := fetched[fetchURL]
		if !ok {
			return nil, fmt.Errorf("font %q: empty response", original)
		}
		out[original] = asset
	}
	return out, nil
}

// MergeFontAssets copies font assets into the main assets map (document font URL keys).
func MergeFontAssets(assets map[string]Asset, fonts map[string]Asset) {
	for k, v := range fonts {
		assets[k] = v
	}
}

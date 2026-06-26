package cache

import (
	"context"
	"net/http"

	"github.com/Dcrispim/psrt.core/compileasset"
)

// FetchURLsWithCache loads assets from disk when mapped; otherwise downloads and caches.
func FetchURLsWithCache(ctx context.Context, client *http.Client, store *Store, urls []string) (map[string]compileasset.Asset, error) {
	if store == nil {
		return compileasset.FetchURLs(client, urls)
	}
	out := make(map[string]compileasset.Asset, len(urls))
	labels := make(map[string]string)
	if store != nil {
		// page labels unknown here; use hash-only label
		for _, u := range urls {
			labels[u] = "asset"
		}
	}
	for _, raw := range urls {
		if asset, ok, err := store.ReadAsset(raw); err != nil {
			return nil, err
		} else if ok {
			out[raw] = asset
			continue
		}
		label := labels[raw]
		if err := store.EnsureCached(ctx, client, raw, label); err != nil {
			return nil, err
		}
		asset, ok, err := store.ReadAsset(raw)
		if err != nil {
			return nil, err
		}
		if !ok {
			fetched, err := compileasset.FetchURLs(client, []string{raw})
			if err != nil {
				return nil, err
			}
			out[raw] = fetched[raw]
			continue
		}
		out[raw] = asset
	}
	return out, nil
}

package cache

import (
	"context"
	"net/http"

	"github.com/Dcrispim/psrt.core/compileasset"
	"github.com/Dcrispim/psrt.core/psrt"
)

// FetchDocumentURLsWithCache uses page labels from doc for any missing cache entries.
func FetchDocumentURLsWithCache(ctx context.Context, client *http.Client, store *Store, doc psrt.Document, urls []string) (map[string]compileasset.Asset, error) {
	if store == nil {
		return compileasset.FetchURLs(client, urls)
	}
	labels := URLPageLabels(doc)
	out := make(map[string]compileasset.Asset, len(urls))
	for _, raw := range urls {
		if asset, ok, err := store.ReadAsset(raw); err != nil {
			return nil, err
		} else if ok {
			out[raw] = asset
			continue
		}
		label := labels[raw]
		if label == "" {
			label = "asset"
		}
		if err := store.EnsureCached(ctx, client, raw, label); err != nil {
			return nil, err
		}
		asset, ok, err := store.ReadAsset(raw)
		if err != nil {
			return nil, err
		}
		if !ok {
			if compileasset.IsLocalAssetRef(raw) {
				baseDir := ""
				if store != nil {
					baseDir = store.PsrtDir
				}
				path, err := compileasset.ResolveAssetPathRelative(raw, baseDir)
				if err != nil {
					return nil, err
				}
				asset, err := compileasset.ReadAssetFile(path)
				if err != nil {
					return nil, err
				}
				out[raw] = asset
				continue
			}
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

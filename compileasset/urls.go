package compileasset

import (
	"strings"

	"github.com/Dcrispim/psrt.core/psrt"
)

// LooksLikeHTTPURL reports whether raw is an http(s) URL.
func LooksLikeHTTPURL(raw string) bool {
	u := strings.TrimSpace(strings.ToLower(raw))
	return strings.HasPrefix(u, "http://") || strings.HasPrefix(u, "https://")
}

// CollectAssetURLs gathers unique URLs in deterministic order:
// pages in order (image then text ImageRefs), then document fonts.
func CollectAssetURLs(doc psrt.Document) []string {
	seen := make(map[string]struct{})
	var out []string
	add := func(u string) {
		u = ResolveAssetReference(u, doc.Consts)
		if u == "" || !IsAssetReference(u) {
			return
		}
		if _, dup := seen[u]; dup {
			return
		}
		seen[u] = struct{}{}
		out = append(out, u)
	}
	for i := range doc.Pages {
		add(doc.Pages[i].ImageURL)
		for j := range doc.Pages[i].Texts {
			add(doc.Pages[i].Texts[j].ImageRef)
		}
		for j := range doc.Pages[i].Masks {
			add(doc.Pages[i].Masks[j].ImageRef)
		}
	}
	for _, f := range doc.Fonts {
		add(f)
	}
	return out
}

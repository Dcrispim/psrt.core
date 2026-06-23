package compileasset

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

var googleFontFileURL = regexp.MustCompile(`url\((https://fonts\.gstatic\.com/[^)]+\.woff2[^)]*)\)`)

// IsGoogleFontsCSSURL reports Google Fonts CSS API links (not direct .woff2 files).
func IsGoogleFontsCSSURL(raw string) bool {
	u, err := url.Parse(strings.TrimSpace(raw))
	if err != nil {
		return false
	}
	host := strings.ToLower(strings.TrimPrefix(u.Hostname(), "www."))
	if host != "fonts.googleapis.com" && host != "fonts.google.com" {
		return false
	}
	return strings.Contains(strings.ToLower(u.Path), "/css")
}

// ResolveGoogleFontsWoff2URLs fetches the CSS API response and returns woff2 file URLs to download.
func ResolveGoogleFontsWoff2URLs(ctx context.Context, client *http.Client, cssURL string) ([]string, error) {
	cssURL = strings.TrimSpace(cssURL)
	if !IsGoogleFontsCSSURL(cssURL) {
		return nil, fmt.Errorf("not a Google Fonts CSS URL")
	}
	if client == nil {
		client = http.DefaultClient
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, cssURL, nil)
	if err != nil {
		return nil, err
	}
	// Google serves @font-face rules with gstatic URLs only for browser-like clients.
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch Google Fonts CSS: %w", err)
	}
	body, err := io.ReadAll(resp.Body)
	_ = resp.Body.Close()
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("fetch Google Fonts CSS: status %s", resp.Status)
	}

	seen := make(map[string]struct{})
	var out []string
	for _, m := range googleFontFileURL.FindAllStringSubmatch(string(body), -1) {
		if len(m) < 2 {
			continue
		}
		u := strings.TrimSpace(m[1])
		if u == "" {
			continue
		}
		if _, dup := seen[u]; dup {
			continue
		}
		seen[u] = struct{}{}
		out = append(out, u)
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("no woff2 URLs in Google Fonts CSS %q", cssURL)
	}
	return out, nil
}

// FontFetchURL returns the URL to download for a document font entry (resolves Google CSS to woff2).
func FontFetchURL(ctx context.Context, client *http.Client, fontURL string) (string, error) {
	fontURL = strings.TrimSpace(fontURL)
	if !IsGoogleFontsCSSURL(fontURL) {
		return fontURL, nil
	}
	urls, err := ResolveGoogleFontsWoff2URLs(ctx, client, fontURL)
	if err != nil {
		return "", err
	}
	return urls[0], nil
}

// ExpandFontURLsForFetch replaces Google Fonts CSS links with downloadable woff2 URLs.
// assetKeys maps each fetch URL to the original document font URL (for assets map lookup).
func ExpandFontURLsForFetch(ctx context.Context, client *http.Client, fontURLs []string) (fetchURLs []string, assetKeys map[string]string, err error) {
	assetKeys = make(map[string]string)
	seenFetch := make(map[string]struct{})
	add := func(fetch, original string) {
		if fetch == "" {
			return
		}
		assetKeys[fetch] = original
		if _, ok := seenFetch[fetch]; ok {
			return
		}
		seenFetch[fetch] = struct{}{}
		fetchURLs = append(fetchURLs, fetch)
	}
	for _, raw := range fontURLs {
		raw = strings.TrimSpace(raw)
		if raw == "" {
			continue
		}
		if IsGoogleFontsCSSURL(raw) {
			woffs, err := ResolveGoogleFontsWoff2URLs(ctx, client, raw)
			if err != nil {
				return nil, nil, fmt.Errorf("font %q: %w", raw, err)
			}
			for _, w := range woffs {
				add(w, raw)
			}
			continue
		}
		add(raw, raw)
	}
	return fetchURLs, assetKeys, nil
}

// ParseGoogleFontFamilies extracts family names from a Google Fonts CSS API URL.
func ParseGoogleFontFamilies(cssURL string) []string {
	u, err := url.Parse(strings.TrimSpace(cssURL))
	if err != nil {
		return nil
	}
	seen := make(map[string]struct{})
	var out []string
	for _, param := range u.Query()["family"] {
		segments := strings.Split(param, "|")
		if len(segments) == 1 {
			segments = []string{param}
		}
		for _, seg := range segments {
			raw := strings.TrimSpace(strings.Split(seg, ":")[0])
			if raw == "" {
				continue
			}
			name := decodeGoogleFamilyName(raw)
			if name == "" {
				continue
			}
			if _, dup := seen[name]; dup {
				continue
			}
			seen[name] = struct{}{}
			out = append(out, name)
		}
	}
	return out
}

func decodeGoogleFamilyName(encoded string) string {
	dec, err := url.QueryUnescape(strings.ReplaceAll(encoded, "+", " "))
	if err != nil {
		return strings.TrimSpace(strings.ReplaceAll(encoded, "+", " "))
	}
	return strings.TrimSpace(dec)
}

// FontFamilyNameForURL returns the CSS font-family name for a $FONTS entry.
func FontFamilyNameForURL(fontURL string, index int) string {
	if names := ParseGoogleFontFamilies(fontURL); len(names) > 0 {
		return names[0]
	}
	return fmt.Sprintf("CompiledFont_%d", index)
}

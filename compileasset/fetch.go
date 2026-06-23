package compileasset

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// FetchURLs downloads each URL via GET with client. Fail-fast on first non-success.
func FetchURLs(client *http.Client, urls []string) (map[string]Asset, error) {
	out := make(map[string]Asset, len(urls))
	for _, raw := range urls {
		raw = strings.TrimSpace(raw)
		if raw == "" {
			continue
		}
		req, err := http.NewRequest(http.MethodGet, raw, nil)
		if err != nil {
			return nil, fmt.Errorf("url %q: %w", raw, err)
		}
		resp, err := client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("fetch %q: %w", raw, err)
		}
		body, err := io.ReadAll(resp.Body)
		_ = resp.Body.Close()
		if err != nil {
			return nil, fmt.Errorf("read body %q: %w", raw, err)
		}
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			return nil, fmt.Errorf("fetch %q: status %s", raw, resp.Status)
		}
		mime := pickMIME(resp.Header.Get("Content-Type"), raw, body)
		out[raw] = Asset{Bytes: body, MIME: mime}
	}
	return out, nil
}

func pickMIME(ctHeader, rawURL string, body []byte) string {
	base := mimeFromHeader(ctHeader)
	if base != "" && base != "application/octet-stream" {
		return base
	}
	if m := mimeFromURLPath(rawURL); m != "" {
		return m
	}
	m := sniffMIME(body)
	if m != "" {
		return m
	}
	if base != "" {
		return base
	}
	return "application/octet-stream"
}

func mimeFromHeader(ct string) string {
	ct = strings.TrimSpace(ct)
	if ct == "" {
		return ""
	}
	if idx := strings.Index(ct, ";"); idx >= 0 {
		ct = strings.TrimSpace(ct[:idx])
	}
	return strings.ToLower(strings.TrimSpace(ct))
}

func sniffMIME(b []byte) string {
	sample := b
	const maxSniff = 512
	if len(sample) > maxSniff {
		sample = sample[:maxSniff]
	}
	switch {
	case len(sample) >= 3 && bytes.Equal(sample[:3], []byte{0xFF, 0xD8, 0xFF}):
		return "image/jpeg"
	case len(sample) >= 8 && bytes.Equal(sample[:8], []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}):
		return "image/png"
	case len(sample) >= 6 && string(sample[:6]) == `GIF87a`:
		return "image/gif"
	case len(sample) >= 6 && string(sample[:6]) == `GIF89a`:
		return "image/gif"
	case len(sample) >= 12 && string(sample[0:4]) == `RIFF` && string(sample[8:12]) == `WEBP`:
		return "image/webp"
	case len(sample) >= 4 && string(sample[:4]) == `wOFF`:
		return "font/woff"
	case len(sample) >= 4 && string(sample[:4]) == `wOF2`:
		return "font/woff2"
	case len(sample) >= 4 && bytes.Equal(sample[:4], []byte{0x00, 0x01, 0x00, 0x00}):
		return "font/ttf"
	case len(sample) >= 4 && string(sample[:4]) == `OTTO`:
		return "font/otf"
	case len(sample) >= 12 && string(sample[4:8]) == `ftyp`:
		brand := string(sample[8:12])
		if brand == "avif" || brand == "avis" || brand == "mif1" {
			return "image/avif"
		}
		return ""
	default:
		return ""
	}
}

package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"path/filepath"
	"strings"
	"time"
	"unicode"
)

func hash8(url string) string {
	sum := sha256.Sum256([]byte(url))
	return hex.EncodeToString(sum[:4])
}

func sanitizeLabel(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return "asset"
	}
	var b strings.Builder
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-' || r == '_' {
			b.WriteRune(r)
		} else if r == ' ' {
			b.WriteRune('_')
		}
	}
	out := b.String()
	if out == "" {
		return "asset"
	}
	return out
}

func extFromMIME(mime string) string {
	mime = strings.ToLower(strings.TrimSpace(mime))
	switch {
	case strings.Contains(mime, "jpeg"), mime == "image/jpg":
		return ".jpg"
	case strings.Contains(mime, "png"):
		return ".png"
	case strings.Contains(mime, "gif"):
		return ".gif"
	case strings.Contains(mime, "webp"):
		return ".webp"
	case strings.Contains(mime, "avif"):
		return ".avif"
	case strings.Contains(mime, "svg"):
		return ".svg"
	case strings.Contains(mime, "woff2"):
		return ".woff2"
	case strings.Contains(mime, "woff"):
		return ".woff"
	case strings.Contains(mime, "ttf"), strings.Contains(mime, "font-sfnt"):
		return ".ttf"
	default:
		return ".bin"
	}
}

func extFromURL(rawURL string) string {
	ext := filepath.Ext(rawURL)
	if idx := strings.Index(ext, "?"); idx >= 0 {
		ext = ext[:idx]
	}
	ext = strings.ToLower(ext)
	if len(ext) > 1 && len(ext) <= 6 {
		return ext
	}
	return ""
}

func assetFilename(pageLabel, psrtBase, url, mime string) string {
	ts := fmt.Sprintf("%s%06d", time.Now().Format("20060102150405"), time.Now().UnixNano()%1000000)
	ext := extFromURL(url)
	if ext == "" {
		ext = extFromMIME(mime)
	}
	if ext == "" {
		ext = ".bin"
	}
	return fmt.Sprintf("%s_%s_%s_%s%s",
		ts, sanitizeLabel(pageLabel), sanitizeLabel(psrtBase), hash8(url), ext)
}

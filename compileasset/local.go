package compileasset

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// IsAssetReference reports whether raw is an http(s) URL or a local file reference.
func IsAssetReference(raw string) bool {
	return LooksLikeHTTPURL(raw) || IsLocalAssetRef(raw)
}

// IsLocalAssetRef reports file:// URLs and filesystem paths (not http).
func IsLocalAssetRef(raw string) bool {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return false
	}
	if LooksLikeHTTPURL(raw) {
		return false
	}
	lower := strings.ToLower(raw)
	if strings.HasPrefix(lower, "file:") {
		return true
	}
	return strings.Contains(raw, `\`) || strings.Contains(raw, `/`) || hasWindowsDrive(raw)
}

func hasWindowsDrive(s string) bool {
	if len(s) < 2 || s[1] != ':' {
		return false
	}
	c := s[0]
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}

// ResolveAssetPath turns a PSRT asset reference into an absolute filesystem path.
func ResolveAssetPath(raw string) (string, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return "", fmt.Errorf("empty asset path")
	}
	if strings.HasPrefix(strings.ToLower(raw), "file:") {
		return pathFromFileURL(raw)
	}
	p := filepath.FromSlash(raw)
	if filepath.IsAbs(p) {
		return filepath.Clean(p), nil
	}
	return filepath.Clean(p), nil
}

// ResolveAssetPathRelative resolves raw against baseDir when raw is not absolute.
func ResolveAssetPathRelative(raw, baseDir string) (string, error) {
	p, err := ResolveAssetPath(raw)
	if err != nil {
		return "", err
	}
	if filepath.IsAbs(p) {
		return p, nil
	}
	if baseDir == "" {
		return "", fmt.Errorf("relative asset path %q requires document directory", raw)
	}
	return filepath.Clean(filepath.Join(baseDir, p)), nil
}

func pathFromFileURL(raw string) (string, error) {
	u, err := url.Parse(raw)
	if err != nil {
		return "", fmt.Errorf("parse file url: %w", err)
	}
	if u.Path == "" && u.Opaque != "" {
		u.Path = u.Opaque
	}
	p := u.Path
	if runtime.GOOS == "windows" {
		p = strings.TrimPrefix(p, "/")
		if strings.HasPrefix(p, "/") {
			p = strings.TrimPrefix(p, "/")
		}
	}
	p = filepath.FromSlash(p)
	if p == "" {
		return "", fmt.Errorf("empty path in file url %q", raw)
	}
	return filepath.Clean(p), nil
}

// ReadAssetFile reads a local file and detects MIME type.
func ReadAssetFile(path string) (Asset, error) {
	body, err := os.ReadFile(path)
	if err != nil {
		return Asset{}, err
	}
	mime := mimeFromFilesystemPath(path)
	if mime == "application/octet-stream" {
		if m := sniffMIME(body); m != "" {
			mime = m
		}
	}
	return Asset{Bytes: body, MIME: mime}, nil
}

func mimeFromFilesystemPath(path string) string {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".webp":
		return "image/webp"
	case ".avif":
		return "image/avif"
	case ".svg":
		return "image/svg+xml"
	case ".woff2":
		return "font/woff2"
	case ".woff":
		return "font/woff"
	case ".ttf":
		return "font/ttf"
	case ".otf":
		return "font/otf"
	default:
		return "application/octet-stream"
	}
}

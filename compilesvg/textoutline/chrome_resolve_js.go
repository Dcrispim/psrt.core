//go:build js

package textoutline

// AddChromeSearchRoots is a no-op in WASM builds.
func AddChromeSearchRoots(_ ...string) {}

// ResolveChromeExec always returns empty under GOOS=js (go-text fallback).
func ResolveChromeExec() string {
	return ""
}

// BundledChromeInDir always returns empty under GOOS=js.
func BundledChromeInDir(_ string) string {
	return ""
}

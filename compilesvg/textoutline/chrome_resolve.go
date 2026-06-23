//go:build !js

package textoutline

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var extraChromeSearchRoots []string

// AddChromeSearchRoots registers extra directories checked for bundled Chrome.
func AddChromeSearchRoots(roots ...string) {
	for _, r := range roots {
		r = strings.TrimSpace(r)
		if r != "" {
			extraChromeSearchRoots = append(extraChromeSearchRoots, r)
		}
	}
}

// ResolveChromeExec returns a Chromium/Edge executable path, or "" if none found.
func ResolveChromeExec() string {
	bases := chromeSearchBases()
	if p := resolveChromeCandidate(os.Getenv("CHROME_PATH"), bases); p != "" {
		return p
	}
	if p := resolveChromeCandidate(os.Getenv("EDGE_PATH"), bases); p != "" {
		return p
	}
	for _, dir := range bases {
		if p := BundledChromeInDir(dir); p != "" {
			return p
		}
	}
	if os.Getenv("PSRT_DEBUG_CHROME") != "" {
		log.Printf("psrt: chrome not found; searched: %v", bases)
	}
	return ""
}

// BundledChromeInDir returns a Chrome/Chromium executable next to dir, if any.
func BundledChromeInDir(dir string) string {
	if dir == "" {
		return ""
	}
	dir = filepath.Clean(dir)
	names := []string{
		"chrome.exe",
		"chromium.exe",
		"chrome-headless-shell.exe",
		"msedge.exe",
		"Google Chrome.exe",
		filepath.Join("browser", "chrome.exe"),
		filepath.Join("browser", "chrome-headless-shell.exe"),
		filepath.Join("chromium", "chrome.exe"),
		"chrome",
		"chromium",
		filepath.Join("browser", "chrome"),
	}
	for _, name := range names {
		if p := absChromePath(filepath.Join(dir, name)); p != "" {
			return p
		}
	}
	if p := findPuppeteerChrome(dir); p != "" {
		return p
	}
	return ""
}

func chromeSearchBases() []string {
	seen := make(map[string]struct{})
	var out []string
	add := func(dir string) {
		dir = strings.TrimSpace(dir)
		if dir == "" {
			return
		}
		abs, err := filepath.Abs(filepath.Clean(dir))
		if err != nil {
			return
		}
		if _, ok := seen[abs]; ok {
			return
		}
		seen[abs] = struct{}{}
		out = append(out, abs)
	}

	for _, r := range extraChromeSearchRoots {
		add(r)
	}

	var seeds []string
	if d := executableDir(); d != "" {
		seeds = append(seeds, d)
	}
	if wd, err := os.Getwd(); err == nil && wd != "" {
		seeds = append(seeds, wd)
	}
	if v := strings.TrimSpace(os.Getenv("PSRT_APP_DIR")); v != "" {
		seeds = append(seeds, v)
	}

	for _, seed := range seeds {
		dir := filepath.Clean(seed)
		for i := 0; i < 12; i++ {
			add(dir)
			add(filepath.Join(dir, "build", "bin"))
			add(filepath.Join(dir, "cmd", "psrt-gui", "build", "bin"))
			parent := filepath.Dir(dir)
			if parent == dir {
				break
			}
			dir = parent
		}
	}
	return out
}

func resolveChromeCandidate(raw string, bases []string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ""
	}
	candidates := []string{raw}
	if !filepath.IsAbs(raw) {
		for _, base := range bases {
			candidates = append(candidates, filepath.Join(base, raw))
		}
	}
	for _, c := range candidates {
		if p := absChromePath(c); p != "" {
			return p
		}
		if p := chromeExeInTree(c); p != "" {
			return p
		}
	}
	return ""
}

func chromeExeInTree(path string) string {
	path = filepath.Clean(path)
	if fileExists(path) {
		return absChromePath(path)
	}
	if !dirExists(path) {
		return ""
	}
	for _, name := range []string{"chrome.exe", "chrome", "Google Chrome.exe"} {
		if p := absChromePath(filepath.Join(path, name)); p != "" {
			return p
		}
	}
	if p := findPuppeteerChrome(path); p != "" {
		return p
	}
	return findChromeUnderPuppeteerDir(path)
}

func findPuppeteerChrome(rootDir string) string {
	patterns := []string{
		filepath.Join(rootDir, "chrome", "win64-*", "chrome-win64", "chrome.exe"),
		filepath.Join(rootDir, "chrome", "win32-*", "chrome-win32", "chrome.exe"),
		filepath.Join(rootDir, "win64-*", "chrome-win64", "chrome.exe"),
		filepath.Join(rootDir, "chrome-win64", "chrome.exe"),
	}
	for _, pat := range patterns {
		matches, err := filepath.Glob(pat)
		if err != nil {
			continue
		}
		for _, m := range matches {
			if p := absChromePath(m); p != "" {
				return p
			}
		}
	}
	return findChromeUnderPuppeteerDir(filepath.Join(rootDir, "chrome"))
}

func findChromeUnderPuppeteerDir(root string) string {
	if !dirExists(root) {
		return ""
	}
	var found string
	_ = filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil || found != "" {
			return nil
		}
		if d.IsDir() {
			return nil
		}
		base := strings.ToLower(d.Name())
		if base == "chrome.exe" || base == "chrome" {
			found = absChromePath(path)
			return fs.SkipAll
		}
		return nil
	})
	return found
}

func absChromePath(path string) string {
	if !fileExists(path) {
		return ""
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return path
	}
	return abs
}

func executableDir() string {
	exe, err := os.Executable()
	if err != nil {
		return ""
	}
	if resolved, err := filepath.EvalSymlinks(exe); err == nil {
		exe = resolved
	}
	return filepath.Dir(exe)
}

func dirExists(path string) bool {
	st, err := os.Stat(path)
	return err == nil && st.IsDir()
}

func fileExists(path string) bool {
	st, err := os.Stat(path)
	return err == nil && !st.IsDir()
}

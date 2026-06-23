//go:build !js

package textoutline

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestBundledChromeInDir_puppeteerLayout(t *testing.T) {
	dir := t.TempDir()
	chromeRoot := filepath.Join(dir, "chrome", "win64-1.0.0", "chrome-win64")
	if err := os.MkdirAll(chromeRoot, 0o755); err != nil {
		t.Fatal(err)
	}
	exe := filepath.Join(chromeRoot, "chrome.exe")
	if err := os.WriteFile(exe, []byte("stub"), 0o644); err != nil {
		t.Fatal(err)
	}
	got := BundledChromeInDir(dir)
	if got != exe {
		t.Fatalf("BundledChromeInDir() = %q, want %q", got, exe)
	}
}

func TestBundledChromeInDir_flatChromeExe(t *testing.T) {
	dir := t.TempDir()
	exe := filepath.Join(dir, "chrome.exe")
	if err := os.WriteFile(exe, []byte("stub"), 0o644); err != nil {
		t.Fatal(err)
	}
	if got := BundledChromeInDir(dir); got != exe {
		t.Fatalf("got %q", got)
	}
}

func TestChromeExeInTree_chromeWin64Directory(t *testing.T) {
	dir := t.TempDir()
	chromeRoot := filepath.Join(dir, "chrome", "win64-148.0.0", "chrome-win64")
	if err := os.MkdirAll(chromeRoot, 0o755); err != nil {
		t.Fatal(err)
	}
	exe := filepath.Join(chromeRoot, "chrome.exe")
	if err := os.WriteFile(exe, []byte("stub"), 0o644); err != nil {
		t.Fatal(err)
	}
	got := chromeExeInTree(chromeRoot)
	if got != exe {
		t.Fatalf("chromeExeInTree(dir) = %q, want %q", got, exe)
	}
}

func TestResolveChromeCandidate_relativePuppeteerDir(t *testing.T) {
	dir := t.TempDir()
	chromeRoot := filepath.Join(dir, "chrome", "win64-148.0.0", "chrome-win64")
	if err := os.MkdirAll(chromeRoot, 0o755); err != nil {
		t.Fatal(err)
	}
	exe := filepath.Join(chromeRoot, "chrome.exe")
	if err := os.WriteFile(exe, []byte("stub"), 0o644); err != nil {
		t.Fatal(err)
	}
	t.Setenv("CHROME_PATH", "")
	got := resolveChromeCandidate(`chrome\win64-148.0.0\chrome-win64`, []string{dir})
	if got != exe {
		t.Fatalf("resolveChromeCandidate = %q, want %q", got, exe)
	}
}

func TestResolveChromeExec_fromRepoRoot(t *testing.T) {
	root := filepath.Join("..", "..")
	abs, err := filepath.Abs(root)
	if err != nil {
		t.Fatal(err)
	}
	t.Chdir(abs)
	got := ResolveChromeExec()
	if got == "" {
		t.Fatalf("ResolveChromeExec() empty from repo root %s", abs)
	}
	if !strings.HasSuffix(strings.ToLower(got), "chrome.exe") {
		t.Fatalf("unexpected chrome path: %s", got)
	}
}

func TestFindPuppeteerChrome_directChromeWin64(t *testing.T) {
	dir := t.TempDir()
	chromeRoot := filepath.Join(dir, "chrome", "win64-9.9.9", "chrome-win64")
	if err := os.MkdirAll(chromeRoot, 0o755); err != nil {
		t.Fatal(err)
	}
	exe := filepath.Join(chromeRoot, "chrome.exe")
	if err := os.WriteFile(exe, []byte("stub"), 0o644); err != nil {
		t.Fatal(err)
	}
	got := findPuppeteerChrome(dir)
	if got != exe {
		t.Fatalf("findPuppeteerChrome = %q, want %q", got, exe)
	}
}

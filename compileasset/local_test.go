package compileasset

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"psrt/psrt"
)

func TestIsLocalAssetRef(t *testing.T) {
	if !IsLocalAssetRef(`file:///d%3A/projs/GO/psrt/grade_referencia.png`) {
		t.Fatal("file url should be local")
	}
	if !IsLocalAssetRef(`d:\projs\GO\psrt\grade.png`) {
		t.Fatal("windows path should be local")
	}
	if IsLocalAssetRef(`https://example.com/x.png`) {
		t.Fatal("http should not be local")
	}
}

func TestResolveAssetPath_fileURL(t *testing.T) {
	got, err := ResolveAssetPath(`file:///d%3A/projs/GO/psrt/grade_referencia.png`)
	if err != nil {
		t.Fatal(err)
	}
	if runtime.GOOS == "windows" {
		if !strings.EqualFold(got, `d:\projs\GO\psrt\grade_referencia.png`) {
			t.Fatalf("got %q", got)
		}
	} else {
		want := filepath.Clean("/d/projs/GO/psrt/grade_referencia.png")
		if got != want {
			t.Fatalf("got %q want %q", got, want)
		}
	}
}

func TestCollectAssetURLs_includesLocal(t *testing.T) {
	doc := psrt.Document{
		Pages: []psrt.Page{{
			ImageURL: `file:///tmp/x.png`,
		}},
	}
	urls := CollectAssetURLs(doc)
	if len(urls) != 1 {
		t.Fatalf("urls: %v", urls)
	}
	if !IsLocalAssetRef(urls[0]) {
		t.Fatalf("expected local ref, got %q", urls[0])
	}
}

func TestReadAssetFile_png(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "x.png")
	if err := os.WriteFile(path, []byte{0x89, 0x50, 0x4e, 0x47}, 0o644); err != nil {
		t.Fatal(err)
	}
	asset, err := ReadAssetFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if asset.MIME != "image/png" {
		t.Fatalf("mime %q", asset.MIME)
	}
}

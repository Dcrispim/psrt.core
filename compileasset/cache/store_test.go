package cache

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestStoreEnsureAndRefresh(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/png")
		_, _ = w.Write([]byte{0x89, 0x50, 0x4e, 0x47})
	}))
	defer srv.Close()

	dir := t.TempDir()
	psrtPath := filepath.Join(dir, "exemplo.psrt")
	if err := os.WriteFile(psrtPath, []byte("$START p\n$END p\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	store, err := NewStore(dir, psrtPath)
	if err != nil {
		t.Fatal(err)
	}
	url := srv.URL + "/img.png"
	if err := store.EnsureCached(context.Background(), srv.Client(), url, "capa"); err != nil {
		t.Fatal(err)
	}
	abs1, ok := store.ResolveAbsolute(url)
	if !ok {
		t.Fatal("expected cache hit")
	}
	if err := store.RefreshAsset(context.Background(), srv.Client(), url, "capa"); err != nil {
		t.Fatal(err)
	}
	abs2, ok := store.ResolveAbsolute(url)
	if !ok || abs1 == abs2 {
		t.Fatalf("refresh should point to new file: %s vs %s", abs1, abs2)
	}
}

func TestStoreEnsureLocalFile(t *testing.T) {
	dir := t.TempDir()
	imgPath := filepath.Join(dir, "grade.png")
	if err := os.WriteFile(imgPath, []byte{0x89, 0x50, 0x4e, 0x47}, 0o644); err != nil {
		t.Fatal(err)
	}
	psrtPath := filepath.Join(dir, "mood.psrt")
	if err := os.WriteFile(psrtPath, []byte("$START p\n$END p\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	store, err := NewStore(dir, psrtPath)
	if err != nil {
		t.Fatal(err)
	}
	ref := "file:///" + filepath.ToSlash(imgPath)
	if err := store.EnsureCached(context.Background(), nil, ref, "capa"); err != nil {
		t.Fatal(err)
	}
	asset, ok, err := store.ReadAsset(ref)
	if err != nil || !ok {
		t.Fatalf("read asset: ok=%v err=%v", ok, err)
	}
	if asset.MIME != "image/png" {
		t.Fatalf("mime %q", asset.MIME)
	}
}

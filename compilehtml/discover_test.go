package compilehtml

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDiscoverPSRTPaths_primaryFirst(t *testing.T) {
	dir := t.TempDir()
	primary := filepath.Join(dir, "main.psrt")
	other := filepath.Join(dir, "alt.psrt")
	for _, p := range []string{primary, other} {
		if err := os.WriteFile(p, []byte("$FONTS\n$ENDFONTS\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.WriteFile(filepath.Join(dir, "readme.txt"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	got, err := DiscoverPSRTPaths(primary)
	if err != nil {
		t.Fatal(err)
	}
	primaryAbs, _ := filepath.Abs(primary)
	if len(got) != 2 || got[0] != primaryAbs {
		t.Fatalf("got %v want primary %s first", got, primaryAbs)
	}
}

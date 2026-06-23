package webconnector

import (
	"os"
	"path/filepath"
	"testing"
)

func TestScanLibraryProjects(t *testing.T) {
	dir := t.TempDir()
	sub := filepath.Join(dir, "manga")
	if err := os.MkdirAll(sub, 0o755); err != nil {
		t.Fatal(err)
	}
	psrtPath := filepath.Join(sub, "cap1.psrt")
	content := `$START p1 | {} | https://x/a.png
>>1-2-3-4 | {} | 0
hi
$END p1
$START p2 | {} | https://x/b.png
>>1-2-3-4 | {} | 0
bye
$END p2
`
	if err := os.WriteFile(psrtPath, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	projects, err := scanLibraryProjects(dir, 4)
	if err != nil {
		t.Fatal(err)
	}
	if len(projects) != 1 {
		t.Fatalf("projects: got %d want 1", len(projects))
	}
	p := projects[0]
	if p.Path != "manga/cap1.psrt" {
		t.Fatalf("path: %q", p.Path)
	}
	if p.Title != "cap1" {
		t.Fatalf("title: %q", p.Title)
	}
	if p.PageCount != 2 {
		t.Fatalf("pageCount: %d", p.PageCount)
	}
}

func TestShouldSkipLibraryDir(t *testing.T) {
	if !shouldSkipLibraryDir("node_modules") {
		t.Fatal("expected skip")
	}
	if shouldSkipLibraryDir("pages") {
		t.Fatal("expected keep")
	}
}

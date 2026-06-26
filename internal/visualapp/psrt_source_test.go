package visualapp

import (
	"strings"
	"testing"

	"github.com/Dcrispim/psrt.core/psrt"
)

func TestGetSetDocumentPSRT(t *testing.T) {
	a := New(nil)
	a.doc = psrt.Document{
		Pages: []psrt.Page{{Name: "intro", ImageURL: "https://example.com/bg.jpg"}},
		Fonts: []string{"https://fonts.example/f.woff2"},
	}
	a.activePage = "intro"
	a.filePath = "/tmp/test.psrt"

	got, err := a.GetDocumentPSRT()
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(got, "intro") {
		t.Fatalf("missing page: %q", got)
	}

	if err := a.SetDocumentFromPSRT(got); err != nil {
		t.Fatal(err)
	}
	if len(a.doc.Pages) != 1 || a.activePage != "intro" {
		t.Fatalf("doc after set: pages=%d active=%q", len(a.doc.Pages), a.activePage)
	}
}

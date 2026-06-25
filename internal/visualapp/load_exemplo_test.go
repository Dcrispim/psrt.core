package visualapp

import (
	"testing"

	"github.com/Dcrispim/psrt.core/psrt/editor"
)

func TestLoadExemploOut(t *testing.T) {
	path := `d:\projs\GO\psrt\exemplo-out.psrt`
	doc, err := editor.LoadDocument(path)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("pages=%d", len(doc.Pages))
	a := New(nil)
	if err := a.OpenFile(path); err != nil {
		t.Fatal(err)
	}
	j, err := a.GetDocumentJSON()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("json len=%d", len(j))
	for _, p := range doc.Pages {
		if p.ImageURL == "" {
			continue
		}
		uri, err := a.GetAssetDataURI(p.ImageURL)
		if err != nil {
			t.Fatalf("asset %s: %v", p.Name, err)
		}
		t.Logf("asset %s uri len=%d", p.Name, len(uri))
	}
}

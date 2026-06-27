package psrt

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestLoadAssetMaskHeights(t *testing.T) {
	root := filepath.Join("..", "assets", "psrts", "o-home-que-nao-desistia-pt-br.psrt")
	data, err := os.ReadFile(root)
	if err != nil {
		t.Skip(err)
	}
	// This fixture predates the comma coordinate separator (hyphen-separated
	// headers); run it through the legacy converter before parsing, same as
	// any real legacy .psrt would need.
	converted, err := ConvertLegacyDocument(string(data))
	if err != nil {
		t.Fatal(err)
	}
	doc, err := Parse(bytes.NewReader([]byte(converted)))
	if err != nil {
		t.Fatal(err)
	}
	page := doc.Pages[0]
	for _, p := range doc.Pages {
		if p.Name == "pagina_1" {
			page = p
			break
		}
	}
	if len(page.Masks) < 5 {
		t.Fatalf("pagina_1 masks=%d want >=5", len(page.Masks))
	}
	for _, m := range page.Masks {
		if m.Height <= 0 {
			t.Fatalf("mask #%d has zero height: %+v", m.Index, m)
		}
	}
	// ==69.46-7.34-26.22-2.64 | ... | 4
	for _, m := range page.Masks {
		if m.Index == 4 && m.Height != 2.64 {
			t.Fatalf("mask #4 height=%v want 2.64", m.Height)
		}
	}
}

package compilehtml

import (
	"strings"
	"testing"

	"psrt/compileasset"
	"psrt/compileopts"
	"psrt/psrt"
)

func TestRenderHTMLBundle_variantSwitcher(t *testing.T) {
	docA := psrt.Document{
		Pages: []psrt.Page{{
			Name: "p", ImageURL: "u",
			Texts: []psrt.Text{psrt.NewText(10, 10, 50, 5, 0, psrt.Style(`{}`), "A", "")},
		}},
	}
	docB := psrt.Document{
		Pages: []psrt.Page{{
			Name: "p", ImageURL: "u",
			Texts: []psrt.Text{psrt.NewText(10, 10, 50, 5, 0, psrt.Style(`{}`), "B", "")},
		}},
	}
	variants := []Variant{
		{Label: "a.psrt", Doc: docA},
		{Label: "b.psrt", Doc: docB},
	}
	assets := map[string]compileasset.Asset{
		"u": {MIME: "image/png", Bytes: tinyPNG},
	}
	out, err := RenderHTMLBundle(variants, assets, compileopts.Options{})
	if err != nil {
		t.Fatal(err)
	}
	html := string(out)
	if !strings.Contains(html, "psrt-text") || !strings.Contains(html, "psrt-v-0") || !strings.Contains(html, "psrt-v-1") {
		t.Fatal("expected psrt-text and per-variant classes")
	}
	if !strings.Contains(html, `labels.push("")`) {
		t.Fatal("expected empty off-state label appended in script")
	}
	if !strings.Contains(html, `e.key.toLowerCase()==="l"`) {
		t.Fatal("expected Ctrl+L handler")
	}
	if !strings.Contains(html, ">A<") || !strings.Contains(html, ">B<") {
		t.Fatal("expected both variant texts in HTML")
	}
	if !strings.Contains(html, "psrt-hidden") {
		t.Fatal("expected hidden class for inactive variant")
	}
}

package compilehtml

import (
	"strings"
	"testing"

	"psrt/compileasset"
	"psrt/compileopts"
	"psrt/psrt"
)

func TestRenderHTMLBundle_noScript(t *testing.T) {
	doc := psrt.Document{
		Pages: []psrt.Page{{
			Name: "p", ImageURL: "u",
			Texts: []psrt.Text{psrt.NewText(10, 10, 50, 5, 0, psrt.Style(`{}`), "A", "")},
		}},
	}
	assets := map[string]compileasset.Asset{
		"u": {MIME: "image/png", Bytes: tinyPNG},
	}
	out, err := RenderHTMLBundle([]Variant{{Label: "PSRT", Doc: doc}}, assets, compileopts.Options{NoScript: true})
	if err != nil {
		t.Fatal(err)
	}
	html := string(out)
	if strings.Contains(html, "<script>") {
		t.Fatal("expected no script with NoScript option")
	}
	if strings.Contains(html, "psrt-variant-hint") {
		t.Fatal("expected no variant hint with NoScript option")
	}
}

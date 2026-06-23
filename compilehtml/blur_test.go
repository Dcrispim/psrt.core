package compilehtml

import (
	"strings"
	"testing"

	"psrt/compileasset"
	"psrt/psrt"
)

func TestRenderHTML_textBoxBlur(t *testing.T) {
	doc := psrt.Document{
		Pages: []psrt.Page{
			{
				Name:     "p1",
				ImageURL: "https://example.com/bg.png",
				Texts: []psrt.Text{
					{
						BaseBlock: psrt.BaseBlock{
							X: 10, Y: 20, Width: 80, Index: 1,
							Style: psrt.Style(`{"background":"#00000088","blur":"10px left","color":"#fff"}`),
						},
						TextSize: 5,
						Content:  "Hi",
					},
				},
			},
		},
	}
	assets := map[string]compileasset.Asset{
		"https://example.com/bg.png": {MIME: "image/png", Bytes: tinyPNG},
	}
	out, err := RenderHTML(doc, assets)
	if err != nil {
		t.Fatal(err)
	}
	html := string(out)
	if !strings.Contains(html, "backdrop-filter:blur(10") {
		t.Fatalf("expected backdrop-filter blur, got:\n%s", snippet(html, "text-layer"))
	}
	if !strings.Contains(html, "mask-image:linear-gradient") {
		t.Fatalf("expected directional mask, got:\n%s", snippet(html, "text-layer"))
	}
}

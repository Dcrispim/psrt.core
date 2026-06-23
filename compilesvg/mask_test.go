package compilesvg

import (
	"strings"
	"testing"

	"psrt/compileasset"
	"psrt/psrt"
)

func TestRenderPageSVG_maskBlock(t *testing.T) {
	url := "https://example.com/bg.png"
	assets := map[string]compileasset.Asset{
		url: {Bytes: tinyPNG, MIME: "image/png"},
	}
	page := &psrt.Page{
		Name:     "p",
		ImageURL: url,
		Masks: []psrt.Mask{
			psrt.NewMask(10, 10, 20, 5, 0, psrt.Style(`{"background":"#eee9b2"}`), ""),
		},
	}
	out, err := RenderPageSVG(page, "p", nil, assets)
	if err != nil {
		t.Fatal(err)
	}
	s := string(out)
	if !strings.Contains(s, `id="psrt-text-p-0-wrap"`) {
		t.Fatalf("missing mask group:\n%s", s)
	}
	if !strings.Contains(s, `<rect`) {
		t.Fatal("expected rect for mask")
	}
}

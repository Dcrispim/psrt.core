package compilehtml

import (
	"strings"
	"testing"

	"github.com/Dcrispim/psrt.core/compileasset"
	"github.com/Dcrispim/psrt.core/psrt"
)

const balloonPathHTML = "M10,50 C10,25 30,10 50,10 C70,10 90,25 90,50 Z"

func TestRenderHTML_pathMaskBlock(t *testing.T) {
	doc := psrt.Document{
		Pages: []psrt.Page{
			{
				Name:     "a",
				Style:    psrt.Style(`{}`),
				ImageURL: "https://example.com/a.png",
				PathMasks: []psrt.PathMask{
					psrt.NewPathMask(10, 10, 20, 5, 0, psrt.Style(`{"background":"#eee9b2"}`), "", balloonPathHTML),
				},
			},
		},
	}
	assets := map[string]compileasset.Asset{
		"https://example.com/a.png": {MIME: "image/png", Bytes: tinyPNG},
	}
	out, err := RenderHTML(doc, assets)
	if err != nil {
		t.Fatal(err)
	}
	html := string(out)
	if !strings.Contains(html, "psrt-path-mask") {
		t.Fatalf("missing path mask layer:\n%s", html)
	}
	if !strings.Contains(html, `viewBox="0 0 100 100"`) {
		t.Fatalf("missing nested viewBox svg:\n%s", html)
	}
	if !strings.Contains(html, `preserveAspectRatio="none"`) {
		t.Fatalf("missing non-uniform scaling:\n%s", html)
	}
	if !strings.Contains(html, "<clipPath") {
		t.Fatalf("missing clipPath:\n%s", html)
	}
	if !strings.Contains(html, balloonPathHTML) {
		t.Fatalf("missing raw path data:\n%s", html)
	}
}

func TestRenderHTML_pathMaskBlock_imageRefClipped(t *testing.T) {
	doc := psrt.Document{
		Pages: []psrt.Page{
			{
				Name:     "a",
				Style:    psrt.Style(`{}`),
				ImageURL: "https://example.com/a.png",
				PathMasks: []psrt.PathMask{
					psrt.NewPathMask(10, 10, 20, 5, 0, psrt.Style(`{}`), "https://example.com/fill.png", balloonPathHTML),
				},
			},
		},
	}
	assets := map[string]compileasset.Asset{
		"https://example.com/a.png":    {MIME: "image/png", Bytes: tinyPNG},
		"https://example.com/fill.png": {MIME: "image/png", Bytes: tinyPNG},
	}
	out, err := RenderHTML(doc, assets)
	if err != nil {
		t.Fatal(err)
	}
	html := string(out)
	if !strings.Contains(html, "<image") || !strings.Contains(html, `clip-path="url(#`) {
		t.Fatalf("expected image clipped to path shape:\n%s", html)
	}
}

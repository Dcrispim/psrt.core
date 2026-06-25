package compilesvg

import (
	"strings"
	"testing"

	"github.com/Dcrispim/psrt.core/compileasset"
	"github.com/Dcrispim/psrt.core/psrt"
)

const balloonPathSVG = "M10,50 C10,25 30,10 50,10 C70,10 90,25 90,50 Z"

func TestRenderPageSVG_pathMaskBlock(t *testing.T) {
	url := "https://example.com/bg.png"
	assets := map[string]compileasset.Asset{
		url: {Bytes: tinyPNG, MIME: "image/png"},
	}
	page := &psrt.Page{
		Name:     "p",
		ImageURL: url,
		PathMasks: []psrt.PathMask{
			psrt.NewPathMask(10, 10, 20, 5, 0, psrt.Style(`{"background":"#eee9b2"}`), "", balloonPathSVG),
		},
	}
	out, err := RenderPageSVG(page, "p", nil, assets)
	if err != nil {
		t.Fatal(err)
	}
	s := string(out)
	if !strings.Contains(s, `id="psrt-pathmask-p-0-wrap"`) {
		t.Fatalf("missing path mask group:\n%s", s)
	}
	if !strings.Contains(s, `viewBox="0 0 100 100"`) {
		t.Fatalf("missing nested viewBox svg:\n%s", s)
	}
	if !strings.Contains(s, `preserveAspectRatio="none"`) {
		t.Fatalf("missing non-uniform scaling:\n%s", s)
	}
	if !strings.Contains(s, "<clipPath") {
		t.Fatalf("missing clipPath:\n%s", s)
	}
	if !strings.Contains(s, balloonPathSVG) {
		t.Fatalf("missing raw path data:\n%s", s)
	}
	if !strings.Contains(s, `fill="#eee9b2"`) {
		t.Fatalf("missing fill decoration on path:\n%s", s)
	}
}

func TestRenderPageSVG_pathMaskBlock_imageRefClipped(t *testing.T) {
	bgURL := "https://example.com/bg.png"
	imgURL := "https://example.com/balloon-fill.png"
	assets := map[string]compileasset.Asset{
		bgURL:  {Bytes: tinyPNG, MIME: "image/png"},
		imgURL: {Bytes: tinyPNG, MIME: "image/png"},
	}
	page := &psrt.Page{
		Name:     "p",
		ImageURL: bgURL,
		PathMasks: []psrt.PathMask{
			psrt.NewPathMask(10, 10, 20, 5, 0, psrt.Style(`{}`), imgURL, balloonPathSVG),
		},
	}
	out, err := RenderPageSVG(page, "p", nil, assets)
	if err != nil {
		t.Fatal(err)
	}
	s := string(out)
	if !strings.Contains(s, "<image") || !strings.Contains(s, `clip-path="url(#`) {
		t.Fatalf("expected image clipped to path shape:\n%s", s)
	}
}

package styleadapter

import (
	"testing"

	"github.com/Dcrispim/psrt.core/psrt"
)

func TestAdaptPathMaskSVG_decoration(t *testing.T) {
	m := psrt.NewPathMask(10, 10, 20, 5, 0, psrt.Style(`{"background":"#000","borderRadius":"4px"}`), "", "M0,0 L1,1 Z")
	ctx := AdaptContext{PathMask: &m, CanvasW: 200, CanvasH: 100, PageSlug: "p", TextIndex: 0}
	frags := AdaptPathMaskSVG(ctx)

	var path StyleFragment
	for _, f := range frags {
		if f.GetString(TypeKey) == TypePath {
			path = f
		}
	}
	if path == nil {
		t.Fatal("missing path fragment")
	}
	if path.GetString("fill") != "#000" {
		t.Fatalf("fill: %q", path.GetString("fill"))
	}
	if _, hasRx := path["rx"]; hasRx {
		t.Fatal("border-radius must be ignored for path masks (no rx on path fragment)")
	}
}

func TestAdaptPathMaskHTML_layoutAndDecoration(t *testing.T) {
	m := psrt.NewPathMask(10, 20, 30, 40, 0, psrt.Style(`{"background":"#fff"}`), "", "M0,0 L1,1 Z")
	ctx := AdaptContext{PathMask: &m, CanvasW: 1000, CanvasH: 500, HTMLCompile: true}
	frags := AdaptPathMaskHTML(ctx)

	var box, path StyleFragment
	for _, f := range frags {
		switch f.GetString(TypeKey) {
		case TypeMotionDiv:
			box = f
		case TypePath:
			path = f
		}
	}
	if box == nil {
		t.Fatal("missing layout box fragment")
	}
	if box.GetString(KeyLeft) != "10%" || box.GetString(KeyTop) != "20%" ||
		box.GetString(KeyWidth) != "30%" || box.GetString(KeyHeight) != "40%" {
		t.Fatalf("layout: %+v", box)
	}
	if path == nil || path.GetString("fill") != "#fff" {
		t.Fatalf("path decoration: %+v", path)
	}
}

func TestPathDecorationAttrs_excludesRx(t *testing.T) {
	f := NewFragment(TypePath)
	f.Set("fill", "#abc")
	f.Set("rx", "4px")
	got := PathDecorationAttrs(f)
	if got != `fill="#abc"` {
		t.Fatalf("PathDecorationAttrs: got %q, want only fill (no rx)", got)
	}
}

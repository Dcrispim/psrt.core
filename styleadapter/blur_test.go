package styleadapter

import (
	"strings"
	"testing"

	"github.com/Dcrispim/psrt.core/psrt"
)

func TestAdaptHTML_blurBackdrop(t *testing.T) {
	ctx := AdaptContext{
		Text: psrt.Text{
			BaseBlock: psrt.BaseBlock{X: 10, Y: 20, Width: 50, Style: psrt.Style(`{"background":"#0008","blur":"8px"}`)},
			TextSize:  5,
		},
		CanvasW: 1000, CanvasH: 500, FontSizePx: 20,
	}
	frags := AdaptHTML(ctx)
	var box StyleFragment
	for _, f := range frags {
		if f.GetString(TypeKey) == TypeMotionDiv {
			box = f
		}
	}
	if box == nil {
		t.Fatal("missing box fragment")
	}
	if !strings.Contains(box.GetString("backdropFilter"), "blur(8") {
		t.Fatalf("backdropFilter: %q", box.GetString("backdropFilter"))
	}
}

func TestAdaptHTML_blurLeftMask(t *testing.T) {
	ctx := AdaptContext{
		Text: psrt.Text{
			BaseBlock: psrt.BaseBlock{Style: psrt.Style(`{"blurLeft":"12px","background":"#fff"}`)},
		},
		CanvasW: 800, CanvasH: 600, FontSizePx: 24,
	}
	frags := AdaptHTML(ctx)
	var box StyleFragment
	for _, f := range frags {
		if f.GetString(TypeKey) == TypeMotionDiv {
			box = f
		}
	}
	if box == nil {
		t.Fatal("missing box")
	}
	if !strings.Contains(box.GetString("maskImage"), "linear-gradient") {
		t.Fatalf("maskImage: %q", box.GetString("maskImage"))
	}
	if !strings.Contains(box.GetString("backdropFilter"), "blur(12") {
		t.Fatalf("backdropFilter: %q", box.GetString("backdropFilter"))
	}
}

func TestAdaptSVG_blurFilterAndMask(t *testing.T) {
	ctx := AdaptContext{
		Text:       psrt.Text{BaseBlock: psrt.BaseBlock{X: 0, Y: 0, Width: 50, Style: psrt.Style(`{"blur":"6px left","background":"#eaede5"}`)}, TextSize: 5},
		CanvasW:    200, CanvasH: 100, FontSizePx: 5,
		PageSlug:   "p", TextIndex: 2,
	}
	frags := AdaptSVG(ctx)
	var rect StyleFragment
	var filterSVG, maskSVG string
	for _, f := range frags {
		switch f.GetString(TypeKey) {
		case TypeRect:
			rect = f
		case TypeFilter:
			filterSVG = FilterFragmentToSVG(f)
		case TypeMask:
			maskSVG = MaskFragmentToSVG(f)
		}
	}
	if rect == nil {
		t.Fatal("missing rect")
	}
	if !strings.Contains(rect.GetString("filter"), "psrt-filter-p-2-blur") {
		t.Fatalf("rect filter: %q", rect.GetString("filter"))
	}
	if !strings.Contains(rect.GetString("mask"), "mask") {
		t.Fatalf("rect mask: %q", rect.GetString("mask"))
	}
	if !strings.Contains(filterSVG, "feGaussianBlur") {
		t.Fatalf("filter svg: %s", filterSVG)
	}
	if !strings.Contains(maskSVG, "linearGradient") {
		t.Fatalf("mask svg: %s", maskSVG)
	}
}

func TestResolveName_blurAliases(t *testing.T) {
	got, ok := ResolveName("blur-left")
	if !ok || got != KeyBlurLeft {
		t.Fatalf("blur-left: %q %v", got, ok)
	}
}

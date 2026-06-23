package compileasset

import (
	"testing"

	"psrt/psrt"
)

func TestBorderRadiusNeedsInnerBoxDecoration_bottomOnly(t *testing.T) {
	c := BorderRadiusCorners{0, 0, 20, 15}
	if !BorderRadiusNeedsInnerBoxDecoration(c) {
		t.Fatal("expected inner box decoration for bottom-only radius")
	}
}

func TestBorderRadiusNeedsInnerBoxDecoration_uniform(t *testing.T) {
	c := BorderRadiusCorners{10, 10, 10, 10}
	if BorderRadiusNeedsInnerBoxDecoration(c) {
		t.Fatal("uniform radius should use SVG rect")
	}
}

func TestParseBorderRadiusCorners_shorthand(t *testing.T) {
	style := psrt.Style(`{"border-top-left-radius":"0px","border-top-right-radius":"0px","border-bottom-right-radius":"20px","border-bottom-left-radius":"15px"}`)
	c := ParseBorderRadiusCorners(style, 1080, 1920, 48)
	if c.TopLeft != 0 || c.BottomRight != 20 || c.BottomLeft != 15 {
		t.Fatalf("got %+v", c)
	}
}

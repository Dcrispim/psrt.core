package compileasset

import (
	"strings"
	"testing"

	"github.com/Dcrispim/psrt.core/psrt"
)

func TestCSSBoxFromStyleJSONNoBackgroundForCanvas_percentPadding(t *testing.T) {
	css := CSSBoxFromStyleJSONNoBackgroundForCanvas(
		psrt.Style(`{"padding":"0.461%"}`),
		700, 6585, 48,
	)
	if strings.Contains(css, "%") {
		t.Fatalf("padding must be resolved to px for SVG foreignObject, got:\n%s", css)
	}
	if !strings.Contains(css, "padding:") || !strings.Contains(css, "px") {
		t.Fatalf("expected px padding, got:\n%s", css)
	}
}

func TestTextBlockDisplayCSS_rightAndJustify(t *testing.T) {
	right := TextBlockDisplayCSS(psrt.Style(`{"text-align":"right"}`))
	if !strings.Contains(right, "align-items:flex-end") {
		t.Fatalf("right align needs flex-end, got %q", right)
	}
	if !strings.Contains(right, "text-align:right") {
		t.Fatalf("expected text-align:right, got %q", right)
	}

	justify := TextBlockDisplayCSS(psrt.Style(`{"text-align":"justify"}`))
	if !strings.Contains(justify, "text-align:justify") {
		t.Fatalf("justify must emit text-align, got %q", justify)
	}
}

func TestTextBlockDisplayCSS_centerAndFlexEnd(t *testing.T) {
	css := TextBlockDisplayCSS(psrt.Style(`{"text-align":"center","align-items":"flex-end"}`))
	if !strings.Contains(css, "display:flex") {
		t.Fatalf("expected flex, got %q", css)
	}
	if !strings.Contains(css, "justify-content:flex-end") {
		t.Fatalf("align-items flex-end must map to justify-content flex-end, got %q", css)
	}
	if !strings.Contains(css, "text-align:center") {
		t.Fatalf("expected text-align center, got %q", css)
	}
}

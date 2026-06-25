package compileasset

import (
	"strings"
	"testing"

	"github.com/Dcrispim/psrt.core/psrt"
)

func TestParseTextBox_background(t *testing.T) {
	box := ParseTextBox(psrt.Style(`{"background":"#000000ff","border-radius":"4px"}`))
	if box.Background != "#000000ff" {
		t.Fatalf("Background = %q", box.Background)
	}
	if !box.HasRect {
		t.Fatal("expected HasRect")
	}
	if box.RadiusPx != 4 {
		t.Fatalf("RadiusPx = %v", box.RadiusPx)
	}
}

func TestCSSBoxFromStyleJSON(t *testing.T) {
	css := CSSBoxFromStyleJSON(psrt.Style(`{"backgroundColor":"#000","font-weight":"600"}`))
	if !strings.Contains(css, "background-color:#000") {
		t.Fatalf("missing background in %q", css)
	}
}

func TestStyleJSONWithoutBox_keepsTypography(t *testing.T) {
	out := StyleJSONWithoutBox(psrt.Style(`{"color":"#fff","background":"#000","font-weight":"600"}`))
	s := string(out)
	if strings.Contains(s, "background") {
		t.Fatalf("box keys should be removed: %s", s)
	}
	if !strings.Contains(s, "font-weight") && !strings.Contains(s, "fontWeight") {
		t.Fatalf("typography should remain: %s", s)
	}
}

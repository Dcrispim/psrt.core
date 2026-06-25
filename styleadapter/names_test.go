package styleadapter

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/Dcrispim/psrt.core/psrt"
	"github.com/Dcrispim/psrt.core/styleadapter/percent"
)

func TestResolveName_aliases(t *testing.T) {
	tests := []struct {
		in   string
		want string
		ok   bool
	}{
		{"br", KeyBorderRadius, true},
		{"border-radius", KeyBorderRadius, true},
		{"borderRadius", KeyBorderRadius, true},
		{"sw", KeyStrokeWidth, true},
		{"stroke-width", KeyStrokeWidth, true},
	}
	for _, tc := range tests {
		got, ok := ResolveName(tc.in)
		if got != tc.want || ok != tc.ok {
			t.Fatalf("ResolveName(%q) = (%q, %v), want (%q, %v)", tc.in, got, ok, tc.want, tc.ok)
		}
	}
}

func TestNormalize_webkitAndCanonical(t *testing.T) {
	m := Normalize(psrt.Style(`{"WebkitTextStrokeWidth":"2","strokeWidth":"1"}`))
	if _, ok := m["WebkitTextStrokeWidth"]; ok {
		t.Fatal("normalized map must not contain Webkit keys")
	}
	if got := StringifyCSSValue(m[KeyStrokeWidth]); got != "1" {
		t.Fatalf("strokeWidth: got %q want 1", got)
	}
}

func TestNormalize_webkitOnly(t *testing.T) {
	m := Normalize(psrt.Style(`{"WebkitTextStrokeWidth":"3"}`))
	if got := StringifyCSSValue(m[KeyStrokeWidth]); got != "3" {
		t.Fatalf("strokeWidth: got %q", got)
	}
}

func TestApplyPercentHandlers_textShadow(t *testing.T) {
	style := map[string]json.RawMessage{
		KeyTextShadow: json.RawMessage(`"1% 2% 3% rgba(0,0,0,0.5)"`),
	}
	out := percent.ApplyPercentHandlers(style, percent.ImageDims{W: 1000, H: 500, Zoom: 1})
	got := StringifyCSSValue(out[KeyTextShadow])
	if !strings.Contains(got, "px") {
		t.Fatalf("expected px in resolved shadow, got %q", got)
	}
}

func TestAdaptHTML_layoutAndStroke(t *testing.T) {
	ctx := AdaptContext{
		Text: psrt.Text{BaseBlock: psrt.BaseBlock{X: 10, Y: 20, Width: 30, Style: psrt.Style(`{"strokeWidth":"2","color":"#fff"}`)}, TextSize: 4},
		CanvasW: 1000, CanvasH: 500, FontSizePx: 20,
	}
	frags := AdaptHTML(ctx)
	var box, span StyleFragment
	for _, f := range frags {
		switch f.GetString(TypeKey) {
		case TypeMotionDiv:
			box = f
		case TypeSpan:
			span = f
		}
	}
	if box.GetString(KeyLeft) != "10%" || box.GetString(KeyTop) != "20%" {
		t.Fatalf("layout: left=%q top=%q", box.GetString(KeyLeft), box.GetString(KeyTop))
	}
	if span.GetString("WebkitTextStrokeWidth") != "2" {
		t.Fatalf("expected WebkitTextStrokeWidth on span, got %q", span.GetString("WebkitTextStrokeWidth"))
	}
	if span.GetString(KeyStrokeWidth) != "" {
		t.Fatal("span must not expose strokeWidth PSRT key")
	}
}

func TestAdaptSVG_splitBoxAndText(t *testing.T) {
	ctx := AdaptContext{
		Text:       psrt.Text{BaseBlock: psrt.BaseBlock{X: 0, Y: 0, Width: 50, Style: psrt.Style(`{"background":"#000","color":"#fff"}`)}, TextSize: 5},
		CanvasW:    200, CanvasH: 100, FontSizePx: 5,
		PageSlug:   "p", TextIndex: 0,
	}
	frags := AdaptSVG(ctx)
	var rect, inner StyleFragment
	for _, f := range frags {
		switch f.GetString(TypeKey) {
		case TypeRect:
			rect = f
		case TypeMotionDiv:
			inner = f
		}
	}
	if rect.GetString("fill") != "#000" {
		t.Fatalf("rect fill: %q", rect.GetString("fill"))
	}
	if inner.GetString(KeyBackground) != "" && inner.GetString("backgroundColor") != "" {
		t.Fatal("inner text host must not have background")
	}
	if inner.GetString(KeyColor) != "#fff" {
		t.Fatalf("text color: %q", inner.GetString(KeyColor))
	}
}

func TestAdaptHTML_glow(t *testing.T) {
	ctx := AdaptContext{
		Text: psrt.Text{BaseBlock: psrt.BaseBlock{Style: psrt.Style(`{"glow":"1% 1% 2% rgba(0,0,0,0.5)"}`)}},
		CanvasW: 100, CanvasH: 100, FontSizePx: 10,
	}
	frags := AdaptHTML(ctx)
	found := false
	for _, f := range frags {
		if f.GetString(TypeKey) == TypeSpan && f.GetString(KeyTextShadow) != "" {
			found = true
		}
	}
	if !found {
		t.Fatal("glow should expand to textShadow on span")
	}
}

package percent

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestTextShadowPercent(t *testing.T) {
	style := map[string]json.RawMessage{
		"textShadow": json.RawMessage(`"1% 2% 3% rgba(0,0,0,0.5)"`),
	}
	out := ApplyPercentHandlers(style, ImageDims{W: 1000, H: 500, Zoom: 1})
	got := string(out["textShadow"])
	if got == string(style["textShadow"]) {
		t.Fatal("expected change")
	}
	if !strings.Contains(got, "px") {
		t.Fatalf("got %s", got)
	}
}

func TestStrokeWidthPercent(t *testing.T) {
	style := map[string]json.RawMessage{
		"strokeWidth": json.RawMessage(`"10%"`),
	}
	out := ApplyPercentHandlers(style, ImageDims{FontSizePx: 20, Zoom: 1})
	got := string(out["strokeWidth"])
	if !strings.Contains(got, "2.") && !strings.Contains(got, "px") {
		t.Fatalf("expected ~2px, got %s", got)
	}
}

func TestUnregisteredKeyUnchanged(t *testing.T) {
	style := map[string]json.RawMessage{
		"color": json.RawMessage(`"#fff"`),
	}
	out := ApplyPercentHandlers(style, ImageDims{W: 100, H: 100})
	if string(out["color"]) != `"#fff"` {
		t.Fatalf("color changed: %s", out["color"])
	}
}

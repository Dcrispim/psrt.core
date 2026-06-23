package styleadapter

import (
	"encoding/json"
	"testing"

	"psrt/psrt"
)

func TestAdaptWebPreview_backgroundAndFontSize(t *testing.T) {
	ctx := AdaptContext{
		Text: psrt.Text{
			X: 10, Y: 20, Width: 50, TextSize: 5,
			Style: psrt.Style(`{"background":"#000000ff","color":"#fff","font-weight":"600"}`),
		},
		CanvasW: 1920, CanvasH: 1080, FontSizePx: 54, Zoom: 2,
	}
	got := AdaptWebPreview(ctx)
	if got.Container["backgroundColor"] != "#000000ff" {
		t.Fatalf("backgroundColor: %q container=%v", got.Container["backgroundColor"], got.Container)
	}
	if got.Container["fontSize"] != "54px" {
		t.Fatalf("fontSize: %q", got.Container["fontSize"])
	}
	if got.Text["color"] != "#fff" {
		t.Fatalf("color: %q", got.Text["color"])
	}
}

func TestAdaptWebPreview_zoomScalesFont(t *testing.T) {
	ctx := AdaptContext{
		Text: psrt.Text{
			TextSize: 10,
			Style:    psrt.Style(`{}`),
		},
		CanvasW: 1000, CanvasH: 1000, Zoom: 2,
	}
	got := AdaptWebPreview(ctx)
	// 10% of 1000 = 100px base; zoom should apply in FontSizePx from caller
	ctx.FontSizePx = 10.0 / 100.0 * 1000 * 2
	got2 := AdaptWebPreview(ctx)
	_ = got
	if got2.Container["fontSize"] != "200px" {
		t.Fatalf("zoom font: %q", got2.Container["fontSize"])
	}
}

func TestAdaptWebPreview_stringifyRoundtrip(t *testing.T) {
	raw := `{"backGround":"#000","color":"#fff"}`
	ctx := AdaptContext{
		Text: psrt.Text{
			Style: psrt.Style(raw),
		},
		CanvasW: 800, CanvasH: 600, FontSizePx: 30,
	}
	got := AdaptWebPreview(ctx)
	b, _ := json.Marshal(got)
	if got.Container["backgroundColor"] == "" {
		t.Fatalf("missing bg: %s", b)
	}
}

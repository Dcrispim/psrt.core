package styleadapter

import (
	"testing"

	"github.com/Dcrispim/psrt.core/psrt"
)

func TestAdaptMaskWebPreview_heightPercent(t *testing.T) {
	ctx := AdaptContext{
		Mask: &psrt.Mask{
			BaseBlock: psrt.BaseBlock{
				X: 16.35, Y: 8.87, Width: 22.49,
				Style: psrt.Style(`{"background":"#edefe2","border-radius":"100px"}`),
			},
			Height: 1.98,
		},
		CanvasW: 1920,
		CanvasH: 1080,
		Zoom:    1,
	}
	got := AdaptMaskWebPreview(ctx)
	if got.Container["height"] != "1.98%" {
		t.Fatalf("height: %q container=%v", got.Container["height"], got.Container)
	}
	if got.Container["width"] != "22.49%" {
		t.Fatalf("width: %q", got.Container["width"])
	}
	if got.Container["backgroundColor"] == "" && got.Container["background"] == "" {
		t.Fatalf("expected background, container=%v", got.Container)
	}
}

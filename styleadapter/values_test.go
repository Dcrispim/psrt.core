package styleadapter

import (
	"encoding/json"
	"testing"

	"psrt/psrt"
)

func TestNormalize_omitsZeroHeight(t *testing.T) {
	m := Normalize(psrt.Style(`{"background":"#000","height":0,"color":"#fff"}`))
	if _, ok := m[KeyHeight]; ok {
		t.Fatal("height 0 should be omitted")
	}
	if !HasStyleValue(KeyBackground, m[KeyBackground]) {
		t.Fatal("background should remain")
	}
}

func TestNormalize_omitsFalseBoolean(t *testing.T) {
	m := Normalize(psrt.Style(`{"height":false,"padding":true}`))
	if _, ok := m[KeyHeight]; ok {
		t.Fatal("false boolean must not become height")
	}
}

func TestAdaptWebPreview_noHeightWhenAbsent(t *testing.T) {
	ctx := AdaptContext{
		Text: psrt.Text{
			BaseBlock: psrt.BaseBlock{Style: psrt.Style(`{"background":"#000000ff","color":"#fff"}`)},
		},
		CanvasW: 1920, CanvasH: 1080, FontSizePx: 54,
	}
	got := AdaptWebPreview(ctx)
	if _, ok := got.Container["height"]; ok {
		t.Fatalf("unexpected height: %v", got.Container)
	}
}

func TestFilterStyleMap_afterPercentZeroPx(t *testing.T) {
	m := map[string]json.RawMessage{
		KeyPadding: []byte(`"0%"`),
	}
	// Simulate percent handler output for 0%
	m[KeyPadding] = []byte(`"0px"`)
	out := FilterStyleMap(m)
	if _, ok := out[KeyPadding]; ok {
		t.Fatal("0px padding should be filtered")
	}
}

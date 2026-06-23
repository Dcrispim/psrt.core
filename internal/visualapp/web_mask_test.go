package visualapp

import (
	"encoding/json"
	"testing"
)

func TestAdaptEntriesForWeb_mask(t *testing.T) {
	inputs := []WebEntryStyleInput{{
		Index: 0,
		Style: `{"background":"#edefe2","border-radius":"100px"}`,
		X:     16.35, Y: 8.87, Width: 22.49,
		Height: 1.98, IsMask: true,
	}}
	raw, err := json.Marshal(inputs)
	if err != nil {
		t.Fatal(err)
	}
	out, err := AdaptEntriesForWeb(string(raw), 1920, 1080, 1)
	if err != nil {
		t.Fatal(err)
	}
	if len(out) != 1 {
		t.Fatalf("len=%d", len(out))
	}
	if out[0].Container["height"] != "1.98%" {
		t.Fatalf("height: %q container=%v", out[0].Container["height"], out[0].Container)
	}
}

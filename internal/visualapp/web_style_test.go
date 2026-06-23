package visualapp

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestAdaptTextStyleForWeb_strokeAndPercent(t *testing.T) {
	got := AdaptTextStyleForWeb(
		`{"strokeWidth":"10%","color":"#fff"}`,
		"",
		10, 20, 50, 5,
		1000, 500,
		1,
	)
	if got.Text["WebkitTextStrokeWidth"] == "" {
		t.Fatalf("expected WebkitTextStrokeWidth, text=%v", got.Text)
	}
	if !strings.Contains(got.Text["WebkitTextStrokeWidth"], "px") {
		t.Fatalf("stroke should be px: %q", got.Text["WebkitTextStrokeWidth"])
	}
	if got.Container["left"] != "10%" || got.Container["top"] != "20%" {
		t.Fatalf("layout: %+v", got.Container)
	}
}

func TestAdaptEntriesForWeb_batch(t *testing.T) {
	inputs := []WebEntryStyleInput{
		{Index: 0, Style: `{"color":"#000"}`, X: 0, Y: 0, Width: 10, TextSize: 5},
	}
	raw, _ := json.Marshal(inputs)
	out, err := AdaptEntriesForWeb(string(raw), 100, 100, 1)
	if err != nil || len(out) != 1 {
		t.Fatalf("batch: %v len=%d", err, len(out))
	}
}

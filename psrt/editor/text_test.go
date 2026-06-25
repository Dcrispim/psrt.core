package editor

import (
	"testing"

	"github.com/Dcrispim/psrt.core/psrt"
)

func TestSetTextStyleAndContent(t *testing.T) {
	doc := sampleDoc(t)
	if err := SetTextStyle(&doc, "p1", 0, "color", `"#1DB954"`, nil); err != nil {
		t.Fatal(err)
	}
	if err := SetTextContent(&doc, "p1", 0, "updated", false); err != nil {
		t.Fatal(err)
	}
	t0, _, _ := FindTextByIndex(&doc.Pages[0], 0)
	if t0.Content != "updated" {
		t.Fatalf("content: %q", t0.Content)
	}
	if err := SetTextContent(&doc, "p1", 1, "!", true); err != nil {
		t.Fatal(err)
	}
	t1, _, _ := FindTextByIndex(&doc.Pages[0], 1)
	if t1.Content != "world!" {
		t.Fatalf("append: %q", t1.Content)
	}
}

func TestAddAndRemoveText(t *testing.T) {
	doc := sampleDoc(t)
	newText := psrt.Text{
		BaseBlock: psrt.BaseBlock{X: 5, Y: 5, Width: 50, Index: 9, Style: psrt.Style("{}")},
		TextSize:  1,
		Content:   "new",
	}
	if err := AddText(&doc, "p1", newText, -1, -1); err != nil {
		t.Fatal(err)
	}
	if _, _, err := FindTextByIndex(&doc.Pages[0], 9); err != nil {
		t.Fatal(err)
	}
	if err := RemoveText(&doc, "p1", 9); err != nil {
		t.Fatal(err)
	}
}

func TestReorderTextRelativeAndDelta(t *testing.T) {
	doc := sampleDoc(t)
	if err := ReorderTextRelative(&doc, "p1", 1, 0, -1); err != nil {
		t.Fatal(err)
	}
	if doc.Pages[0].Texts[0].Index != 1 {
		t.Fatalf("expected index 1 first, got %d", doc.Pages[0].Texts[0].Index)
	}
	doc = sampleDoc(t)
	if err := ReorderTextByDelta(&doc, "p1", 0, 1); err != nil {
		t.Fatal(err)
	}
	if doc.Pages[0].Texts[1].Index != 0 {
		t.Fatalf("expected index 0 second, got %d", doc.Pages[0].Texts[1].Index)
	}
}

func TestReorderTextTo(t *testing.T) {
	doc := sampleDoc(t)
	if err := ReorderTextTo(&doc, "p1", 0, 1); err != nil {
		t.Fatal(err)
	}
	if doc.Pages[0].Texts[1].Index != 0 {
		t.Fatalf("index 0 should be at position 1, got %d", doc.Pages[0].Texts[1].Index)
	}
}

func TestSetAndNudgeTextPosition(t *testing.T) {
	doc := sampleDoc(t)
	x, y, w, s := 50.0, 60.0, 80.0, 2.0
	if err := SetTextPosition(&doc, "p1", 0, PositionFields{X: &x, Y: &y, Width: &w, TextSize: &s}); err != nil {
		t.Fatal(err)
	}
	t0, _, _ := FindTextByIndex(&doc.Pages[0], 0)
	if t0.X != 50 || t0.Y != 60 || t0.Width != 80 || t0.TextSize != 2 {
		t.Fatalf("position: %+v", t0)
	}
	dx, dy := 1.5, -2.0
	if err := NudgeTextPosition(&doc, "p1", 0, PositionFields{X: &dx, Y: &dy}); err != nil {
		t.Fatal(err)
	}
	if t0.X != 51.5 || t0.Y != 58 {
		t.Fatalf("nudge: x=%v y=%v", t0.X, t0.Y)
	}
}

package editor

import (
	"testing"

	"psrt/psrt"
)

const samplePathD = "M10,50 C10,25 30,10 50,10 C70,10 90,25 90,50 Z"

func TestAddAndRemovePathMask(t *testing.T) {
	doc := sampleDoc(t)
	newMask := psrt.NewPathMask(5, 5, 50, 30, 9, psrt.Style("{}"), "", samplePathD)
	if err := AddPathMask(&doc, "p1", newMask, -1, -1); err != nil {
		t.Fatal(err)
	}
	if _, _, err := FindPathMaskByIndex(&doc.Pages[0], 9); err != nil {
		t.Fatal(err)
	}
	if err := RemovePathMask(&doc, "p1", 9); err != nil {
		t.Fatal(err)
	}
	if _, _, err := FindPathMaskByIndex(&doc.Pages[0], 9); err == nil {
		t.Fatal("expected path mask to be removed")
	}
}

func TestAddPathMaskNotFoundCases(t *testing.T) {
	doc := sampleDoc(t)
	newMask := psrt.NewPathMask(0, 0, 10, 10, 9, psrt.Style("{}"), "", samplePathD)
	if err := AddPathMask(&doc, "missing-page", newMask, -1, -1); err == nil {
		t.Fatal("expected error for missing page")
	}
	if err := RemovePathMask(&doc, "p1", 99); err == nil {
		t.Fatal("expected error for missing index")
	}
}

func TestSetPathMaskPosition(t *testing.T) {
	doc := sampleDoc(t)
	newMask := psrt.NewPathMask(5, 5, 50, 30, 9, psrt.Style("{}"), "", samplePathD)
	if err := AddPathMask(&doc, "p1", newMask, -1, -1); err != nil {
		t.Fatal(err)
	}
	x := 12.5
	height := 40.0
	if err := SetPathMaskPosition(&doc, "p1", 9, PathMaskPositionFields{X: &x, Height: &height}); err != nil {
		t.Fatal(err)
	}
	m, _, _ := FindPathMaskByIndex(&doc.Pages[0], 9)
	if m.X != 12.5 || m.Height != 40 {
		t.Fatalf("position not applied: x=%v height=%v", m.X, m.Height)
	}
	if err := SetPathMaskPosition(&doc, "p1", 9, PathMaskPositionFields{}); err == nil {
		t.Fatal("expected error for empty position fields")
	}
}

func TestSetPathMaskStyleAndRemoveKey(t *testing.T) {
	doc := sampleDoc(t)
	newMask := psrt.NewPathMask(5, 5, 50, 30, 9, psrt.Style("{}"), "", samplePathD)
	if err := AddPathMask(&doc, "p1", newMask, -1, -1); err != nil {
		t.Fatal(err)
	}
	if err := SetPathMaskStyle(&doc, "p1", 9, "background", `"#eee9b2"`, nil); err != nil {
		t.Fatal(err)
	}
	m, _, _ := FindPathMaskByIndex(&doc.Pages[0], 9)
	if string(m.Style) != `{"background":"#eee9b2"}` {
		t.Fatalf("style not applied: %s", m.Style)
	}
	if err := RemovePathMaskStyleKey(&doc, "p1", 9, "background"); err != nil {
		t.Fatal(err)
	}
	m, _, _ = FindPathMaskByIndex(&doc.Pages[0], 9)
	if string(m.Style) != `{}` {
		t.Fatalf("style key not removed: %s", m.Style)
	}
}

func TestSetPathMaskPath(t *testing.T) {
	doc := sampleDoc(t)
	newMask := psrt.NewPathMask(5, 5, 50, 30, 9, psrt.Style("{}"), "", samplePathD)
	if err := AddPathMask(&doc, "p1", newMask, -1, -1); err != nil {
		t.Fatal(err)
	}
	const updated = "M0,0 L100,0 L100,100 Z"
	if err := SetPathMaskPath(&doc, "p1", 9, updated); err != nil {
		t.Fatal(err)
	}
	m, _, _ := FindPathMaskByIndex(&doc.Pages[0], 9)
	if m.Path != updated {
		t.Fatalf("path not applied: %s", m.Path)
	}
}

func TestSetPathMaskPathInvalid(t *testing.T) {
	doc := sampleDoc(t)
	newMask := psrt.NewPathMask(5, 5, 50, 30, 9, psrt.Style("{}"), "", samplePathD)
	if err := AddPathMask(&doc, "p1", newMask, -1, -1); err != nil {
		t.Fatal(err)
	}
	if err := SetPathMaskPath(&doc, "p1", 9, ""); err == nil {
		t.Fatal("expected error for empty path")
	}
	if err := SetPathMaskPath(&doc, "p1", 9, "not a path"); err == nil {
		t.Fatal("expected error for invalid svg path syntax")
	}
	if err := SetPathMaskPath(&doc, "p1", 9, "M0,0 L10,10 M20,20 L30,30"); err == nil {
		t.Fatal("expected error for multiple subpaths")
	}
}

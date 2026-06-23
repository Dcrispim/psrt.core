package compileasset

import (
	"testing"

	"psrt/psrt"
)

func TestExplicitHeightPx_percent(t *testing.T) {
	h, ok := ExplicitHeightPx(psrt.Style(`{"height":"20%"}`), 1080, 1920, 48)
	if !ok {
		t.Fatal("expected ok")
	}
	if h != 384 {
		t.Fatalf("got %d want 384", h)
	}
}

func TestExplicitHeightPx_absent(t *testing.T) {
	_, ok := ExplicitHeightPx(psrt.Style(`{"color":"#000"}`), 1080, 1920, 48)
	if ok {
		t.Fatal("expected no height")
	}
}

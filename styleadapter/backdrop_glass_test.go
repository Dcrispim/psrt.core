package styleadapter

import "testing"

func TestSoftenOpaqueColorForBackdrop(t *testing.T) {
	got, ok := softenOpaqueColorForBackdrop("#f3f2e7", backdropGlassAlpha)
	if !ok {
		t.Fatal("expected ok")
	}
	if got != "rgba(243,242,231,0.720)" {
		t.Fatalf("got %q", got)
	}
}

func TestSoftenOpaqueColorForBackdrop_skipsTranslucent(t *testing.T) {
	_, ok := softenOpaqueColorForBackdrop("#00000088", backdropGlassAlpha)
	if ok {
		t.Fatal("expected skip for hex8 with alpha")
	}
	_, ok = softenOpaqueColorForBackdrop("rgba(0,0,0,0.5)", backdropGlassAlpha)
	if ok {
		t.Fatal("expected skip for rgba with alpha")
	}
}

func TestPostProcessBackdropGlass(t *testing.T) {
	box := NewFragment(TypeMotionDiv)
	box.Set("backgroundColor", "#f3f2e7")
	box.Set("backdropFilter", "blur(45px)")
	box.Set("borderRadius", "100px")
	postProcessBackdropGlass(box)
	if box.GetString("overflow") != "hidden" {
		t.Fatalf("overflow: %q", box.GetString("overflow"))
	}
	if box.GetString("backgroundColor") != "rgba(243,242,231,0.720)" {
		t.Fatalf("backgroundColor: %q", box.GetString("backgroundColor"))
	}
}

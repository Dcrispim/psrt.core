package visualapp

import "testing"

func TestAdaptTextStyleForWeb_computedHeightWithPadding(t *testing.T) {
	got := AdaptTextStyleForWeb(
		`{"padding":"2%","text-align":"center","align-items":"flex-end","color":"#000"}`,
		"line one",
		10, 20, 50, 5,
		1000, 2000,
		1,
	)
	if got.Container["height"] == "" {
		t.Fatalf("expected computed height for vertical flex, container=%v", got.Container)
	}
	if got.Container["justifyContent"] != "flex-end" {
		t.Fatalf("justifyContent: %q", got.Container["justifyContent"])
	}
	if got.Container["display"] != "flex" {
		t.Fatalf("display: %q", got.Container["display"])
	}
}

func TestAdaptTextStyleForWeb_explicitHeightNotOverwritten(t *testing.T) {
	got := AdaptTextStyleForWeb(
		`{"padding":"2%","height":"5%","text-align":"center","align-items":"flex-end"}`,
		"",
		10, 20, 50, 5,
		1000, 2000,
		1,
	)
	if got.Container["height"] == "" {
		t.Fatalf("explicit height must stay resolved: %v", got.Container)
	}
}

func TestAdaptTextStyleForWeb_noHeightWithoutPadding(t *testing.T) {
	got := AdaptTextStyleForWeb(
		`{"text-align":"center","align-items":"flex-end"}`,
		"hello",
		10, 20, 50, 5,
		1000, 2000,
		1,
	)
	if got.Container["height"] != "" {
		t.Fatalf("without padding/height, container height must stay unset: %q", got.Container["height"])
	}
}

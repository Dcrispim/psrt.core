package psrt

import "testing"

func TestRoundCoord(t *testing.T) {
	cases := []struct {
		in, want float64
	}{
		{11.6, 11.6},
		{11.649999999999999, 11.65},
		{50, 50},
		{11.123456789, 11.12346},
		{0.000001, 0},
	}
	for _, tc := range cases {
		got := RoundCoord(tc.in)
		if got != tc.want {
			t.Errorf("RoundCoord(%v) = %v, want %v", tc.in, got, tc.want)
		}
	}
}

func TestTextFontSizePx_usesMinDimension(t *testing.T) {
	if got := TextSizeBasisPx(1080, 1920); got != 1080 {
		t.Fatalf("basis: got %v want 1080", got)
	}
	if got := TextFontSizePx(3, 1080, 1920); got != 32.4 {
		t.Fatalf("font: got %v want 32.4", got)
	}
	if got := TextFontSizePx(3, 1920, 1080); got != 32.4 {
		t.Fatalf("font swapped: got %v want 32.4", got)
	}
}

func TestFormatCoord(t *testing.T) {
	cases := []struct {
		in   float64
		want string
	}{
		{11.6, "11.6"},
		{50, "50"},
		{11.65, "11.65"},
		{11.12346, "11.12346"},
		{11.1, "11.1"},
	}
	for _, tc := range cases {
		got := formatCoord(tc.in)
		if got != tc.want {
			t.Errorf("formatCoord(%v) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestParseCoordsMaxDecimals(t *testing.T) {
	const src = `$START p | {} | https://x
>>11.649999999999999,70.01,25.5,3 | {} | 0
hi
$END p`
	doc, err := ParseString(src)
	if err != nil {
		t.Fatal(err)
	}
	t0 := doc.Pages[0].Texts[0]
	if t0.X != 11.65 || t0.Y != 70.01 || t0.Width != 25.5 || t0.TextSize != 3 {
		t.Fatalf("coords: %+v", t0)
	}
	out, err := FormatPSRT(doc, false)
	if err != nil {
		t.Fatal(err)
	}
	if string(out) != `$START p | {} | https://x
>>11.65,70.01,25.5,3 | {} | 0
hi
$END p
` {
		t.Fatalf("formatted:\n%s", string(out))
	}
}

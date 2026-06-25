package svgpath

import "testing"

func TestParse_valid(t *testing.T) {
	cases := []string{
		"M0,0 L10,10 Z",
		"M0,0 L10,10 L20,0 Z",
		"M10,50 C10,25 30,10 50,10 C70,10 90,25 90,50 Z",
		"M0,0 Q50,100 100,0 Z",
		"M0,0 H10 V10 H0 Z",
		"M0,0 S10,10 20,0 Z",
		"M0,0 T10,10 Z",
		// arc with normal separators
		"M10,10 A5,5 0 0,1 20,20 Z",
		// arc with compact, no-separator flags ("0030" => flags 0,0 then x=30)
		"M0,0 A5,5,45,0030,30 Z",
		// scientific notation and signed numbers
		"M-1.5e2,0 L1e-2,3.5 Z",
		// relative commands
		"m0,0 l10,10 z",
	}
	for _, d := range cases {
		if _, err := Parse(d); err != nil {
			t.Errorf("Parse(%q) = %v, want nil", d, err)
		}
	}
}

func TestParse_invalid(t *testing.T) {
	cases := []string{
		"",
		"L10,10 Z",            // doesn't start with M/m
		"M0,0 L10 Z",          // L needs 2 numbers, got 1 before Z
		"M0,0 X10,10 Z",       // unknown command letter
		"M0,0 A5,5 0 2,1 1,1", // large-arc-flag must be 0 or 1
		"M0,0 C1,1 2,2 Z",     // C needs 6 numbers, only 4 given before Z
	}
	for _, d := range cases {
		if _, err := Parse(d); err == nil {
			t.Errorf("Parse(%q) = nil, want error", d)
		}
	}
}

func TestParse_subpathCount(t *testing.T) {
	info, err := Parse("M0,0 L10,10 Z M20,20 L30,30 Z")
	if err != nil {
		t.Fatal(err)
	}
	if info.Subpaths != 2 {
		t.Fatalf("Subpaths: got %d want 2", info.Subpaths)
	}
}

func TestParse_implicitMoveLineto(t *testing.T) {
	// "M10,10 20,20 30,30" is M followed by two implicit linetos, not extra
	// moveto commands — must count as a single subpath.
	info, err := Parse("M10,10 20,20 30,30 Z")
	if err != nil {
		t.Fatalf("implicit lineto after M should be valid: %v", err)
	}
	if info.Subpaths != 1 {
		t.Fatalf("Subpaths: got %d want 1", info.Subpaths)
	}
}

func TestValidate(t *testing.T) {
	if err := Validate("M0,0 L10,10 Z"); err != nil {
		t.Fatalf("Validate valid path: %v", err)
	}
	if err := Validate("not a path"); err == nil {
		t.Fatal("Validate invalid path: want error")
	}
}

func TestSplitCommands(t *testing.T) {
	got := SplitCommands("M10,50 C10,25 30,10 50,10 Z")
	want := []string{"M10,50", "C10,25 30,10 50,10", "Z"}
	if len(got) != len(want) {
		t.Fatalf("SplitCommands: got %v want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("SplitCommands[%d]: got %q want %q", i, got[i], want[i])
		}
	}
}

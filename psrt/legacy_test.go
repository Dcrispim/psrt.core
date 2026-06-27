package psrt

import "testing"

func TestConvertLegacyCoords(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"50-50-80-2", "50,50,80,2"},
		{"22.6-56.11-77-3", "22.6,56.11,77,3"},
		{"0-0-10-10", "0,0,10,10"},
	}
	for _, tc := range cases {
		got, err := ConvertLegacyCoords(tc.in)
		if err != nil {
			t.Fatalf("ConvertLegacyCoords(%q): %v", tc.in, err)
		}
		if got != tc.want {
			t.Errorf("ConvertLegacyCoords(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestConvertLegacyCoordsRejectsWrongPartCount(t *testing.T) {
	cases := []string{
		"1-2-3",
		"1-2-3-4-5",
		"-8.74-1.65-10-20", // a stray leading minus only breaks the part count further; legacy never tolerated this
	}
	for _, in := range cases {
		if _, err := ConvertLegacyCoords(in); err == nil {
			t.Errorf("ConvertLegacyCoords(%q): expected error, got none", in)
		}
	}
}

func TestConvertLegacyCoordsRejectsNonNumeric(t *testing.T) {
	if _, err := ConvertLegacyCoords("a-b-c-d"); err == nil {
		t.Error("expected error for non-numeric segment")
	}
}

func TestConvertLegacyDocument(t *testing.T) {
	const legacy = `$START p | {} | https://x
>>50-50-80-2 | {"color":"#fff"} | 0
hello
==6.58-6.17-22.37-1.89 | {"bg":"#eee9b2"} | 1
~~10-10-20-5 | {"bg":"#fff"} | 2
M10,50 L20,50 Z
$END p`

	converted, err := ConvertLegacyDocument(legacy)
	if err != nil {
		t.Fatalf("ConvertLegacyDocument: %v", err)
	}

	doc, err := ParseString(converted)
	if err != nil {
		t.Fatalf("Parse of converted document: %v\n--- converted ---\n%s", err, converted)
	}
	p := doc.Pages[0]
	if len(p.Texts) != 1 || p.Texts[0].X != 50 || p.Texts[0].Width != 80 {
		t.Fatalf("text: %+v", p.Texts)
	}
	if len(p.Masks) != 1 || p.Masks[0].Width != 22.37 || p.Masks[0].Height != 1.89 {
		t.Fatalf("mask: %+v", p.Masks)
	}
	if len(p.PathMasks) != 1 || p.PathMasks[0].Width != 20 || p.PathMasks[0].Height != 5 {
		t.Fatalf("path mask: %+v", p.PathMasks)
	}

	// Path body content (which legitimately contains commas, e.g. M10,50) must
	// be left untouched by the legacy conversion.
	if p.PathMasks[0].Path == "" {
		t.Fatalf("path mask body lost: %+v", p.PathMasks[0])
	}
}

func TestConvertLegacyDocumentRejectsMalformedHeader(t *testing.T) {
	const legacy = `$START p | {} | https://x
>>50-50-80 | {} | 0
hello
$END p`

	if _, err := ConvertLegacyDocument(legacy); err == nil {
		t.Fatal("expected error for malformed legacy header")
	}
}

// Direct ports of the parser's own legacy fixtures (parser_test.go etc.) used
// to read as hyphen-separated before the comma switch — proves the
// round trip legacy text -> ConvertLegacyDocument -> Parse behaves the same
// as the equivalent comma text parsed directly.
func TestConvertLegacyDocumentMatchesDirectCommaParse(t *testing.T) {
	const legacy = `$START p | {} | u
>>10-10-3-40 | {"color":"black"} | 0
hi
$END p`
	const comma = `$START p | {} | u
>>10,10,3,40 | {"color":"black"} | 0
hi
$END p`

	converted, err := ConvertLegacyDocument(legacy)
	if err != nil {
		t.Fatalf("ConvertLegacyDocument: %v", err)
	}
	gotDoc, err := ParseString(converted)
	if err != nil {
		t.Fatalf("Parse of converted document: %v", err)
	}
	wantDoc, err := ParseString(comma)
	if err != nil {
		t.Fatalf("Parse of comma document: %v", err)
	}
	got, want := gotDoc.Pages[0].Texts[0], wantDoc.Pages[0].Texts[0]
	if got.X != want.X || got.Y != want.Y || got.Width != want.Width || got.TextSize != want.TextSize {
		t.Fatalf("mismatch: got %+v want %+v", got, want)
	}
}

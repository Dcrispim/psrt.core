package psrt

import "testing"

func TestNormalizeTextContent(t *testing.T) {
	tests := []struct {
		in, want string
	}{
		{"  hello  ", "hello"},
		{"\tO Soberano\t\r\n", "O Soberano"},
		{"  line1\n  line2  ", "line1\n  line2"},
		{"", ""},
	}
	for _, tc := range tests {
		if got := NormalizeTextContent(tc.in); got != tc.want {
			t.Fatalf("NormalizeTextContent(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestParse_trimsTextContent(t *testing.T) {
	const src = `$START p | {} | https://example.com/i.png
>>0,0,10,10 | {} | 0
   trimmed text   
$END p`
	doc, err := ParseString(src)
	if err != nil {
		t.Fatal(err)
	}
	if got := doc.Pages[0].Texts[0].Content; got != "trimmed text" {
		t.Fatalf("Content = %q, want %q", got, "trimmed text")
	}
}

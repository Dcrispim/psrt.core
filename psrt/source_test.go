package psrt

import (
	"strings"
	"testing"
)

func TestParseWithSource(t *testing.T) {
	const src = `$START p1 | {} | https://cdn.example/capa.avif
>>10-10-80-3 | {"color":"#fff"} | 0
hello
$END p1

$CONSTS
@ baseURL | https://cdn.example/
$ENDCONSTS

$SOURCE
https://cdn.example/capa.avif | data:image/avif;base64,QUFB
D:\fonts\Foo.woff2 | data:font/woff2;base64,QUFB
$ENDSOURCE
`

	doc, err := ParseString(src)
	if err != nil {
		t.Fatal(err)
	}
	if len(doc.Sources) != 2 {
		t.Fatalf("sources: got %d want 2", len(doc.Sources))
	}
	if doc.Sources["https://cdn.example/capa.avif"] != "data:image/avif;base64,QUFB" {
		t.Fatalf("image source: %q", doc.Sources["https://cdn.example/capa.avif"])
	}

	out, err := FormatPSRT(doc, false)
	if err != nil {
		t.Fatal(err)
	}
	round, err := ParseString(string(out))
	if err != nil {
		t.Fatal(err)
	}
	if len(round.Sources) != 2 {
		t.Fatalf("round-trip sources: %d", len(round.Sources))
	}
}

func TestParseFastAndLoadSource(t *testing.T) {
	const src = `$START p | {} | https://x/img.png
>>1-2-3-4 | {} | 0
hi
$END p

$SOURCE
https://x/img.png | data:image/png;base64,QUJD
$ENDSOURCE
`

	doc, err := ParseFastString(src)
	if err != nil {
		t.Fatal(err)
	}
	if doc.Sources["https://x/img.png"] != "" {
		t.Fatalf("fast parse should not load payload, got %q", doc.Sources["https://x/img.png"])
	}
	data, err := LoadSource(src, "https://x/img.png")
	if err != nil {
		t.Fatal(err)
	}
	if data != "data:image/png;base64,QUJD" {
		t.Fatalf("LoadSource: %q", data)
	}
}

func TestSourceErrors(t *testing.T) {
	cases := []struct {
		name string
		src  string
	}{
		{"duplicate source", `$SOURCE
https://a | data:x;base64,a
https://a | data:x;base64,b
$ENDSOURCE`},
		{"source inside page", `$START p | {} | u
$SOURCE
https://a | data:x;base64,a
$ENDSOURCE
$END p`},
		{"consts after source", `$SOURCE
https://a | data:x;base64,a
$ENDSOURCE
$CONSTS
@x | 1
$ENDCONSTS`},
		{"content after source", `$SOURCE
https://a | data:x;base64,a
$ENDSOURCE
extra line`},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := ParseString(tc.src); err == nil {
				t.Fatal("expected error")
			}
		})
	}
}

func TestParseWithLongSourceLine(t *testing.T) {
	// Exceeds bufio default MaxScanTokenSize (64 KiB).
	payload := strings.Repeat("A", 100*1024)
	src := `$START p | {} | https://x/img.png
>>1-2-3-4 | {} | 0
hi
$END p

$SOURCE
https://x/img.png | data:image/png;base64,` + payload + `
$ENDSOURCE
`
	doc, err := ParseString(src)
	if err != nil {
		t.Fatal(err)
	}
	got := doc.Sources["https://x/img.png"]
	if !strings.HasSuffix(got, payload) {
		t.Fatalf("expected long payload, got len %d", len(got))
	}
}

func TestFormatSourceOrder(t *testing.T) {
	doc := Document{
		Pages: []Page{{
			Name: "p", Style: Style(`{}`), ImageURL: "https://x", Texts: []Text{},
		}},
		Sources: map[string]string{
			"https://b": "data:image/png;base64,b",
			"https://a": "data:image/png;base64,a",
		},
	}
	out, err := FormatPSRT(doc, false)
	if err != nil {
		t.Fatal(err)
	}
	s := string(out)
	idxConsts := strings.Index(s, "$CONSTS")
	idxSource := strings.Index(s, "$SOURCE")
	if idxSource < 0 {
		t.Fatal("missing $SOURCE")
	}
	if idxConsts >= 0 && idxConsts > idxSource {
		t.Fatal("$SOURCE must come after $CONSTS")
	}
}

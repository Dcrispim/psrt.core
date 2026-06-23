package psrt

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestParseFullDocument(t *testing.T) {
	const src = `$START p1 | {"backGround":"#fff"} | https://img.example/a.png
    >>10-10-3-40 | {"color":"black"} | 0
    hello

    >>20-20-5-50 | {} | 1
second block
$END p1

$FONTS
https://fonts/a.woff2
https://fonts/b.woff2
$ENDFONTS

$CONSTS
@shadow | "shadow":"0 0 10px 0 rgba(0, 0, 0, 0.5)"
@note | plain value
$ENDCONSTS
`

	doc, err := ParseString(src)
	if err != nil {
		t.Fatal(err)
	}
	if len(doc.Pages) != 1 {
		t.Fatalf("pages: got %d want 1", len(doc.Pages))
	}
	p := doc.Pages[0]
	if p.Name != "p1" || p.ImageURL != "https://img.example/a.png" {
		t.Fatalf("page fields: %+v", p)
	}
	if len(p.Texts) != 2 {
		t.Fatalf("texts: got %d want 2", len(p.Texts))
	}
	if p.Texts[0].Content != "hello" || p.Texts[0].Index != 0 {
		t.Fatalf("text0: %+v", p.Texts[0])
	}
	if p.Texts[1].Content != "second block" || p.Texts[1].Index != 1 {
		t.Fatalf("text1: %+v", p.Texts[1])
	}
	if len(doc.Fonts) != 2 {
		t.Fatalf("fonts: %v", doc.Fonts)
	}
	if doc.Consts["shadow"] != `"shadow":"0 0 10px 0 rgba(0, 0, 0, 0.5)"` {
		t.Fatalf("const shadow: %q", doc.Consts["shadow"])
	}
	if doc.Consts["note"] != "plain value" {
		t.Fatalf("const note: %q", doc.Consts["note"])
	}

	b, err := ToJSON(doc)
	if err != nil {
		t.Fatal(err)
	}
	var round Document
	if err := json.Unmarshal(b, &round); err != nil {
		t.Fatalf("json round-trip: %v\n%s", err, string(b))
	}
}

func TestJSONStyleIsEmbeddedObjectNotBase64(t *testing.T) {
	const src = `$START p | {"backGround":"#fff"} | https://x
    >>1-2-3-4 | {"color":"red"} | 0
    hi
$END p`
	doc, err := ParseString(src)
	if err != nil {
		t.Fatal(err)
	}
	b, err := ToJSON(doc)
	if err != nil {
		t.Fatal(err)
	}
	out := string(b)
	if strings.Contains(out, `"style": "eyJ`) {
		t.Fatalf("style was base64-encoded; expected embedded JSON object:\n%s", out)
	}
	if !strings.Contains(out, `"backGround"`) || !strings.Contains(out, `"#fff"`) ||
		!strings.Contains(out, `"color"`) || !strings.Contains(out, `"red"`) {
		t.Fatalf("expected raw JSON keys in output:\n%s", out)
	}
}

func TestParseErrors(t *testing.T) {
	cases := []struct {
		name string
		src  string
	}{
		{"nested start", `$START a | {} | u
$START b | {} | u2
$END a`},
		{"end wrong name", `$START a | {} | u
$END b`},
		{"end without page", `$END a`},
		{"content outside page", `hello`},
		{"duplicate const", `$CONSTS
@x | 1
@x | 2
$ENDCONSTS`},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := ParseString(tc.src); err == nil {
				t.Fatal("expected error")
			}
		})
	}
}

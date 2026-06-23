package editor

import (
	"strings"
	"testing"

	"psrt/psrt"
)

func docWithLiteralColor(t *testing.T) psrt.Document {
	t.Helper()
	const src = `$START p1 | {"backGround":"#1DB954"} | https://a.example/1.png
    >>10-10-80-1 | {"color":"#1DB954"} | 0
    Spotify green
$END p1
`
	doc, err := psrt.ParseString(src)
	if err != nil {
		t.Fatal(err)
	}
	return doc
}

func TestAddConstSubstitutes(t *testing.T) {
	doc := docWithLiteralColor(t)
	if err := AddConst(&doc, "accent_spotify", "#1DB954"); err != nil {
		t.Fatal(err)
	}
	p := doc.Pages[0]
	if !strings.Contains(p.ImageURL, "@accent_spotify@") && p.ImageURL == "https://a.example/1.png" {
		// URL unchanged is ok
	}
	if !strings.Contains(string(p.Style), "@accent_spotify@") {
		t.Fatalf("page style: %s", p.Style)
	}
	if !strings.Contains(string(p.Texts[0].Style), "@accent_spotify@") {
		t.Fatalf("text style: %s", p.Texts[0].Style)
	}
	if !strings.Contains(p.Texts[0].Content, "Spotify") {
		t.Fatalf("content unchanged: %q", p.Texts[0].Content)
	}
}

func TestRemoveConstReverts(t *testing.T) {
	doc := docWithLiteralColor(t)
	if err := AddConst(&doc, "accent", "#1DB954"); err != nil {
		t.Fatal(err)
	}
	if err := RemoveConst(&doc, "accent"); err != nil {
		t.Fatal(err)
	}
	if _, ok := doc.Consts["accent"]; ok {
		t.Fatal("const should be removed")
	}
	p := doc.Pages[0]
	if strings.Contains(string(p.Style), "@accent@") {
		t.Fatalf("style still has placeholder: %s", p.Style)
	}
	if !strings.Contains(string(p.Style), "#1DB954") {
		t.Fatalf("style not reverted: %s", p.Style)
	}
}

func TestConstRoundTrip(t *testing.T) {
	doc := docWithLiteralColor(t)
	if err := AddConst(&doc, "accent", "#1DB954"); err != nil {
		t.Fatal(err)
	}
	out, err := psrt.FormatPSRT(doc, false)
	if err != nil {
		t.Fatal(err)
	}
	doc2, err := psrt.ParseString(string(out))
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := doc2.Consts["accent"]; !ok {
		t.Fatal("const missing after round trip")
	}
}

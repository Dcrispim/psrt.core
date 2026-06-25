package editor

import (
	"strings"
	"testing"

	"github.com/Dcrispim/psrt.core/psrt"
)

func sampleDoc(t *testing.T) psrt.Document {
	t.Helper()
	const src = `$START p1 | {"backGround":"#111"} | https://a.example/1.png
    >>10-10-80-1 | {"color":"#fff"} | 0
    hello
    >>20-20-80-1 | {"color":"#aaa"} | 1
    world
$END p1

$START p2 | {} | https://a.example/2.png
$END p2
`
	doc, err := psrt.ParseString(src)
	if err != nil {
		t.Fatal(err)
	}
	return doc
}

func TestRenamePage(t *testing.T) {
	doc := sampleDoc(t)
	if err := RenamePage(&doc, "p1", "cover"); err != nil {
		t.Fatal(err)
	}
	if _, err := FindPage(&doc, "cover"); err != nil {
		t.Fatal(err)
	}
}

func TestSetPagePath(t *testing.T) {
	doc := sampleDoc(t)
	if err := SetPagePath(&doc, "p1", "https://new.example/x.png"); err != nil {
		t.Fatal(err)
	}
	p, _ := FindPage(&doc, "p1")
	if p.ImageURL != "https://new.example/x.png" {
		t.Fatalf("url: %q", p.ImageURL)
	}
}

func TestSetPageStyleKeyAndMerge(t *testing.T) {
	doc := sampleDoc(t)
	if err := SetPageStyle(&doc, "p1", "backGround", `"#000"`, nil); err != nil {
		t.Fatal(err)
	}
	if err := SetPageStyle(&doc, "p1", "", "", []byte(`{"opacity":1}`)); err != nil {
		t.Fatal(err)
	}
	p, _ := FindPage(&doc, "p1")
	if !strings.Contains(string(p.Style), "#000") {
		t.Fatalf("style: %s", p.Style)
	}
}

func TestMovePage(t *testing.T) {
	doc := sampleDoc(t)
	if err := MovePage(&doc, "p2", "p1", ""); err != nil {
		t.Fatal(err)
	}
	if doc.Pages[0].Name != "p2" {
		t.Fatalf("order: %s, %s", doc.Pages[0].Name, doc.Pages[1].Name)
	}
}

func TestPageRoundTrip(t *testing.T) {
	doc := sampleDoc(t)
	if err := RenamePage(&doc, "p1", "renamed"); err != nil {
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
	if doc2.Pages[0].Name != "renamed" {
		t.Fatalf("round trip name: %q", doc2.Pages[0].Name)
	}
}

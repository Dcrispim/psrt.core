package psrt

import (
	"strings"
	"testing"
)

func TestParseMaskBlock(t *testing.T) {
	const src = `$START p | {} | https://x
==6.58-6.17-22.37-1.89 | {"bg":"#eee9b2"} | 0
>>11.94-8.36-13.33-3 | {"bg":"#eeeade"} | 1
hello
$END p`
	doc, err := ParseString(src)
	if err != nil {
		t.Fatal(err)
	}
	p := doc.Pages[0]
	if len(p.Masks) != 1 {
		t.Fatalf("masks: got %d want 1", len(p.Masks))
	}
	m := p.Masks[0]
	if m.Height != 1.89 || m.Width != 22.37 || m.Index != 0 {
		t.Fatalf("mask: %+v", m)
	}
	if len(p.Texts) != 1 || p.Texts[0].Content != "hello" {
		t.Fatalf("texts: %+v", p.Texts)
	}
}

func TestMaskRoundTrip(t *testing.T) {
	const src = `$START p | {} | u
==10-10-20-5 | {"bg":"#fff"} | 0
$END p`
	doc, err := ParseString(src)
	if err != nil {
		t.Fatal(err)
	}
	out, err := FormatPSRT(doc, false)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(out), "==10-10-20-5") {
		t.Fatalf("format missing mask header:\n%s", out)
	}
	doc2, err := ParseString(string(out))
	if err != nil {
		t.Fatal(err)
	}
	m := doc2.Pages[0].Masks[0]
	if m.Height != 5 || m.Width != 20 {
		t.Fatalf("round-trip mask: %+v", m)
	}
}

func TestPromoteEmptyTextsToMasks(t *testing.T) {
	doc := Document{
		Pages: []Page{{
			Name: "p", ImageURL: "u", Style: Style("{}"),
			Texts: []Text{{
				BaseBlock: BaseBlock{X: 1, Y: 2, Width: 3, Style: Style(`{"height":"2.5%","bg":"#000"}`), Index: 0},
				TextSize: 4,
				Content:  "",
			}},
		}},
	}
	PromoteEmptyTextsToMasks(&doc)
	if len(doc.Pages[0].Texts) != 0 {
		t.Fatalf("texts: %+v", doc.Pages[0].Texts)
	}
	if len(doc.Pages[0].Masks) != 1 {
		t.Fatal("expected one mask")
	}
	m := doc.Pages[0].Masks[0]
	if m.Height != 2.5 {
		t.Fatalf("height: %v", m.Height)
	}
	if strings.Contains(string(m.Style), "height") {
		t.Fatalf("height should be removed from style: %s", m.Style)
	}
	data, err := FormatPSRT(doc, false)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(strings.TrimSpace(strings.Split(string(data), "\n")[1]), "==") {
		t.Fatalf("expected == line:\n%s", data)
	}
}

func TestPromoteEmptyTextUsesTextSizeWhenNoStyleHeight(t *testing.T) {
	doc := Document{
		Pages: []Page{{
			Name: "p", ImageURL: "u", Style: Style("{}"),
			Texts: []Text{{
				BaseBlock: BaseBlock{X: 1, Y: 2, Width: 20, Style: Style(`{"bg":"#eee"}`), Index: 3},
				TextSize: 2.64,
				Content:  "",
			}},
		}},
	}
	PromoteEmptyTextsToMasks(&doc)
	m := doc.Pages[0].Masks[0]
	if m.Height != 2.64 {
		t.Fatalf("height from textSize: got %v want 2.64", m.Height)
	}
}

func TestParseMaskCoordsError(t *testing.T) {
	const src = `$START p | {} | u
==1-2-3 | {} | 0
$END p`
	if _, err := ParseString(src); err == nil {
		t.Fatal("expected coord error")
	}
}

package psrt

import (
	"encoding/json"
	"strings"
	"testing"
)

const balloonPath = `M10,50
C10,25 30,10 50,10
C70,10 90,25 90,50
C90,72 75,85 55,82
L48,95
L45,80
C25,77 10,65 10,50
Z`

func TestParsePathMaskBlock(t *testing.T) {
	src := `$START p | {} | https://x
~~6.58,6.17,22.37,8.4 | {"bg":"#eee9b2"} | 0
` + balloonPath + `
>>11.94,8.36,13.33,3 | {"bg":"#eeeade"} | 1
hello
$END p`
	doc, err := ParseString(src)
	if err != nil {
		t.Fatal(err)
	}
	p := doc.Pages[0]
	if len(p.PathMasks) != 1 {
		t.Fatalf("path masks: got %d want 1", len(p.PathMasks))
	}
	pm := p.PathMasks[0]
	if pm.Height != 8.4 || pm.Width != 22.37 || pm.Index != 0 {
		t.Fatalf("path mask: %+v", pm)
	}
	if strings.Contains(pm.Path, "\n") {
		t.Fatalf("path must be normalized to a single line: %q", pm.Path)
	}
	if !strings.HasPrefix(pm.Path, "M10,50 C10,25") {
		t.Fatalf("unexpected normalized path: %q", pm.Path)
	}
	if len(p.Texts) != 1 || p.Texts[0].Content != "hello" {
		t.Fatalf("texts: %+v", p.Texts)
	}
}

func TestPathMaskRoundTrip(t *testing.T) {
	src := `$START p | {} | u
~~10,10,20,5 | {"bg":"#fff"} | 0
` + balloonPath + `
$END p`
	doc, err := ParseString(src)
	if err != nil {
		t.Fatal(err)
	}
	out, err := FormatPSRT(doc, false)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(out), "~~10,10,20,5") {
		t.Fatalf("format missing path mask header:\n%s", out)
	}
	doc2, err := ParseString(string(out))
	if err != nil {
		t.Fatalf("reparse: %v\n%s", err, out)
	}
	pm := doc2.Pages[0].PathMasks[0]
	if pm.Height != 5 || pm.Width != 20 {
		t.Fatalf("round-trip path mask: %+v", pm)
	}
	if pm.Path != doc.Pages[0].PathMasks[0].Path {
		t.Fatalf("round-trip path mismatch:\nbefore=%q\nafter=%q", doc.Pages[0].PathMasks[0].Path, pm.Path)
	}
}

func TestPathMaskJSONRoundTrip(t *testing.T) {
	src := `$START p | {} | u
~~10,10,20,5 | {"bg":"#fff"} | 0
` + balloonPath + `
$END p`
	doc, err := ParseString(src)
	if err != nil {
		t.Fatal(err)
	}
	raw, err := ToJSON(doc)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(raw), `"pathMasks"`) {
		t.Fatalf("json missing pathMasks field:\n%s", raw)
	}
	var doc2 Document
	if err := json.Unmarshal(raw, &doc2); err != nil {
		t.Fatal(err)
	}
	if len(doc2.Pages[0].PathMasks) != 1 {
		t.Fatalf("path masks after json round-trip: %+v", doc2.Pages[0].PathMasks)
	}
	if doc2.Pages[0].PathMasks[0].Path != doc.Pages[0].PathMasks[0].Path {
		t.Fatalf("path mismatch after json round-trip: %q vs %q",
			doc2.Pages[0].PathMasks[0].Path, doc.Pages[0].PathMasks[0].Path)
	}
}

func TestParsePathMaskCoordsError(t *testing.T) {
	src := `$START p | {} | u
~~1,2,3 | {} | 0
` + balloonPath + `
$END p`
	if _, err := ParseString(src); err == nil {
		t.Fatal("expected coord error")
	}
}

func TestParsePathMaskEmptyBodyError(t *testing.T) {
	src := `$START p | {} | u
~~10,10,20,5 | {} | 0
$END p`
	_, err := ParseString(src)
	if err == nil {
		t.Fatal("expected empty body error")
	}
	if !strings.Contains(err.Error(), "path mask body is empty") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParsePathMaskInvalidStyleJSONError(t *testing.T) {
	src := `$START p | {} | u
~~10,10,20,5 | not-json | 0
` + balloonPath + `
$END p`
	if _, err := ParseString(src); err == nil {
		t.Fatal("expected style JSON error")
	}
}

func TestParsePathMaskDuplicateIndexAgainstText(t *testing.T) {
	src := `$START p | {} | u
>>10,10,20,3 | {} | 0
hi
~~10,10,20,5 | {} | 0
` + balloonPath + `
$END p`
	_, err := ParseString(src)
	if err == nil || !strings.Contains(err.Error(), "duplicate index") {
		t.Fatalf("expected duplicate index error, got %v", err)
	}
}

func TestParsePathMaskDuplicateIndexAgainstMask(t *testing.T) {
	src := `$START p | {} | u
==10,10,20,5 | {} | 0
~~10,10,20,5 | {} | 0
` + balloonPath + `
$END p`
	_, err := ParseString(src)
	if err == nil || !strings.Contains(err.Error(), "duplicate index") {
		t.Fatalf("expected duplicate index error, got %v", err)
	}
}

func TestParsePathMaskDuplicateIndexAgainstPathMask(t *testing.T) {
	src := `$START p | {} | u
~~10,10,20,5 | {} | 0
` + balloonPath + `
~~30,30,20,5 | {} | 0
` + balloonPath + `
$END p`
	_, err := ParseString(src)
	if err == nil || !strings.Contains(err.Error(), "duplicate index") {
		t.Fatalf("expected duplicate index error, got %v", err)
	}
}

func TestParsePathMaskInvalidSyntaxError(t *testing.T) {
	src := `$START p | {} | u
~~10,10,20,5 | {} | 0
not a path at all !!
$END p`
	_, err := ParseString(src)
	if err == nil || !strings.Contains(err.Error(), "invalid svg path data") {
		t.Fatalf("expected invalid svg path data error, got %v", err)
	}
}

func TestParsePathMaskMultipleSubpathsError(t *testing.T) {
	src := `$START p | {} | u
~~10,10,20,5 | {} | 0
M0,0 L10,10 Z M20,20 L30,30 Z
$END p`
	_, err := ParseString(src)
	if err == nil || !strings.Contains(err.Error(), "single shape") {
		t.Fatalf("expected single-shape error, got %v", err)
	}
}

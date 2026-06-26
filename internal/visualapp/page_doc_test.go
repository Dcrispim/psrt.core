package visualapp

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Dcrispim/psrt.core/psrt"
)

func TestExtractPageDocument(t *testing.T) {
	doc := psrt.Document{
		Pages: []psrt.Page{
			{Name: "capa", ImageURL: "https://example.com/a.jpg"},
			{Name: "intro", ImageURL: "https://example.com/b.jpg"},
		},
		Fonts:  []string{"https://fonts.example/f.woff2"},
		Consts: map[string]string{"x": "1"},
	}
	sub, err := extractPageDocument(doc, "intro")
	if err != nil {
		t.Fatal(err)
	}
	if len(sub.Pages) != 1 || sub.Pages[0].Name != "intro" {
		t.Fatalf("pages: %+v", sub.Pages)
	}
	if len(sub.Fonts) != 1 || sub.Consts["x"] != "1" {
		t.Fatalf("fonts/consts not preserved")
	}
}

func TestFormatPageDocumentJSON(t *testing.T) {
	doc := psrt.Document{
		Pages: []psrt.Page{
			{Name: "capa", ImageURL: "https://example.com/a.jpg", Style: []byte(`{}`)},
			{Name: "intro", ImageURL: "https://example.com/b.jpg", Style: []byte(`{}`)},
		},
		Fonts:  []string{"https://fonts.example/f.woff2"},
		Consts: map[string]string{"accent": "#1DB954"},
	}
	raw, err := psrt.ToJSON(doc)
	if err != nil {
		t.Fatal(err)
	}
	text, err := FormatPageDocumentJSON(string(raw), "intro")
	if err != nil {
		t.Fatal(err)
	}
	if !contains(text, "$START intro") || !contains(text, "$FONTS") || !contains(text, "$CONSTS") {
		t.Fatalf("unexpected fragment:\n%s", text)
	}
	if contains(text, "$START capa") {
		t.Fatal("other pages must not appear")
	}
}

func TestMergePageDocumentPSRT(t *testing.T) {
	full := psrt.Document{
		Pages: []psrt.Page{
			{Name: "capa", ImageURL: "https://example.com/a.jpg", Style: []byte(`{}`)},
			{
				Name: "intro", ImageURL: "https://example.com/b.jpg", Style: []byte(`{}`),
				Texts: []psrt.Text{{Index: 0, X: 10, Y: 20, Width: 80, TextSize: 2, Style: []byte(`{}`), Content: "old"}},
			},
		},
		Fonts:  []string{"https://old/font.woff2"},
		Consts: map[string]string{"x": "1"},
	}
	raw, err := psrt.ToJSON(full)
	if err != nil {
		t.Fatal(err)
	}
	fragment := `$START intro | {} | https://example.com/b.jpg
>>10-20-80-2 | {} | 0
new line
$END intro
$FONTS
https://new/font.woff2
$ENDFONTS
$CONSTS
@ y | 2
$ENDCONSTS
`
	mergedJSON, err := MergePageDocumentPSRT(string(raw), "intro", fragment)
	if err != nil {
		t.Fatal(err)
	}
	merged, err := psrt.ParseJSON([]byte(mergedJSON))
	if err != nil {
		t.Fatal(err)
	}
	if len(merged.Pages) != 2 {
		t.Fatalf("pages: %d", len(merged.Pages))
	}
	intro := merged.Pages[1]
	if intro.Texts[0].Content != "new line" {
		t.Fatalf("content %q", intro.Texts[0].Content)
	}
	if len(merged.Fonts) != 1 || merged.Fonts[0] != "https://new/font.woff2" {
		t.Fatalf("fonts %+v", merged.Fonts)
	}
	if merged.Consts["y"] != "2" {
		t.Fatalf("consts %+v", merged.Consts)
	}
	if merged.Pages[0].Name != "capa" {
		t.Fatal("capa page must remain")
	}
}

func TestWriteIntermediatePagePSRT(t *testing.T) {
	dir := t.TempDir()
	doc := psrt.Document{
		Pages: []psrt.Page{{Name: "mood-sexta", ImageURL: "https://example.com/x.jpg"}},
	}
	path, err := writeIntermediatePagePSRT(dir, "mood-sexta", doc)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatal(err)
	}
	if filepath.Base(path) != "mood-sexta.psrt" {
		t.Fatalf("path %q", path)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(data) == 0 || !contains(string(data), "mood-sexta") {
		t.Fatalf("unexpected content: %q", string(data))
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(sub) == 0 || indexOf(s, sub) >= 0)
}

func indexOf(s, sub string) int {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}

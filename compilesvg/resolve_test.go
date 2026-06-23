package compilesvg

import (
	"testing"

	"psrt/psrt"
)

func TestResolveDocument_expandsStyleAndContent(t *testing.T) {
	doc := psrt.Document{
		Consts: map[string]string{
			"accent": "#1DB954",
			"shadow": `"textShadow":"0 1px 2px rgba(0,0,0,0.5)"`,
		},
		Pages: []psrt.Page{
			{
				Name:     "p1",
				ImageURL: "https://example.com/@accent@.png",
				Style:    psrt.Style(`{"color":"@accent@"}`),
				Texts: []psrt.Text{
					psrt.NewText(0, 0, 0, 0, 0, psrt.Style(`{"color":"@accent@",@shadow@}`), "Hello @accent@", ""),
				},
			},
		},
	}
	out := ResolveDocument(doc)
	if out.Pages[0].ImageURL != "https://example.com/#1DB954.png" {
		t.Fatalf("ImageURL = %q", out.Pages[0].ImageURL)
	}
	if string(out.Pages[0].Texts[0].Style) == "" {
		t.Fatal("expected expanded text style")
	}
	if out.Pages[0].Texts[0].Content != "Hello #1DB954" {
		t.Fatalf("Content = %q", out.Pages[0].Texts[0].Content)
	}
}

package compilehtml

import (
	"strings"
	"testing"

	"github.com/Dcrispim/psrt.core/compileasset"
	"github.com/Dcrispim/psrt.core/psrt"
)

// minimal 1×1 PNG
var tinyPNG = []byte{
	0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a,
	0x00, 0x00, 0x00, 0x0d, 0x49, 0x48, 0x44, 0x52,
	0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
	0x08, 0x06, 0x00, 0x00, 0x00, 0x1f, 0x15, 0xc4,
	0x89, 0x00, 0x00, 0x00, 0x0a, 0x49, 0x44, 0x41,
	0x54, 0x78, 0x9c, 0x63, 0x00, 0x01, 0x00, 0x00,
	0x05, 0x00, 0x01, 0x0d, 0x0a, 0x2d, 0xb4, 0x00,
	0x00, 0x00, 0x00, 0x49, 0x45, 0x4e, 0x44, 0xae,
	0x42, 0x60, 0x82,
}

func TestRenderHTMLTextVisibleAndStacked(t *testing.T) {
	doc := psrt.Document{
		Pages: []psrt.Page{
			{
				Name:     "a",
				Style:    psrt.Style(`{"backGround":"#000"}`),
				ImageURL: "https://example.com/a.png",
				Texts: []psrt.Text{
					psrt.NewText(10, 20, 80, 12, 0, psrt.Style(`{"color":"#fff"}`), "Hello", ""),
				},
			},
			{
				Name:     "b",
				Style:    psrt.Style(`{}`),
				ImageURL: "https://example.com/b.png",
			},
		},
	}
	assets := map[string]compileasset.Asset{
		"https://example.com/a.png": {MIME: "image/png", Bytes: tinyPNG},
		"https://example.com/b.png": {MIME: "image/png", Bytes: tinyPNG},
	}
	out, err := RenderHTML(doc, assets)
	if err != nil {
		t.Fatal(err)
	}
	html := string(out)
	if !strings.Contains(html, "flex-direction:column") {
		t.Fatal("slides-wrap should stack pages vertically")
	}
	if rule := cssRule(html, ".slide{"); strings.Contains(rule, "container-type") {
		t.Fatal("size containment must not be on .slide (collapses slide height)")
	}
	if !strings.Contains(html, ".slide-overlay{") || !strings.Contains(html, "container-type:size;container-name:slide") {
		t.Fatal("size containment should be on .slide-overlay")
	}
	if !strings.Contains(html, "position:absolute") {
		t.Fatal("text layers must be absolutely positioned")
	}
	if !strings.Contains(html, "font-size:12cqmin") {
		t.Fatalf("font-size should scale with slide container (cqmin), got:\n%s", snippet(html, "font-size"))
	}
	if !strings.Contains(html, `class="slide" style="width:1px`) {
		t.Fatal("slide should use natural image width in px, not width:100%")
	}
	if !strings.Contains(html, ">Hello<") {
		t.Fatal("text content missing from output")
	}
	if strings.Count(html, `class="slide"`) != 2 {
		t.Fatalf("expected 2 slides, got %d", strings.Count(html, `class="slide"`))
	}
}

func TestRenderHTML_textAlignRightAndJustify(t *testing.T) {
	doc := psrt.Document{
		Pages: []psrt.Page{
			{
				Name:     "p",
				ImageURL: "https://example.com/bg.png",
				Texts: []psrt.Text{
					psrt.NewText(5, 5, 90, 5, 0, psrt.Style(`{"color":"#fff","text-align":"right"}`), "Right", ""),
					psrt.NewText(5, 20, 90, 5, 1, psrt.Style(`{"color":"#fff","text-align":"justify"}`), "Justify me now", ""),
				},
			},
		},
	}
	assets := map[string]compileasset.Asset{
		"https://example.com/bg.png": {MIME: "image/png", Bytes: tinyPNG},
	}
	out, err := RenderHTML(doc, assets)
	if err != nil {
		t.Fatal(err)
	}
	html := string(out)
	if !strings.Contains(html, "text-align:right") {
		t.Fatalf("expected text-align:right on text layer, got:\n%s", snippet(html, "text-layer"))
	}
	if !strings.Contains(html, "text-align:justify") {
		t.Fatalf("expected text-align:justify on text layer, got:\n%s", snippet(html, "text-layer"))
	}
}

func TestRenderHTML_shortLabelNoComputedHeight(t *testing.T) {
	doc := psrt.Document{
		Pages: []psrt.Page{
			{
				Name:     "p",
				ImageURL: "https://example.com/bg.png",
				Texts: []psrt.Text{
					psrt.NewText(46.19, 38.69, 6.22, 1.32, 13, psrt.Style(`{"text-align":"center","color":"#000000","font-weight":"700","background":"#ff0000"}`), "Perdidos!", ""),
				},
			},
		},
	}
	assets := map[string]compileasset.Asset{
		"https://example.com/bg.png": {MIME: "image/png", Bytes: tinyPNG},
	}
	out, err := RenderHTML(doc, assets)
	if err != nil {
		t.Fatal(err)
	}
	html := string(out)
	if h := extractStyleProp(html, "text-layer", "height"); h != "" {
		t.Fatalf("short label without padding must not get computed height, got height:%s\n%s", h, snippet(html, "text-layer"))
	}
	if !strings.Contains(html, ">Perdidos!<") {
		t.Fatal("text content missing")
	}
}

func TestRenderHTML_textBoxBackground(t *testing.T) {
	doc := psrt.Document{
		Pages: []psrt.Page{
			{
				Name:     "intro",
				ImageURL: "https://example.com/bg.png",
				Texts: []psrt.Text{
					psrt.NewText(10, 20, 80, 5, 1, psrt.Style(`{"background":"#000000ff","color":"#fff","text-align":"center"}`), "Title", ""),
				},
			},
		},
	}
	assets := map[string]compileasset.Asset{
		"https://example.com/bg.png": {MIME: "image/png", Bytes: tinyPNG},
	}
	out, err := RenderHTML(doc, assets)
	if err != nil {
		t.Fatal(err)
	}
	html := string(out)
	if !strings.Contains(html, "background-color:#000000ff") {
		t.Fatalf("text box background missing:\n%s", snippet(html, "text-layer"))
	}
}

func TestRenderHTML_pageBackgroundColor(t *testing.T) {
	doc := psrt.Document{
		Pages: []psrt.Page{
			{
				Name:     "capa",
				Style:    psrt.Style(`{"backGround":"#1C1C26"}`),
				ImageURL: "https://example.com/bg.png",
			},
		},
	}
	assets := map[string]compileasset.Asset{
		"https://example.com/bg.png": {MIME: "image/png", Bytes: tinyPNG},
	}
	out, err := RenderHTML(doc, assets)
	if err != nil {
		t.Fatal(err)
	}
	html := string(out)
	if !strings.Contains(html, "background:#1C1C26") {
		t.Fatalf("page background missing:\n%s", cssRule(html, `class="slide"`))
	}
}

func cssRule(s, selector string) string {
	i := strings.Index(s, selector)
	if i < 0 {
		return ""
	}
	j := strings.Index(s[i:], "}")
	if j < 0 {
		return s[i:]
	}
	return s[i : i+j+1]
}

func snippet(s, sub string) string {
	i := strings.Index(s, sub)
	if i < 0 {
		return ""
	}
	end := i + 80
	if end > len(s) {
		end = len(s)
	}
	return s[i:end]
}

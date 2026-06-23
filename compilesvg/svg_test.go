package compilesvg

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"psrt/compileasset"
	"psrt/compilesvg/textoutline"
	"psrt/psrt"
)

// minimal 1x1 PNG
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

func TestMain(m *testing.M) {
	code := m.Run()
	textoutline.CloseDefault()
	os.Exit(code)
}

func TestRenderPageSVG_structure(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/png")
		_, _ = w.Write(tinyPNG)
	}))
	defer srv.Close()

	url := srv.URL + "/bg.png"
	assets := map[string]compileasset.Asset{
		url: {Bytes: tinyPNG, MIME: "image/png"},
	}
	page := &psrt.Page{
		Name:     "capa",
		ImageURL: url,
		Style:    psrt.Style(`{"backGround":"#0F0F14","color":"#FFFFFF"}`),
		Texts: []psrt.Text{
			psrt.NewText(10, 20, 80, 12, 0, psrt.Style(`{"color":"#A1A1AA"}`), "Line1\nLine2", ""),
		},
	}
	out, err := RenderPageSVG(page, "capa", nil, assets)
	if err != nil {
		t.Fatal(err)
	}
	s := string(out)
	checks := []string{
		`.psrt-page-capa{`,
		`id="psrt-page-capa"`,
		`id="psrt-text-capa-0-glyphs"`,
		`<path d="`,
		`id="psrt-text-capa-0-a11y" display="none"`,
		`data:image/png;base64,`,
		`Line1`,
	}
	for _, c := range checks {
		if !strings.Contains(s, c) {
			t.Fatalf("missing %q in output", c)
		}
	}
	if strings.Contains(s, `style="`) {
		t.Fatal("SVG must not use inline style attributes on text blocks")
	}
	if strings.Contains(s, `foreignObject`) {
		t.Fatal("compiled SVG must not contain foreignObject text")
	}
	if strings.Contains(s, `.psrt-text-capa-0{`) {
		t.Fatal("text block CSS must not be embedded in compiled SVG")
	}
	if strings.Contains(s, `@font-face`) {
		t.Fatal("font-face must not be embedded in compiled SVG")
	}
}

func canvasPNG(w, h int) []byte {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		panic(err)
	}
	return buf.Bytes()
}

func TestRenderPageSVG_inlineStrongOutlined(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/png")
		_, _ = w.Write(tinyPNG)
	}))
	defer srv.Close()

	url := srv.URL + "/bg.png"
	assets := map[string]compileasset.Asset{url: {Bytes: tinyPNG, MIME: "image/png"}}
	page := &psrt.Page{
		Name:     "capa",
		ImageURL: url,
		Texts: []psrt.Text{
			psrt.NewText(10, 20, 80, 5, 0, psrt.Style(`{"text-align":"center"}`), "estamos **Perdidos!** Metade das", ""),
		},
	}
	out, err := RenderPageSVG(page, "capa", nil, assets)
	if err != nil {
		t.Fatal(err)
	}
	s := string(out)
	if !strings.Contains(s, `id="psrt-text-capa-0-glyphs"`) {
		t.Fatalf("expected glyph group in:\n%s", s)
	}
	pathCount := strings.Count(s, "<path ")
	if pathCount < 2 {
		t.Fatalf("expected multiple paths for inline runs, got %d paths", pathCount)
	}
	if strings.Contains(s, `foreignObject`) {
		t.Fatal("must not emit foreignObject")
	}
	if !strings.Contains(s, `display="none">estamos Perdidos! Metade das`) {
		t.Fatal("missing a11y plain text")
	}
}

func TestRenderPageSVG_textBackgroundRectMatchesGeometry(t *testing.T) {
	const cw, ch = 1080, 1920
	pngBytes := canvasPNG(cw, ch)
	text := psrt.NewText(41, 69, 23, 3, 3, psrt.Style(`{"color":"#ffffff","background":"#ff0000","padding":"8px"}`), "Revisor", "")
	page := &psrt.Page{Name: "intro", ImageURL: "u", Texts: []psrt.Text{text}}
	assets := map[string]compileasset.Asset{
		"u": {Bytes: pngBytes, MIME: "image/png"},
	}
	out, err := RenderPageSVG(page, "intro", nil, assets)
	if err != nil {
		t.Fatal(err)
	}
	s := string(out)
	content := psrt.NormalizeTextContent(text.Content)
	x, y, w, h := TextBlockGeometry(&text, content, cw, ch)

	rectRe := regexp.MustCompile(`<rect[^>]*fill="#ff0000"[^>]*/>`)
	rectM := rectRe.FindString(s)
	if rectM == "" {
		t.Fatalf("missing background rect:\n%s", s)
	}
	for _, attr := range []string{"x", "y", "width", "height"} {
		want := map[string]string{
			"x": fmt.Sprintf("%d", x), "y": fmt.Sprintf("%d", y),
			"width": fmt.Sprintf("%d", w), "height": fmt.Sprintf("%d", h),
		}[attr]
		if got := attrVal(rectM, attr); got != want {
			t.Fatalf("background rect %s=%q want %q", attr, got, want)
		}
	}
	if !strings.Contains(s, `id="psrt-text-intro-3-glyphs"`) {
		t.Fatal("missing glyph paths group")
	}
}

func attrVal(tag, name string) string {
	re := regexp.MustCompile(name + `="([^"]*)"`)
	m := re.FindStringSubmatch(tag)
	if len(m) < 2 {
		return ""
	}
	return m[1]
}

func TestRenderPageSVG_textBackgroundBox(t *testing.T) {
	url := "https://example.com/bg.png"
	assets := map[string]compileasset.Asset{
		url: {Bytes: tinyPNG, MIME: "image/png"},
	}
	text := psrt.NewText(11.6, 55.11, 77.5, 4.9, 1, psrt.Style(`{"background":"#000000ff","text-align":"center"}`), "Cover text", "")
	page := &psrt.Page{
		Name:     "intro",
		ImageURL: url,
		Texts:    []psrt.Text{text},
	}
	out, err := RenderPageSVG(page, "intro", nil, assets)
	if err != nil {
		t.Fatal(err)
	}
	s := string(out)

	slug := "intro"
	canvasW, canvasH := ImageDimensions(tinyPNG, "image/png")
	content := psrt.NormalizeTextContent(text.Content)
	x, y, w, h := TextBlockGeometry(&text, content, canvasW, canvasH)
	bg := compileasset.ParseTextBox(text.Style).Background

	want := []string{
		fmt.Sprintf(`id="%s"`, TextWrapID(slug, text.Index)),
		fmt.Sprintf(`id="%s"`, TextGlyphsID(slug, text.Index)),
		fmt.Sprintf(`x="%d"`, x),
		fmt.Sprintf(`y="%d"`, y),
		fmt.Sprintf(`width="%d"`, w),
		fmt.Sprintf(`height="%d"`, h),
		fmt.Sprintf(`fill="%s"`, bg),
		`<rect`,
		`<path`,
	}
	for _, fragment := range want {
		if !strings.Contains(s, fragment) {
			t.Fatalf("missing %q in output (geometry x=%d y=%d w=%d h=%d)", fragment, x, y, w, h)
		}
	}
	if strings.Contains(s, TextBgID(slug, text.Index)) {
		t.Fatal("must not emit separate SVG background rect")
	}
}

func TestCompile_writesDirectory(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/png")
		_, _ = w.Write(tinyPNG)
	}))
	defer srv.Close()

	doc := psrt.Document{
		Pages: []psrt.Page{
			{
				Name:     "capa",
				ImageURL: srv.URL + "/a.png",
				Texts:    []psrt.Text{psrt.NewText(0, 0, 50, 10, 0, psrt.Style("{}"), "Hi", "")},
			},
			{
				Name:     "page-two",
				ImageURL: srv.URL + "/b.png",
				Texts:    []psrt.Text{psrt.NewText(0, 0, 50, 10, 0, psrt.Style("{}"), "Bye", "")},
			},
		},
	}
	dir := t.TempDir()
	if _, err := Compile(doc, srv.Client(), dir); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"capa.svg", "page-two.svg"} {
		p := filepath.Join(dir, name)
		if _, err := os.Stat(p); err != nil {
			t.Fatalf("missing %s: %v", name, err)
		}
		data, _ := os.ReadFile(p)
		body := string(data)
		if strings.Contains(body, `href="http`) || strings.Contains(body, `src="http`) {
			t.Fatalf("%s should not contain external asset URLs", name)
		}
		if strings.Contains(body, `foreignObject`) {
			t.Fatalf("%s must not contain foreignObject", name)
		}
		if !strings.Contains(body, `<path `) {
			t.Fatalf("%s must contain outlined text paths", name)
		}
	}
}

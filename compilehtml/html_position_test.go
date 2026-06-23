package compilehtml

import (
	"encoding/binary"
	"strings"
	"testing"

	"psrt/compileasset"
	"psrt/psrt"
)

func TestRenderHTML_tallWebtoon_percentPaddingAndSlideWidth(t *testing.T) {
	body := buildWebPVP8X(700, 6556)
	doc := psrt.Document{
		Pages: []psrt.Page{
			{
				Name:     "page",
				ImageURL: "file://test.webp",
				Texts: []psrt.Text{
					{
						BaseBlock: psrt.BaseBlock{
							X: 32.714, Y: 1.427, Width: 55.786, Index: 0,
							Style: psrt.Style(`{"color":"#000000","background":"#ff0000","text-align":"center","border-radius":"100px","padding":"0.6%","font-weight":"700","font-family":"Roboto"}`),
						},
						TextSize: 4.25,
						Content:  "Mesmo se a gente matar você agora, o Mestre de seita não pode nos tocar.",
					},
				},
			},
		},
	}
	assets := map[string]compileasset.Asset{
		"file://test.webp": {MIME: "image/webp", Bytes: body},
	}
	out, err := RenderHTML(doc, assets)
	if err != nil {
		t.Fatal(err)
	}
	html := string(out)
	if !strings.Contains(html, `width:700px`) {
		t.Fatalf("slide width should match image:\n%s", snippet(html, "class=\"slide\""))
	}
	if strings.Contains(html, "padding:39.") {
		t.Fatalf("percent padding must stay %% for HTML, got px:\n%s", snippet(html, "padding"))
	}
	if !strings.Contains(html, "padding:0.6%") {
		t.Fatalf("expected preserved %% padding:\n%s", snippet(html, "padding"))
	}
	hPad := extractStyleProp(html, "text-layer", "height")
	if hPad == "" {
		t.Fatalf("text layer with padding should get computed height, got:\n%s", snippet(html, "text-layer"))
	}
	docNoPad := doc
	docNoPad.Pages[0].Texts[0].Style = psrt.Style(`{"color":"#000000","background":"#ff0000","text-align":"center","font-weight":"700"}`)
	outNoPad, err := RenderHTML(docNoPad, assets)
	if err != nil {
		t.Fatal(err)
	}
	hNoPad := extractStyleProp(string(outNoPad), "text-layer", "height")
	if hNoPad != "" {
		t.Fatalf("text layer without padding must not get computed height, got %q", hNoPad)
	}
	if !strings.Contains(html, "font-size:4.25cqmin") {
		t.Fatalf("expected cqmin font-size:\n%s", snippet(html, "font-size"))
	}
}

func extractStyleProp(html, classFragment, prop string) string {
	i := strings.Index(html, `class="`+classFragment)
	if i < 0 {
		return ""
	}
	j := strings.Index(html[i:], `style="`)
	if j < 0 {
		return ""
	}
	start := i + j + len(`style="`)
	end := strings.Index(html[start:], `"`)
	if end < 0 {
		return ""
	}
	style := html[start : start+end]
	prefix := prop + ":"
	k := strings.Index(style, prefix)
	if k < 0 {
		return ""
	}
	rest := style[k+len(prefix):]
	if semi := strings.Index(rest, ";"); semi >= 0 {
		rest = rest[:semi]
	}
	return rest
}

func buildWebPVP8X(width, height int) []byte {
	vp8x := make([]byte, 8+10)
	copy(vp8x[0:4], "VP8X")
	binary.LittleEndian.PutUint32(vp8x[4:8], 10)
	vp8x[8] = 0x02
	wm1 := width - 1
	hm1 := height - 1
	vp8x[12] = byte(wm1)
	vp8x[13] = byte(wm1 >> 8)
	vp8x[14] = byte(wm1 >> 16)
	vp8x[15] = byte(hm1)
	vp8x[16] = byte(hm1 >> 8)
	vp8x[17] = byte(hm1 >> 16)

	riffSize := 4 + len(vp8x)
	out := make([]byte, 8+4+riffSize)
	copy(out[0:4], "RIFF")
	binary.LittleEndian.PutUint32(out[4:8], uint32(riffSize))
	copy(out[8:12], "WEBP")
	copy(out[12:], vp8x)
	return out
}

package compilehtml

import (
	"fmt"
	"html"
	"strconv"
	"strings"

	"psrt/compileasset"
	"psrt/compileopts"
	"psrt/compilesvg"
	"psrt/psrt"
	"psrt/styleadapter"
)

// RenderHTML produces a standalone HTML document using embedded assets.
func RenderHTML(doc psrt.Document, assets map[string]compileasset.Asset) ([]byte, error) {
	return RenderHTMLBundle([]Variant{{Label: "PSRT", Doc: doc}}, assets, compileopts.Options{})
}

// RenderHTMLBundle renders one or more PSRT variants; Ctrl+L cycles variants and "sem PSRT".
func RenderHTMLBundle(variants []Variant, assets map[string]compileasset.Asset, opts compileopts.Options) ([]byte, error) {
	if len(variants) == 0 {
		return nil, fmt.Errorf("no variants")
	}
	primary := variants[0].Doc
	labels := variantLabels(variants)
	var b strings.Builder
	fontFaces, bodyFontStack := buildFontCSS(primary.Fonts, assets, opts.LinksOnly)

	b.WriteString(`<!DOCTYPE html>
<html lang="pt-BR">
<head>
<meta charset="utf-8"/>
<meta name="viewport" content="width=device-width, initial-scale=1"/>
<title>`)
	title := html.EscapeString("PSRT")
	if len(primary.Pages) > 0 && strings.TrimSpace(primary.Pages[0].Name) != "" {
		title = html.EscapeString(strings.TrimSpace(primary.Pages[0].Name))
	}
	b.WriteString(title)
	b.WriteString(`</title>
<style>`)
	b.WriteString(baseCSS(fontFaces))
	if !opts.NoScript {
		b.WriteString(variantSwitcherCSS())
	}
	b.WriteString(`
body{font-family:`)
	b.WriteString(bodyFontStack)
	b.WriteString(`;margin:0;padding:0;background:#111;overflow-x:auto;}
</style>
</head>
<body>
<main class="slides-wrap">
`)

	for i := range primary.Pages {
		if err := writeSlide(&b, &primary.Pages[i], variants, assets, opts); err != nil {
			return nil, fmt.Errorf("page %q: %w", primary.Pages[i].Name, err)
		}
	}

	b.WriteString(`
</main>
`)
	if !opts.NoScript {
		writeVariantSwitcher(&b, labels)
	}
	b.WriteString(`
</body>
</html>`)

	return []byte(b.String()), nil
}

func variantLabels(variants []Variant) []string {
	labels := make([]string, len(variants))
	for i, v := range variants {
		l := strings.TrimSpace(v.Label)
		if l == "" {
			l = "PSRT"
		}
		labels[i] = l
	}
	return labels
}

func baseCSS(fontFaces string) string {
	return strings.TrimSpace(fontFaces) + `
*{box-sizing:border-box;}
.slides-wrap{margin:0;padding:0;display:flex;width:100%;flex-direction:column;align-items:center;}
.slide{position:relative;display:block;flex:0 0 auto;line-height:0;margin:0;padding:0;}
.slide-img{display:block;width:100%;height:auto;margin:0;padding:0;vertical-align:bottom;}
.slide-overlay{position:absolute;left:0;top:0;right:0;bottom:0;overflow:hidden;container-type:size;container-name:slide;}
.text-layer{position:absolute;box-sizing:border-box;margin:0;padding:0;line-height:1.2;overflow:hidden;overflow-wrap:anywhere;word-wrap:break-word;white-space:pre-wrap;}
.text-ref-img{display:block;max-width:100%;height:auto;margin:0 0 .25em;padding:0;}
`
}

func buildFontCSS(fontURLs []string, assets map[string]compileasset.Asset, linksOnly bool) (fontFacesCSS, bodyStack string) {
	if len(fontURLs) == 0 {
		bodyStack = "-apple-system,BlinkMacSystemFont,Segoe UI,Roboto,sans-serif"
		return "", bodyStack
	}
	var fb strings.Builder
	var names []string
	for i, u := range fontURLs {
		u = strings.TrimSpace(u)
		a, ok := assets[u]
		if !ok && !linksOnly {
			continue
		}
		name := compileasset.FontFamilyNameForURL(u, i)
		src := compileasset.FontSrcURL(u, a, linksOnly)
		format := faceFormat(a.MIME)
		if linksOnly || ok {
			fmt.Fprintf(&fb, "@font-face{font-family:'%s';src:%s%s;font-display:swap;}\n", name, src, format)
		}
		names = append(names, fmt.Sprintf("'%s'", name))
	}
	if len(names) == 0 {
		bodyStack = "-apple-system,BlinkMacSystemFont,Segoe UI,Roboto,sans-serif"
		return fb.String(), bodyStack
	}
	stack := strings.Join(names, ",") + ",-apple-system,BlinkMacSystemFont,Segoe UI,Roboto,sans-serif"
	return fb.String(), stack
}

func faceFormat(m string) string {
	switch strings.TrimSpace(strings.ToLower(m)) {
	case "font/woff2":
		return " format('woff2')"
	case "font/woff":
		return " format('woff')"
	case "font/ttf":
		return " format('truetype')"
	case "font/otf":
		return " format('opentype')"
	default:
		return ""
	}
}

func writeSlide(w *strings.Builder, p *psrt.Page, variants []Variant, assets map[string]compileasset.Asset, opts compileopts.Options) error {
	bg := pageBackgroundCSS(p.Style)
	imgURL := strings.TrimSpace(p.ImageURL)
	a, ok := assets[imgURL]
	if !opts.LinksOnly && !ok {
		return fmt.Errorf("missing fetched asset for page image %q", imgURL)
	}
	src := compileasset.AssetRef(imgURL, a, opts.LinksOnly)
	var canvasW, canvasH int
	if ok {
		canvasW, canvasH = compilesvg.ImageDimensions(a.Bytes, a.MIME)
	} else {
		canvasW, canvasH = compilesvg.ImageDimensions(nil, "")
	}

	slideStyle := fmt.Sprintf("width:%dpx", canvasW)
	if bg != "" {
		slideStyle += ";" + bg
	}
	w.WriteString(`<div class="slide" style="`)
	w.WriteString(html.EscapeString(slideStyle))
	w.WriteString(`">`)
	w.WriteString(`<img class="slide-img" src="`)
	w.WriteString(html.EscapeString(src))
	w.WriteString(`" alt=""/>`)

	w.WriteString(`<div class="slide-overlay">`)
	for vi := range variants {
		vp := pageByName(variants[vi].Doc, p.Name)
		if vp == nil {
			continue
		}
		for _, entry := range psrt.PageBlocksByIndex(vp) {
			switch entry.Kind {
			case psrt.BlockText:
				if entry.Text != nil {
					writeTextLayer(w, entry.Text, assets, canvasW, canvasH, vi, vi > 0, opts.LinksOnly)
				}
			case psrt.BlockMask:
				if entry.Mask != nil {
					writeMaskLayer(w, entry.Mask, assets, canvasW, canvasH, vi, vi > 0, opts.LinksOnly)
				}
			case psrt.BlockPathMask:
				if entry.PathMask != nil {
					writePathMaskLayer(w, entry.PathMask, compilesvg.Slug(p.Name), assets, canvasW, canvasH, vi, vi > 0, opts.LinksOnly)
				}
			}
		}
	}
	w.WriteString(`</div>`)
	w.WriteString(`</div>`)
	return nil
}

func pageBackgroundCSS(style psrt.Style) string {
	bg := compileasset.BackgroundColorFromStyle(style)
	if bg == "" {
		return ""
	}
	return "background:" + bg + ";"
}

func writeTextLayer(w *strings.Builder, t *psrt.Text, assets map[string]compileasset.Asset, canvasW, canvasH, variantIndex int, hidden bool, linksOnly bool) {
	content := psrt.NormalizeTextContent(t.Content)
	fontPx := psrt.TextFontSizePx(t.TextSize, canvasW, canvasH)
	ctx := styleadapter.AdaptContext{
		Text:        *t,
		CanvasW:     canvasW,
		CanvasH:     canvasH,
		FontSizePx:  fontPx,
		HTMLCompile: true,
	}
	frags := styleadapter.AdaptHTML(ctx)
	boxCSS, textCSS := styleadapter.HTMLLayerCSS(frags)
	boxCSS = appendTextLayerGeometryCSS(boxCSS, t, content, canvasW, canvasH)

	classes := "text-layer psrt-text " + variantClass(variantIndex)
	if hidden {
		classes += " psrt-hidden"
	}
	w.WriteString(`<div class="`)
	w.WriteString(classes)
	w.WriteString(`" style="`)
	w.WriteString(html.EscapeString(boxCSS))
	w.WriteString(`">`)

	imgRef := strings.TrimSpace(t.ImageRef)
	if imgRef != "" && compileasset.IsAssetReference(imgRef) {
		if aa, ok := assets[imgRef]; linksOnly || (ok && strings.HasPrefix(aa.MIME, "image/")) {
			refURI := compileasset.AssetRef(imgRef, aa, linksOnly)
			w.WriteString(`<img class="text-ref-img" src="`)
			w.WriteString(html.EscapeString(refURI))
			w.WriteString(`" alt="" style="margin:0 0 .25em 0;display:block;max-width:100%;height:auto"/>`)
		}
	}

	escapedLines := psrt.RenderInlineHTML(content)
	if textCSS != "" {
		w.WriteString(`<span style="`)
		w.WriteString(html.EscapeString(textCSS))
		w.WriteString(`">`)
	} else {
		w.WriteString(`<span>`)
	}
	w.WriteString(escapedLines)
	w.WriteString(`</span></div>`)
}

func escapeHTMLPreserveNewlines(s string) string {
	lines := strings.Split(s, "\n")
	for i := range lines {
		lines[i] = html.EscapeString(lines[i])
	}
	return strings.Join(lines, "<br/>")
}

func pct(v float64) string {
	return strconv.FormatFloat(v, 'f', -1, 64) + "%"
}

// appendTextLayerGeometryCSS adds height only when the PSRT box needs a fixed outer size
// (explicit height or vertical padding/border). Otherwise the layer grows with content and
// cqmin font-size is not clipped by a wrong geometry estimate.
func appendTextLayerGeometryCSS(boxCSS string, t *psrt.Text, content string, canvasW, canvasH int) string {
	if canvasH < 1 || strings.Contains(boxCSS, "height:") {
		return boxCSS
	}
	if !compileasset.TextLayerNeedsComputedHeight(t.Style, canvasW, canvasH, t.TextSize) {
		return boxCSS
	}
	_, _, _, geomH := compilesvg.TextBlockGeometry(t, content, canvasW, canvasH)
	if geomH < 1 {
		return boxCSS
	}
	heightPct := float64(geomH) / float64(canvasH) * 100.0
	return boxCSS + "height:" + pct(heightPct) + ";"
}


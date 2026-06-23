package compilesvg

import (
	"fmt"
	"strings"

	"psrt/compileasset"
	"psrt/compilesvg/textoutline"
	"psrt/psrt"
)

func textFontFamilyStack(documentStack string) string {
	def := "'" + textoutline.DefaultFontFamily + "'"
	if documentStack == "" {
		return def + ",sans-serif"
	}
	return documentStack + "," + def + ",sans-serif"
}

// BuildPageSVGDefsCSS returns minimal CSS embedded in compiled SVG (page background only).
func BuildPageSVGDefsCSS(pageSlug string, pageStyle psrt.Style) string {
	var b strings.Builder
	b.WriteString(ruleBlock(PageClass(pageSlug), pageTypographyCSS(pageStyle)))
	bgClass := PageBgClass(pageSlug)
	if fill := compileasset.BackgroundColorFromStyle(pageStyle); fill != "" {
		b.WriteString(fmt.Sprintf(".%s{fill:%s;}\n", bgClass, fill))
	}
	return b.String()
}

// BuildPageTextStylesheet returns CSS for headless text layout (font-face + text block rules).
func BuildPageTextStylesheet(
	pageSlug string,
	pageStyle psrt.Style,
	texts []psrt.Text,
	canvasW, canvasH int,
	fontURLs []string,
	assets map[string]compileasset.Asset,
) string {
	return BuildPageStylesheet(pageSlug, pageStyle, texts, canvasW, canvasH, fontURLs, assets)
}

// BuildPageStylesheet returns the full CSS for one page (font-face, page class, text classes).
func BuildPageStylesheet(
	pageSlug string,
	pageStyle psrt.Style,
	texts []psrt.Text,
	canvasW, canvasH int,
	fontURLs []string,
	assets map[string]compileasset.Asset,
) string {
	var b strings.Builder
	b.WriteString(buildFontFaceCSS(fontURLs, assets))
	b.WriteString(ruleBlock(PageClass(pageSlug), pageTypographyCSS(pageStyle)))
	bgClass := PageBgClass(pageSlug)
	if fill := compileasset.BackgroundColorFromStyle(pageStyle); fill != "" {
		b.WriteString(fmt.Sprintf(".%s{fill:%s;}\n", bgClass, fill))
	}
	for i := range texts {
		t := &texts[i]
		fontPx := psrt.TextFontSizePx(t.TextSize, canvasW, canvasH)
		corners := compileasset.ParseBorderRadiusCorners(t.Style, canvasW, canvasH, fontPx)
		decl := textBlockLayoutCSS()
		decl += compileasset.TextBlockDisplayCSS(t.Style)
		decl += fmt.Sprintf("font-size:%spx;", compileasset.FormatFloatCSS(fontPx))
		if compileasset.BorderRadiusNeedsInnerBoxDecoration(corners) {
			decl += compileasset.CSSBoxFromStyleJSONForCanvas(t.Style, canvasW, canvasH, fontPx)
		} else {
			decl += compileasset.CSSBoxFromStyleJSONNoBackgroundForCanvas(t.Style, canvasW, canvasH, fontPx)
		}
		decl += compileasset.CSSFromStyleJSON(compileasset.StyleJSONWithoutBox(t.Style))
		decl += textStrokeCSSFromStyle(t.Style)
		stack := fontStackFromAssets(fontURLs, assets)
		decl += "font-family:" + textFontFamilyStack(stack) + ";"
		cls := TextClass(pageSlug, t.Index)
		b.WriteString(ruleBlock(cls, decl))
		innerDecl := textInnerFlowCSS() + textInlineMarkupCSS()
		b.WriteString(ruleBlock(TextInnerClass(pageSlug, t.Index), innerDecl))
	}
	return b.String()
}

func ruleBlock(class, decl string) string {
	decl = strings.TrimSpace(decl)
	if decl == "" {
		return ""
	}
	return fmt.Sprintf(".%s{%s}\n", class, decl)
}

func pageTypographyCSS(style psrt.Style) string {
	return compileasset.CSSFromStyleJSON(compileasset.StyleJSONWithoutBackground(style))
}

// textBlockLayoutCSS fills the foreignObject; padding/border shrink the content box (border-box).
func textBlockLayoutCSS() string {
	return "box-sizing:border-box;width:100%;max-width:100%;min-height:100%;margin:0;line-height:1.2;overflow:hidden;overflow-wrap:break-word;word-wrap:break-word;word-break:normal;white-space:pre-wrap;"
}

func textStrokeCSSFromStyle(style psrt.Style) string {
	sw := strokeWidthFromStyle(style)
	sc := strokeColorFromStyle(style)
	if sw == "" && sc == "" {
		return ""
	}
	var b strings.Builder
	if sw != "" {
		b.WriteString("-webkit-text-stroke-width:")
		b.WriteString(sw)
		b.WriteByte(';')
	}
	if sc != "" {
		b.WriteString("-webkit-text-stroke-color:")
		b.WriteString(sc)
		b.WriteByte(';')
	}
	return b.String()
}

func textInlineMarkupCSS() string {
	return "strong{font-weight:bolder;}em{font-style:italic;}u{text-decoration:underline;}s{text-decoration:line-through;}"
}

// textInnerFlowCSS keeps inline markup (<strong>, <em>, …) on the same line inside flex-centered blocks.
func textInnerFlowCSS() string {
	return "display:block;width:100%;min-width:0;margin:0;white-space:inherit;overflow-wrap:inherit;word-wrap:inherit;word-break:inherit;"
}

func buildFontFaceCSS(fontURLs []string, assets map[string]compileasset.Asset) string {
	var b strings.Builder
	for i, u := range fontURLs {
		u = strings.TrimSpace(u)
		a, ok := assets[u]
		if !ok {
			continue
		}
		name := compileasset.FontFamilyNameForURL(u, i)
		dataURI := compileasset.EncodeDataURI(a.MIME, a.Bytes)
		format := faceFormat(a.MIME)
		fmt.Fprintf(&b, "@font-face{font-family:'%s';src:url(%s)%s;font-display:swap;}\n", name, dataURI, format)
	}
	return b.String()
}

func fontStackFromAssets(fontURLs []string, assets map[string]compileasset.Asset) string {
	var names []string
	for i, u := range fontURLs {
		u = strings.TrimSpace(u)
		if _, ok := assets[u]; !ok {
			continue
		}
		names = append(names, fmt.Sprintf("'%s'", compileasset.FontFamilyNameForURL(u, i)))
	}
	if len(names) == 0 {
		return ""
	}
	return strings.Join(names, ",")
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

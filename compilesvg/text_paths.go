package compilesvg

import (
	"fmt"
	"html"
	"strings"

	"psrt/compilesvg/textoutline"
)

func writeOutlinedGlyphs(
	b *strings.Builder,
	pageSlug string,
	textIndex int,
	outlined *textoutline.OutlinedBlock,
) {
	if outlined == nil || len(outlined.Paths) == 0 {
		writeA11yText(b, pageSlug, textIndex, outlinedPlain(outlined))
		return
	}

	b.WriteString(`<g id="`)
	b.WriteString(TextGlyphsID(pageSlug, textIndex))
	b.WriteString(`"`)
	if outlined.FilterID != "" {
		fmt.Fprintf(b, ` filter="url(#%s)"`, svgAttrEscape(outlined.FilterID))
	}
	if outlined.Transform != "" {
		fmt.Fprintf(b, ` transform="%s"`, svgAttrEscape(outlined.Transform))
	}
	b.WriteString(`>`)

	for _, p := range outlined.Paths {
		b.WriteString(`<path d="`)
		b.WriteString(svgAttrEscape(p.D))
		b.WriteString(`"`)
		if p.Fill != "" {
			fmt.Fprintf(b, ` fill="%s"`, svgAttrEscape(p.Fill))
		}
		if p.Stroke != "" {
			fmt.Fprintf(b, ` stroke="%s"`, svgAttrEscape(p.Stroke))
		}
		if p.StrokeWidth != "" {
			fmt.Fprintf(b, ` stroke-width="%s"`, svgAttrEscape(p.StrokeWidth))
		}
		if p.Opacity != "" && p.Opacity != "1" {
			fmt.Fprintf(b, ` opacity="%s"`, svgAttrEscape(p.Opacity))
		}
		if p.PaintOrder != "" {
			fmt.Fprintf(b, ` paint-order="%s"`, svgAttrEscape(p.PaintOrder))
		}
		b.WriteString(`/>`)
	}
	b.WriteString(`</g>`)
	writeA11yText(b, pageSlug, textIndex, outlinedPlain(outlined))
}

func writeA11yText(b *strings.Builder, pageSlug string, textIndex int, plain string) {
	if plain == "" {
		return
	}
	fmt.Fprintf(b, `<text id="%s" display="none">%s</text>`,
		TextA11yID(pageSlug, textIndex), html.EscapeString(plain))
}

func outlinedPlain(o *textoutline.OutlinedBlock) string {
	if o == nil {
		return ""
	}
	return o.PlainText
}

func writeTextImageRef(b *strings.Builder, imgHTML string) {
	imgHTML = strings.TrimSpace(imgHTML)
	if imgHTML == "" {
		return
	}
	b.WriteString(imgHTML)
}

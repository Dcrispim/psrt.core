package compilesvg

import (
	"fmt"
	"strings"

	"psrt/compileasset"
	"psrt/psrt"
	"psrt/compilesvg/textoutline"
	"psrt/styleadapter"
)

func adaptTextFragments(t *psrt.Text, pageSlug string, canvasW, canvasH int, fontPx float64) []styleadapter.StyleFragment {
	ctx := styleadapter.AdaptContext{
		Text:       *t,
		CanvasW:    canvasW,
		CanvasH:    canvasH,
		FontSizePx: fontPx,
		PageSlug:   pageSlug,
		TextIndex:  t.Index,
	}
	return styleadapter.AdaptSVG(ctx)
}

func collectFilterDefs(pageSlug string, texts []psrt.Text, canvasW, canvasH int) string {
	var b strings.Builder
	for i := range texts {
		t := &texts[i]
		fontPx := psrt.TextFontSizePx(t.TextSize, canvasW, canvasH)
		frags := adaptTextFragments(t, pageSlug, canvasW, canvasH, fontPx)
		for _, f := range frags {
			switch f.GetString(styleadapter.TypeKey) {
			case styleadapter.TypeFilter:
				b.WriteString(styleadapter.FilterFragmentToSVG(f))
			case styleadapter.TypeMask:
				b.WriteString(styleadapter.MaskFragmentToSVG(f))
			}
		}
	}
	return b.String()
}

func fragmentByType(frags []styleadapter.StyleFragment, typ string) (styleadapter.StyleFragment, bool) {
	for _, f := range frags {
		if f.GetString(styleadapter.TypeKey) == typ {
			return f, true
		}
	}
	return nil, false
}

func writeTextBlockFromFragments(
	b *strings.Builder,
	frags []styleadapter.StyleFragment,
	pageSlug string,
	textIndex int,
	imgHTML string,
	outlined *textoutline.OutlinedBlock,
	x, y, width, height int,
	corners compileasset.BorderRadiusCorners,
) {
	rect, hasRect := fragmentByType(frags, styleadapter.TypeRect)

	if hasRect && compileasset.BorderRadiusNeedsInnerBoxDecoration(corners) {
		hasRect = false
	}

	if hasRect {
		fmt.Fprintf(b, `<rect x="%d" y="%d" width="%d" height="%d"`, x, y, width, height)
		if dec := styleadapter.RectDecorationAttrs(rect); dec != "" {
			b.WriteByte(' ')
			b.WriteString(dec)
		}
		b.WriteString(`/>`)
	}

	writeTextImageRef(b, imgHTML)
	writeOutlinedGlyphs(b, pageSlug, textIndex, outlined)
}

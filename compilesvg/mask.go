package compilesvg

import (
	"fmt"
	"strings"

	"psrt/compileasset"
	"psrt/psrt"
	"psrt/styleadapter"
)

// MaskBlockGeometry maps mask percent coords to a pixel rect (fixed height, no text metrics).
func MaskBlockGeometry(m *psrt.Mask, canvasW, canvasH int) (x, y, width, height int) {
	x, y, width, height = 0, 0, 0, 0
	if m == nil {
		return
	}
	x = int(float64(canvasW) * m.X / 100.0)
	y = int(float64(canvasH) * m.Y / 100.0)
	width = textBlockWidthPx(m.Width, canvasW)
	height = int(float64(canvasH) * m.Height / 100.0)
	if width < 1 {
		width = 1
	}
	if height < 1 {
		height = 1
	}
	return x, y, width, height
}

func writeMaskBlock(
	b *strings.Builder,
	m *psrt.Mask,
	pageSlug string,
	canvasW, canvasH int,
	assets map[string]compileasset.Asset,
	linksOnly bool,
) error {
	x, y, width, height := MaskBlockGeometry(m, canvasW, canvasH)
	frags := adaptMaskFragments(m, pageSlug, canvasW, canvasH)
	rect, hasRect := fragmentByType(frags, styleadapter.TypeRect)

	b.WriteString(`<g id="`)
	b.WriteString(TextWrapID(pageSlug, m.Index))
	b.WriteString(`">`)

	if hasRect {
		fmt.Fprintf(b, `<rect x="%d" y="%d" width="%d" height="%d"`, x, y, width, height)
		if dec := styleadapter.RectDecorationAttrs(rect); dec != "" {
			b.WriteByte(' ')
			b.WriteString(dec)
		}
		b.WriteString(`/>`)
	} else {
		fmt.Fprintf(b, `<rect x="%d" y="%d" width="%d" height="%d" fill="transparent"/>`, x, y, width, height)
	}

	imgRef := strings.TrimSpace(m.ImageRef)
	if imgRef != "" && compileasset.IsAssetReference(imgRef) {
		if aa, ok := assets[imgRef]; linksOnly || (ok && IsImageMIME(aa.MIME)) {
			refURI := compileasset.AssetRef(imgRef, aa, linksOnly)
			fmt.Fprintf(b, `<image x="%d" y="%d" width="%d" height="%d" href="%s" preserveAspectRatio="xMidYMid slice"/>`,
				x, y, width, height, refURI)
		}
	}

	b.WriteString(`</g>`)
	return nil
}

func adaptMaskFragments(m *psrt.Mask, pageSlug string, canvasW, canvasH int) []styleadapter.StyleFragment {
	ctx := styleadapter.AdaptContext{
		Mask:       m,
		CanvasW:    canvasW,
		CanvasH:    canvasH,
		PageSlug:   pageSlug,
		TextIndex:  m.Index,
	}
	return styleadapter.AdaptMaskSVG(ctx)
}

package compilesvg

import (
	"fmt"
	"strings"

	"github.com/Dcrispim/psrt.core/compileasset"
	"github.com/Dcrispim/psrt.core/psrt"
	"github.com/Dcrispim/psrt.core/styleadapter"
)

// PathMaskBlockGeometry maps ~~ percent coords to a pixel rect for the nested
// <svg viewBox="0 0 100 100"> wrapper (same percent->pixel conversion as
// MaskBlockGeometry; the path's own 0-100 coordinates are left untouched and
// scaled by the nested viewBox, including non-uniformly when Width != Height).
func PathMaskBlockGeometry(m *psrt.PathMask, canvasW, canvasH int) (x, y, width, height int) {
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

func writePathMaskBlock(
	b *strings.Builder,
	m *psrt.PathMask,
	pageSlug string,
	canvasW, canvasH int,
	assets map[string]compileasset.Asset,
	linksOnly bool,
) error {
	x, y, width, height := PathMaskBlockGeometry(m, canvasW, canvasH)
	frags := adaptPathMaskFragments(m, pageSlug, canvasW, canvasH)
	path, hasPath := fragmentByType(frags, styleadapter.TypePath)

	clipID := PathMaskClipID(pageSlug, m.Index)
	d := svgAttrEscape(m.Path)

	b.WriteString(`<g id="`)
	b.WriteString(PathMaskWrapID(pageSlug, m.Index))
	b.WriteString(`">`)
	fmt.Fprintf(b, `<svg x="%d" y="%d" width="%d" height="%d" viewBox="0 0 100 100" preserveAspectRatio="none">`,
		x, y, width, height)
	fmt.Fprintf(b, `<defs><clipPath id="%s"><path d="%s"/></clipPath></defs>`, clipID, d)

	if hasPath {
		fmt.Fprintf(b, `<path d="%s"`, d)
		if dec := styleadapter.PathDecorationAttrs(path); dec != "" {
			b.WriteByte(' ')
			b.WriteString(dec)
		}
		b.WriteString(`/>`)
	} else {
		fmt.Fprintf(b, `<path d="%s" fill="transparent"/>`, d)
	}

	imgRef := strings.TrimSpace(m.ImageRef)
	if imgRef != "" && compileasset.IsAssetReference(imgRef) {
		if aa, ok := assets[imgRef]; linksOnly || (ok && IsImageMIME(aa.MIME)) {
			refURI := compileasset.AssetRef(imgRef, aa, linksOnly)
			fmt.Fprintf(b, `<image x="0" y="0" width="100" height="100" href="%s" clip-path="url(#%s)" preserveAspectRatio="xMidYMid slice"/>`,
				refURI, clipID)
		}
	}

	b.WriteString(`</svg></g>`)
	return nil
}

func adaptPathMaskFragments(m *psrt.PathMask, pageSlug string, canvasW, canvasH int) []styleadapter.StyleFragment {
	ctx := styleadapter.AdaptContext{
		PathMask:  m,
		CanvasW:   canvasW,
		CanvasH:   canvasH,
		PageSlug:  pageSlug,
		TextIndex: m.Index,
	}
	return styleadapter.AdaptPathMaskSVG(ctx)
}

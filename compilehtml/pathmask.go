package compilehtml

import (
	"fmt"
	"html"
	"strings"

	"psrt/compileasset"
	"psrt/psrt"
	"psrt/styleadapter"
)

// writePathMaskLayer renders a ~~ block as a positioned <div> containing an
// inline <svg> with its own viewBox — the same nested-viewBox technique used
// by compilesvg/pathmask.go, so non-uniform scaling (Width != Height) and
// arbitrary curves/arcs are handled natively by the SVG renderer instead of
// by hand-rolled geometry code.
func writePathMaskLayer(
	w *strings.Builder,
	m *psrt.PathMask,
	pageSlug string,
	assets map[string]compileasset.Asset,
	canvasW, canvasH, variantIndex int,
	hidden bool,
	linksOnly bool,
) {
	ctx := styleadapter.AdaptContext{
		PathMask:    m,
		CanvasW:     canvasW,
		CanvasH:     canvasH,
		HTMLCompile: true,
		PageSlug:    pageSlug,
		TextIndex:   m.Index,
	}
	frags := styleadapter.AdaptPathMaskHTML(ctx)
	boxCSS, _ := styleadapter.HTMLLayerCSS(frags)

	var pathFrag styleadapter.StyleFragment
	for _, f := range frags {
		if f.GetString(styleadapter.TypeKey) == styleadapter.TypePath {
			pathFrag = f
		}
	}
	pathAttrs := styleadapter.PathDecorationAttrs(pathFrag)

	classes := "text-layer psrt-mask psrt-path-mask " + variantClass(variantIndex)
	if hidden {
		classes += " psrt-hidden"
	}
	w.WriteString(`<div class="`)
	w.WriteString(classes)
	w.WriteString(`" style="`)
	w.WriteString(html.EscapeString(boxCSS))
	w.WriteString(`">`)

	clipID := fmt.Sprintf("psrt-pathmask-%s-%d-%d-clip", pageSlug, variantIndex, m.Index)
	d := html.EscapeString(m.Path)

	w.WriteString(`<svg width="100%" height="100%" viewBox="0 0 100 100" preserveAspectRatio="none">`)
	fmt.Fprintf(w, `<defs><clipPath id="%s"><path d="%s"/></clipPath></defs>`, clipID, d)

	if pathAttrs != "" {
		fmt.Fprintf(w, `<path d="%s" %s/>`, d, pathAttrs)
	} else {
		fmt.Fprintf(w, `<path d="%s" fill="transparent"/>`, d)
	}

	imgRef := strings.TrimSpace(m.ImageRef)
	if imgRef != "" && compileasset.IsAssetReference(imgRef) {
		if aa, ok := assets[imgRef]; linksOnly || (ok && strings.HasPrefix(aa.MIME, "image/")) {
			refURI := compileasset.AssetRef(imgRef, aa, linksOnly)
			fmt.Fprintf(w, `<image x="0" y="0" width="100" height="100" href="%s" clip-path="url(#%s)" preserveAspectRatio="xMidYMid slice"/>`,
				html.EscapeString(refURI), clipID)
		}
	}

	w.WriteString(`</svg></div>`)
}

package compilehtml

import (
	"html"
	"strings"

	"psrt/compileasset"
	"psrt/psrt"
	"psrt/styleadapter"
)

func writeMaskLayer(w *strings.Builder, m *psrt.Mask, assets map[string]compileasset.Asset, canvasW, canvasH, variantIndex int, hidden bool, linksOnly bool) {
	ctx := styleadapter.AdaptContext{
		Mask:        m,
		CanvasW:     canvasW,
		CanvasH:     canvasH,
		HTMLCompile: true,
	}
	frags := styleadapter.AdaptMaskHTML(ctx)
	boxCSS, _ := styleadapter.HTMLLayerCSS(frags)

	imgRef := strings.TrimSpace(m.ImageRef)
	if imgRef != "" && compileasset.IsAssetReference(imgRef) {
		if aa, ok := assets[imgRef]; linksOnly || (ok && strings.HasPrefix(aa.MIME, "image/")) {
			refURI := compileasset.AssetRef(imgRef, aa, linksOnly)
			if boxCSS != "" && !strings.HasSuffix(boxCSS, ";") {
				boxCSS += ";"
			}
			boxCSS += "background-image:url(" + refURI + ");background-size:cover;background-position:center;background-repeat:no-repeat;"
		}
	}

	classes := "text-layer psrt-mask " + variantClass(variantIndex)
	if hidden {
		classes += " psrt-hidden"
	}
	w.WriteString(`<div class="`)
	w.WriteString(classes)
	w.WriteString(`" style="`)
	w.WriteString(html.EscapeString(boxCSS))
	w.WriteString(`"></div>`)
}

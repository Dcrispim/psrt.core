package compileasset

import (
	"encoding/json"
	"strings"

	"github.com/Dcrispim/psrt.core/psrt"
	"github.com/Dcrispim/psrt.core/styleadapter"
	"github.com/Dcrispim/psrt.core/styleadapter/percent"
)

// StyleResolvedForCanvas applies percent handlers (padding, border-width, etc.) for render-time CSS/layout.
func StyleResolvedForCanvas(style psrt.Style, canvasW, canvasH int, fontPx float64) psrt.Style {
	m := styleadapter.Normalize(style)
	if len(m) == 0 {
		return style
	}
	m = percent.ApplyPercentHandlers(m, percent.ImageDims{
		W:          canvasW,
		H:          canvasH,
		FontSizePx: fontPx,
		Zoom:       1,
	})
	b, err := json.Marshal(m)
	if err != nil {
		return style
	}
	return psrt.Style(b)
}

// CSSBoxFromStyleJSONNoBackgroundForCanvas emits box CSS with % lengths resolved to px.
func CSSBoxFromStyleJSONNoBackgroundForCanvas(style psrt.Style, canvasW, canvasH int, fontPx float64) string {
	resolved := StyleResolvedForCanvas(style, canvasW, canvasH, fontPx)
	var b strings.Builder
	b.WriteString(cssBoxFromStyle(resolved, false, false))
	b.WriteString(BorderRadiusCSS(ParseBorderRadiusCorners(style, canvasW, canvasH, fontPx)))
	return b.String()
}

// CSSBoxFromStyleJSONForCanvas emits full box CSS (background, border, radius, padding) with % resolved to px.
func CSSBoxFromStyleJSONForCanvas(style psrt.Style, canvasW, canvasH int, fontPx float64) string {
	resolved := StyleResolvedForCanvas(style, canvasW, canvasH, fontPx)
	var b strings.Builder
	b.WriteString(cssBoxFromStyle(resolved, true, false))
	b.WriteString(BorderRadiusCSS(ParseBorderRadiusCorners(style, canvasW, canvasH, fontPx)))
	return b.String()
}

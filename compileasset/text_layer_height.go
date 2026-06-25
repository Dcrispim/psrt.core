package compileasset

import "github.com/Dcrispim/psrt.core/psrt"

// TextLayerNeedsComputedHeight reports whether a text layer needs an explicit outer height
// (explicit height in style or vertical padding/border) for flex vertical alignment to apply.
func TextLayerNeedsComputedHeight(style psrt.Style, canvasW, canvasH int, textSize float64) bool {
	fontPx := psrt.TextFontSizePx(textSize, canvasW, canvasH)
	if _, ok := ExplicitHeightPx(style, canvasW, canvasH, fontPx); ok {
		return true
	}
	insets := TextBoxInsetsForCanvas(style, fontPx, canvasW, canvasH)
	return insets.Vertical() > 0
}

package compilesvg

import (
	"math"
	"strings"

	"psrt/compileasset"
	"psrt/psrt"
)

// TextBlockGeometry maps PSRT percent coords to a pixel rect for layout (SVG foreignObject / HTML text-layer).
// x, y, width match the editor/web border box (box-sizing: border-box); height includes vertical padding.
func TextBlockGeometry(t *psrt.Text, content string, canvasW, canvasH int) (x, y, width, height int) {
	x = int(math.Round(float64(canvasW) * t.X / 100.0))
	y = int(math.Round(float64(canvasH) * t.Y / 100.0))
	outerW := textBlockWidthPx(t.Width, canvasW)
	fontPx := psrt.TextFontSizePx(t.TextSize, canvasW, canvasH)
	insets := compileasset.TextBoxInsetsForCanvas(t.Style, fontPx, canvasW, canvasH)
	padW := int(math.Round(insets.Horizontal()))
	padH := int(math.Round(insets.Vertical()))
	contentW := outerW - padW
	if contentW < 1 {
		contentW = 1
	}
	plain := psrt.PlainTextForLayout(content)
	lines := estimateTextLines(plain, contentW, fontPx, t.Style)
	lh := compileasset.LineHeightMultiplier(t.Style, fontPx)
	linePx := fontPx * lh
	contentH := int(math.Round(linePx * float64(lines)))
	if contentH < int(math.Round(linePx)) {
		contentH = int(math.Round(linePx))
	}
	width = outerW
	height = contentH + padH
	if explicitH, ok := compileasset.ExplicitHeightPx(t.Style, canvasW, canvasH, fontPx); ok {
		if strings.TrimSpace(content) == "" {
			height = explicitH
		} else if explicitH > height {
			height = explicitH
		}
	}
	if width < 1 {
		width = 1
	}
	if height < 1 {
		height = 1
	}
	return x, y, width, height
}

func textBlockWidthPx(widthPct float64, canvasW int) int {
	w := int(math.Round(float64(canvasW) * widthPct / 100.0))
	if w < 1 {
		return 1
	}
	return w
}

func estimateTextLines(content string, widthPx int, fontSizePx float64, style psrt.Style) int {
	if widthPx < 1 {
		widthPx = 1
	}
	charsPerLine := charsPerLineForWidth(widthPx, fontSizePx, style)
	parts := strings.Split(content, "\n")
	total := 0
	for _, part := range parts {
		n := len([]rune(strings.TrimSpace(part)))
		if n == 0 {
			continue
		}
		total += (n + charsPerLine - 1) / charsPerLine
	}
	if total < 1 {
		return 1
	}
	return total
}

// charsPerLineForWidth estimates how many characters fit per line (conservative vs browser wrapping).
func charsPerLineForWidth(widthPx int, fontSizePx float64, style psrt.Style) int {
	em := 0.48
	if compileasset.FontWeightIsBold(style) {
		em = 0.58
	}
	cpl := int(float64(widthPx) / (fontSizePx * em))
	if cpl < 1 {
		return 1
	}
	return cpl
}

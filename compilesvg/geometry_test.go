package compilesvg

import (
	"math"
	"testing"

	"psrt/psrt"
)

func TestTextBlockGeometry_introTitle(t *testing.T) {
	text := psrt.Text{
		BaseBlock: psrt.BaseBlock{X: 11.6, Y: 55.11, Width: 77, Style: psrt.Style(`{"font-weight":"600"}`)},
		TextSize:  3,
		Content:   "O Soberando supremo da eternidade\ns",
	}
	content := "O Soberando supremo da eternidade\ns"
	_, _, w, h := TextBlockGeometry(&text, content, 1080, 1920)
	if w != 832 {
		t.Fatalf("width: got %d want 832", w)
	}
	lines := estimateTextLines(content, 832, psrt.TextFontSizePx(3, 1080, 1920), text.Style)
	if lines < 2 || lines > 4 {
		t.Fatalf("lines: got %d want 2-4", lines)
	}
	linePx := psrt.TextFontSizePx(3, 1080, 1920) * 1.2
	wantH := int(linePx * float64(lines))
	if h < wantH-2 || h > wantH+2 {
		t.Fatalf("height: got %d want ~%d (lines=%d)", h, wantH, lines)
	}
}

func TestTextBlockGeometry_paddingExpandsBox(t *testing.T) {
	text := psrt.Text{
		BaseBlock: psrt.BaseBlock{X: 22.6, Y: 56.11, Width: 77, Style: psrt.Style(`{"padding":"10px","font-weight":"600"}`)},
		TextSize:  3,
		Content:   "One line",
	}
	_, _, w, h := TextBlockGeometry(&text, "One line", 1080, 1920)
	if w != 832 {
		t.Fatalf("outer width with padding: got %d want %d (padding is inside border-box)", w, 832)
	}
	fontPx := psrt.TextFontSizePx(3, 1080, 1920)
	linePx := fontPx * 1.2
	wantH := int(linePx) + 20
	if h < wantH-2 || h > wantH+2 {
		t.Fatalf("height with padding: got %d want ~%d", h, wantH)
	}
}

func TestTextBlockGeometry_percentPaddingExpandsBox(t *testing.T) {
	text := psrt.Text{
		BaseBlock: psrt.BaseBlock{X: 22.6, Y: 56.11, Width: 77, Style: psrt.Style(`{"padding":"0.515%","font-weight":"600"}`)},
		TextSize:  3,
		Content:   "One line",
	}
	canvasW, canvasH := 1080, 1920
	_, _, w, h := TextBlockGeometry(&text, "One line", canvasW, canvasH)
	base := float64(canvasH) // max(W,H) for single-value padding %
	pad := int(math.Round(base * 0.00515 * 2))
	if w != 832 {
		t.Fatalf("outer width with %% padding: got %d want %d", w, 832)
	}
	fontPx := psrt.TextFontSizePx(3, canvasW, canvasH)
	linePx := fontPx * 1.2
	wantH := int(linePx) + pad
	if h < wantH-2 || h > wantH+2 {
		t.Fatalf("height with %% padding: got %d want ~%d", h, wantH)
	}
}

func TestTextBlockGeometry_emptyContentExplicitHeight(t *testing.T) {
	text := psrt.Text{
		BaseBlock: psrt.BaseBlock{X: 10, Y: 10, Width: 50, Style: psrt.Style(`{"height":"20%","background":"#ff0000"}`)},
		TextSize:  3,
		Content:   "",
	}
	_, _, _, h := TextBlockGeometry(&text, "", 1080, 1920)
	wantH := int(1920 * 0.20)
	if h < wantH-2 || h > wantH+2 {
		t.Fatalf("height: got %d want ~%d", h, wantH)
	}
}

func TestCharsPerLineForWidth_bold(t *testing.T) {
	style := psrt.Style(`{"font-weight":"600"}`)
	cpl := charsPerLineForWidth(832, 57.6, style)
	if cpl < 20 || cpl > 28 {
		t.Fatalf("cpl: got %d", cpl)
	}
}

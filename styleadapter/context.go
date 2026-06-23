package styleadapter

import (
	"encoding/json"

	"psrt/psrt"
	"psrt/styleadapter/percent"
)

// AdaptContext carries normalized style and layout for one text or mask block.
type AdaptContext struct {
	Style      map[string]json.RawMessage
	Text       psrt.Text
	Mask       *psrt.Mask
	CanvasW    int
	CanvasH    int
	FontSizePx float64
	Zoom       float64
	// HTMLCompile marks compilehtml output (font-size in px, same as SVG).
	HTMLCompile bool
	// PageSlug and TextIndex for SVG filter ids.
	PageSlug  string
	TextIndex int
}

func (ctx AdaptContext) ImageDims() percent.ImageDims {
	z := ctx.Zoom
	if z <= 0 {
		z = 1
	}
	return percent.ImageDims{
		W:          ctx.CanvasW,
		H:          ctx.CanvasH,
		FontSizePx: ctx.FontSizePx,
		Zoom:       z,
	}
}

func (ctx AdaptContext) FontSizePxOrCompute() float64 {
	if ctx.FontSizePx > 0 {
		return ctx.FontSizePx
	}
	return psrt.TextFontSizePx(ctx.Text.TextSize, ctx.CanvasW, ctx.CanvasH)
}

package visualapp

import (
	"encoding/json"

	"github.com/Dcrispim/psrt.core/psrt"
	"github.com/Dcrispim/psrt.core/styleadapter"
)

// WebPreviewStyle is adapted CSS for one text block in the web preview.
type WebPreviewStyle = styleadapter.WebPreviewStyle

// WebEntryStyleInput identifies one text or mask block for batch adaptation.
type WebEntryStyleInput struct {
	Index    int     `json:"index"`
	Style    string  `json:"style"`
	Content  string  `json:"content"`
	X        float64 `json:"x"`
	Y        float64 `json:"y"`
	Width    float64 `json:"width"`
	TextSize float64 `json:"textSize"`
	Height   float64 `json:"height,omitempty"`
	IsMask   bool    `json:"isMask,omitempty"`
}

// AdaptTextStyleForWeb adapts a single text style for the web preview.
func AdaptTextStyleForWeb(
	styleJSON string,
	content string,
	x, y, width, textSize float64,
	canvasW, canvasH int,
	zoom float64,
) WebPreviewStyle {
	return adaptOneWeb(styleJSON, content, x, y, width, textSize, canvasW, canvasH, zoom)
}

// AdaptEntriesForWeb adapts all text entries for the web preview in one call.
func AdaptEntriesForWeb(entriesJSON string, canvasW, canvasH int, zoom float64) ([]WebPreviewStyle, error) {
	var inputs []WebEntryStyleInput
	if err := json.Unmarshal([]byte(entriesJSON), &inputs); err != nil {
		return nil, err
	}
	out := make([]WebPreviewStyle, 0, len(inputs))
	for _, in := range inputs {
		if in.IsMask {
			out = append(out, adaptMaskWeb(in, canvasW, canvasH, zoom))
		} else {
			out = append(out, adaptOneWeb(in.Style, in.Content, in.X, in.Y, in.Width, in.TextSize, canvasW, canvasH, zoom))
		}
	}
	return out, nil
}

func adaptMaskWeb(in WebEntryStyleInput, canvasW, canvasH int, zoom float64) WebPreviewStyle {
	if zoom <= 0 {
		zoom = 1
	}
	height := in.Height
	if height < 0.5 {
		height = 5
	}
	ctx := styleadapter.AdaptContext{
		Mask: &psrt.Mask{
			BaseBlock: psrt.BaseBlock{
				X: in.X, Y: in.Y, Width: in.Width, Style: psrt.Style(in.Style),
			},
			Height: height,
		},
		CanvasW: canvasW,
		CanvasH: canvasH,
		Zoom:    zoom,
	}
	return styleadapter.AdaptMaskWebPreview(ctx)
}

func adaptOneWeb(
	styleJSON string,
	content string,
	x, y, width, textSize float64,
	canvasW, canvasH int,
	zoom float64,
) WebPreviewStyle {
	if zoom <= 0 {
		zoom = 1
	}
	fontPx := psrt.TextFontSizePx(textSize, canvasW, canvasH) * zoom
	if fontPx < 1 {
		fontPx = 1
	}
	ctx := styleadapter.AdaptContext{
		Text: psrt.Text{
			BaseBlock: psrt.BaseBlock{X: x, Y: y, Width: width, Style: psrt.Style(styleJSON)},
			TextSize:  textSize,
		},
		CanvasW:    canvasW,
		CanvasH:    canvasH,
		FontSizePx: fontPx,
		Zoom:       zoom,
	}
	out := styleadapter.AdaptWebPreview(ctx)
	enrichWebPreviewHeight(&out, ctx.Text, content, canvasW, canvasH)
	return out
}

package styleadapter

import "fmt"

// WebPreviewStyle is CSS for web preview (container = box/layout, text = typography + stroke).
type WebPreviewStyle struct {
	Container map[string]string `json:"container"`
	Text      map[string]string `json:"text"`
	HasStroke bool              `json:"hasStroke"`
}

// AdaptMaskWebPreview adapts a == mask block for live web preview (fixed height %).
func AdaptMaskWebPreview(ctx AdaptContext) WebPreviewStyle {
	if ctx.Mask == nil {
		return WebPreviewStyle{}
	}
	if ctx.Zoom <= 0 {
		ctx.Zoom = 1
	}
	ctx.HTMLCompile = false
	frags := adaptMask(ctx, true)
	var box StyleFragment
	for _, f := range MergeFragments(frags) {
		if f.GetString(TypeKey) == TypeMotionDiv {
			box = f
		}
	}
	postProcessBackdropGlass(box)
	return WebPreviewStyle{
		Container: fragmentToStringMap(box),
		Text:      map[string]string{},
		HasStroke: false,
	}
}

// AdaptWebPreview runs AdaptHTML for live web preview (font-size in px, not cqh).
func AdaptWebPreview(ctx AdaptContext) WebPreviewStyle {
	ctx.HTMLCompile = false
	if ctx.Zoom <= 0 {
		ctx.Zoom = 1
	}
	frags := AdaptHTML(ctx)
	var box, span StyleFragment
	for _, f := range MergeFragments(frags) {
		switch f.GetString(TypeKey) {
		case TypeMotionDiv:
			box = f
		case TypeSpan:
			span = f
		}
	}
	postProcessBackdropGlass(box)
	out := WebPreviewStyle{
		Container: fragmentToStringMap(box),
		Text:      fragmentToStringMap(span),
	}
	out.HasStroke = out.Text["WebkitTextStrokeWidth"] != "" ||
		out.Text["WebkitTextStroke"] != "" ||
		out.Text["WebkitTextStrokeColor"] != ""
	return out
}

func fragmentToStringMap(f StyleFragment) map[string]string {
	if f == nil {
		return nil
	}
	out := make(map[string]string)
	for k, v := range f {
		if k == TypeKey {
			continue
		}
		out[k] = fmt.Sprint(v)
	}
	return out
}

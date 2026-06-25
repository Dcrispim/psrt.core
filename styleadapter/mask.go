package styleadapter

import (
	"encoding/json"
	"strconv"

	"github.com/Dcrispim/psrt.core/styleadapter/percent"
)

// AdaptMaskHTML returns box CSS fragments for a static mask block.
func AdaptMaskHTML(ctx AdaptContext) []StyleFragment {
	return adaptMask(ctx, true)
}

// AdaptMaskSVG returns SVG fragments for a static mask block.
func AdaptMaskSVG(ctx AdaptContext) []StyleFragment {
	return adaptMask(ctx, false)
}

func adaptMask(ctx AdaptContext, html bool) []StyleFragment {
	m := ctx.Mask
	if m == nil {
		return nil
	}
	style := Normalize(m.Style)
	if len(style) == 0 {
		style = map[string]json.RawMessage{}
	}
	dims := ctx.ImageDims()
	if html && ctx.HTMLCompile {
		dims.PreservePercent = true
	}
	style = percent.ApplyPercentHandlers(style, dims)
	style = FilterStyleMap(style)

	filterID := maskFilterID(ctx)
	var effectFrags []StyleFragment
	style, effectFrags = expandEffects(style, dims, !html, filterID)

	var blurMeta BlurAdapt
	var blurFrags []StyleFragment
	style, blurMeta, blurFrags = expandBlur(style, dims, html, filterID)

	var frags []StyleFragment
	if html {
		frags = buildMaskHTMLFragments(ctx, style)
	} else {
		frags = buildMaskSVGFragments(ctx, style, filterID, blurMeta)
	}
	frags = append(frags, effectFrags...)
	frags = append(frags, blurFrags...)
	return MergeFragments(frags)
}

func maskFilterID(ctx AdaptContext) string {
	idx := ctx.TextIndex
	if ctx.Mask != nil {
		idx = ctx.Mask.Index
	}
	return "psrt-filter-" + ctx.PageSlug + "-" + strconv.Itoa(idx)
}

func buildMaskHTMLFragments(ctx AdaptContext, style map[string]json.RawMessage) []StyleFragment {
	box := NewFragment(TypeMotionDiv)
	applyMaskLayout(box, ctx)
	applyTransform(box, style, true)
	for k, raw := range style {
		if !HasStyleValue(k, raw) {
			continue
		}
		val := SanitizeCSSValue(StringifyCSSValue(raw))
		if isBoxKey(k) {
			applyBoxCSS(box, k, val)
		}
	}
	if box.GetString("background-size") == "" {
		box.Set("background-size", "cover")
	}
	return []StyleFragment{box}
}

func buildMaskSVGFragments(ctx AdaptContext, style map[string]json.RawMessage, filterID string, blurMeta BlurAdapt) []StyleFragment {
	g := NewFragment(TypeG)
	rect := NewFragment(TypeRect)
	applyMaskLayoutSVG(g, rect, ctx)
	applyTransform(g, style, false)

	hasBox := false
	for k, raw := range style {
		if !HasStyleValue(k, raw) {
			continue
		}
		val := SanitizeCSSValue(StringifyCSSValue(raw))
		if isBoxKey(k) {
			hasBox = applyBoxSVG(rect, k, val) || hasBox
			if k == KeyBoxShadow && val != "" {
				rect.Set("filter", "url(#"+filterID+")")
			}
		}
	}
	if blurMeta.FilterID != "" {
		applyBlurSVGRect(rect, blurMeta)
		hasBox = true
	}
	var out []StyleFragment
	out = append(out, g)
	if hasBox || rect.GetString(KeyBackground) != "" || rect.GetString("fill") != "" {
		out = append(out, rect)
	}
	return out
}

func applyMaskLayout(box StyleFragment, ctx AdaptContext) {
	m := ctx.Mask
	if m == nil {
		return
	}
	box.Set(KeyPosition, "absolute")
	box.Set(KeyBoxSizing, "border-box")
	box.Set(KeyLeft, pctString(m.X))
	box.Set(KeyTop, pctString(m.Y))
	box.Set(KeyWidth, pctString(m.Width))
	box.Set(KeyHeight, pctString(m.Height))
}

func applyMaskLayoutSVG(g, rect StyleFragment, ctx AdaptContext) {
	m := ctx.Mask
	if m == nil {
		return
	}
	x := int(float64(ctx.CanvasW) * m.X / 100.0)
	y := int(float64(ctx.CanvasH) * m.Y / 100.0)
	w := int(float64(ctx.CanvasW) * m.Width / 100.0)
	if w < 1 {
		w = 1
	}
	g.Set("x", x)
	g.Set("y", y)
	rect.Set("x", x)
	rect.Set("y", y)
	rect.Set(KeyWidth, w)
}

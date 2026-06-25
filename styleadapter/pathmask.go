package styleadapter

import (
	"encoding/json"
	"strconv"

	"psrt/styleadapter/percent"
)

// AdaptPathMaskHTML returns style fragments for a ~~ block (compilehtml pipeline).
// The block is always rendered as an inline <svg> (see compilehtml/pathmask.go),
// so decoration targets SVG presentation attributes in both pipelines.
func AdaptPathMaskHTML(ctx AdaptContext) []StyleFragment {
	return adaptPathMask(ctx, true)
}

// AdaptPathMaskSVG returns style fragments for a ~~ block (compilesvg pipeline).
func AdaptPathMaskSVG(ctx AdaptContext) []StyleFragment {
	return adaptPathMask(ctx, false)
}

func adaptPathMask(ctx AdaptContext, html bool) []StyleFragment {
	m := ctx.PathMask
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

	filterID := pathMaskFilterID(ctx)
	// The block is always emitted as a <path> inside a nested <svg> (HTML or
	// SVG pipeline alike), so effects/blur always target SVG filters, never
	// HTML/CSS — unlike rect masks, which branch on the html flag.
	var effectFrags []StyleFragment
	style, effectFrags = expandEffects(style, dims, true, filterID)

	var blurMeta BlurAdapt
	var blurFrags []StyleFragment
	style, blurMeta, blurFrags = expandBlur(style, dims, false, filterID)

	frags := buildPathMaskFragments(style, blurMeta)
	if html {
		box := NewFragment(TypeMotionDiv)
		applyPathMaskLayout(box, ctx)
		frags = append(frags, box)
	}
	frags = append(frags, effectFrags...)
	frags = append(frags, blurFrags...)
	return MergeFragments(frags)
}

// applyPathMaskLayout positions the wrapping <div> for the compilehtml
// pipeline — the nested <svg> inside it (see compilehtml/pathmask.go) fills
// 100%/100% of this box and handles the path's own scaling via viewBox.
func applyPathMaskLayout(box StyleFragment, ctx AdaptContext) {
	m := ctx.PathMask
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

func pathMaskFilterID(ctx AdaptContext) string {
	idx := ctx.TextIndex
	if ctx.PathMask != nil {
		idx = ctx.PathMask.Index
	}
	return "psrt-filter-" + ctx.PageSlug + "-" + strconv.Itoa(idx)
}

func buildPathMaskFragments(style map[string]json.RawMessage, blurMeta BlurAdapt) []StyleFragment {
	path := NewFragment(TypePath)
	hasDecoration := false
	for k, raw := range style {
		if !HasStyleValue(k, raw) {
			continue
		}
		if isBorderRadiusKey(k) {
			// Border-radius has no meaning for an arbitrary path — the
			// contour is already defined by the path itself.
			continue
		}
		if !isBoxKey(k) {
			continue
		}
		val := SanitizeCSSValue(StringifyCSSValue(raw))
		hasDecoration = applyBoxSVG(path, k, val) || hasDecoration
	}
	if blurMeta.FilterID != "" {
		applyBlurSVGRect(path, blurMeta)
		hasDecoration = true
	}
	if !hasDecoration && path.GetString("fill") == "" {
		return nil
	}
	return []StyleFragment{path}
}

func isBorderRadiusKey(k string) bool {
	switch k {
	case KeyBorderRadius, KeyBorderTopLeftRadius, KeyBorderTopRightRadius,
		KeyBorderBottomRightRadius, KeyBorderBottomLeftRadius:
		return true
	}
	return false
}

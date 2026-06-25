package styleadapter

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Dcrispim/psrt.core/psrt"
	"github.com/Dcrispim/psrt.core/styleadapter/percent"
)

// defaultAdapter implements Adapter.
type defaultAdapter struct{}

// Default is the package-level style adapter.
var Default Adapter = &defaultAdapter{}

// Adapter converts PSRT styles to render fragments.
type Adapter interface {
	ResolveName(raw string) (string, bool)
	Normalize(style psrt.Style) map[string]json.RawMessage
	AdaptHTML(ctx AdaptContext) []StyleFragment
	AdaptSVG(ctx AdaptContext) []StyleFragment
}

func (*defaultAdapter) ResolveName(raw string) (string, bool) {
	return ResolveName(raw)
}

func (*defaultAdapter) Normalize(style psrt.Style) map[string]json.RawMessage {
	return Normalize(style)
}

func (*defaultAdapter) AdaptHTML(ctx AdaptContext) []StyleFragment {
	return adapt(ctx, true)
}

func (*defaultAdapter) AdaptSVG(ctx AdaptContext) []StyleFragment {
	return adapt(ctx, false)
}

func adapt(ctx AdaptContext, html bool) []StyleFragment {
	style := Normalize(ctx.Text.Style)
	if len(style) == 0 {
		style = map[string]json.RawMessage{}
	}
	dims := ctx.ImageDims()
	dims.FontSizePx = ctx.FontSizePxOrCompute()
	if html && ctx.HTMLCompile {
		dims.PreservePercent = true
	}
	style = percent.ApplyPercentHandlers(style, dims)
	style = FilterStyleMap(style)

	filterID := fmt.Sprintf("psrt-filter-%s-%d", ctx.PageSlug, ctx.TextIndex)
	var effectFrags []StyleFragment
	style, effectFrags = expandEffects(style, dims, !html, filterID)

	var blurMeta BlurAdapt
	var blurFrags []StyleFragment
	style, blurMeta, blurFrags = expandBlur(style, dims, html, filterID)

	var frags []StyleFragment
	if html {
		frags = buildHTMLFragments(ctx, style)
	} else {
		frags = buildSVGFragments(ctx, style, filterID, blurMeta)
	}
	frags = append(frags, effectFrags...)
	frags = append(frags, blurFrags...)
	return MergeFragments(frags)
}

func buildHTMLFragments(ctx AdaptContext, style map[string]json.RawMessage) []StyleFragment {
	box := NewFragment(TypeMotionDiv)
	text := NewFragment(TypeSpan)

	applyLayout(box, ctx, style)
	applyTransform(box, style, true)

	for k, raw := range style {
		if !HasStyleValue(k, raw) {
			continue
		}
		val := SanitizeCSSValue(StringifyCSSValue(raw))
		switch {
		case isBoxKey(k):
			applyBoxCSS(box, k, val)
		case isTextKey(k) && k != KeyStroke && k != KeyStrokeWidth && k != KeyStrokeColor:
			applyTextCSS(text, k, val)
		case isTransformKey(k):
			// merged in applyTransform
		}
	}

	applyStrokeHTML(text, style)
	applyFontSize(box, ctx, ctx.HTMLCompile)
	applyTextAlignHTML(box, text, style)
	if text.GetString(KeyColor) == "" {
		if c := textColor(style); c != "" {
			text.Set(KeyColor, c)
		}
	}
	return []StyleFragment{box, text}
}

func buildSVGFragments(ctx AdaptContext, style map[string]json.RawMessage, filterID string, blurMeta BlurAdapt) []StyleFragment {
	g := NewFragment(TypeG)
	rect := NewFragment(TypeRect)
	fo := NewFragment(TypeForeignObject)
	inner := NewFragment(TypeMotionDiv)

	applyLayoutSVG(g, rect, fo, ctx)
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
		} else if isTextKey(k) && k != KeyStroke && k != KeyStrokeWidth && k != KeyStrokeColor {
			applyTextCSS(inner, k, val)
		}
	}
	applyStrokeSVG(inner, style)
	applyFontSize(inner, ctx, false)

	if blurMeta.FilterID != "" {
		applyBlurSVGRect(rect, blurMeta)
		hasBox = true
	}
	if !hasBox {
		// still emit rect only if we need filter from glow
		if rect.GetString("filter") != "" {
			hasBox = true
		}
	}

	var out []StyleFragment
	out = append(out, g)
	if hasBox || rect.GetString(KeyBackground) != "" || rect.GetString("fill") != "" {
		out = append(out, rect)
	}
	out = append(out, fo, inner)
	return out
}

func applyLayout(box StyleFragment, ctx AdaptContext, style map[string]json.RawMessage) {
	box.Set(KeyPosition, "absolute")
	box.Set(KeyBoxSizing, "border-box")
	box.Set(KeyLeft, pctString(ctx.Text.X))
	box.Set(KeyTop, pctString(ctx.Text.Y))
	box.Set(KeyWidth, pctString(ctx.Text.Width))
}

func applyLayoutSVG(g, rect, fo StyleFragment, ctx AdaptContext) {
	// Geometry in px filled by compilesvg; placeholders for merge.
	x := int(float64(ctx.CanvasW) * ctx.Text.X / 100.0)
	y := int(float64(ctx.CanvasH) * ctx.Text.Y / 100.0)
	w := int(float64(ctx.CanvasW) * ctx.Text.Width / 100.0)
	if w < 1 {
		w = 1
	}
	g.Set("x", x)
	g.Set("y", y)
	rect.Set("x", x)
	rect.Set("y", y)
	rect.Set(KeyWidth, w)
	fo.Set("x", x)
	fo.Set("y", y)
	fo.Set(KeyWidth, w)
}

func applyTransform(target StyleFragment, style map[string]json.RawMessage, html bool) {
	var parts []string
	if raw, ok := style[KeyTransform]; ok {
		if v := StringifyCSSValue(raw); v != "" {
			parts = append(parts, v)
		}
	}
	for _, k := range []string{KeyTranslate, KeyRotate, KeyScale, KeySkew, KeyMatrix} {
		if raw, ok := style[k]; ok {
			if v := StringifyCSSValue(raw); v != "" {
				parts = append(parts, k+"("+v+")")
			}
		}
	}
	if len(parts) == 0 {
		return
	}
	prop := "transform"
	if !html {
		prop = "transform"
	}
	target.Set(prop, strings.Join(parts, " "))
	if raw, ok := style[KeyTransformOrigin]; ok {
		if v := StringifyCSSValue(raw); v != "" {
			target.Set(KeyTransformOrigin, v)
		}
	}
	_ = html
}

func applyBoxCSS(f StyleFragment, key, val string) {
	switch key {
	case KeyBackground:
		f.Set("backgroundColor", val)
	case KeyPadding, KeyPaddingTop, KeyPaddingRight, KeyPaddingBottom, KeyPaddingLeft:
		f.Set(key, val)
	case KeyBorder:
		f.Set(KeyBorder, val)
	case KeyBorderWidth, KeyBorderStyle, KeyBorderColor:
		f.Set(key, val)
	case KeyBorderTop, KeyBorderRight, KeyBorderBottom, KeyBorderLeft:
		f.Set(key, val)
	case KeyBorderRadius, KeyBorderTopLeftRadius, KeyBorderTopRightRadius,
		KeyBorderBottomRightRadius, KeyBorderBottomLeftRadius:
		f.Set(key, val)
	case KeyBoxShadow:
		f.Set(KeyBoxShadow, val)
	case KeyHeight:
		f.Set(KeyHeight, val)
	}
}

func applyBoxSVG(rect StyleFragment, key, val string) bool {
	switch key {
	case KeyBackground:
		rect.Set("fill", val)
		return true
	case KeyBorder, KeyBorderColor:
		rect.Set("stroke", extractBorderColor(val))
		return true
	case KeyBorderWidth:
		rect.Set("stroke-width", val)
		return true
	case KeyBorderRadius:
		rect.Set("rx", firstPx(val))
		return true
	case KeyPadding:
		// padding affects FO geometry in compilesvg, not rect fill
		return false
	case KeyBoxShadow:
		return true
	default:
		return false
	}
}

func extractBorderColor(border string) string {
	parts := strings.Fields(border)
	if len(parts) >= 3 {
		return strings.Join(parts[2:], " ")
	}
	if len(parts) == 2 {
		return parts[1]
	}
	return border
}

func firstPx(s string) string {
	s = strings.TrimSpace(s)
	if i := strings.Index(s, " "); i > 0 {
		s = s[:i]
	}
	return s
}

func applyTextCSS(f StyleFragment, key, val string) {
	cssKey := key
	if key == KeyTextAlign {
		cssKey = "textAlign"
	}
	f.Set(cssKey, val)
}

func applyFontSize(f StyleFragment, ctx AdaptContext, htmlCompile bool) {
	if htmlCompile {
		// Scale with .slide-overlay (container-type:size); cqmin = textSize% of min(W,H).
		f.Set(KeyFontSize, fmt.Sprintf("%gcqmin", ctx.Text.TextSize))
		return
	}
	f.Set(KeyFontSize, pxString(ctx.FontSizePxOrCompute()))
}

func applyTextAlignHTML(box, text StyleFragment, style map[string]json.RawMessage) {
	ta := strings.ToLower(strings.TrimSpace(text.GetString(KeyTextAlign)))
	if ta == "" {
		ta = readTextAlignFromStyle(style)
	}
	va := verticalAlignFromStyle(style)
	if ta == "" && va == "" {
		return
	}
	if ta == "" {
		ta = "left"
	}

	text.Set("display", "block")
	text.Set(KeyWidth, "100%")
	text.Set(KeyTextAlign, ta)

	switch ta {
	case "justify":
		box.Set("display", "block")
	case "center", "right", "left", "start":
		box.Set("display", "flex")
		box.Set("flexDirection", "column")
		jc := va
		if jc == "" {
			jc = "center"
		}
		box.Set("justifyContent", jc)
		ai := "stretch"
		switch ta {
		case "center":
			ai = "center"
		case "right":
			ai = "flex-end"
		case "left", "start":
			ai = "flex-start"
		}
		box.Set("alignItems", ai)
	default:
		box.Set("display", "block")
	}
}

func readTextAlignFromStyle(style map[string]json.RawMessage) string {
	if style == nil {
		return ""
	}
	raw, ok := style[KeyTextAlign]
	if !ok {
		for _, k := range []string{"text-align", "ta"} {
			if r, found := style[k]; found {
				raw = r
				ok = true
				break
			}
		}
	}
	if !ok {
		return ""
	}
	return strings.ToLower(SanitizeCSSValue(StringifyCSSValue(raw)))
}

func verticalAlignFromStyle(style map[string]json.RawMessage) string {
	if style == nil {
		return ""
	}
	raw, ok := style[KeyAlignItems]
	if !ok {
		for _, k := range []string{"align-items", "verticalAlign", "vertical-align"} {
			if r, found := style[k]; found {
				raw = r
				ok = true
				break
			}
		}
	}
	if !ok {
		return ""
	}
	v := strings.ToLower(SanitizeCSSValue(StringifyCSSValue(raw)))
	switch v {
	case "flex-start", "start", "top":
		return "flex-start"
	case "flex-end", "end", "bottom":
		return "flex-end"
	case "center", "middle":
		return "center"
	default:
		return ""
	}
}

func textColor(style map[string]json.RawMessage) string {
	if raw, ok := style[KeyColor]; ok {
		return SanitizeCSSValue(StringifyCSSValue(raw))
	}
	return ""
}

// AdaptHTML is a convenience wrapper around Default.AdaptHTML.
func AdaptHTML(ctx AdaptContext) []StyleFragment {
	return Default.AdaptHTML(ctx)
}

// AdaptSVG is a convenience wrapper around Default.AdaptSVG.
func AdaptSVG(ctx AdaptContext) []StyleFragment {
	return Default.AdaptSVG(ctx)
}

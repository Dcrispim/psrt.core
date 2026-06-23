package styleadapter

import (
	"fmt"
	"sort"
	"strings"
)

// HTMLLayerCSS returns box (layout/background) and text (typography) CSS for compilehtml.
func HTMLLayerCSS(fragments []StyleFragment) (boxCSS, textCSS string) {
	merged := MergeFragments(fragments)
	var box, text StyleFragment
	for _, f := range merged {
		switch f.GetString(TypeKey) {
		case TypeMotionDiv:
			box = f
		case TypeSpan:
			text = f
		}
	}
	postProcessBackdropGlass(box)
	return FragmentCSS(box), FragmentCSS(text)
}

// FragmentsToInlineCSS builds a CSS declaration block from fragments (typically div + span).
// Uses the first matching fragment types in order: div, span, div inner.
func FragmentsToInlineCSS(fragments []StyleFragment) string {
	merged := MergeFragments(fragments)
	var b strings.Builder
	for _, typ := range []string{TypeMotionDiv, TypeSpan} {
		for _, f := range merged {
			if f.GetString(TypeKey) != typ {
				continue
			}
			writeFragmentCSS(&b, f)
		}
	}
	return b.String()
}

// FragmentCSS returns CSS declarations for one fragment (excluding __type__).
func FragmentCSS(f StyleFragment) string {
	var b strings.Builder
	writeFragmentCSS(&b, f)
	return b.String()
}

func writeFragmentCSS(b *strings.Builder, f StyleFragment) {
	if f == nil {
		return
	}
	keys := make([]string, 0, len(f))
	for k := range f {
		if k == TypeKey {
			continue
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		v := f[k]
		s := fmt.Sprint(v)
		if s == "" {
			continue
		}
		cssProp := camelToKebab(k)
		b.WriteString(cssProp)
		b.WriteString(":")
		b.WriteString(s)
		b.WriteString(";")
	}
}

func camelToKebab(s string) string {
	if strings.HasPrefix(s, "-") {
		return s
	}
	var b strings.Builder
	for i, r := range s {
		if r >= 'A' && r <= 'Z' {
			if i > 0 {
				b.WriteByte('-')
			}
			b.WriteRune(r + ('a' - 'A'))
			continue
		}
		b.WriteRune(r)
	}
	return b.String()
}

// FilterFragmentToSVG emits a filter element from a filter fragment.
func FilterFragmentToSVG(f StyleFragment) string {
	if f == nil || f.GetString(TypeKey) != TypeFilter {
		return ""
	}
	id := f.GetString("id")
	if id == "" {
		return ""
	}
	if f.GetString("filterKind") == blurFilterKind {
		return blurFilterFragmentToSVG(f)
	}
	dx := f.GetString("feDropShadowDx")
	dy := f.GetString("feDropShadowDy")
	std := f.GetString("feGaussianBlurStd")
	color := f.GetString("floodColor")
	if std == "" {
		std = "1"
	}
	var b strings.Builder
	fmt.Fprintf(&b, `<filter id="%s" x="-50%%" y="-50%%" width="200%%" height="200%%">`, id)
	fmt.Fprintf(&b, `<feDropShadow dx="%s" dy="%s" stdDeviation="%s" flood-color="%s"/>`,
		dx, dy, std, svgEscape(color))
	if dx2, ok := f["feDropShadowDx2"]; ok {
		c2, _ := f["floodColor2"].(string)
		fmt.Fprintf(&b, `<feDropShadow dx="%v" dy="%s" stdDeviation="0.5" flood-color="%s"/>`,
			dx2, dy, svgEscape(c2))
	}
	b.WriteString(`</filter>`)
	return b.String()
}

func blurFilterFragmentToSVG(f StyleFragment) string {
	id := f.GetString("id")
	std := f.GetString("feGaussianBlurStd")
	if std == "" {
		std = "1"
	}
	in := f.GetString("feGaussianBlurIn")
	if in == "" {
		in = "SourceGraphic"
	}
	var b strings.Builder
	fmt.Fprintf(&b, `<filter id="%s" x="-50%%" y="-50%%" width="200%%" height="200%%">`, id)
	fmt.Fprintf(&b, `<feGaussianBlur in="%s" stdDeviation="%s"/>`, svgEscape(in), std)
	b.WriteString(`</filter>`)
	return b.String()
}

// MaskFragmentToSVG emits a linear-gradient mask for directional box blur.
func MaskFragmentToSVG(f StyleFragment) string {
	if f == nil || f.GetString(TypeKey) != TypeMask {
		return ""
	}
	id := f.GetString("id")
	side := f.GetString("maskSide")
	if id == "" || side == "" {
		return ""
	}
	x1, y1, x2, y2 := maskGradientCoords(side)
	var b strings.Builder
	fmt.Fprintf(&b, `<mask id="%s" maskUnits="objectBoundingBox">`, id)
	fmt.Fprintf(&b, `<linearGradient id="%s-grad" x1="%s" y1="%s" x2="%s" y2="%s">`, id, x1, y1, x2, y2)
	b.WriteString(`<stop offset="0%" stop-color="white" stop-opacity="1"/>`)
	b.WriteString(`<stop offset="100%" stop-color="white" stop-opacity="0"/>`)
	b.WriteString(`</linearGradient>`)
	fmt.Fprintf(&b, `<rect width="100%%" height="100%%" fill="url(#%s-grad)"/></mask>`, id)
	return b.String()
}

func maskGradientCoords(side string) (x1, y1, x2, y2 string) {
	switch strings.ToLower(side) {
	case "left":
		return "0%", "0%", "100%", "0%"
	case "right":
		return "100%", "0%", "0%", "0%"
	case "top":
		return "0%", "0%", "0%", "100%"
	case "bottom":
		return "0%", "100%", "0%", "0%"
	default:
		return "0%", "0%", "100%", "0%"
	}
}

func svgEscape(s string) string {
	return strings.ReplaceAll(s, `"`, "&quot;")
}

// RectFragmentAttrs returns SVG attributes for a rect fragment.
func RectFragmentAttrs(f StyleFragment) string {
	if f == nil || f.GetString(TypeKey) != TypeRect {
		return ""
	}
	return rectAttrs(f, true)
}

// RectDecorationAttrs returns fill/stroke/radius/filter only (position and size come from layout).
func RectDecorationAttrs(f StyleFragment) string {
	if f == nil || f.GetString(TypeKey) != TypeRect {
		return ""
	}
	return rectAttrs(f, false)
}

func rectAttrs(f StyleFragment, includeLayout bool) string {
	var parts []string
	keys := []string{"fill", "stroke", "stroke-width", "rx", "filter", "mask"}
	if includeLayout {
		keys = append([]string{"x", "y", "width", "height"}, keys...)
	}
	for _, k := range keys {
		v, ok := f[k]
		if !ok || v == nil {
			continue
		}
		s := fmt.Sprint(v)
		if s == "" {
			continue
		}
		if k == "filter" || k == "mask" {
			parts = append(parts, fmt.Sprintf(`%s="%s"`, k, svgEscape(s)))
			continue
		}
		if k == "x" || k == "y" || k == "width" || k == "height" || k == "stroke-width" || k == "rx" {
			parts = append(parts, fmt.Sprintf(`%s="%s"`, k, svgEscape(s)))
			continue
		}
		parts = append(parts, fmt.Sprintf(`%s="%s"`, k, svgEscape(s)))
	}
	return strings.Join(parts, " ")
}

package compileasset

import (
	"bytes"
	"encoding/json"
	"math"
	"strconv"
	"strings"

	"github.com/Dcrispim/psrt.core/psrt"
)

// CSSFromStyleJSON converts PSRT style JSON to CSS declarations (no selector).
func CSSFromStyleJSON(raw []byte) string {
	raw = bytes.TrimSpace(raw)
	if len(raw) == 0 || string(raw) == "{}" {
		return ""
	}
	m := normalizeStyleMap(psrt.Style(raw))
	if len(m) == 0 {
		return ""
	}

	ordered := [][]string{
		{"fontFamily", "font-family"},
		{"color", "color"},
		{"fontWeight", "font-weight"},
		{"fontStyle", "font-style"},
		{"textAlign", "text-align"},
		{"lineHeight", "line-height"},
		{"letterSpacing", "letter-spacing"},
		{"textShadow", "text-shadow"},
		{"opacity", "opacity"},
		{"textDecoration", "text-decoration"},
		{"margin", "margin"},
	}

	var b strings.Builder
	for _, kv := range ordered {
		prop := kv[0]
		r, ok := m[prop]
		if !ok {
			continue
		}
		val := stringifyJSONCSSValue(r)
		if val == "" {
			continue
		}
		b.WriteString(kv[1])
		b.WriteString(":")
		b.WriteString(SanitizeCSSValue(val))
		b.WriteString(";")
	}
	return b.String()
}

// BackgroundColorFromStyle extracts a CSS color from page style background keys.
func BackgroundColorFromStyle(style psrt.Style) string {
	return textBoxBackground(normalizeStyleMap(style))
}

// StyleJSONWithoutBackground returns style JSON without box/background keys (typography only).
func StyleJSONWithoutBackground(style psrt.Style) []byte {
	return StyleJSONWithoutBox(style)
}

// FontWeightIsBold reports whether the style requests a bold or semibold weight.
func FontWeightIsBold(style psrt.Style) bool {
	m := normalizeStyleMap(style)
	if m == nil {
		return false
	}
	w := strings.ToLower(rawStringProp(m, "fontWeight", "font-weight"))
	switch w {
	case "bold", "bolder", "600", "700", "800", "900":
		return true
	}
	if n, err := strconv.Atoi(w); err == nil && n >= 600 {
		return true
	}
	return false
}

// TextAlignFromStyle returns the text-align value from style JSON, if any.
func TextAlignFromStyle(style psrt.Style) string {
	m := normalizeStyleMap(style)
	if m == nil {
		return ""
	}
	return rawStringProp(m, "textAlign", "text-align")
}

// VerticalAlignFromStyle maps align-items / vertical-align to a flex justify-content keyword.
func VerticalAlignFromStyle(style psrt.Style) string {
	m := normalizeStyleMap(style)
	if m == nil {
		return ""
	}
	v := strings.ToLower(rawStringProp(m, "alignItems", "align-items", "verticalAlign", "vertical-align"))
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

// TextBlockDisplayCSS returns display/flex rules for SVG foreignObject inner blocks.
// foreignObject renderers often ignore % padding and block text-align without a flex formatting context.
func TextBlockDisplayCSS(style psrt.Style) string {
	ta := strings.ToLower(TextAlignFromStyle(style))
	va := VerticalAlignFromStyle(style)

	if (ta == "left" || ta == "start" || ta == "justify") && va == "" {
		var b strings.Builder
		b.WriteString("display:block;")
		if ta != "" {
			b.WriteString("text-align:")
			b.WriteString(ta)
			b.WriteString(";")
		}
		return b.String()
	}

	if ta == "center" || ta == "right" || ta == "left" || ta == "start" || ta == "justify" || va != "" {
		var b strings.Builder
		writeTextBlockFlexCSS(&b, ta, va)
		return b.String()
	}

	return "display:block;"
}

func writeTextBlockFlexCSS(b *strings.Builder, ta, va string) {
	jc := va
	if jc == "" {
		jc = "center"
	}
	ai := "stretch"
	switch ta {
	case "center":
		ai = "center"
	case "right":
		ai = "flex-end"
	case "left", "start":
		ai = "flex-start"
	}
	b.WriteString("display:flex;flex-direction:column;justify-content:")
	b.WriteString(jc)
	b.WriteString(";align-items:")
	b.WriteString(ai)
	b.WriteString(";")
	if ta != "" {
		b.WriteString("text-align:")
		b.WriteString(ta)
		b.WriteString(";")
	}
}

// LineHeightMultiplier returns a unitless line-height factor for layout (default 1.2, matching compilehtml).
func LineHeightMultiplier(style psrt.Style, fontSizePx float64) float64 {
	const def = 1.2
	m := normalizeStyleMap(style)
	if m == nil {
		return def
	}
	raw, ok := m["lineHeight"]
	if !ok {
		return def
	}
	var f float64
	if err := json.Unmarshal(raw, &f); err == nil && f > 0 {
		return f
	}
	val := strings.TrimSpace(stringifyJSONCSSValue(raw))
	if val == "" {
		return def
	}
	if strings.HasSuffix(val, "%") {
		pct, err := strconv.ParseFloat(strings.TrimSuffix(val, "%"), 64)
		if err == nil && pct > 0 {
			return pct / 100.0
		}
		return def
	}
	if strings.HasSuffix(val, "px") {
		px, err := strconv.ParseFloat(strings.TrimSuffix(val, "px"), 64)
		if err == nil && px > 0 && fontSizePx > 0 {
			return px / fontSizePx
		}
		return def
	}
	if v, err := strconv.ParseFloat(val, 64); err == nil && v > 0 {
		return v
	}
	return def
}

func cssColorFromRaw(raw json.RawMessage) string {
	if len(raw) == 0 || string(raw) == "null" {
		return ""
	}
	var s string
	if err := json.Unmarshal(raw, &s); err == nil && strings.TrimSpace(s) != "" {
		return SanitizeCSSValue(s)
	}
	return ""
}

func stringifyJSONCSSValue(r json.RawMessage) string {
	if len(r) == 0 || string(r) == "null" {
		return ""
	}
	var s string
	if err := json.Unmarshal(r, &s); err == nil {
		return strings.TrimSpace(s)
	}
	var f float64
	if err := json.Unmarshal(r, &f); err == nil {
		return formatJSONNumberForCSS(f)
	}
	var b bool
	if err := json.Unmarshal(r, &b); err == nil {
		if b {
			return "1"
		}
		return "0"
	}
	return strings.TrimSpace(string(r))
}

func formatJSONNumberForCSS(f float64) string {
	if math.Abs(f-float64(int64(f))) < 1e-9 && math.Abs(f) < 1e12 {
		return strconv.FormatInt(int64(f), 10)
	}
	return strconv.FormatFloat(f, 'f', -1, 64)
}

// SanitizeCSSValue strips characters that break CSS declarations.
func SanitizeCSSValue(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}
	s = strings.ReplaceAll(s, ";", "")
	return s
}

// FormatFloatCSS formats a float for use in CSS.
func FormatFloatCSS(v float64) string {
	if math.IsNaN(v) || math.IsInf(v, 0) {
		return "0"
	}
	return strconv.FormatFloat(v, 'f', -1, 64)
}

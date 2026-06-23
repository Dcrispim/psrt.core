package styleadapter

import (
	"encoding/json"
	"math"
	"strconv"
	"strings"
)

func StringifyCSSValue(r json.RawMessage) string {
	if len(r) == 0 || string(r) == "null" {
		return ""
	}
	var s string
	if err := json.Unmarshal(r, &s); err == nil {
		return strings.TrimSpace(s)
	}
	var f float64
	if err := json.Unmarshal(r, &f); err == nil {
		return formatJSONNumber(f)
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

func formatJSONNumber(f float64) string {
	if math.Abs(f-float64(int64(f))) < 1e-9 && math.Abs(f) < 1e12 {
		return strconv.FormatInt(int64(f), 10)
	}
	return strconv.FormatFloat(f, 'f', -1, 64)
}

func SanitizeCSSValue(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}
	return strings.ReplaceAll(s, ";", "")
}

// HasStyleValue reports whether a style property was explicitly set to a usable CSS value.
// Absent, null, empty, false, and zero-like dimensions are treated as not passed.
func HasStyleValue(key string, raw json.RawMessage) bool {
	if len(raw) == 0 || string(raw) == "null" {
		return false
	}
	var b bool
	if err := json.Unmarshal(raw, &b); err == nil {
		return b
	}
	val := SanitizeCSSValue(StringifyCSSValue(raw))
	if val == "" {
		return false
	}
	if omitZeroForKey(key) && isZeroLikeCSSValue(val) {
		return false
	}
	return true
}

func omitZeroForKey(key string) bool {
	switch key {
	case KeyHeight, KeyWidth,
		KeyPadding, KeyPaddingTop, KeyPaddingRight, KeyPaddingBottom, KeyPaddingLeft,
		KeyBorderWidth, KeyStrokeWidth,
		KeyBorderRadius, KeyBorderTopLeftRadius, KeyBorderTopRightRadius,
		KeyBorderBottomRightRadius, KeyBorderBottomLeftRadius,
		KeyLetterSpacing, KeyWordSpacing, KeyLineHeight, KeyTextIndent,
		KeyBlur, KeyBlurLeft, KeyBlurRight, KeyBlurTop, KeyBlurBottom:
		return true
	default:
		return false
	}
}

func isZeroLikeCSSValue(s string) bool {
	s = strings.TrimSpace(strings.ToLower(s))
	switch s {
	case "0", "0%", "0px", "0em", "0rem", "0pt", "0cqh", "0ch", "0vw", "0vh", "0vmin", "0vmax":
		return true
	}
	for _, unit := range []string{"px", "%", "em", "rem", "pt", "cqh", "ch"} {
		if strings.HasSuffix(s, unit) {
			n := strings.TrimSuffix(s, unit)
			if f, err := strconv.ParseFloat(n, 64); err == nil && f == 0 {
				return true
			}
		}
	}
	if f, err := strconv.ParseFloat(s, 64); err == nil && f == 0 {
		return true
	}
	return false
}

// FilterStyleMap removes keys that were not effectively passed.
func FilterStyleMap(m map[string]json.RawMessage) map[string]json.RawMessage {
	if len(m) == 0 {
		return m
	}
	out := make(map[string]json.RawMessage, len(m))
	for k, raw := range m {
		if HasStyleValue(k, raw) {
			out[k] = raw
		}
	}
	return out
}

func pctString(v float64) string {
	return strconv.FormatFloat(v, 'f', -1, 64) + "%"
}

func pxString(v float64) string {
	return strconv.FormatFloat(v, 'f', -1, 64) + "px"
}

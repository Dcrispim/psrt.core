package styleadapter

import (
	"encoding/json"
	"strings"
	"unicode"

	"psrt/psrt"
)

var aliasToCanonical = map[string]string{
	// kebab (stored lower in lookup after camel conversion path)
	"border-radius":     KeyBorderRadius,
	"border-width":      KeyBorderWidth,
	"border-style":      KeyBorderStyle,
	"border-color":      KeyBorderColor,
	"box-shadow":        KeyBoxShadow,
	"text-shadow":       KeyTextShadow,
	"text-align":        KeyTextAlign,
	"align-items":       KeyAlignItems,
	"vertical-align":    KeyAlignItems,
	"font-family":       KeyFontFamily,
	"font-weight":       KeyFontWeight,
	"font-style":        KeyFontStyle,
	"line-height":       KeyLineHeight,
	"letter-spacing":    KeyLetterSpacing,
	"text-decoration":   KeyTextDecoration,
	"background-color":  KeyBackground,
	"stroke-width":      KeyStrokeWidth,
	"stroke-color":      KeyStrokeColor,
	"text-stroke":       KeyStroke,
	// legacy PSRT
	"backGround": KeyBackground,
	"backgroundColor": KeyBackground,
	// siglas
	"br":  KeyBorderRadius,
	"bw":  KeyBorderWidth,
	"bc":  KeyBorderColor,
	"bs":  KeyBorderStyle,
	"bg":  KeyBackground,
	"ta":  KeyTextAlign,
	"ts":  KeyTextShadow,
	"bsh": KeyBoxShadow,
	"blur":            KeyBlur,
	"blur-left":       KeyBlurLeft,
	"blur-right":      KeyBlurRight,
	"blur-top":        KeyBlurTop,
	"blur-bottom":     KeyBlurBottom,
	"ff":  KeyFontFamily,
	"fw":  KeyFontWeight,
	"fs":  KeyFontStyle,
	"pd":  KeyPadding,
	"sw":  KeyStrokeWidth,
	"sc":  KeyStrokeColor,
	// webkit legacy → canonical (after strip)
	"textStroke":      KeyStroke,
	"textStrokeWidth": KeyStrokeWidth,
	"textStrokeColor": KeyStrokeColor,
}

// RegisterAlias adds or overrides an alias → canonical mapping.
func RegisterAlias(alias, canonical string) {
	aliasToCanonical[strings.TrimSpace(alias)] = canonical
}

// ResolveName maps raw property names to canonical PSRT camelCase keys.
func ResolveName(raw string) (canonical string, ok bool) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return "", false
	}
	if c, exists := aliasToCanonical[raw]; exists {
		return c, true
	}
	if strings.Contains(raw, "-") {
		camel := kebabToCamel(raw)
		if c, exists := aliasToCanonical[camel]; exists {
			return c, true
		}
		if c, exists := aliasToCanonical[raw]; exists {
			return c, true
		}
		return camel, isKnownCanonical(camel)
	}
	if c, exists := aliasToCanonical[raw]; exists {
		return c, true
	}
	return raw, isKnownCanonical(raw)
}

func isKnownCanonical(name string) bool {
	return isBoxKey(name) || isTextKey(name) || isTransformKey(name) ||
		name == KeyLeft || name == KeyTop || name == KeyWidth || name == KeyHeight ||
		name == KeyGlow || name == KeyBevel ||
		name == KeyBlur || name == KeyBlurLeft || name == KeyBlurRight ||
		name == KeyBlurTop || name == KeyBlurBottom
}

func kebabToCamel(s string) string {
	parts := strings.Split(s, "-")
	if len(parts) == 0 {
		return s
	}
	var b strings.Builder
	for i, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		if i == 0 {
			b.WriteString(strings.ToLower(p))
			continue
		}
		if len(p) > 0 {
			b.WriteString(strings.ToUpper(p[:1]))
			if len(p) > 1 {
				b.WriteString(strings.ToLower(p[1:]))
			}
		}
	}
	return b.String()
}

// Normalize parses style JSON into a canonical-key map (no Webkit* keys).
func Normalize(style psrt.Style) map[string]json.RawMessage {
	raw := bytesTrim(style)
	if len(raw) == 0 || string(raw) == "{}" {
		return nil
	}
	var m map[string]json.RawMessage
	if err := json.Unmarshal(raw, &m); err != nil {
		return nil
	}
	out := make(map[string]json.RawMessage)
	applyKeys := func(vendorOnly bool) {
		for rawKey, val := range m {
			vendor := isVendorRawKey(rawKey)
			if vendorOnly != vendor {
				continue
			}
			canonical := canonicalKeyForRaw(rawKey)
			if canonical == "" || isVendorKey(canonical) {
				continue
			}
			if _, exists := out[canonical]; exists && vendorOnly {
				continue
			}
			out[canonical] = val
		}
	}
	applyKeys(false) // pure / non-vendor first
	applyKeys(true)  // vendor only fills gaps
	mergeBackgroundKeys(out)
	return FilterStyleMap(out)
}

func canonicalKeyForRaw(rawKey string) string {
	for _, cand := range expandRawKey(rawKey) {
		if canonical, ok := ResolveName(cand); ok {
			return canonical
		}
		if !isVendorKey(cand) {
			return cand
		}
	}
	return ""
}

func isVendorRawKey(key string) bool {
	if isVendorKey(key) {
		return true
	}
	_, ok := stripVendorPrefix(key)
	return ok
}

func bytesTrim(style psrt.Style) []byte {
	return []byte(strings.TrimSpace(string(style)))
}

func expandRawKey(rawKey string) []string {
	rawKey = strings.TrimSpace(rawKey)
	if stripped, ok := stripVendorPrefix(rawKey); ok {
		return []string{stripped, rawKey}
	}
	return []string{rawKey}
}

func stripVendorPrefix(key string) (string, bool) {
	key = strings.TrimSpace(key)
	lower := strings.ToLower(key)
	prefixes := []string{"-webkit-", "webkit-", "webkit"}
	for _, p := range prefixes {
		if strings.HasPrefix(lower, p) {
			rest := key[len(p):]
			if rest == "" {
				return "", false
			}
			return decapitalize(rest), true
		}
	}
	if strings.HasPrefix(key, "Webkit") {
		return decapitalize(key[6:]), true
	}
	if strings.HasPrefix(key, "webkit") {
		return decapitalize(key[6:]), true
	}
	return "", false
}

func decapitalize(s string) string {
	if s == "" {
		return s
	}
	runes := []rune(s)
	runes[0] = unicode.ToLower(runes[0])
	return string(runes)
}

func isVendorKey(k string) bool {
	return strings.HasPrefix(k, "Webkit") || strings.HasPrefix(k, "webkit") ||
		strings.HasPrefix(strings.ToLower(k), "-webkit-")
}

func mergeBackgroundKeys(m map[string]json.RawMessage) {
	// Prefer explicit background over backgroundColor/backGround already merged via aliases.
	if bg, ok := m[KeyBackground]; ok {
		m[KeyBackground] = bg
		delete(m, "backgroundColor")
	}
}

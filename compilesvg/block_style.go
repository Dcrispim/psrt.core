package compilesvg

import (
	"encoding/json"
	"math"
	"strings"

	"psrt/compileasset"
	"psrt/compilesvg/textoutline"
	"psrt/psrt"
)

func blockStyleForText(t *psrt.Text, canvasW, canvasH int, fontURLs []string, assets map[string]compileasset.Asset) textoutline.BlockStyle {
	fontPx := psrt.TextFontSizePx(t.TextSize, canvasW, canvasH)
	insets := compileasset.TextBoxInsetsForCanvas(t.Style, fontPx, canvasW, canvasH)
	outerW := textBlockWidthPx(t.Width, canvasW)
	contentW := outerW - int(math.Round(insets.Horizontal()))
	if contentW < 1 {
		contentW = 1
	}
	return textoutline.BlockStyle{
		FontSizePx:  fontPx,
		LineHeight:  compileasset.LineHeightMultiplier(t.Style, fontPx),
		Color:       colorFromStyle(t.Style),
		Stroke:      strokeColorFromStyle(t.Style),
		StrokeWidth: strokeWidthFromStyle(t.Style),
		TextAlign:   compileasset.TextAlignFromStyle(t.Style),
		PadTop:      int(math.Round(insets.Top)),
		PadLeft:     int(math.Round(insets.Left)),
		ContentW:    contentW,
		FontFamily:  firstFontFamily(t.Style, fontURLs, assets),
		Bold:        compileasset.FontWeightIsBold(t.Style),
	}
}

func colorFromStyle(style psrt.Style) string {
	m := normalizeStyleMap(style)
	if m == nil {
		return "#000000"
	}
	if raw, ok := m["color"]; ok {
		if v := stringifyJSONCSSValue(raw); v != "" {
			return compileasset.SanitizeCSSValue(v)
		}
	}
	return "#000000"
}

func strokeColorFromStyle(style psrt.Style) string {
	m := normalizeStyleMap(style)
	if m == nil {
		return ""
	}
	for _, k := range []string{"webkitTextStrokeColor", "WebkitTextStrokeColor", "strokeColor", "stroke-color"} {
		if raw, ok := m[k]; ok {
			if v := stringifyJSONCSSValue(raw); v != "" {
				return compileasset.SanitizeCSSValue(v)
			}
		}
	}
	return ""
}

func strokeWidthFromStyle(style psrt.Style) string {
	m := normalizeStyleMap(style)
	if m == nil {
		return ""
	}
	for _, k := range []string{"webkitTextStrokeWidth", "WebkitTextStrokeWidth", "strokeWidth", "stroke-width"} {
		if raw, ok := m[k]; ok {
			if v := stringifyJSONCSSValue(raw); v != "" {
				return compileasset.SanitizeCSSValue(v)
			}
		}
	}
	return ""
}

func firstFontFamily(style psrt.Style, fontURLs []string, assets map[string]compileasset.Asset) string {
	m := normalizeStyleMap(style)
	if m != nil {
		if raw, ok := m["fontFamily"]; ok {
			if v := stringifyJSONCSSValue(raw); v != "" {
				if name := firstFamilyName(v); name != "" {
					return name
				}
			}
		}
		if raw, ok := m["font-family"]; ok {
			if v := stringifyJSONCSSValue(raw); v != "" {
				if name := firstFamilyName(v); name != "" {
					return name
				}
			}
		}
	}
	for i, u := range fontURLs {
		u = strings.TrimSpace(u)
		if _, ok := assets[u]; ok {
			return compileasset.FontFamilyNameForURL(u, i)
		}
	}
	return ""
}

func firstFamilyName(stack string) string {
	parts := strings.Split(stack, ",")
	if len(parts) == 0 {
		return ""
	}
	return strings.Trim(parts[0], `"' `)
}

func normalizeStyleMap(style psrt.Style) map[string]json.RawMessage {
	raw := compileasset.StyleJSONWithoutBox(style)
	if len(raw) == 0 || string(raw) == "{}" {
		return nil
	}
	var m map[string]json.RawMessage
	if json.Unmarshal(raw, &m) != nil {
		return nil
	}
	return m
}

func stringifyJSONCSSValue(raw json.RawMessage) string {
	s := strings.TrimSpace(string(raw))
	if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
		var v string
		if json.Unmarshal(raw, &v) == nil {
			return v
		}
	}
	return s
}

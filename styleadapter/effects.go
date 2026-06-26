package styleadapter

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Dcrispim/psrt.core/styleadapter/percent"
)

// expandEffects merges glow/bevel into canonical shadow keys on a copy of style.
func expandEffects(style map[string]json.RawMessage, dims percent.ImageDims, targetSVG bool, filterID string) (map[string]json.RawMessage, []StyleFragment) {
	if len(style) == 0 {
		return style, nil
	}
	out := copyStyleMap(style)
	var extra []StyleFragment

	if raw, ok := out[KeyGlow]; ok {
		val := StringifyCSSValue(raw)
		delete(out, KeyGlow)
		if HasStyleValue(KeyGlow, raw) {
			if targetSVG {
				extra = append(extra, glowFilterFragment(filterID, val, dims))
			} else {
				mergeShadowKey(out, KeyTextShadow, val)
			}
		}
	}
	if raw, ok := out[KeyBevel]; ok {
		val := StringifyCSSValue(raw)
		delete(out, KeyBevel)
		if HasStyleValue(KeyBevel, raw) {
			if targetSVG {
				extra = append(extra, bevelFilterFragment(filterID+"-bevel", val, dims))
			} else {
				mergeBoxShadowKey(out, bevelBoxShadows(val))
			}
		}
	}
	return out, extra
}

func copyStyleMap(m map[string]json.RawMessage) map[string]json.RawMessage {
	out := make(map[string]json.RawMessage, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}

func mergeShadowKey(m map[string]json.RawMessage, key, add string) {
	existing := StringifyCSSValue(m[key])
	if existing == "" {
		m[key] = json.RawMessage(jsonQuote(add))
		return
	}
	m[key] = json.RawMessage(jsonQuote(existing + ", " + add))
}

func mergeBoxShadowKey(m map[string]json.RawMessage, add string) {
	existing := StringifyCSSValue(m[KeyBoxShadow])
	if existing == "" {
		m[KeyBoxShadow] = json.RawMessage(jsonQuote(add))
		return
	}
	m[KeyBoxShadow] = json.RawMessage(jsonQuote(existing + ", " + add))
}

func jsonQuote(s string) []byte {
	b, _ := json.Marshal(s)
	return b
}

func bevelBoxShadows(val string) string {
	light := "inset 1px 1px 0 rgba(255,255,255,0.35)"
	dark := "inset -1px -1px 0 rgba(0,0,0,0.35)"
	if strings.TrimSpace(val) != "" {
		dark = "inset -1px -1px 2px " + strings.TrimSpace(val)
	}
	return light + ", " + dark
}

func glowFilterFragment(id, val string, dims percent.ImageDims) StyleFragment {
	dx, dy, blur, color := parseSimpleShadow(val, dims)
	f := NewFragment(TypeFilter)
	f.Set("id", id)
	f.Set("feDropShadowDx", fmt.Sprintf("%.3f", dx))
	f.Set("feDropShadowDy", fmt.Sprintf("%.3f", dy))
	f.Set("feGaussianBlurStd", fmt.Sprintf("%.3f", blur/2))
	f.Set("floodColor", color)
	return f
}

func bevelFilterFragment(id, val string, _ percent.ImageDims) StyleFragment {
	f := NewFragment(TypeFilter)
	f.Set("id", id)
	f.Set("feDropShadowDx", "-1")
	f.Set("feDropShadowDy", "-1")
	f.Set("feGaussianBlurStd", "0.5")
	color := strings.TrimSpace(val)
	if color == "" {
		color = "rgba(0,0,0,0.4)"
	}
	f.Set("floodColor", color)
	f.Set("feDropShadowDx2", "1")
	f.Set("floodColor2", "rgba(255,255,255,0.35)")
	return f
}

func parseSimpleShadow(s string, dims percent.ImageDims) (dx, dy, blur float64, color string) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, 0, 4, "rgba(0,0,0,0.5)"
	}
	h := percent.ApplyPercentHandlers(
		map[string]json.RawMessage{"textShadow": json.RawMessage(jsonQuote(s))},
		dims,
	)
	resolved := StringifyCSSValue(h["textShadow"])
	parts := strings.Fields(resolved)
	if len(parts) >= 3 {
		dx = parsePxNum(parts[0])
		dy = parsePxNum(parts[1])
		blur = parsePxNum(parts[2])
	}
	if len(parts) >= 4 {
		color = parts[3]
	}
	if color == "" {
		color = "rgba(0,0,0,0.5)"
	}
	return dx, dy, blur, color
}

func parsePxNum(s string) float64 {
	s = strings.TrimSuffix(strings.TrimSpace(s), "px")
	var f float64
	fmt.Sscanf(s, "%f", &f)
	return f
}

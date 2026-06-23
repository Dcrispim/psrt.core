package percent

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// resolveShadowList converts % tokens in shadow lists (textShadow/boxShadow).
func resolveShadowList(value string, dims ImageDims) string {
	parts := strings.Split(value, ",")
	for i, shadow := range parts {
		parts[i] = resolveShadowOne(strings.TrimSpace(shadow), dims)
	}
	return strings.Join(parts, ", ")
}

func resolveShadowOne(shadow string, dims ImageDims) string {
	if shadow == "" {
		return shadow
	}
	tokens := strings.Fields(shadow)
	if len(tokens) == 0 {
		return shadow
	}
	colorIdx := findColorTokenIndex(tokens)
	numEnd := colorIdx
	if numEnd < 0 {
		numEnd = len(tokens)
	}
	for j := 0; j < numEnd && j < 3; j++ {
		tokens[j] = resolveShadowToken(tokens[j], j, dims)
	}
	return strings.Join(tokens, " ")
}

func findColorTokenIndex(tokens []string) int {
	for i, t := range tokens {
		if strings.HasPrefix(t, "#") || strings.HasPrefix(strings.ToLower(t), "rgb") {
			return i
		}
	}
	return -1
}

func resolveShadowToken(token string, index int, dims ImageDims) string {
	if !strings.HasSuffix(token, "%") {
		return token
	}
	pct, err := strconv.ParseFloat(strings.TrimSuffix(token, "%"), 64)
	if err != nil {
		return token
	}
	z := dims.Zoom
	if z <= 0 {
		z = 1
	}
	var base float64
	switch index {
	case 0:
		base = float64(dims.W) * z
	case 1:
		base = float64(dims.H) * z
	case 2:
		base = float64(max(dims.W, dims.H)) * z
	default:
		return token
	}
	return fmt.Sprintf("%.3fpx", (pct/100.0)*base)
}

type textShadowHandler struct{}

func (textShadowHandler) Keys() []string { return []string{"textShadow"} }

func (textShadowHandler) Resolve(_ string, value string, dims ImageDims) (string, bool) {
	if !strings.Contains(value, "%") {
		return value, true
	}
	return resolveShadowList(value, dims), true
}

type boxShadowHandler struct{}

func (boxShadowHandler) Keys() []string { return []string{"boxShadow"} }

func (boxShadowHandler) Resolve(_ string, value string, dims ImageDims) (string, bool) {
	if !strings.Contains(value, "%") {
		return value, true
	}
	return resolveShadowList(value, dims), true
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func percentToPx(pctStr string, base float64, zoom float64) string {
	pct, err := strconv.ParseFloat(strings.TrimSuffix(strings.TrimSpace(pctStr), "%"), 64)
	if err != nil {
		return pctStr
	}
	if zoom <= 0 {
		zoom = 1
	}
	v := (pct / 100.0) * base * zoom
	return fmt.Sprintf("%.3fpx", v)
}

func singlePercentToken(value string, base float64, dims ImageDims) (string, bool) {
	v := strings.TrimSpace(value)
	if !strings.HasSuffix(v, "%") {
		return v, true
	}
	z := dims.Zoom
	if z <= 0 {
		z = 1
	}
	return percentToPx(v, base, z), true
}

func roundPx(v float64) string {
	if math.IsNaN(v) || math.IsInf(v, 0) {
		return "0px"
	}
	return fmt.Sprintf("%.3fpx", v)
}

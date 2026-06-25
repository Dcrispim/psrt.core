package compileasset

import (
	"math"
	"strconv"
	"strings"

	"github.com/Dcrispim/psrt.core/psrt"
)

// ExplicitHeightPx returns the outer box height in pixels from style "height" when set.
func ExplicitHeightPx(style psrt.Style, canvasW, canvasH int, fontPx float64) (int, bool) {
	if canvasH < 1 {
		return 0, false
	}
	resolved := StyleResolvedForCanvas(style, canvasW, canvasH, fontPx)
	m := normalizeStyleMap(resolved)
	if len(m) == 0 {
		return 0, false
	}
	val := strings.TrimSpace(rawStringProp(m, "height"))
	if val == "" {
		return 0, false
	}
	px, ok := cssVerticalLengthToPx(val, float64(canvasW), float64(canvasH), fontPx)
	if !ok || px < 1 {
		return 0, false
	}
	return int(math.Round(px)), true
}

func cssVerticalLengthToPx(val string, canvasW, canvasH, fontPx float64) (float64, bool) {
	val = strings.TrimSpace(val)
	if val == "" {
		return 0, false
	}
	if strings.HasSuffix(val, "%") {
		pct, err := strconv.ParseFloat(strings.TrimSuffix(val, "%"), 64)
		if err != nil || pct <= 0 {
			return 0, false
		}
		return canvasH * pct / 100.0, true
	}
	if strings.HasSuffix(val, "px") {
		px, err := strconv.ParseFloat(strings.TrimSuffix(val, "px"), 64)
		if err != nil || px <= 0 {
			return 0, false
		}
		return px, true
	}
	if strings.HasSuffix(val, "em") || strings.HasSuffix(val, "rem") {
		suffix := "em"
		if strings.HasSuffix(val, "rem") {
			suffix = "rem"
		}
		n, err := strconv.ParseFloat(strings.TrimSuffix(val, suffix), 64)
		if err != nil || n <= 0 || fontPx <= 0 {
			return 0, false
		}
		return n * fontPx, true
	}
	if n, err := strconv.ParseFloat(val, 64); err == nil && n > 0 {
		return n, true
	}
	_ = canvasW
	return 0, false
}

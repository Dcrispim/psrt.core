package psrt

import (
	"math"
	"strconv"
	"strings"
)

const coordMaxDecimals = 5

// RoundCoord limits a coordinate to at most coordMaxDecimals decimal places.
func RoundCoord(v float64) float64 {
	pow := math.Pow10(coordMaxDecimals)
	return math.Round(v*pow) / pow
}

// TextSizeBasisPx is the reference length for text-size percentages: min(canvas width, canvas height).
func TextSizeBasisPx(canvasW, canvasH int) float64 {
	w := float64(canvasW)
	h := float64(canvasH)
	if w <= 0 && h <= 0 {
		return 1
	}
	if w <= 0 {
		return h
	}
	if h <= 0 {
		return w
	}
	if w < h {
		return w
	}
	return h
}

// TextFontSizePx converts text-size percent to pixels using TextSizeBasisPx.
func TextFontSizePx(textSizePct float64, canvasW, canvasH int) float64 {
	px := textSizePct / 100.0 * TextSizeBasisPx(canvasW, canvasH)
	if px < 1 {
		return 1
	}
	return px
}

func formatCoord(v float64) string {
	v = RoundCoord(v)
	s := strconv.FormatFloat(v, 'f', coordMaxDecimals, 64)
	s = strings.TrimRight(s, "0")
	s = strings.TrimRight(s, ".")
	if s == "" || s == "-" {
		return "0"
	}
	return s
}

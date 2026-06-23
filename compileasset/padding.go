package compileasset

import (
	"math"
	"strconv"
	"strings"

	"psrt/psrt"
)

// PaddingInsets holds resolved padding in pixels (top, right, bottom, left).
type PaddingInsets struct {
	Top, Right, Bottom, Left float64
}

func (p PaddingInsets) Horizontal() float64 { return p.Left + p.Right }
func (p PaddingInsets) Vertical() float64   { return p.Top + p.Bottom }

func (p PaddingInsets) Add(o PaddingInsets) PaddingInsets {
	return PaddingInsets{
		Top:    p.Top + o.Top,
		Right:  p.Right + o.Right,
		Bottom: p.Bottom + o.Bottom,
		Left:   p.Left + o.Left,
	}
}

// ParsePaddingInsets resolves padding from PSRT style JSON (px, em).
func ParsePaddingInsets(style psrt.Style, refFontPx float64) PaddingInsets {
	box := ParseTextBox(style)
	return parsePaddingCSS(box.PaddingCSS, refFontPx)
}

// ParseBorderInsets resolves uniform border width in pixels for layout expansion.
func ParseBorderInsets(style psrt.Style, refFontPx float64) PaddingInsets {
	box := ParseTextBox(style)
	w := box.BorderWidth
	if w == "" && box.BorderCSS != "" {
		parts := strings.Fields(box.BorderCSS)
		if len(parts) > 0 {
			w = parts[0]
		}
	}
	px := cssLengthToPx(w, refFontPx)
	if px <= 0 {
		return PaddingInsets{}
	}
	return PaddingInsets{Top: px, Right: px, Bottom: px, Left: px}
}

// TextBoxInsets combines padding and border for foreignObject sizing.
func TextBoxInsets(style psrt.Style, refFontPx float64) PaddingInsets {
	return ParsePaddingInsets(style, refFontPx).Add(ParseBorderInsets(style, refFontPx))
}

// TextBoxInsetsForCanvas resolves % padding/border against canvas size, then returns pixel insets.
func TextBoxInsetsForCanvas(style psrt.Style, refFontPx float64, canvasW, canvasH int) PaddingInsets {
	return TextBoxInsets(StyleResolvedForCanvas(style, canvasW, canvasH, refFontPx), refFontPx)
}

func parsePaddingCSS(css string, refFontPx float64) PaddingInsets {
	css = strings.TrimSpace(css)
	if css == "" {
		return PaddingInsets{}
	}
	parts := strings.Fields(css)
	vals := make([]float64, 0, 4)
	for _, p := range parts {
		vals = append(vals, cssLengthToPx(p, refFontPx))
	}
	switch len(vals) {
	case 0:
		return PaddingInsets{}
	case 1:
		return PaddingInsets{vals[0], vals[0], vals[0], vals[0]}
	case 2:
		return PaddingInsets{vals[0], vals[1], vals[0], vals[1]}
	case 3:
		return PaddingInsets{vals[0], vals[1], vals[2], vals[1]}
	default:
		return PaddingInsets{vals[0], vals[1], vals[2], vals[3]}
	}
}

func cssLengthToPx(s string, refFontPx float64) float64 {
	s = strings.TrimSpace(strings.ToLower(s))
	if s == "" {
		return 0
	}
	if strings.HasSuffix(s, "px") {
		v, err := strconv.ParseFloat(strings.TrimSuffix(s, "px"), 64)
		if err != nil {
			return 0
		}
		return math.Max(0, v)
	}
	if strings.HasSuffix(s, "em") {
		v, err := strconv.ParseFloat(strings.TrimSuffix(s, "em"), 64)
		if err != nil {
			return 0
		}
		return math.Max(0, v*refFontPx)
	}
	if strings.HasSuffix(s, "%") {
		return 0
	}
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return math.Max(0, v)
}

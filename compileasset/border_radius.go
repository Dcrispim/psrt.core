package compileasset

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/Dcrispim/psrt.core/psrt"
)

// BorderRadiusCorners holds four corner radii in pixels (after % resolution when applicable).
type BorderRadiusCorners struct {
	TopLeft, TopRight, BottomRight, BottomLeft float64
}

// ParseBorderRadiusCorners reads border-radius from style (resolved for canvas when needed).
func ParseBorderRadiusCorners(style psrt.Style, canvasW, canvasH int, fontPx float64) BorderRadiusCorners {
	cw, ch, fp := float64(canvasW), float64(canvasH), fontPx
	if c, ok := borderRadiusFromLonghands(normalizeStyleMap(style), cw, ch, fp); ok {
		return c
	}
	resolved := StyleResolvedForCanvas(style, canvasW, canvasH, fontPx)
	m := normalizeStyleMap(resolved)
	if c, ok := borderRadiusFromLonghands(m, cw, ch, fp); ok {
		return c
	}
	raw := rawStringProp(m, "borderRadius", "border-radius")
	return parseBorderRadiusCSSValue(raw, cw, ch, fp)
}

func borderRadiusFromLonghands(m map[string]json.RawMessage, canvasW, canvasH, fontPx float64) (BorderRadiusCorners, bool) {
	read := func(keys ...string) (float64, bool) {
		var rawJSON json.RawMessage
		for _, key := range keys {
			if r, ok := m[key]; ok && len(r) > 0 && string(r) != "null" {
				rawJSON = r
				break
			}
		}
		if len(rawJSON) == 0 {
			return 0, false
		}
		raw := strings.TrimSpace(stringifyJSONCSSValue(rawJSON))
		if raw == "" {
			return 0, true
		}
		v, ok := cssVerticalLengthToPx(raw, canvasW, canvasH, fontPx)
		return v, ok
	}
	tl, hasTL := read("borderTopLeftRadius", "border-top-left-radius")
	tr, hasTR := read("borderTopRightRadius", "border-top-right-radius")
	br, hasBR := read("borderBottomRightRadius", "border-bottom-right-radius")
	bl, hasBL := read("borderBottomLeftRadius", "border-bottom-left-radius")
	if !hasTL && !hasTR && !hasBR && !hasBL {
		return BorderRadiusCorners{}, false
	}
	return BorderRadiusCorners{
		TopLeft:      tl,
		TopRight:     tr,
		BottomRight:  br,
		BottomLeft:   bl,
	}, true
}

func parseBorderRadiusCSSValue(raw string, canvasW, canvasH, fontPx float64) BorderRadiusCorners {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return BorderRadiusCorners{}
	}
	tokens := strings.Fields(raw)
	if len(tokens) == 0 {
		return BorderRadiusCorners{}
	}
	px := make([]float64, 0, len(tokens))
	for _, tok := range tokens {
		if v, ok := cssVerticalLengthToPx(tok, canvasW, canvasH, fontPx); ok {
			px = append(px, v)
		}
	}
	if len(px) == 0 {
		return BorderRadiusCorners{}
	}
	var c BorderRadiusCorners
	switch len(px) {
	case 1:
		c = BorderRadiusCorners{px[0], px[0], px[0], px[0]}
	case 2:
		c = BorderRadiusCorners{px[0], px[1], px[0], px[1]}
	case 3:
		c = BorderRadiusCorners{px[0], px[1], px[2], px[1]}
	default:
		c = BorderRadiusCorners{px[0], px[1], px[2], px[3]}
	}
	return c
}

// BorderRadiusIsUniform reports whether all corners share the same radius.
func BorderRadiusIsUniform(c BorderRadiusCorners) bool {
	return c.TopLeft == c.TopRight &&
		c.TopRight == c.BottomRight &&
		c.BottomRight == c.BottomLeft
}

// BorderRadiusHasValue reports whether any corner has a positive radius.
func BorderRadiusHasValue(c BorderRadiusCorners) bool {
	return c.TopLeft > 0 || c.TopRight > 0 || c.BottomRight > 0 || c.BottomLeft > 0
}

// BorderRadiusNeedsInnerBoxDecoration is true when the SVG background rect cannot represent
// the corner radii (e.g. only bottom corners rounded). In that case background + radius
// must live on the foreignObject inner div via CSS.
func BorderRadiusNeedsInnerBoxDecoration(c BorderRadiusCorners) bool {
	if !BorderRadiusHasValue(c) {
		return false
	}
	if !BorderRadiusIsUniform(c) {
		return true
	}
	return false
}

// BorderRadiusCSS returns border-radius shorthand or longhands for CSS.
func BorderRadiusCSS(c BorderRadiusCorners) string {
	if !BorderRadiusHasValue(c) {
		return ""
	}
	if BorderRadiusIsUniform(c) {
		return "border-radius:" + formatPx(c.TopLeft) + ";"
	}
	var b strings.Builder
	b.WriteString("border-top-left-radius:")
	b.WriteString(formatPx(c.TopLeft))
	b.WriteString(";border-top-right-radius:")
	b.WriteString(formatPx(c.TopRight))
	b.WriteString(";border-bottom-right-radius:")
	b.WriteString(formatPx(c.BottomRight))
	b.WriteString(";border-bottom-left-radius:")
	b.WriteString(formatPx(c.BottomLeft))
	b.WriteString(";")
	return b.String()
}

func formatPx(v float64) string {
	if v < 0 {
		v = 0
	}
	return strconv.FormatFloat(v, 'f', -1, 64) + "px"
}

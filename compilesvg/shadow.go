package compilesvg

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"psrt/compileasset"
	"psrt/psrt"
)

var shadowRGBA = regexp.MustCompile(`rgba?\([^)]+\)`)

// buildTextShadowFilters returns SVG filter definitions and maps text index to filter id.
func buildTextShadowFilters(pageSlug string, texts []psrt.Text) (string, map[int]string) {
	var b strings.Builder
	ids := make(map[int]string)
	for i := range texts {
		box := compileasset.ParseTextBox(texts[i].Style)
		if box.BoxShadowCSS == "" {
			continue
		}
		dx, dy, blur, color := parseBoxShadow(box.BoxShadowCSS)
		if color == "" {
			continue
		}
		id := fmt.Sprintf("%s-shadow", TextID(pageSlug, texts[i].Index))
		ids[texts[i].Index] = id
		std := blur / 2
		if std < 0.5 {
			std = 0.5
		}
		fmt.Fprintf(&b, `<filter id="%s" x="-50%%" y="-50%%" width="200%%" height="200%%">`, id)
		fmt.Fprintf(&b, `<feDropShadow dx="%s" dy="%s" stdDeviation="%s" flood-color="%s"/></filter>`,
			compileasset.FormatFloatCSS(dx), compileasset.FormatFloatCSS(dy),
			compileasset.FormatFloatCSS(std), svgAttrEscape(color))
	}
	return b.String(), ids
}

func parseBoxShadow(s string) (dx, dy, blur float64, color string) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, 0, 0, ""
	}
	if colorMatch := shadowRGBA.FindString(s); colorMatch != "" {
		color = colorMatch
		s = strings.TrimSpace(strings.Replace(s, colorMatch, "", 1))
	}
	parts := strings.Fields(s)
	if len(parts) >= 2 {
		dx = parseCSSPx(parts[0])
		dy = parseCSSPx(parts[1])
	}
	if len(parts) >= 3 {
		blur = parseCSSPx(parts[2])
	}
	if color == "" && len(parts) >= 4 {
		color = parts[len(parts)-1]
	}
	return dx, dy, blur, color
}

func parseCSSPx(s string) float64 {
	s = strings.TrimSuffix(strings.TrimSpace(s), "px")
	f, _ := strconv.ParseFloat(s, 64)
	return f
}

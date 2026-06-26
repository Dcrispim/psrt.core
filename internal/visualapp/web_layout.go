package visualapp

import (
	"strconv"
	"strings"

	"github.com/Dcrispim/psrt.core/compilesvg"
	"github.com/Dcrispim/psrt.core/psrt"
	"github.com/Dcrispim/psrt.core/styleadapter"
)

func enrichWebPreviewHeight(
	out *styleadapter.WebPreviewStyle,
	t psrt.Text,
	content string,
	canvasW, canvasH int,
) {
	if out.Container == nil {
		out.Container = make(map[string]string)
	}
	if h := strings.TrimSpace(out.Container["height"]); h != "" && !isZeroLikeHeightCSSValue(h) {
		return
	}
	delete(out.Container, "height")
	if canvasH < 1 {
		return
	}
	_, _, _, geomH := compilesvg.TextBlockGeometry(&t, content, canvasW, canvasH)
	if geomH < 1 {
		return
	}
	heightPct := float64(geomH) / float64(canvasH) * 100.0
	out.Container["height"] = strconv.FormatFloat(heightPct, 'f', -1, 64) + "%"
}

func isZeroLikeHeightCSSValue(s string) bool {
	s = strings.TrimSpace(strings.ToLower(s))
	switch s {
	case "", "0", "0%", "0px", "0em", "0rem", "0pt":
		return true
	}
	if strings.HasSuffix(s, "%") {
		if f, err := strconv.ParseFloat(strings.TrimSuffix(s, "%"), 64); err == nil && f <= 0 {
			return true
		}
	}
	if strings.HasSuffix(s, "px") {
		if f, err := strconv.ParseFloat(strings.TrimSuffix(s, "px"), 64); err == nil && f <= 0 {
			return true
		}
	}
	if f, err := strconv.ParseFloat(s, 64); err == nil && f <= 0 {
		return true
	}
	return false
}

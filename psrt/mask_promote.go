package psrt

import (
	"encoding/json"
	"strconv"
	"strings"
)

const defaultMaskHeightPercent = 5.0

// PromoteEmptyTextsToMasks converts empty >> blocks into == mask blocks before formatting.
func PromoteEmptyTextsToMasks(doc *Document) {
	if doc == nil {
		return
	}
	for i := range doc.Pages {
		promotePageEmptyTexts(&doc.Pages[i])
	}
}

func promotePageEmptyTexts(p *Page) {
	if p == nil {
		return
	}
	var kept []Text
	for j := range p.Texts {
		t := p.Texts[j]
		if !IsTextBlockEmpty(t.Content) {
			kept = append(kept, t)
			continue
		}
		height, style := maskHeightFromText(t)
		m := Mask{
			BaseBlock: BaseBlock{
				X:        t.X,
				Y:        t.Y,
				Width:    t.Width,
				Style:    style,
				Index:    t.Index,
				ImageRef: t.ImageRef,
			},
			Height: height,
		}
		p.Masks = append(p.Masks, m)
	}
	p.Texts = kept
}

// maskHeightFromText derives == height from empty >> block: style height %, else 4th coord (textSize).
func maskHeightFromText(t Text) (height float64, cleaned Style) {
	height = defaultMaskHeightPercent
	styleHadHeight := false
	m, err := styleToMap(t.Style)
	if err == nil && len(m) > 0 {
		if h, ok := parseHeightPercentFromStyle(m); ok {
			height = h
			styleHadHeight = true
		}
	}
	if !styleHadHeight && t.TextSize > 0 {
		height = RoundCoord(t.TextSize)
	}
	return height, stripHeightFromStyle(t.Style)
}

func heightFromTextStyle(style Style) (height float64, cleaned Style) {
	height = defaultMaskHeightPercent
	m, err := styleToMap(style)
	if err == nil && len(m) > 0 {
		if h, ok := parseHeightPercentFromStyle(m); ok {
			height = h
		}
	}
	return height, stripHeightFromStyle(style)
}

func stripHeightFromStyle(style Style) Style {
	m, err := styleToMap(style)
	if err != nil || len(m) == 0 {
		return style
	}
	delete(m, "height")
	delete(m, "Height")
	if len(m) == 0 {
		return Style("{}")
	}
	b, err := json.Marshal(m)
	if err != nil {
		return style
	}
	return Style(b)
}

func parseHeightPercentFromStyle(m map[string]any) (float64, bool) {
	for _, key := range []string{"height", "Height"} {
		if v, ok := m[key]; ok {
			if h, ok := parseHeightValue(v); ok {
				return RoundCoord(h), true
			}
		}
	}
	return 0, false
}

func parseHeightValue(v any) (float64, bool) {
	switch x := v.(type) {
	case float64:
		return x, true
	case int:
		return float64(x), true
	case int64:
		return float64(x), true
	case json.Number:
		f, err := x.Float64()
		return f, err == nil
	case string:
		s := strings.TrimSpace(x)
		s = strings.TrimSuffix(s, "%")
		s = strings.TrimSpace(s)
		if s == "" {
			return 0, false
		}
		f, err := strconv.ParseFloat(s, 64)
		return f, err == nil
	default:
		return 0, false
	}
}


package compileasset

import (
	"encoding/json"
	"regexp"
	"strconv"
	"strings"

	"psrt/psrt"
	"psrt/styleadapter"
)

// TextBox describes background/border decoration for a text block.
type TextBox struct {
	Background   string
	BorderCSS    string
	BorderWidth  string
	BorderColor  string
	RadiusCSS    string
	RadiusPx     float64
	BoxShadowCSS string
	PaddingCSS   string
	HasRect      bool
}

var cssPxNum = regexp.MustCompile(`([\d.]+)\s*px`)

// ParseTextBox reads box-related properties from PSRT style JSON (camelCase or kebab-case).
func ParseTextBox(style psrt.Style) TextBox {
	m := normalizeStyleMap(style)
	if len(m) == 0 {
		return TextBox{}
	}
	var box TextBox
	box.Background = textBoxBackground(m)
	box.BorderCSS = textBoxBorderCSS(m)
	box.BorderWidth, box.BorderColor = textBoxBorderSVG(m)
	box.RadiusCSS, box.RadiusPx = textBoxRadius(m)
	box.BoxShadowCSS = rawStringProp(m, "boxShadow", "box-shadow")
	box.PaddingCSS = rawStringProp(m, "padding")
	box.HasRect = box.Background != "" || box.BorderCSS != "" || box.BorderWidth != ""
	return box
}

// CSSBoxFromStyleJSON returns CSS declarations for text box decoration (background, border, etc.).
func CSSBoxFromStyleJSON(style psrt.Style) string {
	return cssBoxFromStyle(style, true, true)
}

// CSSBoxFromStyleJSONNoBackground returns box CSS without background (use SVG rect instead).
func CSSBoxFromStyleJSONNoBackground(style psrt.Style) string {
	return cssBoxFromStyle(style, false, true)
}

func cssBoxFromStyle(style psrt.Style, includeBackground, includeRadius bool) string {
	box := ParseTextBox(style)
	var b strings.Builder
	if includeBackground && box.Background != "" {
		b.WriteString("background-color:")
		b.WriteString(box.Background)
		b.WriteString(";")
	}
	if box.BorderCSS != "" {
		b.WriteString("border:")
		b.WriteString(box.BorderCSS)
		b.WriteString(";")
	}
	if includeRadius && box.RadiusCSS != "" {
		b.WriteString("border-radius:")
		b.WriteString(box.RadiusCSS)
		b.WriteString(";")
	}
	if box.BoxShadowCSS != "" {
		b.WriteString("box-shadow:")
		b.WriteString(box.BoxShadowCSS)
		b.WriteString(";")
	}
	if box.PaddingCSS != "" {
		b.WriteString("padding:")
		b.WriteString(box.PaddingCSS)
		b.WriteString(";")
	}
	return b.String()
}

// StyleJSONWithoutBox returns typography-only style JSON (no box properties).
func StyleJSONWithoutBox(style psrt.Style) []byte {
	raw := bytesTrim(style)
	if len(raw) == 0 || string(raw) == "{}" {
		return []byte("{}")
	}
	m := normalizeStyleMap(style)
	if len(m) == 0 {
		return raw
	}
	boxKeys := []string{
		"backGround", "background", "backgroundColor",
		"border", "borderWidth", "borderStyle", "borderColor", "borderRadius",
		"borderTopLeftRadius", "borderTopRightRadius", "borderBottomRightRadius", "borderBottomLeftRadius",
		"boxShadow", "padding",
		"background-color", "border-width", "border-style", "border-color",
		"border-radius",
		"border-top-left-radius", "border-top-right-radius",
		"border-bottom-right-radius", "border-bottom-left-radius",
		"box-shadow",
	}
	for _, k := range boxKeys {
		delete(m, k)
	}
	if len(m) == 0 {
		return []byte("{}")
	}
	out, err := json.Marshal(m)
	if err != nil {
		return raw
	}
	return out
}

func bytesTrim(style psrt.Style) []byte {
	return []byte(strings.TrimSpace(string(style)))
}

func normalizeStyleMap(style psrt.Style) map[string]json.RawMessage {
	return styleadapter.Normalize(style)
}

func textBoxBackground(m map[string]json.RawMessage) string {
	for _, key := range []string{"backGround", "background", "backgroundColor"} {
		if c := cssColorFromRaw(m[key]); c != "" {
			return c
		}
	}
	return ""
}

func textBoxBorderCSS(m map[string]json.RawMessage) string {
	if b := rawStringProp(m, "border"); b != "" {
		return SanitizeCSSValue(b)
	}
	w := rawStringProp(m, "borderWidth", "border-width")
	s := rawStringProp(m, "borderStyle", "border-style")
	c := rawStringProp(m, "borderColor", "border-color")
	if w == "" && s == "" && c == "" {
		return ""
	}
	if w == "" {
		w = "1px"
	}
	if s == "" {
		s = "solid"
	}
	return SanitizeCSSValue(strings.TrimSpace(w + " " + s + " " + c))
}

func textBoxBorderSVG(m map[string]json.RawMessage) (width, color string) {
	border := rawStringProp(m, "border")
	if border != "" {
		parts := strings.Fields(border)
		if len(parts) >= 1 {
			width = parts[0]
		}
		if len(parts) >= 3 {
			color = strings.Join(parts[2:], " ")
		} else if len(parts) == 2 {
			color = parts[1]
		}
	}
	if width == "" {
		width = rawStringProp(m, "borderWidth", "border-width")
	}
	if color == "" {
		color = rawStringProp(m, "borderColor", "border-color")
	}
	if width == "" && color != "" {
		width = "1px"
	}
	return SanitizeCSSValue(width), SanitizeCSSValue(color)
}

func textBoxRadius(m map[string]json.RawMessage) (css string, px float64) {
	css = rawStringProp(m, "borderRadius", "border-radius")
	if css == "" {
		return "", 0
	}
	css = SanitizeCSSValue(css)
	if m := cssPxNum.FindStringSubmatch(css); len(m) >= 2 {
		if v, err := strconv.ParseFloat(m[1], 64); err == nil {
			px = v
		}
	}
	return css, px
}

func rawStringProp(m map[string]json.RawMessage, keys ...string) string {
	for _, key := range keys {
		if r, ok := m[key]; ok {
			if s := stringifyJSONCSSValue(r); s != "" {
				return SanitizeCSSValue(s)
			}
		}
	}
	return ""
}

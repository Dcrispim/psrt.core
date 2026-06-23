package psrt

import (
	"encoding/json"
	"errors"
	"strings"
)

var errStyleNotJSON = errors.New("style is not valid JSON")

// emptyTextStyleKeys are removed from text blocks with no visible content.
var emptyTextStyleKeys = []string{
	"color",
	"text-align",
	"textAlign",
	"ta",
	"padding",
}

// IsTextBlockEmpty reports whether a text block has no non-whitespace content.
func IsTextBlockEmpty(content string) bool {
	return strings.TrimSpace(content) == ""
}

// CleanEmptyTextBlockStyles removes useless typography/layout styles from empty text blocks.
func CleanEmptyTextBlockStyles(doc *Document) {
	if doc == nil {
		return
	}
	for i := range doc.Pages {
		texts := doc.Pages[i].Texts
		for j := range texts {
			if !IsTextBlockEmpty(texts[j].Content) {
				continue
			}
			texts[j].Style = cleanEmptyTextStyle(texts[j].Style)
		}
		doc.Pages[i].Texts = texts
	}
}

func cleanEmptyTextStyle(style Style) Style {
	m, err := styleToMap(style)
	if err != nil || len(m) == 0 {
		return style
	}
	for _, key := range emptyTextStyleKeys {
		delete(m, key)
	}
	if len(m) == 0 {
		return Style("{}")
	}
	b, err := json.Marshal(m)
	if err != nil {
		return style
	}
	return Style(b)
}

func styleToMap(style Style) (map[string]any, error) {
	raw := strings.TrimSpace(string(style))
	if raw == "" || raw == "{}" {
		return make(map[string]any), nil
	}
	if !json.Valid([]byte(raw)) {
		return nil, errStyleNotJSON
	}
	var m map[string]any
	if err := json.Unmarshal([]byte(raw), &m); err != nil {
		return nil, err
	}
	if m == nil {
		m = make(map[string]any)
	}
	return m, nil
}

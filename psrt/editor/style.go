package editor

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Dcrispim/psrt.core/psrt"
)

// SetStyleKey sets a single property on a style JSON object.
// value is raw JSON (e.g. "#fff", "\"600\"", "true").
func SetStyleKey(style psrt.Style, key, value string) (psrt.Style, error) {
	key = strings.TrimSpace(key)
	if key == "" {
		return style, fmt.Errorf("style key is empty")
	}
	m, err := styleToMap(style)
	if err != nil {
		return style, err
	}
	parsed, err := parseStyleValue(value)
	if err != nil {
		return style, err
	}
	m[key] = parsed
	return mapToStyle(m)
}

// RemoveStyleKey removes a property from a style JSON object.
func RemoveStyleKey(style psrt.Style, key string) (psrt.Style, error) {
	key = strings.TrimSpace(key)
	if key == "" {
		return style, fmt.Errorf("style key is empty")
	}
	m, err := styleToMap(style)
	if err != nil {
		return style, err
	}
	delete(m, key)
	return mapToStyle(m)
}

// MergeStyle shallow-merges partial JSON into the existing style.
func MergeStyle(style psrt.Style, partial json.RawMessage) (psrt.Style, error) {
	m, err := styleToMap(style)
	if err != nil {
		return style, err
	}
	var patch map[string]any
	if err := json.Unmarshal(partial, &patch); err != nil {
		return nil, fmt.Errorf("partial style is not valid JSON: %w", err)
	}
	for k, v := range patch {
		m[k] = v
	}
	return mapToStyle(m)
}

func styleToMap(style psrt.Style) (map[string]any, error) {
	raw := strings.TrimSpace(string(style))
	if raw == "" || raw == "{}" {
		return make(map[string]any), nil
	}
	if !json.Valid([]byte(raw)) {
		return nil, fmt.Errorf("style is not valid JSON")
	}
	var m map[string]any
	if err := json.Unmarshal([]byte(raw), &m); err != nil {
		return nil, fmt.Errorf("style is not valid JSON: %w", err)
	}
	if m == nil {
		m = make(map[string]any)
	}
	return m, nil
}

// parseStyleValue accepts JSON literals or bare strings (e.g. #fff, 600, center).
// Shells such as PowerShell often strip quotes; bare hex colors are coerced to JSON strings.
func parseStyleValue(value string) (any, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", nil
	}
	if parsed, ok := tryJSONValue(value); ok {
		return parsed, nil
	}
	if unquoted := unwrapOneQuotePair(value); unquoted != value {
		if parsed, ok := tryJSONValue(unquoted); ok {
			return parsed, nil
		}
		value = unquoted
	}
	var parsed any
	if err := json.Unmarshal([]byte(jsonStringLiteral(value)), &parsed); err != nil {
		return nil, fmt.Errorf("style value is not valid JSON: %q", value)
	}
	return parsed, nil
}

func tryJSONValue(value string) (any, bool) {
	var parsed any
	if err := json.Unmarshal([]byte(value), &parsed); err != nil {
		return nil, false
	}
	return parsed, true
}

func unwrapOneQuotePair(s string) string {
	if len(s) < 2 {
		return s
	}
	switch {
	case s[0] == '"' && s[len(s)-1] == '"':
		return s[1 : len(s)-1]
	case s[0] == '\'' && s[len(s)-1] == '\'':
		return s[1 : len(s)-1]
	default:
		return s
	}
}

func jsonStringLiteral(s string) string {
	b, _ := json.Marshal(s)
	return string(b)
}

func mapToStyle(m map[string]any) (psrt.Style, error) {
	b, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return psrt.Style(b), nil
}

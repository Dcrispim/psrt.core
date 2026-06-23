package psrt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

// ExpandConstsInStyle applies @name@ placeholders to compacted style JSON per PSRT.md.
func ExpandConstsInStyle(style Style, consts map[string]string) (Style, error) {
	raw := strings.TrimSpace(string(style))
	if raw == "" || raw == "{}" {
		return style, nil
	}
	compact := compactJSONString(raw)
	expanded := ExpandConsts(compact, consts)
	if !json.Valid([]byte(expanded)) {
		return style, fmt.Errorf("style JSON invalid after constant expansion: %s", expanded)
	}
	return Style(expanded), nil
}

func compactJSONString(s string) string {
	var buf bytes.Buffer
	if err := json.Compact(&buf, []byte(s)); err != nil {
		return strings.ReplaceAll(strings.ReplaceAll(s, "\n", ""), "\r", "")
	}
	return buf.String()
}

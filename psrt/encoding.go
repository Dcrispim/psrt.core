package psrt

import (
	"strings"
	"unicode/utf8"
)

// FixMisencodedUTF8 repairs text that was UTF-8 but interpreted as Latin-1 (e.g. "vocÃª" → "você").
func FixMisencodedUTF8(s string) string {
	if s == "" || (!strings.Contains(s, "Ã") && !strings.Contains(s, "Â")) {
		return s
	}
	buf := make([]byte, 0, len(s))
	for _, r := range s {
		if r > 255 {
			return s
		}
		buf = append(buf, byte(r))
	}
	if utf8.Valid(buf) {
		fixed := string(buf)
		if fixed != s {
			return fixed
		}
	}
	return s
}

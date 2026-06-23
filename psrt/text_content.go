package psrt

import "strings"

// NormalizeTextContent trims spaces and tabs before the first character and after
// the last character of a text block body (per PSRT spec).
func NormalizeTextContent(content string) string {
	content = strings.TrimLeft(content, " \t")
	content = strings.TrimRight(content, " \t\r\n")
	return FixMisencodedUTF8(content)
}

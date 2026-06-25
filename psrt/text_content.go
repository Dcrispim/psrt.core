package psrt

import "strings"

// NormalizeTextContent trims spaces and tabs before the first character and after
// the last character of a text block body (per PSRT spec).
func NormalizeTextContent(content string) string {
	content = strings.TrimLeft(content, " \t")
	content = strings.TrimRight(content, " \t\r\n")
	return FixMisencodedUTF8(content)
}

// NormalizePathData collapses line breaks and runs of whitespace in a ~~ block
// body into single spaces — line breaks in the source file are a readability
// convention only and carry no syntactic meaning (RF-4).
func NormalizePathData(content string) string {
	return strings.Join(strings.Fields(content), " ")
}

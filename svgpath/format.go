package svgpath

import "strings"

// SplitCommands breaks SVG path `d` data into one string per command —
// useful for writers that want one command per line for readable diffs
// (e.g. PSRT's ~~ block serializer). It does not validate d; callers that
// need validity should call Parse/Validate first.
func SplitCommands(d string) []string {
	var lines []string
	start := -1
	for i := 0; i < len(d); i++ {
		if isCommandLetter(d[i]) {
			if start >= 0 {
				if seg := strings.TrimSpace(d[start:i]); seg != "" {
					lines = append(lines, seg)
				}
			}
			start = i
		}
	}
	if start >= 0 {
		if seg := strings.TrimSpace(d[start:]); seg != "" {
			lines = append(lines, seg)
		}
	}
	return lines
}

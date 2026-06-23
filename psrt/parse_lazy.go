package psrt

import (
	"fmt"
	"io"
	"strings"
)

// ParseFast parses a PSRT document without loading embedded asset payloads from $SOURCE.
// Source keys are recorded with empty values; use LoadSource to fetch a payload on demand.
func ParseFast(r io.Reader) (Document, error) {
	return parseDocument(r, parseOptions{skipSourceValues: true})
}

// ParseFastString is ParseFast from a string.
func ParseFastString(input string) (Document, error) {
	return ParseFast(strings.NewReader(input))
}

// LoadSource reads one embedded asset from the $SOURCE block of raw PSRT text.
func LoadSource(raw string, url string) (string, error) {
	url = strings.TrimSpace(url)
	if url == "" {
		return "", fmt.Errorf("empty source url")
	}
	inSource := false
	for _, line := range strings.Split(raw, "\n") {
		if inSource {
			if isEndSourceLine(line) {
				break
			}
			s := strings.TrimSpace(line)
			if s == "" {
				continue
			}
			idx := strings.Index(s, pipeSep)
			if idx < 0 {
				continue
			}
			key := strings.TrimSpace(s[:idx])
			if key == url {
				return strings.TrimSpace(s[idx+len(pipeSep):]), nil
			}
			continue
		}
		if d, ok := directive(line); ok && d.kind == dirSource {
			inSource = true
		}
	}
	return "", fmt.Errorf("source not found for %q", url)
}

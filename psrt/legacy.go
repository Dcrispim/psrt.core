package psrt

import (
	"fmt"
	"strconv"
	"strings"
)

const legacyCoordSep = "-"

// ConvertLegacyCoords rewrites one legacy hyphen-separated coordinate quad
// (X-Y-Width-Height or X-Y-Width-TextSize) into the comma-separated form used
// by the current grammar. The legacy format never supported negative
// numbers — the hyphen doubled as field separator and minus sign, so a
// negative value was already inexpressible — meaning this is a plain
// re-join: split on "-", require exactly four numeric chunks, join on ",".
func ConvertLegacyCoords(s string) (string, error) {
	chunks := strings.Split(strings.TrimSpace(s), legacyCoordSep)
	if len(chunks) != 4 {
		return "", fmt.Errorf("legacy coords want X-Y-Width-Height, got %q", s)
	}
	for i, c := range chunks {
		c = strings.TrimSpace(c)
		if _, err := strconv.ParseFloat(c, 64); err != nil {
			return "", fmt.Errorf("invalid legacy coord segment %q: %w", c, err)
		}
		chunks[i] = c
	}
	return strings.Join(chunks, coordSep), nil
}

// ConvertLegacyDocument rewrites every >>, ==, and ~~ header's coordinate
// quad in raw legacy PSRT text (hyphen-separated) to the comma-separated
// grammar parsed by Parse. Directives, text/path bodies, styles, and
// everything else in the document are left untouched. The result is meant
// to be fed straight into Parse/ParseString.
func ConvertLegacyDocument(raw string) (string, error) {
	lines := strings.Split(raw, "\n")
	for i, line := range lines {
		trimmed := strings.TrimLeft(line, " \t")
		lead := line[:len(line)-len(trimmed)]

		var marker string
		switch {
		case strings.HasPrefix(trimmed, ">>"):
			marker = ">>"
		case strings.HasPrefix(trimmed, "=="):
			marker = "=="
		case strings.HasPrefix(trimmed, "~~"):
			marker = "~~"
		default:
			continue
		}

		body := strings.TrimSpace(trimmed[len(marker):])
		parts := strings.SplitN(body, pipeSep, 2)
		coords, err := ConvertLegacyCoords(parts[0])
		if err != nil {
			return "", fmt.Errorf("line %d: %w", i+1, err)
		}

		rest := ""
		if len(parts) == 2 {
			rest = pipeSep + parts[1]
		}
		lines[i] = lead + marker + coords + rest
	}
	return strings.Join(lines, "\n"), nil
}

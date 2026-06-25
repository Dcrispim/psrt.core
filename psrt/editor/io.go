package editor

import (
	"fmt"
	"os"

	"github.com/Dcrispim/psrt.core/psrt"
)

func readFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

func writeFile(path string, data []byte) error {
	return os.WriteFile(path, data, 0o644)
}

// FormatDocument serialises a document to PSRT bytes.
func FormatDocument(doc *psrt.Document) ([]byte, error) {
	psrt.PromoteEmptyTextsToMasks(doc)
	psrt.CleanEmptyTextBlockStyles(doc)
	data, err := psrt.FormatPSRT(*doc, false)
	if err != nil {
		return nil, fmt.Errorf("format: %w", err)
	}
	return data, nil
}

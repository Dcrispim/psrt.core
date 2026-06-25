package editor

import (
	"fmt"
	"strings"

	"github.com/Dcrispim/psrt.core/psrt"
)

// AddFont appends a font URL if not already present.
func AddFont(doc *psrt.Document, url string) error {
	url = strings.TrimSpace(url)
	if url == "" {
		return fmt.Errorf("font URL is empty")
	}
	for _, f := range doc.Fonts {
		if f == url {
			return nil
		}
	}
	doc.Fonts = append(doc.Fonts, url)
	return nil
}

// RemoveFont removes a font URL.
func RemoveFont(doc *psrt.Document, url string) error {
	url = strings.TrimSpace(url)
	for i, f := range doc.Fonts {
		if f == url {
			doc.Fonts = append(doc.Fonts[:i], doc.Fonts[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("font %q not found", url)
}

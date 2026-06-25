package editor

import (
	"fmt"

	"github.com/Dcrispim/psrt.core/psrt"
)

// FindPage returns the page with the given name.
func FindPage(doc *psrt.Document, name string) (*psrt.Page, error) {
	for i := range doc.Pages {
		if doc.Pages[i].Name == name {
			return &doc.Pages[i], nil
		}
	}
	return nil, fmt.Errorf("page %q not found", name)
}

// FindPageIndex returns the slice index of a page by name.
func FindPageIndex(doc *psrt.Document, name string) (int, error) {
	for i := range doc.Pages {
		if doc.Pages[i].Name == name {
			return i, nil
		}
	}
	return -1, fmt.Errorf("page %q not found", name)
}

// FindTextByIndex returns the text block with the given Index field on a page.
func FindTextByIndex(page *psrt.Page, index int) (*psrt.Text, int, error) {
	for i := range page.Texts {
		if page.Texts[i].Index == index {
			return &page.Texts[i], i, nil
		}
	}
	return nil, -1, fmt.Errorf("text index %d not found on page %q", index, page.Name)
}

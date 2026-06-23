package visualapp

import (
	"fmt"

	"psrt/psrt"
	"psrt/psrt/editor"
)

// GetDocumentPSRT returns the current document serialised as PSRT text.
func (a *App) GetDocumentPSRT() (string, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	data, err := editor.FormatDocument(&a.doc)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// SetDocumentFromPSRT parses PSRT text and replaces the in-memory document.
func (a *App) SetDocumentFromPSRT(text string) error {
	doc, err := psrt.ParseString(text)
	if err != nil {
		return fmt.Errorf("parse: %w", err)
	}
	a.mu.Lock()
	if !a.inEdit {
		a.snapshot()
	}
	prevPage := a.activePage
	a.doc = doc
	if _, err := editor.FindPage(&a.doc, prevPage); err != nil {
		if len(doc.Pages) > 0 {
			a.activePage = doc.Pages[0].Name
		} else {
			a.activePage = ""
		}
	}
	a.selectedIdx = -1
	inEdit := a.inEdit
	a.mu.Unlock()
	if !inEdit {
		a.notify(nil)
		a.maybeAutoCompile()
	}
	return nil
}

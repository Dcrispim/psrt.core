package editor

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Dcrispim/psrt.core/psrt"
)

// RenamePage changes a page name.
func RenamePage(doc *psrt.Document, oldName, newName string) error {
	if oldName == newName {
		return nil
	}
	if _, err := FindPage(doc, newName); err == nil {
		return fmt.Errorf("page %q already exists", newName)
	}
	p, err := FindPage(doc, oldName)
	if err != nil {
		return err
	}
	p.Name = newName
	return nil
}

// SetPagePath updates the page image URL.
func SetPagePath(doc *psrt.Document, pageName, path string) error {
	p, err := FindPage(doc, pageName)
	if err != nil {
		return err
	}
	p.ImageURL = path
	return nil
}

// SetPageStyle merges style properties on a page.
func SetPageStyle(doc *psrt.Document, pageName string, key, value string, partial json.RawMessage) error {
	p, err := FindPage(doc, pageName)
	if err != nil {
		return err
	}
	updated, err := applyStyleUpdate(p.Style, key, value, partial)
	if err != nil {
		return err
	}
	p.Style = updated
	return nil
}

// RemovePageStyleKey removes a style property from a page.
func RemovePageStyleKey(doc *psrt.Document, pageName, key string) error {
	p, err := FindPage(doc, pageName)
	if err != nil {
		return err
	}
	updated, err := RemoveStyleKey(p.Style, key)
	if err != nil {
		return err
	}
	p.Style = updated
	return nil
}

// MovePage reorders a page before or after another page.
func MovePage(doc *psrt.Document, pageName, before, after string) error {
	if before != "" && after != "" {
		return fmt.Errorf("use only one of --before or --after")
	}
	if before == "" && after == "" {
		return fmt.Errorf("one of --before or --after is required")
	}
	from, err := FindPageIndex(doc, pageName)
	if err != nil {
		return err
	}
	refName := before
	if refName == "" {
		refName = after
	}
	ref, err := FindPageIndex(doc, refName)
	if err != nil {
		return err
	}
	if from == ref {
		return fmt.Errorf("cannot move page relative to itself")
	}
	pages := doc.Pages
	if before != "" {
		pages, err = moveBeforeIndex(pages, from, ref)
	} else {
		pages, err = moveAfterIndex(pages, from, ref)
	}
	if err != nil {
		return err
	}
	doc.Pages = pages
	return nil
}

// AddPage inserts a new page (append, or before/after another page by name).
func AddPage(doc *psrt.Document, page psrt.Page, before, after string) error {
	if before != "" && after != "" {
		return fmt.Errorf("use only one of before or after")
	}
	if _, err := FindPage(doc, page.Name); err == nil {
		return fmt.Errorf("page %q already exists", page.Name)
	}
	if strings.TrimSpace(string(page.Style)) == "" {
		page.Style = psrt.Style("{}")
	}
	if before == "" && after == "" {
		doc.Pages = append(doc.Pages, page)
		return nil
	}
	refName := before
	if refName == "" {
		refName = after
	}
	ref, err := FindPageIndex(doc, refName)
	if err != nil {
		return err
	}
	if before != "" {
		doc.Pages = insertPageAt(doc.Pages, ref, page)
	} else {
		doc.Pages = insertPageAt(doc.Pages, ref+1, page)
	}
	return nil
}

// RemovePage deletes a page by name.
func RemovePage(doc *psrt.Document, name string) error {
	idx, err := FindPageIndex(doc, name)
	if err != nil {
		return err
	}
	doc.Pages = append(doc.Pages[:idx], doc.Pages[idx+1:]...)
	return nil
}

func insertPageAt(pages []psrt.Page, pos int, page psrt.Page) []psrt.Page {
	out := make([]psrt.Page, 0, len(pages)+1)
	out = append(out, pages[:pos]...)
	out = append(out, page)
	out = append(out, pages[pos:]...)
	return out
}

func applyStyleUpdate(style psrt.Style, key, value string, partial json.RawMessage) (psrt.Style, error) {
	if len(partial) > 0 {
		if key != "" || value != "" {
			return style, fmt.Errorf("use either --key/--value or --style, not both")
		}
		return MergeStyle(style, partial)
	}
	if key == "" {
		return style, fmt.Errorf("either --key and --value or --style is required")
	}
	return SetStyleKey(style, key, value)
}

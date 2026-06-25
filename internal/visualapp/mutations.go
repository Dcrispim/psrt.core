package visualapp

import (
	"encoding/json"
	"math"

	"github.com/Dcrispim/psrt.core/psrt"
	"github.com/Dcrispim/psrt.core/psrt/editor"
)

// PatchText applies UI patch to a text block.
func (a *App) PatchText(pageName string, index int, patch TextPatch) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	if !a.inEdit {
		a.snapshot()
	}
	if patch.Content != nil {
		if err := editor.SetTextContent(&a.doc, pageName, index, *patch.Content, patch.Append); err != nil {
			return err
		}
	}
	pos := editor.PositionFields{}
	if patch.X != nil {
		pos.X = patch.X
	}
	if patch.Y != nil {
		pos.Y = patch.Y
	}
	if patch.Width != nil {
		pos.Width = patch.Width
	}
	if patch.TextSize != nil {
		pos.TextSize = patch.TextSize
	}
	if !posEmpty(pos) {
		snap := a.snapGrid
		if snap > 0 {
			if pos.X != nil {
				v := snapRound(*pos.X, snap)
				pos.X = &v
			}
			if pos.Y != nil {
				v := snapRound(*pos.Y, snap)
				pos.Y = &v
			}
			if pos.Width != nil {
				v := snapRound(*pos.Width, snap)
				pos.Width = &v
			}
			if pos.TextSize != nil {
				v := snapRound(*pos.TextSize, snap)
				pos.TextSize = &v
			}
		}
		if err := editor.SetTextPosition(&a.doc, pageName, index, pos); err != nil {
			return err
		}
	}
	if patch.ImageRef != nil {
		t, err := findText(&a.doc, pageName, index)
		if err != nil {
			return err
		}
		t.ImageRef = *patch.ImageRef
	}
	for _, k := range patch.StyleRemove {
		if err := editor.RemoveTextStyleKey(&a.doc, pageName, index, k); err != nil {
			return err
		}
	}
	for k, v := range patch.StyleSet {
		if err := editor.SetTextStyle(&a.doc, pageName, index, k, v, nil); err != nil {
			return err
		}
	}
	if !a.inEdit {
		a.notify(nil)
		a.maybeAutoCompile()
	}
	return nil
}

// PatchPage updates page properties.
func (a *App) PatchPage(patch PagePatch) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.snapshot()
	page := a.activePage
	if patch.Name != nil && *patch.Name != page {
		if err := editor.RenamePage(&a.doc, page, *patch.Name); err != nil {
			return err
		}
		a.activePage = *patch.Name
		page = *patch.Name
	}
	if patch.ImageURL != nil {
		if err := editor.SetPagePath(&a.doc, page, *patch.ImageURL); err != nil {
			return err
		}
	}
	for _, k := range patch.StyleRemove {
		if err := editor.RemovePageStyleKey(&a.doc, page, k); err != nil {
			return err
		}
	}
	for k, v := range patch.StyleSet {
		if err := editor.SetPageStyle(&a.doc, page, k, v, nil); err != nil {
			return err
		}
	}
	a.notify(nil)
	a.maybeAutoCompile()
	return nil
}

// AddPage adds a page.
func (a *App) AddPage(name, imageURL, styleJSON string) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.snapshot()
	style := psrt.Style("{}")
	if styleJSON != "" {
		style = psrt.Style(styleJSON)
	}
	p := psrt.Page{Name: name, ImageURL: imageURL, Style: style}
	if err := editor.AddPage(&a.doc, p, "", ""); err != nil {
		return err
	}
	a.activePage = name
	a.selectedIdx = -1
	a.notify(nil)
	return nil
}

// RemovePage removes active or named page.
func (a *App) RemovePage(name string) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	if name == "" {
		name = a.activePage
	}
	a.snapshot()
	if err := editor.RemovePage(&a.doc, name); err != nil {
		return err
	}
	if len(a.doc.Pages) > 0 {
		a.activePage = a.doc.Pages[0].Name
	} else {
		a.activePage = ""
	}
	a.selectedIdx = -1
	a.notify(nil)
	return nil
}

// MovePage reorders relative to ref page (before=true places before ref).
func (a *App) MovePage(name, ref string, before bool) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.snapshot()
	if before {
		if err := editor.MovePage(&a.doc, name, ref, ""); err != nil {
			return err
		}
	} else {
		if err := editor.MovePage(&a.doc, name, "", ref); err != nil {
			return err
		}
	}
	a.notify(nil)
	return nil
}

// AddTextBlock adds a new text on active page.
func (a *App) AddTextBlock(index int, x, y, width, textSize float64, content, styleJSON, imageRef string) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.snapshot()
	t := psrt.Text{
		BaseBlock: psrt.BaseBlock{
			Index: index, X: x, Y: y, Width: width,
			ImageRef: imageRef, Style: psrt.Style("{}"),
		},
		TextSize: textSize,
		Content:  content,
	}
	if styleJSON != "" {
		t.Style = psrt.Style(styleJSON)
	}
	if err := editor.AddText(&a.doc, a.activePage, t, -1, -1); err != nil {
		return err
	}
	a.notify(nil)
	return nil
}

// RemoveText removes text by index.
func (a *App) RemoveText(index int) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.snapshot()
	if err := editor.RemoveText(&a.doc, a.activePage, index); err != nil {
		return err
	}
	if a.selectedIdx == index {
		a.selectedIdx = -1
	}
	a.notify(nil)
	return nil
}

// ReorderText moves text before/after another index.
func (a *App) ReorderText(index, ref int, before bool) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.snapshot()
	b, af := -1, -1
	if before {
		b = ref
	} else {
		af = ref
	}
	if err := editor.ReorderTextRelative(&a.doc, a.activePage, index, b, af); err != nil {
		return err
	}
	a.notify(nil)
	return nil
}

// SetFonts replaces font list.
func (a *App) SetFonts(urls []string) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.snapshot()
	a.doc.Fonts = urls
	a.notify(nil)
	return nil
}

// AddFont appends a font URL.
func (a *App) AddFont(url string) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.snapshot()
	if err := editor.AddFont(&a.doc, url); err != nil {
		return err
	}
	a.notify(nil)
	return nil
}

// RemoveFont removes font URL.
func (a *App) RemoveFont(url string) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.snapshot()
	if err := editor.RemoveFont(&a.doc, url); err != nil {
		return err
	}
	a.notify(nil)
	return nil
}

// AddConst adds constant.
func (a *App) AddConst(name, value string) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.snapshot()
	if err := editor.AddConst(&a.doc, name, value); err != nil {
		return err
	}
	a.notify(nil)
	return nil
}

// RemoveConst removes constant.
func (a *App) RemoveConst(name string) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.snapshot()
	if err := editor.RemoveConst(&a.doc, name); err != nil {
		return err
	}
	a.notify(nil)
	return nil
}

func snapRound(v, step float64) float64 {
	if step <= 0 {
		return psrt.RoundCoord(v)
	}
	return psrt.RoundCoord(math.Round(v/step) * step)
}

func posEmpty(p editor.PositionFields) bool {
	return p.X == nil && p.Y == nil && p.Width == nil && p.TextSize == nil
}

func findText(doc *psrt.Document, pageName string, index int) (*psrt.Text, error) {
	p, err := editor.FindPage(doc, pageName)
	if err != nil {
		return nil, err
	}
	t, _, err := editor.FindTextByIndex(p, index)
	return t, err
}

// ParseStyleMap parses JSON style object to key-value strings for editor.
func ParseStyleMap(raw string) (map[string]string, error) {
	if raw == "" || raw == "{}" {
		return nil, nil
	}
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(raw), &m); err != nil {
		return nil, err
	}
	out := make(map[string]string, len(m))
	for k, v := range m {
		b, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		out[k] = string(b)
	}
	return out, nil
}

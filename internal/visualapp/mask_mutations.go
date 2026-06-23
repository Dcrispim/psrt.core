package visualapp

import (
	"psrt/psrt"
	"psrt/psrt/editor"
)

// PatchMask applies UI patch to a mask block.
func (a *App) PatchMask(pageName string, index int, patch MaskPatch) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	if !a.inEdit {
		a.snapshot()
	}
	pos := editor.MaskPositionFields{}
	if patch.X != nil {
		pos.X = patch.X
	}
	if patch.Y != nil {
		pos.Y = patch.Y
	}
	if patch.Width != nil {
		pos.Width = patch.Width
	}
	if patch.Height != nil {
		pos.Height = patch.Height
	}
	if !pos.IsEmpty() {
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
			if pos.Height != nil {
				v := snapRound(*pos.Height, snap)
				pos.Height = &v
			}
		}
		if err := editor.SetMaskPosition(&a.doc, pageName, index, pos); err != nil {
			return err
		}
	}
	if patch.ImageRef != nil {
		m, err := findMask(&a.doc, pageName, index)
		if err != nil {
			return err
		}
		m.ImageRef = *patch.ImageRef
	}
	for _, k := range patch.StyleRemove {
		if err := editor.RemoveMaskStyleKey(&a.doc, pageName, index, k); err != nil {
			return err
		}
	}
	for k, v := range patch.StyleSet {
		if err := editor.SetMaskStyle(&a.doc, pageName, index, k, v, nil); err != nil {
			return err
		}
	}
	if !a.inEdit {
		a.notify(nil)
		a.maybeAutoCompile()
	}
	return nil
}

func findMask(doc *psrt.Document, pageName string, index int) (*psrt.Mask, error) {
	p, err := editor.FindPage(doc, pageName)
	if err != nil {
		return nil, err
	}
	m, _, err := psrt.FindMaskByIndex(p, index)
	return m, err
}

// AddMaskBlock adds a new mask on active page.
func (a *App) AddMaskBlock(index int, x, y, width, height float64, styleJSON, imageRef string) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.snapshot()
	m := psrt.Mask{
		BaseBlock: psrt.BaseBlock{
			Index: index, X: x, Y: y, Width: width,
			ImageRef: imageRef, Style: psrt.Style("{}"),
		},
		Height: height,
	}
	if styleJSON != "" {
		m.Style = psrt.Style(styleJSON)
	}
	if err := editor.AddMask(&a.doc, a.activePage, m, -1, -1); err != nil {
		return err
	}
	a.notify(nil)
	return nil
}

// RemoveMask removes mask by index.
func (a *App) RemoveMask(index int) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.snapshot()
	if err := editor.RemoveMask(&a.doc, a.activePage, index); err != nil {
		return err
	}
	if a.selectedIdx == index {
		a.selectedIdx = -1
	}
	a.notify(nil)
	return nil
}

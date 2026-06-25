package editor

import (
	"encoding/json"
	"fmt"

	"github.com/Dcrispim/psrt.core/psrt"
)

// MaskPositionFields holds optional mask coordinate updates.
type MaskPositionFields struct {
	X, Y, Width, Height *float64
}

func (p MaskPositionFields) IsEmpty() bool {
	return p.X == nil && p.Y == nil && p.Width == nil && p.Height == nil
}

// SetMaskPosition sets mask X, Y, Width and/or Height (percent).
func SetMaskPosition(doc *psrt.Document, pageName string, maskIndex int, pos MaskPositionFields) error {
	if pos.IsEmpty() {
		return fmt.Errorf("at least one of x, y, width, or height is required")
	}
	m, err := findMask(doc, pageName, maskIndex)
	if err != nil {
		return err
	}
	if pos.X != nil {
		m.X = psrt.RoundCoord(*pos.X)
	}
	if pos.Y != nil {
		m.Y = psrt.RoundCoord(*pos.Y)
	}
	if pos.Width != nil {
		m.Width = psrt.RoundCoord(*pos.Width)
	}
	if pos.Height != nil {
		m.Height = psrt.RoundCoord(*pos.Height)
	}
	return nil
}

// AddMask inserts a mask block on a page.
func AddMask(doc *psrt.Document, pageName string, mask psrt.Mask, beforeIndex, afterIndex int) error {
	p, err := FindPage(doc, pageName)
	if err != nil {
		return err
	}
	if p.Masks == nil {
		p.Masks = []psrt.Mask{}
	}
	if beforeIndex >= 0 {
		_, refPos, err := FindMaskByIndex(p, beforeIndex)
		if err != nil {
			return err
		}
		p.Masks = insertMaskAt(p.Masks, mask, refPos)
	} else if afterIndex >= 0 {
		_, refPos, err := FindMaskByIndex(p, afterIndex)
		if err != nil {
			return err
		}
		p.Masks = insertMaskAt(p.Masks, mask, refPos+1)
	} else {
		p.Masks = append(p.Masks, mask)
	}
	return nil
}

func insertMaskAt(masks []psrt.Mask, m psrt.Mask, at int) []psrt.Mask {
	if at < 0 {
		at = 0
	}
	if at > len(masks) {
		at = len(masks)
	}
	out := make([]psrt.Mask, 0, len(masks)+1)
	out = append(out, masks[:at]...)
	out = append(out, m)
	out = append(out, masks[at:]...)
	return out
}

// RemoveMask deletes a mask by index field.
func RemoveMask(doc *psrt.Document, pageName string, maskIndex int) error {
	p, err := FindPage(doc, pageName)
	if err != nil {
		return err
	}
	_, pos, err := FindMaskByIndex(p, maskIndex)
	if err != nil {
		return err
	}
	p.Masks = append(p.Masks[:pos], p.Masks[pos+1:]...)
	return nil
}

func findMask(doc *psrt.Document, pageName string, maskIndex int) (*psrt.Mask, error) {
	p, err := FindPage(doc, pageName)
	if err != nil {
		return nil, err
	}
	m, _, err := psrt.FindMaskByIndex(p, maskIndex)
	return m, err
}

// FindMaskByIndex delegates to psrt.FindMaskByIndex.
func FindMaskByIndex(page *psrt.Page, index int) (*psrt.Mask, int, error) {
	return psrt.FindMaskByIndex(page, index)
}

// SetMaskStyle merges a style property on a mask block.
func SetMaskStyle(doc *psrt.Document, pageName string, maskIndex int, key, value string, partial json.RawMessage) error {
	m, err := findMask(doc, pageName, maskIndex)
	if err != nil {
		return err
	}
	updated, err := applyStyleUpdate(m.Style, key, value, partial)
	if err != nil {
		return err
	}
	m.Style = updated
	return nil
}

// RemoveMaskStyleKey removes a style property from a mask block.
func RemoveMaskStyleKey(doc *psrt.Document, pageName string, maskIndex int, key string) error {
	m, err := findMask(doc, pageName, maskIndex)
	if err != nil {
		return err
	}
	updated, err := RemoveStyleKey(m.Style, key)
	if err != nil {
		return err
	}
	m.Style = updated
	return nil
}

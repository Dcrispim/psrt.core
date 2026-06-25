package editor

import (
	"encoding/json"
	"fmt"

	"psrt/psrt"
	"psrt/svgpath"
)

// PathMaskPositionFields holds optional path mask coordinate updates.
type PathMaskPositionFields struct {
	X, Y, Width, Height *float64
}

func (p PathMaskPositionFields) IsEmpty() bool {
	return p.X == nil && p.Y == nil && p.Width == nil && p.Height == nil
}

// SetPathMaskPosition sets path mask X, Y, Width and/or Height (percent).
func SetPathMaskPosition(doc *psrt.Document, pageName string, maskIndex int, pos PathMaskPositionFields) error {
	if pos.IsEmpty() {
		return fmt.Errorf("at least one of x, y, width, or height is required")
	}
	m, err := findPathMask(doc, pageName, maskIndex)
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

// AddPathMask inserts a path mask block on a page.
func AddPathMask(doc *psrt.Document, pageName string, mask psrt.PathMask, beforeIndex, afterIndex int) error {
	p, err := FindPage(doc, pageName)
	if err != nil {
		return err
	}
	if p.PathMasks == nil {
		p.PathMasks = []psrt.PathMask{}
	}
	if beforeIndex >= 0 {
		_, refPos, err := FindPathMaskByIndex(p, beforeIndex)
		if err != nil {
			return err
		}
		p.PathMasks = insertPathMaskAt(p.PathMasks, mask, refPos)
	} else if afterIndex >= 0 {
		_, refPos, err := FindPathMaskByIndex(p, afterIndex)
		if err != nil {
			return err
		}
		p.PathMasks = insertPathMaskAt(p.PathMasks, mask, refPos+1)
	} else {
		p.PathMasks = append(p.PathMasks, mask)
	}
	return nil
}

func insertPathMaskAt(masks []psrt.PathMask, m psrt.PathMask, at int) []psrt.PathMask {
	if at < 0 {
		at = 0
	}
	if at > len(masks) {
		at = len(masks)
	}
	out := make([]psrt.PathMask, 0, len(masks)+1)
	out = append(out, masks[:at]...)
	out = append(out, m)
	out = append(out, masks[at:]...)
	return out
}

// RemovePathMask deletes a path mask by index field.
func RemovePathMask(doc *psrt.Document, pageName string, maskIndex int) error {
	p, err := FindPage(doc, pageName)
	if err != nil {
		return err
	}
	_, pos, err := FindPathMaskByIndex(p, maskIndex)
	if err != nil {
		return err
	}
	p.PathMasks = append(p.PathMasks[:pos], p.PathMasks[pos+1:]...)
	return nil
}

func findPathMask(doc *psrt.Document, pageName string, maskIndex int) (*psrt.PathMask, error) {
	p, err := FindPage(doc, pageName)
	if err != nil {
		return nil, err
	}
	m, _, err := psrt.FindPathMaskByIndex(p, maskIndex)
	return m, err
}

// FindPathMaskByIndex delegates to psrt.FindPathMaskByIndex.
func FindPathMaskByIndex(page *psrt.Page, index int) (*psrt.PathMask, int, error) {
	return psrt.FindPathMaskByIndex(page, index)
}

// SetPathMaskStyle merges a style property on a path mask block.
func SetPathMaskStyle(doc *psrt.Document, pageName string, maskIndex int, key, value string, partial json.RawMessage) error {
	m, err := findPathMask(doc, pageName, maskIndex)
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

// RemovePathMaskStyleKey removes a style property from a path mask block.
func RemovePathMaskStyleKey(doc *psrt.Document, pageName string, maskIndex int, key string) error {
	m, err := findPathMask(doc, pageName, maskIndex)
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

// SetPathMaskPath replaces the SVG path `d` data of a path mask block.
// Validation mirrors flushPathMaskBlock's parse-time checks (see
// psrt/parser.go): the path must be non-empty, syntactically valid SVG path
// data, and describe a single shape (RF-7).
func SetPathMaskPath(doc *psrt.Document, pageName string, maskIndex int, path string) error {
	normalized := psrt.NormalizePathData(path)
	if normalized == "" {
		return fmt.Errorf("path mask body is empty")
	}
	info, err := svgpath.Parse(normalized)
	if err != nil {
		return fmt.Errorf("invalid svg path data: %w", err)
	}
	if info.Subpaths > 1 {
		return fmt.Errorf("path mask must be a single shape (multiple M/m commands found)")
	}
	m, err := findPathMask(doc, pageName, maskIndex)
	if err != nil {
		return err
	}
	m.Path = normalized
	return nil
}

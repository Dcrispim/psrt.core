package editor

import (
	"fmt"

	"psrt/psrt"
)

// PositionFields holds optional coordinate updates (nil = leave unchanged).
type PositionFields struct {
	X, Y, Width, TextSize *float64
}

func (p PositionFields) isEmpty() bool {
	return p.X == nil && p.Y == nil && p.Width == nil && p.TextSize == nil
}

// SetTextPosition sets X, Y, Width and/or TextSize to absolute values.
// Only non-nil fields in pos are applied.
func SetTextPosition(doc *psrt.Document, pageName string, textIndex int, pos PositionFields) error {
	if pos.isEmpty() {
		return fmt.Errorf("at least one of x, y, width, or text-size is required")
	}
	t, err := findText(doc, pageName, textIndex)
	if err != nil {
		return err
	}
	applyPositionSet(t, pos)
	return nil
}

// NudgeTextPosition adds deltas to X, Y, Width and/or TextSize.
// Only non-nil fields in delta are applied.
func NudgeTextPosition(doc *psrt.Document, pageName string, textIndex int, delta PositionFields) error {
	if delta.isEmpty() {
		return fmt.Errorf("at least one of x, y, width, or text-size is required")
	}
	t, err := findText(doc, pageName, textIndex)
	if err != nil {
		return err
	}
	applyPositionNudge(t, delta)
	return nil
}

func applyPositionSet(t *psrt.Text, pos PositionFields) {
	if pos.X != nil {
		t.X = psrt.RoundCoord(*pos.X)
	}
	if pos.Y != nil {
		t.Y = psrt.RoundCoord(*pos.Y)
	}
	if pos.Width != nil {
		t.Width = psrt.RoundCoord(*pos.Width)
	}
	if pos.TextSize != nil {
		t.TextSize = psrt.RoundCoord(*pos.TextSize)
	}
}

func applyPositionNudge(t *psrt.Text, delta PositionFields) {
	if delta.X != nil {
		t.X = psrt.RoundCoord(t.X + *delta.X)
	}
	if delta.Y != nil {
		t.Y = psrt.RoundCoord(t.Y + *delta.Y)
	}
	if delta.Width != nil {
		t.Width = psrt.RoundCoord(t.Width + *delta.Width)
	}
	if delta.TextSize != nil {
		t.TextSize = psrt.RoundCoord(t.TextSize + *delta.TextSize)
	}
}

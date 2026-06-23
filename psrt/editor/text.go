package editor

import (
	"encoding/json"
	"fmt"
	"strconv"

	"psrt/psrt"
)

// SetTextStyle merges style properties on a text block.
func SetTextStyle(doc *psrt.Document, pageName string, textIndex int, key, value string, partial json.RawMessage) error {
	t, err := findText(doc, pageName, textIndex)
	if err != nil {
		return err
	}
	updated, err := applyStyleUpdate(t.Style, key, value, partial)
	if err != nil {
		return err
	}
	t.Style = updated
	return nil
}

// RemoveTextStyleKey removes a style property from a text block.
func RemoveTextStyleKey(doc *psrt.Document, pageName string, textIndex int, key string) error {
	t, err := findText(doc, pageName, textIndex)
	if err != nil {
		return err
	}
	updated, err := RemoveStyleKey(t.Style, key)
	if err != nil {
		return err
	}
	t.Style = updated
	return nil
}

// SetTextContent updates text content (replace or append).
func SetTextContent(doc *psrt.Document, pageName string, textIndex int, content string, appendContent bool) error {
	t, err := findText(doc, pageName, textIndex)
	if err != nil {
		return err
	}
	if appendContent {
		t.Content = psrt.NormalizeTextContent(t.Content + content)
	} else {
		t.Content = psrt.NormalizeTextContent(content)
	}
	return nil
}

// AddText inserts a new text block on a page.
func AddText(doc *psrt.Document, pageName string, text psrt.Text, beforeIndex, afterIndex int) error {
	p, err := FindPage(doc, pageName)
	if err != nil {
		return err
	}
	text.Content = psrt.NormalizeTextContent(text.Content)
	if beforeIndex >= 0 && afterIndex >= 0 {
		return fmt.Errorf("use only one of --before or --after")
	}
	if beforeIndex < 0 && afterIndex < 0 {
		p.Texts = append(p.Texts, text)
		return nil
	}
	refIndex := beforeIndex
	if refIndex < 0 {
		refIndex = afterIndex
	}
	_, refPos, err := FindTextByIndex(p, refIndex)
	if err != nil {
		return err
	}
	if beforeIndex >= 0 {
		p.Texts = insertAt(p.Texts, refPos, text)
	} else {
		p.Texts = insertAt(p.Texts, refPos+1, text)
	}
	return nil
}

// RemoveText removes a text block by Index field.
func RemoveText(doc *psrt.Document, pageName string, textIndex int) error {
	p, err := FindPage(doc, pageName)
	if err != nil {
		return err
	}
	_, pos, err := FindTextByIndex(p, textIndex)
	if err != nil {
		return err
	}
	p.Texts = append(p.Texts[:pos], p.Texts[pos+1:]...)
	return nil
}

// ReorderTextRelative moves a text block before or after another text index in the page slice.
func ReorderTextRelative(doc *psrt.Document, pageName string, textIndex int, beforeIndex, afterIndex int) error {
	if beforeIndex >= 0 && afterIndex >= 0 {
		return fmt.Errorf("use only one of --before or --after")
	}
	if beforeIndex < 0 && afterIndex < 0 {
		return fmt.Errorf("one of --before or --after is required")
	}
	p, err := FindPage(doc, pageName)
	if err != nil {
		return err
	}
	_, from, err := FindTextByIndex(p, textIndex)
	if err != nil {
		return err
	}
	refIndex := beforeIndex
	if refIndex < 0 {
		refIndex = afterIndex
	}
	_, ref, err := FindTextByIndex(p, refIndex)
	if err != nil {
		return err
	}
	if from == ref {
		return fmt.Errorf("cannot move text relative to itself")
	}
	texts := p.Texts
	if beforeIndex >= 0 {
		texts, err = moveBeforeIndex(texts, from, ref)
	} else {
		texts, err = moveAfterIndex(texts, from, ref)
	}
	if err != nil {
		return err
	}
	p.Texts = texts
	return nil
}

// ReorderTextTo moves a text block to an absolute slice position.
func ReorderTextTo(doc *psrt.Document, pageName string, textIndex, to int) error {
	p, err := FindPage(doc, pageName)
	if err != nil {
		return err
	}
	_, from, err := FindTextByIndex(p, textIndex)
	if err != nil {
		return err
	}
	texts, err := moveToIndex(p.Texts, from, to)
	if err != nil {
		return err
	}
	p.Texts = texts
	return nil
}

// ReorderTextByDelta moves a text block by a delta in slice order.
func ReorderTextByDelta(doc *psrt.Document, pageName string, textIndex, delta int) error {
	p, err := FindPage(doc, pageName)
	if err != nil {
		return err
	}
	_, from, err := FindTextByIndex(p, textIndex)
	if err != nil {
		return err
	}
	texts, err := moveByDelta(p.Texts, from, delta)
	if err != nil {
		return err
	}
	p.Texts = texts
	return nil
}

// ParseTextIndex parses --index flag value.
func ParseTextIndex(s string) (int, error) {
	return strconv.Atoi(s)
}

func findText(doc *psrt.Document, pageName string, textIndex int) (*psrt.Text, error) {
	p, err := FindPage(doc, pageName)
	if err != nil {
		return nil, err
	}
	t, _, err := FindTextByIndex(p, textIndex)
	return t, err
}

func insertAt(texts []psrt.Text, pos int, text psrt.Text) []psrt.Text {
	out := make([]psrt.Text, 0, len(texts)+1)
	out = append(out, texts[:pos]...)
	out = append(out, text)
	out = append(out, texts[pos:]...)
	return out
}

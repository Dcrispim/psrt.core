package psrt

import (
	"fmt"
	"sort"
)

// BlockKind identifies a page block entry.
type BlockKind int

const (
	BlockText BlockKind = iota
	BlockMask
)

// PageBlockEntry is a text or mask block sorted by Index for emission/compile.
type PageBlockEntry struct {
	Kind BlockKind
	Text *Text
	Mask *Mask
}

// PageBlocksByIndex returns all blocks on p sorted by Index ascending.
func PageBlocksByIndex(p *Page) []PageBlockEntry {
	if p == nil {
		return nil
	}
	out := make([]PageBlockEntry, 0, len(p.Texts)+len(p.Masks))
	for i := range p.Texts {
		t := p.Texts[i]
		out = append(out, PageBlockEntry{Kind: BlockText, Text: &t})
	}
	for i := range p.Masks {
		m := p.Masks[i]
		out = append(out, PageBlockEntry{Kind: BlockMask, Mask: &m})
	}
	sort.Slice(out, func(i, j int) bool {
		return pageBlockIndex(out[i]) < pageBlockIndex(out[j])
	})
	return out
}

func pageBlockIndex(e PageBlockEntry) int {
	switch e.Kind {
	case BlockText:
		if e.Text != nil {
			return e.Text.Index
		}
	case BlockMask:
		if e.Mask != nil {
			return e.Mask.Index
		}
	}
	return 0
}

// FindTextByIndex returns the text block with the given index on p.
func FindTextByIndex(p *Page, index int) (*Text, int, error) {
	if p == nil {
		return nil, -1, fmt.Errorf("page is nil")
	}
	for i := range p.Texts {
		if p.Texts[i].Index == index {
			return &p.Texts[i], i, nil
		}
	}
	return nil, -1, fmt.Errorf("text index %d not found on page %q", index, p.Name)
}

// FindMaskByIndex returns the mask block with the given index on p.
func FindMaskByIndex(p *Page, index int) (*Mask, int, error) {
	if p == nil {
		return nil, -1, fmt.Errorf("page is nil")
	}
	for i := range p.Masks {
		if p.Masks[i].Index == index {
			return &p.Masks[i], i, nil
		}
	}
	return nil, -1, fmt.Errorf("mask index %d not found on page %q", index, p.Name)
}

// FindBlockByIndex returns either a text or mask block by index.
func FindBlockByIndex(p *Page, index int) (PageBlockEntry, error) {
	if t, _, err := FindTextByIndex(p, index); err == nil {
		return PageBlockEntry{Kind: BlockText, Text: t}, nil
	}
	if m, _, err := FindMaskByIndex(p, index); err == nil {
		return PageBlockEntry{Kind: BlockMask, Mask: m}, nil
	}
	return PageBlockEntry{}, fmt.Errorf("block index %d not found on page %q", index, p.Name)
}


package textoutline

import (
	"testing"
)

func TestOutlineGoTextProducesPaths(t *testing.T) {
	in := PageInput{
		CanvasW: 100,
		CanvasH: 100,
		Blocks: []BlockInput{{
			Index:     0,
			X:         0,
			Y:         0,
			Width:     50,
			Height:    20,
			PlainText: "Hello",
			Style:     BlockStyle{FontSizePx: 12, ContentW: 50},
		}},
	}
	blocks, err := outlineGoText(in)
	if err != nil {
		t.Fatal(err)
	}
	if len(blocks) != 1 || len(blocks[0].Paths) == 0 {
		t.Fatalf("expected glyph paths, got blocks=%d paths=%d",
			len(blocks), len(blocks[0].Paths))
	}
}

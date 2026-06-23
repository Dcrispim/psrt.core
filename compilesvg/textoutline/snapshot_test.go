package textoutline

import (
	"strings"
	"testing"
)

func TestBuildSnapshotHTML_noDuplicateBlockClass(t *testing.T) {
	html := BuildSnapshotHTML(PageInput{
		Blocks: []BlockInput{{
			Index:      0,
			X:          1,
			Y:          2,
			Width:      100,
			Height:     50,
			ClassAttr:  "psrt-text-capa-0",
			InnerClass: "psrt-text-capa-0-inner",
			TextHTML:   "Hello",
		}},
	})
	if strings.Contains(html, `class="psrt-text-block psrt-text-capa-0"`) {
		t.Fatal("outer block must not repeat typography class")
	}
	if !strings.Contains(html, `class="psrt-text-block"`) {
		t.Fatal("missing outer psrt-text-block")
	}
	if !strings.Contains(html, `class="psrt-text-capa-0"`) {
		t.Fatal("missing inner typography class")
	}
}

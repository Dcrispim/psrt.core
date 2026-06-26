package compilesvg

import (
	"strings"
	"testing"

	"github.com/Dcrispim/psrt.core/psrt"
)

func TestBuildPageStylesheet_includesTextStroke(t *testing.T) {
	css := BuildPageStylesheet(
		"capa",
		psrt.Style(`{}`),
		[]psrt.Text{{
			BaseBlock: psrt.BaseBlock{Index: 0, Width: 50, Style: psrt.Style(`{"strokeWidth":"2px","strokeColor":"#ff0000"}`)},
			TextSize:  5,
		}},
		100, 100,
		nil, nil,
	)
	if !strings.Contains(css, "-webkit-text-stroke-width:2px") {
		t.Fatalf("missing stroke width in stylesheet:\n%s", css)
	}
	if !strings.Contains(css, "-webkit-text-stroke-color:#ff0000") {
		t.Fatalf("missing stroke color in stylesheet:\n%s", css)
	}
}

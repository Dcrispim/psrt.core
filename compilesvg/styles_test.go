package compilesvg

import (
	"strings"
	"testing"

	"psrt/psrt"
)

func TestBuildPageSVGDefsCSS_pageOnly(t *testing.T) {
	css := BuildPageSVGDefsCSS("capa", psrt.Style(`{"color":"#fff","background":"#000"}`))
	if !strings.Contains(css, ".psrt-page-capa{") {
		t.Fatalf("expected page rule:\n%s", css)
	}
	if strings.Contains(css, ".psrt-text-") {
		t.Fatal("defs CSS must not include text block rules")
	}
}

func TestBuildPageStylesheet_orderAndCascade(t *testing.T) {
	css := BuildPageStylesheet("capa", psrt.Style(`{"color":"#fff"}`), []psrt.Text{
		{BaseBlock: psrt.BaseBlock{Index: 0, Style: psrt.Style(`{"color":"#aaa"}`)}, TextSize: 12},
	}, 1080, 1920, nil, nil)

	pageIdx := strings.Index(css, ".psrt-page-capa{")
	textIdx := strings.Index(css, ".psrt-text-capa-0{")
	if pageIdx < 0 || textIdx < 0 {
		t.Fatalf("missing selectors in:\n%s", css)
	}
	if pageIdx > textIdx {
		t.Fatal("page rule must come before text rule for cascade")
	}
	if !strings.Contains(css, "color:#fff") || !strings.Contains(css, "color:#aaa") {
		t.Fatalf("expected colors in css:\n%s", css)
	}
}

func TestBuildPageStylesheet_textBackground(t *testing.T) {
	css := BuildPageStylesheet("intro", psrt.Style(`{}`), []psrt.Text{
		{BaseBlock: psrt.BaseBlock{Index: 1, Width: 77, Style: psrt.Style(`{"background":"#000000ff","padding":"10px","text-align":"center","font-weight":"600"}`)}, TextSize: 5},
	}, 1080, 1920, nil, nil)
	if strings.Contains(css, "background-color:#000000ff") {
		t.Fatalf("box background must be on SVG rect, not text class CSS:\n%s", css)
	}
	if strings.Contains(css, ".psrt-text-intro-1-bg{") {
		t.Fatalf("must not use separate SVG bg class:\n%s", css)
	}
	if !strings.Contains(css, "width:100%") || !strings.Contains(css, "min-height:100%") {
		t.Fatal("text block should fill foreignObject (width/min-height 100%)")
	}
	if !strings.Contains(css, "display:flex") || !strings.Contains(css, "justify-content:center") {
		t.Fatalf("centered text must use flex in foreignObject:\n%s", css)
	}
	if !strings.Contains(css, ".psrt-text-intro-1-inner{display:block") {
		t.Fatalf("centered text must wrap inline markup in inner span CSS:\n%s", css)
	}
	if strings.Contains(css, "padding:0.15em") || strings.Contains(css, "-ref-img") {
		t.Fatalf("must not add default padding or ref-img rules:\n%s", css)
	}
	if !strings.Contains(css, "padding:10px") {
		t.Fatalf("text block padding must be in CSS for foreignObject inset:\n%s", css)
	}
}

func TestBuildPageStylesheet_individualBorderRadiusLonghands(t *testing.T) {
	css := BuildPageStylesheet("p", psrt.Style(`{}`), []psrt.Text{
		{
			BaseBlock: psrt.BaseBlock{Index: 0, Style: psrt.Style(`{
				"background":"#edeee2",
				"border-top-left-radius":"0px",
				"border-top-right-radius":"0px",
				"border-bottom-right-radius":"48px",
				"border-bottom-left-radius":"32px"
			}`)},
			TextSize: 3,
		},
	}, 1080, 1920, nil, nil)
	if !strings.Contains(css, "border-bottom-right-radius:48px") {
		t.Fatalf("expected individual corner radius in SVG CSS:\n%s", css)
	}
	if !strings.Contains(css, "border-bottom-left-radius:32px") {
		t.Fatalf("expected bottom-left radius in SVG CSS:\n%s", css)
	}
	if strings.Contains(css, "border-radius:0px 0px") {
		t.Fatalf("should use longhands, not broken shorthand with stripped zeros:\n%s", css)
	}
}

func TestBuildPageStylesheet_textAlignRightAndJustify(t *testing.T) {
	css := BuildPageStylesheet("p", psrt.Style(`{}`), []psrt.Text{
		{BaseBlock: psrt.BaseBlock{Index: 0, Width: 80, Style: psrt.Style(`{"text-align":"right","color":"#fff"}`)}, TextSize: 5},
		{BaseBlock: psrt.BaseBlock{Index: 1, Width: 80, Style: psrt.Style(`{"text-align":"justify","color":"#fff"}`)}, TextSize: 5},
	}, 1080, 1920, nil, nil)
	if !strings.Contains(css, "text-align:right") {
		t.Fatalf("expected text-align:right in SVG CSS:\n%s", css)
	}
	if !strings.Contains(css, "align-items:flex-end") {
		t.Fatalf("right align needs flex-end in SVG CSS:\n%s", css)
	}
	if !strings.Contains(css, "text-align:justify") {
		t.Fatalf("expected text-align:justify in SVG CSS:\n%s", css)
	}
}

func TestBuildPageStylesheet_percentPaddingResolved(t *testing.T) {
	css := BuildPageStylesheet("p", psrt.Style(`{}`), []psrt.Text{
		{BaseBlock: psrt.BaseBlock{Index: 0, Style: psrt.Style(`{"padding":"0.461%","text-align":"center"}`)}, TextSize: 3},
	}, 700, 6585, nil, nil)
	if strings.Contains(css, "padding:0.461%") {
		t.Fatalf("percent padding must be px in SVG CSS:\n%s", css)
	}
	if !strings.Contains(css, "padding:") || !strings.Contains(css, "px") {
		t.Fatalf("expected resolved px padding:\n%s", css)
	}
}

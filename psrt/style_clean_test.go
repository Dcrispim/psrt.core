package psrt

import "testing"

func TestCleanEmptyTextBlockStyles_removesColorAlignPadding(t *testing.T) {
	doc := Document{
		Pages: []Page{{
			Name: "p1",
			Texts: []Text{
				{
					BaseBlock: BaseBlock{Index: 1, Style: Style(`{"color":"#000","text-align":"center","padding":"1%"}`)},
					Content:   "hello",
				},
				{
					BaseBlock: BaseBlock{Index: 2, Style: Style(`{"color":"#000","text-align":"center","padding":"0%","background":"#fff"}`)},
					Content:   "   \n  ",
				},
			},
		}},
	}
	CleanEmptyTextBlockStyles(&doc)

	got := string(doc.Pages[0].Texts[0].Style)
	want := `{"color":"#000","text-align":"center","padding":"1%"}`
	if got != want {
		t.Fatalf("non-empty block style changed:\ngot  %s\nwant %s", got, want)
	}

	got = string(doc.Pages[0].Texts[1].Style)
	want = `{"background":"#fff"}`
	if got != want {
		t.Fatalf("empty block style:\ngot  %s\nwant %s", got, want)
	}
}

func TestCleanEmptyTextBlockStyles_aliases(t *testing.T) {
	doc := Document{
		Pages: []Page{{
			Texts: []Text{{
				BaseBlock: BaseBlock{Style: Style(`{"color":"#fff","textAlign":"right","ta":"left","padding":"2px"}`)},
				Content:   "",
			}},
		}},
	}
	CleanEmptyTextBlockStyles(&doc)
	if string(doc.Pages[0].Texts[0].Style) != "{}" {
		t.Fatalf("expected empty style object, got %s", doc.Pages[0].Texts[0].Style)
	}
}

func TestIsTextBlockEmpty(t *testing.T) {
	if !IsTextBlockEmpty("") {
		t.Fatal("expected empty")
	}
	if !IsTextBlockEmpty("  \n\t  ") {
		t.Fatal("expected whitespace-only as empty")
	}
	if IsTextBlockEmpty("x") {
		t.Fatal("expected non-empty")
	}
}

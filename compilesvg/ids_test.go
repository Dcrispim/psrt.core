package compilesvg

import "testing"

func TestSlug(t *testing.T) {
	if got := Slug("Mood Sexta"); got != "mood-sexta" {
		t.Fatalf("Slug = %q, want mood-sexta", got)
	}
	if got := Slug(""); got != "page" {
		t.Fatalf("empty Slug = %q", got)
	}
}

func TestUniqueSlugs(t *testing.T) {
	got := UniqueSlugs([]string{"capa", "capa", "other"})
	want := []string{"capa", "capa-2", "other"}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("UniqueSlugs[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}

func TestPageAndTextIDs(t *testing.T) {
	s := "capa"
	if PageID(s) != "psrt-page-capa" {
		t.Fatal(PageID(s))
	}
	if TextID(s, 1) != "psrt-text-capa-1" {
		t.Fatal(TextID(s, 1))
	}
	if TextClassAttr(s, 0) != "psrt-page-capa psrt-text-capa-0" {
		t.Fatal(TextClassAttr(s, 0))
	}
	if TextInnerClass(s, 0) != "psrt-text-capa-0-inner" {
		t.Fatal(TextInnerClass(s, 0))
	}
}

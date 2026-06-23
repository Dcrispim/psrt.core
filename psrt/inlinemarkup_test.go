package psrt

import (
	"strings"
	"testing"
)

func TestRenderInlineHTML_bold(t *testing.T) {
	got := RenderInlineHTML("**bold**")
	want := "<strong>bold</strong>"
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

func TestRenderInlineHTML_italic(t *testing.T) {
	got := RenderInlineHTML("*italic*")
	want := "<em>italic</em>"
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

func TestRenderInlineHTML_boldItalic(t *testing.T) {
	got := RenderInlineHTML("***both***")
	want := "<strong><em>both</em></strong>"
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

func TestRenderInlineHTML_underlineStrike(t *testing.T) {
	got := RenderInlineHTML("_u_ ~s~")
	if !strings.Contains(got, "<u>u</u>") || !strings.Contains(got, "<s>s</s>") {
		t.Fatalf("got %q", got)
	}
}

func TestRenderInlineHTML_unclosedLiteral(t *testing.T) {
	got := RenderInlineHTML("**open")
	if !strings.Contains(got, "**") {
		t.Fatalf("expected literal asterisks, got %q", got)
	}
}

func TestRenderInlineHTML_escape(t *testing.T) {
	got := RenderInlineHTML(`\*\*literal\*\*`)
	if got != "**literal**" {
		t.Fatalf("got %q", got)
	}
}

func TestRenderInlineHTML_multiline(t *testing.T) {
	got := RenderInlineHTML("a\n**b**")
	want := "a<br/><strong>b</strong>"
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

func TestRenderInlineHTML_ampersand(t *testing.T) {
	got := RenderInlineHTML("a & b")
	if got != "a &amp; b" {
		t.Fatalf("got %q", got)
	}
}

func TestPlainTextForLayout(t *testing.T) {
	got := PlainTextForLayout("**a** ~b~")
	want := "a b"
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

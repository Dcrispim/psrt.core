package compileasset

import "testing"

func TestParseGoogleFontFamilies_css2(t *testing.T) {
	url := "https://fonts.googleapis.com/css2?family=Roboto:wght@400&family=Schoolbell&display=swap"
	got := ParseGoogleFontFamilies(url)
	if len(got) != 2 || got[0] != "Roboto" || got[1] != "Schoolbell" {
		t.Fatalf("got %v", got)
	}
}

func TestIsGoogleFontsCSSURL(t *testing.T) {
	if !IsGoogleFontsCSSURL("https://fonts.googleapis.com/css2?family=Roboto&display=swap") {
		t.Fatal("expected google css url")
	}
	if IsGoogleFontsCSSURL("https://cdn.example/font.woff2") {
		t.Fatal("woff2 is not css api")
	}
}

func TestFontFamilyNameForURL_google(t *testing.T) {
	url := "https://fonts.googleapis.com/css2?family=Schoolbell&display=swap"
	if got := FontFamilyNameForURL(url, 0); got != "Schoolbell" {
		t.Fatalf("got %q", got)
	}
}

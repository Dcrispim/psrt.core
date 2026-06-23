package editor

import "testing"

func TestAddRemoveFont(t *testing.T) {
	doc := sampleDoc(t)
	url := "https://fonts.example/new.woff2"
	if err := AddFont(&doc, url); err != nil {
		t.Fatal(err)
	}
	if err := AddFont(&doc, url); err != nil {
		t.Fatal(err)
	}
	if len(doc.Fonts) != 1 {
		t.Fatalf("fonts: %v", doc.Fonts)
	}
	if err := RemoveFont(&doc, url); err != nil {
		t.Fatal(err)
	}
	if len(doc.Fonts) != 0 {
		t.Fatalf("fonts after remove: %v", doc.Fonts)
	}
}

package editor

import "testing"

func TestAddAndRemovePage(t *testing.T) {
	doc := sampleDoc(t)
	newPage := doc.Pages[0]
	newPage.Name = "extra"
	newPage.Texts = nil
	if err := AddPage(&doc, newPage, "", ""); err != nil {
		t.Fatal(err)
	}
	if _, err := FindPage(&doc, "extra"); err != nil {
		t.Fatal(err)
	}
	if err := RemovePage(&doc, "extra"); err != nil {
		t.Fatal(err)
	}
	if _, err := FindPage(&doc, "extra"); err == nil {
		t.Fatal("page should be removed")
	}
}

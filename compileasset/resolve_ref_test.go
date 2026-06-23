package compileasset

import "testing"

func TestResolveAssetReference_expandsConstPrefix(t *testing.T) {
	consts := map[string]string{
		"baseURL": "file:///D:/images/chapter/",
	}
	got := ResolveAssetReference("@baseURL@01.webp", consts)
	want := "file:///D:/images/chapter/01.webp"
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
	if !IsAssetReference(got) {
		t.Fatalf("resolved ref should be asset reference: %q", got)
	}
}

func TestResolveAssetReference_noConsts(t *testing.T) {
	raw := "@baseURL@01.webp"
	if got := ResolveAssetReference(raw, nil); got != raw {
		t.Fatalf("got %q", got)
	}
}

package webconnector

import "testing"

func TestNormalizeAllowedOriginsList(t *testing.T) {
	got, err := normalizeAllowedOriginsList("https://dcrispim.github.io/psrt-gui-web/, http://localhost:5174")
	if err != nil {
		t.Fatal(err)
	}
	want := "https://dcrispim.github.io,http://localhost:5174"
	if joinAllowedOrigins(got) != want {
		t.Fatalf("got %q want %q", joinAllowedOrigins(got), want)
	}
}

func TestOriginsMatch(t *testing.T) {
	allowed := "https://dcrispim.github.io,http://localhost:5174"
	if !originsMatch("https://dcrispim.github.io", allowed) {
		t.Fatal("expected github origin to match")
	}
	if !originsMatch("http://localhost:5174", allowed) {
		t.Fatal("expected localhost origin to match")
	}
	if originsMatch("https://evil.example", allowed) {
		t.Fatal("unexpected match")
	}
}

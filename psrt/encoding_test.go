package psrt

import "testing"

func TestFixMisencodedUTF8(t *testing.T) {
	const want = "Mesmo se a gente matar você agora, o Mestre de seita não pode nos tocar."
	// UTF-8 read as Latin-1 / Windows-1252
	broken := "Mesmo se a gente matar voc\u00c3\u00aa agora, o Mestre de seita n\u00c3\u00a3o pode nos tocar."
	if got := FixMisencodedUTF8(broken); got != want {
		t.Fatalf("got %q want %q", got, want)
	}
	if got := FixMisencodedUTF8(want); got != want {
		t.Fatalf("already valid: got %q", got)
	}
}

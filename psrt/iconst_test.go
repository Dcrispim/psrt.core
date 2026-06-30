package psrt

import (
	"strings"
	"testing"
)

const iconstSrc = `$START p1 | {} | https://img.example/a.png
>>10,10,3,40 | {} | 0
Considere uma @link:Nosso discord@ — fala com @desc:Akira@.
$END p1

$CONSTS
@ note | plain value
@desc:Akira | Personagem principal, introduzido no capítulo 1
@link:Nosso discord | https://discord.example/invite
$ENDCONSTS
`

func TestParseInteractiveConsts(t *testing.T) {
	doc, err := ParseString(iconstSrc)
	if err != nil {
		t.Fatal(err)
	}
	if doc.Consts["note"] != "plain value" {
		t.Fatalf("plain const broken: %q", doc.Consts["note"])
	}
	link, ok := doc.IConst["link:Nosso discord"]
	if !ok {
		t.Fatalf("iConst keys: %+v", doc.IConst)
	}
	if link.Type != "link" || link.Render != "Nosso discord" || link.Value != "https://discord.example/invite" {
		t.Fatalf("link entry: %+v", link)
	}
	desc := doc.IConst["desc:Akira"]
	if desc.Type != "desc" || desc.Render != "Akira" || desc.Value != "Personagem principal, introduzido no capítulo 1" {
		t.Fatalf("desc entry: %+v", desc)
	}
}

func TestInteractiveConstNotMisreadAsPlain(t *testing.T) {
	// A plain const whose VALUE contains a colon must stay plain.
	doc, err := ParseString("$CONSTS\n@ shadow | \"boxShadow\":\"0 0 1px\"\n$ENDCONSTS\n")
	if err != nil {
		t.Fatal(err)
	}
	if len(doc.IConst) != 0 {
		t.Fatalf("value colon misread as interactive: %+v", doc.IConst)
	}
	if doc.Consts["shadow"] != `"boxShadow":"0 0 1px"` {
		t.Fatalf("shadow: %q", doc.Consts["shadow"])
	}
}

func TestInteractiveConstRoundTrip(t *testing.T) {
	doc, err := ParseString(iconstSrc)
	if err != nil {
		t.Fatal(err)
	}
	out, err := FormatPSRT(doc, false)
	if err != nil {
		t.Fatal(err)
	}
	reparsed, err := ParseString(string(out))
	if err != nil {
		t.Fatalf("reparse: %v\n%s", err, out)
	}
	if len(reparsed.IConst) != 2 {
		t.Fatalf("round-trip lost iConst: %+v", reparsed.IConst)
	}
	if reparsed.IConst["link:Nosso discord"] != doc.IConst["link:Nosso discord"] {
		t.Fatalf("link mismatch: %+v vs %+v", reparsed.IConst, doc.IConst)
	}
	if reparsed.IConst["desc:Akira"] != doc.IConst["desc:Akira"] {
		t.Fatalf("desc mismatch: %+v", reparsed.IConst)
	}
}

func TestInteractiveConstSnapshotStrip(t *testing.T) {
	doc, err := ParseString(iconstSrc)
	if err != nil {
		t.Fatal(err)
	}
	consts := ConstsWithInteractive(doc.Consts, doc.IConst)
	got := ExpandConsts(doc.Pages[0].Texts[0].Content, consts)
	if strings.Contains(got, "@link:") || strings.Contains(got, "@desc:") {
		t.Fatalf("interactive tokens leaked into snapshot: %q", got)
	}
	if !strings.Contains(got, "Nosso discord") || !strings.Contains(got, "Akira") {
		t.Fatalf("base render text missing: %q", got)
	}
}

func TestDuplicateInteractiveConst(t *testing.T) {
	_, err := ParseString("$CONSTS\n@link:x | a\n@link:x | b\n$ENDCONSTS\n")
	if err == nil {
		t.Fatal("expected duplicate error")
	}
}

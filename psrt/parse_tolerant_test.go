package psrt

import (
	"strings"
	"testing"
)

func TestParseTolerant_discardsInvalidBlockKeepsRest(t *testing.T) {
	src := `$START p | {} | u
~~10-10-20-5 | {} | 0
M0,0 L10,10 Z
~~30-30-20-5 | {} | 1
$END p`
	// The second ~~ has an empty body (malformed). The strict parser must
	// fail the whole document; the tolerant one must keep the first block
	// and report the second as an accumulated error.
	if _, err := ParseString(src); err == nil {
		t.Fatal("strict Parse should fail on the malformed second ~~ block")
	}

	doc, errs := ParseTolerant(strings.NewReader(src))
	if len(errs) != 1 {
		t.Fatalf("errs: got %d want 1: %+v", len(errs), errs)
	}
	if !strings.Contains(errs[0].Message, "path mask body is empty") {
		t.Fatalf("unexpected error: %+v", errs[0])
	}
	if len(doc.Pages) != 1 || len(doc.Pages[0].PathMasks) != 1 {
		t.Fatalf("expected 1 surviving path mask, got: %+v", doc.Pages)
	}
	if doc.Pages[0].PathMasks[0].Index != 0 {
		t.Fatalf("unexpected surviving path mask: %+v", doc.Pages[0].PathMasks[0])
	}
}

func TestParseTolerant_unknownMarkerDiscarded(t *testing.T) {
	src := `$START p | {} | u
##not-a-real-marker | {} | 1
>>10-10-20-3 | {} | 0
hello
$END p`
	if _, err := ParseString(src); err == nil {
		t.Fatal("strict Parse should fail on the unrecognized ## marker")
	}

	doc, errs := ParseTolerant(strings.NewReader(src))
	if len(errs) == 0 {
		t.Fatal("expected at least one accumulated error for the unknown marker")
	}
	if len(doc.Pages) != 1 || len(doc.Pages[0].Texts) != 1 || doc.Pages[0].Texts[0].Content != "hello" {
		t.Fatalf("expected the >> block to survive: %+v", doc.Pages)
	}
}

func TestParseTolerant_strictModeUnaffected(t *testing.T) {
	src := `$START p | {} | u
~~10-10-20-5 | {} | 0
M0,0 L10,10 Z
$END p`
	doc, err := ParseString(src)
	if err != nil {
		t.Fatal(err)
	}
	if len(doc.Pages[0].PathMasks) != 1 {
		t.Fatalf("strict Parse must still work for valid documents: %+v", doc.Pages[0])
	}
}

func TestParseTolerant_structuralErrorStillSurfaces(t *testing.T) {
	src := `$START p | {} | u
>>10-10-20-3 | {} | 0
hello`
	// No matching $END before EOF — a structural error, not a per-block one;
	// must still be reported even in tolerant mode.
	_, errs := ParseTolerant(strings.NewReader(src))
	if len(errs) == 0 {
		t.Fatal("expected the structural EOF error to be reported")
	}
	found := false
	for _, e := range errs {
		if strings.Contains(e.Message, "EOF with open page") {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected an EOF-with-open-page error, got: %+v", errs)
	}
}

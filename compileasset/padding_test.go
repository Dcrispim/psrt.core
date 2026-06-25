package compileasset

import (
	"testing"

	"github.com/Dcrispim/psrt.core/psrt"
)

func TestParsePaddingInsets_shorthand(t *testing.T) {
	p := ParsePaddingInsets(psrt.Style(`{"padding":"10px"}`), 48)
	if p.Top != 10 || p.Horizontal() != 20 {
		t.Fatalf("got %+v", p)
	}
}

func TestParsePaddingInsets_fourSides(t *testing.T) {
	p := ParsePaddingInsets(psrt.Style(`{"padding":"10px 20px 30px 40px"}`), 48)
	if p.Top != 10 || p.Right != 20 || p.Bottom != 30 || p.Left != 40 {
		t.Fatalf("got %+v", p)
	}
}

func TestTextBoxInsets_paddingAndBorder(t *testing.T) {
	in := TextBoxInsets(psrt.Style(`{"padding":"10px","border":"2px solid #fff"}`), 48)
	if in.Horizontal() != 24 || in.Vertical() != 24 {
		t.Fatalf("got h=%v v=%v", in.Horizontal(), in.Vertical())
	}
}

func TestTextBoxInsetsForCanvas_percentPadding(t *testing.T) {
	in := TextBoxInsetsForCanvas(psrt.Style(`{"padding":"0.515%"}`), 48, 700, 6585)
	// single-value padding % uses max(W,H) as base
	want := 6585 * 0.00515
	if in.Top < want-0.01 || in.Top > want+0.01 {
		t.Fatalf("top padding: got %v want ~%v", in.Top, want)
	}
	if in.Horizontal() < 2*want-0.02 || in.Horizontal() > 2*want+0.02 {
		t.Fatalf("horizontal padding: got %v want ~%v", in.Horizontal(), 2*want)
	}
}

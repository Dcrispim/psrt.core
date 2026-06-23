package psrt

import "testing"

func TestExpandConstsInStyle(t *testing.T) {
	consts := map[string]string{
		"accent": "#1DB954",
	}
	style := Style(`{"color":"@accent@"}`)
	out, err := ExpandConstsInStyle(style, consts)
	if err != nil {
		t.Fatal(err)
	}
	if string(out) != `{"color":"#1DB954"}` {
		t.Fatalf("got %s", out)
	}
}

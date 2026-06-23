package editor

import (
	"encoding/json"
	"testing"

	"psrt/psrt"
)

func TestSetStyleKeyAndRemove(t *testing.T) {
	style := psrt.Style(`{"color":"#fff"}`)
	updated, err := SetStyleKey(style, "fontWeight", `"700"`)
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]any
	if err := json.Unmarshal(updated, &m); err != nil {
		t.Fatal(err)
	}
	if m["fontWeight"] != "700" {
		t.Fatalf("fontWeight: %v", m["fontWeight"])
	}
	updated, err = RemoveStyleKey(updated, "color")
	if err != nil {
		t.Fatal(err)
	}
	m = make(map[string]any)
	if err := json.Unmarshal(updated, &m); err != nil {
		t.Fatal(err)
	}
	if _, ok := m["color"]; ok {
		t.Fatal("color should be removed")
	}
}

func TestParseStyleValue(t *testing.T) {
	cases := []struct {
		in   string
		want any
	}{
		{`"#fff"`, "#fff"},
		{"#fff", "#fff"},
		{"#ffffffff", "#ffffffff"},
		{`600`, float64(600)},
		{"center", "center"},
		{`true`, true},
	}
	for _, tc := range cases {
		got, err := parseStyleValue(tc.in)
		if err != nil {
			t.Fatalf("parseStyleValue(%q): %v", tc.in, err)
		}
		if got != tc.want {
			t.Errorf("parseStyleValue(%q) = %v (%T), want %v (%T)", tc.in, got, got, tc.want, tc.want)
		}
	}
}

func TestSetStyleKeyBareHex(t *testing.T) {
	style := psrt.Style(`{}`)
	updated, err := SetStyleKey(style, "color", "#ffffffff")
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]any
	if err := json.Unmarshal(updated, &m); err != nil {
		t.Fatal(err)
	}
	if m["color"] != "#ffffffff" {
		t.Fatalf("color: %v", m["color"])
	}
}

func TestMergeStyle(t *testing.T) {
	style := psrt.Style(`{"color":"#fff"}`)
	updated, err := MergeStyle(style, json.RawMessage(`{"fontWeight":"600"}`))
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]any
	if err := json.Unmarshal(updated, &m); err != nil {
		t.Fatal(err)
	}
	if m["color"] != "#fff" || m["fontWeight"] != "600" {
		t.Fatalf("merge: %v", m)
	}
}

package editor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Dcrispim/psrt.core/psrt"
)

// AddConst registers a constant and replaces literal value occurrences with @name@.
func AddConst(doc *psrt.Document, name, value string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return fmt.Errorf("const name is empty")
	}
	if doc.Consts == nil {
		doc.Consts = make(map[string]string)
	}
	if _, exists := doc.Consts[name]; exists {
		return fmt.Errorf("const %q already exists", name)
	}
	doc.Consts[name] = value
	if value != "" {
		SubstituteConstReferences(doc, name, value)
	}
	return nil
}

// RemoveConst reverts @name@ to the stored value, then deletes the constant.
func RemoveConst(doc *psrt.Document, name string) error {
	if doc.Consts == nil {
		return fmt.Errorf("const %q not found", name)
	}
	value, ok := doc.Consts[name]
	if !ok {
		return fmt.Errorf("const %q not found", name)
	}
	RevertConstReferences(doc, name, value)
	delete(doc.Consts, name)
	return nil
}

// SubstituteConstReferences replaces literal value with @name@ across the document.
func SubstituteConstReferences(doc *psrt.Document, name, value string) {
	placeholder := "@" + name + "@"
	for i := range doc.Pages {
		p := &doc.Pages[i]
		p.ImageURL = strings.ReplaceAll(p.ImageURL, value, placeholder)
		p.Style = substituteInStyle(p.Style, value, placeholder)
		for j := range p.Texts {
			t := &p.Texts[j]
			t.Content = strings.ReplaceAll(t.Content, value, placeholder)
			t.ImageRef = strings.ReplaceAll(t.ImageRef, value, placeholder)
			t.Style = substituteInStyle(t.Style, value, placeholder)
		}
	}
}

// RevertConstReferences replaces @name@ with the original value across the document.
func RevertConstReferences(doc *psrt.Document, name, value string) {
	placeholder := "@" + name + "@"
	for i := range doc.Pages {
		p := &doc.Pages[i]
		p.ImageURL = strings.ReplaceAll(p.ImageURL, placeholder, value)
		p.Style = substituteInStyle(p.Style, placeholder, value)
		for j := range p.Texts {
			t := &p.Texts[j]
			t.Content = strings.ReplaceAll(t.Content, placeholder, value)
			t.ImageRef = strings.ReplaceAll(t.ImageRef, placeholder, value)
			t.Style = substituteInStyle(t.Style, placeholder, value)
		}
	}
}

func substituteInStyle(style psrt.Style, old, new string) psrt.Style {
	raw := strings.TrimSpace(string(style))
	if raw == "" || raw == "{}" {
		return style
	}
	compact := compactJSONString(raw)
	replaced := strings.ReplaceAll(compact, old, new)
	if replaced == compact {
		return style
	}
	if !json.Valid([]byte(replaced)) {
		return psrt.Style(replaced)
	}
	return psrt.Style(replaced)
}

func compactJSONString(s string) string {
	var buf bytes.Buffer
	if err := json.Compact(&buf, []byte(s)); err != nil {
		return strings.ReplaceAll(strings.ReplaceAll(s, "\n", ""), "\r", "")
	}
	return buf.String()
}

// SaveDocument formats and writes a PSRT document to path.
func SaveDocument(doc *psrt.Document, path string) error {
	data, err := psrt.FormatPSRT(*doc, false)
	if err != nil {
		return fmt.Errorf("format: %w", err)
	}
	return writeFile(path, data)
}

// LoadDocument parses a PSRT file from path.
func LoadDocument(path string) (psrt.Document, error) {
	data, err := readFile(path)
	if err != nil {
		return psrt.Document{}, err
	}
	doc, err := psrt.Parse(bytes.NewReader(data))
	if err != nil {
		return psrt.Document{}, fmt.Errorf("parse: %w", err)
	}
	return doc, nil
}

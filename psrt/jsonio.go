package psrt

import (
	"encoding/json"
)

// ParseJSON decodes a Document from JSON (inverse of ToJSON).
func ParseJSON(b []byte) (Document, error) {
	var doc Document
	if err := json.Unmarshal(b, &doc); err != nil {
		return doc, err
	}
	if doc.Consts == nil {
		doc.Consts = make(map[string]string)
	}
	if doc.IConst == nil {
		doc.IConst = make(map[string]InteractiveConst)
	}
	if doc.Sources == nil {
		doc.Sources = make(map[string]string)
	}
	if doc.Fonts == nil {
		doc.Fonts = []string{}
	}
	for i := range doc.Pages {
		if doc.Pages[i].Texts == nil {
			doc.Pages[i].Texts = []Text{}
		}
		if doc.Pages[i].Masks == nil {
			doc.Pages[i].Masks = []Mask{}
		}
	}
	return doc, nil
}

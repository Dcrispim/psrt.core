package compilesvg

import (
	"fmt"
	"strings"

	"github.com/Dcrispim/psrt.core/psrt"
)

// ResolveDocument expands all @const@ placeholders in styles, content, and URLs.
func ResolveDocument(doc psrt.Document) psrt.Document {
	consts := doc.Consts
	if len(consts) == 0 {
		return doc
	}
	out := doc
	out.Pages = make([]psrt.Page, len(doc.Pages))
	for i := range doc.Pages {
		p := doc.Pages[i]
		p.Style, _ = psrt.ExpandConstsInStyle(p.Style, consts)
		p.ImageURL = psrt.ExpandConsts(strings.TrimSpace(p.ImageURL), consts)
		p.Texts = make([]psrt.Text, len(doc.Pages[i].Texts))
		for j := range doc.Pages[i].Texts {
			t := doc.Pages[i].Texts[j]
			t.Style, _ = psrt.ExpandConstsInStyle(t.Style, consts)
			t.Content = psrt.NormalizeTextContent(psrt.ExpandConsts(t.Content, consts))
			t.ImageRef = psrt.ExpandConsts(strings.TrimSpace(t.ImageRef), consts)
			p.Texts[j] = t
		}
		p.Masks = make([]psrt.Mask, len(doc.Pages[i].Masks))
		for j := range doc.Pages[i].Masks {
			m := doc.Pages[i].Masks[j]
			m.Style, _ = psrt.ExpandConstsInStyle(m.Style, consts)
			m.ImageRef = psrt.ExpandConsts(strings.TrimSpace(m.ImageRef), consts)
			p.Masks[j] = m
		}
		p.PathMasks = make([]psrt.PathMask, len(doc.Pages[i].PathMasks))
		for j := range doc.Pages[i].PathMasks {
			pm := doc.Pages[i].PathMasks[j]
			pm.Style, _ = psrt.ExpandConstsInStyle(pm.Style, consts)
			pm.ImageRef = psrt.ExpandConsts(strings.TrimSpace(pm.ImageRef), consts)
			p.PathMasks[j] = pm
		}
		out.Pages[i] = p
	}
	return out
}

// ResolveDocumentStrict is like ResolveDocument but returns an error on invalid style JSON.
func ResolveDocumentStrict(doc psrt.Document) (psrt.Document, error) {
	consts := doc.Consts
	if len(consts) == 0 {
		return doc, nil
	}
	out := doc
	out.Pages = make([]psrt.Page, len(doc.Pages))
	for i := range doc.Pages {
		p := doc.Pages[i]
		var err error
		p.Style, err = psrt.ExpandConstsInStyle(p.Style, consts)
		if err != nil {
			return doc, fmt.Errorf("page %q style: %w", p.Name, err)
		}
		p.ImageURL = psrt.ExpandConsts(strings.TrimSpace(p.ImageURL), consts)
		p.Texts = make([]psrt.Text, len(doc.Pages[i].Texts))
		for j := range doc.Pages[i].Texts {
			t := doc.Pages[i].Texts[j]
			t.Style, err = psrt.ExpandConstsInStyle(t.Style, consts)
			if err != nil {
				return doc, fmt.Errorf("page %q text %d style: %w", p.Name, t.Index, err)
			}
			t.Content = psrt.NormalizeTextContent(psrt.ExpandConsts(t.Content, consts))
			t.ImageRef = psrt.ExpandConsts(strings.TrimSpace(t.ImageRef), consts)
			p.Texts[j] = t
		}
		p.Masks = make([]psrt.Mask, len(doc.Pages[i].Masks))
		for j := range doc.Pages[i].Masks {
			m := doc.Pages[i].Masks[j]
			m.Style, err = psrt.ExpandConstsInStyle(m.Style, consts)
			if err != nil {
				return doc, fmt.Errorf("page %q mask %d style: %w", p.Name, m.Index, err)
			}
			m.ImageRef = psrt.ExpandConsts(strings.TrimSpace(m.ImageRef), consts)
			p.Masks[j] = m
		}
		p.PathMasks = make([]psrt.PathMask, len(doc.Pages[i].PathMasks))
		for j := range doc.Pages[i].PathMasks {
			pm := doc.Pages[i].PathMasks[j]
			pm.Style, err = psrt.ExpandConstsInStyle(pm.Style, consts)
			if err != nil {
				return doc, fmt.Errorf("page %q path mask %d style: %w", p.Name, pm.Index, err)
			}
			pm.ImageRef = psrt.ExpandConsts(strings.TrimSpace(pm.ImageRef), consts)
			p.PathMasks[j] = pm
		}
		out.Pages[i] = p
	}
	return out, nil
}

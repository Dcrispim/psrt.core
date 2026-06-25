package psrt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"psrt/svgpath"
)

// FormatPSRT serialises a full document to PSRT syntax.
func FormatPSRT(doc Document, compact bool) ([]byte, error) {
	var b strings.Builder
	for i := range doc.Pages {
		if err := writePage(&b, &doc.Pages[i], compact, doc.Consts); err != nil {
			return nil, err
		}
	}
	if len(doc.Fonts) > 0 {
		b.WriteString("$FONTS\n")
		for _, f := range doc.Fonts {
			b.WriteString(strings.TrimSpace(f))
			b.WriteByte('\n')
		}
		b.WriteString("$ENDFONTS\n")
	}
	if len(doc.Consts) > 0 {
		keys := sortedStringKeys(doc.Consts)
		b.WriteString("$CONSTS\n")
		for _, k := range keys {
			b.WriteString("@ ")
			b.WriteString(k)
			b.WriteString(pipeSep)
			b.WriteString(doc.Consts[k])
			b.WriteByte('\n')
		}
		b.WriteString("$ENDCONSTS\n")
	}
	if len(doc.Sources) > 0 {
		keys := sortedStringKeys(doc.Sources)
		b.WriteString("$SOURCE\n")
		for _, k := range keys {
			b.WriteString(k)
			b.WriteString(pipeSep)
			b.WriteString(doc.Sources[k])
			b.WriteByte('\n')
		}
		b.WriteString("$ENDSOURCE\n")
	}
	return []byte(b.String()), nil
}

// FormatPagePSRT writes a single page as a PSRT fragment ($START … $END).
func FormatPagePSRT(p *Page) ([]byte, error) {
	if p == nil {
		return nil, fmt.Errorf("page is nil")
	}
	var b strings.Builder
	if err := writePage(&b, p, false, nil); err != nil {
		return nil, err
	}
	return []byte(b.String()), nil
}

// FormatTextPSRT writes one text block (>> header + indented body).
func FormatTextPSRT(t *Text) ([]byte, error) {
	if t == nil {
		return nil, fmt.Errorf("text is nil")
	}
	var b strings.Builder
	if err := writeText(&b, t, false, nil); err != nil {
		return nil, err
	}
	return []byte(b.String()), nil
}

// FormatConstPSRT writes one constant line suitable inside $CONSTS … $ENDCONSTS.
func FormatConstPSRT(name, value string) []byte {
	var b strings.Builder
	b.WriteString("@ ")
	b.WriteString(name)
	b.WriteString(pipeSep)
	b.WriteString(value)
	b.WriteByte('\n')
	return []byte(b.String())
}

func writePage(b *strings.Builder, p *Page, compact bool, consts map[string]string) error {
	style := string(p.Style)
	if strings.TrimSpace(style) == "" {
		style = "{}"
	}
	if !json.Valid([]byte(style)) {
		return fmt.Errorf("page %q: style is not valid JSON", p.Name)
	}
	b.WriteString("$START ")
	b.WriteString(p.Name)
	b.WriteString(pipeSep)
	b.WriteString(style)
	b.WriteString(pipeSep)
	b.WriteString(p.ImageURL)
	b.WriteByte('\n')
	entries := PageBlocksByIndex(p)
	for _, e := range entries {
		switch e.Kind {
		case BlockText:
			if e.Text != nil {
				if err := writeText(b, e.Text, compact, consts); err != nil {
					return err
				}
			}
		case BlockMask:
			if e.Mask != nil {
				if err := writeMask(b, e.Mask, consts); err != nil {
					return err
				}
			}
		case BlockPathMask:
			if e.PathMask != nil {
				if err := writePathMask(b, e.PathMask, consts); err != nil {
					return err
				}
			}
		}
	}
	b.WriteString("$END ")
	b.WriteString(p.Name)
	b.WriteByte('\n')
	return nil
}

func writeText(b *strings.Builder, t *Text, compact bool, consts map[string]string) error {
	style := string(t.Style)
	if strings.TrimSpace(style) == "" {
		style = "{}"
	}
	coord := formatCoordQuad(t.X, t.Y, t.Width, t.TextSize)
	b.WriteString(">>")
	b.WriteString(coord)
	b.WriteString(pipeSep)
	b.WriteString(style)
	b.WriteString(pipeSep)
	b.WriteString(strconv.Itoa(t.Index))
	if strings.TrimSpace(t.ImageRef) != "" {
		b.WriteString(pipeSep)
		b.WriteString(t.ImageRef)
	}
	b.WriteByte('\n')
	lines := strings.Split(t.Content, "\n")
	for _, line := range lines {
		if compact {
			line = strings.ReplaceAll(line, "@", "")
			for constant := range consts {
				line = strings.ReplaceAll(line, "@"+constant+"@", consts[constant])
			}
			line = strings.ReplaceAll(line, "@", "")
		}
		b.WriteString(line)
		if !compact {
			b.WriteByte('\n')
		}
	}
	return nil
}

func writeMask(b *strings.Builder, m *Mask, consts map[string]string) error {
	_ = consts
	style := string(m.Style)
	if strings.TrimSpace(style) == "" {
		style = "{}"
	}
	coord := formatCoordQuad(m.X, m.Y, m.Width, m.Height)
	b.WriteString("==")
	b.WriteString(coord)
	b.WriteString(pipeSep)
	b.WriteString(style)
	b.WriteString(pipeSep)
	b.WriteString(strconv.Itoa(m.Index))
	if strings.TrimSpace(m.ImageRef) != "" {
		b.WriteString(pipeSep)
		b.WriteString(m.ImageRef)
	}
	b.WriteByte('\n')
	return nil
}

func writePathMask(b *strings.Builder, m *PathMask, consts map[string]string) error {
	_ = consts
	style := string(m.Style)
	if strings.TrimSpace(style) == "" {
		style = "{}"
	}
	coord := formatCoordQuad(m.X, m.Y, m.Width, m.Height)
	b.WriteString("~~")
	b.WriteString(coord)
	b.WriteString(pipeSep)
	b.WriteString(style)
	b.WriteString(pipeSep)
	b.WriteString(strconv.Itoa(m.Index))
	if strings.TrimSpace(m.ImageRef) != "" {
		b.WriteString(pipeSep)
		b.WriteString(m.ImageRef)
	}
	b.WriteByte('\n')
	for _, line := range svgpath.SplitCommands(m.Path) {
		b.WriteString(line)
		b.WriteByte('\n')
	}
	return nil
}

func formatCoordQuad(x, y, w, ts float64) string {
	return strings.Join([]string{
		formatCoord(x),
		formatCoord(y),
		formatCoord(w),
		formatCoord(ts),
	}, "-")
}

func sortedStringKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// FormatDocumentMarkdown renders the whole document as Markdown.
func FormatDocumentMarkdown(doc Document) string {
	var b strings.Builder
	b.WriteString("# PSRT document\n\n")
	for i := range doc.Pages {
		b.WriteString(formatPageMarkdown(&doc.Pages[i]))
		b.WriteByte('\n')
	}
	if len(doc.Fonts) > 0 {
		b.WriteString("## Fonts\n\n")
		for _, f := range doc.Fonts {
			b.WriteString("- ")
			b.WriteString(f)
			b.WriteByte('\n')
		}
		b.WriteByte('\n')
	}
	if len(doc.Consts) > 0 {
		b.WriteString("## Constants\n\n")
		for _, k := range sortedStringKeys(doc.Consts) {
			b.WriteString(FormatConstMarkdown(k, doc.Consts[k]))
			b.WriteByte('\n')
		}
	}
	return b.String()
}

// FormatPageMarkdown renders one page as Markdown.
func FormatPageMarkdown(p *Page) string {
	return formatPageMarkdown(p)
}

func formatPageMarkdown(p *Page) string {
	var b strings.Builder
	b.WriteString("## Page `")
	b.WriteString(p.Name)
	b.WriteString("`\n\n")
	b.WriteString("- **Image:** ")
	b.WriteString(p.ImageURL)
	b.WriteString("\n- **Page style (JSON):**\n\n```json\n")
	b.WriteString(prettyJSONOrRaw([]byte(p.Style)))
	b.WriteString("\n```\n\n### Texts\n\n")
	for i := range p.Texts {
		b.WriteString(formatTextMarkdown(&p.Texts[i]))
		b.WriteByte('\n')
	}
	return b.String()
}

// FormatTextMarkdown renders a single text block as Markdown.
func FormatTextMarkdown(t *Text) string {
	return formatTextMarkdown(t)
}

func formatTextMarkdown(t *Text) string {
	var b strings.Builder
	b.WriteString("#### Text index ")
	b.WriteString(strconv.Itoa(t.Index))
	b.WriteString("\n\n")
	b.WriteString("- **Position:** x=")
	b.WriteString(strconv.FormatFloat(t.X, 'f', -1, 64))
	b.WriteString(", y=")
	b.WriteString(strconv.FormatFloat(t.Y, 'f', -1, 64))
	b.WriteString(", width=")
	b.WriteString(strconv.FormatFloat(t.Width, 'f', -1, 64))
	b.WriteString(", textSize=")
	b.WriteString(strconv.FormatFloat(t.TextSize, 'f', -1, 64))
	b.WriteString("\n- **Image ref:** ")
	if strings.TrimSpace(t.ImageRef) == "" {
		b.WriteString("—")
	} else {
		b.WriteString(t.ImageRef)
	}
	b.WriteString("\n- **Style (JSON):**\n\n```json\n")
	b.WriteString(prettyJSONOrRaw([]byte(t.Style)))
	b.WriteString("\n```\n\n**Content:**\n\n")
	b.WriteString(t.Content)
	b.WriteString("\n")
	return b.String()
}

// FormatConstMarkdown renders a named constant as Markdown.
func FormatConstMarkdown(name, value string) string {
	var b strings.Builder
	b.WriteString("### `")
	b.WriteString(name)
	b.WriteString("`\n\n```\n")
	b.WriteString(value)
	b.WriteString("\n```\n")
	return b.String()
}

func prettyJSONOrRaw(raw []byte) string {
	if len(bytes.TrimSpace(raw)) == 0 {
		return "{}"
	}
	var buf bytes.Buffer
	if err := json.Indent(&buf, raw, "", "  "); err != nil {
		return string(raw)
	}
	return buf.String()
}

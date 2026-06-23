package psrt

import (
	"html"
	"strings"
)

type inlineDelim struct {
	open, close, tagOpen, tagClose string
}

var inlineDelims = []inlineDelim{
	{open: "***", close: "***", tagOpen: "<strong><em>", tagClose: "</em></strong>"},
	{open: "**", close: "**", tagOpen: "<strong>", tagClose: "</strong>"},
	{open: "*", close: "*", tagOpen: "<em>", tagClose: "</em>"},
	{open: "_", close: "_", tagOpen: "<u>", tagClose: "</u>"},
	{open: "~", close: "~", tagOpen: "<s>", tagClose: "</s>"},
}

// RenderInlineHTML converts PSRT inline markup to safe HTML. Newlines become <br/>.
func RenderInlineHTML(content string) string {
	if content == "" {
		return ""
	}
	lines := strings.Split(content, "\n")
	for i := range lines {
		lines[i] = renderInlineLine(lines[i])
	}
	return strings.Join(lines, "<br/>")
}

// PlainTextForLayout strips inline markup delimiters for line-width estimation.
func PlainTextForLayout(content string) string {
	if content == "" {
		return ""
	}
	lines := strings.Split(content, "\n")
	for i := range lines {
		lines[i] = plainInlineLine(lines[i])
	}
	return strings.Join(lines, "\n")
}

func renderInlineLine(line string) string {
	var b strings.Builder
	renderInlineSegment(line, &b)
	return b.String()
}

func plainInlineLine(line string) string {
	var b strings.Builder
	plainInlineSegment(line, &b)
	return b.String()
}

func renderInlineSegment(s string, b *strings.Builder) {
	i := 0
	for i < len(s) {
		if s[i] == '\\' && i+1 < len(s) {
			b.WriteString(html.EscapeString(s[i+1 : i+2]))
			i += 2
			continue
		}
		matched := false
		for _, d := range inlineDelims {
			if !strings.HasPrefix(s[i:], d.open) {
				continue
			}
			innerStart := i + len(d.open)
			closeAt := strings.Index(s[innerStart:], d.close)
			if closeAt <= 0 {
				continue
			}
			innerEnd := innerStart + closeAt
			b.WriteString(d.tagOpen)
			renderInlineSegment(s[innerStart:innerEnd], b)
			b.WriteString(d.tagClose)
			i = innerEnd + len(d.close)
			matched = true
			break
		}
		if matched {
			continue
		}
		next := i + 1
		if next > len(s) {
			next = len(s)
		}
		b.WriteString(html.EscapeString(s[i:next]))
		i = next
	}
}

func plainInlineSegment(s string, b *strings.Builder) {
	i := 0
	for i < len(s) {
		if s[i] == '\\' && i+1 < len(s) {
			b.WriteByte(s[i+1])
			i += 2
			continue
		}
		matched := false
		for _, d := range inlineDelims {
			if !strings.HasPrefix(s[i:], d.open) {
				continue
			}
			innerStart := i + len(d.open)
			closeAt := strings.Index(s[innerStart:], d.close)
			if closeAt <= 0 {
				continue
			}
			innerEnd := innerStart + closeAt
			plainInlineSegment(s[innerStart:innerEnd], b)
			i = innerEnd + len(d.close)
			matched = true
			break
		}
		if matched {
			continue
		}
		b.WriteByte(s[i])
		i++
	}
}

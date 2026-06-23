package psrt

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

const pipeSep = " | "

// maxPSRTLineBytes is the per-line limit for bufio.Scanner ($SOURCE base64 can exceed 64 KiB).
const maxPSRTLineBytes = 64 << 20 // 64 MiB

func newLineScanner(r io.Reader) *bufio.Scanner {
	sc := bufio.NewScanner(r)
	buf := make([]byte, 0, 64*1024)
	sc.Buffer(buf, maxPSRTLineBytes)
	sc.Split(scanLinesKeepTerminatorNone)
	return sc
}

type parseOptions struct {
	skipSourceValues bool
}

// Parse reads a PSRT document line by line and returns its structured form.
func Parse(r io.Reader) (Document, error) {
	return parseDocument(r, parseOptions{})
}

func parseDocument(r io.Reader, opts parseOptions) (Document, error) {
	var doc Document
	doc.Consts = make(map[string]string)
	doc.Sources = make(map[string]string)

	sc := newLineScanner(r)

	var cur *Page
	var active *textBuilder
	var skipMaskBody bool
	var inFonts, inConsts, inSource bool
	var sourceClosed bool
	lineNo := 0

	for sc.Scan() {
		lineNo++
		raw := sc.Bytes()
		line := string(raw)

		switch {
		case inFonts:
			if isEndFontsLine(line) {
				inFonts = false
				continue
			}
			s := strings.TrimSpace(line)
			if s == "" {
				continue
			}
			doc.Fonts = append(doc.Fonts, s)
			continue

		case inConsts:
			if isEndConstsLine(line) {
				inConsts = false
				continue
			}
			s := strings.TrimSpace(line)
			if s == "" {
				continue
			}
			if err := parseConstLine(s, doc.Consts, lineNo); err != nil {
				return doc, fmt.Errorf("line %d: %w", lineNo, err)
			}
			continue

		case inSource:
			if isEndSourceLine(line) {
				inSource = false
				sourceClosed = true
				continue
			}
			s := strings.TrimSpace(line)
			if s == "" {
				continue
			}
			if err := parseSourceLine(s, doc.Sources, lineNo, opts.skipSourceValues); err != nil {
				return doc, fmt.Errorf("line %d: %w", lineNo, err)
			}
			continue
		}

		if sourceClosed {
			s := strings.TrimSpace(line)
			if s != "" {
				return doc, fmt.Errorf("line %d: content after $ENDSOURCE; $SOURCE must be the last block", lineNo)
			}
			continue
		}

		if d, ok := directive(line); ok {
			switch d.kind {
			case dirFonts:
				if cur != nil {
					return doc, fmt.Errorf("line %d: $FONTS inside page %q", lineNo, cur.Name)
				}
				if inSource || sourceClosed {
					return doc, fmt.Errorf("line %d: $FONTS after $SOURCE; $SOURCE must be the last block", lineNo)
				}
				inFonts = true
				continue

			case dirConsts:
				if cur != nil {
					return doc, fmt.Errorf("line %d: $CONSTS inside page %q", lineNo, cur.Name)
				}
				if inSource || sourceClosed {
					return doc, fmt.Errorf("line %d: $CONSTS after $SOURCE; $SOURCE must be the last block", lineNo)
				}
				inConsts = true
				continue

			case dirSource:
				if cur != nil {
					return doc, fmt.Errorf("line %d: $SOURCE inside page %q", lineNo, cur.Name)
				}
				if inFonts || inConsts {
					return doc, fmt.Errorf("line %d: $SOURCE before closing $FONTS/$CONSTS", lineNo)
				}
				if inSource || sourceClosed {
					return doc, fmt.Errorf("line %d: duplicate $SOURCE block", lineNo)
				}
				inSource = true
				continue

			case dirStart:
				if cur != nil {
					return doc, fmt.Errorf("line %d: $START while page %q is still open", lineNo, cur.Name)
				}
				p, err := parsePageStart(d.rest, lineNo)
				if err != nil {
					return doc, err
				}
				doc.Pages = append(doc.Pages, p)
				cur = &doc.Pages[len(doc.Pages)-1]
				continue

			case dirEnd:
				if cur == nil {
					return doc, fmt.Errorf("line %d: $END without open page", lineNo)
				}
				name := strings.TrimSpace(d.rest)
				if name == "" {
					return doc, fmt.Errorf("line %d: $END missing page name", lineNo)
				}
				if name != cur.Name {
					return doc, fmt.Errorf("line %d: $END %q does not match open page %q", lineNo, name, cur.Name)
				}
				if err := flushTextBlock(cur, &active, lineNo); err != nil {
					return doc, err
				}
				cur = nil
				continue
			}
		}

		if cur == nil {
			s := strings.TrimSpace(line)
			if s == "" {
				continue
			}
			return doc, fmt.Errorf("line %d: unexpected content outside a page: %q", lineNo, trimForErr(s))
		}

		trimmedForHeader := strings.TrimLeft(line, " \t")
		if strings.HasPrefix(trimmedForHeader, ">>") {
			if err := flushTextBlock(cur, &active, lineNo); err != nil {
				return doc, err
			}
			skipMaskBody = false
			t, err := parseTextHeader(strings.TrimSpace(trimmedForHeader), lineNo)
			if err != nil {
				return doc, err
			}
			active = &textBuilder{text: t}
			continue
		}

		if strings.HasPrefix(trimmedForHeader, "==") {
			if err := flushTextBlock(cur, &active, lineNo); err != nil {
				return doc, err
			}
			m, err := parseMaskHeader(strings.TrimSpace(trimmedForHeader), lineNo)
			if err != nil {
				return doc, err
			}
			cur.Masks = append(cur.Masks, m)
			skipMaskBody = true
			continue
		}

		if skipMaskBody {
			continue
		}

		if active == nil {
			s := strings.TrimSpace(line)
			if s == "" {
				continue
			}
			return doc, fmt.Errorf("line %d: text content without active >> block in page %q", lineNo, cur.Name)
		}
		if active.buf.Len() > 0 {
			active.buf.WriteByte('\n')
		}
		active.buf.WriteString(strings.TrimRight(line, "\r"))
	}

	if err := sc.Err(); err != nil {
		return doc, err
	}
	if cur != nil {
		return doc, fmt.Errorf("line %d: EOF with open page %q", lineNo, cur.Name)
	}
	if inFonts {
		return doc, errors.New("EOF inside $FONTS")
	}
	if inConsts {
		return doc, errors.New("EOF inside $CONSTS")
	}
	if inSource {
		return doc, errors.New("EOF inside $SOURCE")
	}
	return doc, nil
}

// ParseString parses a PSRT document from a string.
func ParseString(input string) (Document, error) {
	return Parse(strings.NewReader(input))
}

// ToJSON returns the canonical JSON encoding of doc.
func ToJSON(doc Document) ([]byte, error) {
	return json.MarshalIndent(doc, "", "  ")
}

type textBuilder struct {
	text Text
	buf  strings.Builder
}

func flushTextBlock(cur *Page, active **textBuilder, lineNo int) error {
	if *active == nil {
		return nil
	}
	b := *active
	b.text.Content = NormalizeTextContent(b.buf.String())
	cur.Texts = append(cur.Texts, b.text)
	*active = nil
	return nil
}

func parsePageStart(rest string, lineNo int) (Page, error) {
	parts := strings.Split(rest, pipeSep)
	if len(parts) < 3 {
		return Page{}, fmt.Errorf("line %d: $START expects name | style | image-url (pipe+space separated)", lineNo)
	}
	name := strings.TrimSpace(parts[0])
	styleStr := strings.TrimSpace(parts[1])
	imageURL := strings.TrimSpace(parts[2])
	if name == "" {
		return Page{}, fmt.Errorf("line %d: $START missing page name", lineNo)
	}
	if !json.Valid([]byte(styleStr)) {
		return Page{}, fmt.Errorf("line %d: page style is not valid JSON", lineNo)
	}
	return Page{
		Name:     name,
		Style:    Style(styleStr),
		ImageURL: imageURL,
		Texts:    []Text{},
		Masks:    []Mask{},
	}, nil
}

func parseTextHeader(line string, lineNo int) (Text, error) {
	// line is trimmed; starts with >>
	body := strings.TrimSpace(line[2:])
	parts := strings.Split(body, pipeSep)
	if len(parts) < 3 {
		return Text{}, fmt.Errorf("line %d: text header needs coords | style | index", lineNo)
	}
	coords := strings.TrimSpace(parts[0])
	styleStr := strings.TrimSpace(parts[1])
	idxStr := strings.TrimSpace(parts[2])
	var imageRef string
	if len(parts) >= 4 {
		imageRef = strings.TrimSpace(parts[3])
	}

	x, y, w, ts, err := parseCoords(coords, lineNo)
	if err != nil {
		return Text{}, err
	}
	if !json.Valid([]byte(styleStr)) {
		return Text{}, fmt.Errorf("line %d: text style is not valid JSON", lineNo)
	}
	idx, err := strconv.Atoi(idxStr)
	if err != nil {
		return Text{}, fmt.Errorf("line %d: invalid text index %q: %w", lineNo, idxStr, err)
	}
	return Text{
		BaseBlock: BaseBlock{
			X: x, Y: y, Width: w,
			Style: Style(styleStr), Index: idx, ImageRef: imageRef,
		},
		TextSize: ts,
	}, nil
}

func parseMaskHeader(line string, lineNo int) (Mask, error) {
	body := strings.TrimSpace(line[2:])
	parts := strings.Split(body, pipeSep)
	if len(parts) < 3 {
		return Mask{}, fmt.Errorf("line %d: mask header needs coords | style | index", lineNo)
	}
	coords := strings.TrimSpace(parts[0])
	styleStr := strings.TrimSpace(parts[1])
	idxStr := strings.TrimSpace(parts[2])
	var imageRef string
	if len(parts) >= 4 {
		imageRef = strings.TrimSpace(parts[3])
	}

	x, y, w, h, err := parseMaskCoords(coords, lineNo)
	if err != nil {
		return Mask{}, err
	}
	if !json.Valid([]byte(styleStr)) {
		return Mask{}, fmt.Errorf("line %d: mask style is not valid JSON", lineNo)
	}
	idx, err := strconv.Atoi(idxStr)
	if err != nil {
		return Mask{}, fmt.Errorf("line %d: invalid mask index %q: %w", lineNo, idxStr, err)
	}
	return Mask{
		BaseBlock: BaseBlock{
			X: x, Y: y, Width: w,
			Style: Style(styleStr), Index: idx, ImageRef: imageRef,
		},
		Height: h,
	}, nil
}

func parseCoords(s string, lineNo int) (x, y, w, ts float64, err error) {
	chunks := strings.Split(s, "-")
	if len(chunks) != 4 {
		return 0, 0, 0, 0, fmt.Errorf("line %d: coords want X-Y-Width-TextSize, got %q", lineNo, s)
	}
	vals := make([]float64, 4)
	for i, c := range chunks {
		vals[i], err = strconv.ParseFloat(strings.TrimSpace(c), 64)
		if err != nil {
			return 0, 0, 0, 0, fmt.Errorf("line %d: invalid coord segment %q: %w", lineNo, c, err)
		}
		vals[i] = RoundCoord(vals[i])
	}
	return vals[0], vals[1], vals[2], vals[3], nil
}

func parseMaskCoords(s string, lineNo int) (x, y, w, h float64, err error) {
	chunks := strings.Split(s, "-")
	if len(chunks) != 4 {
		return 0, 0, 0, 0, fmt.Errorf("line %d: coords want X-Y-Width-Height, got %q", lineNo, s)
	}
	vals := make([]float64, 4)
	for i, c := range chunks {
		vals[i], err = strconv.ParseFloat(strings.TrimSpace(c), 64)
		if err != nil {
			return 0, 0, 0, 0, fmt.Errorf("line %d: invalid coord segment %q: %w", lineNo, c, err)
		}
		vals[i] = RoundCoord(vals[i])
	}
	return vals[0], vals[1], vals[2], vals[3], nil
}

func parseSourceLine(line string, dst map[string]string, lineNo int, skipValue bool) error {
	idx := strings.Index(line, pipeSep)
	if idx < 0 {
		return fmt.Errorf("source line must use %q between url and data-uri: %q", pipeSep, line)
	}
	key := strings.TrimSpace(line[:idx])
	val := strings.TrimSpace(line[idx+len(pipeSep):])
	if key == "" {
		return fmt.Errorf("empty source url")
	}
	if !skipValue && val == "" {
		return fmt.Errorf("empty source data-uri for %q", key)
	}
	if _, exists := dst[key]; exists {
		return fmt.Errorf("duplicate source %q", key)
	}
	if skipValue {
		dst[key] = ""
	} else {
		dst[key] = val
	}
	return nil
}

func parseConstLine(line string, dst map[string]string, lineNo int) error {
	s := strings.TrimSpace(line)
	s = strings.TrimPrefix(s, "@")
	s = strings.TrimSpace(s)
	idx := strings.Index(s, pipeSep)
	if idx < 0 {
		return fmt.Errorf("const line must use %q between name and value: %q", pipeSep, line)
	}
	key := strings.TrimSpace(s[:idx])
	val := strings.TrimSpace(s[idx+len(pipeSep):])
	if key == "" {
		return fmt.Errorf("empty const name")
	}
	if _, exists := dst[key]; exists {
		return fmt.Errorf("duplicate const %q", key)
	}
	dst[key] = val
	return nil
}

type dirKind int

const (
	dirNone dirKind = iota
	dirStart
	dirEnd
	dirFonts
	dirConsts
	dirSource
)

type directiveResult struct {
	kind dirKind
	rest string
}

func directive(line string) (directiveResult, bool) {
	s := strings.TrimSpace(line)
	if !strings.HasPrefix(s, "$") {
		return directiveResult{}, false
	}
	after := strings.TrimSpace(s[1:])
	afterUpper := strings.ToUpper(after)

	if strings.HasPrefix(afterUpper, "ENDCONSTS") {
		return directiveResult{}, false
	}
	if strings.HasPrefix(afterUpper, "ENDFONTS") {
		return directiveResult{}, false
	}
	if strings.HasPrefix(afterUpper, "ENDSOURCE") {
		return directiveResult{}, false
	}

	const (
		kwSTART  = "START"
		kwEND    = "END"
		kwFONTS  = "FONTS"
		kwCONSTS = "CONSTS"
		kwSOURCE = "SOURCE"
	)

	if hasKeyword(after, afterUpper, kwSTART) {
		rest := strings.TrimSpace(after[len(kwSTART):])
		return directiveResult{kind: dirStart, rest: rest}, true
	}
	if hasKeyword(after, afterUpper, kwEND) {
		rest := strings.TrimSpace(after[len(kwEND):])
		return directiveResult{kind: dirEnd, rest: rest}, true
	}
	if hasKeyword(after, afterUpper, kwFONTS) {
		return directiveResult{kind: dirFonts}, true
	}
	if hasKeyword(after, afterUpper, kwCONSTS) {
		return directiveResult{kind: dirConsts}, true
	}
	if hasKeyword(after, afterUpper, kwSOURCE) {
		return directiveResult{kind: dirSource}, true
	}
	return directiveResult{}, false
}

// hasKeyword reports whether afterUpper starts with kw (ASCII, upper) and the keyword
// is followed by end-of-string or a space (so $STARTPAGE is not $START).
func hasKeyword(after, afterUpper, kw string) bool {
	if !strings.HasPrefix(afterUpper, kw) {
		return false
	}
	if len(after) == len(kw) {
		return true
	}
	if len(after) <= len(kw) {
		return false
	}
	r, _ := utf8.DecodeRuneInString(after[len(kw):])
	return r == utf8.RuneError || unicode.IsSpace(r)
}

func isEndFontsLine(line string) bool {
	s := strings.TrimSpace(line)
	if !strings.HasPrefix(s, "$") {
		return false
	}
	after := strings.TrimSpace(strings.TrimPrefix(s, "$"))
	return strings.HasPrefix(strings.ToUpper(after), "ENDFONTS")
}

func isEndConstsLine(line string) bool {
	s := strings.TrimSpace(line)
	if !strings.HasPrefix(s, "$") {
		return false
	}
	after := strings.TrimSpace(strings.TrimPrefix(s, "$"))
	return strings.HasPrefix(strings.ToUpper(after), "ENDCONSTS")
}

func isEndSourceLine(line string) bool {
	s := strings.TrimSpace(line)
	if !strings.HasPrefix(s, "$") {
		return false
	}
	after := strings.TrimSpace(strings.TrimPrefix(s, "$"))
	return strings.HasPrefix(strings.ToUpper(after), "ENDSOURCE")
}

func trimForErr(s string) string {
	const max = 80
	if len(s) <= max {
		return s
	}
	return s[:max] + "…"
}

// scanLinesKeepTerminatorNone splits on \n and \r\n without including newline in token.
func scanLinesKeepTerminatorNone(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := bytes.IndexByte(data, '\n'); i >= 0 {
		return i + 1, dropCR(data[0:i]), nil
	}
	if atEOF {
		return len(data), dropCR(data), nil
	}
	return 0, nil, nil
}

func dropCR(data []byte) []byte {
	if len(data) > 0 && data[len(data)-1] == '\r' {
		return data[0 : len(data)-1]
	}
	return data
}

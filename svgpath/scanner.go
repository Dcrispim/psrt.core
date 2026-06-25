package svgpath

import "fmt"

type scanner struct {
	s   string
	pos int
}

func (p *scanner) eof() bool { return p.pos >= len(p.s) }

func (p *scanner) skipWSP() {
	for !p.eof() {
		switch p.s[p.pos] {
		case ' ', '\t', '\n', '\r', '\f':
			p.pos++
		default:
			return
		}
	}
}

// skipSep skips an optional comma-wsp separator (wsp* ","? wsp*). It is a
// no-op when no separator is present, which is what allows the compact
// no-separator notation between arc flags and the following number.
func (p *scanner) skipSep() {
	p.skipWSP()
	if !p.eof() && p.s[p.pos] == ',' {
		p.pos++
		p.skipWSP()
	}
}

func isCommandLetter(b byte) bool {
	switch b {
	case 'M', 'm', 'L', 'l', 'H', 'h', 'V', 'v', 'C', 'c', 'S', 's', 'Q', 'q', 'T', 't', 'A', 'a', 'Z', 'z':
		return true
	}
	return false
}

func (p *scanner) readCommandLetter() (byte, bool) {
	if p.eof() || !isCommandLetter(p.s[p.pos]) {
		return 0, false
	}
	c := p.s[p.pos]
	p.pos++
	return c, true
}

func isDigit(b byte) bool { return b >= '0' && b <= '9' }

// readNumber consumes one SVG "number" token (sign? digits? ("." digits?)? exponent?)
// at the current position, requiring at least one digit overall.
func (p *scanner) readNumber() (string, bool) {
	start := p.pos
	if !p.eof() && (p.s[p.pos] == '+' || p.s[p.pos] == '-') {
		p.pos++
	}
	digitsBefore := 0
	for !p.eof() && isDigit(p.s[p.pos]) {
		p.pos++
		digitsBefore++
	}
	digitsAfter := 0
	if !p.eof() && p.s[p.pos] == '.' {
		p.pos++
		for !p.eof() && isDigit(p.s[p.pos]) {
			p.pos++
			digitsAfter++
		}
	}
	if digitsBefore == 0 && digitsAfter == 0 {
		p.pos = start
		return "", false
	}
	if !p.eof() && (p.s[p.pos] == 'e' || p.s[p.pos] == 'E') {
		ePos := p.pos
		p.pos++
		if !p.eof() && (p.s[p.pos] == '+' || p.s[p.pos] == '-') {
			p.pos++
		}
		expDigits := 0
		for !p.eof() && isDigit(p.s[p.pos]) {
			p.pos++
			expDigits++
		}
		if expDigits == 0 {
			p.pos = ePos
		}
	}
	return p.s[start:p.pos], true
}

// readFlag consumes exactly one '0' or '1' character — used for the arc
// command's large-arc-flag/sweep-flag, which are always a single digit even
// when glued to the following number with no separator.
func (p *scanner) readFlag() (byte, bool) {
	if p.eof() {
		return 0, false
	}
	c := p.s[p.pos]
	if c != '0' && c != '1' {
		return 0, false
	}
	p.pos++
	return c, true
}

// readNumbers consumes n comma-wsp-separated numbers (one argument group).
func (p *scanner) readNumbers(n int) error {
	for i := 0; i < n; i++ {
		if i == 0 {
			p.skipWSP()
		} else {
			p.skipSep()
		}
		if _, ok := p.readNumber(); !ok {
			return fmt.Errorf("expected number at offset %d", p.pos)
		}
	}
	return nil
}

// readArcArgs consumes one elliptical-arc argument group:
// rx ry x-axis-rotation large-arc-flag sweep-flag x y.
func (p *scanner) readArcArgs() error {
	p.skipWSP()
	if _, ok := p.readNumber(); !ok {
		return fmt.Errorf("expected arc radius rx at offset %d", p.pos)
	}
	p.skipSep()
	if _, ok := p.readNumber(); !ok {
		return fmt.Errorf("expected arc radius ry at offset %d", p.pos)
	}
	p.skipSep()
	if _, ok := p.readNumber(); !ok {
		return fmt.Errorf("expected arc x-axis-rotation at offset %d", p.pos)
	}
	p.skipSep()
	if _, ok := p.readFlag(); !ok {
		return fmt.Errorf("expected arc large-arc-flag (0 or 1) at offset %d", p.pos)
	}
	p.skipSep()
	if _, ok := p.readFlag(); !ok {
		return fmt.Errorf("expected arc sweep-flag (0 or 1) at offset %d", p.pos)
	}
	p.skipSep()
	if _, ok := p.readNumber(); !ok {
		return fmt.Errorf("expected arc endpoint x at offset %d", p.pos)
	}
	p.skipSep()
	if _, ok := p.readNumber(); !ok {
		return fmt.Errorf("expected arc endpoint y at offset %d", p.pos)
	}
	return nil
}

// readArgs consumes exactly one argument group for cmd (case-insensitive —
// relative vs. absolute only affects semantics, not the parameter count).
func (p *scanner) readArgs(cmd byte) error {
	switch cmd {
	case 'H', 'h', 'V', 'v':
		return p.readNumbers(1)
	case 'M', 'm', 'L', 'l', 'T', 't':
		return p.readNumbers(2)
	case 'S', 's', 'Q', 'q':
		return p.readNumbers(4)
	case 'C', 'c':
		return p.readNumbers(6)
	case 'A', 'a':
		return p.readArcArgs()
	case 'Z', 'z':
		return nil
	}
	return fmt.Errorf("unsupported command %q", string(cmd))
}

// Package svgpath parses and validates SVG path `d` attribute data — the
// grammar shared by <path>, <clipPath>, and CSS clip-path: path(). It has no
// knowledge of PSRT; callers compose it with their own policies (e.g. PSRT's
// ~~ block requiring a single shape — see psrt.flushPathMaskBlock).
//
// This is a purpose-built parser rather than a third-party SVG library:
// available Go libraries either bundle unrelated rasterization machinery
// (oksvg embeds rasterx, which pulls in fill/stroke/scan code no caller here
// needs) or have unproven handling of the compact elliptical-arc flag
// notation (e.g. "A5,5,45,0030,30", where the flags are single digits glued
// to the following number) — handled explicitly here via readFlag.
package svgpath

import "fmt"

// Info describes basic structural facts about a parsed path.
type Info struct {
	// Subpaths is the number of moveto (M/m) commands in the path. More than
	// one means the path describes multiple disconnected shapes.
	Subpaths int
}

// Parse validates d against the SVG path data grammar
// (https://www.w3.org/TR/SVG11/paths.html#PathData) and returns basic
// structural info about it.
func Parse(d string) (Info, error) {
	s := &scanner{s: d}
	s.skipWSP()
	if s.eof() {
		return Info{}, fmt.Errorf("empty path")
	}

	cmd, ok := s.readCommandLetter()
	if !ok || (cmd != 'M' && cmd != 'm') {
		return Info{}, fmt.Errorf("path must start with a moveto command (M or m)")
	}
	info := Info{Subpaths: 1}
	if err := s.readArgs(cmd); err != nil {
		return Info{}, err
	}

	cur := cmd
	for {
		s.skipWSP()
		if s.eof() {
			break
		}
		if nc, ok := s.readCommandLetter(); ok {
			cur = nc
			if cur == 'Z' || cur == 'z' {
				continue
			}
			if cur == 'M' || cur == 'm' {
				info.Subpaths++
			}
			if err := s.readArgs(cur); err != nil {
				return Info{}, err
			}
			continue
		}
		if cur == 'Z' || cur == 'z' {
			return Info{}, fmt.Errorf("unexpected content after Z at offset %d", s.pos)
		}
		if err := s.readArgs(cur); err != nil {
			return Info{}, err
		}
	}
	return info, nil
}

// Validate checks that d is syntactically valid SVG path data.
func Validate(d string) error {
	_, err := Parse(d)
	return err
}

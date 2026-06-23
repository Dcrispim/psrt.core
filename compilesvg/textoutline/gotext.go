package textoutline

import (
	"bytes"
	"fmt"
	"math"
	"strings"
	"unicode"

	"github.com/go-text/typesetting/di"
	"github.com/go-text/typesetting/font"
	ot "github.com/go-text/typesetting/font/opentype"
	"github.com/go-text/typesetting/language"
	"github.com/go-text/typesetting/shaping"
	"golang.org/x/image/math/fixed"
)

func outlineGoText(in PageInput) ([]OutlinedBlock, error) {
	faces, err := loadFaces(in.Fonts)
	if err != nil {
		return nil, err
	}
	defer faces.close()

	var shaper shaping.HarfbuzzShaper
	out := make([]OutlinedBlock, 0, len(in.Blocks))
	for _, blk := range in.Blocks {
		ob, err := outlineBlockGoText(&shaper, faces, blk)
		if err != nil {
			return nil, fmt.Errorf("text block %d: %w", blk.Index, err)
		}
		out = append(out, ob)
	}
	return out, nil
}

type faceSet struct {
	byFamily map[string]*font.Face
	defaultF *font.Face
}

func (fs *faceSet) close() {
	seen := make(map[*font.Face]bool)
	for _, f := range fs.byFamily {
		if f != nil && !seen[f] {
			seen[f] = true
		}
	}
}

func (fs *faceSet) pick(family string) *font.Face {
	family = strings.Trim(family, `"' `)
	if family != "" {
		if f, ok := fs.byFamily[family]; ok && f != nil {
			return f
		}
	}
	return fs.defaultF
}

func loadFaces(fonts map[string]FontBytes) (*faceSet, error) {
	fs := &faceSet{byFamily: make(map[string]*font.Face)}
	for name, fb := range fonts {
		face, err := font.ParseTTF(bytes.NewReader(fb.Bytes))
		if err != nil {
			continue
		}
		fs.byFamily[name] = face
		if fs.defaultF == nil {
			fs.defaultF = face
		}
	}
	if fs.defaultF == nil {
		face, err := font.ParseTTF(bytes.NewReader(DefaultFontBytes()))
		if err != nil {
			return nil, fmt.Errorf("parse default font: %w", err)
		}
		fs.byFamily[DefaultFontFamily] = face
		fs.defaultF = face
	}
	return fs, nil
}

func outlineBlockGoText(shaper *shaping.HarfbuzzShaper, faces *faceSet, blk BlockInput) (OutlinedBlock, error) {
	st := blk.Style
	face := faces.pick(st.FontFamily)
	if face == nil {
		return OutlinedBlock{Index: blk.Index, PlainText: blk.PlainText}, nil
	}

	fontSize := st.FontSizePx
	if fontSize <= 0 {
		fontSize = 12
	}
	lh := st.LineHeight
	if lh <= 0 {
		lh = 1.2
	}
	linePx := fontSize * lh

	lines := wrapPlainLines(shaper, face, fontSize, st.ContentW, blk.PlainText)
	if len(lines) == 0 {
		return OutlinedBlock{Index: blk.Index, PlainText: blk.PlainText}, nil
	}

	var paths []OutlinedPath
	baseX := float64(blk.X + st.PadLeft)
	baseY := float64(blk.Y + st.PadTop)

	for li, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		lineW := measureLineWidth(shaper, face, fontSize, line)
		x := baseX
		switch strings.ToLower(st.TextAlign) {
		case "center":
			x += (float64(st.ContentW) - lineW) / 2
		case "right", "end":
			x += float64(st.ContentW) - lineW
		}
		baselineY := baseY + fontSize*0.8 + float64(li)*linePx
		linePaths := shapeLineToPaths(shaper, face, fontSize, line, x, baselineY, st)
		paths = append(paths, linePaths...)
	}

	return OutlinedBlock{
		Index:     blk.Index,
		Paths:     paths,
		PlainText: blk.PlainText,
		FilterID:  blk.FilterID,
		Transform: blk.Transform,
	}, nil
}

func wrapPlainLines(shaper *shaping.HarfbuzzShaper, face *font.Face, fontSize float64, maxW int, text string) []string {
	if maxW < 1 {
		maxW = 1
	}
	paras := strings.Split(text, "\n")
	var out []string
	for _, para := range paras {
		para = strings.TrimRight(para, "\r")
		if para == "" {
			out = append(out, "")
			continue
		}
		words := splitWords(para)
		var line strings.Builder
		for _, w := range words {
			candidate := line.String()
			if candidate != "" {
				candidate += " "
			}
			candidate += w
			if line.Len() > 0 && measureLineWidth(shaper, face, fontSize, candidate) > float64(maxW) {
				out = append(out, strings.TrimSpace(line.String()))
				line.Reset()
				line.WriteString(w)
			} else {
				if line.Len() > 0 {
					line.WriteByte(' ')
				}
				line.WriteString(w)
			}
		}
		if line.Len() > 0 {
			out = append(out, strings.TrimSpace(line.String()))
		}
	}
	if len(out) == 0 {
		out = []string{""}
	}
	return out
}

func splitWords(s string) []string {
	var words []string
	var cur strings.Builder
	flush := func() {
		if cur.Len() > 0 {
			words = append(words, cur.String())
			cur.Reset()
		}
	}
	for _, r := range s {
		if unicode.IsSpace(r) {
			flush()
			continue
		}
		cur.WriteRune(r)
	}
	flush()
	if len(words) == 0 {
		return []string{s}
	}
	return words
}

func shapeInput(face *font.Face, fontSize float64, text []rune) shaping.Input {
	return shaping.Input{
		Text:      text,
		RunStart:  0,
		RunEnd:    len(text),
		Face:      face,
		Direction: di.DirectionLTR,
		Size:      fontSizeFixed(fontSize),
		Script:    language.Latin,
		Language:  language.NewLanguage("en"),
	}
}

func fontSizeFixed(fontSize float64) fixed.Int26_6 {
	return fixed.Int26_6(math.Round(fontSize * 64))
}

func measureLineWidth(shaper *shaping.HarfbuzzShaper, face *font.Face, fontSize float64, line string) float64 {
	out := shaper.Shape(shapeInput(face, fontSize, []rune(line)))
	return fixedToFloat(out.Advance)
}

func shapeLineToPaths(
	shaper *shaping.HarfbuzzShaper,
	face *font.Face,
	fontSize float64,
	line string,
	startX, baselineY float64,
	st BlockStyle,
) []OutlinedPath {
	out := shaper.Shape(shapeInput(face, fontSize, []rune(line)))
	attrs := pathAttrsFromStyle(st)
	var paths []OutlinedPath
	penX := startX
	for _, g := range out.Glyphs {
		gx := penX + fixedToFloat(g.XOffset)
		gy := baselineY + fixedToFloat(g.YOffset)
		outline, ok := face.GlyphData(g.GlyphID).(font.GlyphOutline)
		if !ok {
			penX += fixedToFloat(g.Advance)
			continue
		}
		d := segmentsToSVG(outline.Segments, gx, gy, &out)
		if d == "" {
			penX += fixedToFloat(g.Advance)
			continue
		}
		p := OutlinedPath{D: d, Fill: attrs.Fill, Stroke: attrs.Stroke, StrokeWidth: attrs.StrokeWidth, PaintOrder: attrs.PaintOrder}
		paths = append(paths, p)
		penX += fixedToFloat(g.Advance)
	}
	return paths
}

func pathAttrsFromStyle(st BlockStyle) OutlinedPath {
	fill := st.Color
	if fill == "" {
		fill = "#000000"
	}
	p := OutlinedPath{Fill: fill}
	if st.Stroke != "" && st.StrokeWidth != "" && st.StrokeWidth != "0" && st.StrokeWidth != "0px" {
		p.Stroke = st.Stroke
		p.StrokeWidth = st.StrokeWidth
		p.PaintOrder = "stroke fill"
	}
	return p
}

func fixedToFloat(v fixed.Int26_6) float64 {
	return float64(v) / 64.0
}

func segmentsToSVG(segments []font.Segment, ox, oy float64, shaped *shaping.Output) string {
	var b strings.Builder
	var cur ot.SegmentPoint
	hasCur := false
	for _, seg := range segments {
		pts := seg.ArgsSlice()
		switch seg.Op {
		case ot.SegmentOpMoveTo:
			px := ox + fixedToFloat(shaped.FromFontUnit(pts[0].X))
			py := oy - fixedToFloat(shaped.FromFontUnit(pts[0].Y))
			fmt.Fprintf(&b, "M%.3f,%.3f", px, py)
			cur = pts[0]
			hasCur = true
		case ot.SegmentOpLineTo:
			if !hasCur {
				continue
			}
			px := ox + fixedToFloat(shaped.FromFontUnit(pts[0].X))
			py := oy - fixedToFloat(shaped.FromFontUnit(pts[0].Y))
			fmt.Fprintf(&b, "L%.3f,%.3f", px, py)
			cur = pts[0]
		case ot.SegmentOpQuadTo:
			if !hasCur {
				continue
			}
			cx := ox + fixedToFloat(shaped.FromFontUnit(pts[0].X))
			cy := oy - fixedToFloat(shaped.FromFontUnit(pts[0].Y))
			px := ox + fixedToFloat(shaped.FromFontUnit(pts[1].X))
			py := oy - fixedToFloat(shaped.FromFontUnit(pts[1].Y))
			fmt.Fprintf(&b, "Q%.3f,%.3f %.3f,%.3f", cx, cy, px, py)
			cur = pts[1]
		case ot.SegmentOpCubeTo:
			if !hasCur {
				continue
			}
			c1x := ox + fixedToFloat(shaped.FromFontUnit(pts[0].X))
			c1y := oy - fixedToFloat(shaped.FromFontUnit(pts[0].Y))
			c2x := ox + fixedToFloat(shaped.FromFontUnit(pts[1].X))
			c2y := oy - fixedToFloat(shaped.FromFontUnit(pts[1].Y))
			px := ox + fixedToFloat(shaped.FromFontUnit(pts[2].X))
			py := oy - fixedToFloat(shaped.FromFontUnit(pts[2].Y))
			fmt.Fprintf(&b, "C%.3f,%.3f %.3f,%.3f %.3f,%.3f", c1x, c1y, c2x, c2y, px, py)
			cur = pts[2]
		}
	}
	_ = cur
	return b.String()
}

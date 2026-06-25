package styleadapter

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Dcrispim/psrt.core/styleadapter/percent"
)

const blurFilterKind = "blur"

// BlurAdapt carries SVG ids produced for backdrop-style box blur.
type BlurAdapt struct {
	FilterID string
	MaskID   string
}

// BlurSpec is a resolved blur amount and optional single side.
type BlurSpec struct {
	AmountPx float64
	Side     string // "", "left", "right", "top", "bottom"
}

var blurSideKeys = []struct {
	key  string
	side string
}{
	{KeyBlurLeft, "left"},
	{KeyBlurRight, "right"},
	{KeyBlurTop, "top"},
	{KeyBlurBottom, "bottom"},
}

// expandBlur removes blur keys from style and returns SVG filter/mask fragments and/or an HTML box patch.
func expandBlur(style map[string]json.RawMessage, dims percent.ImageDims, html bool, filterID string) (map[string]json.RawMessage, BlurAdapt, []StyleFragment) {
	if len(style) == 0 {
		return style, BlurAdapt{}, nil
	}
	spec, ok := parseBlurFromStyle(style, dims)
	if !ok {
		return style, BlurAdapt{}, nil
	}
	out := copyStyleMap(style)
	for _, bk := range blurSideKeys {
		delete(out, bk.key)
	}
	delete(out, KeyBlur)

	var frags []StyleFragment
	meta := BlurAdapt{FilterID: filterID + "-blur"}
	if html {
		patch := NewFragment(TypeMotionDiv)
		applyBlurHTML(patch, spec)
		frags = append(frags, patch)
		return out, meta, frags
	}
	frags = append(frags, blurFilterFragment(meta.FilterID, spec))
	if spec.Side != "" {
		meta.MaskID = meta.FilterID + "-mask"
		frags = append(frags, blurMaskFragment(meta.MaskID, spec.Side))
	}
	return out, meta, frags
}

func parseBlurFromStyle(style map[string]json.RawMessage, dims percent.ImageDims) (BlurSpec, bool) {
	for _, sk := range blurSideKeys {
		if raw, ok := style[sk.key]; ok && HasStyleValue(sk.key, raw) {
			amount, ok := parseBlurAmount(StringifyCSSValue(raw), dims)
			if !ok || amount <= 0 {
				return BlurSpec{}, false
			}
			return BlurSpec{AmountPx: amount, Side: sk.side}, true
		}
	}
	if raw, ok := style[KeyBlur]; ok && HasStyleValue(KeyBlur, raw) {
		return parseBlurCSSValue(StringifyCSSValue(raw), dims)
	}
	return BlurSpec{}, false
}

func parseBlurCSSValue(val string, dims percent.ImageDims) (BlurSpec, bool) {
	val = strings.TrimSpace(val)
	if val == "" {
		return BlurSpec{}, false
	}
	parts := strings.Fields(val)
	if len(parts) == 0 {
		return BlurSpec{}, false
	}
	side := ""
	var amountParts []string
	for _, p := range parts {
		low := strings.ToLower(p)
		switch low {
		case "left", "right", "top", "bottom":
			if side == "" {
				side = low
			}
		default:
			amountParts = append(amountParts, p)
		}
	}
	amountStr := strings.Join(amountParts, " ")
	if amountStr == "" && side != "" && len(parts) >= 2 {
		// "left 8px" with side first — retry treating first token as side only
		for i, p := range parts {
			low := strings.ToLower(p)
			if low == side && i+1 < len(parts) {
				amountStr = strings.Join(parts[i+1:], " ")
				break
			}
		}
	}
	if amountStr == "" {
		amountStr = val
		side = ""
	}
	amount, ok := parseBlurAmount(amountStr, dims)
	if !ok || amount <= 0 {
		return BlurSpec{}, false
	}
	return BlurSpec{AmountPx: amount, Side: side}, true
}

func parseBlurAmount(s string, dims percent.ImageDims) (float64, bool) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, false
	}
	if strings.Contains(s, "%") {
		resolved := percent.ApplyPercentHandlers(
			map[string]json.RawMessage{KeyBlur: json.RawMessage(jsonQuote(s))},
			dims,
		)
		s = StringifyCSSValue(resolved[KeyBlur])
	}
	px := parsePxNum(s)
	return px, px > 0
}

func applyBlurHTML(box StyleFragment, spec BlurSpec) {
	if box == nil || spec.AmountPx <= 0 {
		return
	}
	px := fmt.Sprintf("%.3fpx", spec.AmountPx)
	blurVal := "blur(" + px + ")"
	box.Set("backdropFilter", blurVal)
	box.Set("WebkitBackdropFilter", blurVal)
	if spec.Side == "" {
		return
	}
	mask := blurMaskCSS(spec.Side)
	box.Set("maskImage", mask)
	box.Set("WebkitMaskImage", mask)
	box.Set("maskSize", "100% 100%")
	box.Set("WebkitMaskSize", "100% 100%")
}

func blurMaskCSS(side string) string {
	switch strings.ToLower(side) {
	case "left":
		return "linear-gradient(to right, rgba(0,0,0,1) 0%, rgba(0,0,0,0) 100%)"
	case "right":
		return "linear-gradient(to left, rgba(0,0,0,1) 0%, rgba(0,0,0,0) 100%)"
	case "top":
		return "linear-gradient(to bottom, rgba(0,0,0,1) 0%, rgba(0,0,0,0) 100%)"
	case "bottom":
		return "linear-gradient(to top, rgba(0,0,0,1) 0%, rgba(0,0,0,0) 100%)"
	default:
		return ""
	}
}

func blurFilterFragment(id string, spec BlurSpec) StyleFragment {
	f := NewFragment(TypeFilter)
	f.Set("id", id)
	f.Set("filterKind", blurFilterKind)
	std := spec.AmountPx / 2
	if std < 0.5 {
		std = 0.5
	}
	f.Set("feGaussianBlurStd", fmt.Sprintf("%.3f", std))
	f.Set("feGaussianBlurIn", "SourceGraphic")
	return f
}

func blurMaskFragment(id, side string) StyleFragment {
	f := NewFragment(TypeMask)
	f.Set("id", id)
	f.Set("maskSide", side)
	return f
}

func applyBlurSVGRect(rect StyleFragment, meta BlurAdapt) {
	if rect == nil || meta.FilterID == "" {
		return
	}
	rect.Set("filter", "url(#"+meta.FilterID+")")
	if meta.MaskID != "" {
		rect.Set("mask", "url(#"+meta.MaskID+")")
	}
}

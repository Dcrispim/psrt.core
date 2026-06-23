package percent

import "strings"

type paddingHandler struct{}

func (paddingHandler) Keys() []string {
	return []string{
		"padding", "paddingTop", "paddingRight", "paddingBottom", "paddingLeft",
	}
}

func (paddingHandler) Resolve(key string, value string, dims ImageDims) (string, bool) {
	if !strings.Contains(value, "%") {
		return value, true
	}
	tokens := strings.Fields(value)
	if len(tokens) == 0 {
		return value, true
	}
	// Single-value padding or side-specific key.
	if len(tokens) == 1 || key != "padding" {
		base := paddingAxisBase(key, len(tokens), 0, dims)
		if len(tokens) == 1 {
			return singlePercentToken(tokens[0], base, dims)
		}
	}
	// Shorthand 1-4 values.
	bases := shorthandBases(len(tokens))
	out := make([]string, len(tokens))
	for i, tok := range tokens {
		base := paddingAxisBase(key, len(tokens), bases[i], dims)
		resolved, _ := singlePercentToken(tok, base, dims)
		out[i] = resolved
	}
	return strings.Join(out, " "), true
}

func shorthandBases(n int) []int {
	switch n {
	case 1:
		return []int{0, 0, 0, 0}
	case 2:
		return []int{0, 1, 0, 1} // vertical, horizontal
	case 3:
		return []int{0, 1, 2, 1}
	default:
		return []int{0, 1, 2, 3}
	}
}

// axis: 0=top, 1=right, 2=bottom, 3=left for shorthand; side keys use W/H directly.
func paddingAxisBase(key string, count int, axis int, dims ImageDims) float64 {
	switch key {
	case "paddingTop", "paddingBottom":
		return float64(dims.H)
	case "paddingLeft", "paddingRight":
		return float64(dims.W)
	}
	// shorthand
	if count == 1 {
		return float64(max(dims.W, dims.H))
	}
	switch axis {
	case 0, 2:
		return float64(dims.H)
	case 1, 3:
		return float64(dims.W)
	default:
		return float64(dims.W)
	}
}

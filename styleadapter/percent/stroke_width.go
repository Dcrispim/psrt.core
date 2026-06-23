package percent

import "strings"

type strokeWidthHandler struct{}

func (strokeWidthHandler) Keys() []string { return []string{"strokeWidth"} }

func (strokeWidthHandler) Resolve(_ string, value string, dims ImageDims) (string, bool) {
	if !strings.Contains(value, "%") {
		return value, true
	}
	base := dims.FontSizePx
	if base <= 0 {
		base = 1
	}
	return singlePercentToken(value, base, dims)
}

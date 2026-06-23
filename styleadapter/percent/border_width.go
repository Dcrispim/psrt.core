package percent

import "strings"

type borderWidthHandler struct{}

func (borderWidthHandler) Keys() []string {
	return []string{
		"borderWidth", "borderTopWidth", "borderRightWidth",
		"borderBottomWidth", "borderLeftWidth",
	}
}

func (borderWidthHandler) Resolve(_ string, value string, dims ImageDims) (string, bool) {
	if !strings.Contains(value, "%") {
		return value, true
	}
	return singlePercentToken(value, float64(dims.W), dims)
}

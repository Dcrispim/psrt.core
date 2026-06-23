package percent

import "strings"

type lineHeightHandler struct{}

func (lineHeightHandler) Keys() []string { return []string{"lineHeight"} }

func (lineHeightHandler) Resolve(_ string, value string, dims ImageDims) (string, bool) {
	if !strings.HasSuffix(strings.TrimSpace(value), "%") {
		return value, true
	}
	return singlePercentToken(value, float64(dims.H), dims)
}

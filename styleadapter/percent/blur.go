package percent

import "strings"

type blurHandler struct{}

func (blurHandler) Keys() []string {
	return []string{"blur", "blurLeft", "blurRight", "blurTop", "blurBottom"}
}

func (blurHandler) Resolve(_ string, value string, dims ImageDims) (string, bool) {
	if !strings.Contains(value, "%") {
		return value, true
	}
	base := float64(dims.W)
	if dims.H > 0 && dims.H < dims.W {
		base = float64(dims.H)
	}
	return singlePercentToken(value, base, dims)
}

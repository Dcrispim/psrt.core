package percent

import "strings"

type dimensionHandler struct{}

func (dimensionHandler) Keys() []string { return []string{"height", "width"} }

func (dimensionHandler) Resolve(key string, value string, dims ImageDims) (string, bool) {
	if !strings.Contains(value, "%") {
		return value, true
	}
	base := float64(dims.H)
	if key == "width" {
		base = float64(dims.W)
	}
	return singlePercentToken(value, base, dims)
}

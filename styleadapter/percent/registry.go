package percent

import (
	"encoding/json"
	"strings"
)

// ImageDims holds canvas/image dimensions for percent resolution at render time.
type ImageDims struct {
	W, H       int
	FontSizePx float64
	Zoom       float64
	// PreservePercent keeps % values for HTML compile so padding scales with the slide overlay.
	PreservePercent bool
}

// PercentHandler resolves % values to px for specific properties.
type PercentHandler interface {
	Keys() []string
	Resolve(key string, value string, dims ImageDims) (resolved string, ok bool)
}

var defaultHandlers []PercentHandler

func init() {
	defaultHandlers = []PercentHandler{
		textShadowHandler{},
		boxShadowHandler{},
		blurHandler{},
		borderWidthHandler{},
		strokeWidthHandler{},
		lineHeightHandler{},
		paddingHandler{},
		dimensionHandler{},
	}
}

// RegisterPercentHandler adds a handler to the registry.
func RegisterPercentHandler(h PercentHandler) {
	defaultHandlers = append(defaultHandlers, h)
}

// ApplyPercentHandlers converts % to px for registered keys (render-time only).
func ApplyPercentHandlers(style map[string]json.RawMessage, dims ImageDims) map[string]json.RawMessage {
	if len(style) == 0 {
		return style
	}
	if dims.Zoom <= 0 {
		dims.Zoom = 1
	}
	keyHandler := make(map[string]PercentHandler)
	for _, h := range defaultHandlers {
		for _, k := range h.Keys() {
			keyHandler[k] = h
		}
	}
	out := make(map[string]json.RawMessage, len(style))
	for k, raw := range style {
		val := stringifyRaw(raw)
		if dims.PreservePercent && strings.Contains(val, "%") {
			out[k] = raw
			continue
		}
		if h, ok := keyHandler[k]; ok && val != "" {
			if resolved, ok := h.Resolve(k, val, dims); ok {
				out[k] = json.RawMessage(mustQuote(resolved))
				continue
			}
		}
		out[k] = raw
	}
	return out
}

func stringifyRaw(r json.RawMessage) string {
	if len(r) == 0 || string(r) == "null" {
		return ""
	}
	var s string
	if err := json.Unmarshal(r, &s); err == nil {
		return s
	}
	return string(r)
}

func mustQuote(s string) []byte {
	b, _ := json.Marshal(s)
	return b
}

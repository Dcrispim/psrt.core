package styleadapter

import "encoding/json"

func applyStrokeHTML(span StyleFragment, style map[string]json.RawMessage) {
	if span == nil {
		return
	}
	if raw, ok := style[KeyStroke]; ok {
		if v := StringifyCSSValue(raw); v != "" {
			span.Set("WebkitTextStroke", v)
		}
	}
	if raw, ok := style[KeyStrokeWidth]; ok {
		if v := StringifyCSSValue(raw); v != "" {
			span.Set("WebkitTextStrokeWidth", v)
		}
	}
	if raw, ok := style[KeyStrokeColor]; ok {
		if v := StringifyCSSValue(raw); v != "" {
			span.Set("WebkitTextStrokeColor", v)
		}
	}
}

func applyStrokeSVG(host StyleFragment, style map[string]json.RawMessage) {
	if host == nil {
		return
	}
	if raw, ok := style[KeyStroke]; ok {
		if v := StringifyCSSValue(raw); v != "" {
			host.Set("-webkit-text-stroke", v)
		}
	}
	if raw, ok := style[KeyStrokeWidth]; ok {
		if v := StringifyCSSValue(raw); v != "" {
			host.Set("-webkit-text-stroke-width", v)
		}
	}
	if raw, ok := style[KeyStrokeColor]; ok {
		if v := StringifyCSSValue(raw); v != "" {
			host.Set("-webkit-text-stroke-color", v)
		}
	}
}

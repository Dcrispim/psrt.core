//go:build js

package wasmbridge

import (
	"encoding/json"
	"syscall/js"

	"psrt/psrt"
	"psrt/psrt/editor"
)

func HandleSetStyleKey() js.Func {
	return wrap(func(args []js.Value) ([]byte, error) {
		var style psrt.Style
		if err := parseJSONArg(args, 0, &style); err != nil {
			return nil, err
		}
		key, err := stringArg(args, 1)
		if err != nil {
			return nil, err
		}
		value, err := stringArg(args, 2)
		if err != nil {
			return nil, err
		}
		out, err := editor.SetStyleKey(style, key, value)
		if err != nil {
			return nil, err
		}
		return json.Marshal(out)
	})
}

func HandleRemoveStyleKey() js.Func {
	return wrap(func(args []js.Value) ([]byte, error) {
		var style psrt.Style
		if err := parseJSONArg(args, 0, &style); err != nil {
			return nil, err
		}
		key, err := stringArg(args, 1)
		if err != nil {
			return nil, err
		}
		out, err := editor.RemoveStyleKey(style, key)
		if err != nil {
			return nil, err
		}
		return json.Marshal(out)
	})
}

func HandleMergeStyle() js.Func {
	return wrap(func(args []js.Value) ([]byte, error) {
		var style psrt.Style
		if err := parseJSONArg(args, 0, &style); err != nil {
			return nil, err
		}
		var partial json.RawMessage
		if err := parseJSONArg(args, 1, &partial); err != nil {
			return nil, err
		}
		out, err := editor.MergeStyle(style, partial)
		if err != nil {
			return nil, err
		}
		return json.Marshal(out)
	})
}

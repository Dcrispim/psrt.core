//go:build js

package wasmbridge

import (
	"encoding/json"
	"syscall/js"

	"psrt/psrt"
	"psrt/psrt/editor"
)

type positionFieldsJS struct {
	X        *float64 `json:"x"`
	Y        *float64 `json:"y"`
	Width    *float64 `json:"width"`
	TextSize *float64 `json:"textSize"`
}

type maskPositionFieldsJS struct {
	X      *float64 `json:"x"`
	Y      *float64 `json:"y"`
	Width  *float64 `json:"width"`
	Height *float64 `json:"height"`
}

type pathMaskPositionFieldsJS struct {
	X      *float64 `json:"x"`
	Y      *float64 `json:"y"`
	Width  *float64 `json:"width"`
	Height *float64 `json:"height"`
}

func toPositionFields(p positionFieldsJS) editor.PositionFields {
	return editor.PositionFields{X: p.X, Y: p.Y, Width: p.Width, TextSize: p.TextSize}
}

func toMaskPositionFields(p maskPositionFieldsJS) editor.MaskPositionFields {
	return editor.MaskPositionFields{X: p.X, Y: p.Y, Width: p.Width, Height: p.Height}
}

func toPathMaskPositionFields(p pathMaskPositionFieldsJS) editor.PathMaskPositionFields {
	return editor.PathMaskPositionFields{X: p.X, Y: p.Y, Width: p.Width, Height: p.Height}
}

func HandleAddConst() js.Func {
	return wrap(func(args []js.Value) ([]byte, error) {
		name, err := stringArg(args, 1)
		if err != nil {
			return nil, err
		}
		value, err := stringArg(args, 2)
		if err != nil {
			return nil, err
		}
		return mutateDocJSON(args, func(d *psrt.Document) error {
			return editor.AddConst(d, name, value)
		})
	})
}

func HandleRemoveConst() js.Func {
	return wrap(func(args []js.Value) ([]byte, error) {
		name, err := stringArg(args, 1)
		if err != nil {
			return nil, err
		}
		return mutateDocJSON(args, func(d *psrt.Document) error {
			return editor.RemoveConst(d, name)
		})
	})
}

func HandleSubstituteConstReferences() js.Func {
	return wrap(func(args []js.Value) ([]byte, error) {
		name, err := stringArg(args, 1)
		if err != nil {
			return nil, err
		}
		value, err := stringArg(args, 2)
		if err != nil {
			return nil, err
		}
		return mutateDocJSON(args, func(d *psrt.Document) error {
			editor.SubstituteConstReferences(d, name, value)
			return nil
		})
	})
}

func HandleRevertConstReferences() js.Func {
	return wrap(func(args []js.Value) ([]byte, error) {
		name, err := stringArg(args, 1)
		if err != nil {
			return nil, err
		}
		value, err := stringArg(args, 2)
		if err != nil {
			return nil, err
		}
		return mutateDocJSON(args, func(d *psrt.Document) error {
			editor.RevertConstReferences(d, name, value)
			return nil
		})
	})
}

func HandleAddFont() js.Func {
	return wrap(func(args []js.Value) ([]byte, error) {
		url, err := stringArg(args, 1)
		if err != nil {
			return nil, err
		}
		return mutateDocJSON(args, func(d *psrt.Document) error {
			return editor.AddFont(d, url)
		})
	})
}

func HandleRemoveFont() js.Func {
	return wrap(func(args []js.Value) ([]byte, error) {
		url, err := stringArg(args, 1)
		if err != nil {
			return nil, err
		}
		return mutateDocJSON(args, func(d *psrt.Document) error {
			return editor.RemoveFont(d, url)
		})
	})
}

func HandleRenamePage() js.Func {
	return wrap(func(args []js.Value) ([]byte, error) {
		oldName, err := stringArg(args, 1)
		if err != nil {
			return nil, err
		}
		newName, err := stringArg(args, 2)
		if err != nil {
			return nil, err
		}
		return mutateDocJSON(args, func(d *psrt.Document) error {
			return editor.RenamePage(d, oldName, newName)
		})
	})
}

func HandleSetPagePath() js.Func {
	return wrap(func(args []js.Value) ([]byte, error) {
		pageName, err := stringArg(args, 1)
		if err != nil {
			return nil, err
		}
		path, err := stringArg(args, 2)
		if err != nil {
			return nil, err
		}
		return mutateDocJSON(args, func(d *psrt.Document) error {
			return editor.SetPagePath(d, pageName, path)
		})
	})
}

func HandleSetPageStyle() js.Func {
	return wrap(func(args []js.Value) ([]byte, error) {
		pageName, err := stringArg(args, 1)
		if err != nil {
			return nil, err
		}
		key, err := stringArg(args, 2)
		if err != nil {
			return nil, err
		}
		value, err := stringArg(args, 3)
		if err != nil {
			return nil, err
		}
		var partial json.RawMessage
		if len(args) > 4 {
			b, _ := bytesArg(args, 4)
			partial = json.RawMessage(b)
		}
		return mutateDocJSON(args, func(d *psrt.Document) error {
			return editor.SetPageStyle(d, pageName, key, value, partial)
		})
	})
}

func HandleRemovePageStyleKey() js.Func {
	return wrap(func(args []js.Value) ([]byte, error) {
		pageName, err := stringArg(args, 1)
		if err != nil {
			return nil, err
		}
		key, err := stringArg(args, 2)
		if err != nil {
			return nil, err
		}
		return mutateDocJSON(args, func(d *psrt.Document) error {
			return editor.RemovePageStyleKey(d, pageName, key)
		})
	})
}

func HandleMovePage() js.Func {
	return wrap(func(args []js.Value) ([]byte, error) {
		pageName, _ := stringArg(args, 1)
		before, _ := stringArg(args, 2)
		after, _ := stringArg(args, 3)
		return mutateDocJSON(args, func(d *psrt.Document) error {
			return editor.MovePage(d, pageName, before, after)
		})
	})
}

func HandleAddPage() js.Func {
	return wrap(func(args []js.Value) ([]byte, error) {
		var page psrt.Page
		if err := parseJSONArg(args, 1, &page); err != nil {
			return nil, err
		}
		before, _ := stringArg(args, 2)
		after, _ := stringArg(args, 3)
		return mutateDocJSON(args, func(d *psrt.Document) error {
			return editor.AddPage(d, page, before, after)
		})
	})
}

func HandleRemovePage() js.Func {
	return wrap(func(args []js.Value) ([]byte, error) {
		name, err := stringArg(args, 1)
		if err != nil {
			return nil, err
		}
		return mutateDocJSON(args, func(d *psrt.Document) error {
			return editor.RemovePage(d, name)
		})
	})
}

func HandleSetTextStyle() js.Func {
	return wrap(func(args []js.Value) ([]byte, error) {
		pageName, err := stringArg(args, 1)
		if err != nil {
			return nil, err
		}
		textIndex := intArg(args, 2, -1)
		key, err := stringArg(args, 3)
		if err != nil {
			return nil, err
		}
		value, err := stringArg(args, 4)
		if err != nil {
			return nil, err
		}
		var partial json.RawMessage
		if len(args) > 5 {
			b, _ := bytesArg(args, 5)
			partial = json.RawMessage(b)
		}
		return mutateDocJSON(args, func(d *psrt.Document) error {
			return editor.SetTextStyle(d, pageName, textIndex, key, value, partial)
		})
	})
}

func HandleRemoveTextStyleKey() js.Func {
	return wrap(func(args []js.Value) ([]byte, error) {
		pageName, err := stringArg(args, 1)
		if err != nil {
			return nil, err
		}
		textIndex := intArg(args, 2, -1)
		key, err := stringArg(args, 3)
		if err != nil {
			return nil, err
		}
		return mutateDocJSON(args, func(d *psrt.Document) error {
			return editor.RemoveTextStyleKey(d, pageName, textIndex, key)
		})
	})
}

func HandleSetTextContent() js.Func {
	return wrap(func(args []js.Value) ([]byte, error) {
		pageName, err := stringArg(args, 1)
		if err != nil {
			return nil, err
		}
		textIndex := intArg(args, 2, -1)
		content, err := stringArg(args, 3)
		if err != nil {
			return nil, err
		}
		appendContent := boolArg(args, 4, false)
		return mutateDocJSON(args, func(d *psrt.Document) error {
			return editor.SetTextContent(d, pageName, textIndex, content, appendContent)
		})
	})
}

func HandleAddText() js.Func {
	return wrap(func(args []js.Value) ([]byte, error) {
		pageName, err := stringArg(args, 1)
		if err != nil {
			return nil, err
		}
		var text psrt.Text
		if err := parseJSONArg(args, 2, &text); err != nil {
			return nil, err
		}
		beforeIndex := optionalIndexArg(args, 3)
		afterIndex := optionalIndexArg(args, 4)
		return mutateDocJSON(args, func(d *psrt.Document) error {
			return editor.AddText(d, pageName, text, beforeIndex, afterIndex)
		})
	})
}

func HandleRemoveText() js.Func {
	return wrap(func(args []js.Value) ([]byte, error) {
		pageName, err := stringArg(args, 1)
		if err != nil {
			return nil, err
		}
		textIndex := intArg(args, 2, -1)
		return mutateDocJSON(args, func(d *psrt.Document) error {
			return editor.RemoveText(d, pageName, textIndex)
		})
	})
}

func HandleReorderTextRelative() js.Func {
	return wrap(func(args []js.Value) ([]byte, error) {
		pageName, _ := stringArg(args, 1)
		textIndex := intArg(args, 2, -1)
		beforeIndex := optionalIndexArg(args, 3)
		afterIndex := optionalIndexArg(args, 4)
		return mutateDocJSON(args, func(d *psrt.Document) error {
			return editor.ReorderTextRelative(d, pageName, textIndex, beforeIndex, afterIndex)
		})
	})
}

func HandleReorderTextTo() js.Func {
	return wrap(func(args []js.Value) ([]byte, error) {
		pageName, _ := stringArg(args, 1)
		textIndex := intArg(args, 2, -1)
		to := intArg(args, 3, 0)
		return mutateDocJSON(args, func(d *psrt.Document) error {
			return editor.ReorderTextTo(d, pageName, textIndex, to)
		})
	})
}

func HandleReorderTextByDelta() js.Func {
	return wrap(func(args []js.Value) ([]byte, error) {
		pageName, _ := stringArg(args, 1)
		textIndex := intArg(args, 2, -1)
		delta := intArg(args, 3, 0)
		return mutateDocJSON(args, func(d *psrt.Document) error {
			return editor.ReorderTextByDelta(d, pageName, textIndex, delta)
		})
	})
}

func HandleSetTextPosition() js.Func {
	return wrap(func(args []js.Value) ([]byte, error) {
		pageName, _ := stringArg(args, 1)
		textIndex := intArg(args, 2, -1)
		var pos positionFieldsJS
		if err := parseJSONArg(args, 3, &pos); err != nil {
			return nil, err
		}
		return mutateDocJSON(args, func(d *psrt.Document) error {
			return editor.SetTextPosition(d, pageName, textIndex, toPositionFields(pos))
		})
	})
}

func HandleNudgeTextPosition() js.Func {
	return wrap(func(args []js.Value) ([]byte, error) {
		pageName, _ := stringArg(args, 1)
		textIndex := intArg(args, 2, -1)
		var pos positionFieldsJS
		if err := parseJSONArg(args, 3, &pos); err != nil {
			return nil, err
		}
		return mutateDocJSON(args, func(d *psrt.Document) error {
			return editor.NudgeTextPosition(d, pageName, textIndex, toPositionFields(pos))
		})
	})
}

func HandleParseTextIndex() js.Func {
	return wrap(func(args []js.Value) ([]byte, error) {
		s, err := stringArg(args, 0)
		if err != nil {
			return nil, err
		}
		n, err := editor.ParseTextIndex(s)
		if err != nil {
			return nil, err
		}
		return json.Marshal(map[string]int{"index": n})
	})
}

func HandleSetMaskPosition() js.Func {
	return wrap(func(args []js.Value) ([]byte, error) {
		pageName, _ := stringArg(args, 1)
		maskIndex := intArg(args, 2, -1)
		var pos maskPositionFieldsJS
		if err := parseJSONArg(args, 3, &pos); err != nil {
			return nil, err
		}
		return mutateDocJSON(args, func(d *psrt.Document) error {
			return editor.SetMaskPosition(d, pageName, maskIndex, toMaskPositionFields(pos))
		})
	})
}

func HandleAddMask() js.Func {
	return wrap(func(args []js.Value) ([]byte, error) {
		pageName, err := stringArg(args, 1)
		if err != nil {
			return nil, err
		}
		var mask psrt.Mask
		if err := parseJSONArg(args, 2, &mask); err != nil {
			return nil, err
		}
		beforeIndex := optionalIndexArg(args, 3)
		afterIndex := optionalIndexArg(args, 4)
		return mutateDocJSON(args, func(d *psrt.Document) error {
			return editor.AddMask(d, pageName, mask, beforeIndex, afterIndex)
		})
	})
}

func HandleRemoveMask() js.Func {
	return wrap(func(args []js.Value) ([]byte, error) {
		pageName, _ := stringArg(args, 1)
		maskIndex := intArg(args, 2, -1)
		return mutateDocJSON(args, func(d *psrt.Document) error {
			return editor.RemoveMask(d, pageName, maskIndex)
		})
	})
}

func HandleSetMaskStyle() js.Func {
	return wrap(func(args []js.Value) ([]byte, error) {
		pageName, _ := stringArg(args, 1)
		maskIndex := intArg(args, 2, -1)
		key, err := stringArg(args, 3)
		if err != nil {
			return nil, err
		}
		value, err := stringArg(args, 4)
		if err != nil {
			return nil, err
		}
		var partial json.RawMessage
		if len(args) > 5 {
			b, _ := bytesArg(args, 5)
			partial = json.RawMessage(b)
		}
		return mutateDocJSON(args, func(d *psrt.Document) error {
			return editor.SetMaskStyle(d, pageName, maskIndex, key, value, partial)
		})
	})
}

func HandleRemoveMaskStyleKey() js.Func {
	return wrap(func(args []js.Value) ([]byte, error) {
		pageName, _ := stringArg(args, 1)
		maskIndex := intArg(args, 2, -1)
		key, err := stringArg(args, 3)
		if err != nil {
			return nil, err
		}
		return mutateDocJSON(args, func(d *psrt.Document) error {
			return editor.RemoveMaskStyleKey(d, pageName, maskIndex, key)
		})
	})
}

func HandleSetPathMaskPosition() js.Func {
	return wrap(func(args []js.Value) ([]byte, error) {
		pageName, _ := stringArg(args, 1)
		maskIndex := intArg(args, 2, -1)
		var pos pathMaskPositionFieldsJS
		if err := parseJSONArg(args, 3, &pos); err != nil {
			return nil, err
		}
		return mutateDocJSON(args, func(d *psrt.Document) error {
			return editor.SetPathMaskPosition(d, pageName, maskIndex, toPathMaskPositionFields(pos))
		})
	})
}

func HandleAddPathMask() js.Func {
	return wrap(func(args []js.Value) ([]byte, error) {
		pageName, err := stringArg(args, 1)
		if err != nil {
			return nil, err
		}
		var mask psrt.PathMask
		if err := parseJSONArg(args, 2, &mask); err != nil {
			return nil, err
		}
		beforeIndex := optionalIndexArg(args, 3)
		afterIndex := optionalIndexArg(args, 4)
		return mutateDocJSON(args, func(d *psrt.Document) error {
			return editor.AddPathMask(d, pageName, mask, beforeIndex, afterIndex)
		})
	})
}

func HandleRemovePathMask() js.Func {
	return wrap(func(args []js.Value) ([]byte, error) {
		pageName, _ := stringArg(args, 1)
		maskIndex := intArg(args, 2, -1)
		return mutateDocJSON(args, func(d *psrt.Document) error {
			return editor.RemovePathMask(d, pageName, maskIndex)
		})
	})
}

func HandleSetPathMaskStyle() js.Func {
	return wrap(func(args []js.Value) ([]byte, error) {
		pageName, _ := stringArg(args, 1)
		maskIndex := intArg(args, 2, -1)
		key, err := stringArg(args, 3)
		if err != nil {
			return nil, err
		}
		value, err := stringArg(args, 4)
		if err != nil {
			return nil, err
		}
		var partial json.RawMessage
		if len(args) > 5 {
			b, _ := bytesArg(args, 5)
			partial = json.RawMessage(b)
		}
		return mutateDocJSON(args, func(d *psrt.Document) error {
			return editor.SetPathMaskStyle(d, pageName, maskIndex, key, value, partial)
		})
	})
}

func HandleRemovePathMaskStyleKey() js.Func {
	return wrap(func(args []js.Value) ([]byte, error) {
		pageName, _ := stringArg(args, 1)
		maskIndex := intArg(args, 2, -1)
		key, err := stringArg(args, 3)
		if err != nil {
			return nil, err
		}
		return mutateDocJSON(args, func(d *psrt.Document) error {
			return editor.RemovePathMaskStyleKey(d, pageName, maskIndex, key)
		})
	})
}

func HandleSetPathMaskPath() js.Func {
	return wrap(func(args []js.Value) ([]byte, error) {
		pageName, _ := stringArg(args, 1)
		maskIndex := intArg(args, 2, -1)
		path, err := stringArg(args, 3)
		if err != nil {
			return nil, err
		}
		return mutateDocJSON(args, func(d *psrt.Document) error {
			return editor.SetPathMaskPath(d, pageName, maskIndex, path)
		})
	})
}

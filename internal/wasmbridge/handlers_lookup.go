//go:build js

package wasmbridge

import (
	"encoding/json"
	"syscall/js"

	"psrt/psrt/editor"
)

func HandleFindPage() js.Func {
	return wrap(func(args []js.Value) ([]byte, error) {
		doc, err := loadDocJSON(args)
		if err != nil {
			return nil, err
		}
		name, err := stringArg(args, 1)
		if err != nil {
			return nil, err
		}
		p, err := editor.FindPage(&doc, name)
		if err != nil {
			return nil, err
		}
		return json.Marshal(p)
	})
}

func HandleFindPageIndex() js.Func {
	return wrap(func(args []js.Value) ([]byte, error) {
		doc, err := loadDocJSON(args)
		if err != nil {
			return nil, err
		}
		name, err := stringArg(args, 1)
		if err != nil {
			return nil, err
		}
		idx, err := editor.FindPageIndex(&doc, name)
		if err != nil {
			return nil, err
		}
		return json.Marshal(map[string]int{"index": idx})
	})
}

func HandleFindTextByIndex() js.Func {
	return wrap(func(args []js.Value) ([]byte, error) {
		doc, err := loadDocJSON(args)
		if err != nil {
			return nil, err
		}
		pageName, err := stringArg(args, 1)
		if err != nil {
			return nil, err
		}
		textIndex := intArg(args, 2, -1)
		page, err := editor.FindPage(&doc, pageName)
		if err != nil {
			return nil, err
		}
		t, _, err := editor.FindTextByIndex(page, textIndex)
		if err != nil {
			return nil, err
		}
		return json.Marshal(t)
	})
}

func HandleFindMaskByIndex() js.Func {
	return wrap(func(args []js.Value) ([]byte, error) {
		doc, err := loadDocJSON(args)
		if err != nil {
			return nil, err
		}
		pageName, err := stringArg(args, 1)
		if err != nil {
			return nil, err
		}
		maskIndex := intArg(args, 2, -1)
		page, err := editor.FindPage(&doc, pageName)
		if err != nil {
			return nil, err
		}
		m, _, err := editor.FindMaskByIndex(page, maskIndex)
		if err != nil {
			return nil, err
		}
		return json.Marshal(m)
	})
}

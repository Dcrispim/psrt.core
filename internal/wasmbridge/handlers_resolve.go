//go:build js

package wasmbridge

import (
	"syscall/js"

	"psrt/compilesvg"
	"psrt/psrt"
)

func HandleResolveDocument() js.Func {
	return wrap(func(args []js.Value) ([]byte, error) {
		b, err := bytesArg(args, 0)
		if err != nil {
			return nil, err
		}
		doc, err := psrt.ParseJSON(b)
		if err != nil {
			return nil, err
		}
		resolved := compilesvg.ResolveDocument(doc)
		return exportDocJSON(resolved)
	})
}

func HandleResolveDocumentStrict() js.Func {
	return wrap(func(args []js.Value) ([]byte, error) {
		b, err := bytesArg(args, 0)
		if err != nil {
			return nil, err
		}
		doc, err := psrt.ParseJSON(b)
		if err != nil {
			return nil, err
		}
		resolved, err := compilesvg.ResolveDocumentStrict(doc)
		if err != nil {
			return nil, err
		}
		return exportDocJSON(resolved)
	})
}

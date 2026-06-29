//go:build js

package wasmbridge

import (
	"syscall/js"

	"github.com/Dcrispim/psrt.core/compilesvg"
	"github.com/Dcrispim/psrt.core/psrt"
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
		// Keep @type:render@ tokens so interactive readers (react-image) render them.
		resolved := compilesvg.ResolveDocumentKeepInteractive(doc)
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
		resolved, err := compilesvg.ResolveDocumentStrictKeepInteractive(doc)
		if err != nil {
			return nil, err
		}
		return exportDocJSON(resolved)
	})
}

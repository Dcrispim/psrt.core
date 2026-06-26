//go:build js

package wasmbridge

import (
	"encoding/json"
	"syscall/js"

	"github.com/Dcrispim/psrt.core/internal/visualapp"
)

func HandleAdaptEntriesForWeb() js.Func {
	return wrap(func(args []js.Value) ([]byte, error) {
		entriesJSON, err := stringArg(args, 0)
		if err != nil {
			return nil, err
		}
		canvasW := intArg(args, 1, 0)
		canvasH := intArg(args, 2, 0)
		zoom := floatArg(args, 3, 1)
		out, err := visualapp.AdaptEntriesForWeb(entriesJSON, canvasW, canvasH, zoom)
		if err != nil {
			return nil, err
		}
		return json.Marshal(out)
	})
}

func HandleMergePageDocumentPSRT() js.Func {
	return wrap(func(args []js.Value) ([]byte, error) {
		fullDocJSON, err := stringArg(args, 0)
		if err != nil {
			return nil, err
		}
		pageName, err := stringArg(args, 1)
		if err != nil {
			return nil, err
		}
		psrtText, err := stringArg(args, 2)
		if err != nil {
			return nil, err
		}
		out, err := visualapp.MergePageDocumentPSRT(fullDocJSON, pageName, psrtText)
		if err != nil {
			return nil, err
		}
		return []byte(out), nil
	})
}

func HandleFormatPageDocumentJSON() js.Func {
	return wrap(func(args []js.Value) ([]byte, error) {
		docJSON, err := stringArg(args, 0)
		if err != nil {
			return nil, err
		}
		pageName, err := stringArg(args, 1)
		if err != nil {
			return nil, err
		}
		out, err := visualapp.FormatPageDocumentJSON(docJSON, pageName)
		if err != nil {
			return nil, err
		}
		return []byte(out), nil
	})
}

func floatArg(args []js.Value, i int, def float64) float64 {
	if i >= len(args) || args[i].Type() != js.TypeNumber {
		return def
	}
	return args[i].Float()
}

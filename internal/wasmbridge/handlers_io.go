//go:build js

package wasmbridge

import (
	"encoding/json"
	"syscall/js"

	"github.com/Dcrispim/psrt.core/psrt"
	"github.com/Dcrispim/psrt.core/psrt/editor"
)

func HandleParse() js.Func {
	return wrap(func(args []js.Value) ([]byte, error) {
		s, err := stringArg(args, 0)
		if err != nil {
			return nil, err
		}
		doc, err := psrt.ParseString(s)
		if err != nil {
			return nil, err
		}
		return exportDocJSON(doc)
	})
}

func HandleParseFast() js.Func {
	return wrap(func(args []js.Value) ([]byte, error) {
		s, err := stringArg(args, 0)
		if err != nil {
			return nil, err
		}
		doc, err := psrt.ParseFastString(s)
		if err != nil {
			return nil, err
		}
		return exportDocJSON(doc)
	})
}

func HandleLoadSource() js.Func {
	return wrapString(func(args []js.Value) (string, error) {
		raw, err := stringArg(args, 0)
		if err != nil {
			return "", err
		}
		url, err := stringArg(args, 1)
		if err != nil {
			return "", err
		}
		return psrt.LoadSource(raw, url)
	})
}

func HandleStringify() js.Func {
	return wrap(func(args []js.Value) ([]byte, error) {
		b, err := bytesArg(args, 0)
		if err != nil {
			return nil, err
		}
		doc, err := psrt.ParseJSON(b)
		if err != nil {
			return nil, err
		}
		return psrt.FormatPSRT(doc, false)
	})
}

func HandleFormatDocument() js.Func {
	return wrap(func(args []js.Value) ([]byte, error) {
		b, err := bytesArg(args, 0)
		if err != nil {
			return nil, err
		}
		doc, err := psrt.ParseJSON(b)
		if err != nil {
			return nil, err
		}
		d := &doc
		return editor.FormatDocument(d)
	})
}

// HandleConfigure applies process-wide PSRT formatting options. Meant to be
// called once, right after boot — see initPsrt() in the JS SDK.
func HandleConfigure() js.Func {
	return wrap(func(args []js.Value) ([]byte, error) {
		var opts psrt.FormatOptions
		if b, err := bytesArg(args, 0); err == nil && len(b) > 0 {
			if err := json.Unmarshal(b, &opts); err != nil {
				return nil, err
			}
		}
		psrt.Configure(opts)
		return nil, nil
	})
}

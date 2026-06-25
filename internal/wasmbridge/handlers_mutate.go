//go:build js

package wasmbridge

import (
	"syscall/js"

	"github.com/Dcrispim/psrt.core/psrt"
)

func loadDocJSON(args []js.Value) (psrt.Document, error) {
	b, err := bytesArg(args, 0)
	if err != nil {
		return psrt.Document{}, err
	}
	return psrt.ParseJSON(b)
}

func saveDocJSON(doc psrt.Document) ([]byte, error) {
	return exportDocJSON(cloneDoc(doc))
}

func mutateDocJSON(args []js.Value, fn func(*psrt.Document) error) ([]byte, error) {
	doc, err := loadDocJSON(args)
	if err != nil {
		return nil, err
	}
	d := &doc
	if err := fn(d); err != nil {
		return nil, err
	}
	return saveDocJSON(doc)
}

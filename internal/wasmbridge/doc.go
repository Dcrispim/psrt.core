//go:build js

package wasmbridge

import (
	"encoding/json"

	"github.com/Dcrispim/psrt.core/psrt"
)

func exportDocJSON(doc psrt.Document) ([]byte, error) {
	return json.Marshal(doc)
}

func cloneDoc(doc psrt.Document) psrt.Document {
	b, err := json.Marshal(doc)
	if err != nil {
		return doc
	}
	out, err := psrt.ParseJSON(b)
	if err != nil {
		return doc
	}
	return out
}

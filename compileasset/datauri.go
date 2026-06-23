package compileasset

import (
	"encoding/base64"
	"fmt"
)

// EncodeDataURI returns a data URI for the given MIME type and bytes.
func EncodeDataURI(mime string, b []byte) string {
	return fmt.Sprintf("data:%s;base64,%s", mime, base64.StdEncoding.EncodeToString(b))
}

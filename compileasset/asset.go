package compileasset

// Asset holds downloaded bytes with a MIME type suitable for data URIs.
type Asset struct {
	Bytes []byte
	MIME  string
}

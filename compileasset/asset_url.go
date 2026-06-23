package compileasset

// AssetRef returns either the resolved URL (linksOnly) or a data URI for embedded assets.
func AssetRef(resolvedURL string, asset Asset, linksOnly bool) string {
	if linksOnly {
		return resolvedURL
	}
	return EncodeDataURI(asset.MIME, asset.Bytes)
}

// FontSrcURL returns a CSS url(...) value for @font-face src.
func FontSrcURL(fontURL string, asset Asset, linksOnly bool) string {
	if linksOnly {
		return "url(" + fontURL + ")"
	}
	return "url(" + EncodeDataURI(asset.MIME, asset.Bytes) + ")"
}

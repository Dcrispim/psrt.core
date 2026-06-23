// Package compileopts holds shared flags for HTML and SVG compilation.
package compileopts

// Options configures PSRT compile output.
type Options struct {
	// LinksOnly keeps original asset URLs instead of embedding data URIs.
	LinksOnly bool
	// NoScript omits the HTML variant switcher script (Ctrl+L) for cleaner output.
	NoScript bool
}

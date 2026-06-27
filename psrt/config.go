package psrt

// FormatOptions configures process-wide PSRT serialization behavior. It is
// set once via Configure during process startup — psrt formatting has no
// reload mechanism, mirroring how the JS SDK's initPsrt() boots the WASM
// core exactly once per process.
type FormatOptions struct {
	// PathCommandsPerLine caps how many SVG path commands are written per
	// line in a ~~ block body. Zero/negative falls back to the default.
	PathCommandsPerLine int `json:"pathCommandsPerLine"`
}

const defaultPathCommandsPerLine = 20

var formatOptions = FormatOptions{PathCommandsPerLine: defaultPathCommandsPerLine}

// Configure overrides process-wide formatting options. Call once during
// startup, before any FormatPSRT call — it is not safe to call concurrently
// with formatting.
func Configure(opts FormatOptions) {
	if opts.PathCommandsPerLine <= 0 {
		opts.PathCommandsPerLine = defaultPathCommandsPerLine
	}
	formatOptions = opts
}

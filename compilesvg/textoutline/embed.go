package textoutline

import _ "embed"

const DefaultFontFamily = "PSRTDefault"

//go:embed opentype.min.js
var opentypeJS string

//go:embed extract.js
var extractJS string

//go:embed default.ttf
var defaultFontBytes []byte

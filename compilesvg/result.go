package compilesvg

// PageSVGResult is the output of rendering one page to SVG.
type PageSVGResult struct {
	Data               []byte
	UsedGoTextFallback bool
}

// GoTextFallbackNotice is shown in the GUI when the go-text fallback is used.
const GoTextFallbackNotice = "Compilação SVG com motor go-text (Chromium não encontrado). Pode haver pequenas diferenças em relação ao preview."

package textoutline

// OutlinedPath is one SVG path for a text run.
type OutlinedPath struct {
	D           string  `json:"d"`
	Fill        string  `json:"fill,omitempty"`
	Stroke      string  `json:"stroke,omitempty"`
	StrokeWidth string  `json:"strokeWidth,omitempty"`
	Opacity     string  `json:"opacity,omitempty"`
	PaintOrder  string  `json:"paintOrder,omitempty"`
}

// OutlinedBlock holds vector paths and metadata for one PSRT text block.
type OutlinedBlock struct {
	Index     int            `json:"index"`
	Paths     []OutlinedPath `json:"paths"`
	PlainText string         `json:"plainText"`
	FilterID  string         `json:"filterId,omitempty"`
	Transform string         `json:"transform,omitempty"`
}

// BlockStyle carries layout hints for the go-text fallback.
type BlockStyle struct {
	FontSizePx  float64
	LineHeight  float64
	Color       string
	Stroke      string
	StrokeWidth string
	TextAlign   string
	PadTop      int
	PadLeft     int
	ContentW    int
	FontFamily  string
	Bold        bool
}

// BlockInput describes one text block in the layout snapshot.
type BlockInput struct {
	Index      int
	X, Y       int
	Width      int
	Height     int
	ClassAttr  string
	InnerClass string
	TextHTML   string
	PlainText  string
	FilterID   string
	Transform  string
	Style      BlockStyle
}

// PageInput is the layout snapshot for one page.
type PageInput struct {
	CanvasW  int
	CanvasH  int
	CSS      string
	Blocks   []BlockInput
	Fonts    map[string]FontBytes // family name -> raw font bytes
}

// FontBytes wraps font file bytes for embedding.
type FontBytes struct {
	MIME  string
	Bytes []byte
}

// outlineResult is the JSON shape returned from the browser script.
type outlineResult struct {
	Blocks []OutlinedBlock `json:"blocks"`
}

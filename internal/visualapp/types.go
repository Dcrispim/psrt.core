package visualapp

// OpenFileResult is returned when the user opens a PSRT file.
type OpenFileResult struct {
	FilePath string `json:"filePath"`
	Document string `json:"document"`
}

// UIState is sent to the frontend after changes.
type UIState struct {
	FilePath      string            `json:"filePath"`
	ActivePage    string            `json:"activePage"`
	SelectedIndex int               `json:"selectedIndex"` // -1 none
	Pages         []PageSummary     `json:"pages"`
	Page          *PageDetail       `json:"page,omitempty"`
	Texts         []TextDetail      `json:"texts,omitempty"`
	Masks         []MaskDetail      `json:"masks,omitempty"`
	Text          *TextDetail       `json:"text,omitempty"`
	Mask          *MaskDetail       `json:"mask,omitempty"`
	Fonts         []string          `json:"fonts"`
	Consts        map[string]string `json:"consts"`
	AutoCompile   bool              `json:"autoCompile"`
}

type PageSummary struct {
	Name     string `json:"name"`
	ImageURL string `json:"imageUrl,omitempty"`
}

type PageDetail struct {
	Name     string `json:"name"`
	ImageURL string `json:"imageUrl"`
	Style    string `json:"style"`
}

type TextDetail struct {
	Index    int     `json:"index"`
	X        float64 `json:"x"`
	Y        float64 `json:"y"`
	Width    float64 `json:"width"`
	TextSize float64 `json:"textSize"`
	Content  string  `json:"content"`
	ImageRef string  `json:"imageRef"`
	Style    string  `json:"style"`
}

// TextPatch updates a text block from the UI.
type TextPatch struct {
	Content      *string            `json:"content,omitempty"`
	Append       bool               `json:"append"`
	X            *float64           `json:"x,omitempty"`
	Y            *float64           `json:"y,omitempty"`
	Width        *float64           `json:"width,omitempty"`
	TextSize     *float64           `json:"textSize,omitempty"`
	ImageRef     *string            `json:"imageRef,omitempty"`
	StyleSet     map[string]string  `json:"styleSet,omitempty"`
	StyleRemove  []string           `json:"styleRemove,omitempty"`
}

// MaskDetail is a static == block for the UI.
type MaskDetail struct {
	Index    int     `json:"index"`
	X        float64 `json:"x"`
	Y        float64 `json:"y"`
	Width    float64 `json:"width"`
	Height   float64 `json:"height"`
	ImageRef string  `json:"imageRef"`
	Style    string  `json:"style"`
}

// MaskPatch updates a mask block from the UI.
type MaskPatch struct {
	X           *float64          `json:"x,omitempty"`
	Y           *float64          `json:"y,omitempty"`
	Width       *float64          `json:"width,omitempty"`
	Height      *float64          `json:"height,omitempty"`
	ImageRef    *string           `json:"imageRef,omitempty"`
	StyleSet    map[string]string `json:"styleSet,omitempty"`
	StyleRemove []string          `json:"styleRemove,omitempty"`
}

// PagePatch updates page fields.
type PagePatch struct {
	Name        *string           `json:"name,omitempty"`
	ImageURL    *string           `json:"imageUrl,omitempty"`
	StyleSet    map[string]string `json:"styleSet,omitempty"`
	StyleRemove []string          `json:"styleRemove,omitempty"`
}

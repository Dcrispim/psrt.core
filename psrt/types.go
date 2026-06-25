package psrt



import "encoding/json"



// Document is the in-memory representation of a complete PSRT file after parsing.

type Document struct {

	Pages  []Page            `json:"pages"`

	Fonts  []string          `json:"fonts"`

	Consts map[string]string `json:"consts"`

	// Sources maps original asset URLs/paths to embedded data: URIs ($SOURCE block).
	Sources map[string]string `json:"sources,omitempty"`

}



// Page is a logical canvas: background/page style JSON, base image URL, and blocks.

// Corresponds to $ START <name> | <page-style> | <image-url> ... $ END <name> .

type Page struct {

	Name     string `json:"name"`

	Style    Style  `json:"style"`

	ImageURL string `json:"imageUrl"`

	Texts    []Text `json:"texts"`

	Masks    []Mask `json:"masks,omitempty"`

	PathMasks []PathMask `json:"pathMasks,omitempty"`

}



// BaseBlock holds fields shared by text and mask blocks.

type BaseBlock struct {

	X        float64 `json:"x"`

	Y        float64 `json:"y"`

	Width    float64 `json:"width"`

	Style    Style   `json:"style"`

	Index    int     `json:"index"`

	ImageRef string  `json:"imageRef,omitempty"`

}



// Text is a positioned text block relative to the page image (percent 0–100).

// Header: >> <X>-<Y>-<Width>-<TextSize> | <style-text> | <index> [| optional image ref].

// TextSize is a percent of min(canvas width, canvas height); see TextFontSizePx.

type Text struct {

	BaseBlock

	TextSize float64 `json:"textSize"`

	Content  string  `json:"content"`

}



// Mask is a static coverage block (no text body).

// Header: == <X>-<Y>-<Width>-<Height> | <style-text> | <index> [| optional image ref].

// Height is a percent of page image height (0–100).

type Mask struct {

	BaseBlock

	Height float64 `json:"height"`

}



// PathMask is a coverage block whose shape is an arbitrary SVG path (single sub-path).

// Header: ~~ <X>-<Y>-<Width>-<Height> | <style-text> | <index> [| optional image ref].

// X, Y, Width, Height define the local box (percent of page image); Path is interpreted

// as a 0-100 viewBox relative to that box, not to the page.

type PathMask struct {

	BaseBlock

	Height float64 `json:"height"`

	// Path is the SVG path `d` attribute content, already concatenated/normalized.

	Path string `json:"path"`

}



// Style is JSON object bytes; aliasing json.RawMessage keeps Marshal/Unmarshal inline JSON (not base64).

type Style = json.RawMessage



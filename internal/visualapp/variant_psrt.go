package visualapp

// VariantPSRT is extra PSRT text bundled into HTML export (from file picker when paths are unavailable).
type VariantPSRT struct {
	Label   string `json:"label"`
	Content string `json:"content"`
}

package psrt

// NewText builds a text block (convenience for tests and callers).
func NewText(x, y, width, textSize float64, index int, style Style, content, imageRef string) Text {
	return Text{
		BaseBlock: BaseBlock{X: x, Y: y, Width: width, Style: style, Index: index, ImageRef: imageRef},
		TextSize:  textSize,
		Content:   content,
	}
}

// NewMask builds a mask block.
func NewMask(x, y, width, height float64, index int, style Style, imageRef string) Mask {
	return Mask{
		BaseBlock: BaseBlock{X: x, Y: y, Width: width, Style: style, Index: index, ImageRef: imageRef},
		Height:    height,
	}
}

// NewPathMask builds a path mask block.
func NewPathMask(x, y, width, height float64, index int, style Style, imageRef, path string) PathMask {
	return PathMask{
		BaseBlock: BaseBlock{X: x, Y: y, Width: width, Style: style, Index: index, ImageRef: imageRef},
		Height:    height,
		Path:      path,
	}
}

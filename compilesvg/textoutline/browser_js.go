//go:build js

package textoutline

import (
	"context"
	"fmt"
)

// Outliner is unavailable under GOOS=js; text layout uses go-text fallback.
type Outliner struct{}

// NewOutliner is not supported in WASM builds.
func NewOutliner(_ context.Context, _ string) (*Outliner, error) {
	return nil, fmt.Errorf("textoutline: chromium not available in wasm")
}

// Close is a no-op in WASM builds.
func (o *Outliner) Close() {}

// OutlinePage is not supported in WASM builds.
func (o *Outliner) OutlinePage(_ context.Context, _ PageInput) ([]OutlinedBlock, error) {
	return nil, fmt.Errorf("textoutline: chromium not available in wasm")
}

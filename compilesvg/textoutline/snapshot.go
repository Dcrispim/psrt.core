package textoutline

import (
	"encoding/base64"
	"fmt"
	"strings"
)

// BuildSnapshotHTML returns a self-contained HTML document for headless layout.
func BuildSnapshotHTML(in PageInput) string {
	var b strings.Builder
	b.WriteString("<!DOCTYPE html><html><head><meta charset=\"utf-8\"><style>")
	b.WriteString(in.CSS)
	b.WriteString(`html,body{margin:0;padding:0;overflow:hidden;}
.psrt-snapshot-root{position:relative;margin:0;padding:0;}
.psrt-text-block{position:absolute;box-sizing:border-box;overflow:hidden;margin:0;}
`)
	b.WriteString("</style>")
	if len(in.Fonts) > 0 {
		b.WriteString("<script>window.__psrtFonts={")
		first := true
		for name, fb := range in.Fonts {
			if !first {
				b.WriteByte(',')
			}
			first = false
			enc := base64.StdEncoding.EncodeToString(fb.Bytes)
			fmt.Fprintf(&b, "%q:%q", name, enc)
		}
		b.WriteString("};</script>")
	}
	b.WriteString("</head><body><div class=\"psrt-snapshot-root\">")
	for _, blk := range in.Blocks {
		fmt.Fprintf(&b,
			`<div id="psrt-block-%d" class="psrt-text-block" data-block-index="%d" style="left:%dpx;top:%dpx;width:%dpx;height:%dpx">`,
			blk.Index, blk.Index, blk.X, blk.Y, blk.Width, blk.Height)
		b.WriteString(`<div class="`)
		b.WriteString(blk.ClassAttr)
		b.WriteString(`">`)
		if blk.TextHTML != "" {
			b.WriteString(`<span class="`)
			b.WriteString(blk.InnerClass)
			b.WriteString(`">`)
			b.WriteString(blk.TextHTML)
			b.WriteString(`</span>`)
		}
		b.WriteString(`</div></div>`)
	}
	b.WriteString(`</div></body></html>`)
	return b.String()
}

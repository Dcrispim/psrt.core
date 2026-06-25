package compilesvg

import (
	"context"
	"fmt"
	"strings"

	"psrt/compileasset"
	"psrt/compileopts"
	"psrt/compilesvg/textoutline"
	"psrt/psrt"
)

// RenderPageSVG produces a standalone SVG for one page.
func RenderPageSVG(
	p *psrt.Page,
	pageSlug string,
	fontURLs []string,
	assets map[string]compileasset.Asset,
) ([]byte, error) {
	res, err := RenderPageSVGWithContext(context.Background(), p, pageSlug, fontURLs, assets, compileopts.Options{})
	return res.Data, err
}

// RenderPageSVGWithContext compiles one page using ctx for text outlining.
func RenderPageSVGWithContext(
	ctx context.Context,
	p *psrt.Page,
	pageSlug string,
	fontURLs []string,
	assets map[string]compileasset.Asset,
	opts compileopts.Options,
) (PageSVGResult, error) {
	imgURL := strings.TrimSpace(p.ImageURL)
	a, ok := assets[imgURL]
	if !opts.LinksOnly && !ok {
		return PageSVGResult{}, fmt.Errorf("missing fetched asset for page image %q", imgURL)
	}
	var w, h int
	if ok {
		w, h = ImageDimensions(a.Bytes, a.MIME)
	} else {
		w, h = ImageDimensions(nil, "")
	}
	css := BuildPageSVGDefsCSS(pageSlug, p.Style)
	filters := collectFilterDefs(pageSlug, p.Texts, w, h)
	legacyFilters, _ := buildTextShadowFilters(pageSlug, p.Texts)
	filters += legacyFilters
	imgData := compileasset.AssetRef(imgURL, a, opts.LinksOnly)

	outlined, usedGoText, err := outlinePageTexts(ctx, p, pageSlug, w, h, fontURLs, assets)
	if err != nil {
		return PageSVGResult{}, err
	}

	var b strings.Builder
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" width="%d" height="%d" viewBox="0 0 %d %d" id="%s" class="%s">`,
		w, h, w, h, PageID(pageSlug), PageClass(pageSlug))
	if css != "" || filters != "" {
		b.WriteString(`<defs>`)
		if css != "" {
			b.WriteString(`<style type="text/css">`)
			b.WriteString(css)
			b.WriteString(`</style>`)
		}
		b.WriteString(filters)
		b.WriteString(`</defs>`)
	}

	if fill := compileasset.BackgroundColorFromStyle(p.Style); fill != "" {
		fmt.Fprintf(&b, `<rect id="%s" class="%s" x="0" y="0" width="%d" height="%d"/>`,
			PageBgID(pageSlug), PageBgClass(pageSlug), w, h)
	}

	fmt.Fprintf(&b, `<image id="%s" xlink:href="%s" x="0" y="0" width="%d" height="%d" preserveAspectRatio="xMidYMid meet"/>`,
		PageImageID(pageSlug), imgData, w, h)

	b.WriteString(`<g id="`)
	b.WriteString(PageTextsID(pageSlug))
	b.WriteString(`">`)

	for _, entry := range psrt.PageBlocksByIndex(p) {
		switch entry.Kind {
		case psrt.BlockText:
			if entry.Text == nil {
				continue
			}
			ob := outlined[entry.Text.Index]
			if err := writeTextBlock(&b, entry.Text, pageSlug, w, h, assets, &ob, opts.LinksOnly); err != nil {
				return PageSVGResult{}, err
			}
		case psrt.BlockMask:
			if entry.Mask == nil {
				continue
			}
			if err := writeMaskBlock(&b, entry.Mask, pageSlug, w, h, assets, opts.LinksOnly); err != nil {
				return PageSVGResult{}, err
			}
		case psrt.BlockPathMask:
			if entry.PathMask == nil {
				continue
			}
			if err := writePathMaskBlock(&b, entry.PathMask, pageSlug, w, h, assets, opts.LinksOnly); err != nil {
				return PageSVGResult{}, err
			}
		}
	}

	b.WriteString(`</g></svg>`)
	return PageSVGResult{Data: []byte(b.String()), UsedGoTextFallback: usedGoText}, nil
}

func writeTextBlock(
	b *strings.Builder,
	t *psrt.Text,
	pageSlug string,
	canvasW, canvasH int,
	assets map[string]compileasset.Asset,
	outlined *textoutline.OutlinedBlock,
	linksOnly bool,
) error {
	content := psrt.NormalizeTextContent(t.Content)
	x, y, width, foH := TextBlockGeometry(t, content, canvasW, canvasH)
	fontPx := psrt.TextFontSizePx(t.TextSize, canvasW, canvasH)
	corners := compileasset.ParseBorderRadiusCorners(t.Style, canvasW, canvasH, fontPx)
	frags := adaptTextFragments(t, pageSlug, canvasW, canvasH, fontPx)

	b.WriteString(`<g id="`)
	b.WriteString(TextWrapID(pageSlug, t.Index))
	b.WriteString(`">`)

	imgRef := strings.TrimSpace(t.ImageRef)
	var imgHTML string
	if imgRef != "" && compileasset.IsAssetReference(imgRef) {
		if aa, ok := assets[imgRef]; linksOnly || (ok && IsImageMIME(aa.MIME)) {
			refURI := compileasset.AssetRef(imgRef, aa, linksOnly)
			imgHTML = fmt.Sprintf(`<image href="%s" alt=""/>`, refURI)
		}
	}

	writeTextBlockFromFragments(b, frags, pageSlug, t.Index, imgHTML, outlined, x, y, width, foH, corners)
	b.WriteString(`</g>`)
	return nil
}

func svgAttrEscape(s string) string {
	return strings.ReplaceAll(s, `"`, "&quot;")
}

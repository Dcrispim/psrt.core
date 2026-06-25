package compilesvg

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/Dcrispim/psrt.core/compileasset"
	"github.com/Dcrispim/psrt.core/compilesvg/textoutline"
	"github.com/Dcrispim/psrt.core/psrt"
	"github.com/Dcrispim/psrt.core/styleadapter"
)

func buildPageOutlineInput(
	p *psrt.Page,
	pageSlug string,
	canvasW, canvasH int,
	fontURLs []string,
	assets map[string]compileasset.Asset,
) textoutline.PageInput {
	css := BuildPageTextStylesheet(pageSlug, p.Style, p.Texts, canvasW, canvasH, fontURLs, assets)
	css += defaultFontFaceCSS()
	fonts := fontMapForOutline(fontURLs, assets)

	var blocks []textoutline.BlockInput
	for i := range p.Texts {
		t := &p.Texts[i]
		content := psrt.NormalizeTextContent(t.Content)
		x, y, width, height := TextBlockGeometry(t, content, canvasW, canvasH)
		fontPx := psrt.TextFontSizePx(t.TextSize, canvasW, canvasH)
		frags := adaptTextFragments(t, pageSlug, canvasW, canvasH, fontPx)

		blocks = append(blocks, textoutline.BlockInput{
			Index:      t.Index,
			X:          x,
			Y:          y,
			Width:      width,
			Height:     height,
			ClassAttr:  TextClassAttr(pageSlug, t.Index),
			InnerClass: TextInnerClass(pageSlug, t.Index),
			TextHTML:   psrt.RenderInlineHTML(content),
			PlainText:  psrt.PlainTextForLayout(content),
			FilterID:   filterIDFromFragments(frags),
			Transform:  transformFromFragments(frags),
			Style:      blockStyleForText(t, canvasW, canvasH, fontURLs, assets),
		})
	}

	return textoutline.PageInput{
		CanvasW: canvasW,
		CanvasH: canvasH,
		CSS:     css,
		Blocks:  blocks,
		Fonts:   fonts,
	}
}

func defaultFontFaceCSS() string {
	data := base64.StdEncoding.EncodeToString(textoutline.DefaultFontBytes())
	return fmt.Sprintf(
		"@font-face{font-family:'%s';src:url(data:font/ttf;base64,%s) format('truetype');font-display:swap;}\n",
		textoutline.DefaultFontFamily,
		data,
	)
}

func fontMapForOutline(fontURLs []string, assets map[string]compileasset.Asset) map[string]textoutline.FontBytes {
	out := make(map[string]textoutline.FontBytes)
	out[textoutline.DefaultFontFamily] = textoutline.FontBytes{MIME: "font/ttf", Bytes: textoutline.DefaultFontBytes()}
	for i, u := range fontURLs {
		u = strings.TrimSpace(u)
		a, ok := assets[u]
		if !ok {
			continue
		}
		name := compileasset.FontFamilyNameForURL(u, i)
		out[name] = textoutline.FontBytes{MIME: a.MIME, Bytes: a.Bytes}
	}
	return out
}

func filterIDFromFragments(frags []styleadapter.StyleFragment) string {
	for _, f := range frags {
		if f.GetString(styleadapter.TypeKey) != styleadapter.TypeFilter {
			continue
		}
		if id := f.GetString("id"); id != "" {
			return id
		}
	}
	for _, f := range frags {
		if f.GetString(styleadapter.TypeKey) != styleadapter.TypeRect {
			continue
		}
		filter := f.GetString("filter")
		if strings.HasPrefix(filter, "url(#") && strings.HasSuffix(filter, ")") {
			return strings.TrimSuffix(strings.TrimPrefix(filter, "url(#"), ")")
		}
	}
	return ""
}

func transformFromFragments(frags []styleadapter.StyleFragment) string {
	for _, f := range frags {
		if f.GetString(styleadapter.TypeKey) != styleadapter.TypeG {
			continue
		}
		if tr := f.GetString("transform"); tr != "" && tr != "none" {
			return tr
		}
	}
	return ""
}

func outlinePageTexts(
	ctx context.Context,
	p *psrt.Page,
	pageSlug string,
	canvasW, canvasH int,
	fontURLs []string,
	assets map[string]compileasset.Asset,
) (map[int]textoutline.OutlinedBlock, bool, error) {
	in := buildPageOutlineInput(p, pageSlug, canvasW, canvasH, fontURLs, assets)
	if len(in.Blocks) == 0 {
		return map[int]textoutline.OutlinedBlock{}, false, nil
	}
	res, err := textoutline.Outline(ctx, in)
	if err != nil {
		return nil, false, fmt.Errorf("outline page %q: %w", p.Name, err)
	}
	out := make(map[int]textoutline.OutlinedBlock, len(res.Blocks))
	for _, b := range res.Blocks {
		out[b.Index] = b
	}
	return out, res.UsedGoTextFallback, nil
}

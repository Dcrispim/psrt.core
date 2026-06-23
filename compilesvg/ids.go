package compilesvg

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var slugSanitize = regexp.MustCompile(`[^a-zA-Z0-9_-]+`)

// Slug returns a filesystem- and XML-safe slug from a page name.
func Slug(name string) string {
	s := strings.ToLower(strings.TrimSpace(name))
	s = strings.ReplaceAll(s, " ", "-")
	s = slugSanitize.ReplaceAllString(s, "")
	s = strings.Trim(s, "-")
	if s == "" {
		return "page"
	}
	return s
}

// PageID returns the standardized page element id.
func PageID(slug string) string {
	return "psrt-page-" + slug
}

// PageClass returns the standardized page CSS class name (without dot).
func PageClass(slug string) string {
	return PageID(slug)
}

// PageBgID returns the background rect id.
func PageBgID(slug string) string {
	return PageID(slug) + "-bg"
}

// PageBgClass returns the background rect class name.
func PageBgClass(slug string) string {
	return PageBgID(slug)
}

// PageImageID returns the page background image id.
func PageImageID(slug string) string {
	return PageID(slug) + "-image"
}

// PageTextsID returns the texts group id.
func PageTextsID(slug string) string {
	return PageID(slug) + "-texts"
}

// TextID returns the standardized text block id.
func TextID(slug string, index int) string {
	return fmt.Sprintf("psrt-text-%s-%d", slug, index)
}

// TextClass returns the standardized text CSS class name.
func TextClass(slug string, index int) string {
	return TextID(slug, index)
}

// TextInnerClass is the inner span that holds inline markup inside a flex foreignObject div.
func TextInnerClass(slug string, index int) string {
	return TextID(slug, index) + "-inner"
}

// TextBgID returns the SVG background rect id for a text block.
func TextBgID(slug string, index int) string {
	return TextID(slug, index) + "-bg"
}

// TextBgClass returns the CSS class for the text background rect.
func TextBgClass(slug string, index int) string {
	return TextBgID(slug, index)
}

// TextWrapID returns the wrapper group id (rect + glyph paths).
func TextWrapID(slug string, index int) string {
	return TextID(slug, index) + "-wrap"
}

// TextGlyphsID returns the group id wrapping outlined text paths.
func TextGlyphsID(slug string, index int) string {
	return TextID(slug, index) + "-glyphs"
}

// TextA11yID returns the hidden accessibility text element id.
func TextA11yID(slug string, index int) string {
	return TextID(slug, index) + "-a11y"
}

// TextClassAttr returns both page and text classes for cascade (page first, text second).
func TextClassAttr(pageSlug string, index int) string {
	return PageClass(pageSlug) + " " + TextClass(pageSlug, index)
}

// UniqueSlugs assigns a unique slug per page name (suffix -2, -3 on collision).
func UniqueSlugs(pages []string) []string {
	seen := make(map[string]int)
	out := make([]string, len(pages))
	for i, name := range pages {
		base := Slug(name)
		n := seen[base]
		seen[base]++
		if n == 0 {
			out[i] = base
		} else {
			out[i] = base + "-" + strconv.Itoa(n+1)
		}
	}
	return out
}

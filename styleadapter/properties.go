package styleadapter

// Canonical property name constants (camelCase).
const (
	KeyPosition     = "position"
	KeyLeft         = "left"
	KeyTop          = "top"
	KeyWidth        = "width"
	KeyHeight       = "height"
	KeyBoxSizing    = "boxSizing"
	KeyPadding      = "padding"
	KeyPaddingTop   = "paddingTop"
	KeyPaddingRight = "paddingRight"
	KeyPaddingBottom = "paddingBottom"
	KeyPaddingLeft  = "paddingLeft"

	KeyBorder      = "border"
	KeyBorderTop   = "borderTop"
	KeyBorderRight = "borderRight"
	KeyBorderBottom = "borderBottom"
	KeyBorderLeft  = "borderLeft"
	KeyBorderWidth = "borderWidth"
	KeyBorderStyle = "borderStyle"
	KeyBorderColor = "borderColor"

	KeyBorderRadius            = "borderRadius"
	KeyBorderTopLeftRadius     = "borderTopLeftRadius"
	KeyBorderTopRightRadius    = "borderTopRightRadius"
	KeyBorderBottomRightRadius = "borderBottomRightRadius"
	KeyBorderBottomLeftRadius  = "borderBottomLeftRadius"

	KeyBoxShadow  = "boxShadow"
	KeyTextShadow = "textShadow"
	KeyGlow       = "glow"
	KeyBevel      = "bevel"
	KeyBlur       = "blur"
	KeyBlurLeft   = "blurLeft"
	KeyBlurRight  = "blurRight"
	KeyBlurTop    = "blurTop"
	KeyBlurBottom = "blurBottom"

	KeyBackground = "background"
	KeyColor      = "color"

	KeyTextAlign        = "textAlign"
	KeyAlignItems       = "alignItems"
	KeyTextDecoration   = "textDecoration"
	KeyTextDecorationLine = "textDecorationLine"
	KeyLetterSpacing    = "letterSpacing"
	KeyLineHeight       = "lineHeight"
	KeyWordSpacing      = "wordSpacing"
	KeyWhiteSpace       = "whiteSpace"
	KeyTextTransform    = "textTransform"
	KeyTextIndent       = "textIndent"
	KeyTextOverflow     = "textOverflow"
	KeyOpacity          = "opacity"

	KeyStroke      = "stroke"
	KeyStrokeWidth = "strokeWidth"
	KeyStrokeColor = "strokeColor"

	KeyTransform       = "transform"
	KeyTransformOrigin = "transformOrigin"
	KeyTranslate       = "translate"
	KeyRotate          = "rotate"
	KeyScale           = "scale"
	KeySkew            = "skew"
	KeyMatrix          = "matrix"

	KeyFontFamily = "fontFamily"
	KeyFontSize   = "fontSize"
	KeyFontWeight = "fontWeight"
	KeyFontStyle  = "fontStyle"
	KeyFontVariant = "fontVariant"
	KeyFontStretch = "fontStretch"
)

var boxKeys = map[string]struct{}{
	KeyHeight:     {},
	KeyBackground: {},
	KeyPadding:    {}, KeyPaddingTop: {}, KeyPaddingRight: {}, KeyPaddingBottom: {}, KeyPaddingLeft: {},
	KeyBorder: {}, KeyBorderTop: {}, KeyBorderRight: {}, KeyBorderBottom: {}, KeyBorderLeft: {},
	KeyBorderWidth: {}, KeyBorderStyle: {}, KeyBorderColor: {},
	KeyBorderRadius: {}, KeyBorderTopLeftRadius: {}, KeyBorderTopRightRadius: {},
	KeyBorderBottomRightRadius: {}, KeyBorderBottomLeftRadius: {},
	KeyBoxShadow: {}, KeyGlow: {}, KeyBevel: {},
	KeyBlur: {}, KeyBlurLeft: {}, KeyBlurRight: {}, KeyBlurTop: {}, KeyBlurBottom: {},
}

var textKeys = map[string]struct{}{
	KeyColor: {}, KeyTextAlign: {}, KeyAlignItems: {}, KeyTextDecoration: {}, KeyTextDecorationLine: {},
	KeyLetterSpacing: {}, KeyLineHeight: {}, KeyWordSpacing: {}, KeyWhiteSpace: {},
	KeyTextTransform: {}, KeyTextIndent: {}, KeyTextOverflow: {}, KeyOpacity: {},
	KeyTextShadow: {}, KeyStroke: {}, KeyStrokeWidth: {}, KeyStrokeColor: {},
	KeyFontFamily: {}, KeyFontSize: {}, KeyFontWeight: {}, KeyFontStyle: {},
	KeyFontVariant: {}, KeyFontStretch: {},
}

var transformKeys = map[string]struct{}{
	KeyTransform: {}, KeyTransformOrigin: {}, KeyTranslate: {}, KeyRotate: {},
	KeyScale: {}, KeySkew: {}, KeyMatrix: {},
}

func isBoxKey(k string) bool {
	_, ok := boxKeys[k]
	return ok
}

func isTextKey(k string) bool {
	_, ok := textKeys[k]
	return ok
}

func isTransformKey(k string) bool {
	_, ok := transformKeys[k]
	return ok
}

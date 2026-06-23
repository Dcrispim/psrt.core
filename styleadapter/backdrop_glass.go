package styleadapter

import (
	"fmt"
	"strconv"
	"strings"
)

const backdropGlassAlpha = 0.72

// postProcessBackdropGlass makes backdrop-filter blur visible through rounded edges.
// backdrop-filter blurs content *behind* the box; an opaque background hides that effect.
func postProcessBackdropGlass(box StyleFragment) {
	if box == nil {
		return
	}
	if box.GetString("backdropFilter") == "" && box.GetString("WebkitBackdropFilter") == "" {
		return
	}
	box.Set("overflow", "hidden")
	if bg := box.GetString("backgroundColor"); bg != "" {
		if softened, ok := softenOpaqueColorForBackdrop(bg, backdropGlassAlpha); ok {
			box.Set("backgroundColor", softened)
		}
	}
}

func softenOpaqueColorForBackdrop(color string, alpha float64) (string, bool) {
	color = strings.TrimSpace(color)
	if color == "" {
		return "", false
	}
	lower := strings.ToLower(color)
	if strings.HasPrefix(lower, "rgba(") {
		return softenRGBA(color, alpha)
	}
	if strings.HasPrefix(lower, "rgb(") {
		return softenRGB(color, alpha)
	}
	if strings.HasPrefix(color, "#") {
		return softenHexColor(color, alpha)
	}
	return "", false
}

func softenRGBA(color string, targetAlpha float64) (string, bool) {
	inner := strings.TrimSuffix(strings.TrimPrefix(color, "rgba("), ")")
	parts := strings.Split(inner, ",")
	if len(parts) != 4 {
		return "", false
	}
	a, err := parseAlphaToken(strings.TrimSpace(parts[3]))
	if err != nil || a < 0.99 {
		return "", false
	}
	r := strings.TrimSpace(parts[0])
	g := strings.TrimSpace(parts[1])
	b := strings.TrimSpace(parts[2])
	return fmt.Sprintf("rgba(%s,%s,%s,%.3f)", r, g, b, targetAlpha), true
}

func softenRGB(color string, targetAlpha float64) (string, bool) {
	inner := strings.TrimSuffix(strings.TrimPrefix(color, "rgb("), ")")
	parts := strings.Split(inner, ",")
	if len(parts) != 3 {
		return "", false
	}
	r := strings.TrimSpace(parts[0])
	g := strings.TrimSpace(parts[1])
	b := strings.TrimSpace(parts[2])
	return fmt.Sprintf("rgba(%s,%s,%s,%.3f)", r, g, b, targetAlpha), true
}

func softenHexColor(color string, targetAlpha float64) (string, bool) {
	hex := strings.TrimPrefix(color, "#")
	if len(hex) == 3 {
		hex = string([]byte{hex[0], hex[0], hex[1], hex[1], hex[2], hex[2]})
	}
	if len(hex) == 4 {
		// #rgba shorthand — already has alpha
		return "", false
	}
	if len(hex) == 8 {
		a, err := strconv.ParseUint(hex[6:8], 16, 8)
		if err != nil || a < 250 {
			return "", false
		}
		hex = hex[:6]
	}
	if len(hex) != 6 {
		return "", false
	}
	r, err1 := strconv.ParseUint(hex[0:2], 16, 8)
	g, err2 := strconv.ParseUint(hex[2:4], 16, 8)
	b, err3 := strconv.ParseUint(hex[4:6], 16, 8)
	if err1 != nil || err2 != nil || err3 != nil {
		return "", false
	}
	return fmt.Sprintf("rgba(%d,%d,%d,%.3f)", r, g, b, targetAlpha), true
}

func parseAlphaToken(s string) (float64, error) {
	s = strings.TrimSpace(s)
	if strings.HasSuffix(s, "%") {
		pct, err := strconv.ParseFloat(strings.TrimSuffix(s, "%"), 64)
		if err != nil {
			return 0, err
		}
		return pct / 100, nil
	}
	return strconv.ParseFloat(s, 64)
}

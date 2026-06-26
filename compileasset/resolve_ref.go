package compileasset

import (
	"strings"

	"github.com/Dcrispim/psrt.core/psrt"
)

// ResolveAssetReference expands @name@ placeholders in raw using consts, then trims.
func ResolveAssetReference(raw string, consts map[string]string) string {
	return strings.TrimSpace(psrt.ExpandConsts(raw, consts))
}

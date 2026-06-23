package psrt

import (
	"sort"
	"strings"
)

// ExpandConsts replaces @name@ placeholders in content using consts.
// Longer keys are substituted first so names that are prefixes of other names do not break.
func ExpandConsts(content string, consts map[string]string) string {
	if len(consts) == 0 {
		return content
	}
	keys := sortedStringKeys(consts)
	sort.Slice(keys, func(i, j int) bool {
		if len(keys[i]) != len(keys[j]) {
			return len(keys[i]) > len(keys[j])
		}
		return keys[i] < keys[j]
	})
	out := content
	for _, k := range keys {
		out = strings.ReplaceAll(out, "@"+k+"@", consts[k])
	}
	return out
}

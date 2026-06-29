package psrt

import (
	"sort"
	"strings"
)

// ConstsWithInteractive returns a substitution map combining plain consts with
// interactive ones flattened to their Render text (keyed by the `type:render`
// reference token). Feeding it to ExpandConsts collapses `@type:render@` to its
// base label in universal snapshots (SVG/HTML), dropping the interactive behaviour.
func ConstsWithInteractive(consts map[string]string, iConst map[string]InteractiveConst) map[string]string {
	if len(iConst) == 0 {
		return consts
	}
	merged := make(map[string]string, len(consts)+len(iConst))
	for k, v := range consts {
		merged[k] = v
	}
	for token, ic := range iConst {
		merged[token] = ic.Render
	}
	return merged
}

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

package compilehtml

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// DiscoverPSRTPaths returns sibling .psrt files in the same directory as primaryPath,
// with the primary file first (stable order for the rest).
func DiscoverPSRTPaths(primaryPath string) ([]string, error) {
	primaryPath = strings.TrimSpace(primaryPath)
	if primaryPath == "" {
		return nil, nil
	}
	primaryPath, err := filepath.Abs(primaryPath)
	if err != nil {
		return nil, err
	}
	dir := filepath.Dir(primaryPath)
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var others []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if !strings.HasSuffix(strings.ToLower(name), ".psrt") {
			continue
		}
		full := filepath.Join(dir, name)
		full, err := filepath.Abs(full)
		if err != nil {
			continue
		}
		if full == primaryPath {
			continue
		}
		others = append(others, full)
	}
	sort.Strings(others)
	out := make([]string, 0, 1+len(others))
	out = append(out, primaryPath)
	out = append(out, others...)
	return out, nil
}

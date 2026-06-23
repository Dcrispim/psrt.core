package cache

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

// LoadMap reads URL → relative path entries from assets.map.
func LoadMap(path string) (map[string]string, error) {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return make(map[string]string), nil
		}
		return nil, err
	}
	defer f.Close()

	out := make(map[string]string)
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.Split(line, "\t")
		if len(parts) < 2 {
			continue
		}
		url := strings.TrimSpace(parts[0])
		rel := strings.TrimSpace(parts[1])
		if url != "" && rel != "" {
			out[url] = rel
		}
	}
	return out, sc.Err()
}

// SaveMap writes assets.map with a header comment.
func SaveMap(path, psrtBase string, m map[string]string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	var b strings.Builder
	b.WriteString("# psrt: ")
	b.WriteString(psrtBase)
	b.WriteByte('\n')
	for url, rel := range m {
		b.WriteString(url)
		b.WriteByte('\t')
		b.WriteString(rel)
		b.WriteByte('\n')
	}
	return os.WriteFile(path, []byte(b.String()), 0o644)
}

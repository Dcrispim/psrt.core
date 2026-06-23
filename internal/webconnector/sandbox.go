package webconnector

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"

	"psrt/compileasset"
)

func ResolveWithinBase(baseDir, rawPath string) (string, error) {
	rawPath = strings.TrimSpace(rawPath)
	if rawPath == "" {
		return "", fmt.Errorf("empty path")
	}
	if !compileasset.IsLocalAssetRef(rawPath) {
		return "", fmt.Errorf("not a local asset path")
	}
	p, err := compileasset.ResolveAssetPath(rawPath)
	if err != nil {
		return "", err
	}
	if !filepath.IsAbs(p) {
		if baseDir == "" {
			return "", fmt.Errorf("relative path %q requires base_dir", rawPath)
		}
		p = filepath.Join(baseDir, p)
	}
	clean := filepath.Clean(p)
	resolved, err := filepath.EvalSymlinks(clean)
	if err != nil {
		return "", fmt.Errorf("resolve path: %w", err)
	}
	if !isUnderBase(baseDir, resolved) {
		return "", &SandboxError{Requested: resolved, Base: baseDir}
	}
	return resolved, nil
}

func isUnderBase(baseDir, target string) bool {
	base := filepath.Clean(baseDir)
	target = filepath.Clean(target)
	if runtime.GOOS == "windows" {
		base = strings.ToLower(base)
		target = strings.ToLower(target)
	}
	sep := string(filepath.Separator)
	if !strings.HasSuffix(base, sep) {
		base += sep
	}
	return strings.HasPrefix(target, base) || target == strings.TrimSuffix(base, sep)
}

type SandboxError struct {
	Requested string
	Base      string
}

func (e *SandboxError) Error() string {
	return "path outside shared base_dir"
}

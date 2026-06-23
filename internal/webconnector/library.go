package webconnector

import (
	"errors"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	errMissingPath = errors.New("missing path query")
	errNotPsrtFile = errors.New("path must be a .psrt file")
)

const libraryMaxDepth = 4

type LibraryProject struct {
	Path       string `json:"path"`
	Title      string `json:"title"`
	PageCount  int    `json:"pageCount"`
	ModifiedAt string `json:"modifiedAt"`
}

type LibraryProjectsResponse struct {
	Projects []LibraryProject `json:"projects"`
}

func (s *Server) handleLibraryProjects(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	base := s.BaseDir()
	projects, err := scanLibraryProjects(base, libraryMaxDepth)
	if err != nil {
		WriteErr(w, http.StatusInternalServerError, err)
		return
	}
	WriteJSON(w, http.StatusOK, LibraryProjectsResponse{Projects: projects})
}

func (s *Server) handleLibraryFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	rel := strings.TrimSpace(r.URL.Query().Get("path"))
	if rel == "" {
		WriteErr(w, http.StatusBadRequest, errMissingPath)
		return
	}
	if !strings.EqualFold(filepath.Ext(rel), ".psrt") {
		WriteErr(w, http.StatusBadRequest, errNotPsrtFile)
		return
	}
	abs, err := ResolveWithinBase(s.BaseDir(), rel)
	if err != nil {
		writeHandlerErr(w, s, err, r.RemoteAddr)
		return
	}
	data, err := os.ReadFile(abs)
	if err != nil {
		writeHandlerErr(w, s, err, r.RemoteAddr)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)
}

func scanLibraryProjects(baseDir string, maxDepth int) ([]LibraryProject, error) {
	var out []LibraryProject
	baseDir = filepath.Clean(baseDir)

	err := filepath.WalkDir(baseDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			if path != baseDir && shouldSkipLibraryDir(d.Name()) {
				return filepath.SkipDir
			}
			rel, relErr := filepath.Rel(baseDir, path)
			if relErr != nil {
				return nil
			}
			depth := len(strings.Split(filepath.ToSlash(rel), "/"))
			if depth > maxDepth {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.EqualFold(filepath.Ext(d.Name()), ".psrt") {
			return nil
		}
		rel, relErr := filepath.Rel(baseDir, path)
		if relErr != nil {
			return nil
		}
		rel = filepath.ToSlash(rel)
		info, statErr := d.Info()
		if statErr != nil {
			return nil
		}
		pages, _ := countPageStarts(path)
		out = append(out, LibraryProject{
			Path:       rel,
			Title:      strings.TrimSuffix(d.Name(), filepath.Ext(d.Name())),
			PageCount:  pages,
			ModifiedAt: info.ModTime().UTC().Format(time.RFC3339),
		})
		return nil
	})
	if err != nil {
		return nil, err
	}
	if out == nil {
		out = []LibraryProject{}
	}
	return out, nil
}

func shouldSkipLibraryDir(name string) bool {
	lower := strings.ToLower(name)
	if strings.HasPrefix(name, ".") {
		return true
	}
	switch lower {
	case "node_modules", "dist", "build", ".git":
		return true
	default:
		return false
	}
}

func countPageStarts(path string) (int, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return 0, err
	}
	count := 0
	for _, line := range strings.Split(string(data), "\n") {
		s := strings.TrimSpace(line)
		if strings.HasPrefix(strings.ToUpper(s), "$START") {
			after := strings.TrimSpace(s[1:])
			if hasLibraryKeyword(after, "START") {
				count++
			}
		}
	}
	return count, nil
}

func hasLibraryKeyword(after, kw string) bool {
	upper := strings.ToUpper(after)
	if !strings.HasPrefix(upper, kw) {
		return false
	}
	if len(after) == len(kw) {
		return true
	}
	if len(after) <= len(kw) {
		return false
	}
	return after[len(kw)] == ' ' || after[len(kw)] == '\t'
}

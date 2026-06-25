package cache

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/Dcrispim/psrt.core/compileasset"
	"github.com/Dcrispim/psrt.core/psrt"
)

func readAll(resp *http.Response) ([]byte, error) {
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

// Store manages on-disk cache for one PSRT file's remote assets.
type Store struct {
	RootDir  string
	MapPath  string
	ImagesDir string
	PsrtBase string
	PsrtDir  string
	urlMap   map[string]string
}

// NewStore creates or opens cache for psrtPath under configDir/cache/{basename}/.
func NewStore(configDir, psrtPath string) (*Store, error) {
	if configDir == "" {
		configDir = ConfigDir()
	}
	base := strings.TrimSuffix(filepath.Base(psrtPath), filepath.Ext(psrtPath))
	if base == "" {
		base = "document"
	}
	root := filepath.Join(configDir, "cache", base)
	s := &Store{
		RootDir:   root,
		MapPath:   filepath.Join(root, "assets.map"),
		ImagesDir: filepath.Join(root, "images"),
		PsrtBase:  base,
		PsrtDir:   filepath.Dir(psrtPath),
	}
	if err := os.MkdirAll(s.ImagesDir, 0o755); err != nil {
		return nil, err
	}
	m, err := LoadMap(s.MapPath)
	if err != nil {
		return nil, err
	}
	s.urlMap = m
	return s, nil
}

// ResolveAbsolute returns full path for a cached URL if file exists.
func (s *Store) ResolveAbsolute(url string) (string, bool) {
	rel, ok := s.urlMap[strings.TrimSpace(url)]
	if !ok {
		return "", false
	}
	abs := filepath.Join(s.RootDir, filepath.FromSlash(rel))
	if _, err := os.Stat(abs); err != nil {
		return "", false
	}
	return abs, true
}

func (s *Store) persistMap() error {
	return SaveMap(s.MapPath, s.PsrtBase+".psrt", s.urlMap)
}

// EnsureCached downloads or copies url if not cached locally.
func (s *Store) EnsureCached(ctx context.Context, client *http.Client, url, pageLabel string) error {
	url = strings.TrimSpace(url)
	if url == "" {
		return nil
	}
	if compileasset.IsLocalAssetRef(url) {
		return s.ensureLocalCached(url, pageLabel, false)
	}
	if !compileasset.LooksLikeHTTPURL(url) {
		return nil
	}
	if path, ok := s.ResolveAbsolute(url); ok {
		_ = path
		return nil
	}
	return s.downloadAndMap(ctx, client, url, pageLabel)
}

// RefreshAsset forces re-download or re-read from disk and updates the map.
func (s *Store) RefreshAsset(ctx context.Context, client *http.Client, url, pageLabel string) error {
	url = strings.TrimSpace(url)
	if url == "" {
		return fmt.Errorf("empty asset reference")
	}
	if compileasset.IsLocalAssetRef(url) {
		return s.ensureLocalCached(url, pageLabel, true)
	}
	if !compileasset.LooksLikeHTTPURL(url) {
		return fmt.Errorf("not an http(s) url or local path")
	}
	return s.downloadAndMap(ctx, client, url, pageLabel)
}

func (s *Store) ensureLocalCached(url, pageLabel string, force bool) error {
	if !force {
		if path, ok := s.ResolveAbsolute(url); ok {
			_ = path
			return nil
		}
	}
	src, err := compileasset.ResolveAssetPathRelative(url, s.PsrtDir)
	if err != nil {
		return err
	}
	if _, err := os.Stat(src); err != nil {
		return fmt.Errorf("local asset %q: %w", url, err)
	}
	asset, err := compileasset.ReadAssetFile(src)
	if err != nil {
		return err
	}
	name := assetFilename(pageLabel, s.PsrtBase, url, asset.MIME)
	abs := filepath.Join(s.ImagesDir, name)
	if err := os.WriteFile(abs, asset.Bytes, 0o644); err != nil {
		return err
	}
	rel := filepath.ToSlash(filepath.Join("images", name))
	s.urlMap[url] = rel
	return s.persistMap()
}

func (s *Store) downloadAndMap(ctx context.Context, client *http.Client, url, pageLabel string) error {
	if client == nil {
		client = http.DefaultClient
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		resp.Body.Close()
		return fmt.Errorf("fetch %q: status %s", url, resp.Status)
	}
	body, err := readAll(resp)
	if err != nil {
		return err
	}
	mime := resp.Header.Get("Content-Type")
	if mime == "" {
		mime = "application/octet-stream"
	}
	if idx := strings.Index(mime, ";"); idx >= 0 {
		mime = strings.TrimSpace(mime[:idx])
	}
	name := assetFilename(pageLabel, s.PsrtBase, url, mime)
	abs := filepath.Join(s.ImagesDir, name)
	if err := os.WriteFile(abs, body, 0o644); err != nil {
		return err
	}
	rel := filepath.ToSlash(filepath.Join("images", name))
	s.urlMap[url] = rel
	return s.persistMap()
}

// URLPageLabels builds url → page label for EnsureDocument.
func URLPageLabels(doc psrt.Document) map[string]string {
	labels := make(map[string]string)
	for i := range doc.Pages {
		p := &doc.Pages[i]
		img := compileasset.ResolveAssetReference(p.ImageURL, doc.Consts)
		if compileasset.IsAssetReference(img) {
			labels[img] = p.Name
		}
		for j := range p.Texts {
			ref := compileasset.ResolveAssetReference(p.Texts[j].ImageRef, doc.Consts)
			if compileasset.IsAssetReference(ref) {
				labels[ref] = fmt.Sprintf("%s_text_%d", p.Name, p.Texts[j].Index)
			}
		}
		for j := range p.Masks {
			ref := compileasset.ResolveAssetReference(p.Masks[j].ImageRef, doc.Consts)
			if compileasset.IsAssetReference(ref) {
				labels[ref] = fmt.Sprintf("%s_mask_%d", p.Name, p.Masks[j].Index)
			}
		}
	}
	for _, f := range doc.Fonts {
		font := compileasset.ResolveAssetReference(f, doc.Consts)
		if compileasset.IsAssetReference(font) {
			labels[font] = "font"
		}
	}
	return labels
}

// EnsureDocument warms cache for all HTTP URLs in doc.
func (s *Store) EnsureDocument(ctx context.Context, client *http.Client, doc psrt.Document) error {
	labels := URLPageLabels(doc)
	for url, label := range labels {
		if err := s.EnsureCached(ctx, client, url, label); err != nil {
			return fmt.Errorf("%s: %w", url, err)
		}
	}
	return nil
}

// ReadAsset loads bytes for url from cache or returns false.
func (s *Store) ReadAsset(url string) (compileasset.Asset, bool, error) {
	abs, ok := s.ResolveAbsolute(url)
	if !ok {
		return compileasset.Asset{}, false, nil
	}
	body, err := os.ReadFile(abs)
	if err != nil {
		return compileasset.Asset{}, false, err
	}
	mime := mimeFromPath(abs)
	return compileasset.Asset{Bytes: body, MIME: mime}, true, nil
}

func mimeFromPath(path string) string {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".webp":
		return "image/webp"
	case ".avif":
		return "image/avif"
	case ".svg":
		return "image/svg+xml"
	case ".woff2":
		return "font/woff2"
	case ".woff":
		return "font/woff"
	default:
		return "application/octet-stream"
	}
}

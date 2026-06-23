package webconnector

import (
	"net/http"
	"sync"
)

type Server struct {
	configPath string
	cfgMu      sync.RWMutex
	cfg        *Config
	auth       *Auth
	audit      *Audit
}

func NewServer(configPath string, cfg *Config, audit *Audit) *Server {
	return &Server{
		configPath: configPath,
		cfg:        cfg,
		auth:       NewAuth(audit),
		audit:      audit,
	}
}

func (s *Server) Auth() *Auth {
	return s.auth
}

func (s *Server) Config() *Config {
	s.cfgMu.RLock()
	defer s.cfgMu.RUnlock()
	return s.cfg.Clone()
}

func (s *Server) BaseDir() string {
	s.cfgMu.RLock()
	defer s.cfgMu.RUnlock()
	return s.cfg.BaseDir
}

func (s *Server) AllowedOrigin() string {
	s.cfgMu.RLock()
	defer s.cfgMu.RUnlock()
	return s.cfg.AllowedOrigin
}

func (s *Server) IsOriginAllowed(origin string) bool {
	if origin == "" {
		return false
	}
	return originsMatch(origin, s.AllowedOrigin())
}

func (s *Server) ConfigPath() string {
	return s.configPath
}

func (s *Server) ReloadConfig() error {
	cfg, err := LoadConfig(s.configPath)
	if err != nil {
		return err
	}
	s.cfgMu.RLock()
	oldPort := s.cfg.Port
	s.cfgMu.RUnlock()

	s.cfgMu.Lock()
	s.cfg = cfg
	s.cfgMu.Unlock()

	s.audit.ConfigReloaded(s.configPath, oldPort != cfg.Port)
	return nil
}

func (s *Server) UpdateConfig(next *Config) error {
	if err := next.Validate(); err != nil {
		return err
	}
	if err := next.NormalizeBaseDir(); err != nil {
		return err
	}
	if err := next.Save(s.configPath); err != nil {
		return err
	}
	s.cfgMu.Lock()
	s.cfg = next
	s.cfgMu.Unlock()
	return nil
}

func (s *Server) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /health", s.handleHealth)
	mux.HandleFunc("POST /pair", s.handlePair)
	mux.HandleFunc("GET /config", s.requireAuth(s.handleGetConfig))
	mux.HandleFunc("PUT /config", s.requireAuth(s.handlePutConfig))
	mux.HandleFunc("GET /image", s.requireAuth(s.handleImage))
	mux.HandleFunc("GET /library/projects", s.requireAuth(s.handleLibraryProjects))
	mux.HandleFunc("GET /library/file", s.requireAuth(s.handleLibraryFile))
	mux.HandleFunc("POST /api/adapt-entries-for-web", s.requireAuth(s.handleAdaptEntries))
	mux.HandleFunc("POST /api/get-asset-data-uri", s.requireAuth(s.handleAssetURI))
	mux.HandleFunc("POST /api/format-document-json", s.requireAuth(s.handleFormatJSON))
	mux.HandleFunc("POST /api/format-page-document-json", s.requireAuth(s.handleFormatPageJSON))
	mux.HandleFunc("POST /api/merge-page-document-psrt", s.requireAuth(s.handleMergePagePSRT))
	mux.HandleFunc("POST /api/parse-document-psrt", s.requireAuth(s.handleParsePSRT))
	mux.HandleFunc("POST /api/compile-page-svg-from-document", s.requireAuth(s.handleCompileSVG))
	mux.HandleFunc("POST /api/compile-page-html-from-document", s.requireAuth(s.handleCompileHTML))
	mux.HandleFunc("POST /api/refresh-asset-url", s.requireAuth(s.handleRefreshAsset))
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	s.RegisterRoutes(mux)
	return WithStrictCORS(s, mux)
}

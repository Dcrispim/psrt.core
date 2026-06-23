package webconnector

import (
	"fmt"
	"net/http"
	"path/filepath"
)

type configResponse struct {
	BaseDir             string `json:"base_dir"`
	AllowedOrigin       string `json:"allowed_origin"`
	Port                int    `json:"port"`
	ConfigPath          string `json:"config_path"`
	PortRestartRequired bool   `json:"port_restart_required"`
}

type configUpdateRequest struct {
	BaseDir       *string `json:"base_dir"`
	AllowedOrigin *string `json:"allowed_origin"`
	Port          *int    `json:"port"`
}

func (s *Server) configResponse() configResponse {
	cfg := s.Config()
	abs, _ := filepath.Abs(s.configPath)
	return configResponse{
		BaseDir:             cfg.BaseDir,
		AllowedOrigin:       cfg.AllowedOrigin,
		Port:                cfg.Port,
		ConfigPath:          abs,
		PortRestartRequired: false,
	}
}

func (s *Server) handleGetConfig(w http.ResponseWriter, _ *http.Request) {
	WriteJSON(w, http.StatusOK, s.configResponse())
}

func (s *Server) handlePutConfig(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req configUpdateRequest
	if !DecodeJSON(w, r, &req) {
		return
	}
	if req.Port != nil {
		WriteErr(w, http.StatusBadRequest, fmt.Errorf("alteração de porta exige reiniciar o conector"))
		return
	}
	cur := s.Config()
	updated := cur.Clone()
	var fields []string
	if req.BaseDir != nil {
		updated.BaseDir = *req.BaseDir
		fields = append(fields, "base_dir")
	}
	if req.AllowedOrigin != nil {
		updated.AllowedOrigin = *req.AllowedOrigin
		fields = append(fields, "allowed_origin")
	}
	if len(fields) == 0 {
		WriteErr(w, http.StatusBadRequest, fmt.Errorf("no fields to update"))
		return
	}
	if err := s.UpdateConfig(updated); err != nil {
		WriteErr(w, http.StatusBadRequest, err)
		return
	}
	s.audit.ConfigUpdated(fields)
	WriteJSON(w, http.StatusOK, s.configResponse())
}

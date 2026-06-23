package webconnector

import (
	"fmt"
	"net/http"
	"strings"
)

func (s *Server) handleHealth(w http.ResponseWriter, _ *http.Request) {
	WriteJSON(w, http.StatusOK, map[string]string{
		"status":  "ok",
		"version": Version,
	})
}

type pairRequest struct {
	Code string `json:"code"`
}

type pairResponse struct {
	Token string `json:"token"`
}

func (s *Server) handlePair(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req pairRequest
	if !DecodeJSON(w, r, &req) {
		return
	}
	origin := r.Header.Get("Origin")
	token, ok := s.auth.Pair(strings.TrimSpace(req.Code), origin, r.RemoteAddr)
	if !ok {
		WriteErr(w, http.StatusUnauthorized, fmt.Errorf("invalid or expired pairing code"))
		return
	}
	WriteJSON(w, http.StatusOK, pairResponse{Token: token})
}

package webconnector

import (
	"bytes"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"time"
)

func (s *Server) handleImage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	expanded := strings.TrimSpace(r.URL.Query().Get("path"))
	if expanded == "" {
		WriteErr(w, http.StatusBadRequest, fmt.Errorf("missing path query parameter"))
		return
	}
	absPath, err := ResolveWithinBase(s.BaseDir(), expanded)
	if err != nil {
		if _, ok := err.(*SandboxError); ok {
			s.audit.SandboxViolation(expanded, r.RemoteAddr)
			WriteErr(w, http.StatusForbidden, err)
			return
		}
		WriteErr(w, http.StatusBadRequest, err)
		return
	}
	asset, err := ReadAllowedImage(absPath)
	if err != nil {
		if me, ok := err.(*MIMEError); ok {
			s.audit.MimeRejected(absPath, me.MIME, r.RemoteAddr)
			WriteErr(w, http.StatusUnsupportedMediaType, err)
			return
		}
		WriteErr(w, http.StatusBadRequest, err)
		return
	}
	w.Header().Set("Content-Type", asset.MIME)
	http.ServeContent(w, r, filepath.Base(absPath), time.Time{}, bytes.NewReader(asset.Bytes))
}

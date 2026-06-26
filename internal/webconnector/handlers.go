package webconnector

import (
	"errors"
	"net/http"

	"github.com/Dcrispim/psrt.core/compileasset"
	"github.com/Dcrispim/psrt.core/internal/visualapp"
)

type AdaptEntriesReq struct {
	EntriesJSON string  `json:"entriesJSON"`
	CanvasW     int     `json:"canvasW"`
	CanvasH     int     `json:"canvasH"`
	Zoom        float64 `json:"zoom"`
}

type AssetURIReq struct {
	URL string `json:"url"`
}

type DocJSONReq struct {
	DocJSON string `json:"docJSON"`
}

type CompileReq struct {
	DocJSON  string `json:"docJSON"`
	PageName string `json:"pageName"`
}

type MergePageReq struct {
	FullDocJSON string `json:"fullDocJSON"`
	PageName    string `json:"pageName"`
	PsrtText    string `json:"psrtText"`
}

func (s *Server) handleAdaptEntries(w http.ResponseWriter, r *http.Request) {
	var req AdaptEntriesReq
	if !DecodeJSONMethod(w, r, http.MethodPost, &req) {
		return
	}
	out, err := visualapp.AdaptEntriesForWeb(req.EntriesJSON, req.CanvasW, req.CanvasH, req.Zoom)
	if err != nil {
		WriteErr(w, http.StatusBadRequest, err)
		return
	}
	WriteJSON(w, http.StatusOK, out)
}

func (s *Server) handleAssetURI(w http.ResponseWriter, r *http.Request) {
	var req AssetURIReq
	if !DecodeJSONMethod(w, r, http.MethodPost, &req) {
		return
	}
	uri, err := s.fetchAssetDataURI(req.URL, r.RemoteAddr)
	if err != nil {
		writeHandlerErr(w, s, err, r.RemoteAddr)
		return
	}
	WriteJSON(w, http.StatusOK, map[string]string{"uri": uri})
}

func (s *Server) fetchAssetDataURI(url, remote string) (string, error) {
	if url == "" {
		return "", nil
	}
	if compileasset.IsLocalAssetRef(url) {
		path, err := ResolveWithinBase(s.BaseDir(), url)
		if err != nil {
			return "", err
		}
		asset, err := ReadAllowedImage(path)
		if err != nil {
			return "", err
		}
		return compileasset.EncodeDataURI(asset.MIME, asset.Bytes), nil
	}
	fetched, err := compileasset.FetchURLs(http.DefaultClient, []string{url})
	if err != nil {
		return "", err
	}
	asset, ok := fetched[url]
	if !ok {
		return "", nil
	}
	return compileasset.EncodeDataURI(asset.MIME, asset.Bytes), nil
}

func writeHandlerErr(w http.ResponseWriter, s *Server, err error, remote string) {
	var sb *SandboxError
	var me *MIMEError
	switch {
	case errors.As(err, &sb):
		s.audit.SandboxViolation(sb.Requested, remote)
		WriteErr(w, http.StatusForbidden, err)
	case errors.As(err, &me):
		s.audit.MimeRejected("", me.MIME, remote)
		WriteErr(w, http.StatusUnsupportedMediaType, err)
	default:
		WriteErr(w, http.StatusBadRequest, err)
	}
}

func (s *Server) handleFormatJSON(w http.ResponseWriter, r *http.Request) {
	var req DocJSONReq
	if !DecodeJSONMethod(w, r, http.MethodPost, &req) {
		return
	}
	app := visualapp.New(nil)
	out, err := app.FormatDocumentJSON(req.DocJSON)
	if err != nil {
		WriteErr(w, http.StatusBadRequest, err)
		return
	}
	WriteJSON(w, http.StatusOK, map[string]string{"text": out})
}

func (s *Server) handleFormatPageJSON(w http.ResponseWriter, r *http.Request) {
	var req CompileReq
	if !DecodeJSONMethod(w, r, http.MethodPost, &req) {
		return
	}
	out, err := visualapp.FormatPageDocumentJSON(req.DocJSON, req.PageName)
	if err != nil {
		WriteErr(w, http.StatusBadRequest, err)
		return
	}
	WriteJSON(w, http.StatusOK, map[string]string{"text": out})
}

func (s *Server) handleMergePagePSRT(w http.ResponseWriter, r *http.Request) {
	var req MergePageReq
	if !DecodeJSONMethod(w, r, http.MethodPost, &req) {
		return
	}
	out, err := visualapp.MergePageDocumentPSRT(req.FullDocJSON, req.PageName, req.PsrtText)
	if err != nil {
		WriteErr(w, http.StatusBadRequest, err)
		return
	}
	WriteJSON(w, http.StatusOK, map[string]string{"document": out})
}

func (s *Server) handleParsePSRT(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Text string `json:"text"`
	}
	if !DecodeJSONMethod(w, r, http.MethodPost, &body) {
		return
	}
	app := visualapp.New(nil)
	out, err := app.ParseDocumentPSRT(body.Text)
	if err != nil {
		WriteErr(w, http.StatusBadRequest, err)
		return
	}
	WriteJSON(w, http.StatusOK, map[string]string{"document": out})
}

func (s *Server) handleCompileSVG(w http.ResponseWriter, r *http.Request) {
	var req CompileReq
	if !DecodeJSONMethod(w, r, http.MethodPost, &req) {
		return
	}
	app := visualapp.New(nil)
	res, err := app.CompilePageSVGFromDocument(req.DocJSON, req.PageName)
	if err != nil {
		WriteErr(w, http.StatusBadRequest, err)
		return
	}
	WriteJSON(w, http.StatusOK, map[string]any{
		"uri":                res.URI,
		"usedGoTextFallback": res.UsedGoTextFallback,
	})
}

func (s *Server) handleCompileHTML(w http.ResponseWriter, r *http.Request) {
	var req CompileReq
	if !DecodeJSONMethod(w, r, http.MethodPost, &req) {
		return
	}
	app := visualapp.New(nil)
	uri, err := app.CompilePageHTMLFromDocument(req.DocJSON, req.PageName)
	if err != nil {
		WriteErr(w, http.StatusBadRequest, err)
		return
	}
	WriteJSON(w, http.StatusOK, map[string]string{"uri": uri})
}

func (s *Server) handleRefreshAsset(w http.ResponseWriter, _ *http.Request) {
	WriteJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

// psrt-gui-dev is a small HTTP API for browser-only frontend development (Vite mode "web").
// Run: go run ./cmd/psrt-gui-dev
// Then: cd cmd/psrt-gui/frontend && npm run dev:web
package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"psrt/compileasset"
	"psrt/internal/visualapp"
)

func main() {
	addr := envOr("PSRT_DEV_API_ADDR", ":8787")
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", healthOK)
	mux.HandleFunc("GET /api/health", healthOK)
	mux.HandleFunc("POST /api/adapt-entries-for-web", handleAdaptEntries)
	mux.HandleFunc("POST /api/get-asset-data-uri", handleAssetURI)
	mux.HandleFunc("POST /api/format-document-json", handleFormatJSON)
	mux.HandleFunc("POST /api/format-page-document-json", handleFormatPageJSON)
	mux.HandleFunc("POST /api/merge-page-document-psrt", handleMergePagePSRT)
	mux.HandleFunc("POST /api/parse-document-psrt", handleParsePSRT)
	mux.HandleFunc("POST /api/compile-page-svg-from-document", handleCompileSVG)
	mux.HandleFunc("POST /api/compile-page-html-from-document", handleCompileHTML)
	mux.HandleFunc("POST /api/refresh-asset-url", handleRefreshAsset)

	log.Printf("PSRT dev API listening on http://localhost%s", addr)
	log.Fatal(http.ListenAndServe(addr, withCORS(mux)))
}

func healthOK(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"ok": "true"})
}

func envOr(key, fallback string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return fallback
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

type adaptEntriesReq struct {
	EntriesJSON string  `json:"entriesJSON"`
	CanvasW     int     `json:"canvasW"`
	CanvasH     int     `json:"canvasH"`
	Zoom        float64 `json:"zoom"`
}

func handleAdaptEntries(w http.ResponseWriter, r *http.Request) {
	var req adaptEntriesReq
	if !decodeJSON(w, r, &req) {
		return
	}
	out, err := visualapp.AdaptEntriesForWeb(req.EntriesJSON, req.CanvasW, req.CanvasH, req.Zoom)
	if err != nil {
		writeErr(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, out)
}

type assetURIReq struct {
	URL     string            `json:"url"`
	BaseDir string            `json:"baseDir"`
	Consts  map[string]string `json:"consts"`
}

func handleAssetURI(w http.ResponseWriter, r *http.Request) {
	var req assetURIReq
	if !decodeJSON(w, r, &req) {
		return
	}
	url := compileasset.ResolveAssetReference(req.URL, req.Consts)
	uri, err := fetchAssetDataURI(url, req.BaseDir)
	if err != nil {
		writeErr(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"uri": uri})
}

func fetchAssetDataURI(url, baseDir string) (string, error) {
	if url == "" {
		return "", nil
	}
	if compileasset.IsLocalAssetRef(url) {
		path, err := compileasset.ResolveAssetPathRelative(url, baseDir)
		if err != nil {
			return "", err
		}
		asset, err := compileasset.ReadAssetFile(path)
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

type docJSONReq struct {
	DocJSON string `json:"docJSON"`
}

type compileReq struct {
	DocJSON  string `json:"docJSON"`
	PageName string `json:"pageName"`
}

type mergePageReq struct {
	FullDocJSON string `json:"fullDocJSON"`
	PageName    string `json:"pageName"`
	PsrtText    string `json:"psrtText"`
}

func handleFormatJSON(w http.ResponseWriter, r *http.Request) {
	var req docJSONReq
	if !decodeJSON(w, r, &req) {
		return
	}
	app := visualapp.New(nil)
	out, err := app.FormatDocumentJSON(req.DocJSON)
	if err != nil {
		writeErr(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"text": out})
}

func handleFormatPageJSON(w http.ResponseWriter, r *http.Request) {
	var req compileReq
	if !decodeJSON(w, r, &req) {
		return
	}
	out, err := visualapp.FormatPageDocumentJSON(req.DocJSON, req.PageName)
	if err != nil {
		writeErr(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"text": out})
}

func handleMergePagePSRT(w http.ResponseWriter, r *http.Request) {
	var req mergePageReq
	if !decodeJSON(w, r, &req) {
		return
	}
	out, err := visualapp.MergePageDocumentPSRT(req.FullDocJSON, req.PageName, req.PsrtText)
	if err != nil {
		writeErr(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"document": out})
}

func handleParsePSRT(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Text string `json:"text"`
	}
	if !decodeJSON(w, r, &body) {
		return
	}
	app := visualapp.New(nil)
	out, err := app.ParseDocumentPSRT(body.Text)
	if err != nil {
		writeErr(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"document": out})
}

func handleCompileSVG(w http.ResponseWriter, r *http.Request) {
	var req compileReq
	if !decodeJSON(w, r, &req) {
		return
	}
	app := visualapp.New(nil)
	res, err := app.CompilePageSVGFromDocument(req.DocJSON, req.PageName)
	if err != nil {
		writeErr(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"uri":                res.URI,
		"usedGoTextFallback": res.UsedGoTextFallback,
	})
}

func handleCompileHTML(w http.ResponseWriter, r *http.Request) {
	var req compileReq
	if !decodeJSON(w, r, &req) {
		return
	}
	app := visualapp.New(nil)
	uri, err := app.CompilePageHTMLFromDocument(req.DocJSON, req.PageName)
	if err != nil {
		writeErr(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"uri": uri})
}

func handleRefreshAsset(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func decodeJSON(w http.ResponseWriter, r *http.Request, dst any) bool {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return false
	}
	defer r.Body.Close()
	if err := json.NewDecoder(io.LimitReader(r.Body, 32<<20)).Decode(dst); err != nil {
		writeErr(w, http.StatusBadRequest, err)
		return false
	}
	return true
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeErr(w http.ResponseWriter, status int, err error) {
	writeJSON(w, status, map[string]string{"error": err.Error()})
}

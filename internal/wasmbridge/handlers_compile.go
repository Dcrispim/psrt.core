//go:build js

package wasmbridge

import (
	"context"
	"net/http"
	"strings"
	"syscall/js"
	"time"

	"github.com/Dcrispim/psrt.core/compilehtml"
	"github.com/Dcrispim/psrt.core/compileopts"
	"github.com/Dcrispim/psrt.core/compilesvg"
	"github.com/Dcrispim/psrt.core/psrt"
)

func loadDocFlexible(args []js.Value) (psrt.Document, error) {
	b, err := bytesArg(args, 0)
	if err != nil {
		return psrt.Document{}, err
	}
	if len(b) == 0 {
		return psrt.Document{}, errMissing("doc")
	}
	if looksLikePSRT(b) {
		return psrt.ParseString(string(b))
	}
	return psrt.ParseJSON(b)
}

func looksLikePSRT(b []byte) bool {
	s := strings.TrimSpace(string(b))
	return strings.HasPrefix(s, "$START") || strings.HasPrefix(s, "$FONTS") || strings.HasPrefix(s, "$CONSTS") || strings.HasPrefix(s, "$SOURCE")
}

func HandleCompileToHtml() js.Func {
	return wrap(func(args []js.Value) ([]byte, error) {
		doc, err := loadDocFlexible(args)
		if err != nil {
			return nil, err
		}
		doc = compilesvg.ResolveDocument(doc)
		jsOpts := compileOptsFromArg(args, 1)
		opts := compileopts.Options{LinksOnly: jsOpts.LinksOnly, NoScript: jsOpts.NoScript}
		client := &http.Client{Timeout: 120 * time.Second}
		return compilehtml.CompileWithOptions(context.Background(), doc, client, nil, opts)
	})
}

func HandleCompileToSvg() js.Func {
	return wrap(func(args []js.Value) ([]byte, error) {
		doc, err := loadDocFlexible(args)
		if err != nil {
			return nil, err
		}
		pageName, err := stringArg(args, 1)
		if err != nil {
			return nil, err
		}
		jsOpts := compileOptsFromArg(args, 2)
		opts := compileopts.Options{LinksOnly: jsOpts.LinksOnly, NoScript: jsOpts.NoScript}
		client := &http.Client{Timeout: 120 * time.Second}
		res, err := compilesvg.CompilePageSVGWithOptions(context.Background(), doc, pageName, client, nil, opts)
		if err != nil {
			return nil, err
		}
		return res.Data, nil
	})
}

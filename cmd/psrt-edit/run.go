package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/Dcrispim/psrt.core/compileasset/cache"
	"github.com/Dcrispim/psrt.core/compilehtml"
	"github.com/Dcrispim/psrt.core/compileopts"
	"github.com/Dcrispim/psrt.core/compilesvg"
	"github.com/Dcrispim/psrt.core/psrt"
	"github.com/Dcrispim/psrt.core/psrt/editor"
)

var (
	cat             bool
	compileSVG      bool
	compileHTML     bool
	compileSVGDir   string
	compileHTMLOut  string
	compileTimeout  time.Duration
)

func initRunFlags(root *cobra.Command) {
	f := root.PersistentFlags()
	f.BoolVar(&cat, "cat", false, "print resulting PSRT to stdout instead of saving")
	f.BoolVar(&compileSVG, "compile-svg", false, "compile resulting PSRT to SVG files after edit")
	f.BoolVar(&compileHTML, "compile-html", false, "compile resulting PSRT to HTML after edit")
	f.StringVar(&compileSVGDir, "compile-svg-dir", "", "output directory for SVG (--compile-svg; default: <input-dir>/out-svg)")
	f.StringVar(&compileHTMLOut, "compile-html-out", "", "output HTML path (--compile-html; default: <input>.html)")
	f.DurationVar(&compileTimeout, "compile-timeout", 30*time.Second, "HTTP timeout for compile steps")
}

func runEdit(mutate func(*psrt.Document) error) error {
	doc, err := editor.LoadDocument(inputPath)
	if err != nil {
		return err
	}
	if err := mutate(&doc); err != nil {
		return err
	}

	data, err := editor.FormatDocument(&doc)
	if err != nil {
		return err
	}

	if cat {
		if _, err := os.Stdout.Write(data); err != nil {
			return err
		}
	} else {
		out := outputPath
		if out == "" {
			out = inputPath
		}
		if err := os.WriteFile(out, data, 0o644); err != nil {
			return err
		}
	}

	if !compileSVG && !compileHTML {
		return nil
	}

	client := &http.Client{Timeout: compileTimeout}
	ctx := context.Background()
	var store *cache.Store
	if compileSVG || compileHTML {
		store, err = cache.NewStore("", inputPath)
		if err != nil {
			return fmt.Errorf("cache: %w", err)
		}
		_ = store.EnsureDocument(ctx, client, doc)
	}
	if compileSVG {
		dir := compileSVGDir
		if dir == "" {
			dir = defaultCompileSVGDir(inputPath)
		}
		if batch, err := compilesvg.CompileWithCache(ctx, doc, client, dir, store); err != nil {
			return fmt.Errorf("compile-svg: %w", err)
		} else if batch.UsedGoTextFallback {
			fmt.Fprintf(os.Stderr, "%s\n", compilesvg.GoTextFallbackNotice)
		}
	}
	if compileHTML {
		html, err := compilehtml.CompileWithCacheFrom(ctx, doc, inputPath, nil, nil, client, store, compileopts.Options{})
		if err != nil {
			return fmt.Errorf("compile-html: %w", err)
		}
		out := compileHTMLOut
		if out == "" {
			out = defaultCompileHTMLOut(inputPath)
		}
		if err := os.WriteFile(out, html, 0o644); err != nil {
			return fmt.Errorf("compile-html: %w", err)
		}
	}
	return nil
}

func defaultCompileSVGDir(input string) string {
	return filepath.Join(filepath.Dir(input), "out-svg")
}

func defaultCompileHTMLOut(input string) string {
	base := strings.TrimSuffix(filepath.Base(input), filepath.Ext(input))
	return filepath.Join(filepath.Dir(input), base+".html")
}

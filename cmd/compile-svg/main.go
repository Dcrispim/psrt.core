// Command psrt-compile-svg turns a PSRT file into standalone SVG files (one per page).
package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"psrt/compilesvg"
	"psrt/compileopts"
	"psrt/psrt"
)

const version = "0.1.0"

var (
	input     string
	outputDir string
	timeout   time.Duration
	linksOnly bool
)

var rootCmd = &cobra.Command{
	Use:   "psrt-compile-svg",
	Short: "Compile a PSRT file into standalone SVG files (one per page, embedded assets).",
	Long: `Reads a PSRT document, resolves constants, downloads HTTP(S) assets (images, fonts),
embeds them as data URIs, and writes one self-contained SVG per page into --output-dir.`,
	Version:       version,
	SilenceUsage:  true,
	SilenceErrors: true,
	Args:          cobra.NoArgs,
	RunE:          runRoot,
	Example: `  psrt-compile-svg --input=doc.psrt --output-dir=./out-svg
  psrt-compile-svg --input=doc.psrt --output-dir=./out-svg --timeout=60s`,
}

func init() {
	f := rootCmd.Flags()
	f.StringVar(&input, "input", "", "input .psrt path (use - for stdin)")
	f.StringVar(&outputDir, "output-dir", "", "output directory for .svg files")
	f.DurationVar(&timeout, "timeout", 30*time.Second, "timeout per HTTP request")
	f.BoolVar(&linksOnly, "links-only", false, "keep original asset URLs instead of embedding data URIs")
	_ = rootCmd.MarkFlagRequired("input")
	_ = rootCmd.MarkFlagRequired("output-dir")
	rootCmd.SetVersionTemplate("{{.Version}}\n")
}

func main() {
	rootCmd.SetOut(os.Stdout)
	rootCmd.SetErr(os.Stderr)
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "psrt-compile-svg: %v\n", err)
		os.Exit(1)
	}
}

func runRoot(_ *cobra.Command, _ []string) error {
	raw, err := readInput(input)
	if err != nil {
		return err
	}
	doc, err := psrt.Parse(strings.NewReader(string(raw)))
	if err != nil {
		return err
	}
	client := &http.Client{Timeout: timeout}
	opts := compileopts.Options{LinksOnly: linksOnly}
	res, err := compilesvg.CompileWithOptions(context.Background(), doc, client, strings.TrimSpace(outputDir), nil, opts)
	if err != nil {
		return err
	}
	if res.UsedGoTextFallback {
		fmt.Fprintln(os.Stderr, compilesvg.GoTextFallbackNotice)
	}
	return nil
}

func readInput(path string) ([]byte, error) {
	path = strings.TrimSpace(path)
	if path == "" {
		return nil, fmt.Errorf("--input is required")
	}
	if path == "-" {
		return io.ReadAll(os.Stdin)
	}
	return os.ReadFile(path)
}

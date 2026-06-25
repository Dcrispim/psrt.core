// Command psrt-compile turns a PSRT file into a single self-contained offline HTML page.
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
	"github.com/Dcrispim/psrt.core/compilehtml"
	"github.com/Dcrispim/psrt.core/compileopts"
	"github.com/Dcrispim/psrt.core/psrt"
)

const version = "0.1.0"

var (
	inputs  []string
	output  string
	timeout   time.Duration
	linksOnly bool
	noScript  bool
)

var rootCmd = &cobra.Command{
	Use:   "psrt-compile",
	Short: "Compile a PSRT file into a standalone HTML document (embedded assets).",
	Long: `Reads a PSRT document, downloads all HTTP(S) assets (page images, optional text images, fonts),
embeds them as data URIs, and writes one HTML file suitable for offline viewing.`,
	Version:       version,
	SilenceUsage:  true,
	SilenceErrors: true,
	Args:          cobra.NoArgs,
	RunE:          runRoot,
	Example: `  psrt-compile --input=doc.psrt --output=out.html
  psrt-compile --input=base.psrt --input=alt.psrt --output=out.html
  psrt-compile --input=doc.psrt --output=- --timeout=60s`,
}

func init() {
	f := rootCmd.Flags()
	f.StringSliceVar(&inputs, "input", nil, "one or more input .psrt paths (use - for stdin on the first only)")
	f.StringVar(&output, "output", "", "output .html path (- or empty for stdout)")
	f.DurationVar(&timeout, "timeout", 30*time.Second, "timeout per HTTP request")
	f.BoolVar(&linksOnly, "links-only", false, "keep original asset URLs instead of embedding data URIs")
	f.BoolVar(&noScript, "no-script", false, "omit variant switcher script for cleaner HTML")
	_ = rootCmd.MarkFlagRequired("input")
	rootCmd.SetVersionTemplate("{{.Version}}\n")
}

func main() {
	rootCmd.SetOut(os.Stdout)
	rootCmd.SetErr(os.Stderr)
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "psrt-compile: %v\n", err)
		os.Exit(1)
	}
}

func runRoot(_ *cobra.Command, _ []string) error {
	if len(inputs) == 0 {
		return fmt.Errorf("--input is required")
	}
	primary := strings.TrimSpace(inputs[0])
	raw, err := readInput(primary)
	if err != nil {
		return err
	}
	doc, err := psrt.Parse(strings.NewReader(string(raw)))
	if err != nil {
		return err
	}
	var more []string
	if len(inputs) > 1 {
		more = inputs[1:]
	}
	client := &http.Client{Timeout: timeout}
	opts := compileopts.Options{LinksOnly: linksOnly, NoScript: noScript}
	out, err := compilehtml.CompileWithCacheFrom(context.Background(), doc, primary, more, nil, client, nil, opts)
	if err != nil {
		return err
	}
	return writeOutput(output, out)
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

func writeOutput(path string, data []byte) error {
	if path == "" || path == "-" {
		_, err := os.Stdout.Write(data)
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

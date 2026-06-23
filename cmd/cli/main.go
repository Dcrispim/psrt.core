// Command psrt manipulates PSRT files from the command line (Cobra).
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"psrt/psrt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

const version = "0.1.0"

var (
	from      string
	to        string
	input     string
	output    string
	page      string
	textIdx   string
	constName string
)

var rootCmd = &cobra.Command{
	Use:   "psrt",
	Short: "Convert and slice PSRT documents (JSON / PSRT / Markdown).",
	Long: `Read a PSRT or JSON document, then write JSON, PSRT, or Markdown.

Optional filters export a single page, a text block (with --page and --text),
or a named constant. --const cannot be combined with --page or --text;
--text requires --page.`,
	Version:       version,
	SilenceUsage:  true,
	SilenceErrors: true,
	Args:          cobra.NoArgs,
	RunE:          runRoot,
	Example: `  psrt --input=doc.psrt --output=out.json
  psrt --input=doc.psrt --to=md --output=out.md
  psrt --input=doc.json --from=json --to=psrt --output=out.psrt
  psrt --input=doc.psrt --page=my-page --to=json
  psrt --input=doc.psrt --page=my-page --text=0 --to=psrt
  psrt --input=doc.psrt --const=shadow --to=md`,
}

func init() {
	f := rootCmd.Flags()
	f.StringVar(&from, "from", "psrt", "input format: psrt or json")
	f.StringVar(&to, "to", "json", "output format: json, psrt, or md")
	f.StringVar(&input, "input", "", "input file path (use - for stdin)")
	f.StringVar(&output, "output", "", "output file path (- or empty for stdout)")
	f.StringVar(&page, "page", "", "export only this page name")
	f.StringVar(&textIdx, "text", "", "with --page, export the text block whose Index matches")
	f.StringVar(&constName, "const", "", "export only this named constant")
	_ = rootCmd.MarkFlagRequired("input")

	rootCmd.SetVersionTemplate("{{.Version}}\n")
}

func main() {
	rootCmd.SetOut(os.Stdout)
	rootCmd.SetErr(os.Stderr)
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "psrt: %v\n", err)
		os.Exit(1)
	}
}

func runRoot(_ *cobra.Command, _ []string) error {
	cfg := config{
		from: from, to: to, input: input, output: output,
		page: page, constName: constName, textIdx: textIdx,
	}
	return run(cfg)
}

type config struct {
	from, to, input, output  string
	page, constName, textIdx string
}

func run(c config) error {
	doc, err := loadDocument(c.from, c.input)
	if err != nil {
		return err
	}

	payload, err := selectPayload(&doc, c)
	if err != nil {
		return err
	}

	out, err := encodePayload(payload, c.to)
	if err != nil {
		return err
	}

	return writeOutput(c.output, out)
}

func loadDocument(from, input string) (psrt.Document, error) {
	r, closer, err := openInput(input)
	if err != nil {
		return psrt.Document{}, err
	}
	if closer != nil {
		defer closer.Close()
	}
	raw, err := io.ReadAll(r)
	if err != nil {
		return psrt.Document{}, err
	}
	switch strings.ToLower(strings.TrimSpace(from)) {
	case "psrt":
		return psrt.Parse(strings.NewReader(string(raw)))
	case "json":
		return psrt.ParseJSON(raw)
	default:
		return psrt.Document{}, fmt.Errorf("unsupported --from=%q (use psrt or json)", from)
	}
}

func openInput(path string) (io.Reader, io.Closer, error) {
	if path == "-" {
		return os.Stdin, nil, nil
	}
	f, err := os.Open(path)
	if err != nil {
		return nil, nil, err
	}
	return f, f, nil
}

type payloadKind int

const (
	payloadDoc payloadKind = iota
	payloadPage
	payloadText
	payloadConst
)

type payload struct {
	kind  payloadKind
	doc   *psrt.Document
	page  *psrt.Page
	text  *psrt.Text
	cname string
	cval  string
}

func selectPayload(doc *psrt.Document, c config) (payload, error) {
	if strings.TrimSpace(c.constName) != "" {
		if strings.TrimSpace(c.page) != "" || strings.TrimSpace(c.textIdx) != "" {
			return payload{}, fmt.Errorf("--const cannot be combined with --page or --text")
		}
		v, ok := doc.Consts[c.constName]
		if !ok {
			return payload{}, fmt.Errorf("unknown constant %q", c.constName)
		}
		return payload{kind: payloadConst, cname: c.constName, cval: v}, nil
	}
	if strings.TrimSpace(c.textIdx) != "" {
		if strings.TrimSpace(c.page) == "" {
			return payload{}, fmt.Errorf("--text requires --page")
		}
	}
	if strings.TrimSpace(c.page) != "" {
		p := findPage(doc, c.page)
		if p == nil {
			return payload{}, fmt.Errorf("unknown page %q", c.page)
		}
		if strings.TrimSpace(c.textIdx) == "" {
			return payload{kind: payloadPage, doc: doc, page: p}, nil
		}
		idx, err := strconv.Atoi(strings.TrimSpace(c.textIdx))
		if err != nil {
			return payload{}, fmt.Errorf("--text: invalid index %q: %w", c.textIdx, err)
		}
		t := findTextByIndex(p, idx)
		if t == nil {
			return payload{}, fmt.Errorf("page %q has no text with index %d", c.page, idx)
		}
		return payload{kind: payloadText, doc: doc, page: p, text: t}, nil
	}
	return payload{kind: payloadDoc, doc: doc}, nil
}

func findPage(doc *psrt.Document, name string) *psrt.Page {
	for i := range doc.Pages {
		if doc.Pages[i].Name == name {
			return &doc.Pages[i]
		}
	}
	return nil
}

func findTextByIndex(p *psrt.Page, index int) *psrt.Text {
	for i := range p.Texts {
		if p.Texts[i].Index == index {
			return &p.Texts[i]
		}
	}
	return nil
}

func encodePayload(p payload, to string) ([]byte, error) {
	to = strings.ToLower(strings.TrimSpace(to))
	switch to {
	case "json":
		return encodeJSON(p)
	case "psrt":
		return encodePSRT(p)
	case "md", "markdown":
		return []byte(encodeMD(p)), nil
	default:
		return nil, fmt.Errorf("unsupported --to=%q (use json, psrt, or md)", to)
	}
}

func encodeJSON(p payload) ([]byte, error) {
	switch p.kind {
	case payloadDoc:
		return psrt.ToJSON(*p.doc)
	case payloadPage:
		return json.MarshalIndent(p.page, "", "  ")
	case payloadText:
		return json.MarshalIndent(p.text, "", "  ")
	case payloadConst:
		v := struct {
			Name  string `json:"name"`
			Value string `json:"value"`
		}{Name: p.cname, Value: p.cval}
		return json.MarshalIndent(v, "", "  ")
	default:
		return nil, fmt.Errorf("internal: unknown payload")
	}
}

func encodePSRT(p payload) ([]byte, error) {
	switch p.kind {
	case payloadDoc:
		return psrt.FormatPSRT(*p.doc, false)
	case payloadPage:
		return psrt.FormatPagePSRT(p.page)
	case payloadText:
		return psrt.FormatTextPSRT(p.text)
	case payloadConst:
		var b strings.Builder
		b.WriteString("$CONSTS\n")
		b.Write(psrt.FormatConstPSRT(p.cname, p.cval))
		b.WriteString("$ENDCONSTS\n")
		return []byte(b.String()), nil
	default:
		return nil, fmt.Errorf("internal: unknown payload")
	}
}

func encodeMD(p payload) string {
	switch p.kind {
	case payloadDoc:
		return psrt.FormatDocumentMarkdown(*p.doc)
	case payloadPage:
		return psrt.FormatPageMarkdown(p.page)
	case payloadText:
		return psrt.FormatTextMarkdown(p.text)
	case payloadConst:
		return psrt.FormatConstMarkdown(p.cname, p.cval)
	default:
		return ""
	}
}

func writeOutput(path string, data []byte) error {
	if path == "" || path == "-" {
		_, err := os.Stdout.Write(data)
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

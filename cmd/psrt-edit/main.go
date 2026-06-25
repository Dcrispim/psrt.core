// Command psrt-edit mutates PSRT documents from the command line.
package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/Dcrispim/psrt.core/psrt"
	"github.com/Dcrispim/psrt.core/psrt/editor"

	"github.com/spf13/cobra"
)

var (
	inputPath  string
	outputPath string
	pageName   string
	textIndex  string
)

func main() {
	root := &cobra.Command{
		Use:           "psrt-edit",
		Short:         "Edit PSRT documents (pages, texts, fonts, constants).",
		Long:          rootLong,
		Example:       rootExample,
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	root.PersistentFlags().StringVar(&inputPath, "input", "", "input PSRT file path (required)")
	root.PersistentFlags().StringVar(&outputPath, "output", "", "output file path (default: overwrite --input)")
	initRunFlags(root)
	_ = root.MarkPersistentFlagRequired("input")

	root.AddCommand(pageCmd(), textCmd(), maskCmd(), fontCmd(), constCmd())

	if err := root.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "psrt-edit: %v\n", err)
		os.Exit(1)
	}
}

func pageCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "page",
		Short:   "Edit pages (--page required on every subcommand)",
		Long:    pageLong,
		Example: pageExample,
	}

	var newName, path, styleKey, styleValue, styleJSON, beforePage, afterPage string

	rename := &cobra.Command{
		Use:     "rename",
		Short:   "Rename a page",
		Example: `  psrt-edit --input=doc.psrt page rename --page=capa --name=cover`,
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := requirePage(); err != nil {
				return err
			}
			if strings.TrimSpace(newName) == "" {
				return fmt.Errorf("--name is required")
			}
			return runEdit(func(doc *psrt.Document) error {
				return editor.RenamePage(doc, pageName, newName)
			})
		},
	}
	rename.Flags().StringVar(&newName, "name", "", "new page name")

	setPath := &cobra.Command{
		Use:     "set-path",
		Short:   "Set page background image URL",
		Example: `  psrt-edit --input=doc.psrt page set-path --page=intro --path=https://example.com/bg.avif`,
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := requirePage(); err != nil {
				return err
			}
			if strings.TrimSpace(path) == "" {
				return fmt.Errorf("--path is required")
			}
			return runEdit(func(doc *psrt.Document) error {
				return editor.SetPagePath(doc, pageName, path)
			})
		},
	}
	setPath.Flags().StringVar(&path, "path", "", "image URL")

	styleSet := &cobra.Command{
		Use:     "style-set",
		Short:   "Set or merge page style JSON (--key/--value or --style)",
		Example: `  psrt-edit --input=doc.psrt page style-set --page=capa --key=backGround --value='"#0F0F14"'`,
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := requirePage(); err != nil {
				return err
			}
			return runEdit(func(doc *psrt.Document) error {
				return editor.SetPageStyle(doc, pageName, styleKey, styleValue, []byte(styleJSON))
			})
		},
	}
	styleSet.Flags().StringVar(&styleKey, "key", "", "style property key")
	styleSet.Flags().StringVar(&styleValue, "value", "", "style property value (JSON literal, e.g. \"#fff\" or \"600\")")
	styleSet.Flags().StringVar(&styleJSON, "style", "", "partial style JSON object to merge")

	styleRemove := &cobra.Command{
		Use:     "style-remove",
		Short:   "Remove one property from page style JSON",
		Example: `  psrt-edit --input=doc.psrt page style-remove --page=capa --key=backGround`,
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := requirePage(); err != nil {
				return err
			}
			if strings.TrimSpace(styleKey) == "" {
				return fmt.Errorf("--key is required")
			}
			return runEdit(func(doc *psrt.Document) error {
				return editor.RemovePageStyleKey(doc, pageName, styleKey)
			})
		},
	}
	styleRemove.Flags().StringVar(&styleKey, "key", "", "style property key")

	move := &cobra.Command{
		Use:     "move",
		Short:   "Reorder page (--before or --after another page name)",
		Example: `  psrt-edit --input=doc.psrt page move --page=intro --before=capa`,
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := requirePage(); err != nil {
				return err
			}
			return runEdit(func(doc *psrt.Document) error {
				return editor.MovePage(doc, pageName, beforePage, afterPage)
			})
		},
	}
	move.Flags().StringVar(&beforePage, "before", "", "place before this page name")
	move.Flags().StringVar(&afterPage, "after", "", "place after this page name")

	for _, c := range []*cobra.Command{rename, setPath, styleSet, styleRemove, move} {
		c.Flags().StringVar(&pageName, "page", "", "page name")
		_ = c.MarkFlagRequired("page")
	}

	cmd.AddCommand(rename, setPath, styleSet, styleRemove, move)
	return cmd
}

func textCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "text",
		Short:   "Edit text blocks (see: psrt-edit text --help)",
		Long:    textLong,
		Example: textExample,
	}

	var (
		styleKey, styleValue, stylePartial, addStyleJSON, content string
		appendContent                                             bool
		x, y, width, textSize                                     float64
		posX, posY, posW, posS                                    float64
		newIndex                                                  int
		imageRef                                                  string
		beforeIdx, afterIdx                                       int
		reorderTo, reorderBy                                      int
	)

	styleSet := &cobra.Command{
		Use:     "style-set",
		Short:   "Set or merge text style JSON (--key/--value or --style)",
		Example: `  psrt-edit --input=doc.psrt text style-set --page=intro --index=1 --key=color --value='"#fff"'`,
		RunE: textRun(func(doc *psrt.Document, idx int) error {
			return editor.SetTextStyle(doc, pageName, idx, styleKey, styleValue, []byte(stylePartial))
		}),
	}
	styleSet.Flags().StringVar(&styleKey, "key", "", "style property key")
	styleSet.Flags().StringVar(&styleValue, "value", "", "style property value (JSON literal)")
	styleSet.Flags().StringVar(&stylePartial, "style", "", "partial style JSON object to merge")

	styleRemove := &cobra.Command{
		Use:     "style-remove",
		Short:   "Remove one property from text style JSON",
		Example: `  psrt-edit --input=doc.psrt text style-remove --page=intro --index=1 --key=background`,
		RunE: textRun(func(doc *psrt.Document, idx int) error {
			if strings.TrimSpace(styleKey) == "" {
				return fmt.Errorf("--key is required")
			}
			return editor.RemoveTextStyleKey(doc, pageName, idx, styleKey)
		}),
	}
	styleRemove.Flags().StringVar(&styleKey, "key", "", "style property key")

	setContent := &cobra.Command{
		Use:     "set-content",
		Short:   "Replace or append the text body (--content required)",
		Long:    "Updates the lines under the >> header. Use --index from the header (e.g. >>…| 1 → --index=1).",
		Example: `  psrt-edit --input=doc.psrt text set-content --page=intro --index=1 --content=revisor`,
		RunE: textRun(func(doc *psrt.Document, idx int) error {
			return editor.SetTextContent(doc, pageName, idx, content, appendContent)
		}),
	}
	setContent.Flags().StringVar(&content, "content", "", "new text body (required)")
	setContent.Flags().BoolVar(&appendContent, "append", false, "append to existing content instead of replacing")

	add := &cobra.Command{
		Use:     "add",
		Short:   "Add a new text block on a page",
		Example: `  psrt-edit --input=doc.psrt text add --page=intro --x=10 --y=20 --width=80 --text-size=3 --index=3 --content=Footer`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if err := requirePage(); err != nil {
				return err
			}
			before := -1
			after := -1
			if cmd.Flags().Changed("before") {
				before = beforeIdx
			}
			if cmd.Flags().Changed("after") {
				after = afterIdx
			}
			t := psrt.Text{
				BaseBlock: psrt.BaseBlock{
					X: x, Y: y, Width: width, Index: newIndex, ImageRef: imageRef,
				},
				TextSize: textSize,
				Content:  content,
			}
			if strings.TrimSpace(addStyleJSON) != "" {
				t.Style = psrt.Style(addStyleJSON)
			} else {
				t.Style = psrt.Style("{}")
			}
			return runEdit(func(doc *psrt.Document) error {
				return editor.AddText(doc, pageName, t, before, after)
			})
		},
	}
	add.Flags().StringVar(&pageName, "page", "", "page name")
	add.Flags().Float64Var(&x, "x", 0, "X coordinate (percent)")
	add.Flags().Float64Var(&y, "y", 0, "Y coordinate (percent)")
	add.Flags().Float64VarP(&width, "width","w", 0, "width (percent)")
	add.Flags().Float64VarP(&textSize, "text-size", "s", 0, "text size (percent)")

	add.Flags().IntVar(&newIndex, "index", 0, "text index id")
	add.Flags().StringVar(&content, "content", "", "text body")
	add.Flags().StringVar(&addStyleJSON, "style", "{}", "style JSON")
	add.Flags().StringVar(&imageRef, "image-ref", "", "optional image reference")
	add.Flags().IntVar(&beforeIdx, "before", 0, "insert before text with this index")
	add.Flags().IntVar(&afterIdx, "after", 0, "insert after text with this index")
	_ = add.MarkFlagRequired("page")

	remove := &cobra.Command{
		Use:     "remove",
		Short:   "Remove a text block by --index",
		Example: `  psrt-edit --input=doc.psrt text remove --page=intro --index=2`,
		RunE: textRun(func(doc *psrt.Document, idx int) error {
			return editor.RemoveText(doc, pageName, idx)
		}),
	}

	reorder := &cobra.Command{
		Use:     "reorder",
		Short:   "Change block order in file (--before or --after another --index)",
		Example: `  psrt-edit --input=doc.psrt text reorder --page=intro --index=2 --before=1`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if err := requireTextIndex(); err != nil {
				return err
			}
			idx, err := parseTextIndexFlag()
			if err != nil {
				return fmt.Errorf("--index: %w", err)
			}
			before := -1
			after := -1
			if cmd.Flags().Changed("before") {
				before = beforeIdx
			}
			if cmd.Flags().Changed("after") {
				after = afterIdx
			}
			return runEdit(func(doc *psrt.Document) error {
				return editor.ReorderTextRelative(doc, pageName, idx, before, after)
			})
		},
	}
	reorder.Flags().IntVar(&beforeIdx, "before", 0, "place before text with this index")
	reorder.Flags().IntVar(&afterIdx, "after", 0, "place after text with this index")

	reorderToCmd := &cobra.Command{
		Use:     "reorder-to",
		Short:   "Move block to absolute position in file order (0-based --to)",
		Example: `  psrt-edit --input=doc.psrt text reorder-to --page=intro --index=2 --to=0`,
		RunE: textRun(func(doc *psrt.Document, idx int) error {
			return editor.ReorderTextTo(doc, pageName, idx, reorderTo)
		}),
	}
	reorderToCmd.Flags().IntVar(&reorderTo, "to", 0, "target slice position (0-based)")

	reorderByCmd := &cobra.Command{
		Use:     "reorder-by",
		Short:   "Move block by delta in file order (--by, may be negative)",
		Example: `  psrt-edit --input=doc.psrt text reorder-by --page=intro --index=0 --by=1`,
		RunE: textRun(func(doc *psrt.Document, idx int) error {
			return editor.ReorderTextByDelta(doc, pageName, idx, reorderBy)
		}),
	}
	reorderByCmd.Flags().IntVar(&reorderBy, "by", 0, "order delta (may be negative)")

	positionSet := &cobra.Command{
		Use:     "position-set",
		Short:   "Set canvas coordinates in >> header (percent); only flags you pass are changed",
		Example: `  psrt-edit --input=doc.psrt text position-set --page=intro --index=1 --x=11.6 --y=56.11`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if err := requireTextIndex(); err != nil {
				return err
			}
			idx, err := parseTextIndexFlag()
			if err != nil {
				return fmt.Errorf("--index: %w", err)
			}
			pos, err := positionFieldsFromFlags(cmd, posX, posY, posW, posS)
			if err != nil {
				return err
			}
			return runEdit(func(doc *psrt.Document) error {
				return editor.SetTextPosition(doc, pageName, idx, pos)
			})
		},
	}
	bindPositionFlags(positionSet, &posX, &posY, &posW, &posS)

	positionNudge := &cobra.Command{
		Use:     "position-nudge",
		Short:   "Nudge canvas coordinates by delta; only flags you pass are applied",
		Example: `  psrt-edit --input=doc.psrt text position-nudge --page=intro --index=1 --y=-0.5`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if err := requireTextIndex(); err != nil {
				return err
			}
			idx, err := parseTextIndexFlag()
			if err != nil {
				return fmt.Errorf("--index: %w", err)
			}
			pos, err := positionFieldsFromFlags(cmd, posX, posY, posW, posS)
			if err != nil {
				return err
			}
			return runEdit(func(doc *psrt.Document) error {
				return editor.NudgeTextPosition(doc, pageName, idx, pos)
			})
		},
	}
	bindPositionFlags(positionNudge, &posX, &posY, &posW, &posS)

	for _, c := range []*cobra.Command{styleSet, styleRemove, setContent, remove, reorder, reorderToCmd, reorderByCmd, positionSet, positionNudge} {
		bindTextFlags(c)
	}

	cmd.AddCommand(styleSet, styleRemove, setContent, add, remove, reorder, reorderToCmd, reorderByCmd, positionSet, positionNudge)
	return cmd
}

func fontCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "font",
		Short:   "Edit font URLs in $FONTS",
		Long:    fontLong,
		Example: fontExample,
	}
	var url string

	add := &cobra.Command{
		Use:   "add",
		Short: "Add a font URL",
		RunE: func(_ *cobra.Command, _ []string) error {
			if strings.TrimSpace(url) == "" {
				return fmt.Errorf("--url is required")
			}
			return runEdit(func(doc *psrt.Document) error {
				return editor.AddFont(doc, url)
			})
		},
	}
	add.Flags().StringVar(&url, "url", "", "font URL")

	remove := &cobra.Command{
		Use:   "remove",
		Short: "Remove a font URL",
		RunE: func(_ *cobra.Command, _ []string) error {
			if strings.TrimSpace(url) == "" {
				return fmt.Errorf("--url is required")
			}
			return runEdit(func(doc *psrt.Document) error {
				return editor.RemoveFont(doc, url)
			})
		},
	}
	remove.Flags().StringVar(&url, "url", "", "font URL")

	cmd.AddCommand(add, remove)
	return cmd
}

func constCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "const",
		Short:   "Edit named constants in $CONSTS",
		Long:    constLong,
		Example: constExample,
	}
	var name, value string

	add := &cobra.Command{
		Use:   "add",
		Short: "Add a constant and replace matching literals with @name@",
		RunE: func(_ *cobra.Command, _ []string) error {
			if strings.TrimSpace(name) == "" {
				return fmt.Errorf("--name is required")
			}
			return runEdit(func(doc *psrt.Document) error {
				return editor.AddConst(doc, name, value)
			})
		},
	}
	add.Flags().StringVar(&name, "name", "", "constant name (without @)")
	add.Flags().StringVar(&value, "value", "", "constant value")
	_ = add.MarkFlagRequired("name")

	remove := &cobra.Command{
		Use:   "remove",
		Short: "Revert @name@ references and remove the constant",
		RunE: func(_ *cobra.Command, _ []string) error {
			if strings.TrimSpace(name) == "" {
				return fmt.Errorf("--name is required")
			}
			return runEdit(func(doc *psrt.Document) error {
				return editor.RemoveConst(doc, name)
			})
		},
	}
	remove.Flags().StringVar(&name, "name", "", "constant name")

	cmd.AddCommand(add, remove)
	return cmd
}

func requirePage() error {
	if strings.TrimSpace(pageName) == "" {
		return fmt.Errorf("--page is required")
	}
	return nil
}

func requireTextIndex() error {
	if err := requirePage(); err != nil {
		return err
	}
	if strings.TrimSpace(textIndex) == "" {
		return fmt.Errorf("--index is required")
	}
	return nil
}

func parseTextIndexFlag() (int, error) {
	return editor.ParseTextIndex(strings.TrimSpace(textIndex))
}

type textMutator func(*psrt.Document, int) error

func textRun(mutate textMutator) func(*cobra.Command, []string) error {
	return func(_ *cobra.Command, _ []string) error {
		if err := requireTextIndex(); err != nil {
			return err
		}
		idx, err := parseTextIndexFlag()
		if err != nil {
			return fmt.Errorf("--index: %w", err)
		}
		return runEdit(func(doc *psrt.Document) error {
			return mutate(doc, idx)
		})
	}
}

func bindTextFlags(c *cobra.Command) {
	c.Flags().StringVar(&pageName, "page", "", "page name ($START name)")
	c.Flags().StringVar(&textIndex, "index", "", "text id from >> header (third pipe field)")
	_ = c.MarkFlagRequired("page")
	_ = c.MarkFlagRequired("index")
}

func bindPositionFlags(c *cobra.Command, x, y, w, s *float64) {
	c.Flags().Float64Var(x, "x", 0, "X coordinate (percent)")
	c.Flags().Float64Var(y, "y", 0, "Y coordinate (percent)")
	c.Flags().Float64Var(w, "width", 0, "width (percent)")
	c.Flags().Float64VarP(s, "text-size", "s", 0, "text size (percent)")
}

func positionFieldsFromFlags(cmd *cobra.Command, x, y, w, s float64) (editor.PositionFields, error) {
	var pos editor.PositionFields
	if cmd.Flags().Changed("x") {
		pos.X = &x
	}
	if cmd.Flags().Changed("y") {
		pos.Y = &y
	}
	if cmd.Flags().Changed("width") {
		pos.Width = &w
	}
	if cmd.Flags().Changed("text-size") || cmd.Flags().Changed("s") {
		pos.TextSize = &s
	}
	if pos.X == nil && pos.Y == nil && pos.Width == nil && pos.TextSize == nil {
		return pos, fmt.Errorf("at least one of --x, --y, --width, --text-size (or -s) is required")
	}
	return pos, nil
}

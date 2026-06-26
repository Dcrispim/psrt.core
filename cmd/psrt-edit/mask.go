package main

import (
	"github.com/spf13/cobra"
	"github.com/Dcrispim/psrt.core/psrt"
	"github.com/Dcrispim/psrt.core/psrt/editor"
)

func maskCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "mask",
		Short:   "Edit mask blocks (== …)",
		Long:    maskLong,
		Example: maskExample,
	}
	cmd.AddCommand(maskPositionSetCmd())
	cmd.AddCommand(maskAddCmd())
	cmd.AddCommand(maskRemoveCmd())
	return cmd
}

func maskPositionSetCmd() *cobra.Command {
	var x, y, width, height float64
	c := &cobra.Command{
		Use:   "position-set",
		Short: "Set mask X/Y/width/height (percent)",
		RunE: func(cmd *cobra.Command, _ []string) error {
			if err := requireTextIndex(); err != nil {
				return err
			}
			idx, err := parseTextIndexFlag()
			if err != nil {
				return err
			}
			pos := editor.MaskPositionFields{}
			if cmd.Flags().Changed("x") {
				pos.X = &x
			}
			if cmd.Flags().Changed("y") {
				pos.Y = &y
			}
			if cmd.Flags().Changed("width") {
				pos.Width = &width
			}
			if cmd.Flags().Changed("height") {
				pos.Height = &height
			}
			return runEdit(func(doc *psrt.Document) error {
				return editor.SetMaskPosition(doc, pageName, idx, pos)
			})
		},
	}
	bindTextFlags(c)
	c.Flags().Float64Var(&x, "x", 0, "X (percent)")
	c.Flags().Float64Var(&y, "y", 0, "Y (percent)")
	c.Flags().Float64VarP(&width, "width", "w", 0, "width (percent)")
	c.Flags().Float64Var(&height, "height", 0, "height (percent of page)")
	return c
}

func maskAddCmd() *cobra.Command {
	var x, y, width, height float64
	var newIndex int
	var styleJSON, imageRef string
	c := &cobra.Command{
		Use:   "add",
		Short: "Add a mask block",
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := requirePage(); err != nil {
				return err
			}
			m := psrt.Mask{
				BaseBlock: psrt.BaseBlock{
					X: x, Y: y, Width: width, Index: newIndex, ImageRef: imageRef,
				},
				Height: height,
			}
			if styleJSON != "" {
				m.Style = psrt.Style(styleJSON)
			} else {
				m.Style = psrt.Style("{}")
			}
			return runEdit(func(doc *psrt.Document) error {
				return editor.AddMask(doc, pageName, m, -1, -1)
			})
		},
	}
	c.Flags().StringVar(&pageName, "page", "", "page name")
	c.Flags().IntVar(&newIndex, "index", 0, "mask index in == header")
	c.Flags().Float64Var(&x, "x", 10, "X (percent)")
	c.Flags().Float64Var(&y, "y", 10, "Y (percent)")
	c.Flags().Float64VarP(&width, "width", "w", 20, "width (percent)")
	c.Flags().Float64Var(&height, "height", 5, "height (percent)")
	c.Flags().StringVar(&styleJSON, "style", "", "style JSON")
	c.Flags().StringVar(&imageRef, "image-ref", "", "optional image ref")
	_ = c.MarkFlagRequired("page")
	return c
}

func maskRemoveCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "remove",
		Short: "Remove a mask block",
		RunE: textRun(func(doc *psrt.Document, idx int) error {
			return editor.RemoveMask(doc, pageName, idx)
		}),
	}
	bindTextFlags(c)
	return c
}

const maskLong = `Mask block commands (== header).

--page     Page name from $START … $END.
--index    Numeric id in the == header (same index space as >> blocks).

Subcommands:
  position-set   Set X/Y/width/height (percent); pass only flags to change
  add            Insert a new mask block
  remove         Delete a mask block`

const maskExample = `  psrt-edit --input=doc.psrt mask position-set --page=p1 --index=0 --height=8
  psrt-edit --input=doc.psrt mask add --page=p1 --index=2 --x=10 --y=10 --width=20 --height=5`

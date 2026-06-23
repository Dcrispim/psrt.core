package main

// Help text for psrt-edit (shown via --help on each command).

const rootLong = `Edit PSRT files from the command line.

Syntax:
  psrt-edit --input=FILE <group> <action> [flags]

Groups:
  page   Edit pages ($START … $END)
  text   Edit text blocks (>> …)
  mask   Edit mask blocks (== …)
  font   Edit $FONTS … $ENDFONTS
  const  Edit $CONSTS … $ENDCONSTS

Global output:
  By default the input file is overwritten. Use --output to write elsewhere.
  --cat prints the resulting .psrt to stdout instead of saving.
  --compile-svg / --compile-html run compilation after the edit (see global flags).

Important:
  <action> must come immediately after <group>. Flags usually follow <action>.
  Wrong:  psrt-edit --input=f.psrt text --page=p set-content hello
  Right:  psrt-edit --input=f.psrt text set-content --page=p --index=0 --content=hello`

const rootExample = `  # Rename a page (overwrites input)
  psrt-edit --input=exemplo.psrt page rename --page=capa --name=cover

  # Change text body; --index is the number in the >> header, not line number
  psrt-edit --input=exemplo.psrt text set-content --page=intro --index=1 --content=revisor

  # Preview result without saving
  psrt-edit --input=exemplo.psrt --cat text set-content --page=intro --index=1 --content=revisor

  # Edit and compile SVG
  psrt-edit --input=exemplo.psrt --compile-svg text position-nudge --page=intro --index=1 --y=-1`

const pageLong = `Page commands. Every subcommand requires --page=<name> (the $START name).

Subcommands:
  rename        Change page name (--name)
  set-path      Change background image URL (--path)
  style-set     Merge page style JSON (--key/--value or --style)
  style-remove  Remove one style property (--key)
  move          Reorder pages (--before or --after another page name)`

const pageExample = `  psrt-edit --input=doc.psrt page rename --page=capa --name=cover
  psrt-edit --input=doc.psrt page set-path --page=intro --path=https://example.com/bg.avif
  psrt-edit --input=doc.psrt page style-set --page=capa --key=backGround --value='"#000"'
  psrt-edit --input=doc.psrt page move --page=intro --before=capa`

const textLong = `Text block commands.

Syntax:
  psrt-edit --input=FILE text <subcommand> --page=PAGE --index=N [flags]

--page     Page name from $START … $END.
--index    The numeric id in the text header (>> … | index), NOT a line number.
           Example: >>11.6-56-11-77-3 | {...} | 1  →  --index=1

Subcommands:
  set-content     Replace or append body text (--content, optional --append)
  style-set       Merge CSS JSON on the text (--key/--value or --style)
  style-remove    Remove one style key (--key)
  add             Insert a new text block (needs --x --y --width --text-size --index --content)
  remove          Delete a text block
  reorder         Change block order vs another index (--before or --after)
  reorder-to      Move to absolute position in file order (--to)
  reorder-by      Move by delta in file order (--by, may be negative)
  position-set    Set X/Y/width/text-size (percent); only pass flags you want to change
  position-nudge  Add deltas to coordinates; only pass flags you want to nudge

Reorder vs position:
  reorder-*     Changes the order of blocks in the .psrt file.
  position-*    Changes >> X-Y-Width-TextSize coordinates on the canvas.

Style --value accepts JSON literals or bare tokens (PowerShell-friendly):
  --key=color --value=#fff
  --key=fontWeight --value=600
  Quoted JSON still works: --value='"#1DB954"'`

const textExample = `  psrt-edit --input=doc.psrt text set-content --page=intro --index=1 --content=revisor
  psrt-edit --input=doc.psrt text style-set --page=intro --index=1 --key=color --value='"#fff"'
  psrt-edit --input=doc.psrt text position-nudge --page=intro --index=2 --y=0.5
  psrt-edit --input=doc.psrt text reorder --page=intro --index=2 --before=1`

const fontLong = `Font URL commands ($FONTS section).

Subcommands:
  add      Append a font URL (--url); duplicates are ignored
  remove   Remove a font URL (--url)`

const fontExample = `  psrt-edit --input=doc.psrt font add --url=https://cdn.example/font.woff2
  psrt-edit --input=doc.psrt font remove --url=https://cdn.example/font.woff2`

const constLong = `Named constant commands ($CONSTS section).

Usage in PSRT: @name@ in text, URLs, or style JSON.

Subcommands:
  add      Define a constant (--name, --value); replaces matching literals with @name@
  remove   Reverts all @name@ to the stored value, then deletes the constant`

const constExample = `  psrt-edit --input=doc.psrt const add --name=accent --value=#1DB954
  psrt-edit --input=doc.psrt const remove --name=accent`

package main

import (
	"context"
	"embed"
	"log"
	"os"
	"path/filepath"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/runtime"

	"psrt/compilesvg/textoutline"
	"psrt/internal/crashlog"
	"psrt/internal/visualapp"
	"psrt/psrt"
)

func init() {
	if dir := findPSRTGUIProjectDir(); dir != "" {
		textoutline.AddChromeSearchRoots(filepath.Join(dir, "build", "bin"))
	}
	if os.Getenv("PSRT_DEBUG_CHROME") != "" {
		if p := textoutline.ResolveChromeExec(); p != "" {
			log.Printf("psrt: using chrome %s", p)
		} else {
			log.Printf("psrt: chrome not found")
		}
	}
}

func findPSRTGUIProjectDir() string {
	var seeds []string
	if exe, err := os.Executable(); err == nil {
		seeds = append(seeds, filepath.Dir(exe))
	}
	if wd, err := os.Getwd(); err == nil {
		seeds = append(seeds, wd)
	}
	for _, seed := range seeds {
		dir := filepath.Clean(seed)
		for i := 0; i < 12; i++ {
			if _, err := os.Stat(filepath.Join(dir, "wails.json")); err == nil {
				return dir
			}
			parent := filepath.Dir(dir)
			if parent == dir {
				break
			}
			dir = parent
		}
	}
	return ""
}

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	defer crashlog.RecoverMain()
	crashlog.Install()
	gui := NewGUIApp()
	err := wails.Run(&options.App{
		Title:  "PSRT Visual Editor",
		Width:  1280,
		Height: 800,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		OnStartup:  gui.startup,
		OnShutdown: gui.shutdown,
		Bind: []interface{}{
			gui,
		},
	})
	if err != nil {
		log.Fatal(err)
	}
}

// CompileSVGResult is returned by SVG compile/preview bindings.
type CompileSVGResult struct {
	URI                string `json:"uri"`
	UsedGoTextFallback bool   `json:"usedGoTextFallback"`
}

// GUIApp exposes methods to the frontend (Wails bindings).
type GUIApp struct {
	ctx context.Context
	app *visualapp.App
}

func NewGUIApp() *GUIApp {
	return &GUIApp{}
}

func (g *GUIApp) startup(ctx context.Context) {
	g.ctx = ctx
	g.app = visualapp.New(func(event string, data interface{}) {
		runtime.EventsEmit(ctx, event, data)
	})
}

func (g *GUIApp) shutdown(context.Context) {}

func (g *GUIApp) OpenFileDialog() (res visualapp.OpenFileResult, err error) {
	var path string
	defer guardOpenFile(g, &path, &err)
	path, err = runtime.OpenFileDialog(g.ctx, runtime.OpenDialogOptions{
		Title: "Open PSRT",
		Filters: []runtime.FileFilter{
			{DisplayName: "PSRT", Pattern: "*.psrt"},
		},
	})
	if err != nil || path == "" {
		return visualapp.OpenFileResult{}, err
	}
	if err = g.app.OpenFile(path); err != nil {
		crashlog.WriteError("OpenFile", path, err)
		return visualapp.OpenFileResult{}, err
	}
	docJSON, err := g.app.GetDocumentJSON()
	if err != nil {
		crashlog.WriteError("GetDocumentJSON", path, err)
		return visualapp.OpenFileResult{}, err
	}
	return visualapp.OpenFileResult{FilePath: path, Document: docJSON}, nil
}

func (g *GUIApp) OpenImageFileDialog() (path string, err error) {
	defer guardErr(g, "OpenImageFileDialog", &err)
	path, err = runtime.OpenFileDialog(g.ctx, runtime.OpenDialogOptions{
		Title: "Selecionar imagem",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "Imagens",
				Pattern:     "*.png;*.jpg;*.jpeg;*.gif;*.webp;*.avif;*.svg",
			},
		},
	})
	return path, err
}

func (g *GUIApp) Save() (err error) {
	defer guardErr(g, "Save", &err)
	err = g.app.Save()
	if err != nil {
		crashlog.WriteError("Save", g.currentFilePath(), err)
	}
	return err
}

func (g *GUIApp) SaveDocumentJSON(docJSON string) (err error) {
	defer guardErr(g, "SaveDocumentJSON", &err)
	err = g.app.SaveDocumentJSON(docJSON)
	if err != nil {
		crashlog.WriteError("SaveDocumentJSON", g.currentFilePath(), err)
	}
	return err
}

func (g *GUIApp) SaveAsDocumentJSON(docJSON string) (path string, err error) {
	defer guardErr(g, "SaveAsDocumentJSON", &err)
	path, err = runtime.SaveFileDialog(g.ctx, runtime.SaveDialogOptions{
		Title:           "Save PSRT",
		DefaultFilename: "document.psrt",
		Filters: []runtime.FileFilter{
			{DisplayName: "PSRT", Pattern: "*.psrt"},
		},
	})
	if err != nil || path == "" {
		return "", err
	}
	err = g.app.SaveDocumentJSONTo(docJSON, path)
	if err != nil {
		crashlog.WriteError("SaveAsDocumentJSON", path, err)
	}
	return path, err
}

func (g *GUIApp) GetState() (st visualapp.UIState, err error) {
	defer guardErr(g, "GetState", &err)
	return g.app.GetState()
}

func (g *GUIApp) SetActivePage(name string) (err error) {
	defer guardErr(g, "SetActivePage", &err)
	return g.app.SetActivePage(name)
}

func (g *GUIApp) BeginEdit() { defer guardVoid(g, "BeginEdit"); g.app.BeginEdit() }

func (g *GUIApp) EndEdit() { defer guardVoid(g, "EndEdit"); g.app.EndEdit() }

func (g *GUIApp) SelectText(index int) { defer guardVoid(g, "SelectText"); g.app.SelectText(index) }

func (g *GUIApp) PatchText(pageName string, index int, patch visualapp.TextPatch) (err error) {
	defer guardErr(g, "PatchText", &err)
	return g.app.PatchText(pageName, index, patch)
}

func (g *GUIApp) PatchMask(pageName string, index int, patch visualapp.MaskPatch) (err error) {
	defer guardErr(g, "PatchMask", &err)
	return g.app.PatchMask(pageName, index, patch)
}

func (g *GUIApp) PatchPage(patch visualapp.PagePatch) (err error) {
	defer guardErr(g, "PatchPage", &err)
	return g.app.PatchPage(patch)
}

func (g *GUIApp) AddPage(name, imageURL, styleJSON string) (err error) {
	defer guardErr(g, "AddPage", &err)
	return g.app.AddPage(name, imageURL, styleJSON)
}

func (g *GUIApp) RemovePage(name string) (err error) {
	defer guardErr(g, "RemovePage", &err)
	return g.app.RemovePage(name)
}

func (g *GUIApp) MovePage(name, ref string, before bool) (err error) {
	defer guardErr(g, "MovePage", &err)
	return g.app.MovePage(name, ref, before)
}

func (g *GUIApp) AddTextBlock(index int, x, y, width, textSize float64, content, styleJSON, imageRef string) (err error) {
	defer guardErr(g, "AddTextBlock", &err)
	return g.app.AddTextBlock(index, x, y, width, textSize, content, styleJSON, imageRef)
}

func (g *GUIApp) RemoveText(index int) (err error) {
	defer guardErr(g, "RemoveText", &err)
	return g.app.RemoveText(index)
}

func (g *GUIApp) ReorderText(index, ref int, before bool) (err error) {
	defer guardErr(g, "ReorderText", &err)
	return g.app.ReorderText(index, ref, before)
}

func (g *GUIApp) GetAssetDataURI(url string) (uri string, err error) {
	defer guardErr(g, "GetAssetDataURI", &err)
	uri, err = g.app.GetAssetDataURI(url)
	if err != nil {
		crashlog.WriteError("GetAssetDataURI", g.currentFilePath(), err)
	}
	return uri, err
}

func (g *GUIApp) RefreshPageImage() (err error) {
	defer guardErr(g, "RefreshPageImage", &err)
	return g.app.RefreshPageImage()
}

func (g *GUIApp) RefreshAssetURL(url string) (err error) {
	defer guardErr(g, "RefreshAssetURL", &err)
	return g.app.RefreshAssetURL(url)
}

func (g *GUIApp) CompilePageSVG(page string) (CompileSVGResult, error) {
	var err error
	defer guardErr(g, "CompilePageSVG", &err)
	res, e := g.app.CompilePageSVG(page)
	if e != nil {
		return CompileSVGResult{}, e
	}
	return CompileSVGResult{URI: res.URI, UsedGoTextFallback: res.UsedGoTextFallback}, nil
}

func (g *GUIApp) CompileDocumentHTML() (out string, err error) {
	defer guardErr(g, "CompileDocumentHTML", &err)
	return g.app.CompileDocumentHTML()
}

func (g *GUIApp) CompilePageHTML(page string) (out string, err error) {
	defer guardErr(g, "CompilePageHTML", &err)
	return g.app.CompilePageHTML(page)
}

func (g *GUIApp) ExportSVG(dir string) (CompileSVGResult, error) {
	var err error
	defer guardErr(g, "ExportSVG", &err)
	res, e := g.app.ExportSVG(dir)
	if e != nil {
		return CompileSVGResult{}, e
	}
	return CompileSVGResult{UsedGoTextFallback: res.UsedGoTextFallback}, nil
}

func (g *GUIApp) ExportSVGFromDocument(docJSON string) (CompileSVGResult, error) {
	var err error
	defer guardErr(g, "ExportSVGFromDocument", &err)
	parentDir, err := runtime.OpenDirectoryDialog(g.ctx, runtime.OpenDialogOptions{
		Title: "Selecionar pasta para exportar SVG",
	})
	if err != nil || parentDir == "" {
		return CompileSVGResult{}, err
	}
	res, e := g.app.ExportSVGFromDocument(docJSON, parentDir, exportBaseName(g.currentFilePath()))
	if e != nil {
		crashlog.WriteError("ExportSVGFromDocument", res.URI, e)
		return CompileSVGResult{}, e
	}
	return CompileSVGResult{URI: res.URI, UsedGoTextFallback: res.UsedGoTextFallback}, nil
}

func (g *GUIApp) ExportHTMLFromDocument(docJSON string, variantPaths []string, variantBodies []visualapp.VariantPSRT) (outPath string, err error) {
	defer guardErr(g, "ExportHTMLFromDocument", &err)
	dir, err := runtime.OpenDirectoryDialog(g.ctx, runtime.OpenDialogOptions{
		Title: "Selecionar pasta para salvar HTML",
	})
	if err != nil || dir == "" {
		return "", err
	}
	outPath, err = g.app.ExportHTMLFromDocument(docJSON, dir, exportBaseName(g.currentFilePath()), variantPaths, variantBodies)
	if err != nil {
		crashlog.WriteError("ExportHTMLFromDocument", outPath, err)
	}
	return outPath, err
}

func exportBaseName(path string) string {
	if path == "" {
		return "document"
	}
	base := filepath.Base(path)
	ext := filepath.Ext(base)
	if ext != "" {
		return base[:len(base)-len(ext)]
	}
	return base
}

func (g *GUIApp) SetAutoCompile(on bool) { defer guardVoid(g, "SetAutoCompile"); g.app.SetAutoCompile(on) }

func (g *GUIApp) Undo() (err error) {
	defer guardErr(g, "Undo", &err)
	return g.app.Undo()
}

func (g *GUIApp) Redo() (err error) {
	defer guardErr(g, "Redo", &err)
	return g.app.Redo()
}

func (g *GUIApp) AddFont(url string) (err error) {
	defer guardErr(g, "AddFont", &err)
	return g.app.AddFont(url)
}

func (g *GUIApp) RemoveFont(url string) (err error) {
	defer guardErr(g, "RemoveFont", &err)
	return g.app.RemoveFont(url)
}

func (g *GUIApp) AddConst(name, value string) (err error) {
	defer guardErr(g, "AddConst", &err)
	return g.app.AddConst(name, value)
}

func (g *GUIApp) RemoveConst(name string) (err error) {
	defer guardErr(g, "RemoveConst", &err)
	return g.app.RemoveConst(name)
}

func (g *GUIApp) GetDocumentPSRT() (out string, err error) {
	defer guardErr(g, "GetDocumentPSRT", &err)
	return g.app.GetDocumentPSRT()
}

func (g *GUIApp) SetDocumentFromPSRT(text string) (err error) {
	defer guardErr(g, "SetDocumentFromPSRT", &err)
	return g.app.SetDocumentFromPSRT(text)
}

func (g *GUIApp) GetDocumentJSON() (out string, err error) {
	defer guardErr(g, "GetDocumentJSON", &err)
	return g.app.GetDocumentJSON()
}

func (g *GUIApp) ParseDocumentPSRT(text string) (out string, err error) {
	defer guardErr(g, "ParseDocumentPSRT", &err)
	return g.app.ParseDocumentPSRT(text)
}

func (g *GUIApp) FormatDocumentJSON(docJSON string) (out string, err error) {
	defer guardErr(g, "FormatDocumentJSON", &err)
	return g.app.FormatDocumentJSON(docJSON)
}

func (g *GUIApp) FormatPageDocumentJSON(docJSON, pageName string) (out string, err error) {
	defer guardErr(g, "FormatPageDocumentJSON", &err)
	return g.app.FormatPageDocumentJSON(docJSON, pageName)
}

func (g *GUIApp) RenderTextContentHTML(content string) string {
	return psrt.RenderInlineHTML(content)
}

func (g *GUIApp) MergePageDocumentPSRT(fullDocJSON, pageName, psrtText string) (out string, err error) {
	defer guardErr(g, "MergePageDocumentPSRT", &err)
	return g.app.MergePageDocumentPSRT(fullDocJSON, pageName, psrtText)
}

func (g *GUIApp) CompilePageSVGFromDocument(docJSON, page string) (CompileSVGResult, error) {
	var err error
	defer guardErr(g, "CompilePageSVGFromDocument", &err)
	res, e := g.app.CompilePageSVGFromDocument(docJSON, page)
	if e != nil {
		return CompileSVGResult{}, e
	}
	return CompileSVGResult{URI: res.URI, UsedGoTextFallback: res.UsedGoTextFallback}, nil
}

func (g *GUIApp) CompilePageHTMLFromDocument(docJSON, page string) (out string, err error) {
	defer guardErr(g, "CompilePageHTMLFromDocument", &err)
	return g.app.CompilePageHTMLFromDocument(docJSON, page)
}

// AdaptEntriesForWeb returns adapted CSS (container + text) per text block for the web preview.
func (g *GUIApp) AdaptEntriesForWeb(entriesJSON string, canvasW, canvasH int, zoom float64) (out []visualapp.WebPreviewStyle, err error) {
	defer guardErr(g, "AdaptEntriesForWeb", &err)
	return visualapp.AdaptEntriesForWeb(entriesJSON, canvasW, canvasH, zoom)
}

// AdaptTextStyleForWeb adapts one text block style for the web preview.
func (g *GUIApp) AdaptTextStyleForWeb(
	styleJSON string,
	content string,
	x, y, width, textSize float64,
	canvasW, canvasH int,
	zoom float64,
) (out visualapp.WebPreviewStyle) {
	defer guardVoid(g, "AdaptTextStyleForWeb")
	return visualapp.AdaptTextStyleForWeb(styleJSON, content, x, y, width, textSize, canvasW, canvasH, zoom)
}

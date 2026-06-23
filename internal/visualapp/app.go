package visualapp

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"psrt/compileasset"
	"psrt/compileasset/cache"
	"psrt/internal/crashlog"
	"psrt/psrt"
	"psrt/psrt/editor"
)

// App holds editor state for the visual GUI.
type App struct {
	mu           sync.Mutex
	doc          psrt.Document
	filePath     string
	activePage   string
	selectedIdx  int
	store        *cache.Store
	client       *http.Client
	autoCompile  bool
	undo          *undoStack
	redo          *undoStack
	emit          func(event string, data interface{})
	compileTimer *time.Timer
	compileMu    sync.Mutex
	snapGrid     float64
	inEdit       bool
}

// New creates an App with optional event emitter (Wails runtime).
func New(emit func(string, interface{})) *App {
	return &App{
		client:      &http.Client{Timeout: 30 * time.Second},
		selectedIdx: -1,
		undo:        newUndoStack(50),
		redo:        newUndoStack(50),
		emit:        emit,
		snapGrid:    0.5,
	}
}

func (a *App) snapshot() {
	a.undo.push(a.doc)
	a.redo = newUndoStack(50)
}

func (a *App) notify(err error) {
	if a.emit == nil {
		return
	}
	if err != nil {
		a.emit("error", err.Error())
		return
	}
	st, _ := a.buildState()
	a.emit("document:changed", st)
}

// FilePath returns the path of the currently open document, if any.
func (a *App) FilePath() string {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.filePath
}

// OpenFile loads a PSRT document.
func (a *App) OpenFile(path string) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	doc, err := editor.LoadDocument(path)
	if err != nil {
		return err
	}
	store, err := cache.NewStore("", path)
	if err != nil {
		return err
	}
	a.doc = doc
	a.filePath = path
	a.store = store
	if len(doc.Pages) > 0 {
		a.activePage = doc.Pages[0].Name
	}
	a.selectedIdx = -1
	a.undo = newUndoStack(50)
	a.redo = newUndoStack(50)
	fp := path
	crashlog.Go("EnsureDocument", fp, func() {
		_ = store.EnsureDocument(context.Background(), a.client, doc)
		a.notify(nil)
	})
	return nil
}

// Save writes the document to filePath.
func (a *App) Save() error {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.filePath == "" {
		return fmt.Errorf("no file open")
	}
	return a.saveTo(a.filePath)
}

// SaveAs saves to a new path.
func (a *App) SaveAs(path string) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	if err := a.saveTo(path); err != nil {
		return err
	}
	a.filePath = path
	store, err := cache.NewStore("", path)
	if err != nil {
		return err
	}
	a.store = store
	return nil
}

func (a *App) saveTo(path string) error {
	data, err := editor.FormatDocument(&a.doc)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

// GetState returns current UI state.
func (a *App) GetState() (UIState, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.buildState()
}

func (a *App) buildState() (UIState, error) {
	st := UIState{
		FilePath:      a.filePath,
		ActivePage:    a.activePage,
		SelectedIndex: a.selectedIdx,
		Fonts:         a.doc.Fonts,
		Consts:        a.doc.Consts,
		AutoCompile:   a.autoCompile,
	}
	for i := range a.doc.Pages {
		p := &a.doc.Pages[i]
		st.Pages = append(st.Pages, PageSummary{Name: p.Name, ImageURL: p.ImageURL})
	}
	if a.activePage == "" {
		return st, nil
	}
	p, err := editor.FindPage(&a.doc, a.activePage)
	if err != nil {
		return st, nil
	}
	st.Page = &PageDetail{
		Name:     p.Name,
		ImageURL: p.ImageURL,
		Style:    string(p.Style),
	}
	for i := range p.Texts {
		t := &p.Texts[i]
		st.Texts = append(st.Texts, TextDetail{
			Index: t.Index, X: t.X, Y: t.Y, Width: t.Width, TextSize: t.TextSize,
			Content: t.Content, ImageRef: t.ImageRef, Style: string(t.Style),
		})
	}
	for i := range p.Masks {
		m := &p.Masks[i]
		st.Masks = append(st.Masks, MaskDetail{
			Index: m.Index, X: m.X, Y: m.Y, Width: m.Width, Height: m.Height,
			ImageRef: m.ImageRef, Style: string(m.Style),
		})
	}
	if a.selectedIdx >= 0 {
		t, _, err := editor.FindTextByIndex(p, a.selectedIdx)
		if err == nil {
			st.Text = &TextDetail{
				Index: t.Index, X: t.X, Y: t.Y, Width: t.Width, TextSize: t.TextSize,
				Content: t.Content, ImageRef: t.ImageRef, Style: string(t.Style),
			}
		} else if m, _, err := editor.FindMaskByIndex(p, a.selectedIdx); err == nil {
			st.Mask = &MaskDetail{
				Index: m.Index, X: m.X, Y: m.Y, Width: m.Width, Height: m.Height,
				ImageRef: m.ImageRef, Style: string(m.Style),
			}
		}
	}
	return st, nil
}

// SetActivePage switches the current page.
func (a *App) SetActivePage(name string) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	if _, err := editor.FindPage(&a.doc, name); err != nil {
		return err
	}
	a.activePage = name
	a.selectedIdx = -1
	a.notify(nil)
	return nil
}

// BeginEdit starts a drag session (single undo snapshot).
func (a *App) BeginEdit() {
	a.mu.Lock()
	defer a.mu.Unlock()
	if !a.inEdit {
		a.snapshot()
		a.inEdit = true
	}
}

// EndEdit finishes a drag session.
func (a *App) EndEdit() {
	a.mu.Lock()
	a.inEdit = false
	a.mu.Unlock()
	a.notify(nil)
	a.maybeAutoCompile()
}

// SelectText selects a text index (-1 clears).
func (a *App) SelectText(index int) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.selectedIdx == index {
		return
	}
	a.selectedIdx = index
	a.notify(nil)
}

// SetAutoCompile toggles auto compile.
func (a *App) SetAutoCompile(on bool) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.autoCompile = on
}

// GetAssetDataURI returns a data URI for canvas display.
func (a *App) GetAssetDataURI(url string) (string, error) {
	a.mu.Lock()
	store := a.store
	client := a.client
	doc := a.doc
	filePath := a.filePath
	a.mu.Unlock()

	url = compileasset.ResolveAssetReference(url, doc.Consts)
	if !compileasset.IsAssetReference(url) {
		return "", fmt.Errorf("not an asset reference")
	}
	ctx := context.Background()
	label := "asset"
	if store != nil {
		labels := cache.URLPageLabels(doc)
		if l, ok := labels[url]; ok {
			label = l
		}
		if err := store.EnsureCached(ctx, client, url, label); err != nil {
			return "", err
		}
		if asset, ok, err := store.ReadAsset(url); err != nil {
			return "", err
		} else if ok {
			return compileasset.EncodeDataURI(asset.MIME, asset.Bytes), nil
		}
	}
	if compileasset.IsLocalAssetRef(url) {
		baseDir := ""
		if filePath != "" {
			baseDir = filepath.Dir(filePath)
		}
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
	fetched, err := compileasset.FetchURLs(client, []string{url})
	if err != nil {
		return "", err
	}
	asset := fetched[url]
	return compileasset.EncodeDataURI(asset.MIME, asset.Bytes), nil
}

// RefreshAssetURL forces online refresh for url.
func (a *App) RefreshAssetURL(url string) error {
	a.mu.Lock()
	store := a.store
	client := a.client
	doc := a.doc
	a.mu.Unlock()
	url = compileasset.ResolveAssetReference(url, doc.Consts)
	if store == nil {
		return fmt.Errorf("no cache store")
	}
	labels := cache.URLPageLabels(doc)
	label := labels[url]
	if label == "" {
		label = "asset"
	}
	if err := store.RefreshAsset(context.Background(), client, url, label); err != nil {
		return err
	}
	if a.emit != nil {
		a.emit("asset:refreshed", url)
	}
	a.notify(nil)
	return nil
}

// RefreshPageImage refreshes active page background.
func (a *App) RefreshPageImage() error {
	a.mu.Lock()
	page := a.activePage
	doc := a.doc
	a.mu.Unlock()
	p, err := editor.FindPage(&doc, page)
	if err != nil {
		return err
	}
	return a.RefreshAssetURL(p.ImageURL)
}

// Undo reverts last change.
func (a *App) Undo() error {
	a.mu.Lock()
	defer a.mu.Unlock()
	doc, ok := a.undo.pop()
	if !ok {
		return fmt.Errorf("nothing to undo")
	}
	a.redo.push(a.doc)
	a.doc = doc
	a.notify(nil)
	return nil
}

// Redo reapplies undone change.
func (a *App) Redo() error {
	a.mu.Lock()
	defer a.mu.Unlock()
	doc, ok := a.redo.pop()
	if !ok {
		return fmt.Errorf("nothing to redo")
	}
	a.undo.push(a.doc)
	a.doc = doc
	a.notify(nil)
	return nil
}

func encodePreview(data []byte, mime string) string {
	if mime == "" {
		mime = "image/svg+xml"
	}
	return "data:" + mime + ";base64," + base64.StdEncoding.EncodeToString(data)
}

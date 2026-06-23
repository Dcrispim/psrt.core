//go:build !js

package textoutline

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
)

// Outliner runs headless Chromium to convert laid-out text to SVG paths.
type Outliner struct {
	allocCtx context.Context
	cancel   context.CancelFunc
	mu       sync.Mutex
}

// NewOutliner starts a headless browser using the given Chromium executable.
func NewOutliner(parent context.Context, execPath string) (*Outliner, error) {
	if execPath == "" {
		return nil, fmt.Errorf("textoutline: empty chrome path")
	}
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-dev-shm-usage", true),
	)
	opts = append(opts, chromedp.ExecPath(execPath))
	allocCtx, cancel := chromedp.NewExecAllocator(parent, opts...)
	return &Outliner{allocCtx: allocCtx, cancel: cancel}, nil
}

// Close shuts down the browser allocator.
func (o *Outliner) Close() {
	if o == nil || o.cancel == nil {
		return
	}
	o.cancel()
}

// OutlinePage layouts texts in Chromium and returns vector paths per block.
func (o *Outliner) OutlinePage(ctx context.Context, in PageInput) ([]OutlinedBlock, error) {
	if o == nil {
		return nil, fmt.Errorf("textoutline: nil outliner")
	}
	html := BuildSnapshotHTML(in)
	htmlJSON, err := json.Marshal(html)
	if err != nil {
		return nil, err
	}
	loadJS := fmt.Sprintf(`(function(){document.open();document.write(%s);document.close();})()`, string(htmlJSON))

	o.mu.Lock()
	defer o.mu.Unlock()

	tabCtx, tabCancel := chromedp.NewContext(o.allocCtx)
	defer tabCancel()

	runCtx, runCancel := context.WithTimeout(tabCtx, 60*time.Second)
	defer runCancel()

	var blocks []OutlinedBlock
	script := opentypeJS + "\n" + extractJS

	err = chromedp.Run(runCtx,
		chromedp.EmulateViewport(int64(in.CanvasW), int64(in.CanvasH)),
		chromedp.Navigate("about:blank"),
		chromedp.Evaluate(loadJS, nil),
		chromedp.WaitReady("body", chromedp.ByQuery),
		chromedp.Evaluate(script, &blocks, func(p *runtime.EvaluateParams) *runtime.EvaluateParams {
			return p.WithAwaitPromise(true)
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("textoutline: chromedp: %w", err)
	}

	mergeBlockMeta(blocks, in)
	return blocks, nil
}

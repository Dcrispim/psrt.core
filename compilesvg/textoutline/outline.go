package textoutline

import (
	"context"
	"log"
	"os"
	"sync"
)

var (
	browserMu     sync.Mutex
	browserInst   *Outliner
	browserPath   string
	browserInitOK bool
)

// Outline layouts text blocks with Chromium when available, otherwise go-text.
func Outline(ctx context.Context, in PageInput) (Result, error) {
	if len(in.Blocks) == 0 {
		return Result{Blocks: nil}, nil
	}
	if path := ResolveChromeExec(); path != "" {
		o, err := getBrowser(ctx, path)
		if err != nil {
			logChromeDebug("browser init failed (%s): %v", path, err)
		} else {
			blocks, err := o.OutlinePage(ctx, in)
			if err != nil {
				logChromeDebug("outline failed (%s): %v", path, err)
			} else {
				mergeBlockMeta(blocks, in)
				return Result{Blocks: blocks, UsedGoTextFallback: false}, nil
			}
		}
	} else {
		logChromeDebug("no chrome executable found")
	}
	blocks, err := outlineGoText(in)
	if err != nil {
		return Result{}, err
	}
	mergeBlockMeta(blocks, in)
	return Result{Blocks: blocks, UsedGoTextFallback: true}, nil
}

func getBrowser(ctx context.Context, path string) (*Outliner, error) {
	browserMu.Lock()
	defer browserMu.Unlock()
	if browserInst != nil && browserPath == path && browserInitOK {
		return browserInst, nil
	}
	if browserInst != nil {
		browserInst.Close()
		browserInst = nil
	}
	o, err := NewOutliner(ctx, path)
	if err != nil {
		browserInitOK = false
		return nil, err
	}
	browserInst = o
	browserPath = path
	browserInitOK = true
	return browserInst, nil
}

func mergeBlockMeta(blocks []OutlinedBlock, in PageInput) {
	byIndex := make(map[int]*BlockInput, len(in.Blocks))
	for i := range in.Blocks {
		byIndex[in.Blocks[i].Index] = &in.Blocks[i]
	}
	for i := range blocks {
		inp, ok := byIndex[blocks[i].Index]
		if !ok {
			continue
		}
		if inp.FilterID != "" {
			blocks[i].FilterID = inp.FilterID
		}
		if inp.Transform != "" {
			blocks[i].Transform = inp.Transform
		}
		if blocks[i].PlainText == "" {
			blocks[i].PlainText = inp.PlainText
		}
	}
}

func logChromeDebug(format string, args ...any) {
	if os.Getenv("PSRT_DEBUG_CHROME") != "" {
		log.Printf("psrt: "+format, args...)
	}
}

// CloseBrowser releases the shared Chromium allocator.
func CloseBrowser() {
	browserMu.Lock()
	defer browserMu.Unlock()
	if browserInst != nil {
		browserInst.Close()
		browserInst = nil
	}
	browserPath = ""
	browserInitOK = false
}

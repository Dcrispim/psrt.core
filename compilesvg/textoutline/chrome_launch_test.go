//go:build chromiumtest

package textoutline

import (
	"context"
	"testing"
	"time"
)

func TestChromeLaunchFromResolve(t *testing.T) {
	path := ResolveChromeExec()
	if path == "" {
		t.Fatal("ResolveChromeExec returned empty")
	}
	t.Log("chrome:", path)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	o, err := NewOutliner(ctx, path)
	if err != nil {
		t.Fatalf("NewOutliner: %v", err)
	}
	defer o.Close()
}

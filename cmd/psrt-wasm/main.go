//go:build js

package main

import (
	"syscall/js"

	"github.com/Dcrispim/psrt.core/internal/wasmbridge"
)

func main() {
	exports := js.Global().Get("Object").New()
	wasmbridge.Register(exports)
	js.Global().Set("psrtWasm", exports)
	<-make(chan struct{})
}

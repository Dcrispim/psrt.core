// psrt-web-connector is the local HTTP bridge for psrt-gui-web (browser editor).
// Run: psrt-web-connector -config psrt-connector.ini
//
// Windows release build (no console on double-click):
//
//	go build -ldflags="-H windowsgui" -o psrt-web-connector.exe ./cmd/psrt-web-connector
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"runtime"

	"psrt/internal/webconnector"
)

func exitFatal(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	if !consoleAttached {
		showErrorBox(msg)
	}
	log.Fatal(msg)
}

func main() {
	initWinConsole()

	configPath := flag.String("config", "psrt-connector.ini", "path to connector INI config (required)")
	flag.Parse()

	absConfig, err := filepath.Abs(*configPath)
	if err != nil {
		exitFatal("config path: %v", err)
	}

	cfg, err := webconnector.LoadOrCreateConfig(absConfig)
	if err != nil {
		exitFatal("invalid config: %v", err)
	}

	audit := webconnector.NewAudit()
	srv := webconnector.NewServer(absConfig, cfg, audit)
	addr := fmt.Sprintf("127.0.0.1:%d", cfg.Port)

	audit.Startup(addr, cfg.BaseDir, cfg.AllowedOrigin)
	if consoleAttached {
		fmt.Printf("Código de conexão: %s (válido por 5 minutos)\n", srv.Auth().PairCodeForDisplay())
	}
	log.Printf("PSRT web connector listening on http://%s", addr)
	log.Printf("shared base_dir=%s allowed_origin=%s", cfg.BaseDir, cfg.AllowedOrigin)

	go func() {
		if err := http.ListenAndServe(addr, srv.Handler()); err != nil {
			exitFatal("listen: %v", err)
		}
	}()

	if runtime.GOOS == "windows" {
		webconnector.RunTray(srv, addr, consoleAttached)
		return
	}

	select {}
}

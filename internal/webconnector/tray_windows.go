//go:build windows

package webconnector

import (
	_ "embed"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/getlantern/systray"
)

//go:embed icons/tray.ico
var trayIcon []byte

// RunTray blocks on the Windows notification area until the user quits from the tray menu.
func RunTray(srv *Server, listenAddr string, consoleAttached bool) {
	systray.Run(func() { onTrayReady(srv, listenAddr, consoleAttached) }, onTrayExit)
}

func onTrayReady(srv *Server, listenAddr string, consoleAttached bool) {
	systray.SetIcon(trayIcon)
	systray.SetTitle("PSRT Connector")
	systray.SetTooltip(fmt.Sprintf("PSRT Web Connector — %s", listenAddr))

	mStatus := systray.AddMenuItem("PSRT Web Connector", "Servidor local ativo")
	mStatus.Disable()

	mCode := systray.AddMenuItem(codeMenuTitle(srv), "Código de pareamento (válido por 5 minutos)")
	mCode.Disable()

	mCopy := systray.AddMenuItem("Copiar código", "Copia o código de pareamento para a área de transferência")

	mRefresh := systray.AddMenuItem("Gerar novo código", "Invalida o código atual e gera outro")
	mReloadINI := systray.AddMenuItem("Recarregar INI", "Recarrega psrt-connector.ini do disco")
	systray.AddSeparator()
	mQuit := systray.AddMenuItem("Sair", "Encerra o conector")

	go refreshCodeMenuItem(srv, mCode)

	go func() {
		for {
			select {
			case <-mCopy.ClickedCh:
				code := srv.Auth().PairCodeForDisplay()
				if err := copyTextToClipboard(code); err != nil {
					log.Printf("copiar código: %v", err)
					mCopy.SetTitle("Falha ao copiar")
				} else {
					mCopy.SetTitle("Código copiado!")
				}
				go resetMenuLabel(mCopy, "Copiar código", 2*time.Second)
			case <-mRefresh.ClickedCh:
				code := srv.Auth().RefreshPairCode()
				mCode.SetTitle(codeMenuTitleWith(code))
				if consoleAttached {
					fmt.Printf("Código de conexão: %s (válido por 5 minutos)\n", code)
				}
			case <-mReloadINI.ClickedCh:
				old := srv.Config()
				if err := srv.ReloadConfig(); err != nil {
					log.Printf("recarregar INI: %v", err)
					continue
				}
				next := srv.Config()
				log.Printf("INI recarregado: base_dir=%s allowed_origin=%s port=%d",
					next.BaseDir, next.AllowedOrigin, next.Port)
				if old.Port != next.Port {
					log.Printf("porta alterada no INI (%d → %d); reinicie o conector para aplicar", old.Port, next.Port)
				}
			case <-mQuit.ClickedCh:
				systray.Quit()
				return
			}
		}
	}()
}

func onTrayExit() {
	os.Exit(0)
}

func codeMenuTitle(srv *Server) string {
	return codeMenuTitleWith(srv.Auth().PairCodeForDisplay())
}

func codeMenuTitleWith(code string) string {
	return fmt.Sprintf("Código: %s", code)
}

func resetMenuLabel(item *systray.MenuItem, label string, after time.Duration) {
	time.Sleep(after)
	item.SetTitle(label)
}

func refreshCodeMenuItem(srv *Server, item *systray.MenuItem) {
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		item.SetTitle(codeMenuTitle(srv))
	}
}

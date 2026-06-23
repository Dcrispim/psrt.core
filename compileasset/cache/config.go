package cache

import (
	"os"
	"path/filepath"
	"runtime"
)

// ConfigDir returns the PSRT config root (PSRT_CONFIG_DIR or OS default).
func ConfigDir() string {
	if d := os.Getenv("PSRT_CONFIG_DIR"); d != "" {
		return d
	}
	switch runtime.GOOS {
	case "windows":
		if d := os.Getenv("APPDATA"); d != "" {
			return filepath.Join(d, "psrt")
		}
		return filepath.Join(os.Getenv("USERPROFILE"), "AppData", "Roaming", "psrt")
	default:
		if d := os.Getenv("XDG_CONFIG_HOME"); d != "" {
			return filepath.Join(d, "psrt")
		}
		if home, err := os.UserHomeDir(); err == nil {
			return filepath.Join(home, ".config", "psrt")
		}
		return "psrt"
	}
}

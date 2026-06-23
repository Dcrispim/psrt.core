package crashlog

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime/debug"
	"sync"
	"time"
)

const logFileName = "psrt-gui.log"

var (
	mu      sync.Mutex
	logPath string
)

// Install resolves the log file path (next to the executable, or user config dir).
func Install() string {
	mu.Lock()
	defer mu.Unlock()
	if logPath != "" {
		return logPath
	}
	logPath = resolveLogPath()
	return logPath
}

// Path returns the active log file path, calling Install if needed.
func Path() string {
	mu.Lock()
	if logPath != "" {
		p := logPath
		mu.Unlock()
		return p
	}
	mu.Unlock()
	return Install()
}

func resolveLogPath() string {
	if exec, err := os.Executable(); err == nil {
		dir := filepath.Dir(exec)
		if dir != "" {
			return filepath.Join(dir, logFileName)
		}
	}
	if dir, err := os.UserConfigDir(); err == nil {
		return filepath.Join(dir, "psrt-gui", logFileName)
	}
	return logFileName
}

func appendRecord(operation, filePath, detail string, withStack bool) {
	mu.Lock()
	defer mu.Unlock()
	if logPath == "" {
		logPath = resolveLogPath()
	}
	_ = os.MkdirAll(filepath.Dir(logPath), 0o755)

	f, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return
	}
	defer f.Close()

	ts := time.Now().UTC().Format(time.RFC3339)
	_, _ = fmt.Fprintf(f, "\n---\n%s | operation=%s\n", ts, operation)
	if filePath != "" {
		_, _ = fmt.Fprintf(f, "file=%s\n", filePath)
	}
	_, _ = fmt.Fprintf(f, "error=%s\n", detail)
	if withStack {
		_, _ = f.Write(debug.Stack())
	}
}

// WritePanic records a recovered panic (never logs file content).
func WritePanic(operation, filePath string, recovered interface{}) {
	appendRecord(operation, filePath, fmt.Sprint(recovered), true)
}

// WriteError records a returned error without a stack trace.
func WriteError(operation, filePath string, err error) {
	if err == nil {
		return
	}
	appendRecord(operation, filePath, err.Error(), false)
}

// Guard recovers panics from Wails bindings and turns them into errors.
func Guard(operation string, filePathFn func() string, errOut *error) {
	if r := recover(); r != nil {
		fp := ""
		if filePathFn != nil {
			fp = filePathFn()
		}
		WritePanic(operation, fp, r)
		if errOut != nil {
			*errOut = fmt.Errorf("%s: %v (detalhes em %s)", operation, r, Path())
		}
	}
}

// Go runs fn in a goroutine with panic recovery.
func Go(operation, filePath string, fn func()) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				WritePanic(operation, filePath, r)
			}
		}()
		fn()
	}()
}

// RecoverMain should be deferred from main() to log fatal panics.
func RecoverMain() {
	if r := recover(); r != nil {
		WritePanic("main", "", r)
		os.Exit(1)
	}
}

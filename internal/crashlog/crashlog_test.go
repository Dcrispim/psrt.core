package crashlog

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWritePanicAndError(t *testing.T) {
	dir := t.TempDir()
	mu.Lock()
	old := logPath
	logPath = filepath.Join(dir, logFileName)
	mu.Unlock()
	t.Cleanup(func() {
		mu.Lock()
		logPath = old
		mu.Unlock()
	})

	WritePanic("TestOp", `D:\docs\test.psrt`, "simulated panic")
	WriteError("TestOp", `D:\docs\test.psrt`, os.ErrInvalid)

	data, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatal(err)
	}
	s := string(data)
	if !strings.Contains(s, "operation=TestOp") {
		t.Fatalf("missing operation: %s", s)
	}
	if !strings.Contains(s, `file=D:\docs\test.psrt`) {
		t.Fatalf("missing file path: %s", s)
	}
	if !strings.Contains(s, "simulated panic") {
		t.Fatalf("missing panic detail: %s", s)
	}
	if !strings.Contains(s, "invalid argument") {
		t.Fatalf("missing error detail: %s", s)
	}
}

package main

import "psrt/internal/crashlog"

func (g *GUIApp) currentFilePath() string {
	if g.app == nil {
		return ""
	}
	return g.app.FilePath()
}

func guardErr(g *GUIApp, operation string, err *error) {
	crashlog.Guard(operation, g.currentFilePath, err)
}

func guardVoid(g *GUIApp, operation string) {
	crashlog.Guard(operation, g.currentFilePath, nil)
}

func guardOpenFile(g *GUIApp, path *string, err *error) {
	crashlog.Guard("OpenFileDialog", func() string {
		if path != nil && *path != "" {
			return *path
		}
		return g.currentFilePath()
	}, err)
}

// ReportClientError logs frontend errors to psrt-gui.log (Wails binding).
func (g *GUIApp) ReportClientError(operation, filePath, message string) {
	fp := filePath
	if fp == "" {
		fp = g.currentFilePath()
	}
	crashlog.WriteError(operation, fp, fmtError(message))
}

func fmtError(msg string) error {
	return &clientError{msg: msg}
}

type clientError struct{ msg string }

func (e *clientError) Error() string { return e.msg }

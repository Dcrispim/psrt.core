//go:build !windows

package webconnector

// RunTray is a no-op on non-Windows platforms.
func RunTray(_ *Server, _ string, _ bool) {}

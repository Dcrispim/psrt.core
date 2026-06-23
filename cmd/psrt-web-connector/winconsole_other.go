//go:build !windows

package main

var consoleAttached = true

func initWinConsole() bool {
	return true
}

func showErrorBox(_ string) {}

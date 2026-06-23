//go:build windows

package main

import (
	"log"
	"os"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

const attachParentProcess = ^uint32(0)

var consoleAttached bool

func initWinConsole() bool {
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	getConsole := kernel32.NewProc("GetConsoleWindow")
	if hwnd, _, _ := getConsole.Call(); hwnd != 0 {
		consoleAttached = true
		return true
	}

	attach := kernel32.NewProc("AttachConsole")
	if ret, _, _ := attach.Call(uintptr(attachParentProcess)); ret != 0 {
		reopenConsoleHandles()
		consoleAttached = true
		return true
	}

	silenceConsoleOutput()
	consoleAttached = false
	return false
}

func reopenConsoleHandles() {
	out, err := os.OpenFile("CONOUT$", os.O_RDWR, 0)
	if err == nil {
		os.Stdout = out
		os.Stderr = out
		log.SetOutput(out)
	}
	in, err := os.OpenFile("CONIN$", os.O_RDWR, 0)
	if err == nil {
		os.Stdin = in
	}
}

func silenceConsoleOutput() {
	devNull, err := os.OpenFile(os.DevNull, os.O_WRONLY, 0)
	if err != nil {
		return
	}
	os.Stdout = devNull
	os.Stderr = devNull
	log.SetOutput(devNull)
}

func showErrorBox(message string) {
	title, _ := windows.UTF16PtrFromString("PSRT Web Connector")
	body, _ := windows.UTF16PtrFromString(message)
	_, _, _ = syscall.NewLazyDLL("user32.dll").
		NewProc("MessageBoxW").
		Call(0, uintptr(unsafe.Pointer(body)), uintptr(unsafe.Pointer(title)), uintptr(windows.MB_ICONERROR))
}

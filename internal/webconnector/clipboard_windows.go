//go:build windows

package webconnector

import (
	"fmt"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

const (
	cfUnicodeText = 13
	gmemMoveable  = 0x0002
)

var (
	clipboardUser32   = syscall.NewLazyDLL("user32.dll")
	clipboardKernel32 = syscall.NewLazyDLL("kernel32.dll")
	procOpenClipboard = clipboardUser32.NewProc("OpenClipboard")
	procCloseClipboard = clipboardUser32.NewProc("CloseClipboard")
	procEmptyClipboard = clipboardUser32.NewProc("EmptyClipboard")
	procSetClipboardData = clipboardUser32.NewProc("SetClipboardData")
	procGlobalAlloc = clipboardKernel32.NewProc("GlobalAlloc")
	procGlobalLock = clipboardKernel32.NewProc("GlobalLock")
	procGlobalUnlock = clipboardKernel32.NewProc("GlobalUnlock")
)

func copyTextToClipboard(text string) error {
	ret, _, err := procOpenClipboard.Call(0)
	if ret == 0 {
		return fmt.Errorf("open clipboard: %w", err)
	}
	defer procCloseClipboard.Call()

	procEmptyClipboard.Call()

	utf16, err := windows.UTF16FromString(text)
	if err != nil {
		return err
	}
	size := len(utf16) * 2

	hMem, _, err := procGlobalAlloc.Call(gmemMoveable, uintptr(size))
	if hMem == 0 {
		return fmt.Errorf("global alloc: %w", err)
	}

	pMem, _, err := procGlobalLock.Call(hMem)
	if pMem == 0 {
		return fmt.Errorf("global lock: %w", err)
	}

	dst := unsafe.Slice((*uint16)(unsafe.Pointer(pMem)), len(utf16))
	copy(dst, utf16)
	procGlobalUnlock.Call(hMem)

	setRet, _, err := procSetClipboardData.Call(cfUnicodeText, hMem)
	if setRet == 0 {
		return fmt.Errorf("set clipboard data: %w", err)
	}
	return nil
}

//go:build windows

package main

import (
	"os"
	"syscall"
	"unsafe"
)

// ── Windows Console API constants ───────────────────────────

const (
	// Input mode flags
	ENABLE_PROCESSED_INPUT        = 0x0001
	ENABLE_LINE_INPUT             = 0x0002
	ENABLE_ECHO_INPUT             = 0x0004
	ENABLE_WINDOW_INPUT           = 0x0008
	ENABLE_MOUSE_INPUT            = 0x0010
	ENABLE_VIRTUAL_TERMINAL_INPUT = 0x0200

	// Output mode flags
	ENABLE_PROCESSED_OUTPUT            = 0x0001
	ENABLE_WRAP_AT_EOL_OUTPUT          = 0x0002
	ENABLE_VIRTUAL_TERMINAL_PROCESSING = 0x0004 // enables ANSI escape codes
)

var (
	kernel32             = syscall.NewLazyDLL("kernel32.dll")
	procGetConsoleMode   = kernel32.NewProc("GetConsoleMode")
	procSetConsoleMode   = kernel32.NewProc("SetConsoleMode")
	procGetConsoleWindow = kernel32.NewProc("GetConsoleWindow")

	// COORD / CONSOLE_SCREEN_BUFFER_INFO for size detection
	procGetConsoleScreenBufferInfo = kernel32.NewProc("GetConsoleScreenBufferInfo")

	savedInputMode  uint32
	savedOutputMode uint32
)

type COORD struct{ X, Y int16 }
type SMALL_RECT struct{ Left, Top, Right, Bottom int16 }
type CONSOLE_SCREEN_BUFFER_INFO struct {
	Size              COORD
	CursorPosition    COORD
	Attributes        uint16
	Window            SMALL_RECT
	MaximumWindowSize COORD
}

func getConsoleMode(handle syscall.Handle, mode *uint32) bool {
	r, _, _ := procGetConsoleMode.Call(uintptr(handle), uintptr(unsafe.Pointer(mode)))
	return r != 0
}

func setConsoleMode(handle syscall.Handle, mode uint32) {
	procSetConsoleMode.Call(uintptr(handle), uintptr(mode))
}

func enableRaw() {
	inHandle, _ := syscall.GetStdHandle(syscall.STD_INPUT_HANDLE)
	outHandle, _ := syscall.GetStdHandle(syscall.STD_OUTPUT_HANDLE)

	// Save original modes so we can restore on exit
	getConsoleMode(inHandle, &savedInputMode)
	getConsoleMode(outHandle, &savedOutputMode)

	// Raw input: disable line buffering, echo, processed input
	// Keep ENABLE_VIRTUAL_TERMINAL_INPUT for arrow key escape sequences
	rawInput := uint32(ENABLE_VIRTUAL_TERMINAL_INPUT)
	setConsoleMode(inHandle, rawInput)

	// Enable ANSI/VT100 output processing (required on Windows 10+)
	rawOutput := savedOutputMode | ENABLE_VIRTUAL_TERMINAL_PROCESSING | ENABLE_PROCESSED_OUTPUT
	setConsoleMode(outHandle, rawOutput)
}

func restoreTerminal() {
	inHandle, _ := syscall.GetStdHandle(syscall.STD_INPUT_HANDLE)
	outHandle, _ := syscall.GetStdHandle(syscall.STD_OUTPUT_HANDLE)
	setConsoleMode(inHandle, savedInputMode)
	setConsoleMode(outHandle, savedOutputMode)
}

// termSize returns (cols, rows) using GetConsoleScreenBufferInfo.
func termSize() (int, int) {
	outHandle, err := syscall.GetStdHandle(syscall.STD_OUTPUT_HANDLE)
	if err != nil {
		return 0, 0
	}
	var info CONSOLE_SCREEN_BUFFER_INFO
	r, _, _ := procGetConsoleScreenBufferInfo.Call(
		uintptr(outHandle),
		uintptr(unsafe.Pointer(&info)),
	)
	if r == 0 {
		return 0, 0
	}
	cols := int(info.Window.Right-info.Window.Left) + 1
	rows := int(info.Window.Bottom-info.Window.Top) + 1
	return cols, rows
}

// openInputDevice returns os.Stdin on Windows (no /dev/tty).
// With ENABLE_VIRTUAL_TERMINAL_INPUT set, arrow keys arrive as ESC sequences.
func openInputDevice() *os.File {
	return os.Stdin
}

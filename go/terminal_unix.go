//go:build !windows

package main

import (
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// ── Raw mode via stty (works on macOS + Linux) ──────────────

func stty(args ...string) {
	cmd := exec.Command("stty", args...)
	cmd.Stdin = os.Stdin // stty reads terminal settings via stdin fd
	cmd.Run()
}

func enableRaw() {
	// raw    → disable canonical mode (no line buffering) + signals
	// -echo  → don't print typed chars
	// min 0  → read() returns immediately even with 0 bytes
	// time 1 → 0.1s read timeout (so the goroutine isn't blocked forever)
	stty("raw", "-echo", "min", "0", "time", "1")
}

func restoreTerminal() {
	stty("sane") // restore all terminal settings to sane defaults
}

// termSize returns (cols, rows) via "stty size".
// Output format: "rows cols\n"
func termSize() (int, int) {
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	out, err := cmd.Output()
	if err != nil {
		return 0, 0
	}
	parts := strings.Fields(strings.TrimSpace(string(out)))
	if len(parts) != 2 {
		return 0, 0
	}
	rows, _ := strconv.Atoi(parts[0])
	cols, _ := strconv.Atoi(parts[1])
	return cols, rows
}

// openInputDevice returns /dev/tty for reading keyboard input.
// Falls back to os.Stdin if /dev/tty is unavailable.
func openInputDevice() *os.File {
	f, err := os.OpenFile("/dev/tty", os.O_RDONLY, 0)
	if err != nil {
		return os.Stdin
	}
	return f
}

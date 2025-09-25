package platform

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"golang.org/x/term"
)

var originalTerminalState *term.State

func IsCapsLockOn() bool {
	files, err := filepath.Glob("/sys/class/leds/input*::capslock/brightness")
	if err != nil || len(files) == 0 {
		return false
	}
	content, err := os.ReadFile(files[0])
	if err != nil {
		return false
	}
	return strings.TrimSpace(string(content)) == "1"
}

func IsKoreanInputMode(lastChar byte) bool {
	return lastChar >= 0x80
}

func RestoreTerminal() {
	if originalTerminalState != nil {
		term.Restore(int(syscall.Stdin), originalTerminalState)
		originalTerminalState = nil
	}
}

func SetupSignalHandler() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)
	go func() {
		<-c
		fmt.Println("\nProgram interrupted. Restoring terminal...")
		RestoreTerminal()
		os.Exit(0)
	}()
}

func GetTerminalState() error {
	var err error
	originalTerminalState, err = term.GetState(int(syscall.Stdin))
	return err
}

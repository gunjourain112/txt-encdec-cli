package ui

import (
	"fmt"
	"os"
	"strings"
	"syscall"
	"txt-encdec-cli/platform"

	"golang.org/x/term"
)

func ReadPasswordWithStars(prompt string) (string, error) {
	fmt.Print(prompt)

	err := platform.GetTerminalState()
	if err != nil {
		return "", fmt.Errorf("failed to get terminal state: %v", err)
	}

	platform.SetupSignalHandler()

	rawState, err := term.MakeRaw(int(syscall.Stdin))
	if err != nil {
		return "", fmt.Errorf("failed to make terminal raw: %v", err)
	}

	defer func() {
		term.Restore(int(syscall.Stdin), rawState)
		platform.RestoreTerminal()
	}()

	var password []byte
	var b [1]byte
	var lastInputChar byte = 0

	for {
		n, err := os.Stdin.Read(b[:])
		if err != nil {
			return "", fmt.Errorf("failed to read input: %v", err)
		}
		if n == 0 {
			continue
		}

		char := b[0]

		if char == 13 || char == 10 {
			break
		}

		if char == 127 || char == 8 {
			if len(password) > 0 {
				password = password[:len(password)-1]
				lastInputChar = 0
				redrawPasswordLineWithMode(prompt, password, lastInputChar)
			}
			continue
		}

		if char == 3 {
			fmt.Println("\nProgram interrupted. Restoring terminal...")
			platform.RestoreTerminal()
			os.Exit(0)
		}

		if char >= 32 && char <= 126 {
			password = append(password, char)
			lastInputChar = char
			redrawPasswordLineWithMode(prompt, password, lastInputChar)
		} else if char >= 0x80 {
			password = append(password, char)
			lastInputChar = char
			redrawPasswordLineWithMode(prompt, password, lastInputChar)
		}
	}

	fmt.Println()
	return string(password), nil
}

func redrawPasswordLineWithMode(prompt string, password []byte, lastChar byte) {
	capsOn := platform.IsCapsLockOn()
	koreanOn := platform.IsKoreanInputMode(lastChar)

	fmt.Print("\r\033[K")
	fmt.Print(prompt)

	var statusParts []string
	if koreanOn {
		statusParts = append(statusParts, "\033[31m[한글]\033[0m")
	}
	if capsOn {
		statusParts = append(statusParts, "\033[31m[CAPS]\033[0m")
	}

	if len(statusParts) > 0 {
		fmt.Print(" " + strings.Join(statusParts, " ") + " ")
	} else {
		fmt.Print(" ")
	}

	for i := 0; i < len(password); i++ {
		fmt.Print("*")
	}
}

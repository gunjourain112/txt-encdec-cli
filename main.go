package main

import (
	"bufio"
	"fmt"
	"os"
	"txt-encdec-cli/core"
	"txt-encdec-cli/platform"
	"txt-encdec-cli/tui"
	"txt-encdec-cli/ui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	defer platform.RestoreTerminal()

	m := tui.New()
	p := tea.NewProgram(m)

	finalModel, err := p.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "TUI error: %v", err)
		os.Exit(1)
	}

	tuiModel := finalModel.(tui.Model)
	chosenMode, selected := tuiModel.GetSelectedMode()

	if !selected {
		return
	}

	processCLI(chosenMode)
}

func processCLI(mode string) {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Printf("\nSelected mode: %s\n", mode)

	if platform.IsCapsLockOn() {
		fmt.Println("\033[31m⚠️  WARNING: CAPS LOCK is ON\033[0m")
	}

	secretKey, err := ui.ReadPasswordWithStars("Enter Secret Key: ")
	if err != nil {
		fmt.Printf("Error reading secret key: %v\n", err)
		os.Exit(1)
	}

	var promptText string
	if mode == "Encrypt" {
		promptText = "Enter text to encrypt: "
	} else {
		promptText = "Enter text to decrypt: "
	}

	fmt.Print(promptText)
	if !scanner.Scan() {
		fmt.Println("Error reading input text")
		os.Exit(1)
	}
	inputText := scanner.Text()

	var result string
	if mode == "Encrypt" {
		result, err = core.Encrypt(secretKey, inputText)
		if err != nil {
			fmt.Printf("Encryption error: %v\n", err)
			os.Exit(1)
		}
	} else {
		result, err = core.Decrypt(secretKey, inputText)
		if err != nil {
			fmt.Printf("Decryption error: %v\n", err)
			os.Exit(1)
		}
	}

	err = platform.CopyToClipboard(result)
	if err != nil {
		fmt.Printf("Clipboard error: %v\n", err)
		fmt.Println("Result:", result)
		os.Exit(1)
	}

	fmt.Println("Success! Result copied to clipboard.")
}

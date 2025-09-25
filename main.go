package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"txt-encdec-cli/core"
	"txt-encdec-cli/platform"
	"txt-encdec-cli/ui"
)

func main() {
	defer platform.RestoreTerminal()

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("=== Text Encryption Tool ===")
	fmt.Println("1. Encrypt")
	fmt.Println("2. Decrypt")
	fmt.Print("Mode: ")

	if !scanner.Scan() {
		fmt.Println("Error reading input")
		os.Exit(1)
	}

	modeStr := strings.TrimSpace(scanner.Text())
	mode, err := strconv.Atoi(modeStr)
	if err != nil || (mode != 1 && mode != 2) {
		fmt.Println("Error: Invalid mode. Please enter 1 or 2.")
		os.Exit(1)
	}

	if platform.IsCapsLockOn() {
		fmt.Println("\033[31mWARNING: CAPS LOCK is ON\033[0m")
	}

	secretKey, err := ui.ReadPasswordWithStars("Enter Secret Key: ")
	if err != nil {
		fmt.Printf("Error reading secret key: %v\n", err)
		os.Exit(1)
	}

	var promptText string
	if mode == 1 {
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
	if mode == 1 {
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

package main

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"

	"golang.org/x/term"
)

var originalTerminalState *term.State

func deriveKey(secret string) []byte {
	h := sha256.Sum256([]byte(secret))
	return h[:]
}

func encrypt(secret, plaintext string) (string, error) {
	key := deriveKey(secret)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func decrypt(secret, encoded string) (string, error) {
	key := deriveKey(secret)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", err
	}
	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", fmt.Errorf("ciphertext too short")
	}
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}

func copyToClipboard(text string) error {
	tools := []struct {
		name string
		args []string
	}{
		{"wl-copy", nil},
		{"xclip", []string{"-selection", "clipboard"}},
		{"xsel", []string{"--clipboard", "--input"}},
	}

	for _, tool := range tools {
		cmd := exec.Command(tool.name, tool.args...)
		stdin, err := cmd.StdinPipe()
		if err != nil {
			continue
		}
		if err := cmd.Start(); err != nil {
			continue
		}
		if _, err := stdin.Write([]byte(text)); err != nil {
			stdin.Close()
			continue
		}
		stdin.Close()
		if err := cmd.Wait(); err == nil {
			return nil
		}
	}
	return fmt.Errorf("no clipboard tool available")
}

func isCapsLockOn() bool {
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

// 터미널 상태 복원
func restoreTerminal() {
	if originalTerminalState != nil {
		term.Restore(int(syscall.Stdin), originalTerminalState)
		originalTerminalState = nil
	}
}

// 시그널 핸들러
func setupSignalHandler() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)
	go func() {
		<-c
		fmt.Println("\nProgram interrupted. Restoring terminal...")
		restoreTerminal()
		os.Exit(0)
	}()
}

func readPasswordWithStars(prompt string) (string, error) {
	fmt.Print(prompt)

	// 원본 터미널 상태 저장
	var err error
	originalTerminalState, err = term.GetState(int(syscall.Stdin))
	if err != nil {
		return "", fmt.Errorf("failed to get terminal state: %v", err)
	}

	setupSignalHandler()

	// Raw 모드로 전환
	rawState, err := term.MakeRaw(int(syscall.Stdin))
	if err != nil {
		return "", fmt.Errorf("failed to make terminal raw: %v", err)
	}

	defer func() {
		term.Restore(int(syscall.Stdin), rawState)
		restoreTerminal()
	}()

	var password []byte
	var b [1]byte

	for {
		n, err := os.Stdin.Read(b[:])
		if err != nil {
			return "", fmt.Errorf("failed to read input: %v", err)
		}
		if n == 0 {
			continue
		}

		char := b[0]

		// 엔터 (13 or 10)
		if char == 13 || char == 10 {
			break
		}

		// 백스페이스 (127 or 8)
		if char == 127 || char == 8 {
			if len(password) > 0 {
				password = password[:len(password)-1]
				redrawPasswordLine(prompt, password) // 라인 다시 그림
			}
			continue
		}

		// Ctrl+C (3)
		if char == 3 {
			fmt.Println("\nProgram interrupted. Restoring terminal...")
			restoreTerminal()
			os.Exit(0)
		}

		// 일반
		if char >= 32 && char <= 126 {
			password = append(password, char)
			redrawPasswordLine(prompt, password) // 라인 다시 그림
		}
	}

	fmt.Println()
	return string(password), nil
}

func redrawPasswordLine(prompt string, password []byte) {
	capsOn := isCapsLockOn()

	fmt.Print("\r\033[K")
	fmt.Print(prompt)

	if capsOn {
		fmt.Print(" \033[31m[CAPS]\033[0m ")
	} else {
		fmt.Print(" ")
	}

	for i := 0; i < len(password); i++ {
		fmt.Print("*")
	}
}

func main() {
	defer restoreTerminal()

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

	if isCapsLockOn() {
		fmt.Println("\033[31mWARNING: CAPS LOCK is ON\033[0m")
	}

	secretKey, err := readPasswordWithStars("Enter Secret Key: ")
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
		result, err = encrypt(secretKey, inputText)
		if err != nil {
			fmt.Printf("Encryption error: %v\n", err)
			os.Exit(1)
		}
	} else {
		result, err = decrypt(secretKey, inputText)
		if err != nil {
			fmt.Printf("Decryption error: %v\n", err)
			os.Exit(1)
		}
	}

	err = copyToClipboard(result)
	if err != nil {
		fmt.Printf("Clipboard error: %v\n", err)
		fmt.Println("Result:", result)
		os.Exit(1)
	}

	fmt.Println("Success! Result copied to clipboard.")
}

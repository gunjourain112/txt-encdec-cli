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
	"strconv"
	"strings"
	"syscall"

	"golang.org/x/term"
)

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

func readPasswordWithStars(prompt string) (string, error) {
	fmt.Print(prompt)

	oldState, err := term.MakeRaw(int(syscall.Stdin))
	if err != nil {
		return "", err
	}
	defer term.Restore(int(syscall.Stdin), oldState)

	var password []byte
	var b [1]byte

	for {
		n, err := os.Stdin.Read(b[:])
		if err != nil {
			return "", err
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
				fmt.Print("\b \b")
			}
			continue
		}

		// Ctrl+C (3)
		if char == 3 {
			fmt.Println()
			os.Exit(1)
		}

		// 일반
		if char >= 32 && char <= 126 {
			password = append(password, char)
			fmt.Print("*")
		}
	}

	fmt.Println()
	return string(password), nil
}

func main() {
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

	secretKey, err := readPasswordWithStars("Enter Secret Key: ")
	if err != nil {
		fmt.Println("Error reading secret key")
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

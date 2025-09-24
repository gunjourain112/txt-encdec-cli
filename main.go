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

func writeToFile(content string) error {
	file, err := os.Create("out.txt")
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(content)
	return err
}

func readPassword(prompt string) (string, error) {
	fmt.Print(prompt)
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", err
	}
	fmt.Println()
	return string(bytePassword), nil
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

	// 마스킹 처리된 비밀키 입력
	secretKey, err := readPassword("Enter Secret Key: ")
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

	err = writeToFile(result)
	if err != nil {
		fmt.Printf("File write error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Success! Result saved to out.txt")
}

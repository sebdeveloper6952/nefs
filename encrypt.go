package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"math"
	"os"

	"github.com/nbd-wtf/go-nostr/nip44"
)

const (
	MaxPlaintextSize = int(0xffff)
)

func readFileToBase64(filePath string) (string, error) {
	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("readFileToBase64: %w", err)
	}

	return base64.StdEncoding.EncodeToString(fileBytes), nil
}

func encrypt(privkey string, pubkey string, filePath string) error {
	fileBase64, err := readFileToBase64(filePath)
	if err != nil {
		return fmt.Errorf("encrypt: %w", err)
	}

	base64Parts := splitString(fileBase64, MaxPlaintextSize)

	convKey, err := nip44.GenerateConversationKey(pubkey, privkey)
	if err != nil {
		return fmt.Errorf("encrypt: GenerateConversationKey: %w", err)
	}

	salt := make([]byte, 32)
	if _, err := rand.Read(salt); err != nil {
		return fmt.Errorf("encrypt: make salt: %w", err)
	}

	for i, part := range base64Parts {
		base64Encrypted, err := nip44.Encrypt(part, convKey, nip44.WithCustomSalt(salt))
		if err != nil {
			return fmt.Errorf("encrypt: Encrypt: %w", err)
		}

		file, err := os.Create(fmt.Sprintf("encrypted_%d", i))
		if err != nil {
			return fmt.Errorf("encrypt: create file: %w", err)
		}
		defer file.Close()

		_, err = file.WriteString(base64Encrypted)
		if err != nil {
			return fmt.Errorf("encrypt: write to file: %w", err)
		}
	}

	return nil
}

func splitString(s string, MaxPlaintextSize int) []string {
	need := int(math.Ceil(float64(len(s)) / float64(MaxPlaintextSize)))
	parts := make([]string, need)

	for i := 0; i < need; i++ {
		start := i * MaxPlaintextSize
		end := start + MaxPlaintextSize
		if len(s) < end {
			end = len(s)
		}

		parts[i] = s[i*MaxPlaintextSize : end]
	}

	return parts
}

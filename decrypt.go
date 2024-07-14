package main

import (
	"encoding/base64"
	"fmt"
	"os"

	"github.com/nbd-wtf/go-nostr/nip44"
)

func decrypt(sk string, pk string, filePath string) error {
	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("decrypt: read file: %w", err)
	}

	convKey, err := nip44.GenerateConversationKey(pk, sk)
	if err != nil {
		return fmt.Errorf("decrypt: compute conversation key: %w", err)
	}

	plaintextBase64, err := nip44.Decrypt(string(fileBytes), convKey)
	if err != nil {
		return fmt.Errorf("decrypt: %w", err)
	}

	plaintextBytes, err := base64.StdEncoding.DecodeString(plaintextBase64)
	if err != nil {
		return fmt.Errorf("decrypt: decode base64: %w", err)
	}

	decryptedFile, err := os.Create("decrypted.png")
	if err != nil {
		return fmt.Errorf("decrypt: create: %w", err)
	}

	_, err = decryptedFile.Write(plaintextBytes)
	if err != nil {
		return fmt.Errorf("decrypt: write to file: %w", err)
	}

	return nil
}

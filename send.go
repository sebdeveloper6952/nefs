package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"

	"github.com/nbd-wtf/go-nostr"
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

func send(sk string, pubkey string, filePath string) (string, error) {
	var (
		pk, _ = nostr.GetPublicKey(sk)
	)

	fileBase64, err := readFileToBase64(filePath)
	if err != nil {
		return "", fmt.Errorf("encrypt: %w", err)
	}

	base64Parts := splitString(fileBase64, MaxPlaintextSize)

	convKey, err := nip44.GenerateConversationKey(pubkey, sk)
	if err != nil {
		return "", fmt.Errorf("encrypt: GenerateConversationKey: %w", err)
	}

	salt := make([]byte, 32)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("encrypt: make salt: %w", err)
	}

	events := make([]nostr.Event, len(base64Parts))
	for i := len(base64Parts) - 1; i >= 0; i-- {
		base64Encrypted, err := nip44.Encrypt(base64Parts[i], convKey, nip44.WithCustomSalt(salt))
		if err != nil {
			return "", fmt.Errorf("encrypt: Encrypt: %w", err)
		}

		event := nostr.Event{
			Kind:      69999,
			PubKey:    pk,
			Content:   base64Encrypted,
			CreatedAt: nostr.Now(),
		}

		if i < len(base64Parts)-1 {
			nextEventID := events[i+1].ID
			event.Tags = nostr.Tags{
				nostr.Tag{"e", nextEventID},
			}
		}

		event.Sign(sk)
		events[i] = event
	}

	return events[0].ID, nil
}

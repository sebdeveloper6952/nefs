package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"

	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/nip44"
	blossomClient "github.com/sebdeveloper6952/blossom-server/client"
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
		pk, _       = nostr.GetPublicKey(sk)
		client, err = blossomClient.New([]string{"http://localhost:8000"}, sk)
	)

	fileBase64, err := readFileToBase64(filePath)
	if err != nil {
		return "", fmt.Errorf("encrypt: %w", err)
	}

	base64Parts := splitString(fileBase64, nip44.MaxPlaintextSize)

	convKey, err := nip44.GenerateConversationKey(pubkey, sk)
	if err != nil {
		return "", fmt.Errorf("encrypt: GenerateConversationKey: %w", err)
	}

	salt := make([]byte, 32)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("encrypt: make salt: %w", err)
	}

	var (
		event = nostr.Event{
			Kind:      70000,
			PubKey:    pk,
			CreatedAt: nostr.Now(),
			Tags:      make([]nostr.Tag, len(base64Parts)),
		}
	)

	for i := len(base64Parts) - 1; i >= 0; i-- {
		base64Encrypted, err := nip44.Encrypt(base64Parts[i], convKey, nip44.WithCustomSalt(salt))
		if err != nil {
			return "", fmt.Errorf("encrypt: Encrypt: %w", err)
		}

		blob, err := client.Upload([]byte(base64Encrypted))
		if err != nil {
			return "", fmt.Errorf("encrypt: upload: %w", err)
		}
		event.Tags[i] = nostr.Tag{"chunk", blob.Sha256, fmt.Sprintf("%d", i)}
	}

	event.Sign(sk)
	publishEvents([]nostr.Event{event})

	return event.ID, nil
}

package main

import (
	"crypto/rand"
	"fmt"

	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/nip44"
	blossomClient "github.com/sebdeveloper6952/blossom-server/client"
)

type sendResult struct {
	Chunks  int
	EventID string
}

func send(sk string, pubkey string, filePath string) (*sendResult, error) {
	pk, err := nostr.GetPublicKey(sk)
	if err != nil {
		return nil, fmt.Errorf("send: invalid private key: %w\n", err)
	}

	cdnList, err := fetchPubkeyCDNList(pk)
	if err != nil || len(cdnList) == 0 {
		return nil, fmt.Errorf("send: fetch cdn list: %w", err)
	}

	fileBase64, err := readFileToBase64(filePath)
	if err != nil {
		return nil, fmt.Errorf("send: %w", err)
	}
	base64Parts := splitString(fileBase64, nip44.MaxPlaintextSize)

	convKey, err := nip44.GenerateConversationKey(pubkey, sk)
	if err != nil {
		return nil, fmt.Errorf("send: GenerateConversationKey: %w", err)
	}

	salt := make([]byte, 32)
	if _, err := rand.Read(salt); err != nil {
		return nil, fmt.Errorf("send: make salt: %w", err)
	}

	event := nostr.Event{
		Kind:      70000,
		PubKey:    pk,
		CreatedAt: nostr.Now(),
		Tags:      make([]nostr.Tag, len(base64Parts)),
	}
	blossomClient, err := blossomClient.New(cdnList, sk)
	if err != nil {
		return nil, fmt.Errorf("send: init blossom client: %w\n", err)
	}

	// for each file chunk, encrypt and upload chunk to specified blossom server
	// TODO: support multiple servers
	for i := len(base64Parts) - 1; i >= 0; i-- {
		base64Encrypted, err := nip44.Encrypt(base64Parts[i], convKey, nip44.WithCustomSalt(salt))
		if err != nil {
			return nil, fmt.Errorf("send: Encrypt: %w", err)
		}

		blob, err := blossomClient.Upload([]byte(base64Encrypted))
		if err != nil {
			return nil, fmt.Errorf("send: upload: %w", err)
		}
		event.Tags[i] = nostr.Tag{"chunk", blob.Sha256, fmt.Sprintf("%d", i)}
	}

	if err := event.Sign(sk); err != nil {
		return nil, fmt.Errorf("send: sign chunks event: %w\n", err)
	}

	if err := publishEvents([]nostr.Event{event}); err != nil {
		return nil, fmt.Errorf("send: publish chunk event: %w\n", err)
	}

	return &sendResult{EventID: event.ID, Chunks: len(base64Parts)}, nil
}

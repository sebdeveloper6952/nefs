package nefs

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"github.com/nbd-wtf/go-nostr/nip44"
	blossomClient "github.com/sebdeveloper6952/blossom-server/client"
)

type ReceiveResult struct {
	Chunks    int
	FileBytes []byte
}

func Receive(sk string, chunksEventID string) (*ReceiveResult, error) {
	chunksEvent, err := FetchEventByID(chunksEventID)
	if err != nil {
		return nil, fmt.Errorf("receive: fetch summary chunksEvent: %w", err)
	}

	convKey, err := nip44.GenerateConversationKey(chunksEvent.PubKey, sk)
	if err != nil {
		return nil, fmt.Errorf("receive: compute conversation key: %w", err)
	}

	serverUrlTag := chunksEvent.Tags.GetFirst([]string{"server"})
	if serverUrlTag == nil || len(*serverUrlTag) != 2 {
		return nil, errors.New("receive: event must have at least one 'server' tag")
	}

	blossomClient, _ := blossomClient.New((*serverUrlTag)[1], sk)
	chunkTags := chunksEvent.Tags.GetAll([]string{"chunk"})
	decryptedBase64 := make([]string, len(chunkTags))
	chunkNumber := 0

	for _, chunk := range chunkTags {
		if len(chunk) < 2 {
			return nil, fmt.Errorf("receive: malformed chunk tag\n")
		}

		blobBytes, err := blossomClient.Get(chunk[1])
		if err != nil {
			fmt.Println(err)
			continue
		}

		plaintextBase64, err := nip44.Decrypt(string(blobBytes), convKey)
		if err != nil {
			return nil, fmt.Errorf("decrypt: %w", err)
		}
		decryptedBase64 = append(decryptedBase64, plaintextBase64)
		chunkNumber++
	}

	fileBytes, err := base64.StdEncoding.DecodeString(strings.Join(decryptedBase64, ""))
	if err != nil {
		return nil, fmt.Errorf("decrypt: decode file base64: %w", err)
	}

	return &ReceiveResult{
		FileBytes: fileBytes,
		Chunks:    chunkNumber,
	}, err
}

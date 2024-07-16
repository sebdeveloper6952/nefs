package main

import (
	"encoding/base64"
	"fmt"
	"os"
	"strings"

	"github.com/nbd-wtf/go-nostr/nip44"
	blossomClient "github.com/sebdeveloper6952/blossom-server/client"
)

func receive(sk string, pk string, chunksEventID string) error {
	convKey, err := nip44.GenerateConversationKey(pk, sk)
	if err != nil {
		return fmt.Errorf("receive: compute conversation key: %w", err)
	}

	chunksEvent, err := fetchEventByID(eventID)
	if err != nil {
		return fmt.Errorf("receive: fetch summary chunksEvent: %w", err)
	}

	cdnList, err := fetchPubkeyCDNList(chunksEvent.PubKey)
	if err != nil || len(cdnList) == 0 {
		return fmt.Errorf("receive: fetch cdn list: %w", err)
	}

	blossomClient, _ := blossomClient.New(cdnList, sk)
	chunkTags := chunksEvent.Tags.GetAll([]string{"chunk"})
	decryptedBase64 := make([]string, len(chunkTags))
	chunkNumber := 0

	for _, chunk := range chunkTags {
		if len(chunk) < 2 {
			return fmt.Errorf("receive: malformed chunk tag\n")
		}

		blobBytes, err := blossomClient.Get(chunk[1])
		if err != nil {
			fmt.Println(err)
			continue
		}

		plaintextBase64, err := nip44.Decrypt(string(blobBytes), convKey)
		if err != nil {
			return fmt.Errorf("decrypt: %w", err)
		}
		decryptedBase64 = append(decryptedBase64, plaintextBase64)
		chunkNumber++
	}

	fileBytes, err := base64.StdEncoding.DecodeString(strings.Join(decryptedBase64, ""))
	if err != nil {
		return fmt.Errorf("decrypt: decode file base64: %w", err)
	}
	file, _ := os.Create("decrypted.png")
	defer file.Close()
	_, err = file.Write(fileBytes)

	return err
}

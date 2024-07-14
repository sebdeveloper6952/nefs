package main

import (
	"encoding/base64"
	"fmt"
	"os"
	"strings"

	"github.com/nbd-wtf/go-nostr/nip44"
)

func receive(sk string, pk string, eventID string) error {
	convKey, err := nip44.GenerateConversationKey(pk, sk)
	if err != nil {
		return fmt.Errorf("receive: compute conversation key: %w", err)
	}

	// fetch start event
	currentEvent, err := fetchEvent(eventID, true)
	if err != nil {
		return fmt.Errorf("receive: fetch start event: %w", err)
	}

	// to incrementally append the decrypted base64 content
	decryptedBase64 := make([]string, 1)
	partNumber := 0

	plaintextBase64, err := nip44.Decrypt(currentEvent.Content, convKey)
	if err != nil {
		return fmt.Errorf("decrypt: %w", err)
	}
	decryptedBase64 = append(decryptedBase64, plaintextBase64)

	// loop over next parts until done
	for {
		// needs more validation, for now, if there is no next event we assume this was the last part of the file.
		nextEventID, ok := extractNextEventID(currentEvent)
		if !ok {
			break
		}
		partNumber++

		currentEvent, err = fetchEvent(nextEventID, false)
		if err != nil {
			return fmt.Errorf("receive: fetch event: %w", err)
		}

		plaintextBase64, err := nip44.Decrypt(currentEvent.Content, convKey)
		if err != nil {
			return fmt.Errorf("decrypt part number %d: %w", partNumber, err)
		}

		decryptedBase64 = append(decryptedBase64, plaintextBase64)
	}

	plaintextBytes, err := base64.StdEncoding.DecodeString(strings.Join(decryptedBase64, ""))
	if err != nil {
		return fmt.Errorf("decrypt: decode base64: %w", err)
	}

	file, _ := os.Create("decrypted.png")
	defer file.Close()
	_, err = file.WriteString(string(plaintextBytes))
	if err != nil {
		return fmt.Errorf("decrypt: write to file: %w", err)
	}

	return nil
}

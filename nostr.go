package main

import (
	"context"
	"fmt"

	"github.com/nbd-wtf/go-nostr"
)

func publishEvents(events []nostr.Event) error {
	ctx := context.Background()
	relay, err := nostr.RelayConnect(ctx, "ws://localhost:4869")
	if err != nil {
		fmt.Println(err)
	}

	for _, ev := range events {
		if err := relay.Publish(ctx, ev); err != nil {
			fmt.Println(err)
		}
	}

	return nil
}

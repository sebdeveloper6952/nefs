package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/nbd-wtf/go-nostr"
)

func publishEvents(events []nostr.Event) error {
	ctx := context.Background()
	relay, err := nostr.RelayConnect(ctx, "ws://localhost:4869")
	if err != nil {
		return err
	}

	for _, ev := range events {
		if err := relay.Publish(ctx, ev); err != nil {
			return err
		}
		fmt.Printf("%s\n", ev)
	}

	return nil
}

func fetchEvent(id string, isStart bool) (*nostr.Event, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := nostr.Filter{
		IDs: []string{id},
	}
	if isStart {
		filter.Tags = nostr.TagMap{
			"t": []string{"start"},
		}
	}

	relay, err := nostr.RelayConnect(ctx, "ws://localhost:4869")
	if err != nil {
		fmt.Println(err)
	}

	sub, err := relay.Subscribe(ctx, nostr.Filters{filter})
	if err != nil {
		return nil, fmt.Errorf("fetchEventByID: %w\n", err)
	}

	select {
	case ev := <-sub.Events:
		return ev, nil
	case <-ctx.Done():
		return nil, errors.New("fetchEventByID: timeout\n")
	}
}

func extractNextEventID(event *nostr.Event) (string, bool) {
	tag := event.Tags.GetFirst([]string{"e"})
	if tag == nil || len(*tag) < 2 {
		return "", false
	}

	return (*tag)[1], true
}

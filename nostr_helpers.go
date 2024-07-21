package nefs

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/nbd-wtf/go-nostr"
)

func PublishEvents(events []nostr.Event, relays []string) error {
	ctx := context.Background()

	for _, url := range relays {
		relay, err := nostr.RelayConnect(ctx, url)
		if err != nil {
			return err
		}

		for _, ev := range events {
			if err := relay.Publish(ctx, ev); err != nil {
				return err
			}
		}
	}

	return nil
}

func PublishCDNList(sk string, cdns []string) error {
	ctx := context.Background()
	relay, err := nostr.RelayConnect(ctx, "ws://localhost:4869")
	if err != nil {
		return err
	}

	pk, err := nostr.GetPublicKey(sk)
	event := nostr.Event{
		Kind:   10063,
		PubKey: pk,
		Tags:   make([]nostr.Tag, len(cdns)),
	}
	for i := range cdns {
		event.Tags[i] = nostr.Tag{"server", cdns[i]}
	}
	event.Sign(sk)

	return relay.Publish(ctx, event)
}

func FetchEvent(filters nostr.Filters) (*nostr.Event, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	relay, err := nostr.RelayConnect(ctx, "ws://localhost:4869")
	if err != nil {
		fmt.Println(err)
	}

	sub, err := relay.Subscribe(ctx, filters)
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

func FetchEventByID(id string) (*nostr.Event, error) {
	return FetchEvent(nostr.Filters{
		nostr.Filter{
			IDs: []string{id},
		},
	})
}

func FetchPubkeyCDNList(pk string) ([]string, error) {
	filter := nostr.Filter{
		Authors: []string{pk},
		Kinds:   []int{10063},
	}

	event, err := FetchEvent(nostr.Filters{filter})
	if err != nil {
		return nil, err
	}

	serverTags := event.Tags.GetAll([]string{"server"})
	serverUrls := make([]string, len(serverTags))
	for i := range serverTags {
		if len(serverTags[i]) == 2 {
			serverUrls[i] = serverTags[i][1]
		}
	}

	return serverUrls, nil
}

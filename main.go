package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

var (
	sk       string
	pk       string
	filePath string
	eventID  string
)

func encryptAction() cli.ActionFunc {
	return func(ctx *cli.Context) error {
		eventID, err := send(sk, pk, filePath)
		if err != nil {
			return err
		}

		fmt.Printf("share this event ID with the recipient: %s\n", eventID)

		return nil
	}
}

func decryptAction() cli.ActionFunc {
	return func(ctx *cli.Context) error {
		return receive(sk, pk, eventID)
	}
}

func main() {
	app := &cli.App{
		Name:        "nefs",
		Description: "send/receive encrypted files over nostr",
		Flags:       []cli.Flag{},
		Commands: []*cli.Command{
			{
				Name:    "encrypt",
				Aliases: []string{"e"},
				Action:  encryptAction(),
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "file",
						Aliases:     []string{"f"},
						Usage:       "file to encrypt",
						Required:    true,
						Destination: &filePath,
					},
					&cli.StringFlag{
						Name:        "sk",
						Usage:       "private key to encrypt/decrypt",
						Required:    true,
						Destination: &sk,
					},
					&cli.StringFlag{
						Name:        "pk",
						Usage:       "public key to encrypt/decrypt",
						Required:    true,
						Destination: &pk,
					},
				},
			},
			{
				Name:    "decrypt",
				Aliases: []string{"d"},
				Action:  decryptAction(),
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "event",
						Aliases:     []string{"e"},
						Usage:       "event ID of start",
						Required:    true,
						Destination: &eventID,
					},
					&cli.StringFlag{
						Name:        "sk",
						Usage:       "private key to encrypt/decrypt",
						Required:    true,
						Destination: &sk,
					},
					&cli.StringFlag{
						Name:        "pk",
						Usage:       "public key to encrypt/decrypt",
						Required:    true,
						Destination: &pk,
					},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

var (
	sk         string
	pk         string
	filePath   string
	eventID    string
	serverUrls cli.StringSlice
)

func sendCmdAction() cli.ActionFunc {
	return func(ctx *cli.Context) error {
		res, err := send(sk, pk, filePath)
		if err != nil {
			return err
		}

		fmt.Printf("uploaded %d chunks\nshare this event ID with the recipient: %s\n", res.Chunks, res.EventID)

		return nil
	}
}

func receiveCmdAction() cli.ActionFunc {
	return func(ctx *cli.Context) error {
		res, err := receive(sk, pk, eventID)

		fmt.Printf("received %d chunks\n", res.Chunks)

		return err
	}
}

func serverListAction() cli.ActionFunc {
	return func(ctx *cli.Context) error {
		return publishCDNList(sk, serverUrls.Value())
	}
}

func main() {
	app := &cli.App{
		Name:        "nefs",
		Description: "send/receive encrypted files over nostr",
		Flags:       []cli.Flag{},
		Commands: []*cli.Command{
			{
				Name:    "send",
				Aliases: []string{"s"},
				Action:  sendCmdAction(),
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
				Name:    "receive",
				Aliases: []string{"r"},
				Action:  receiveCmdAction(),
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
			{
				Name:        "serverlist",
				Description: "publish server list event (10063)",
				Action:      serverListAction(),
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "sk",
						Usage:       "private key to encrypt/decrypt",
						Required:    true,
						Destination: &sk,
					},
					&cli.StringSliceFlag{
						Name:        "servers",
						Aliases:     []string{"s"},
						Usage:       "server urls",
						Required:    true,
						Destination: &serverUrls,
					},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

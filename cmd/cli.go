package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gabriel-vasile/mimetype"
	"github.com/sebdeveloper6952/nefs"
	"github.com/urfave/cli/v2"
)

var (
	sk         string
	pk         string
	filePath   string
	outputFile string
	relayUrl   string
	serverUrl  string
	eventID    string
)

func sendCmdAction() cli.ActionFunc {
	return func(ctx *cli.Context) error {
		res, err := nefs.Send(
			sk,
			pk,
			filePath,
			relayUrl,
			serverUrl,
		)
		if err != nil {
			return err
		}

		fmt.Printf("uploaded %d chunks\nshare this event ID with the recipient: %s\n", res.Chunks, res.EventID)

		return nil
	}
}

func receiveCmdAction() cli.ActionFunc {
	return func(ctx *cli.Context) error {
		res, err := nefs.Receive(sk, eventID)
		if err != nil {
			return err
		}

		mimeType := mimetype.Detect(res.FileBytes)
		fileName := outputFile
		if mimeType.Extension() != "" {
			fileName += mimeType.Extension()
		}
		file, err := os.Create(fileName)
		if err != nil {
			return err
		}

		_, err = file.Write(res.FileBytes)
		if err != nil {
			return err
		}

		fmt.Printf("received %d chunks\n", res.Chunks)

		return err
	}
}

func main() {
	app := &cli.App{
		Name:  "nefs",
		Usage: "send/receive encrypted files over nostr",
		Flags: []cli.Flag{},
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
						Name:        "privkey",
						Aliases:     []string{"sk"},
						Usage:       "private key used to sign event",
						Required:    true,
						Destination: &sk,
					},
					&cli.StringFlag{
						Name:        "pubkey",
						Aliases:     []string{"pk"},
						Usage:       "recipient of file",
						Required:    true,
						Destination: &pk,
					},
					&cli.StringFlag{
						Name:        "relay",
						Aliases:     []string{"r"},
						Usage:       "relay where file chunks event will be published",
						Required:    true,
						Destination: &relayUrl,
					},
					&cli.StringFlag{
						Name:        "server",
						Aliases:     []string{"s"},
						Usage:       "URL of blossom server where file chunks will be uploaded",
						Required:    true,
						Destination: &serverUrl,
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
						Name:        "privkey",
						Aliases:     []string{"sk"},
						Usage:       "private key to decrypt file chunks",
						Required:    true,
						Destination: &sk,
					},
					&cli.StringFlag{
						Name:        "relay",
						Aliases:     []string{"r"},
						Usage:       "relay where event was published",
						Required:    true,
						Destination: &relayUrl,
					},
					&cli.StringFlag{
						Name:        "output",
						Aliases:     []string{"o"},
						Usage:       "name of output file",
						Required:    true,
						Destination: &outputFile,
					},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

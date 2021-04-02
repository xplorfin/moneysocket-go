package cmd

import (
	"fmt"

	cli "github.com/urfave/cli/v2"
	messagebase "github.com/xplorfin/moneysocket-go/moneysocket/message/base"
	"github.com/xplorfin/moneysocket-go/relay"
	"github.com/xplorfin/moneysocket-go/terminus"
)

// Start takes a list of args (from os.Args()), parse a command and runs it.
func Start(args []string) {
	// TODO go requires flags to appear before positional args
	// this should be fixed for parity w/ moneysocket-py
	app := cli.NewApp()
	app.Name = "terminus cli"
	app.Version = messagebase.Version

	app.Commands = []*cli.Command{
		{
			Name:  "getinfo",
			Usage: "gets list of beacons",
			Flags: configOptions,
			Action: func(context *cli.Context) error {
				cfg, err := fetchConfig(context)
				if err != nil {
					return err
				}
				client := terminus.NewClient(&cfg)

				res, err := client.GetInfo()
				if err != nil {
					return err
				}
				fmt.Println(res)
				return nil
			},
		},
		{
			Name:      "listen",
			Usage:     "listen to a websocket",
			Flags:     configOptions,
			ArgsUsage: "--config=x [account name]",
			Action: func(context *cli.Context) error {
				cfg, err := fetchConfig(context)
				if err != nil {
					return err
				}

				client := terminus.NewClient(&cfg)

				// TODO replace this with https://git.io/Jt1LR
				accountName := context.Args().Get(0)
				if accountName == "" {
					return fmt.Errorf("acocunt name is required")
				}

				res, err := client.Listen(accountName)
				if err != nil {
					return err
				}
				fmt.Println(res)
				return nil
			},
		},
		{
			Name:  "start",
			Usage: "start the server",
			Flags: configOptions,
			Action: func(context *cli.Context) error {
				cfg, err := fetchConfig(context)
				if err != nil {
					return err
				}

				server, err := terminus.NewTerminus(&cfg)
				if err != nil {
					return err
				}
				return server.Start(context.Context)
			},
		},
		{
			Name:  "relay",
			Usage: "start the relay server",
			Flags: configOptions,
			Action: func(context *cli.Context) error {
				cfg, err := fetchConfig(context)
				if err != nil {
					return err
				}

				server := relay.NewRelay(&cfg)

				return server.RunApp()
			},
		},
	}
	if err := app.Run(args); err != nil {
		panic(err)
	}
}

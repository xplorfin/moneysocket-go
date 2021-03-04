package cmd

import (
	"fmt"
	cli "github.com/urfave/cli/v2"
	messagebase "github.com/xplorfin/moneysocket-go/moneysocket/message/base"
	"github.com/xplorfin/moneysocket-go/terminus"
)

// TODO go requires flags to appear before positional args
// this should be fixed for parity
func Start(args []string) {
	app := cli.NewApp()
	app.Name = "terminus cli"
	app.Version = messagebase.VERSION

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

				server := terminus.NewTerminus(&cfg)
				return server.Start(context.Context)
			},
		},
	}
	if err := app.Run(args); err != nil {
		panic(err)
	}
}

package cmd

import (
	"github.com/urfave/cli/v2"
	"github.com/xplorfin/moneysocket-go/moneysocket/config"
)

var configOptions = []cli.Flag{
	&cli.StringFlag{
		Name:    "config",
		Aliases: []string{"c"},
		Usage:   "set a config path",
	},
}

func fetchConfigFromPath(path string) (config.Config, error) {
	cfg, err := config.ParseConfigFromFile(path)
	if err != nil {
		panic(err)
	}
	err = cfg.Validate()
	return cfg, err
}

// global options must come before command options (https://git.io/Jt1Ij)
// this is a workaround
func fetchConfig(c *cli.Context) (config.Config, error) {
	return fetchConfigFromPath(c.String("config"))
}

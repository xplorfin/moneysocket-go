package cmd

import (
	"github.com/urfave/cli/v2"
	"github.com/xplorfin/moneysocket-go/moneysocket/config"
)

// configOptions defines the flag used across commands to define a terminus config file
var configOptions = []cli.Flag{
	&cli.StringFlag{
		Name:    "config",
		Aliases: []string{"c"},
		Usage:   "set a config path",
	},
}

// fetchConfigFromPath fetches a config.Config from a path,
//	validates it and returns an error if invalid otherwise it
// will return the parsed config
func fetchConfigFromPath(path string) (config.Config, error) {
	cfg, err := config.ParseConfigFromFile(path)
	if err != nil {
		panic(err)
	}
	err = cfg.Validate()
	return cfg, err
}

// fetchConfig gets the parsed config from a common "config" flag
// this is a workaround to deal with a lack of global
// options (see: https://git.io/Jt1Ij)
func fetchConfig(c *cli.Context) (config.Config, error) {
	return fetchConfigFromPath(c.String("config"))
}

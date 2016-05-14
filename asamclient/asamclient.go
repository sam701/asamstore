package main

import (
	"log"
	"os"

	"github.com/codegangsta/cli"
	"github.com/sam701/asamstore/asamclient/initialize"
	"github.com/sam701/asamstore/asamclient/put"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	app := cli.NewApp()
	app.Name = "asamclient"
	app.Version = "0.1.0"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config, c",
			Usage: "path to config.toml",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:   "init",
			Usage:  "initalize the asamstore",
			Action: initialize.Initialize,
		},
		{
			Name:      "put",
			Usage:     "put content into storage",
			ArgsUsage: "<path to content>",
			Action:    put.PutAction,
		},
	}
	app.Run(os.Args)
}

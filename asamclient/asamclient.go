package main

import (
	"log"
	"os"

	"github.com/codegangsta/cli"
	"github.com/sam701/asamstore/asamclient/initialize"
	"github.com/sam701/asamstore/asamclient/mount"
	"github.com/sam701/asamstore/asamclient/put"
	"github.com/sam701/asamstore/asamclient/root"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	app := cli.NewApp()
	app.Name = "asamclient"
	app.Version = "0.2.0"
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
			Name:  "put",
			Usage: "put content into storage",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "root",
					Usage: "root node",
				},
			},
			ArgsUsage: "<path to content>",
			Action:    put.PutAction,
		},
		{
			Name:      "root",
			Usage:     "create new root node",
			ArgsUsage: "<root node name>",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "new",
					Usage: "create new root node",
				},
			},
			Action: root.Root,
		},
		{
			Name:      "mount",
			Usage:     "mount storage",
			ArgsUsage: "<mount point>",
			Action:    mount.Mount,
		},
	}
	app.Run(os.Args)
}

package main

import (
	"log"
	"os"

	"github.com/urfave/cli"
	"github.com/sam701/asamstore/asamclient/get"
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
			Name:  "init",
			Usage: "initalize the asamstore",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "client",
					Usage: "initialize client",
				},
				cli.BoolFlag{
					Name:  "blob-server",
					Usage: "initialize blob server",
				},
				cli.StringFlag{
					Name:  "dest-dir, d",
					Usage: "destination directory for configuration and certificates",
				},
			},
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
			Name:      "get",
			Usage:     "print ref content",
			ArgsUsage: "<ref>",
			Action:    get.Get,
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
				cli.BoolFlag{
					Name:  "list, l",
					Usage: "list roots",
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

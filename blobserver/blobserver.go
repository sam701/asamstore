package main

import (
	"fmt"
	"log"
	"os"

	"github.com/codegangsta/cli"
)

var store *DataStore
var config *configuration

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	app := cli.NewApp()
	app.Version = "0.1.0"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config, c",
			Usage: "path to config.yaml",
		},
	}
	app.Action = startServer
	app.Run(os.Args)

}

func startServer(c *cli.Context) error {
	cfgPath := c.String("config")
	if cfgPath == "" {
		fmt.Println("No config file provided")
		cli.ShowAppHelp(c)
		return nil
	}
	config = readConfig(cfgPath)
	store = OpenDataStore(config.StorageDir)

	initTlsClient()
	go syncWithAllRemotes()

	setupHttpHandlers()
	startHttpsServer()
	return nil
}

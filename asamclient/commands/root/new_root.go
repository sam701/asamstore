package root

import (
	"log"

	"github.com/sam701/asamstore/asamclient/client"
	"github.com/sam701/asamstore/asamclient/config"
	"github.com/sam701/asamstore/asamclient/index"
	"github.com/sam701/asamstore/asamclient/schema"
	"github.com/urfave/cli"
)

func Root(c *cli.Context) error {
	newRootName := c.String("new")
	if c.Bool("list") {
		listRoots(c)
		return nil
	}
	if newRootName == "" {
		cli.ShowCommandHelp(c, "root")
		return nil
	}

	s := schema.NewSchema(schema.ContentTypeRoot)
	s.RootName = newRootName

	conf := config.ReadConfig(c.GlobalString("config"))
	cl := client.NewClient(conf)
	ref := cl.PutSchema(s)

	ix := index.OpenIndex(conf.IndexDir)
	ix.AddRoot(newRootName, ref)

	log.Println("Root ref:", ref)

	cl.UpdateServerState()

	return nil
}

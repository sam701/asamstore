package root

import (
	"log"

	"github.com/codegangsta/cli"
	"github.com/sam701/asamstore/asamclient/client"
	"github.com/sam701/asamstore/asamclient/config"
	"github.com/sam701/asamstore/asamclient/index"
	"github.com/sam701/asamstore/asamclient/schema"
)

func Root(c *cli.Context) error {
	newRootName := c.String("new")
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

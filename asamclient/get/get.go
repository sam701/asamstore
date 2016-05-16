package get

import (
	"log"
	"os"

	"github.com/codegangsta/cli"
	"github.com/sam701/asamstore/asamclient/client"
	"github.com/sam701/asamstore/asamclient/config"
	"github.com/sam701/asamstore/asamclient/schema"
)

func Get(c *cli.Context) error {
	ref := c.Args().First()
	if ref == "" {
		log.Fatalln("No ref provided")
	}
	cfg := config.ReadConfig(c.GlobalString("config"))
	bsClient := client.NewClient(cfg)

	ok := bsClient.Get(schema.BlobRef(ref), os.Stdout)

	if !ok {
		log.Println("No content exists with such ref")
	}
	return nil
}

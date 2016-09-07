package get

import (
	"log"
	"os"

	"github.com/urfave/cli"
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

	data := bsClient.Get(schema.BlobRef(ref))

	if data == nil {
		log.Println("No content exists with such ref")
	} else {
		os.Stdout.Write(data)
	}
	return nil
}
